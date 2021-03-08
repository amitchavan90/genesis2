package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
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

// DistributorFactory creates distributors
func DistributorFactory() *db.Distributor {
	companyName := faker.Company().Name()
	u := &db.Distributor{
		ID:   uuid.Must(uuid.NewV4()).String(),
		Name: companyName,
		Code: companyName[0:3],
	}
	return u
}

// NewDistributorStore handle distributor methods
func NewDistributorStore(conn *sqlx.DB) *Distributor {
	ts := &Distributor{conn}
	return ts
}

// Distributor for persistence
type Distributor struct {
	Conn *sqlx.DB
}

// Get distributor
func (s *Distributor) Get(id uuid.UUID) (*db.Distributor, error) {
	return db.FindDistributor(s.Conn, id.String())
}

// GetMany distributors
func (s *Distributor) GetMany(keys []string) (db.DistributorSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Distributors(db.DistributorWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Distributor{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Distributor{}
	for _, key := range keys {
		var row *db.Distributor
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

// GetByCode gets a Distributor by their code
func (s *Distributor) GetByCode(code string) (*db.Distributor, error) {
	return db.Distributors(db.DistributorWhere.Code.EQ(code)).One(s.Conn)
}

// All distributors
func (s *Distributor) All() (db.DistributorSlice, error) {
	return db.Distributors().All(s.Conn)
}

// SearchSelect searchs/selects Distributors
func (s *Distributor) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Distributor, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("LOWER(%s) LIKE ?",
						db.DistributorColumns.Name,
					),
					"%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.DistributorWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.DistributorWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.Distributors(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.DistributorColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.DistributorColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.DistributorColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.DistributorColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Distributors(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// Insert distributor
func (s *Distributor) Insert(record *db.Distributor, txes ...*sql.Tx) (*db.Distributor, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update distributor
func (s *Distributor) Update(record *db.Distributor, txes ...*sql.Tx) (*db.Distributor, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive distributors
func (s *Distributor) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Distributor, error) {
	u, err := db.FindDistributor(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.DistributorColumns.Archived, db.DistributorColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive distributors
func (s *Distributor) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Distributor, error) {
	u, err := db.FindDistributor(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.DistributorColumns.Archived, db.DistributorColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of distributors
func (s *Distributor) Count() (int64, error) {
	dat, err := db.Distributors().Count(s.Conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return dat, nil
}
