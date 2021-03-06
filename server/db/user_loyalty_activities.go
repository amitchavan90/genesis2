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

// UserLoyaltyActivity is an object representing the database table.
type UserLoyaltyActivity struct {
	ID              string      `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID          string      `db:"user_id" boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	ProductID       null.String `db:"product_id" boil:"product_id" json:"product_id,omitempty" toml:"product_id" yaml:"product_id,omitempty"`
	Amount          int         `db:"amount" boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Bonus           int         `db:"bonus" boil:"bonus" json:"bonus" toml:"bonus" yaml:"bonus"`
	Message         string      `db:"message" boil:"message" json:"message" toml:"message" yaml:"message"`
	TransactionHash null.String `db:"transaction_hash" boil:"transaction_hash" json:"transaction_hash,omitempty" toml:"transaction_hash" yaml:"transaction_hash,omitempty"`
	CreatedAt       time.Time   `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *userLoyaltyActivityR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userLoyaltyActivityL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserLoyaltyActivityColumns = struct {
	ID              string
	UserID          string
	ProductID       string
	Amount          string
	Bonus           string
	Message         string
	TransactionHash string
	CreatedAt       string
}{
	ID:              "id",
	UserID:          "user_id",
	ProductID:       "product_id",
	Amount:          "amount",
	Bonus:           "bonus",
	Message:         "message",
	TransactionHash: "transaction_hash",
	CreatedAt:       "created_at",
}

// Generated where

var UserLoyaltyActivityWhere = struct {
	ID              whereHelperstring
	UserID          whereHelperstring
	ProductID       whereHelpernull_String
	Amount          whereHelperint
	Bonus           whereHelperint
	Message         whereHelperstring
	TransactionHash whereHelpernull_String
	CreatedAt       whereHelpertime_Time
}{
	ID:              whereHelperstring{field: "\"user_loyalty_activities\".\"id\""},
	UserID:          whereHelperstring{field: "\"user_loyalty_activities\".\"user_id\""},
	ProductID:       whereHelpernull_String{field: "\"user_loyalty_activities\".\"product_id\""},
	Amount:          whereHelperint{field: "\"user_loyalty_activities\".\"amount\""},
	Bonus:           whereHelperint{field: "\"user_loyalty_activities\".\"bonus\""},
	Message:         whereHelperstring{field: "\"user_loyalty_activities\".\"message\""},
	TransactionHash: whereHelpernull_String{field: "\"user_loyalty_activities\".\"transaction_hash\""},
	CreatedAt:       whereHelpertime_Time{field: "\"user_loyalty_activities\".\"created_at\""},
}

// UserLoyaltyActivityRels is where relationship names are stored.
var UserLoyaltyActivityRels = struct {
	Product string
	User    string
}{
	Product: "Product",
	User:    "User",
}

// userLoyaltyActivityR is where relationships are stored.
type userLoyaltyActivityR struct {
	Product *Product
	User    *User
}

// NewStruct creates a new relationship struct
func (*userLoyaltyActivityR) NewStruct() *userLoyaltyActivityR {
	return &userLoyaltyActivityR{}
}

// userLoyaltyActivityL is where Load methods for each relationship are stored.
type userLoyaltyActivityL struct{}

var (
	userLoyaltyActivityAllColumns            = []string{"id", "user_id", "product_id", "amount", "bonus", "message", "transaction_hash", "created_at"}
	userLoyaltyActivityColumnsWithoutDefault = []string{"user_id", "product_id", "amount", "transaction_hash"}
	userLoyaltyActivityColumnsWithDefault    = []string{"id", "bonus", "message", "created_at"}
	userLoyaltyActivityPrimaryKeyColumns     = []string{"id"}
)

type (
	// UserLoyaltyActivitySlice is an alias for a slice of pointers to UserLoyaltyActivity.
	// This should generally be used opposed to []UserLoyaltyActivity.
	UserLoyaltyActivitySlice []*UserLoyaltyActivity
	// UserLoyaltyActivityHook is the signature for custom UserLoyaltyActivity hook methods
	UserLoyaltyActivityHook func(boil.Executor, *UserLoyaltyActivity) error

	userLoyaltyActivityQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userLoyaltyActivityType                 = reflect.TypeOf(&UserLoyaltyActivity{})
	userLoyaltyActivityMapping              = queries.MakeStructMapping(userLoyaltyActivityType)
	userLoyaltyActivityPrimaryKeyMapping, _ = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, userLoyaltyActivityPrimaryKeyColumns)
	userLoyaltyActivityInsertCacheMut       sync.RWMutex
	userLoyaltyActivityInsertCache          = make(map[string]insertCache)
	userLoyaltyActivityUpdateCacheMut       sync.RWMutex
	userLoyaltyActivityUpdateCache          = make(map[string]updateCache)
	userLoyaltyActivityUpsertCacheMut       sync.RWMutex
	userLoyaltyActivityUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userLoyaltyActivityBeforeInsertHooks []UserLoyaltyActivityHook
var userLoyaltyActivityBeforeUpdateHooks []UserLoyaltyActivityHook
var userLoyaltyActivityBeforeDeleteHooks []UserLoyaltyActivityHook
var userLoyaltyActivityBeforeUpsertHooks []UserLoyaltyActivityHook

var userLoyaltyActivityAfterInsertHooks []UserLoyaltyActivityHook
var userLoyaltyActivityAfterSelectHooks []UserLoyaltyActivityHook
var userLoyaltyActivityAfterUpdateHooks []UserLoyaltyActivityHook
var userLoyaltyActivityAfterDeleteHooks []UserLoyaltyActivityHook
var userLoyaltyActivityAfterUpsertHooks []UserLoyaltyActivityHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserLoyaltyActivity) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserLoyaltyActivity) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserLoyaltyActivity) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserLoyaltyActivity) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserLoyaltyActivity) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserLoyaltyActivity) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserLoyaltyActivity) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserLoyaltyActivity) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserLoyaltyActivity) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userLoyaltyActivityAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserLoyaltyActivityHook registers your hook function for all future operations.
func AddUserLoyaltyActivityHook(hookPoint boil.HookPoint, userLoyaltyActivityHook UserLoyaltyActivityHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userLoyaltyActivityBeforeInsertHooks = append(userLoyaltyActivityBeforeInsertHooks, userLoyaltyActivityHook)
	case boil.BeforeUpdateHook:
		userLoyaltyActivityBeforeUpdateHooks = append(userLoyaltyActivityBeforeUpdateHooks, userLoyaltyActivityHook)
	case boil.BeforeDeleteHook:
		userLoyaltyActivityBeforeDeleteHooks = append(userLoyaltyActivityBeforeDeleteHooks, userLoyaltyActivityHook)
	case boil.BeforeUpsertHook:
		userLoyaltyActivityBeforeUpsertHooks = append(userLoyaltyActivityBeforeUpsertHooks, userLoyaltyActivityHook)
	case boil.AfterInsertHook:
		userLoyaltyActivityAfterInsertHooks = append(userLoyaltyActivityAfterInsertHooks, userLoyaltyActivityHook)
	case boil.AfterSelectHook:
		userLoyaltyActivityAfterSelectHooks = append(userLoyaltyActivityAfterSelectHooks, userLoyaltyActivityHook)
	case boil.AfterUpdateHook:
		userLoyaltyActivityAfterUpdateHooks = append(userLoyaltyActivityAfterUpdateHooks, userLoyaltyActivityHook)
	case boil.AfterDeleteHook:
		userLoyaltyActivityAfterDeleteHooks = append(userLoyaltyActivityAfterDeleteHooks, userLoyaltyActivityHook)
	case boil.AfterUpsertHook:
		userLoyaltyActivityAfterUpsertHooks = append(userLoyaltyActivityAfterUpsertHooks, userLoyaltyActivityHook)
	}
}

