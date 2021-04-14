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
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// NewUserActivityStore handle userActivity methods
func NewUserActivityStore(conn *sqlx.DB) *UserActivity {
	ts := &UserActivity{conn}
	return ts
}

// UserActivity for persistence
type UserActivity struct {
	Conn *sqlx.DB
}

// Get userActivity
func (s *UserActivity) Get(id uuid.UUID) (*db.UserActivity, error) {
	return db.FindUserActivity(s.Conn, id.String())
}

// GetMany userActivitys
func (s *UserActivity) GetMany(keys []string) (db.UserActivitySlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.UserActivities(db.UserActivityWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.UserActivity{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.UserActivity{}
	for _, key := range keys {
		var row *db.UserActivity
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

// All userActivitys
func (s *UserActivity) All() (db.UserActivitySlice, error) {
	return db.UserActivities().All(s.Conn)
}

// SearchSelect searchs/selects Products
func (s *UserActivity) SearchSelect(
	search graphql.SearchFilter,
	limit int,
	offset int,
	userID null.String,
) (int64, db.UserActivitySlice, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.UserActivityColumns.Action,
						db.UserActivityColumns.ObjectType,
						db.UserActivityColumns.ObjectCode,
					),
					"%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	if userID.Valid {
		queries = append(queries, db.UserActivityWhere.UserID.EQ(userID.String))
	}

	// Get Count
	count, err := db.UserActivities(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.UserActivityColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.UserActivityColumns.CreatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.UserActivityColumns.Action+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.UserActivityColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.UserActivities(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// Insert userActivity
func (s *UserActivity) Insert(record *db.UserActivity, txes ...*sql.Tx) (*db.UserActivity, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}
