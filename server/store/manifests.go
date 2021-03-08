package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// NewManifestStore handle Manifest methods
func NewManifestStore(conn *sqlx.DB) *Manifest {
	ts := &Manifest{conn}
	return ts
}

// Manifest for persistence
type Manifest struct {
	Conn *sqlx.DB
}

// BeginTransaction will start a new transaction for use with other stores
func (m *Manifest) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := m.Conn.Beginx()
	if err != nil {
		return nil, terror.New(err, "")
	}
	return tx, nil
}

// Get Manifest
func (m *Manifest) Get(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.FindManifest(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByLineHash by line hash
func (m *Manifest) GetByLineHash(lineHash string, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	// sanity check
	if lineHash == "" {
		return nil, terror.New(terror.ErrInvalidInput, "")
	}

	dt, err := db.Transactions(
		db.TransactionWhere.ManifestLineSha256.EQ(null.StringFrom(lineHash)),
		db.TransactionWhere.Archived.EQ(false),
	).One(conn)
	if err != nil {
		return nil, terror.New(err, "get db.transaction by line hash")
	}
	// sanity check
	if dt.Archived || !dt.ManifestID.Valid || dt.ManifestID.String == "" {
		return nil, terror.New(fmt.Errorf("invalid db.Transaction data"), "")
	}

	mn, err := db.FindManifest(conn, dt.ManifestID.String)
	if err != nil {
		return nil, terror.New(err, "get manifest by line hash")
	}
	// sanity check
	if mn.Archived || mn.Pending || !mn.Confirmed {
		return nil, terror.New(fmt.Errorf("invalid manifest data"), "")
	}

	return mn, nil
}

// GetByBlockchainByMerkleRootHash by Merkle Root hash
func (m *Manifest) GetByBlockchainByMerkleRootHash(merkleRootHash string, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	// sanity check
	if merkleRootHash == "" {
		return nil, terror.New(terror.ErrInvalidInput, "")
	}

	dat, err := db.Manifests(
		db.ManifestWhere.MerkleRootSha256.EQ(null.StringFrom(merkleRootHash)),
		db.ManifestWhere.Archived.EQ(false),
	).One(conn)
	if err != nil {
		return nil, terror.New(err, "get manifest by blockchain transaction hash")
	}
	return dat, nil
}

// GetByBlockchainByTransactionHash by tx hash
func (m *Manifest) GetByBlockchainByTransactionHash(bcTxHash string, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	// sanity check
	if bcTxHash == "" {
		return nil, terror.New(terror.ErrInvalidInput, "")
	}

	dat, err := db.Manifests(
		db.ManifestWhere.TransactionHash.EQ(null.StringFrom(bcTxHash)),
		db.ManifestWhere.Archived.EQ(false),
	).One(conn)
	if err != nil {
		return nil, terror.New(err, "get manifest by blockchain transaction hash")
	}
	return dat, nil
}

// GetMany Manifests
func (m *Manifest) GetMany(keys []string, txes ...*sqlx.Tx) (db.ManifestSlice, []error) {
	conn := prepConn(m.Conn, txes...)

	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Manifests(db.ManifestWhere.ID.IN(keys)).All(conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Manifest{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Manifest{}
	for _, key := range keys {
		var row *db.Manifest
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

// All Manifests
func (m *Manifest) All(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		qm.OrderBy(db.ManifestColumns.TransactionNonce),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// AllUnarchived Manifests
func (m *Manifest) AllUnarchived(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		db.ManifestWhere.Archived.EQ(false),
		qm.OrderBy(db.ManifestColumns.TransactionNonce),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// AllPending Manifests
func (m *Manifest) AllPending(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		db.ManifestWhere.Pending.EQ(true),
		db.ManifestWhere.Confirmed.EQ(false),
		db.ManifestWhere.Archived.EQ(false),
		qm.OrderBy(db.ManifestColumns.TransactionNonce),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "get confirmed manifest")
	}

	return dat, nil
}

// AllConfirmed Manifests
func (m *Manifest) AllConfirmed(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		db.ManifestWhere.Pending.EQ(false),
		db.ManifestWhere.Confirmed.EQ(true),
		db.ManifestWhere.Archived.EQ(false),
		qm.OrderBy(db.ManifestColumns.TransactionNonce),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "get confirmed manifest")
	}

	return dat, nil
}

// AllUnconfirmed Manifests
func (m *Manifest) AllUnconfirmed(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		db.ManifestWhere.Pending.EQ(false),
		db.ManifestWhere.Confirmed.EQ(false),
		db.ManifestWhere.Archived.EQ(false),
		qm.OrderBy(db.ManifestColumns.TransactionNonce),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "get unconfirmed manifest")
	}

	return dat, nil
}

// AllUnfinished Manifests
func (m *Manifest) AllUnfinished(txes ...*sqlx.Tx) (db.ManifestSlice, error) {
	conn := prepConn(m.Conn, txes...)

	dat, err := db.Manifests(
		// sql query select
		// select * where
		qm.Expr(
			db.ManifestWhere.Pending.EQ(true),
			db.ManifestWhere.Archived.EQ(false),
			qm.OrderBy(db.ManifestColumns.TransactionNonce),
		),
		// or
		qm.Or2(
			qm.Expr(
				db.ManifestWhere.Confirmed.EQ(false),
				db.ManifestWhere.Archived.EQ(false),
				qm.OrderBy(db.ManifestColumns.TransactionNonce),
			),
		),
	).All(conn)
	if err != nil {
		return nil, terror.New(err, "get unfinished manifest")
	}

	return dat, nil
}

// Insert Manifest
func (m *Manifest) Insert(record *db.Manifest, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	err := record.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "failed to insert")
	}
	return record, nil
}

// Update Manifest
func (m *Manifest) Update(record *db.Manifest, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	record.UpdatedAt = time.Now()
	_, err := record.Update(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive Manifests
func (m *Manifest) Archive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	u, err := db.FindManifest(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(conn, boil.Whitelist(db.ManifestColumns.Archived, db.ManifestColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive Manifests
func (m *Manifest) Unarchive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error) {
	conn := prepConn(m.Conn, txes...)

	u, err := db.FindManifest(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(conn, boil.Whitelist(db.ManifestColumns.Archived, db.ManifestColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}
