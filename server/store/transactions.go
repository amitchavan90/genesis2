package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// NewTransactionStore handle Transaction methods
func NewTransactionStore(conn *sqlx.DB) *Transaction {
	ts := &Transaction{conn}
	return ts
}

// Transaction for persistence
type Transaction struct {
	Conn *sqlx.DB
}

// ManifestLineJSON struct for creating manifest_line_json for transaction tables
type ManifestLineJSON struct {
	Time            *time.Time `json:"time,omitempty"`
	CartonID        *string    `json:"cartonID,omitempty"`
	CartonCode      *string    `json:"cartonCode,omitempty"`
	ProductID       *string    `json:"productID,omitempty"`
	ProductCode     *string    `json:"productCode,omitempty"`
	TractActionName *string    `json:"tractAction,omitempty"`
	EntityName      *string    `json:"entityName,omitempty"`
	Location        *string    `json:"location,omitempty"`
	Cordinate       *string    `json:"cordinate,omitempty"`
	Note            *string    `json:"note,omitempty"`
}

// Get Transaction
func (s *Transaction) Get(id uuid.UUID) (*db.Transaction, error) {
	dat, err := db.FindTransaction(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany Transactions
func (s *Transaction) GetMany(keys []string) (db.TransactionSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{fmt.Errorf("no keys provided")}
	}
	records, err := db.Transactions(db.TransactionWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Transaction{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Transaction{}
	for _, key := range keys {
		var row *db.Transaction
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

// GetByCartonID gets Transactions by Carton ID
func (s *Transaction) GetByCartonID(id uuid.UUID) (db.TransactionSlice, error) {
	dat, err := db.Transactions(
		db.TransactionWhere.CartonID.EQ(null.StringFrom(id.String())),
		db.TransactionWhere.ProductID.IsNull(),
		db.TransactionWhere.Archived.EQ(false),
		qm.Load(db.TransactionRels.TrackAction),
	).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return dat, nil
}

// GetByProductID gets Transactions by Product ID
func (s *Transaction) GetByProductID(id uuid.UUID) (db.TransactionSlice, error) {
	dat, err := db.Transactions(
		db.TransactionWhere.ProductID.EQ(null.StringFrom(id.String())),
		db.TransactionWhere.Archived.EQ(false),
		qm.Load(db.TransactionRels.TrackAction),
	).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return dat, nil
}

// InsertByCarton attaches a transaction to a carton
func (s *Transaction) InsertByCarton(
	carton *db.Carton,
	trackAction *db.TrackAction,
	user *db.User,
	createdByName string,
	scannedAt time.Time,
	cartonPhotoBlobID null.String,
	productPhotoBlobID null.String,
	optTransaction *db.Transaction,
	txes ...*sqlx.Tx,
) (*db.Transaction, error) {
	// sanity check
	if carton == nil {
		return nil, terror.New(terror.ErrDataBlank, "carton is nil")
	}
	if user == nil {
		return nil, terror.New(terror.ErrDataBlank, "user is nil")
	}
	if trackAction == nil {
		return nil, terror.New(terror.ErrDataBlank, "trackAction is nil")
	}
	// just incase want to be anonymous
	if createdByName == "" {
		createdByName = "Anonymous"
	}

	// init data
	var t *db.Transaction
	if optTransaction != nil {
		// use optional transaction if provided
		t = optTransaction
	} else {
		// start blank transaction
		t = &db.Transaction{}
	}

	// fill in data (for admin use)
	t.TrackActionID = trackAction.ID
	t.CartonID = null.StringFrom(carton.ID)
	t.ScannedAt = null.TimeFrom(scannedAt)
	t.CartonPhotoBlobID = cartonPhotoBlobID
	t.ProductPhotoBlobID = productPhotoBlobID
	t.CreatedByID = null.StringFrom(user.ID)
	t.CreatedByName = createdByName

	if !trackAction.Blockchain {
		t.TransactionHash = null.StringFrom("-")
	} else {
		// manifest json, to be published publicly
		mj := ManifestLineJSON{
			Time:            &scannedAt,
			CartonID:        &carton.ID,
			CartonCode:      &carton.Code,
			TractActionName: &trackAction.Name,
			EntityName:      &createdByName,
		}
		// make gps coord more private
		if t.LocationGeohash.Valid && len(t.LocationGeohash.String) > 0 {
			mjcord := helpers.LimitCoordinate(t.LocationGeohash.String)
			mj.Location = &mjcord
		}
		mjs, err := json.Marshal(mj)
		if err != nil {
			return nil, terror.New(err, "json marshal fail")
		}

		t.ManifestLineJSON = null.StringFrom(string(mjs))
	}

	conn := prepConn(s.Conn, txes...)

	err := t.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "failed to insert")
	}

	return t, nil
}

// InsertByProduct add a transaction by product, accept name instead of user to help overwrite name sometimes in manifest line json
func (s *Transaction) InsertByProduct(
	product *db.Product,
	trackAction *db.TrackAction,
	user *db.User,
	createdByName string,
	optTransaction *db.Transaction,
	txes ...*sqlx.Tx,
) (*db.Transaction, error) {
	// sanity check
	if product == nil {
		return nil, terror.New(terror.ErrInvalidInput, "product is nil")
	}
	if trackAction == nil {
		return nil, terror.New(terror.ErrInvalidInput, "tractAction is nil")
	}
	// just incase want to be anonymous
	if createdByName == "" {
		createdByName = "Anonymous"
	}

	conn := prepConn(s.Conn, txes...)

	// pick time, use if provided or use now
	tm := time.Now()
	if optTransaction != nil && optTransaction.ScannedAt.Valid {
		tm = optTransaction.ScannedAt.Time
	}

	// init data
	var t *db.Transaction
	if optTransaction != nil {
		// use optional transaction if provided
		t = optTransaction
	} else {
		// start blank transaction
		t = &db.Transaction{}
	}

	// fill in data (for admin use)
	t.ProductID = null.StringFrom(product.ID)
	t.TrackActionID = trackAction.ID
	t.CreatedByName = createdByName
	if user != nil && user.ID != "" {
		t.CreatedByID = null.StringFrom(user.ID)
	}

	// manifest json, to be published publicly
	mj := ManifestLineJSON{
		Time:            &tm,
		ProductID:       &product.ID,
		ProductCode:     &product.Code,
		TractActionName: &trackAction.Name,
		EntityName:      &createdByName,
	}
	// fill in data for manifest
	// make gps coord more private
	if t.LocationGeohash.Valid && len(t.LocationGeohash.String) > 0 {
		mjcord := helpers.LimitCoordinate(t.LocationGeohash.String)
		mj.Location = &mjcord
	}
	// build string
	mjs, err := json.Marshal(mj)
	if err != nil {
		return nil, terror.New(err, "json marshal fail")
	}
	// copy into transaction
	t.ManifestLineJSON = null.StringFrom(string(mjs))

	// insert
	err = t.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "failed to insert")
	}

	return t, nil
}

// AttachManyToProduct copies many carton transactions to a product, tx a requirement because failure could cause half insert state
func (s *Transaction) AttachManyToProduct(transactions db.TransactionSlice, product *db.Product, tx *sqlx.Tx) error {
	var err error
	if tx == nil {
		return terror.New(fmt.Errorf("tx is nil"), "")
	}

	for _, tran := range transactions {
		// sanity check
		if tran.Archived {
			continue
		}
		// not copying other product history transaction
		if tran.ProductID.Valid || !tran.ProductID.IsZero() {
			continue
		}

		blankStr := null.StringFromPtr(new(string))

		// dereference to duplicate
		t := *tran

		t.ProductID = null.StringFrom(product.ID)
		// maybe redundant, due to use of boil.Whitelist()
		// these should not be copied, it needs new blockchain commit
		t.ID = ""
		t.TransactionHash = blankStr
		t.ManifestID = blankStr
		t.ManifestLineSha256 = blankStr

		if t.ManifestLineJSON.Valid && !t.ManifestLineJSON.IsZero() {
			// repack json data
			mj := ManifestLineJSON{}
			err = json.Unmarshal([]byte(t.ManifestLineJSON.String), &mj)
			if err != nil {
				return terror.New(err, "json unmarshal")
			}
			mj.ProductID = &product.ID
			mj.ProductCode = &product.Code
			mjs, err := json.Marshal(mj)
			if err != nil {
				return terror.New(err, "json marshal")
			}
			t.ManifestLineJSON = null.StringFrom(string(mjs))
		}

		// because product inherent carton history might sometimes copied same history again from prior history
		// using upsert
		err = t.Upsert(tx, false, nil, boil.Whitelist(),
			// need to use only permitted column to insert,
			// or boil.Infer() wrong column and blow up
			// because required field is blank, even though it shouldnt included
			boil.Whitelist(
				db.TransactionColumns.TrackActionID,
				db.TransactionColumns.Memo,
				db.TransactionColumns.ProductID,
				db.TransactionColumns.CartonID,
				db.TransactionColumns.ScannedAt,
				db.TransactionColumns.LocationGeohash,
				db.TransactionColumns.LocationName,
				db.TransactionColumns.ManifestLineJSON,
				db.TransactionColumns.ProductPhotoBlobID,
				db.TransactionColumns.CartonPhotoBlobID,
				db.TransactionColumns.UpdatedAt,
				db.TransactionColumns.CreatedAt,
				db.TransactionColumns.CreatedByID,
				db.TransactionColumns.CreatedByName,
			),
		)
		if err != nil {
			return terror.New(err, "upsert transactions")
		}
	}

	return nil
}

// All Transactions
func (s *Transaction) All() (db.TransactionSlice, error) {
	dat, err := db.Transactions().All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// AllPending Transactions that will be made into public manifest
func (s *Transaction) AllPending() (db.TransactionSlice, error) {
	dat, err := db.Transactions(
		db.TransactionWhere.ManifestID.EQ(null.StringFromPtr(nil)),
		db.TransactionWhere.TransactionHash.IsNull(),
		db.TransactionWhere.Archived.EQ(false),
		qm.OrderBy(db.TransactionColumns.UpdatedAt),
	).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// AllPendingCount a count of transactions that will be made into public manifest, otherwise it will kill the server for large query
func (s *Transaction) AllPendingCount() (int64, error) {
	dat, err := db.Transactions(
		db.TransactionWhere.ManifestID.EQ(null.StringFromPtr(nil)),
		db.TransactionWhere.TransactionHash.IsNull(),
		db.TransactionWhere.Archived.EQ(false),
	).Count(s.Conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return dat, nil
}

type countResult struct {
	Count int64 `db:"count"`
}

// SearchSelect searchs/selects Products
func (s *Transaction) SearchSelect(
	search graphql.SearchFilter,
	limit int,
	offset int,
	productID null.String,
	cartonID null.String,
	trackActionID null.String,
) (int64, []*db.Transaction, error) {
	// holds where condition
	qms := []qm.QueryMod{}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			qms = append(qms, db.TransactionWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			qms = append(qms, db.TransactionWhere.Archived.EQ(true))
		case graphql.FilterOptionPending:
			qms = append(qms, db.TransactionWhere.TransactionHash.IsNull())
		}
	}

	// Hide non-blockchain transactions
	if search.Filter == nil || *search.Filter != graphql.FilterOptionPending {
		// TransactionHash == '-'
		// or
		// TransactionHash is null
		qms = append(
			qms,
			qm.And(
				fmt.Sprintf(
					"(%s != '-' OR %s IS NULL)",
					db.TransactionColumns.TransactionHash,
					db.TransactionColumns.TransactionHash,
				),
			),
		)
	}

	// Filter by product/carton/track action
	if productID.Valid {
		qms = append(qms, db.TransactionWhere.ProductID.EQ(productID))
	}
	if cartonID.Valid {
		qms = append(qms, db.TransactionWhere.CartonID.EQ(cartonID))
	}
	if trackActionID.Valid {
		qms = append(qms, db.TransactionWhere.TrackActionID.EQ(trackActionID.String))
	}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			qms = append(qms,
				qm.And(
					fmt.Sprintf(
						"((LOWER(%s) LIKE '%s') OR (LOWER(%s) LIKE '%s'))",
						db.TransactionColumns.CreatedByName,
						"%"+searchText+"%",
						db.TransactionColumns.LocationName,
						"%"+searchText+"%",
					),
				),
			)
		}
	}

	// Sort
	sortDir := "ASC"
	if search.SortDir != nil && *search.SortDir == graphql.SortDirDescending {
		sortDir = "DESC"
	}

	sort := db.TransactionColumns.ScannedAt
	if search.SortBy != nil {
		switch *search.SortBy {
		case graphql.SortByOptionDateUpdated:
			sort = db.TransactionColumns.UpdatedAt
		case graphql.SortByOptionAlphabetical:
			sort = db.TransactionColumns.CreatedAt
		}
	}

	count, err := db.Transactions(qms...).Count(s.Conn)
	if err != nil {
		return 0, nil, terror.New(err, "count transactions")
	}

	qms = append(qms, qm.OrderBy(sort+" "+sortDir))
	qms = append(qms, qm.Limit(limit))
	qms = append(qms, qm.Offset(offset))

	result, err := db.Transactions(qms...).All(s.Conn)
	if err != nil {
		return 0, nil, terror.New(err, "count transactions")
	}

	return count, result, nil
}

// Insert Transaction
func (s *Transaction) Insert(record *db.Transaction, txes ...*sqlx.Tx) (*db.Transaction, error) {
	conn := prepConn(s.Conn, txes...)

	err := record.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update Transaction
func (s *Transaction) Update(record *db.Transaction, txes ...*sql.Tx) (*db.Transaction, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive Transactions
func (s *Transaction) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Transaction, error) {
	u, err := db.FindTransaction(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.TransactionColumns.Archived, db.TransactionColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive Transactions
func (s *Transaction) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Transaction, error) {
	u, err := db.FindTransaction(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.TransactionColumns.Archived, db.TransactionColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// GetSettings returns the single settings record
func (s *Transaction) GetSettings() (*db.Setting, error) {
	dat, err := db.Settings().One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// UpdateSettings updates the setting record
func (s *Transaction) UpdateSettings(record *db.Setting) (*db.Setting, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}
