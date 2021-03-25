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

// UserTaskFactory creates tasks
func UserTaskFactory() *db.UserTask {
	t := &db.UserTask{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
	return t
}

// NewUserTaskStore returns a new task repo that implements UserTaskMutator, UserTaskArchiver and UserTaskQueryer
func NewUserTaskStore(conn *sqlx.DB) *UserTasks {
	r := &UserTasks{conn}
	return r
}

// UserTasks for persistence
type UserTasks struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *UserTasks) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All tasks
func (s *UserTasks) All(txes ...*sql.Tx) (db.UserTaskSlice, error) {
	return db.UserTasks().All(s.Conn)
}

// Count gives the amount of tasks
func (s *UserTasks) Count() (int64, error) {
	return db.UserTasks().Count(s.Conn)
}

// SearchSelect searchs/selects tasks
func (s *UserTasks) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.UserTask, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?)",
						db.UserTaskColumns.ID,
					),
					"%"+searchText+"%",
				))
		}
	}

	// Get Count
	count, err := db.UserTasks(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.UserTaskColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.UserTaskColumns.UpdatedAt+sortDir))
			// case graphql.SortByOptionAlphabetical:
			// 	queries = append(queries, qm.OrderBy(db.UserTaskColumns.FirstName+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.UserTaskColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.UserTasks(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetByCode by code
func (s *UserTasks) GetByCode(code string, txes ...*sql.Tx) (*db.UserTask, error) {
	dat, err := db.UserTasks(db.UserTaskWhere.Code.EQ(code)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetTask by taskID
func (s *UserTasks) GetTask(taskID string, txes ...*sql.Tx) (*db.Task, error) {
	dat, err := db.Tasks(db.TaskWhere.ID.EQ(taskID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetUserSubtask by subTaskID
func (s *UserTasks) GetSubtask(subTaskID string, txes ...*sql.Tx) (*db.Subtask, error) {
	dat, err := db.Subtasks(db.SubtaskWhere.ID.EQ(subTaskID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetUser by userID
func (s *UserTasks) GetUser(userID string, txes ...*sql.Tx) (*db.User, error) {
	dat, err := db.Users(db.UserWhere.ID.EQ(userID)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetSubtasks by userTaskID
func (s *UserTasks) GetSubtasks(userTaskID string, txes ...*sql.Tx) (db.UserSubtaskSlice, error) {
	dat, err := db.UserSubtasks(db.UserSubtaskWhere.UserTaskID.EQ(null.StringFrom(userTaskID))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany tasks given a list of IDs
func (s *UserTasks) GetMany(keys []string, txes ...*sql.Tx) (db.UserTaskSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.UserTasks(db.UserTaskWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.UserTask{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.UserTask{}
	for _, key := range keys {
		var row *db.UserTask
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
func (s *UserTasks) Get(id uuid.UUID, txes ...*sql.Tx) (*db.UserTask, error) {
	dat, err := db.UserTasks(db.UserTaskWhere.ID.EQ(id.String()),
		qm.Load(db.UserTaskRels.UserSubtasks, qm.Select(db.UserSubtaskColumns.ID, db.UserSubtaskColumns.UserTaskID, db.UserSubtaskColumns.SubtaskID)),
	).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a task
func (s *UserTasks) Insert(t *db.UserTask, txes ...*sql.Tx) (*db.UserTask, error) {
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
func (s *UserTasks) InsertSubtask(st *db.UserSubtask, txes ...*sql.Tx) (*db.UserSubtask, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return st.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return st, nil
}

// Update a task
func (s *UserTasks) Update(u *db.UserTask, txes ...*sql.Tx) (*db.UserTask, error) {
	u.UpdatedAt = time.Now()
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Create a task
func (s *UserTasks) Create(input *db.UserTask, txes ...*sql.Tx) (*db.UserTask, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}
