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

// Order is an object representing the database table.
type Order struct {
	ID          string    `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	Code        string    `db:"code" boil:"code" json:"code" toml:"code" yaml:"code"`
	Archived    bool      `db:"archived" boil:"archived" json:"archived" toml:"archived" yaml:"archived"`
	ArchivedAt  null.Time `db:"archived_at" boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedAt   time.Time `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	CreatedByID string    `db:"created_by_id" boil:"created_by_id" json:"created_by_id" toml:"created_by_id" yaml:"created_by_id"`

	R *orderR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L orderL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OrderColumns = struct {
	ID          string
	Code        string
	Archived    string
	ArchivedAt  string
	UpdatedAt   string
	CreatedAt   string
	CreatedByID string
}{
	ID:          "id",
	Code:        "code",
	Archived:    "archived",
	ArchivedAt:  "archived_at",
	UpdatedAt:   "updated_at",
	CreatedAt:   "created_at",
	CreatedByID: "created_by_id",
}

// Generated where

var OrderWhere = struct {
	ID          whereHelperstring
	Code        whereHelperstring
	Archived    whereHelperbool
	ArchivedAt  whereHelpernull_Time
	UpdatedAt   whereHelpertime_Time
	CreatedAt   whereHelpertime_Time
	CreatedByID whereHelperstring
}{
	ID:          whereHelperstring{field: "\"orders\".\"id\""},
	Code:        whereHelperstring{field: "\"orders\".\"code\""},
	Archived:    whereHelperbool{field: "\"orders\".\"archived\""},
	ArchivedAt:  whereHelpernull_Time{field: "\"orders\".\"archived_at\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"orders\".\"updated_at\""},
	CreatedAt:   whereHelpertime_Time{field: "\"orders\".\"created_at\""},
	CreatedByID: whereHelperstring{field: "\"orders\".\"created_by_id\""},
}

// OrderRels is where relationship names are stored.
var OrderRels = struct {
	CreatedBy string
	Products  string
}{
	CreatedBy: "CreatedBy",
	Products:  "Products",
}

// orderR is where relationships are stored.
type orderR struct {
	CreatedBy *User
	Products  ProductSlice
}

// NewStruct creates a new relationship struct
func (*orderR) NewStruct() *orderR {
	return &orderR{}
}

// orderL is where Load methods for each relationship are stored.
type orderL struct{}

var (
	orderAllColumns            = []string{"id", "code", "archived", "archived_at", "updated_at", "created_at", "created_by_id"}
	orderColumnsWithoutDefault = []string{"code", "archived_at", "created_by_id"}
	orderColumnsWithDefault    = []string{"id", "archived", "updated_at", "created_at"}
	orderPrimaryKeyColumns     = []string{"id"}
)

type (
	// OrderSlice is an alias for a slice of pointers to Order.
	// This should generally be used opposed to []Order.
	OrderSlice []*Order
	// OrderHook is the signature for custom Order hook methods
	OrderHook func(boil.Executor, *Order) error

	orderQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	orderType                 = reflect.TypeOf(&Order{})
	orderMapping              = queries.MakeStructMapping(orderType)
	orderPrimaryKeyMapping, _ = queries.BindMapping(orderType, orderMapping, orderPrimaryKeyColumns)
	orderInsertCacheMut       sync.RWMutex
	orderInsertCache          = make(map[string]insertCache)
	orderUpdateCacheMut       sync.RWMutex
	orderUpdateCache          = make(map[string]updateCache)
	orderUpsertCacheMut       sync.RWMutex
	orderUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var orderBeforeInsertHooks []OrderHook
var orderBeforeUpdateHooks []OrderHook
var orderBeforeDeleteHooks []OrderHook
var orderBeforeUpsertHooks []OrderHook

var orderAfterInsertHooks []OrderHook
var orderAfterSelectHooks []OrderHook
var orderAfterUpdateHooks []OrderHook
var orderAfterDeleteHooks []OrderHook
var orderAfterUpsertHooks []OrderHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Order) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range orderBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Order) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range orderBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Order) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range orderBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Order) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range orderBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Order) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range orderAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Order) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range orderAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Order) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range orderAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Order) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range orderAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Order) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range orderAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOrderHook registers your hook function for all future operations.
func AddOrderHook(hookPoint boil.HookPoint, orderHook OrderHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		orderBeforeInsertHooks = append(orderBeforeInsertHooks, orderHook)
	case boil.BeforeUpdateHook:
		orderBeforeUpdateHooks = append(orderBeforeUpdateHooks, orderHook)
	case boil.BeforeDeleteHook:
		orderBeforeDeleteHooks = append(orderBeforeDeleteHooks, orderHook)
	case boil.BeforeUpsertHook:
		orderBeforeUpsertHooks = append(orderBeforeUpsertHooks, orderHook)
	case boil.AfterInsertHook:
		orderAfterInsertHooks = append(orderAfterInsertHooks, orderHook)
	case boil.AfterSelectHook:
		orderAfterSelectHooks = append(orderAfterSelectHooks, orderHook)
	case boil.AfterUpdateHook:
		orderAfterUpdateHooks = append(orderAfterUpdateHooks, orderHook)
	case boil.AfterDeleteHook:
		orderAfterDeleteHooks = append(orderAfterDeleteHooks, orderHook)
	case boil.AfterUpsertHook:
		orderAfterUpsertHooks = append(orderAfterUpsertHooks, orderHook)
	}
}

// One returns a single order record from the query.
func (q orderQuery) One(exec boil.Executor) (*Order, error) {
	o := &Order{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for orders")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Order records from the query.
func (q orderQuery) All(exec boil.Executor) (OrderSlice, error) {
	var o []*Order

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to Order slice")
	}

	if len(orderAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Order records in the query.
func (q orderQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count orders rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q orderQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if orders exists")
	}

	return count > 0, nil
}

// CreatedBy pointed to by the foreign key.
func (o *Order) CreatedBy(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.CreatedByID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// Products retrieves all the product's Products with an executor.
func (o *Order) Products(mods ...qm.QueryMod) productQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"products\".\"order_id\"=?", o.ID),
	)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"products\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"products\".*"})
	}

	return query
}

// LoadCreatedBy allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (orderL) LoadCreatedBy(e boil.Executor, singular bool, maybeOrder interface{}, mods queries.Applicator) error {
	var slice []*Order
	var object *Order

	if singular {
		object = maybeOrder.(*Order)
	} else {
		slice = *maybeOrder.(*[]*Order)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &orderR{}
		}
		args = append(args, object.CreatedByID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &orderR{}
			}

			for _, a := range args {
				if a == obj.CreatedByID {
					continue Outer
				}
			}

			args = append(args, obj.CreatedByID)

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

	if len(orderAfterSelectHooks) != 0 {
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
		object.R.CreatedBy = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.CreatedByOrders = append(foreign.R.CreatedByOrders, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CreatedByID == foreign.ID {
				local.R.CreatedBy = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.CreatedByOrders = append(foreign.R.CreatedByOrders, local)
				break
			}
		}
	}

	return nil
}

// LoadProducts allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (orderL) LoadProducts(e boil.Executor, singular bool, maybeOrder interface{}, mods queries.Applicator) error {
	var slice []*Order
	var object *Order

	if singular {
		object = maybeOrder.(*Order)
	} else {
		slice = *maybeOrder.(*[]*Order)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &orderR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &orderR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`products`), qm.WhereIn(`products.order_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load products")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice products")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on products")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for products")
	}

	if len(productAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Products = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &productR{}
			}
			foreign.R.Order = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.OrderID) {
				local.R.Products = append(local.R.Products, foreign)
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.Order = local
				break
			}
		}
	}

	return nil
}

