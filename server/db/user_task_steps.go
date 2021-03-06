// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// UserTaskStep is an object representing the database table.
type UserTaskStep struct {
	ID          string      `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UserTaskID  null.String `db:"user_task_id" boil:"user_task_id" json:"user_task_id,omitempty" toml:"user_task_id" yaml:"user_task_id,omitempty"`
	Name        string      `db:"name" boil:"name" json:"name" toml:"name" yaml:"name"`
	Description string      `db:"description" boil:"description" json:"description" toml:"description" yaml:"description"`
	Status      string      `db:"status" boil:"status" json:"status" toml:"status" yaml:"status"`
	IsComplete  bool        `db:"is_complete" boil:"is_complete" json:"is_complete" toml:"is_complete" yaml:"is_complete"`
	IsActive    bool        `db:"is_active" boil:"is_active" json:"is_active" toml:"is_active" yaml:"is_active"`
	UpdatedAt   time.Time   `db:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedAt   time.Time   `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *userTaskStepR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userTaskStepL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserTaskStepColumns = struct {
	ID          string
	UserTaskID  string
	Name        string
	Description string
	Status      string
	IsComplete  string
	IsActive    string
	UpdatedAt   string
	CreatedAt   string
}{
	ID:          "id",
	UserTaskID:  "user_task_id",
	Name:        "name",
	Description: "description",
	Status:      "status",
	IsComplete:  "is_complete",
	IsActive:    "is_active",
	UpdatedAt:   "updated_at",
	CreatedAt:   "created_at",
}

// Generated where

var UserTaskStepWhere = struct {
	ID          whereHelperstring
	UserTaskID  whereHelpernull_String
	Name        whereHelperstring
	Description whereHelperstring
	Status      whereHelperstring
	IsComplete  whereHelperbool
	IsActive    whereHelperbool
	UpdatedAt   whereHelpertime_Time
	CreatedAt   whereHelpertime_Time
}{
	ID:          whereHelperstring{field: "\"user_task_steps\".\"id\""},
	UserTaskID:  whereHelpernull_String{field: "\"user_task_steps\".\"user_task_id\""},
	Name:        whereHelperstring{field: "\"user_task_steps\".\"name\""},
	Description: whereHelperstring{field: "\"user_task_steps\".\"description\""},
	Status:      whereHelperstring{field: "\"user_task_steps\".\"status\""},
	IsComplete:  whereHelperbool{field: "\"user_task_steps\".\"is_complete\""},
	IsActive:    whereHelperbool{field: "\"user_task_steps\".\"is_active\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"user_task_steps\".\"updated_at\""},
	CreatedAt:   whereHelpertime_Time{field: "\"user_task_steps\".\"created_at\""},
}

// UserTaskStepRels is where relationship names are stored.
var UserTaskStepRels = struct {
	UserTask string
}{
	UserTask: "UserTask",
}

// userTaskStepR is where relationships are stored.
type userTaskStepR struct {
	UserTask *UserTask
}

// NewStruct creates a new relationship struct
func (*userTaskStepR) NewStruct() *userTaskStepR {
	return &userTaskStepR{}
}

// userTaskStepL is where Load methods for each relationship are stored.
type userTaskStepL struct{}

var (
	userTaskStepAllColumns            = []string{"id", "user_task_id", "name", "description", "status", "is_complete", "is_active", "updated_at", "created_at"}
	userTaskStepColumnsWithoutDefault = []string{"user_task_id", "name", "description", "status"}
	userTaskStepColumnsWithDefault    = []string{"id", "is_complete", "is_active", "updated_at", "created_at"}
	userTaskStepPrimaryKeyColumns     = []string{"id"}
)

type (
	// UserTaskStepSlice is an alias for a slice of pointers to UserTaskStep.
	// This should generally be used opposed to []UserTaskStep.
	UserTaskStepSlice []*UserTaskStep
	// UserTaskStepHook is the signature for custom UserTaskStep hook methods
	UserTaskStepHook func(boil.Executor, *UserTaskStep) error

	userTaskStepQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userTaskStepType                 = reflect.TypeOf(&UserTaskStep{})
	userTaskStepMapping              = queries.MakeStructMapping(userTaskStepType)
	userTaskStepPrimaryKeyMapping, _ = queries.BindMapping(userTaskStepType, userTaskStepMapping, userTaskStepPrimaryKeyColumns)
	userTaskStepInsertCacheMut       sync.RWMutex
	userTaskStepInsertCache          = make(map[string]insertCache)
	userTaskStepUpdateCacheMut       sync.RWMutex
	userTaskStepUpdateCache          = make(map[string]updateCache)
	userTaskStepUpsertCacheMut       sync.RWMutex
	userTaskStepUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userTaskStepBeforeInsertHooks []UserTaskStepHook
var userTaskStepBeforeUpdateHooks []UserTaskStepHook
var userTaskStepBeforeDeleteHooks []UserTaskStepHook
var userTaskStepBeforeUpsertHooks []UserTaskStepHook

var userTaskStepAfterInsertHooks []UserTaskStepHook
var userTaskStepAfterSelectHooks []UserTaskStepHook
var userTaskStepAfterUpdateHooks []UserTaskStepHook
var userTaskStepAfterDeleteHooks []UserTaskStepHook
var userTaskStepAfterUpsertHooks []UserTaskStepHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserTaskStep) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserTaskStep) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserTaskStep) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserTaskStep) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserTaskStep) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserTaskStep) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserTaskStep) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserTaskStep) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserTaskStep) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userTaskStepAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserTaskStepHook registers your hook function for all future operations.
func AddUserTaskStepHook(hookPoint boil.HookPoint, userTaskStepHook UserTaskStepHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userTaskStepBeforeInsertHooks = append(userTaskStepBeforeInsertHooks, userTaskStepHook)
	case boil.BeforeUpdateHook:
		userTaskStepBeforeUpdateHooks = append(userTaskStepBeforeUpdateHooks, userTaskStepHook)
	case boil.BeforeDeleteHook:
		userTaskStepBeforeDeleteHooks = append(userTaskStepBeforeDeleteHooks, userTaskStepHook)
	case boil.BeforeUpsertHook:
		userTaskStepBeforeUpsertHooks = append(userTaskStepBeforeUpsertHooks, userTaskStepHook)
	case boil.AfterInsertHook:
		userTaskStepAfterInsertHooks = append(userTaskStepAfterInsertHooks, userTaskStepHook)
	case boil.AfterSelectHook:
		userTaskStepAfterSelectHooks = append(userTaskStepAfterSelectHooks, userTaskStepHook)
	case boil.AfterUpdateHook:
		userTaskStepAfterUpdateHooks = append(userTaskStepAfterUpdateHooks, userTaskStepHook)
	case boil.AfterDeleteHook:
		userTaskStepAfterDeleteHooks = append(userTaskStepAfterDeleteHooks, userTaskStepHook)
	case boil.AfterUpsertHook:
		userTaskStepAfterUpsertHooks = append(userTaskStepAfterUpsertHooks, userTaskStepHook)
	}
}

// One returns a single userTaskStep record from the query.
func (q userTaskStepQuery) One(exec boil.Executor) (*UserTaskStep, error) {
	o := &UserTaskStep{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for user_task_steps")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserTaskStep records from the query.
func (q userTaskStepQuery) All(exec boil.Executor) (UserTaskStepSlice, error) {
	var o []*UserTaskStep

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to UserTaskStep slice")
	}

	if len(userTaskStepAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserTaskStep records in the query.
func (q userTaskStepQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count user_task_steps rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userTaskStepQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if user_task_steps exists")
	}

	return count > 0, nil
}

// UserTask pointed to by the foreign key.
func (o *UserTaskStep) UserTask(mods ...qm.QueryMod) userTaskQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserTaskID),
	}

	queryMods = append(queryMods, mods...)

	query := UserTasks(queryMods...)
	queries.SetFrom(query.Query, "\"user_tasks\"")

	return query
}

// LoadUserTask allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userTaskStepL) LoadUserTask(e boil.Executor, singular bool, maybeUserTaskStep interface{}, mods queries.Applicator) error {
	var slice []*UserTaskStep
	var object *UserTaskStep

	if singular {
		object = maybeUserTaskStep.(*UserTaskStep)
	} else {
		slice = *maybeUserTaskStep.(*[]*UserTaskStep)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userTaskStepR{}
		}
		if !queries.IsNil(object.UserTaskID) {
			args = append(args, object.UserTaskID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userTaskStepR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.UserTaskID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.UserTaskID) {
				args = append(args, obj.UserTaskID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`user_tasks`), qm.WhereIn(`user_tasks.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UserTask")
	}

	var resultSlice []*UserTask
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UserTask")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for user_tasks")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_tasks")
	}

	if len(userTaskStepAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.UserTask = foreign
		if foreign.R == nil {
			foreign.R = &userTaskR{}
		}
		foreign.R.UserTaskSteps = append(foreign.R.UserTaskSteps, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.UserTaskID, foreign.ID) {
				local.R.UserTask = foreign
				if foreign.R == nil {
					foreign.R = &userTaskR{}
				}
				foreign.R.UserTaskSteps = append(foreign.R.UserTaskSteps, local)
				break
			}
		}
	}

	return nil
}

