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

// Manifest is an object representing the database table.
type Manifest struct {
	ID               string      `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	ContractAddress  string      `db:"contract_address" boil:"contract_address" json:"contract_address" toml:"contract_address" yaml:"contract_address"`
	TransactionNonce int         `db:"transaction_nonce" boil:"transaction_nonce" json:"transaction_nonce" toml:"transaction_nonce" yaml:"transaction_nonce"`
	TransactionHash  null.String `db:"transaction_hash" boil:"transaction_hash" json:"transaction_hash,omitempty" toml:"transaction_hash" yaml:"transaction_hash,omitempty"`
	Confirmed        bool        `db:"confirmed" boil:"confirmed" json:"confirmed" toml:"confirmed" yaml:"confirmed"`
	MerkleRootSha256 null.String `db:"merkle_root_sha256" boil:"merkle_root_sha256" json:"merkle_root_sha256,omitempty" toml:"merkle_root_sha256" yaml:"merkle_root_sha256,omitempty"`
	CompiledText     null.Bytes  `db:"compiled_text" boil:"compiled_text" json:"compiled_text,omitempty" toml:"compiled_text" yaml:"compiled_text,omitempty"`
	Pending          bool        `db:"pending" boil:"pending" json:"pending" toml:"pending" yaml:"pending"`
	Archived         bool        `db:"archived" boil:"archived" json:"archived" toml:"archived" yaml:"archived"`
	ArchivedAt       null.Time   `db:"archived_at" boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`
	UpdatedAt        time.Time   `db:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedAt        time.Time   `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *manifestR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L manifestL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ManifestColumns = struct {
	ID               string
	ContractAddress  string
	TransactionNonce string
	TransactionHash  string
	Confirmed        string
	MerkleRootSha256 string
	CompiledText     string
	Pending          string
	Archived         string
	ArchivedAt       string
	UpdatedAt        string
	CreatedAt        string
}{
	ID:               "id",
	ContractAddress:  "contract_address",
	TransactionNonce: "transaction_nonce",
	TransactionHash:  "transaction_hash",
	Confirmed:        "confirmed",
	MerkleRootSha256: "merkle_root_sha256",
	CompiledText:     "compiled_text",
	Pending:          "pending",
	Archived:         "archived",
	ArchivedAt:       "archived_at",
	UpdatedAt:        "updated_at",
	CreatedAt:        "created_at",
}

// Generated where

type whereHelpernull_Bytes struct{ field string }

func (w whereHelpernull_Bytes) EQ(x null.Bytes) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Bytes) NEQ(x null.Bytes) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Bytes) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Bytes) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Bytes) LT(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Bytes) LTE(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Bytes) GT(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Bytes) GTE(x null.Bytes) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ManifestWhere = struct {
	ID               whereHelperstring
	ContractAddress  whereHelperstring
	TransactionNonce whereHelperint
	TransactionHash  whereHelpernull_String
	Confirmed        whereHelperbool
	MerkleRootSha256 whereHelpernull_String
	CompiledText     whereHelpernull_Bytes
	Pending          whereHelperbool
	Archived         whereHelperbool
	ArchivedAt       whereHelpernull_Time
	UpdatedAt        whereHelpertime_Time
	CreatedAt        whereHelpertime_Time
}{
	ID:               whereHelperstring{field: "\"manifests\".\"id\""},
	ContractAddress:  whereHelperstring{field: "\"manifests\".\"contract_address\""},
	TransactionNonce: whereHelperint{field: "\"manifests\".\"transaction_nonce\""},
	TransactionHash:  whereHelpernull_String{field: "\"manifests\".\"transaction_hash\""},
	Confirmed:        whereHelperbool{field: "\"manifests\".\"confirmed\""},
	MerkleRootSha256: whereHelpernull_String{field: "\"manifests\".\"merkle_root_sha256\""},
	CompiledText:     whereHelpernull_Bytes{field: "\"manifests\".\"compiled_text\""},
	Pending:          whereHelperbool{field: "\"manifests\".\"pending\""},
	Archived:         whereHelperbool{field: "\"manifests\".\"archived\""},
	ArchivedAt:       whereHelpernull_Time{field: "\"manifests\".\"archived_at\""},
	UpdatedAt:        whereHelpertime_Time{field: "\"manifests\".\"updated_at\""},
	CreatedAt:        whereHelpertime_Time{field: "\"manifests\".\"created_at\""},
}

// ManifestRels is where relationship names are stored.
var ManifestRels = struct {
	Transactions string
}{
	Transactions: "Transactions",
}

// manifestR is where relationships are stored.
type manifestR struct {
	Transactions TransactionSlice
}

// NewStruct creates a new relationship struct
func (*manifestR) NewStruct() *manifestR {
	return &manifestR{}
}

// manifestL is where Load methods for each relationship are stored.
type manifestL struct{}

var (
	manifestAllColumns            = []string{"id", "contract_address", "transaction_nonce", "transaction_hash", "confirmed", "merkle_root_sha256", "compiled_text", "pending", "archived", "archived_at", "updated_at", "created_at"}
	manifestColumnsWithoutDefault = []string{"contract_address", "transaction_nonce", "transaction_hash", "merkle_root_sha256", "compiled_text", "archived_at"}
	manifestColumnsWithDefault    = []string{"id", "confirmed", "pending", "archived", "updated_at", "created_at"}
	manifestPrimaryKeyColumns     = []string{"id"}
)

type (
	// ManifestSlice is an alias for a slice of pointers to Manifest.
	// This should generally be used opposed to []Manifest.
	ManifestSlice []*Manifest
	// ManifestHook is the signature for custom Manifest hook methods
	ManifestHook func(boil.Executor, *Manifest) error

	manifestQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	manifestType                 = reflect.TypeOf(&Manifest{})
	manifestMapping              = queries.MakeStructMapping(manifestType)
	manifestPrimaryKeyMapping, _ = queries.BindMapping(manifestType, manifestMapping, manifestPrimaryKeyColumns)
	manifestInsertCacheMut       sync.RWMutex
	manifestInsertCache          = make(map[string]insertCache)
	manifestUpdateCacheMut       sync.RWMutex
	manifestUpdateCache          = make(map[string]updateCache)
	manifestUpsertCacheMut       sync.RWMutex
	manifestUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var manifestBeforeInsertHooks []ManifestHook
var manifestBeforeUpdateHooks []ManifestHook
var manifestBeforeDeleteHooks []ManifestHook
var manifestBeforeUpsertHooks []ManifestHook

var manifestAfterInsertHooks []ManifestHook
var manifestAfterSelectHooks []ManifestHook
var manifestAfterUpdateHooks []ManifestHook
var manifestAfterDeleteHooks []ManifestHook
var manifestAfterUpsertHooks []ManifestHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Manifest) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Manifest) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Manifest) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Manifest) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Manifest) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Manifest) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Manifest) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Manifest) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Manifest) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range manifestAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddManifestHook registers your hook function for all future operations.
func AddManifestHook(hookPoint boil.HookPoint, manifestHook ManifestHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		manifestBeforeInsertHooks = append(manifestBeforeInsertHooks, manifestHook)
	case boil.BeforeUpdateHook:
		manifestBeforeUpdateHooks = append(manifestBeforeUpdateHooks, manifestHook)
	case boil.BeforeDeleteHook:
		manifestBeforeDeleteHooks = append(manifestBeforeDeleteHooks, manifestHook)
	case boil.BeforeUpsertHook:
		manifestBeforeUpsertHooks = append(manifestBeforeUpsertHooks, manifestHook)
	case boil.AfterInsertHook:
		manifestAfterInsertHooks = append(manifestAfterInsertHooks, manifestHook)
	case boil.AfterSelectHook:
		manifestAfterSelectHooks = append(manifestAfterSelectHooks, manifestHook)
	case boil.AfterUpdateHook:
		manifestAfterUpdateHooks = append(manifestAfterUpdateHooks, manifestHook)
	case boil.AfterDeleteHook:
		manifestAfterDeleteHooks = append(manifestAfterDeleteHooks, manifestHook)
	case boil.AfterUpsertHook:
		manifestAfterUpsertHooks = append(manifestAfterUpsertHooks, manifestHook)
	}
}

// One returns a single manifest record from the query.
func (q manifestQuery) One(exec boil.Executor) (*Manifest, error) {
	o := &Manifest{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: failed to execute a one query for manifests")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Manifest records from the query.
func (q manifestQuery) All(exec boil.Executor) (ManifestSlice, error) {
	var o []*Manifest

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "db: failed to assign all query results to Manifest slice")
	}

	if len(manifestAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Manifest records in the query.
func (q manifestQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to count manifests rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q manifestQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "db: failed to check if manifests exists")
	}

	return count > 0, nil
}

// Transactions retrieves all the transaction's Transactions with an executor.
func (o *Manifest) Transactions(mods ...qm.QueryMod) transactionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"transactions\".\"manifest_id\"=?", o.ID),
	)

	query := Transactions(queryMods...)
	queries.SetFrom(query.Query, "\"transactions\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"transactions\".*"})
	}

	return query
}

// LoadTransactions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (manifestL) LoadTransactions(e boil.Executor, singular bool, maybeManifest interface{}, mods queries.Applicator) error {
	var slice []*Manifest
	var object *Manifest

	if singular {
		object = maybeManifest.(*Manifest)
	} else {
		slice = *maybeManifest.(*[]*Manifest)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &manifestR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &manifestR{}
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

	query := NewQuery(qm.From(`transactions`), qm.WhereIn(`transactions.manifest_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load transactions")
	}

	var resultSlice []*Transaction
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice transactions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on transactions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for transactions")
	}

	if len(transactionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Transactions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &transactionR{}
			}
			foreign.R.Manifest = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.ManifestID) {
				local.R.Transactions = append(local.R.Transactions, foreign)
				if foreign.R == nil {
					foreign.R = &transactionR{}
				}
				foreign.R.Manifest = local
				break
			}
		}
	}

	return nil
}

