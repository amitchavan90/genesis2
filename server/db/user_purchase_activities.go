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

// UserPurchaseActivity is an object representing the database table.
type UserPurchaseActivity struct {
	ID              string      `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID          string      `db:"user_id" boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	ProductID       null.String `db:"product_id" boil:"product_id" json:"product_id,omitempty" toml:"product_id" yaml:"product_id,omitempty"`
	LoyaltyPoints   int         `db:"loyalty_points" boil:"loyalty_points" json:"loyalty_points" toml:"loyalty_points" yaml:"loyalty_points"`
	Message         string      `db:"message" boil:"message" json:"message" toml:"message" yaml:"message"`
	TransactionHash null.String `db:"transaction_hash" boil:"transaction_hash" json:"transaction_hash,omitempty" toml:"transaction_hash" yaml:"transaction_hash,omitempty"`
	CreatedAt       time.Time   `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *userPurchaseActivityR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userPurchaseActivityL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserPurchaseActivityColumns = struct {
	ID              string
	UserID          string
	ProductID       string
	LoyaltyPoints   string
	Message         string
	TransactionHash string
	CreatedAt       string
}{
	ID:              "id",
	UserID:          "user_id",
	ProductID:       "product_id",
	LoyaltyPoints:   "loyalty_points",
	Message:         "message",
	TransactionHash: "transaction_hash",
	CreatedAt:       "created_at",
}

// Generated where

var UserPurchaseActivityWhere = struct {
	ID              whereHelperstring
	UserID          whereHelperstring
	ProductID       whereHelpernull_String
	LoyaltyPoints   whereHelperint
	Message         whereHelperstring
	TransactionHash whereHelpernull_String
	CreatedAt       whereHelpertime_Time
}{
	ID:              whereHelperstring{field: "\"user_purchase_activities\".\"id\""},
	UserID:          whereHelperstring{field: "\"user_purchase_activities\".\"user_id\""},
	ProductID:       whereHelpernull_String{field: "\"user_purchase_activities\".\"product_id\""},
	LoyaltyPoints:   whereHelperint{field: "\"user_purchase_activities\".\"loyalty_points\""},
	Message:         whereHelperstring{field: "\"user_purchase_activities\".\"message\""},
	TransactionHash: whereHelpernull_String{field: "\"user_purchase_activities\".\"transaction_hash\""},
	CreatedAt:       whereHelpertime_Time{field: "\"user_purchase_activities\".\"created_at\""},
}

// UserPurchaseActivityRels is where relationship names are stored.
var UserPurchaseActivityRels = struct {
	Product string
	User    string
}{
	Product: "Product",
	User:    "User",
}

// userPurchaseActivityR is where relationships are stored.
type userPurchaseActivityR struct {
	Product *Product
	User    *User
}

// NewStruct creates a new relationship struct
func (*userPurchaseActivityR) NewStruct() *userPurchaseActivityR {
	return &userPurchaseActivityR{}
}

// userPurchaseActivityL is where Load methods for each relationship are stored.
type userPurchaseActivityL struct{}

var (
	userPurchaseActivityAllColumns            = []string{"id", "user_id", "product_id", "loyalty_points", "message", "transaction_hash", "created_at"}
	userPurchaseActivityColumnsWithoutDefault = []string{"user_id", "product_id", "loyalty_points", "transaction_hash"}
	userPurchaseActivityColumnsWithDefault    = []string{"id", "message", "created_at"}
	userPurchaseActivityPrimaryKeyColumns     = []string{"id"}
)

type (
	// UserPurchaseActivitySlice is an alias for a slice of pointers to UserPurchaseActivity.
	// This should generally be used opposed to []UserPurchaseActivity.
	UserPurchaseActivitySlice []*UserPurchaseActivity
	// UserPurchaseActivityHook is the signature for custom UserPurchaseActivity hook methods
	UserPurchaseActivityHook func(boil.Executor, *UserPurchaseActivity) error

	userPurchaseActivityQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userPurchaseActivityType                 = reflect.TypeOf(&UserPurchaseActivity{})
	userPurchaseActivityMapping              = queries.MakeStructMapping(userPurchaseActivityType)
	userPurchaseActivityPrimaryKeyMapping, _ = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, userPurchaseActivityPrimaryKeyColumns)
	userPurchaseActivityInsertCacheMut       sync.RWMutex
	userPurchaseActivityInsertCache          = make(map[string]insertCache)
	userPurchaseActivityUpdateCacheMut       sync.RWMutex
	userPurchaseActivityUpdateCache          = make(map[string]updateCache)
	userPurchaseActivityUpsertCacheMut       sync.RWMutex
	userPurchaseActivityUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userPurchaseActivityBeforeInsertHooks []UserPurchaseActivityHook
var userPurchaseActivityBeforeUpdateHooks []UserPurchaseActivityHook
var userPurchaseActivityBeforeDeleteHooks []UserPurchaseActivityHook
var userPurchaseActivityBeforeUpsertHooks []UserPurchaseActivityHook

var userPurchaseActivityAfterInsertHooks []UserPurchaseActivityHook
var userPurchaseActivityAfterSelectHooks []UserPurchaseActivityHook
var userPurchaseActivityAfterUpdateHooks []UserPurchaseActivityHook
var userPurchaseActivityAfterDeleteHooks []UserPurchaseActivityHook
var userPurchaseActivityAfterUpsertHooks []UserPurchaseActivityHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserPurchaseActivity) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserPurchaseActivity) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserPurchaseActivity) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserPurchaseActivity) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserPurchaseActivity) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserPurchaseActivity) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserPurchaseActivity) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserPurchaseActivity) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserPurchaseActivity) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userPurchaseActivityAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserPurchaseActivityHook registers your hook function for all future operations.
func AddUserPurchaseActivityHook(hookPoint boil.HookPoint, userPurchaseActivityHook UserPurchaseActivityHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userPurchaseActivityBeforeInsertHooks = append(userPurchaseActivityBeforeInsertHooks, userPurchaseActivityHook)
	case boil.BeforeUpdateHook:
		userPurchaseActivityBeforeUpdateHooks = append(userPurchaseActivityBeforeUpdateHooks, userPurchaseActivityHook)
	case boil.BeforeDeleteHook:
		userPurchaseActivityBeforeDeleteHooks = append(userPurchaseActivityBeforeDeleteHooks, userPurchaseActivityHook)
	case boil.BeforeUpsertHook:
		userPurchaseActivityBeforeUpsertHooks = append(userPurchaseActivityBeforeUpsertHooks, userPurchaseActivityHook)
	case boil.AfterInsertHook:
		userPurchaseActivityAfterInsertHooks = append(userPurchaseActivityAfterInsertHooks, userPurchaseActivityHook)
	case boil.AfterSelectHook:
		userPurchaseActivityAfterSelectHooks = append(userPurchaseActivityAfterSelectHooks, userPurchaseActivityHook)
	case boil.AfterUpdateHook:
		userPurchaseActivityAfterUpdateHooks = append(userPurchaseActivityAfterUpdateHooks, userPurchaseActivityHook)
	case boil.AfterDeleteHook:
		userPurchaseActivityAfterDeleteHooks = append(userPurchaseActivityAfterDeleteHooks, userPurchaseActivityHook)
	case boil.AfterUpsertHook:
		userPurchaseActivityAfterUpsertHooks = append(userPurchaseActivityAfterUpsertHooks, userPurchaseActivityHook)
	}
}

// One returns a single userPurchaseActivity record from the query.
func (q userPurchaseActivityQuery) One(exec boil.Executor) (*UserPurchaseActivity, error) {
	o := &UserPurchaseActivity{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for user_purchase_activities")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserPurchaseActivity records from the query.
func (q userPurchaseActivityQuery) All(exec boil.Executor) (UserPurchaseActivitySlice, error) {
	var o []*UserPurchaseActivity

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to UserPurchaseActivity slice")
	}

	if len(userPurchaseActivityAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserPurchaseActivity records in the query.
func (q userPurchaseActivityQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count user_purchase_activities rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userPurchaseActivityQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if user_purchase_activities exists")
	}

	return count > 0, nil
}

// Product pointed to by the foreign key.
func (o *UserPurchaseActivity) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"products\"")

	return query
}

// User pointed to by the foreign key.
func (o *UserPurchaseActivity) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userPurchaseActivityL) LoadProduct(e boil.Executor, singular bool, maybeUserPurchaseActivity interface{}, mods queries.Applicator) error {
	var slice []*UserPurchaseActivity
	var object *UserPurchaseActivity

	if singular {
		object = maybeUserPurchaseActivity.(*UserPurchaseActivity)
	} else {
		slice = *maybeUserPurchaseActivity.(*[]*UserPurchaseActivity)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userPurchaseActivityR{}
		}
		if !queries.IsNil(object.ProductID) {
			args = append(args, object.ProductID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userPurchaseActivityR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ProductID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.ProductID) {
				args = append(args, obj.ProductID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`products`), qm.WhereIn(`products.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Product")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Product")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for products")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for products")
	}

	if len(userPurchaseActivityAfterSelectHooks) != 0 {
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
		object.R.Product = foreign
		if foreign.R == nil {
			foreign.R = &productR{}
		}
		foreign.R.UserPurchaseActivities = append(foreign.R.UserPurchaseActivities, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.ProductID, foreign.ID) {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.UserPurchaseActivities = append(foreign.R.UserPurchaseActivities, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userPurchaseActivityL) LoadUser(e boil.Executor, singular bool, maybeUserPurchaseActivity interface{}, mods queries.Applicator) error {
	var slice []*UserPurchaseActivity
	var object *UserPurchaseActivity

	if singular {
		object = maybeUserPurchaseActivity.(*UserPurchaseActivity)
	} else {
		slice = *maybeUserPurchaseActivity.(*[]*UserPurchaseActivity)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userPurchaseActivityR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userPurchaseActivityR{}
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

	if len(userPurchaseActivityAfterSelectHooks) != 0 {
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
		foreign.R.UserPurchaseActivities = append(foreign.R.UserPurchaseActivities, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.UserPurchaseActivities = append(foreign.R.UserPurchaseActivities, local)
				break
			}
		}
	}

	return nil
}

// SetProduct of the userPurchaseActivity to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.UserPurchaseActivities.
func (o *UserPurchaseActivity) SetProduct(exec boil.Executor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_purchase_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, userPurchaseActivityPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.ProductID, related.ID)
	if o.R == nil {
		o.R = &userPurchaseActivityR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			UserPurchaseActivities: UserPurchaseActivitySlice{o},
		}
	} else {
		related.R.UserPurchaseActivities = append(related.R.UserPurchaseActivities, o)
	}

	return nil
}