// SetUserTask of the userTaskStep to the related item.
// Sets o.R.UserTask to related.
// Adds o to related.R.UserTaskSteps.
func (o *UserTaskStep) SetUserTask(exec boil.Executor, insert bool, related *UserTask) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_task_steps\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_task_id"}),
		strmangle.WhereClause("\"", "\"", 2, userTaskStepPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.UserTaskID, related.ID)
	if o.R == nil {
		o.R = &userTaskStepR{
			UserTask: related,
		}
	} else {
		o.R.UserTask = related
	}

	if related.R == nil {
		related.R = &userTaskR{
			UserTaskSteps: UserTaskStepSlice{o},
		}
	} else {
		related.R.UserTaskSteps = append(related.R.UserTaskSteps, o)
	}

	return nil
}

// RemoveUserTask relationship.
// Sets o.R.UserTask to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *UserTaskStep) RemoveUserTask(exec boil.Executor, related *UserTask) error {
	var err error

	queries.SetScanner(&o.UserTaskID, nil)
	if _, err = o.Update(exec, boil.Whitelist("user_task_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.UserTask = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.UserTaskSteps {
		if queries.Equal(o.UserTaskID, ri.UserTaskID) {
			continue
		}

		ln := len(related.R.UserTaskSteps)
		if ln > 1 && i < ln-1 {
			related.R.UserTaskSteps[i] = related.R.UserTaskSteps[ln-1]
		}
		related.R.UserTaskSteps = related.R.UserTaskSteps[:ln-1]
		break
	}
	return nil
}

