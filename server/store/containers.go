package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"syreclabs.com/go/faker"
)

// ContainerFactory creates containers
func ContainerFactory() *db.Container {
	u := &db.Container{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Code:        faker.Numerify("CON#####"),
		Description: faker.Lorem().Sentence(10),
	}
	return u
}

// NewContainerStore handle container methods
func NewContainerStore(conn *sqlx.DB) *Container {
	ts := &Container{conn}
	return ts
}

// Container for persistence
type Container struct {
	Conn *sqlx.DB
}

// GetByCode gets a Container by their code
func (s *Container) GetByCode(code string) (*db.Container, error) {
	dat, err := db.Containers(db.ContainerWhere.Code.EQ(code)).One(s.Conn)
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}

// Get containers
func (s *Container) Get(id uuid.UUID) (*db.Container, error) {
	dat, err := db.FindContainer(s.Conn, id.String())
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}

// All containers
func (s *Container) All() (db.ContainerSlice, error) {
	dat, err := db.Containers().All(s.Conn)
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}

// SearchSelect searchs/selects Containers
func (s *Container) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Container, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.ContainerColumns.Code,
						db.ContainerColumns.Description),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.ContainerWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.ContainerWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.Containers(queries...).Count(s.Conn)
	if err != nil {
		return 0, nil, terror.New(err, "")
	}

	// Sort
	sortDir := " ASC"
	if search.SortDir != nil && *search.SortDir == graphql.SortDirDescending {
		sortDir = " DESC"
	}

	if search.SortBy != nil {
		switch *search.SortBy {
		case graphql.SortByOptionDateCreated:
			queries = append(queries, qm.OrderBy(db.ContainerColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.ContainerColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.ContainerColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.ContainerColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Containers(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany containers
func (s *Container) GetMany(keys []string) (db.ContainerSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Containers(db.ContainerWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Container{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Container{}
	for _, key := range keys {
		var row *db.Container
		for _, record := range records {
			if record.ID == key {
				row = record
				break
			}
		}
		result = append(result, row)
	}
	return result, nil
}

// GetRange returns all containers between two container codes
func (s *Container) GetRange(from string, to string) (db.ContainerSlice, error) {
	start, err := strconv.Atoi(strings.Replace(from, "CON", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}
	end, err := strconv.Atoi(strings.Replace(to, "CON", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}

	if start < 0 || end < start {
		return nil, terror.New(fmt.Errorf("invalid arguments"), "")
	}

	keys := []string{}
	for i := start; i <= end; i++ {
		keys = append(keys, fmt.Sprintf("CON%05d", i))
	}

	records, err := db.Containers(db.ContainerWhere.Code.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Container{}, nil
	}
	if err != nil {
		return nil, terror.New(err, "")
	}

	result := []*db.Container{}
	for _, key := range keys {
		var row *db.Container
		for _, record := range records {
			if record.Code == key {
				row = record
				break
			}
		}
		result = append(result, row)
	}
	return result, nil
}

// Insert containers
func (s *Container) Insert(record *db.Container, txes ...*sql.Tx) (*db.Container, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update containers
func (s *Container) Update(record *db.Container, txes ...*sql.Tx) (*db.Container, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive containers
func (s *Container) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Container, error) {
	u, err := db.FindContainer(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.ContainerColumns.Archived, db.ContainerColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive containers
func (s *Container) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Container, error) {
	u, err := db.FindContainer(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.ContainerColumns.Archived, db.ContainerColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of containers
func (s *Container) Count() (int64, error) {
	dat, err := db.Containers().Count(s.Conn)
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}

// PalletCount returns the count of pallets in the container
func (s *Container) PalletCount(record *db.Container) (int64, error) {
	dat, err := db.Pallets(db.PalletWhere.ContainerID.EQ(null.StringFrom(record.ID))).Count(s.Conn)
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}
