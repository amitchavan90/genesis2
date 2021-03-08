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

// NewTrackActionStore handle trackAction methods
func NewTrackActionStore(conn *sqlx.DB) *TrackAction {
	ts := &TrackAction{conn}
	return ts
}

// TrackAction for persistence
type TrackAction struct {
	Conn *sqlx.DB
}

// Get trackAction
func (s *TrackAction) Get(id uuid.UUID) (*db.TrackAction, error) {
	return db.FindTrackAction(s.Conn, id.String())
}

// GetMany trackActions
func (s *TrackAction) GetMany(keys []string) (db.TrackActionSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.TrackActions(db.TrackActionWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.TrackAction{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.TrackAction{}
	for _, key := range keys {
		var row *db.TrackAction
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

// GetByCode get trackAction by code
func (s *TrackAction) GetByCode(code string) (*db.TrackAction, error) {
	return db.TrackActions(db.TrackActionWhere.Code.EQ(code)).One(s.Conn)
}

// All trackActions
func (s *TrackAction) All() (db.TrackActionSlice, error) {
	return db.TrackActions().All(s.Conn)
}

// SearchSelect searchs/selects TrackActions
func (s *TrackAction) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.TrackAction, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.TrackActionColumns.Name,
						db.TrackActionColumns.NameChinese,
					),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.TrackActionWhere.Archived.EQ(false), db.TrackActionWhere.System.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.TrackActionWhere.Archived.EQ(true))
		case graphql.FilterOptionSystem:
			queries = append(queries, db.TrackActionWhere.System.EQ(true))
		case graphql.FilterOptionBlockchain:
			queries = append(queries, db.TrackActionWhere.Blockchain.EQ(true))
		}
	}

	// Get Count
	count, err := db.TrackActions(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.TrackActionColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.TrackActionColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.TrackActionColumns.Name+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.TrackActionColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.TrackActions(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// Insert trackAction
func (s *TrackAction) Insert(record *db.TrackAction, txes ...*sql.Tx) (*db.TrackAction, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update trackAction
func (s *TrackAction) Update(record *db.TrackAction, txes ...*sql.Tx) (*db.TrackAction, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive trackActions
func (s *TrackAction) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.TrackAction, error) {
	u, err := db.FindTrackAction(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.System {
		return nil, fmt.Errorf("cannot archive system track action")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.TrackActionColumns.Archived, db.TrackActionColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive trackActions
func (s *TrackAction) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.TrackAction, error) {
	u, err := db.FindTrackAction(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.TrackActionColumns.Archived, db.TrackActionColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of trackActions
func (s *TrackAction) Count() (int64, error) {
	return db.TrackActions().Count(s.Conn)
}
