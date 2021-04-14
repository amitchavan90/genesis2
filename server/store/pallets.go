package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"syreclabs.com/go/faker"
)

// PalletFactory creates pallets
func PalletFactory() *db.Pallet {
	u := &db.Pallet{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Code:        faker.Numerify("PAL#####"),
		Description: faker.Lorem().Sentence(10),
	}
	return u
}

// NewPalletStore handle pallet methods
func NewPalletStore(conn *sqlx.DB) *Pallet {
	ts := &Pallet{conn}
	return ts
}

// Pallet for persistence
type Pallet struct {
	Conn *sqlx.DB
}

// GetByCode gets a Pallet by their code
func (s *Pallet) GetByCode(code string) (*db.Pallet, error) {
	dat, err := db.Pallets(db.PalletWhere.Code.EQ(code)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Get pallets
func (s *Pallet) Get(id uuid.UUID) (*db.Pallet, error) {
	dat, err := db.FindPallet(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// All pallets
func (s *Pallet) All() (db.PalletSlice, error) {
	dat, err := db.Pallets().All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// SearchSelect searchs/selects Pallets
func (s *Pallet) SearchSelect(
	search graphql.SearchFilter,
	limit int,
	offset int,
	containerID null.String,
	trackActionID null.String,
) (int64, []*db.Pallet, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.PalletColumns.Code,
						db.PalletColumns.Description),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.PalletWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.PalletWhere.Archived.EQ(true))
		case graphql.FilterOptionPalletWithoutContainer:
			queries = append(queries, db.PalletWhere.ContainerID.EQ(null.StringFromPtr(nil)))
		}
	}

	if containerID.Valid {
		queries = append(queries, db.PalletWhere.ContainerID.EQ(containerID))
	}
	if trackActionID.Valid {
		// using combination of pallet_latest_transactions view and the inner join lateral
		// to complete this complex search
		// which is to search by the pallet's last track action
		// it will not search pallet's previous track action
		// note: t1.pallet_id only exist in the view
		queries = append(
			queries,
			qm.InnerJoin(
				`LATERAL
				(
					SELECT pallet_id, track_action_id FROM pallet_latest_transactions
					WHERE track_action_id = ?
				) t1 ON pallets.id = t1.pallet_id
				`,
				trackActionID.String,
			),
		)
	}

	// Get Count
	count, err := db.Pallets(queries...).Count(s.Conn)
	if err != nil {
		return 0, nil, terror.New(err, "")
	}

	// Sort
	sortDir := " ASC"
	if search.SortDir != nil && *search.SortDir == graphql.SortDirDescending {
		sortDir = " DESC"
	}

	if search.SortBy != nil {
		switch *search.SortBy {
		case graphql.SortByOptionDateCreated:
			queries = append(queries, qm.OrderBy(db.PalletColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.PalletColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.PalletColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.PalletColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Pallets(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany pallets
func (s *Pallet) GetMany(keys []string) (db.PalletSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{fmt.Errorf("no keys provided")}
	}
	records, err := db.Pallets(db.PalletWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Pallet{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Pallet{}
	for _, key := range keys {
		var row *db.Pallet
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

// GetRange returns all pallets between two pallet codes
func (s *Pallet) GetRange(from string, to string) (db.PalletSlice, error) {
	start, err := strconv.Atoi(strings.Replace(from, "PAL", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}
	end, err := strconv.Atoi(strings.Replace(to, "PAL", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}

	if start < 0 || end < start {
		return nil, terror.New(terror.ErrInvalidInput, "")
	}

	keys := []string{}
	for i := start; i <= end; i++ {
		keys = append(keys, fmt.Sprintf("PAL%05d", i))
	}

	records, err := db.Pallets(db.PalletWhere.Code.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Pallet{}, nil
	}
	if err != nil {
		return nil, terror.New(err, "")
	}

	result := []*db.Pallet{}
	for _, key := range keys {
		var row *db.Pallet
		for _, record := range records {
			if record.Code == key {
				row = record
				break
			}
		}
		result = append(result, row)
	}
	return result, nil
}

// Insert pallets
func (s *Pallet) Insert(record *db.Pallet, txes ...*sql.Tx) (*db.Pallet, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update pallets
func (s *Pallet) Update(record *db.Pallet, txes ...*sql.Tx) (*db.Pallet, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive pallets
func (s *Pallet) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Pallet, error) {
	u, err := db.FindPallet(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.PalletColumns.Archived, db.PalletColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive pallets
func (s *Pallet) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Pallet, error) {
	u, err := db.FindPallet(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.PalletColumns.Archived, db.PalletColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of pallets
func (s *Pallet) Count() (int64, error) {
	dat, err := db.Pallets().Count(s.Conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return dat, nil
}

// CartonCount returns the count of cartons in the pallet
func (s *Pallet) CartonCount(record *db.Pallet) (int64, error) {
	dat, err := db.Cartons(db.CartonWhere.PalletID.EQ(null.StringFrom(record.ID))).Count(s.Conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return dat, nil
}

// LatestTrackAction returns the latest transaction track action name
func (s *Pallet) LatestTrackAction(record *db.Pallet) (*graphql.LatestTransactionInfo, error) {
	// get carton
	carton, err := record.Cartons().One(s.Conn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, terror.New(err, "failed to get pallet carton")
	}

	// get carton's latest transaction track action name
	transaction, err := db.Transactions(
		db.TransactionWhere.CartonID.EQ(null.StringFrom(carton.ID)),
		qm.OrderBy(
			fmt.Sprintf("COALESCE(%s,%s) DESC",
				db.TransactionColumns.ScannedAt,
				db.TransactionColumns.CreatedAt,
			),
		),
		qm.Load(
			db.TransactionRels.TrackAction,
		),
	).One(s.Conn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, terror.New(err, "")
	}

	date := transaction.CreatedAt
	if transaction.ScannedAt.Valid {
		date = transaction.ScannedAt.Time
	}
	info := &graphql.LatestTransactionInfo{
		Name:      transaction.R.TrackAction.Name,
		CreatedAt: date,
	}

	return info, nil
}
