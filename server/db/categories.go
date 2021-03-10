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

// Category is an object representing the database table.
type Category struct {
	ID        string    `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	SkuID     string    `db:"sku_id" boil:"sku_id" json:"sku_id" toml:"sku_id" yaml:"sku_id"`
	Name      string    `db:"name" boil:"name" json:"name" toml:"name" yaml:"name"`
	CreatedAt time.Time `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *categoryR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L categoryL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CategoryColumns = struct {
	ID        string
	SkuID     string
	Name      string
	CreatedAt string
}{
	ID:        "id",
	SkuID:     "sku_id",
	Name:      "name",
	CreatedAt: "created_at",
}

// Generated where

var CategoryWhere = struct {
	ID        whereHelperstring
	SkuID     whereHelperstring
	Name      whereHelperstring
	CreatedAt whereHelpertime_Time
}{
	ID:        whereHelperstring{field: "\"categories\".\"id\""},
	SkuID:     whereHelperstring{field: "\"categories\".\"sku_id\""},
	Name:      whereHelperstring{field: "\"categories\".\"name\""},
	CreatedAt: whereHelpertime_Time{field: "\"categories\".\"created_at\""},
}

// CategoryRels is where relationship names are stored.
var CategoryRels = struct {
	Sku string
}{
	Sku: "Sku",
}

// categoryR is where relationships are stored.
type categoryR struct {
	Sku *StockKeepingUnit
}

// NewStruct creates a new relationship struct
func (*categoryR) NewStruct() *categoryR {
	return &categoryR{}
}

// categoryL is where Load methods for each relationship are stored.
type categoryL struct{}

var (
	categoryAllColumns            = []string{"id", "sku_id", "name", "created_at"}
	categoryColumnsWithoutDefault = []string{"sku_id", "name"}
	categoryColumnsWithDefault    = []string{"id", "created_at"}
	categoryPrimaryKeyColumns     = []string{"id"}
)

type (
	// CategorySlice is an alias for a slice of pointers to Category.
	// This should generally be used opposed to []Category.
	CategorySlice []*Category
	// CategoryHook is the signature for custom Category hook methods
	CategoryHook func(boil.Executor, *Category) error

	categoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	categoryType                 = reflect.TypeOf(&Category{})
	categoryMapping              = queries.MakeStructMapping(categoryType)
	categoryPrimaryKeyMapping, _ = queries.BindMapping(categoryType, categoryMapping, categoryPrimaryKeyColumns)
	categoryInsertCacheMut       sync.RWMutex
	categoryInsertCache          = make(map[string]insertCache)
	categoryUpdateCacheMut       sync.RWMutex
	categoryUpdateCache          = make(map[string]updateCache)
	categoryUpsertCacheMut       sync.RWMutex
	categoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var categoryBeforeInsertHooks []CategoryHook
var categoryBeforeUpdateHooks []CategoryHook
var categoryBeforeDeleteHooks []CategoryHook
var categoryBeforeUpsertHooks []CategoryHook

var categoryAfterInsertHooks []CategoryHook
var categoryAfterSelectHooks []CategoryHook
var categoryAfterUpdateHooks []CategoryHook
var categoryAfterDeleteHooks []CategoryHook
var categoryAfterUpsertHooks []CategoryHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Category) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Category) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Category) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Category) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Category) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Category) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Category) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Category) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Category) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range categoryAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCategoryHook registers your hook function for all future operations.
func AddCategoryHook(hookPoint boil.HookPoint, categoryHook CategoryHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		categoryBeforeInsertHooks = append(categoryBeforeInsertHooks, categoryHook)
	case boil.BeforeUpdateHook:
		categoryBeforeUpdateHooks = append(categoryBeforeUpdateHooks, categoryHook)
	case boil.BeforeDeleteHook:
		categoryBeforeDeleteHooks = append(categoryBeforeDeleteHooks, categoryHook)
	case boil.BeforeUpsertHook:
		categoryBeforeUpsertHooks = append(categoryBeforeUpsertHooks, categoryHook)
	case boil.AfterInsertHook:
		categoryAfterInsertHooks = append(categoryAfterInsertHooks, categoryHook)
	case boil.AfterSelectHook:
		categoryAfterSelectHooks = append(categoryAfterSelectHooks, categoryHook)
	case boil.AfterUpdateHook:
		categoryAfterUpdateHooks = append(categoryAfterUpdateHooks, categoryHook)
	case boil.AfterDeleteHook:
		categoryAfterDeleteHooks = append(categoryAfterDeleteHooks, categoryHook)
	case boil.AfterUpsertHook:
		categoryAfterUpsertHooks = append(categoryAfterUpsertHooks, categoryHook)
	}
}

