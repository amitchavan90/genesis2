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

// CartonFactory creates cartons
func CartonFactory() *db.Carton {
	u := &db.Carton{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Code:        faker.Numerify("CAR#####"),
		Description: faker.Lorem().Sentence(15),
	}
	return u
}

// NewCartonStore handle carton methods
func NewCartonStore(conn *sqlx.DB) *Carton {
	ts := &Carton{conn}
	return ts
}

// Carton for persistence
type Carton struct {
	Conn *sqlx.DB
}

// Get cartons
func (s *Carton) Get(id uuid.UUID) (*db.Carton, error) {
	dat, err := db.FindCarton(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByCode gets a Carton by their code
func (s *Carton) GetByCode(code string) (*db.Carton, error) {
	dat, err := db.Cartons(db.CartonWhere.Code.EQ(code)).One(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetManyByPalletID gets cartons by pallet id
func (s *Carton) GetManyByPalletID(palletID uuid.UUID) (db.CartonSlice, error) {
	dat, err := db.Cartons(db.CartonWhere.PalletID.EQ(null.StringFrom(palletID.String()))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetManyByContainerID get cartons by container id
func (s *Carton) GetManyByContainerID(containerID uuid.UUID) (db.CartonSlice, error) {
	pallet, err := db.Pallets(db.PalletWhere.ContainerID.EQ(null.StringFrom(containerID.String()))).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	cartons := db.CartonSlice{}
	for _, pallet := range pallet {
		c, err := db.Cartons(db.CartonWhere.PalletID.EQ(null.StringFrom(pallet.ID))).All(s.Conn)
		if err != nil {
			return nil, terror.New(err, "")
		}
		cartons = append(cartons, c...)
	}

	return cartons, nil
}

// All cartons
func (s *Carton) All() (db.CartonSlice, error) {
	dat, err := db.Cartons().All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// SearchSelect searchs/selects Cartons
func (s *Carton) SearchSelect(
	search graphql.SearchFilter,
	limit int,
	offset int,
	palletID null.String,
	trackActionID null.String,
) (int64, []*db.Carton, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(
				queries,
				qm.Where(
					fmt.Sprintf(
						"(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.CartonColumns.Code,
						db.CartonColumns.Description,
					),
					"%"+searchText+"%", "%"+searchText+"%",
				),
			)
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.CartonWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.CartonWhere.Archived.EQ(true))
		case graphql.FilterOptionCartonWithoutPallet:
			queries = append(queries, db.CartonWhere.PalletID.EQ(null.StringFromPtr(nil)))
		}
	}

	if palletID.Valid {
		queries = append(queries, db.CartonWhere.PalletID.EQ(palletID))
	}
	if trackActionID.Valid {
		// using combination of carton_latest_transactions view and the inner join lateral
		// to complete this complex search
		// which is to search by the carton's last track action
		// it will not search carton's previous track action
		queries = append(
			queries,
			qm.InnerJoin(
				`LATERAL
				(
					SELECT carton_id, track_action_id FROM carton_latest_transactions
					WHERE track_action_id = ?
				) t1 ON cartons.id = t1.carton_id
				`,
				trackActionID.String,
			),
		)
	}

	// Get Count
	count, err := db.Cartons(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.CartonColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.CartonColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.CartonColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.CartonColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Cartons(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany cartons
func (s *Carton) GetMany(keys []string) (db.CartonSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Cartons(db.CartonWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Carton{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Carton{}
	for _, key := range keys {
		var row *db.Carton
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

// GetRange returns all cartons between two carton codes
func (s *Carton) GetRange(from string, to string) (db.CartonSlice, error) {
	start, err := strconv.Atoi(strings.Replace(from, "CAR", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}
	end, err := strconv.Atoi(strings.Replace(to, "CAR", "", 1))
	if err != nil {
		return nil, terror.New(err, "")
	}

	if start < 0 || end < start {
		return nil, terror.New(err, "")
	}

	keys := []string{}
	for i := start; i <= end; i++ {
		keys = append(keys, fmt.Sprintf("CAR%05d", i))
	}

	records, err := db.Cartons(db.CartonWhere.Code.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Carton{}, nil
	}
	if err != nil {
		return nil, terror.New(err, "")
	}

	result := []*db.Carton{}
	for _, key := range keys {
		var row *db.Carton
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

// Insert cartons
func (s *Carton) Insert(record *db.Carton, txes ...*sql.Tx) (*db.Carton, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update cartons
func (s *Carton) Update(record *db.Carton, txes ...*sql.Tx) (*db.Carton, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive cartons
func (s *Carton) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Carton, error) {
	u, err := db.FindCarton(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.CartonColumns.Archived, db.CartonColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive cartons
func (s *Carton) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Carton, error) {
	u, err := db.FindCarton(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.CartonColumns.Archived, db.CartonColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of cartons
func (s *Carton) Count() (int64, error) {
	count, err := db.Cartons().Count(s.Conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return count, nil
}

// ProductCount returns the count of products in the carton
func (s *Carton) ProductCount(record *db.Carton) (int64, error) {
	return db.Products(db.ProductWhere.CartonID.EQ(null.StringFrom(record.ID))).Count(s.Conn)
}

// LatestTrackAction returns the latest transaction track action name
func (s *Carton) LatestTrackAction(record *db.Carton) (*graphql.LatestTransactionInfo, error) {
	transaction, err := db.Transactions(
		db.TransactionWhere.CartonID.EQ(null.StringFrom(record.ID)),
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

// OrderID returns the order id related to the carton based on it's products (returns nil if not all products are from the same order)
func (s *Carton) OrderID(record *db.Carton) (*string, error) {
	products, err := db.Products(
		db.ProductWhere.CartonID.EQ(null.StringFrom(record.ID)),
		qm.Select(db.ProductColumns.OrderID),
		qm.GroupBy(db.ProductColumns.OrderID),
	).All(s.Conn)

	if err != nil {
		return nil, terror.New(fmt.Errorf("invalid arguments"), "")
	}
	if len(products) != 1 {
		return nil, nil
	}
	if !products[0].OrderID.Valid {
		return nil, nil
	}

	return &products[0].OrderID.String, nil
}

// DistributorID returns the distributor id related to the carton based on it's products (returns nil if not all products are from the same distributor)
func (s *Carton) DistributorID(record *db.Carton) (*string, error) {
	products, err := db.Products(
		db.ProductWhere.CartonID.EQ(null.StringFrom(record.ID)),
		qm.Select(db.ProductColumns.DistributorID),
		qm.GroupBy(db.ProductColumns.DistributorID),
	).All(s.Conn)

	if err != nil {
		return nil, terror.New(err, "")
	}
	if len(products) != 1 {
		return nil, nil
	}
	if !products[0].DistributorID.Valid {
		return nil, nil
	}

	return &products[0].DistributorID.String, nil
}

// SkuID returns the sku id related to the carton based on it's products (returns nil if not all products are from the same sku)
func (s *Carton) SkuID(record *db.Carton) (*string, error) {
	products, err := db.Products(
		db.ProductWhere.CartonID.EQ(null.StringFrom(record.ID)),
		qm.Select(db.ProductColumns.SkuID),
		qm.GroupBy(db.ProductColumns.SkuID),
	).All(s.Conn)

	if err != nil {
		return nil, terror.New(err, "")
	}
	if len(products) != 1 {
		return nil, nil
	}
	if !products[0].SkuID.Valid {
		return nil, nil
	}

	return &products[0].SkuID.String, nil
}
