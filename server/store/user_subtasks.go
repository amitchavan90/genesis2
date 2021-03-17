package store

import (
	"database/sql"
	"errors"
	"genesis/db"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/sqlboiler/boil"
)

// UserSubtaskFactory creates tasks
func UserSubtaskFactory() *db.UserTask {
	t := &db.UserTask{
		ID: uuid.Must(uuid.NewV4()).String(),
	}
	return t
}

// NewUserSubtaskStore returns a new task repo that implements UserTaskMutator, UserTaskArchiver and UserTaskQueryer
func NewUserSubtaskStore(conn *sqlx.DB) *UserSubtasks {
	r := &UserSubtasks{conn}
	return r
}

// UserSubtasks for persistence
type UserSubtasks struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (s *UserSubtasks) BeginTransaction() (*sql.Tx, error) {
	return s.Conn.Begin()
}

// All tasks
func (s *UserSubtasks) All(txes ...*sql.Tx) (db.UserSubtaskSlice, error) {
	return db.UserSubtasks().All(s.Conn)
}

// Count gives the amount of tasks
func (s *UserSubtasks) Count() (int64, error) {
	return db.StockKeepingUnits().Count(s.Conn)
}

// GetMany tasks given a list of IDs
func (s *UserSubtasks) GetMany(keys []string, txes ...*sql.Tx) (db.UserSubtaskSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.UserSubtasks(db.UserSubtaskWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.UserSubtask{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.UserSubtask{}
	for _, key := range keys {
		var row *db.UserSubtask
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

// Get a user subtask given their ID
func (s *UserSubtasks) Get(id uuid.UUID, txes ...*sql.Tx) (*db.UserSubtask, error) {
	dat, err := db.UserSubtasks(db.UserSubtaskWhere.ID.EQ(id.String())).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Insert a task
func (s *UserSubtasks) Insert(st *db.UserSubtask, txes ...*sql.Tx) (*db.UserSubtask, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return st.Insert(tx, boil.Infer())
	}, txes...)

	err = st.Reload(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return st, nil
}

// Update a task
func (s *UserSubtasks) Update(u *db.UserTask, txes ...*sql.Tx) (*db.UserTask, error) {
	u.UpdatedAt = time.Now()
	_, err := u.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Create a task
func (s *UserSubtasks) Create(input *db.UserTask, txes ...*sql.Tx) (*db.UserTask, error) {
	err := input.Insert(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return input, nil
}
