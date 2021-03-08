package store

import (
	"database/sql"
	"genesis/db"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"syreclabs.com/go/faker"
)

// Prospects for persistence
type Prospects struct {
	Conn boil.Executor
}

// ProspectFactory creates prosects
func ProspectFactory() *db.Prospect {
	u := &db.Prospect{
		ID:                 uuid.Must(uuid.NewV4()).String(),
		Email:              faker.Internet().Email(),
		FirstName:          null.StringFrom(faker.Name().FirstName()),
		LastName:           null.StringFrom(faker.Name().LastName()),
		OnboardingComplete: false,
	}
	return u
}

// Get a prospect
func (s *Prospects) Get(id uuid.UUID, txes ...*sql.Tx) (*db.Prospect, error) {
	return db.FindProspect(s.Conn, id.String())
}

// Start a prospect
func (s *Prospects) Start(email string, txes ...*sql.Tx) (*db.Prospect, error) {
	p := &db.Prospect{
		Email: email,
	}
	err := p.Insert(s.Conn, boil.Whitelist(db.ProspectColumns.Email))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return p, nil
}

// Update a prospect
func (s *Prospects) Update(p *db.Prospect, txes ...*sql.Tx) (*db.Prospect, error) {
	_, err := p.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return p, nil
}

// Finish a prospect
func (s *Prospects) Finish(id uuid.UUID, txes ...*sql.Tx) (*db.Prospect, error) {
	p, err := db.FindProspect(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	p.OnboardingComplete = true
	_, err = p.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return p, nil

}