// RemoveProduct relationship.
// Sets o.R.Product to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *UserPurchaseActivity) RemoveProduct(exec boil.Executor, related *Product) error {
	var err error

	queries.SetScanner(&o.ProductID, nil)
	if _, err = o.Update(exec, boil.Whitelist("product_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Product = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.UserPurchaseActivities {
		if queries.Equal(o.ProductID, ri.ProductID) {
			continue
		}

		ln := len(related.R.UserPurchaseActivities)
		if ln > 1 && i < ln-1 {
			related.R.UserPurchaseActivities[i] = related.R.UserPurchaseActivities[ln-1]
		}
		related.R.UserPurchaseActivities = related.R.UserPurchaseActivities[:ln-1]
		break
	}
	return nil
}

// SetUser of the userPurchaseActivity to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserPurchaseActivities.
func (o *UserPurchaseActivity) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_purchase_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userPurchaseActivityPrimaryKeyColumns),
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
		o.R = &userPurchaseActivityR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			UserPurchaseActivities: UserPurchaseActivitySlice{o},
		}
	} else {
		related.R.UserPurchaseActivities = append(related.R.UserPurchaseActivities, o)
	}

	return nil
}

// UserPurchaseActivities retrieves all the records using an executor.
func UserPurchaseActivities(mods ...qm.QueryMod) userPurchaseActivityQuery {
	mods = append(mods, qm.From("\"user_purchase_activities\""))
	return userPurchaseActivityQuery{NewQuery(mods...)}
}

