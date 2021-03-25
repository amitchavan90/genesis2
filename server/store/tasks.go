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

// TaskFactory creates tasks
func TaskFactory() *db.Task {
	t := &db.Task{
		ID:            uuid.Must(uuid.NewV4()).String(),
		Title:         faker.Lorem().Word(),
		Description:   faker.Lorem().Word(),
		LoyaltyPoints: faker.Number().NumberInt(3),
	}
	return t
}

// NewTaskStore returns a new task repo that implements TaskMutator, TaskArchiver and TaskQueryer
func NewTaskStore(conn *sqlx.DB) *Tasks {
	r := &Tasks{conn}
	return r
}

// Tasks for persistence
type Tasks struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *Tasks) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All tasks
func (s *Tasks) All(txes ...*sql.Tx) (db.TaskSlice, error) {
	return db.Tasks().All(s.Conn)
}

// Count gives the amount of tasks
func (s *Tasks) Count() (int64, error) {
	return db.Tasks().Count(s.Conn)
}

// SearchSelect searchs/selects tasks
func (s *Tasks) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Task, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?)",
						db.TaskColumns.ID,
					),
					"%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.TaskWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.TaskWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.Tasks(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.TaskColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.TaskColumns.UpdatedAt+sortDir))
			// case graphql.SortByOptionAlphabetical:
			// 	queries = append(queries, qm.OrderBy(db.TaskColumns.FirstName+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.TaskColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Tasks(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetByCode by code
func (s *Tasks) GetByCode(code string, txes ...*sql.Tx) (*db.Task, error) {
	dat, err := db.Tasks(db.TaskWhere.Code.EQ(code)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetSku by skuID
func (s *Tasks) GetSku(skuID string, txes ...*sql.Tx) (*db.StockKeepingUnit, error) {
	dat, err := db.StockKeepingUnits(db.StockKeepingUnitWhere.ID.EQ(skuID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetSubtasks by taskID
func (s *Tasks) GetSubtasks(taskID string, txes ...*sql.Tx) (db.SubtaskSlice, error) {
	dat, err := db.Subtasks(db.SubtaskWhere.TaskID.EQ(null.StringFrom(taskID))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany tasks given a list of IDs
func (s *Tasks) GetMany(keys []string, txes ...*sql.Tx) (db.TaskSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Tasks(db.TaskWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Task{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Task{}
	for _, key := range keys {
		var row *db.Task
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

// Get a task given their ID
func (s *Tasks) Get(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error) {
	dat, err := db.Tasks(db.TaskWhere.ID.EQ(id.String()),
		qm.Load(db.TaskRels.Subtasks, qm.Select(db.SubtaskColumns.ID, db.SubtaskColumns.Title, db.SubtaskColumns.Description)),
	).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a task
func (s *Tasks) Insert(t *db.Task, txes ...*sql.Tx) (*db.Task, error) {
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

// InsertSubtask skus
func (s *Tasks) InsertSubtask(st *db.Subtask, txes ...*sql.Tx) (*db.Subtask, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return st.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return st, nil
}

// Update a task
func (s *Tasks) Update(u *db.Task, txes ...*sql.Tx) (*db.Task, error) {
	u.UpdatedAt = time.Now()
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Create a task
func (s *Tasks) Create(input *db.Task, txes ...*sql.Tx) (*db.Task, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}

// Archive will archive tasks
func (s *Tasks) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error) {
	u, err := db.FindTask(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.TaskColumns.Archived, db.TaskColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive tasks
func (s *Tasks) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error) {
	u, err := db.FindTask(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.TaskColumns.Archived, db.TaskColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}