// SetCreatedBy of the order to the related item.
// Sets o.R.CreatedBy to related.
// Adds o to related.R.CreatedByOrders.
func (o *Order) SetCreatedBy(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"orders\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"created_by_id"}),
		strmangle.WhereClause("\"", "\"", 2, orderPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CreatedByID = related.ID
	if o.R == nil {
		o.R = &orderR{
			CreatedBy: related,
		}
	} else {
		o.R.CreatedBy = related
	}

	if related.R == nil {
		related.R = &userR{
			CreatedByOrders: OrderSlice{o},
		}
	} else {
		related.R.CreatedByOrders = append(related.R.CreatedByOrders, o)
	}

	return nil
}

// AddProducts adds the given related objects to the existing relationships
// of the order, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Order appropriately.
func (o *Order) AddProducts(exec boil.Executor, insert bool, related ...*Product) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.OrderID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"products\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"order_id"}),
				strmangle.WhereClause("\"", "\"", 2, productPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.OrderID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &orderR{
			Products: related,
		}
	} else {
		o.R.Products = append(o.R.Products, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &productR{
				Order: o,
			}
		} else {
			rel.R.Order = o
		}
	}
	return nil
}

// SetProducts removes all previously related items of the
// order replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Order's Products accordingly.
// Replaces o.R.Products with related.
// Sets related.R.Order's Products accordingly.
func (o *Order) SetProducts(exec boil.Executor, insert bool, related ...*Product) error {
	query := "update \"products\" set \"order_id\" = null where \"order_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Products {
			queries.SetScanner(&rel.OrderID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Order = nil
		}

		o.R.Products = nil
	}
	return o.AddProducts(exec, insert, related...)
}