// One returns a single userLoyaltyActivity record from the query.
func (q userLoyaltyActivityQuery) One(exec boil.Executor) (*UserLoyaltyActivity, error) {
	o := &UserLoyaltyActivity{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for user_loyalty_activities")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserLoyaltyActivity records from the query.
func (q userLoyaltyActivityQuery) All(exec boil.Executor) (UserLoyaltyActivitySlice, error) {
	var o []*UserLoyaltyActivity

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to UserLoyaltyActivity slice")
	}

	if len(userLoyaltyActivityAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserLoyaltyActivity records in the query.
func (q userLoyaltyActivityQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count user_loyalty_activities rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userLoyaltyActivityQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if user_loyalty_activities exists")
	}

	return count > 0, nil
}

// Product pointed to by the foreign key.
func (o *UserLoyaltyActivity) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"products\"")

	return query
}

// User pointed to by the foreign key.
func (o *UserLoyaltyActivity) User(mods ...qm.QueryMod) userQuery {
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
func (userLoyaltyActivityL) LoadProduct(e boil.Executor, singular bool, maybeUserLoyaltyActivity interface{}, mods queries.Applicator) error {
	var slice []*UserLoyaltyActivity
	var object *UserLoyaltyActivity

	if singular {
		object = maybeUserLoyaltyActivity.(*UserLoyaltyActivity)
	} else {
		slice = *maybeUserLoyaltyActivity.(*[]*UserLoyaltyActivity)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userLoyaltyActivityR{}
		}
		if !queries.IsNil(object.ProductID) {
			args = append(args, object.ProductID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userLoyaltyActivityR{}
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

	if len(userLoyaltyActivityAfterSelectHooks) != 0 {
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
		foreign.R.UserLoyaltyActivities = append(foreign.R.UserLoyaltyActivities, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.ProductID, foreign.ID) {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.UserLoyaltyActivities = append(foreign.R.UserLoyaltyActivities, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userLoyaltyActivityL) LoadUser(e boil.Executor, singular bool, maybeUserLoyaltyActivity interface{}, mods queries.Applicator) error {
	var slice []*UserLoyaltyActivity
	var object *UserLoyaltyActivity

	if singular {
		object = maybeUserLoyaltyActivity.(*UserLoyaltyActivity)
	} else {
		slice = *maybeUserLoyaltyActivity.(*[]*UserLoyaltyActivity)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userLoyaltyActivityR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userLoyaltyActivityR{}
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

	if len(userLoyaltyActivityAfterSelectHooks) != 0 {
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
		foreign.R.UserLoyaltyActivities = append(foreign.R.UserLoyaltyActivities, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.UserLoyaltyActivities = append(foreign.R.UserLoyaltyActivities, local)
				break
			}
		}
	}

	return nil
}

// SetProduct of the userLoyaltyActivity to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.UserLoyaltyActivities.
func (o *UserLoyaltyActivity) SetProduct(exec boil.Executor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_loyalty_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, userLoyaltyActivityPrimaryKeyColumns),
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
		o.R = &userLoyaltyActivityR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			UserLoyaltyActivities: UserLoyaltyActivitySlice{o},
		}
	} else {
		related.R.UserLoyaltyActivities = append(related.R.UserLoyaltyActivities, o)
	}

	return nil
}

// RemoveProduct relationship.
// Sets o.R.Product to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *UserLoyaltyActivity) RemoveProduct(exec boil.Executor, related *Product) error {
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

	for i, ri := range related.R.UserLoyaltyActivities {
		if queries.Equal(o.ProductID, ri.ProductID) {
			continue
		}

		ln := len(related.R.UserLoyaltyActivities)
		if ln > 1 && i < ln-1 {
			related.R.UserLoyaltyActivities[i] = related.R.UserLoyaltyActivities[ln-1]
		}
		related.R.UserLoyaltyActivities = related.R.UserLoyaltyActivities[:ln-1]
		break
	}
	return nil
}

// SetUser of the userLoyaltyActivity to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserLoyaltyActivities.
func (o *UserLoyaltyActivity) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_loyalty_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userLoyaltyActivityPrimaryKeyColumns),
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
		o.R = &userLoyaltyActivityR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			UserLoyaltyActivities: UserLoyaltyActivitySlice{o},
		}
	} else {
		related.R.UserLoyaltyActivities = append(related.R.UserLoyaltyActivities, o)
	}

	return nil
}

