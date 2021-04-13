package store

import (
	"database/sql"
	"errors"
	"genesis/db"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/sqlboiler/boil"
)

// Blob for persistence
type Blob struct {
	Conn *sqlx.DB
}

// NewBlobStore handle blob methods
func NewBlobStore(conn *sqlx.DB) *Blob {
	ts := &Blob{conn}
	return ts
}

// Get blobs
func (s *Blob) Get(id uuid.UUID) (*db.Blob, error) {
	return db.FindBlob(s.Conn, id.String())
}

// GetMany blobs given a list of IDs
func (s *Blob) GetMany(keys []string) (db.BlobSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Blobs(db.BlobWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Blob{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Blob{}
	for _, key := range keys {
		var row *db.Blob
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

// All blobs
func (s *Blob) All() (db.BlobSlice, error) {
	return db.Blobs().All(s.Conn)
}

// Insert blobs
func (s *Blob) Insert(record *db.Blob, txes ...*sql.Tx) (*db.Blob, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update blobs
func (s *Blob) Update(record *db.Blob, txes ...*sql.Tx) (*db.Blob, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Delete blobs
func (s *Blob) Delete(record *db.Blob, txes ...*sql.Tx) error {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		_, err := record.Delete(tx)
		return err
	}, txes...)
	if err != nil {
		return err
	}
	return nil
}

// Exists return whether a blob exists
func (s *Blob) Exists(id string) (bool, error) {
	return db.BlobExists(s.Conn, id)
}
