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

// ContractFactory creates contracts
func ContractFactory() *db.Contract {
	companyName := faker.Company().Name()
	u := &db.Contract{
		ID:           uuid.Must(uuid.NewV4()).String(),
		Name:         fmt.Sprintf("%s Contract", companyName),
		SupplierName: companyName,
		DateSigned:   null.TimeFrom(time.Now()),
	}
	return u
}

// NewContractStore handle contract methods
func NewContractStore(conn *sqlx.DB) *Contract {
	ts := &Contract{conn}
	return ts
}

// Contract for persistence
type Contract struct {
	Conn *sqlx.DB
}

// Get contract
func (s *Contract) Get(id uuid.UUID) (*db.Contract, error) {
	return db.FindContract(s.Conn, id.String())
}

// GetMany contracts
func (s *Contract) GetMany(keys []string) (db.ContractSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Contracts(db.ContractWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Contract{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Contract{}
	for _, key := range keys {
		var row *db.Contract
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

// GetByCode gets a Contract by their code
func (s *Contract) GetByCode(code string) (*db.Contract, error) {
	return db.Contracts(db.ContractWhere.Code.EQ(code)).One(s.Conn)
}

// All contracts
func (s *Contract) All() (db.ContractSlice, error) {
	return db.Contracts().All(s.Conn)
}

// SearchSelect searchs/selects Contracts
func (s *Contract) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Contract, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.ContractColumns.Name,
						db.ContractColumns.SupplierName,
					),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.ContractWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.ContractWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.Contracts(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.ContractColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.ContractColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.ContractColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.ContractColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Contracts(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// Insert contract
func (s *Contract) Insert(record *db.Contract, txes ...*sql.Tx) (*db.Contract, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update contract
func (s *Contract) Update(record *db.Contract, txes ...*sql.Tx) (*db.Contract, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive contracts
func (s *Contract) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Contract, error) {
	u, err := db.FindContract(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.ContractColumns.Archived, db.ContractColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive contracts
func (s *Contract) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Contract, error) {
	u, err := db.FindContract(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.ContractColumns.Archived, db.ContractColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of contracts
func (s *Contract) Count() (int64, error) {
	dat, err := db.Contracts().Count(s.Conn)
	if err != nil {
		return dat, terror.New(err, "")
	}
	return dat, nil
}