// FindUserPurchaseActivity retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserPurchaseActivity(exec boil.Executor, iD string, selectCols ...string) (*UserPurchaseActivity, error) {
	userPurchaseActivityObj := &UserPurchaseActivity{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_purchase_activities\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, userPurchaseActivityObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from user_purchase_activities")
	}

	return userPurchaseActivityObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserPurchaseActivity) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_purchase_activities provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userPurchaseActivityColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userPurchaseActivityInsertCacheMut.RLock()
	cache, cached := userPurchaseActivityInsertCache[key]
	userPurchaseActivityInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userPurchaseActivityAllColumns,
			userPurchaseActivityColumnsWithDefault,
			userPurchaseActivityColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_purchase_activities\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_purchase_activities\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into user_purchase_activities")
	}

	if !cached {
		userPurchaseActivityInsertCacheMut.Lock()
		userPurchaseActivityInsertCache[key] = cache
		userPurchaseActivityInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the UserPurchaseActivity.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserPurchaseActivity) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userPurchaseActivityUpdateCacheMut.RLock()
	cache, cached := userPurchaseActivityUpdateCache[key]
	userPurchaseActivityUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userPurchaseActivityAllColumns,
			userPurchaseActivityPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update user_purchase_activities, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_purchase_activities\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userPurchaseActivityPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, append(wl, userPurchaseActivityPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update user_purchase_activities row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for user_purchase_activities")
	}

	if !cached {
		userPurchaseActivityUpdateCacheMut.Lock()
		userPurchaseActivityUpdateCache[key] = cache
		userPurchaseActivityUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userPurchaseActivityQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for user_purchase_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for user_purchase_activities")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserPurchaseActivitySlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPurchaseActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_purchase_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userPurchaseActivityPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in userPurchaseActivity slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all userPurchaseActivity")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserPurchaseActivity) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_purchase_activities provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userPurchaseActivityColumnsWithDefault, o)

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

	userPurchaseActivityUpsertCacheMut.RLock()
	cache, cached := userPurchaseActivityUpsertCache[key]
	userPurchaseActivityUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userPurchaseActivityAllColumns,
			userPurchaseActivityColumnsWithDefault,
			userPurchaseActivityColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userPurchaseActivityAllColumns,
			userPurchaseActivityPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert user_purchase_activities, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userPurchaseActivityPrimaryKeyColumns))
			copy(conflict, userPurchaseActivityPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_purchase_activities\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userPurchaseActivityType, userPurchaseActivityMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert user_purchase_activities")
	}

	if !cached {
		userPurchaseActivityUpsertCacheMut.Lock()
		userPurchaseActivityUpsertCache[key] = cache
		userPurchaseActivityUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single UserPurchaseActivity record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserPurchaseActivity) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no UserPurchaseActivity provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userPurchaseActivityPrimaryKeyMapping)
	sql := "DELETE FROM \"user_purchase_activities\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from user_purchase_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for user_purchase_activities")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userPurchaseActivityQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no userPurchaseActivityQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from user_purchase_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_purchase_activities")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserPurchaseActivitySlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userPurchaseActivityBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPurchaseActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_purchase_activities\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPurchaseActivityPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from userPurchaseActivity slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_purchase_activities")
	}

	if len(userPurchaseActivityAfterDeleteHooks) != 0 {
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
func (o *UserPurchaseActivity) Reload(exec boil.Executor) error {
	ret, err := FindUserPurchaseActivity(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserPurchaseActivitySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserPurchaseActivitySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPurchaseActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_purchase_activities\".* FROM \"user_purchase_activities\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPurchaseActivityPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in UserPurchaseActivitySlice")
	}

	*o = slice

	return nil
}

// UserPurchaseActivityExists checks if the UserPurchaseActivity row exists.
func UserPurchaseActivityExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_purchase_activities\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if user_purchase_activities exists")
	}

	return exists, nil
}
