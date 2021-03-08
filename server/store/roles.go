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
)

// Role for persistence
type Role struct {
	Conn *sqlx.DB
}

// NewRoleStore handle role methods
func NewRoleStore(conn *sqlx.DB) *Role {
	ts := &Role{conn}
	return ts
}

// Get roles
func (s *Role) Get(id uuid.UUID) (*db.Role, error) {
	return db.FindRole(s.Conn, id.String())
}

// GetWithTrackActions gets a role with their track actions
func (s *Role) GetWithTrackActions(id uuid.UUID) (*db.Role, error) {
	return db.Roles(
		db.RoleWhere.ID.EQ(id.String()),
		qm.Load(db.RoleRels.TrackActions),
	).One(s.Conn)
}

// GetMany roles given a list of IDs
func (s *Role) GetMany(keys []string) (db.RoleSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Roles(db.RoleWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Role{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Role{}
	for _, key := range keys {
		var row *db.Role
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

// GetByName returns the role given its name
func (s *Role) GetByName(name string) (*db.Role, error) {
	return db.Roles(db.RoleWhere.Name.EQ(name)).One(s.Conn)
}

// GetByUser role by user id
func (s *Role) GetByUser(userID uuid.UUID) (*db.Role, error) {
	u, err := db.FindUser(s.Conn, userID.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u.Role().One(s.Conn)
}

// All roles
func (s *Role) All() (db.RoleSlice, error) {
	return db.Roles().All(s.Conn)
}

// SearchSelect searchs/selects roles
func (s *Role) SearchSelect(search graphql.SearchFilter, limit int, offset int, excludeSuper bool) (int64, []*db.Role, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("LOWER(%s) LIKE ?", db.RoleColumns.Name),
					"%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.RoleWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.RoleWhere.Archived.EQ(true))
		}
	}

	if excludeSuper {
		queries = append(queries, db.RoleWhere.Name.NEQ("SUPERADMIN"))
	}

	// Get Count
	count, err := db.Roles(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.RoleColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.RoleColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.RoleColumns.Name+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.RoleColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Roles(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// Insert roles
func (s *Role) Insert(record *db.Role, txes ...*sql.Tx) (*db.Role, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update roles
func (s *Role) Update(record *db.Role, txes ...*sql.Tx) (*db.Role, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive roles
func (s *Role) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Role, error) {
	u, err := db.FindRole(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.RoleColumns.Archived, db.RoleColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive roles
func (s *Role) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Role, error) {
	u, err := db.FindRole(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.RoleColumns.Archived, db.RoleColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// GetTrackActions returns the role's track actions
func (s *Role) GetTrackActions(record *db.Role) (db.TrackActionSlice, error) {
	if record.Tier == 1 {
		// Super Admin - return all track actions
		result, err := db.TrackActions(db.TrackActionWhere.Archived.EQ(false), db.TrackActionWhere.System.EQ(false)).All(s.Conn)
		if err != nil {
			return nil, terror.New(err, "")
		}
		return result, nil
	}

	result, err := record.TrackActions(db.TrackActionWhere.Archived.EQ(false), db.TrackActionWhere.System.EQ(false)).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return result, nil
}

// SetTrackActions sets the role's track actions
func (s *Role) SetTrackActions(record *db.Role, actions db.TrackActionSlice) error {
	return record.SetTrackActions(s.Conn, false, actions...)
}