// RemoveProducts relationships from objects passed in.
// Removes related items from R.Products (uses pointer comparison, removal does not keep order)
// Sets related.R.Order.
func (o *Order) RemoveProducts(exec boil.Executor, related ...*Product) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.OrderID, nil)
		if rel.R != nil {
			rel.R.Order = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("order_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Products {
			if rel != ri {
				continue
			}

			ln := len(o.R.Products)
			if ln > 1 && i < ln-1 {
				o.R.Products[i] = o.R.Products[ln-1]
			}
			o.R.Products = o.R.Products[:ln-1]
			break
		}
	}

	return nil
}

// Orders retrieves all the records using an executor.
func Orders(mods ...qm.QueryMod) orderQuery {
	mods = append(mods, qm.From("\"orders\""))
	return orderQuery{NewQuery(mods...)}
}

// FindOrder retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOrder(exec boil.Executor, iD string, selectCols ...string) (*Order, error) {
	orderObj := &Order{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"orders\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, orderObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from orders")
	}

	return orderObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Order) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no orders provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(orderColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	orderInsertCacheMut.RLock()
	cache, cached := orderInsertCache[key]
	orderInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			orderAllColumns,
			orderColumnsWithDefault,
			orderColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(orderType, orderMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(orderType, orderMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"orders\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"orders\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into orders")
	}

	if !cached {
		orderInsertCacheMut.Lock()
		orderInsertCache[key] = cache
		orderInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the Order.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Order) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	orderUpdateCacheMut.RLock()
	cache, cached := orderUpdateCache[key]
	orderUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			orderAllColumns,
			orderPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update orders, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"orders\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, orderPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(orderType, orderMapping, append(wl, orderPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update orders row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for orders")
	}

	if !cached {
		orderUpdateCacheMut.Lock()
		orderUpdateCache[key] = cache
		orderUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q orderQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for orders")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for orders")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OrderSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"orders\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, orderPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in order slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all order")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Order) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no orders provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(orderColumnsWithDefault, o)

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

	orderUpsertCacheMut.RLock()
	cache, cached := orderUpsertCache[key]
	orderUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			orderAllColumns,
			orderColumnsWithDefault,
			orderColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			orderAllColumns,
			orderPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert orders, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(orderPrimaryKeyColumns))
			copy(conflict, orderPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"orders\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(orderType, orderMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(orderType, orderMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert orders")
	}

	if !cached {
		orderUpsertCacheMut.Lock()
		orderUpsertCache[key] = cache
		orderUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single Order record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Order) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no Order provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), orderPrimaryKeyMapping)
	sql := "DELETE FROM \"orders\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from orders")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for orders")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q orderQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no orderQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from orders")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for orders")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OrderSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(orderBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"orders\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, orderPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from order slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for orders")
	}

	if len(orderAfterDeleteHooks) != 0 {
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
func (o *Order) Reload(exec boil.Executor) error {
	ret, err := FindOrder(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OrderSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OrderSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"orders\".* FROM \"orders\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, orderPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in OrderSlice")
	}

	*o = slice

	return nil
}

// OrderExists checks if the Order row exists.
func OrderExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"orders\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if orders exists")
	}

	return exists, nil
}
