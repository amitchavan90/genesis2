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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// IssuedToken is an object representing the database table.
type IssuedToken struct {
	ID           string    `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID       string    `db:"user_id" boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Device       string    `db:"device" boil:"device" json:"device" toml:"device" yaml:"device"`
	TokenCreated time.Time `db:"token_created" boil:"token_created" json:"token_created" toml:"token_created" yaml:"token_created"`
	TokenExpires time.Time `db:"token_expires" boil:"token_expires" json:"token_expires" toml:"token_expires" yaml:"token_expires"`
	Blacklisted  bool      `db:"blacklisted" boil:"blacklisted" json:"blacklisted" toml:"blacklisted" yaml:"blacklisted"`

	R *issuedTokenR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L issuedTokenL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var IssuedTokenColumns = struct {
	ID           string
	UserID       string
	Device       string
	TokenCreated string
	TokenExpires string
	Blacklisted  string
}{
	ID:           "id",
	UserID:       "user_id",
	Device:       "device",
	TokenCreated: "token_created",
	TokenExpires: "token_expires",
	Blacklisted:  "blacklisted",
}

// Generated where

var IssuedTokenWhere = struct {
	ID           whereHelperstring
	UserID       whereHelperstring
	Device       whereHelperstring
	TokenCreated whereHelpertime_Time
	TokenExpires whereHelpertime_Time
	Blacklisted  whereHelperbool
}{
	ID:           whereHelperstring{field: "\"issued_tokens\".\"id\""},
	UserID:       whereHelperstring{field: "\"issued_tokens\".\"user_id\""},
	Device:       whereHelperstring{field: "\"issued_tokens\".\"device\""},
	TokenCreated: whereHelpertime_Time{field: "\"issued_tokens\".\"token_created\""},
	TokenExpires: whereHelpertime_Time{field: "\"issued_tokens\".\"token_expires\""},
	Blacklisted:  whereHelperbool{field: "\"issued_tokens\".\"blacklisted\""},
}

// IssuedTokenRels is where relationship names are stored.
var IssuedTokenRels = struct {
	User string
}{
	User: "User",
}

// issuedTokenR is where relationships are stored.
type issuedTokenR struct {
	User *User
}

// NewStruct creates a new relationship struct
func (*issuedTokenR) NewStruct() *issuedTokenR {
	return &issuedTokenR{}
}

// issuedTokenL is where Load methods for each relationship are stored.
type issuedTokenL struct{}

var (
	issuedTokenAllColumns            = []string{"id", "user_id", "device", "token_created", "token_expires", "blacklisted"}
	issuedTokenColumnsWithoutDefault = []string{"user_id", "device", "token_expires"}
	issuedTokenColumnsWithDefault    = []string{"id", "token_created", "blacklisted"}
	issuedTokenPrimaryKeyColumns     = []string{"id"}
)

type (
	// IssuedTokenSlice is an alias for a slice of pointers to IssuedToken.
	// This should generally be used opposed to []IssuedToken.
	IssuedTokenSlice []*IssuedToken
	// IssuedTokenHook is the signature for custom IssuedToken hook methods
	IssuedTokenHook func(boil.Executor, *IssuedToken) error

	issuedTokenQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	issuedTokenType                 = reflect.TypeOf(&IssuedToken{})
	issuedTokenMapping              = queries.MakeStructMapping(issuedTokenType)
	issuedTokenPrimaryKeyMapping, _ = queries.BindMapping(issuedTokenType, issuedTokenMapping, issuedTokenPrimaryKeyColumns)
	issuedTokenInsertCacheMut       sync.RWMutex
	issuedTokenInsertCache          = make(map[string]insertCache)
	issuedTokenUpdateCacheMut       sync.RWMutex
	issuedTokenUpdateCache          = make(map[string]updateCache)
	issuedTokenUpsertCacheMut       sync.RWMutex
	issuedTokenUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var issuedTokenBeforeInsertHooks []IssuedTokenHook
var issuedTokenBeforeUpdateHooks []IssuedTokenHook
var issuedTokenBeforeDeleteHooks []IssuedTokenHook
var issuedTokenBeforeUpsertHooks []IssuedTokenHook

var issuedTokenAfterInsertHooks []IssuedTokenHook
var issuedTokenAfterSelectHooks []IssuedTokenHook
var issuedTokenAfterUpdateHooks []IssuedTokenHook
var issuedTokenAfterDeleteHooks []IssuedTokenHook
var issuedTokenAfterUpsertHooks []IssuedTokenHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *IssuedToken) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *IssuedToken) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *IssuedToken) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *IssuedToken) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *IssuedToken) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *IssuedToken) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *IssuedToken) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *IssuedToken) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *IssuedToken) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range issuedTokenAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddIssuedTokenHook registers your hook function for all future operations.
func AddIssuedTokenHook(hookPoint boil.HookPoint, issuedTokenHook IssuedTokenHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		issuedTokenBeforeInsertHooks = append(issuedTokenBeforeInsertHooks, issuedTokenHook)
	case boil.BeforeUpdateHook:
		issuedTokenBeforeUpdateHooks = append(issuedTokenBeforeUpdateHooks, issuedTokenHook)
	case boil.BeforeDeleteHook:
		issuedTokenBeforeDeleteHooks = append(issuedTokenBeforeDeleteHooks, issuedTokenHook)
	case boil.BeforeUpsertHook:
		issuedTokenBeforeUpsertHooks = append(issuedTokenBeforeUpsertHooks, issuedTokenHook)
	case boil.AfterInsertHook:
		issuedTokenAfterInsertHooks = append(issuedTokenAfterInsertHooks, issuedTokenHook)
	case boil.AfterSelectHook:
		issuedTokenAfterSelectHooks = append(issuedTokenAfterSelectHooks, issuedTokenHook)
	case boil.AfterUpdateHook:
		issuedTokenAfterUpdateHooks = append(issuedTokenAfterUpdateHooks, issuedTokenHook)
	case boil.AfterDeleteHook:
		issuedTokenAfterDeleteHooks = append(issuedTokenAfterDeleteHooks, issuedTokenHook)
	case boil.AfterUpsertHook:
		issuedTokenAfterUpsertHooks = append(issuedTokenAfterUpsertHooks, issuedTokenHook)
	}
}

