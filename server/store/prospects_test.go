package store_test

import (
	"genesis/store"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"syreclabs.com/go/faker"
)

func TestProspects_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := &store.Prospects{conn}
	t.Run("no records", func(t *testing.T) {
		drop()
		migrate()
		randomID := uuid.Must(uuid.NewV4())
		result, err := repo.Get(randomID)
		if err == nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}
		if result != nil {
			t.Errorf("result: got %v, expected %v", result, nil)
		}
	})
}

func TestProspects_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := &store.Prospects{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()
		email := faker.Internet().Email()
		result, err := repo.Start(email)
		if err != nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}
		if result == nil {
			t.Errorf("result: got %v, expected %v", result, nil)
		}
	})
}

func TestProspects_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := &store.Prospects{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		prevName := null.StringFrom(faker.Name().FirstName())
		nextName := null.StringFrom(faker.Name().FirstName())

		p := store.ProspectFactory()
		p.FirstName = prevName
		err := p.Insert(conn, boil.Infer())
		if err != nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}

		p.FirstName = nextName

		result, err := repo.Update(p)
		if err != nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}

		if result.FirstName != nextName {
			t.Errorf("result.FirstName: got %v, expected %v", result.FirstName, nextName)
		}
	})
}

func TestProspects_Finish(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	conn, drop, migrate, teardown := Setup(t)
	defer teardown()
	repo := &store.Prospects{conn}
	t.Run("happy path", func(t *testing.T) {
		drop()
		migrate()

		prevName := null.StringFrom(faker.Name().FirstName())
		nextName := null.StringFrom(faker.Name().FirstName())

		p := store.ProspectFactory()
		p.FirstName = prevName
		err := p.Insert(conn, boil.Infer())
		if err != nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}

		p.FirstName = nextName

		_, err = repo.Finish(uuid.Must(uuid.FromString(p.ID)))
		if err != nil {
			t.Errorf("err: got %v, expected %v", err, nil)
		}

		p, err = repo.Get(uuid.Must(uuid.FromString(p.ID)))
		if !p.OnboardingComplete {
			t.Errorf("p.OnboardingComplete: got %v, expected %v", p.OnboardingComplete, true)
		}
	})
}