// UserTaskSteps retrieves all the records using an executor.
func UserTaskSteps(mods ...qm.QueryMod) userTaskStepQuery {
	mods = append(mods, qm.From("\"user_task_steps\""))
	return userTaskStepQuery{NewQuery(mods...)}
}

// FindUserTaskStep retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserTaskStep(exec boil.Executor, iD string, selectCols ...string) (*UserTaskStep, error) {
	userTaskStepObj := &UserTaskStep{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_task_steps\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, userTaskStepObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from user_task_steps")
	}

	return userTaskStepObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserTaskStep) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_task_steps provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userTaskStepColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userTaskStepInsertCacheMut.RLock()
	cache, cached := userTaskStepInsertCache[key]
	userTaskStepInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userTaskStepAllColumns,
			userTaskStepColumnsWithDefault,
			userTaskStepColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userTaskStepType, userTaskStepMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userTaskStepType, userTaskStepMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_task_steps\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_task_steps\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "db: unable to insert into user_task_steps")
	}

	if !cached {
		userTaskStepInsertCacheMut.Lock()
		userTaskStepInsertCache[key] = cache
		userTaskStepInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the UserTaskStep.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserTaskStep) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userTaskStepUpdateCacheMut.RLock()
	cache, cached := userTaskStepUpdateCache[key]
	userTaskStepUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userTaskStepAllColumns,
			userTaskStepPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update user_task_steps, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_task_steps\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userTaskStepPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userTaskStepType, userTaskStepMapping, append(wl, userTaskStepPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update user_task_steps row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for user_task_steps")
	}

	if !cached {
		userTaskStepUpdateCacheMut.Lock()
		userTaskStepUpdateCache[key] = cache
		userTaskStepUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userTaskStepQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for user_task_steps")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for user_task_steps")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserTaskStepSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("db: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userTaskStepPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_task_steps\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userTaskStepPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in userTaskStep slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all userTaskStep")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserTaskStep) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_task_steps provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userTaskStepColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userTaskStepUpsertCacheMut.RLock()
	cache, cached := userTaskStepUpsertCache[key]
	userTaskStepUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userTaskStepAllColumns,
			userTaskStepColumnsWithDefault,
			userTaskStepColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userTaskStepAllColumns,
			userTaskStepPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert user_task_steps, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userTaskStepPrimaryKeyColumns))
			copy(conflict, userTaskStepPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_task_steps\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userTaskStepType, userTaskStepMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userTaskStepType, userTaskStepMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "db: unable to upsert user_task_steps")
	}

	if !cached {
		userTaskStepUpsertCacheMut.Lock()
		userTaskStepUpsertCache[key] = cache
		userTaskStepUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single UserTaskStep record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserTaskStep) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no UserTaskStep provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userTaskStepPrimaryKeyMapping)
	sql := "DELETE FROM \"user_task_steps\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from user_task_steps")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for user_task_steps")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userTaskStepQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no userTaskStepQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from user_task_steps")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_task_steps")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserTaskStepSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userTaskStepBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userTaskStepPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_task_steps\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userTaskStepPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from userTaskStep slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_task_steps")
	}

	if len(userTaskStepAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserTaskStep) Reload(exec boil.Executor) error {
	ret, err := FindUserTaskStep(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserTaskStepSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserTaskStepSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userTaskStepPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_task_steps\".* FROM \"user_task_steps\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userTaskStepPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in UserTaskStepSlice")
	}

	*o = slice

	return nil
}

// UserTaskStepExists checks if the UserTaskStep row exists.
func UserTaskStepExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_task_steps\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if user_task_steps exists")
	}

	return exists, nil
}