// One returns a single issuedToken record from the query.
func (q issuedTokenQuery) One(exec boil.Executor) (*IssuedToken, error) {
	o := &IssuedToken{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for issued_tokens")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all IssuedToken records from the query.
func (q issuedTokenQuery) All(exec boil.Executor) (IssuedTokenSlice, error) {
	var o []*IssuedToken

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to IssuedToken slice")
	}

	if len(issuedTokenAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all IssuedToken records in the query.
func (q issuedTokenQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count issued_tokens rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q issuedTokenQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if issued_tokens exists")
	}

	return count > 0, nil
}

// User pointed to by the foreign key.
func (o *IssuedToken) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (issuedTokenL) LoadUser(e boil.Executor, singular bool, maybeIssuedToken interface{}, mods queries.Applicator) error {
	var slice []*IssuedToken
	var object *IssuedToken

	if singular {
		object = maybeIssuedToken.(*IssuedToken)
	} else {
		slice = *maybeIssuedToken.(*[]*IssuedToken)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &issuedTokenR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &issuedTokenR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`users`), qm.WhereIn(`users.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(issuedTokenAfterSelectHooks) != 0 {
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
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.IssuedTokens = append(foreign.R.IssuedTokens, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.IssuedTokens = append(foreign.R.IssuedTokens, local)
				break
			}
		}
	}

	return nil
}

// SetUser of the issuedToken to the related item.
// Sets o.R.User to related.
// Adds o to related.R.IssuedTokens.
func (o *IssuedToken) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"issued_tokens\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, issuedTokenPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &issuedTokenR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			IssuedTokens: IssuedTokenSlice{o},
		}
	} else {
		related.R.IssuedTokens = append(related.R.IssuedTokens, o)
	}

	return nil
}

// IssuedTokens retrieves all the records using an executor.
func IssuedTokens(mods ...qm.QueryMod) issuedTokenQuery {
	mods = append(mods, qm.From("\"issued_tokens\""))
	return issuedTokenQuery{NewQuery(mods...)}
}

// FindIssuedToken retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindIssuedToken(exec boil.Executor, iD string, selectCols ...string) (*IssuedToken, error) {
	issuedTokenObj := &IssuedToken{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"issued_tokens\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, issuedTokenObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from issued_tokens")
	}

	return issuedTokenObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *IssuedToken) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no issued_tokens provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(issuedTokenColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	issuedTokenInsertCacheMut.RLock()
	cache, cached := issuedTokenInsertCache[key]
	issuedTokenInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			issuedTokenAllColumns,
			issuedTokenColumnsWithDefault,
			issuedTokenColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(issuedTokenType, issuedTokenMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(issuedTokenType, issuedTokenMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"issued_tokens\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"issued_tokens\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into issued_tokens")
	}

	if !cached {
		issuedTokenInsertCacheMut.Lock()
		issuedTokenInsertCache[key] = cache
		issuedTokenInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the IssuedToken.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *IssuedToken) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	issuedTokenUpdateCacheMut.RLock()
	cache, cached := issuedTokenUpdateCache[key]
	issuedTokenUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			issuedTokenAllColumns,
			issuedTokenPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update issued_tokens, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"issued_tokens\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, issuedTokenPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(issuedTokenType, issuedTokenMapping, append(wl, issuedTokenPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update issued_tokens row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for issued_tokens")
	}

	if !cached {
		issuedTokenUpdateCacheMut.Lock()
		issuedTokenUpdateCache[key] = cache
		issuedTokenUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q issuedTokenQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for issued_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for issued_tokens")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o IssuedTokenSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuedTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"issued_tokens\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, issuedTokenPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in issuedToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all issuedToken")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *IssuedToken) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no issued_tokens provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(issuedTokenColumnsWithDefault, o)

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

	issuedTokenUpsertCacheMut.RLock()
	cache, cached := issuedTokenUpsertCache[key]
	issuedTokenUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			issuedTokenAllColumns,
			issuedTokenColumnsWithDefault,
			issuedTokenColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			issuedTokenAllColumns,
			issuedTokenPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert issued_tokens, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(issuedTokenPrimaryKeyColumns))
			copy(conflict, issuedTokenPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"issued_tokens\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(issuedTokenType, issuedTokenMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(issuedTokenType, issuedTokenMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert issued_tokens")
	}

	if !cached {
		issuedTokenUpsertCacheMut.Lock()
		issuedTokenUpsertCache[key] = cache
		issuedTokenUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single IssuedToken record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *IssuedToken) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no IssuedToken provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), issuedTokenPrimaryKeyMapping)
	sql := "DELETE FROM \"issued_tokens\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from issued_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for issued_tokens")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q issuedTokenQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no issuedTokenQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from issued_tokens")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for issued_tokens")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o IssuedTokenSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(issuedTokenBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuedTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"issued_tokens\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, issuedTokenPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from issuedToken slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for issued_tokens")
	}

	if len(issuedTokenAfterDeleteHooks) != 0 {
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
func (o *IssuedToken) Reload(exec boil.Executor) error {
	ret, err := FindIssuedToken(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *IssuedTokenSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := IssuedTokenSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), issuedTokenPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"issued_tokens\".* FROM \"issued_tokens\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, issuedTokenPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in IssuedTokenSlice")
	}

	*o = slice

	return nil
}

// IssuedTokenExists checks if the IssuedToken row exists.
func IssuedTokenExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"issued_tokens\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if issued_tokens exists")
	}

	return exists, nil
}