// UserLoyaltyActivities retrieves all the records using an executor.
func UserLoyaltyActivities(mods ...qm.QueryMod) userLoyaltyActivityQuery {
	mods = append(mods, qm.From("\"user_loyalty_activities\""))
	return userLoyaltyActivityQuery{NewQuery(mods...)}
}

// FindUserLoyaltyActivity retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserLoyaltyActivity(exec boil.Executor, iD string, selectCols ...string) (*UserLoyaltyActivity, error) {
	userLoyaltyActivityObj := &UserLoyaltyActivity{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_loyalty_activities\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, userLoyaltyActivityObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from user_loyalty_activities")
	}

	return userLoyaltyActivityObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserLoyaltyActivity) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_loyalty_activities provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userLoyaltyActivityColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userLoyaltyActivityInsertCacheMut.RLock()
	cache, cached := userLoyaltyActivityInsertCache[key]
	userLoyaltyActivityInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userLoyaltyActivityAllColumns,
			userLoyaltyActivityColumnsWithDefault,
			userLoyaltyActivityColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_loyalty_activities\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_loyalty_activities\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into user_loyalty_activities")
	}

	if !cached {
		userLoyaltyActivityInsertCacheMut.Lock()
		userLoyaltyActivityInsertCache[key] = cache
		userLoyaltyActivityInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the UserLoyaltyActivity.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserLoyaltyActivity) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userLoyaltyActivityUpdateCacheMut.RLock()
	cache, cached := userLoyaltyActivityUpdateCache[key]
	userLoyaltyActivityUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userLoyaltyActivityAllColumns,
			userLoyaltyActivityPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update user_loyalty_activities, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_loyalty_activities\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userLoyaltyActivityPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, append(wl, userLoyaltyActivityPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update user_loyalty_activities row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for user_loyalty_activities")
	}

	if !cached {
		userLoyaltyActivityUpdateCacheMut.Lock()
		userLoyaltyActivityUpdateCache[key] = cache
		userLoyaltyActivityUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userLoyaltyActivityQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for user_loyalty_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for user_loyalty_activities")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserLoyaltyActivitySlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userLoyaltyActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_loyalty_activities\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userLoyaltyActivityPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in userLoyaltyActivity slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all userLoyaltyActivity")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserLoyaltyActivity) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no user_loyalty_activities provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userLoyaltyActivityColumnsWithDefault, o)

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

	userLoyaltyActivityUpsertCacheMut.RLock()
	cache, cached := userLoyaltyActivityUpsertCache[key]
	userLoyaltyActivityUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userLoyaltyActivityAllColumns,
			userLoyaltyActivityColumnsWithDefault,
			userLoyaltyActivityColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userLoyaltyActivityAllColumns,
			userLoyaltyActivityPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert user_loyalty_activities, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userLoyaltyActivityPrimaryKeyColumns))
			copy(conflict, userLoyaltyActivityPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_loyalty_activities\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userLoyaltyActivityType, userLoyaltyActivityMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert user_loyalty_activities")
	}

	if !cached {
		userLoyaltyActivityUpsertCacheMut.Lock()
		userLoyaltyActivityUpsertCache[key] = cache
		userLoyaltyActivityUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single UserLoyaltyActivity record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserLoyaltyActivity) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no UserLoyaltyActivity provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userLoyaltyActivityPrimaryKeyMapping)
	sql := "DELETE FROM \"user_loyalty_activities\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from user_loyalty_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for user_loyalty_activities")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userLoyaltyActivityQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no userLoyaltyActivityQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from user_loyalty_activities")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_loyalty_activities")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserLoyaltyActivitySlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userLoyaltyActivityBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userLoyaltyActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_loyalty_activities\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userLoyaltyActivityPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from userLoyaltyActivity slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for user_loyalty_activities")
	}

	if len(userLoyaltyActivityAfterDeleteHooks) != 0 {
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
func (o *UserLoyaltyActivity) Reload(exec boil.Executor) error {
	ret, err := FindUserLoyaltyActivity(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserLoyaltyActivitySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserLoyaltyActivitySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userLoyaltyActivityPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_loyalty_activities\".* FROM \"user_loyalty_activities\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userLoyaltyActivityPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in UserLoyaltyActivitySlice")
	}

	*o = slice

	return nil
}

// UserLoyaltyActivityExists checks if the UserLoyaltyActivity row exists.
func UserLoyaltyActivityExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_loyalty_activities\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if user_loyalty_activities exists")
	}

	return exists, nil
}
