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

// Setting is an object representing the database table.
type Setting struct {
	ID                   string `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	SmartContractAddress string `db:"smart_contract_address" boil:"smart_contract_address" json:"smart_contract_address" toml:"smart_contract_address" yaml:"smart_contract_address"`

	R *settingR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L settingL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SettingColumns = struct {
	ID                   string
	SmartContractAddress string
}{
	ID:                   "id",
	SmartContractAddress: "smart_contract_address",
}

// Generated where

var SettingWhere = struct {
	ID                   whereHelperstring
	SmartContractAddress whereHelperstring
}{
	ID:                   whereHelperstring{field: "\"settings\".\"id\""},
	SmartContractAddress: whereHelperstring{field: "\"settings\".\"smart_contract_address\""},
}

// SettingRels is where relationship names are stored.
var SettingRels = struct {
}{}

// settingR is where relationships are stored.
type settingR struct {
}

// NewStruct creates a new relationship struct
func (*settingR) NewStruct() *settingR {
	return &settingR{}
}

// settingL is where Load methods for each relationship are stored.
type settingL struct{}

var (
	settingAllColumns            = []string{"id", "smart_contract_address"}
	settingColumnsWithoutDefault = []string{}
	settingColumnsWithDefault    = []string{"id", "smart_contract_address"}
	settingPrimaryKeyColumns     = []string{"id"}
)

type (
	// SettingSlice is an alias for a slice of pointers to Setting.
	// This should generally be used opposed to []Setting.
	SettingSlice []*Setting
	// SettingHook is the signature for custom Setting hook methods
	SettingHook func(boil.Executor, *Setting) error

	settingQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	settingType                 = reflect.TypeOf(&Setting{})
	settingMapping              = queries.MakeStructMapping(settingType)
	settingPrimaryKeyMapping, _ = queries.BindMapping(settingType, settingMapping, settingPrimaryKeyColumns)
	settingInsertCacheMut       sync.RWMutex
	settingInsertCache          = make(map[string]insertCache)
	settingUpdateCacheMut       sync.RWMutex
	settingUpdateCache          = make(map[string]updateCache)
	settingUpsertCacheMut       sync.RWMutex
	settingUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var settingBeforeInsertHooks []SettingHook
var settingBeforeUpdateHooks []SettingHook
var settingBeforeDeleteHooks []SettingHook
var settingBeforeUpsertHooks []SettingHook

var settingAfterInsertHooks []SettingHook
var settingAfterSelectHooks []SettingHook
var settingAfterUpdateHooks []SettingHook
var settingAfterDeleteHooks []SettingHook
var settingAfterUpsertHooks []SettingHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Setting) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range settingBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Setting) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range settingBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Setting) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range settingBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Setting) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range settingBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Setting) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range settingAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Setting) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range settingAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Setting) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range settingAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Setting) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range settingAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Setting) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range settingAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSettingHook registers your hook function for all future operations.
func AddSettingHook(hookPoint boil.HookPoint, settingHook SettingHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		settingBeforeInsertHooks = append(settingBeforeInsertHooks, settingHook)
	case boil.BeforeUpdateHook:
		settingBeforeUpdateHooks = append(settingBeforeUpdateHooks, settingHook)
	case boil.BeforeDeleteHook:
		settingBeforeDeleteHooks = append(settingBeforeDeleteHooks, settingHook)
	case boil.BeforeUpsertHook:
		settingBeforeUpsertHooks = append(settingBeforeUpsertHooks, settingHook)
	case boil.AfterInsertHook:
		settingAfterInsertHooks = append(settingAfterInsertHooks, settingHook)
	case boil.AfterSelectHook:
		settingAfterSelectHooks = append(settingAfterSelectHooks, settingHook)
	case boil.AfterUpdateHook:
		settingAfterUpdateHooks = append(settingAfterUpdateHooks, settingHook)
	case boil.AfterDeleteHook:
		settingAfterDeleteHooks = append(settingAfterDeleteHooks, settingHook)
	case boil.AfterUpsertHook:
		settingAfterUpsertHooks = append(settingAfterUpsertHooks, settingHook)
	}
}

// One returns a single setting record from the query.
func (q settingQuery) One(exec boil.Executor) (*Setting, error) {
	o := &Setting{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for settings")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Setting records from the query.
func (q settingQuery) All(exec boil.Executor) (SettingSlice, error) {
	var o []*Setting

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to Setting slice")
	}

	if len(settingAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Setting records in the query.
func (q settingQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count settings rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q settingQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if settings exists")
	}

	return count > 0, nil
}

// Settings retrieves all the records using an executor.
func Settings(mods ...qm.QueryMod) settingQuery {
	mods = append(mods, qm.From("\"settings\""))
	return settingQuery{NewQuery(mods...)}
}

// FindSetting retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSetting(exec boil.Executor, iD string, selectCols ...string) (*Setting, error) {
	settingObj := &Setting{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"settings\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, settingObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from settings")
	}

	return settingObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Setting) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no settings provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(settingColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	settingInsertCacheMut.RLock()
	cache, cached := settingInsertCache[key]
	settingInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			settingAllColumns,
			settingColumnsWithDefault,
			settingColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(settingType, settingMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(settingType, settingMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"settings\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"settings\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into settings")
	}

	if !cached {
		settingInsertCacheMut.Lock()
		settingInsertCache[key] = cache
		settingInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the Setting.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Setting) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	settingUpdateCacheMut.RLock()
	cache, cached := settingUpdateCache[key]
	settingUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			settingAllColumns,
			settingPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update settings, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"settings\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, settingPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(settingType, settingMapping, append(wl, settingPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update settings row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for settings")
	}

	if !cached {
		settingUpdateCacheMut.Lock()
		settingUpdateCache[key] = cache
		settingUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q settingQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for settings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for settings")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SettingSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), settingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"settings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, settingPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in setting slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all setting")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Setting) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no settings provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(settingColumnsWithDefault, o)

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

	settingUpsertCacheMut.RLock()
	cache, cached := settingUpsertCache[key]
	settingUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			settingAllColumns,
			settingColumnsWithDefault,
			settingColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			settingAllColumns,
			settingPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert settings, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(settingPrimaryKeyColumns))
			copy(conflict, settingPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"settings\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(settingType, settingMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(settingType, settingMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert settings")
	}

	if !cached {
		settingUpsertCacheMut.Lock()
		settingUpsertCache[key] = cache
		settingUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single Setting record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Setting) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no Setting provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), settingPrimaryKeyMapping)
	sql := "DELETE FROM \"settings\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from settings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for settings")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q settingQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no settingQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from settings")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for settings")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SettingSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(settingBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), settingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"settings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, settingPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from setting slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for settings")
	}

	if len(settingAfterDeleteHooks) != 0 {
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
func (o *Setting) Reload(exec boil.Executor) error {
	ret, err := FindSetting(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SettingSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SettingSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), settingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"settings\".* FROM \"settings\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, settingPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in SettingSlice")
	}

	*o = slice

	return nil
}

// SettingExists checks if the Setting row exists.
func SettingExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"settings\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if settings exists")
	}

	return exists, nil
}
