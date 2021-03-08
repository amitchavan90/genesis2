package store

import (
	"database/sql"
	"errors"
	"genesis/db"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/sqlboiler/boil"
	"syreclabs.com/go/faker"
)

// OrganisationFactory creates orgs
func OrganisationFactory() *db.Organisation {
	u := &db.Organisation{
		ID:   uuid.Must(uuid.NewV4()).String(),
		Name: faker.Company().Name(),
	}
	return u
}

// Organisation for persistence
type Organisation struct {
	Conn *sqlx.DB
}

// NewOrganisationStore returns a new store
func NewOrganisationStore(conn *sqlx.DB) *Organisation {
	os := &Organisation{conn}
	return os
}

// All organisations
func (s *Organisation) All() (db.OrganisationSlice, error) {
	dat, err := db.Organisations().All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Get organisations
func (s *Organisation) Get(id uuid.UUID) (*db.Organisation, error) {
	dat, err := db.FindOrganisation(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetMany organisations
func (s *Organisation) GetMany(keys []string) (db.OrganisationSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Organisations(db.OrganisationWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Organisation{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Organisation{}
	for _, key := range keys {
		var row *db.Organisation
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

// Insert organisations
func (s *Organisation) Insert(record *db.Organisation, txes ...*sql.Tx) (*db.Organisation, error) {
	var err error

	handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)

	err = record.Reload(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update organisations
func (s *Organisation) Update(record *db.Organisation, txes ...*sql.Tx) (*db.Organisation, error) {
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}
