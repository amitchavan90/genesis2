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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// ReferralFactory creates referrals
func ReferralFactory() *db.Referral {
	u := &db.Referral{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
	return u
}

// NewReferralStore returns a new referral repo that implements ReferralMutator, ReferralArchiver and ReferralQueryer
func NewReferralStore(conn *sqlx.DB) *Referrals {
	r := &Referrals{conn}
	return r
}

// Referrals for persistence
type Referrals struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *Referrals) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All referrals
func (s *Referrals) All(txes ...*sql.Tx) (db.ReferralSlice, error) {
	return db.Referrals().All(s.Conn)
}

// SearchSelect searchs/selects referrals
func (s *Referrals) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Referral, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.ReferralColumns.UserID,
						db.ReferralColumns.ReferredByID,
					),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Get Count
	count, err := db.Referrals(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.ReferralColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.ReferralColumns.UpdatedAt+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.ReferralColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Referrals(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetByReferee referrals by org
// func (s *Referrals) GetByReferee(orgID uuid.UUID, txes ...*sql.Tx) (db.ReferralSlice, error) {
// 	dat, err := db.Referrals(db.ReferralWhere.ReferredByID.EQ(null.StringFrom(orgID.String()))).All(s.Conn)
// 	if err != nil {
// 		return nil, terror.New(err, "")
// 	}
// 	return dat, nil
// }

// GetMany referrals given a list of IDs
func (s *Referrals) GetMany(keys []string, txes ...*sql.Tx) (db.ReferralSlice, []error) {
	fmt.Println("Referral Resolvers: GetMany")
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Referrals(db.ReferralWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Referral{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Referral{}
	for _, key := range keys {
		var row *db.Referral
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

// Get a referral given their ID
func (s *Referrals) Get(id uuid.UUID, txes ...*sql.Tx) (*db.Referral, error) {
	dat, err := db.Referrals(db.ReferralWhere.ID.EQ(id.String()),
		qm.Load(db.ReferralRels.ReferredBy, qm.Select(db.ReferralColumns.ID)),
	).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByUserID returns a referral given an userID
func (s *Referrals) GetByUserID(userID string, txes ...*sql.Tx) (*db.Referral, error) {
	dat, err := db.Referrals(db.ReferralWhere.UserID.EQ(userID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a referral
func (s *Referrals) Insert(u *db.Referral, txes ...*sql.Tx) (*db.Referral, error) {
	var err error

	fmt.Println("Referral Store: Insert")

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return u.Insert(tx, boil.Infer())
	}, txes...)

	err = u.Reload(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return u, nil
}

// Update a referral
func (s *Referrals) Update(u *db.Referral, txes ...*sql.Tx) (*db.Referral, error) {
	u.UpdatedAt = time.Now()
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Create a referral
func (s *Referrals) Create(input *db.Referral, txes ...*sql.Tx) (*db.Referral, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}
