package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"syreclabs.com/go/faker"
)

// WalletTransactionFactory creates wallet transactions
func WalletTransactionFactory() *db.WalletTransaction {
	t := &db.WalletTransaction{
		ID:            uuid.Must(uuid.NewV4()).String(),
		LoyaltyPoints: faker.Number().NumberInt(3),
	}
	return t
}

// NewWalletTransactionStore returns a new wallet transaction repo that implements WalletTransactionMutator, WalletTransactionArchiver and WalletTransactionQueryer
func NewWalletTransactionStore(conn *sqlx.DB) *WalletTransactions {
	r := &WalletTransactions{conn}
	return r
}

// WalletTransactions for persistence
type WalletTransactions struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *WalletTransactions) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All wallet transactions
func (s *WalletTransactions) All(txes ...*sql.Tx) (db.WalletTransactionSlice, error) {
	return db.WalletTransactions().All(s.Conn)
}

// Count gives the amount of wallet transactions
func (s *WalletTransactions) Count() (int64, error) {
	return db.WalletTransactions().Count(s.Conn)
}

// SearchSelect searchs/selects wallet transactions
func (s *WalletTransactions) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.WalletTransaction, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?)",
						db.WalletTransactionColumns.ID,
					),
					"%"+searchText+"%",
				))
		}
	}

	// Get Count
	count, err := db.WalletTransactions(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.WalletTransactionColumns.CreatedAt+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.WalletTransactionColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.WalletTransactions(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany wallet transactions given a list of IDs
func (s *WalletTransactions) GetMany(keys []string, txes ...*sql.Tx) (db.WalletTransactionSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.WalletTransactions(db.WalletTransactionWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.WalletTransaction{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.WalletTransaction{}
	for _, key := range keys {
		var row *db.WalletTransaction
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

// Get a wallet transaction given their ID
func (s *WalletTransactions) Get(id uuid.UUID, txes ...*sql.Tx) (*db.WalletTransaction, error) {
	dat, err := db.WalletTransactions(db.WalletTransactionWhere.ID.EQ(id.String())).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a wallet transaction
func (s *WalletTransactions) Insert(t *db.WalletTransaction, txes ...*sql.Tx) (*db.WalletTransaction, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return t.Insert(tx, boil.Infer())
	}, txes...)

	err = t.Reload(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return t, nil
}

// Create a wallet transaction
func (s *WalletTransactions) Create(input *db.WalletTransaction, txes ...*sql.Tx) (*db.WalletTransaction, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}