// AddTransactions adds the given related objects to the existing relationships
// of the manifest, optionally inserting them as new records.
// Appends related to o.R.Transactions.
// Sets related.R.Manifest appropriately.
func (o *Manifest) AddTransactions(exec boil.Executor, insert bool, related ...*Transaction) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.ManifestID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"transactions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"manifest_id"}),
				strmangle.WhereClause("\"", "\"", 2, transactionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.ManifestID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &manifestR{
			Transactions: related,
		}
	} else {
		o.R.Transactions = append(o.R.Transactions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &transactionR{
				Manifest: o,
			}
		} else {
			rel.R.Manifest = o
		}
	}
	return nil
}

// SetTransactions removes all previously related items of the
// manifest replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Manifest's Transactions accordingly.
// Replaces o.R.Transactions with related.
// Sets related.R.Manifest's Transactions accordingly.
func (o *Manifest) SetTransactions(exec boil.Executor, insert bool, related ...*Transaction) error {
	query := "update \"transactions\" set \"manifest_id\" = null where \"manifest_id\" = $1"
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
		for _, rel := range o.R.Transactions {
			queries.SetScanner(&rel.ManifestID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Manifest = nil
		}

		o.R.Transactions = nil
	}
	return o.AddTransactions(exec, insert, related...)
}