// One returns a single category record from the query.
func (q categoryQuery) One(exec boil.Executor) (*Category, error) {
	o := &Category{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for categories")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Category records from the query.
func (q categoryQuery) All(exec boil.Executor) (CategorySlice, error) {
	var o []*Category

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to Category slice")
	}

	if len(categoryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Category records in the query.
func (q categoryQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count categories rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q categoryQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if categories exists")
	}

	return count > 0, nil
}

// Sku pointed to by the foreign key.
func (o *Category) Sku(mods ...qm.QueryMod) stockKeepingUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SkuID),
	}

	queryMods = append(queryMods, mods...)

	query := StockKeepingUnits(queryMods...)
	queries.SetFrom(query.Query, "\"stock_keeping_units\"")

	return query
}

// LoadSku allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (categoryL) LoadSku(e boil.Executor, singular bool, maybeCategory interface{}, mods queries.Applicator) error {
	var slice []*Category
	var object *Category

	if singular {
		object = maybeCategory.(*Category)
	} else {
		slice = *maybeCategory.(*[]*Category)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &categoryR{}
		}
		args = append(args, object.SkuID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &categoryR{}
			}

			for _, a := range args {
				if a == obj.SkuID {
					continue Outer
				}
			}

			args = append(args, obj.SkuID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`stock_keeping_units`), qm.WhereIn(`stock_keeping_units.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load StockKeepingUnit")
	}

	var resultSlice []*StockKeepingUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice StockKeepingUnit")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for stock_keeping_units")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for stock_keeping_units")
	}

	if len(categoryAfterSelectHooks) != 0 {
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
		object.R.Sku = foreign
		if foreign.R == nil {
			foreign.R = &stockKeepingUnitR{}
		}
		foreign.R.SkuCategories = append(foreign.R.SkuCategories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SkuID == foreign.ID {
				local.R.Sku = foreign
				if foreign.R == nil {
					foreign.R = &stockKeepingUnitR{}
				}
				foreign.R.SkuCategories = append(foreign.R.SkuCategories, local)
				break
			}
		}
	}

	return nil
}

// SetSku of the category to the related item.
// Sets o.R.Sku to related.
// Adds o to related.R.SkuCategories.
func (o *Category) SetSku(exec boil.Executor, insert bool, related *StockKeepingUnit) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sku_id"}),
		strmangle.WhereClause("\"", "\"", 2, categoryPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SkuID = related.ID
	if o.R == nil {
		o.R = &categoryR{
			Sku: related,
		}
	} else {
		o.R.Sku = related
	}

	if related.R == nil {
		related.R = &stockKeepingUnitR{
			SkuCategories: CategorySlice{o},
		}
	} else {
		related.R.SkuCategories = append(related.R.SkuCategories, o)
	}

	return nil
}

// Categories retrieves all the records using an executor.
func Categories(mods ...qm.QueryMod) categoryQuery {
	mods = append(mods, qm.From("\"categories\""))
	return categoryQuery{NewQuery(mods...)}
}

// FindCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCategory(exec boil.Executor, iD string, selectCols ...string) (*Category, error) {
	categoryObj := &Category{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"categories\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, categoryObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from categories")
	}

	return categoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Category) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no categories provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(categoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	categoryInsertCacheMut.RLock()
	cache, cached := categoryInsertCache[key]
	categoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			categoryAllColumns,
			categoryColumnsWithDefault,
			categoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(categoryType, categoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"categories\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"categories\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into categories")
	}

	if !cached {
		categoryInsertCacheMut.Lock()
		categoryInsertCache[key] = cache
		categoryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the Category.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Category) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	categoryUpdateCacheMut.RLock()
	cache, cached := categoryUpdateCache[key]
	categoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			categoryAllColumns,
			categoryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update categories, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"categories\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, categoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, append(wl, categoryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update categories row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for categories")
	}

	if !cached {
		categoryUpdateCacheMut.Lock()
		categoryUpdateCache[key] = cache
		categoryUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q categoryQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for categories")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CategorySlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, categoryPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in category slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all category")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Category) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no categories provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(categoryColumnsWithDefault, o)

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

	categoryUpsertCacheMut.RLock()
	cache, cached := categoryUpsertCache[key]
	categoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			categoryAllColumns,
			categoryColumnsWithDefault,
			categoryColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			categoryAllColumns,
			categoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert categories, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(categoryPrimaryKeyColumns))
			copy(conflict, categoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"categories\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(categoryType, categoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(categoryType, categoryMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert categories")
	}

	if !cached {
		categoryUpsertCacheMut.Lock()
		categoryUpsertCache[key] = cache
		categoryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single Category record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Category) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no Category provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), categoryPrimaryKeyMapping)
	sql := "DELETE FROM \"categories\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for categories")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q categoryQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no categoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for categories")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CategorySlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(categoryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, categoryPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from category slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for categories")
	}

	if len(categoryAfterDeleteHooks) != 0 {
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
func (o *Category) Reload(exec boil.Executor) error {
	ret, err := FindCategory(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CategorySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), categoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"categories\".* FROM \"categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, categoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in CategorySlice")
	}

	*o = slice

	return nil
}

// CategoryExists checks if the Category row exists.
func CategoryExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"categories\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if categories exists")
	}

	return exists, nil
}
