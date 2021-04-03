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

// UserFactory creates users
func UserFactory() *db.User {
	u := &db.User{
		ID:           uuid.Must(uuid.NewV4()).String(),
		FirstName:    null.StringFrom(faker.Name().FirstName()),
		LastName:     null.StringFrom(faker.Name().LastName()),
		Email:        null.StringFrom(faker.Internet().Email()),
		PasswordHash: faker.Internet().Password(8, 20),
	}
	return u
}

// NewUserStore returns a new user repo that implements UserMutator, UserArchiver and UserQueryer
func NewUserStore(conn *sqlx.DB) *Users {
	r := &Users{conn}
	return r
}

// Users for persistence
type Users struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *Users) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// GetByVerifyToken returns a user with the matching verify token
func (s *Users) GetByVerifyToken(token string, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.Users(db.UserWhere.VerifyToken.EQ(token)).One(s.Conn)
	if err != nil {
		return nil, fmt.Errorf("Invalid token")
	}
	if u.VerifyTokenExpires.Before(time.Now()) {
		return nil, fmt.Errorf("Token expired")
	}
	return u, nil
}

// GetByResetToken returns a user with the matching reset token
func (s *Users) GetByResetToken(token string, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.Users(db.UserWhere.ResetToken.EQ(token)).One(s.Conn)
	if err != nil {
		return nil, fmt.Errorf("Invalid token")
	}
	if u.ResetTokenExpires.Before(time.Now()) {
		return nil, fmt.Errorf("Token expired")
	}
	return u, nil
}

// All users
func (s *Users) All(txes ...*sql.Tx) (db.UserSlice, error) {
	return db.Users().All(s.Conn)
}

// SearchSelect searchs/selects users
func (s *Users) SearchSelect(search graphql.SearchFilter, limit int, offset int, consumers bool) (int64, []*db.User, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.UserColumns.FirstName,
						db.UserColumns.LastName,
						db.UserColumns.Email,
					),
					"%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.UserWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.UserWhere.Archived.EQ(true))
		}
	}

	if consumers {
		queries = append(queries, db.UserWhere.WechatID.NEQ(null.StringFromPtr(nil)))
	} else {
		queries = append(queries, db.UserWhere.WechatID.EQ(null.StringFromPtr(nil)))
	}

	// Get Count
	count, err := db.Users(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.UserColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.UserColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.UserColumns.FirstName+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.UserColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Users(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetByOrganisation users by org
func (s *Users) GetByOrganisation(orgID uuid.UUID, txes ...*sql.Tx) (db.UserSlice, error) {
	dat, err := db.Users(db.UserWhere.OrganisationID.EQ(null.StringFrom(orgID.String()))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetReferrals users by referredByID
func (s *Users) GetReferrals(refByID string, txes ...*sql.Tx) (db.ReferralSlice, error) {
	dat, err := db.Referrals(db.ReferralWhere.ReferredByID.EQ(null.StringFrom(refByID))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetWalletHistory users by referredByID
func (s *Users) GetWalletHistory(userID string, txes ...*sql.Tx) (db.WalletHistorySlice, error) {
	dat, err := db.WalletHistories(db.WalletHistoryWhere.UserID.EQ(userID)).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany users given a list of IDs
func (s *Users) GetMany(keys []string, txes ...*sql.Tx) (db.UserSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Users(db.UserWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.User{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.User{}
	for _, key := range keys {
		var row *db.User
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

// Get a user given their ID
func (s *Users) Get(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(db.UserWhere.ID.EQ(id.String()),
		qm.Load(db.UserRels.Role, qm.Select(db.RoleColumns.ID, db.RoleColumns.Name, db.RoleColumns.Tier)),
	).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetWithPermissions gets a user given their ID w/ role permissions
func (s *Users) GetWithPermissions(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(
		db.UserWhere.ID.EQ(id.String()),
		qm.Load(db.UserRels.Role,
			qm.Select(
				db.RoleColumns.ID,
				db.RoleColumns.Name,
				db.RoleColumns.Tier,
				db.RoleColumns.Permissions,
			),
		),
	).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a user
func (s *Users) Insert(u *db.User, txes ...*sql.Tx) (*db.User, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return u.Insert(tx, boil.Infer())
	}, txes...)

	err = u.Reload(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return u, nil
}

// Update a user
func (s *Users) Update(u *db.User, txes ...*sql.Tx) (*db.User, error) {
	u.UpdatedAt = time.Now()
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// GetByEmail returns a user given an email
func (s *Users) GetByEmail(email string, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(db.UserWhere.Email.EQ(null.StringFrom(strings.ToLower(email)))).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByReferralCode returns a user given an referralCode
func (s *Users) GetByReferralCode(referralCode string, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(db.UserWhere.ReferralCode.EQ(null.StringFrom(referralCode))).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByWechatID returns a user given a wechat_id
func (s *Users) GetByWechatID(wechatID string, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(db.UserWhere.WechatID.EQ(null.StringFrom(wechatID))).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Create a user
func (s *Users) Create(input *db.User, txes ...*sql.Tx) (*db.User, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}

// Archive will archive users
func (s *Users) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.FindUser(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.UserColumns.Archived, db.UserColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive users
func (s *Users) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error) {
	u, err := db.FindUser(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.UserColumns.Archived, db.UserColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}