// RemoveTransactions relationships from objects passed in.
// Removes related items from R.Transactions (uses pointer comparison, removal does not keep order)
// Sets related.R.Manifest.
func (o *Manifest) RemoveTransactions(exec boil.Executor, related ...*Transaction) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.ManifestID, nil)
		if rel.R != nil {
			rel.R.Manifest = nil
		}
		if _, err = rel.Update(exec, boil.Whitelist("manifest_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Transactions {
			if rel != ri {
				continue
			}

			ln := len(o.R.Transactions)
			if ln > 1 && i < ln-1 {
				o.R.Transactions[i] = o.R.Transactions[ln-1]
			}
			o.R.Transactions = o.R.Transactions[:ln-1]
			break
		}
	}

	return nil
}

// Manifests retrieves all the records using an executor.
func Manifests(mods ...qm.QueryMod) manifestQuery {
	mods = append(mods, qm.From("\"manifests\""))
	return manifestQuery{NewQuery(mods...)}
}

// FindManifest retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindManifest(exec boil.Executor, iD string, selectCols ...string) (*Manifest, error) {
	manifestObj := &Manifest{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"manifests\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, manifestObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "db: unable to select from manifests")
	}

	return manifestObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Manifest) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("db: no manifests provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(manifestColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	manifestInsertCacheMut.RLock()
	cache, cached := manifestInsertCache[key]
	manifestInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			manifestAllColumns,
			manifestColumnsWithDefault,
			manifestColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(manifestType, manifestMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(manifestType, manifestMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"manifests\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"manifests\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "db: unable to insert into manifests")
	}

	if !cached {
		manifestInsertCacheMut.Lock()
		manifestInsertCache[key] = cache
		manifestInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// Update uses an executor to update the Manifest.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Manifest) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	manifestUpdateCacheMut.RLock()
	cache, cached := manifestUpdateCache[key]
	manifestUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			manifestAllColumns,
			manifestPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("db: unable to update manifests, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"manifests\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, manifestPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(manifestType, manifestMapping, append(wl, manifestPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "db: unable to update manifests row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by update for manifests")
	}

	if !cached {
		manifestUpdateCacheMut.Lock()
		manifestUpdateCache[key] = cache
		manifestUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q manifestQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all for manifests")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected for manifests")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ManifestSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), manifestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"manifests\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, manifestPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to update all in manifest slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to retrieve rows affected all in update all manifest")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Manifest) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("db: no manifests provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(manifestColumnsWithDefault, o)

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

	manifestUpsertCacheMut.RLock()
	cache, cached := manifestUpsertCache[key]
	manifestUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			manifestAllColumns,
			manifestColumnsWithDefault,
			manifestColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			manifestAllColumns,
			manifestPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("db: unable to upsert manifests, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(manifestPrimaryKeyColumns))
			copy(conflict, manifestPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"manifests\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(manifestType, manifestMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(manifestType, manifestMapping, ret)
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
		return errors.Wrap(err, "db: unable to upsert manifests")
	}

	if !cached {
		manifestUpsertCacheMut.Lock()
		manifestUpsertCache[key] = cache
		manifestUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// Delete deletes a single Manifest record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Manifest) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("db: no Manifest provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), manifestPrimaryKeyMapping)
	sql := "DELETE FROM \"manifests\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete from manifests")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by delete for manifests")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q manifestQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("db: no manifestQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from manifests")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for manifests")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ManifestSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(manifestBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), manifestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"manifests\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, manifestPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "db: unable to delete all from manifest slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "db: failed to get rows affected by deleteall for manifests")
	}

	if len(manifestAfterDeleteHooks) != 0 {
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
func (o *Manifest) Reload(exec boil.Executor) error {
	ret, err := FindManifest(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ManifestSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ManifestSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), manifestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"manifests\".* FROM \"manifests\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, manifestPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "db: unable to reload all in ManifestSlice")
	}

	*o = slice

	return nil
}

// ManifestExists checks if the Manifest row exists.
func ManifestExists(exec boil.Executor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"manifests\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "db: unable to check if manifests exists")
	}

	return exists, nil
}
