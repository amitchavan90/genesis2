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

// UserPurchaseActivityFactory creates userPurchaseActivities
func UserPurchaseActivityFactory() *db.UserPurchaseActivity {
	t := &db.UserPurchaseActivity{
		ID:            uuid.Must(uuid.NewV4()).String(),
		LoyaltyPoints: faker.Number().NumberInt(3),
	}
	return t
}

// NewUserPurchaseActivityStore returns a new user purchase activity repo that implements UserPurchaseActivityMutator, UserPurchaseActivityArchiver and UserPurchaseActivityQueryer
func NewUserPurchaseActivityStore(conn *sqlx.DB) *UserPurchaseActivities {
	r := &UserPurchaseActivities{conn}
	return r
}

// UserPurchaseActivities for persistence
type UserPurchaseActivities struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *UserPurchaseActivities) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All userPurchaseActivities
func (s *UserPurchaseActivities) All(txes ...*sql.Tx) (db.UserPurchaseActivitySlice, error) {
	return db.UserPurchaseActivities().All(s.Conn)
}

// Count gives the amount of userPurchaseActivities
func (s *UserPurchaseActivities) Count() (int64, error) {
	return db.UserPurchaseActivities().Count(s.Conn)
}

// SearchSelect searchs/selects userPurchaseActivities
func (s *UserPurchaseActivities) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.UserPurchaseActivity, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?)",
						db.UserPurchaseActivityColumns.ID,
					),
					"%"+searchText+"%",
				))
		}
	}

	// Get Count
	count, err := db.UserPurchaseActivities(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.UserPurchaseActivityColumns.CreatedAt+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.UserPurchaseActivityColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.UserPurchaseActivities(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetProduct by productID
func (s *UserPurchaseActivities) GetProduct(productID string, txes ...*sql.Tx) (*db.Product, error) {
	dat, err := db.Products(db.ProductWhere.ID.EQ(productID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany userPurchaseActivities given a list of IDs
func (s *UserPurchaseActivities) GetMany(keys []string, txes ...*sql.Tx) (db.UserPurchaseActivitySlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.UserPurchaseActivities(db.UserPurchaseActivityWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.UserPurchaseActivity{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.UserPurchaseActivity{}
	for _, key := range keys {
		var row *db.UserPurchaseActivity
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

// Get by id
func (s *UserPurchaseActivities) Get(id uuid.UUID, txes ...*sql.Tx) (*db.UserPurchaseActivity, error) {
	dat, err := db.UserPurchaseActivities(db.UserPurchaseActivityWhere.ID.EQ(id.String())).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a user purchase activity
func (s *UserPurchaseActivities) Insert(t *db.UserPurchaseActivity, txes ...*sql.Tx) (*db.UserPurchaseActivity, error) {
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

// Update a user purchase activity
func (s *UserPurchaseActivities) Update(u *db.UserPurchaseActivity, txes ...*sql.Tx) (*db.UserPurchaseActivity, error) {
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Create a user purchase activity
func (s *UserPurchaseActivities) Create(input *db.UserPurchaseActivity, txes ...*sql.Tx) (*db.UserPurchaseActivity, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}
