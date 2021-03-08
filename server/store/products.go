package store

import (
	"database/sql"
	"errors"
	"fmt"
	"genesis/db"
	"genesis/graphql"
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

// ProductFactory creates products
func ProductFactory() *db.Product {
	u := &db.Product{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Code:        faker.Numerify("P#####"),
		Description: faker.Lorem().Sentence(10),
	}
	return u
}

// NewProductStore handle product methods
func NewProductStore(conn *sqlx.DB) *Product {
	ts := &Product{conn}
	return ts
}

// Product for persistence
type Product struct {
	Conn *sqlx.DB
}

// GetByCode gets a Product by their code
func (s *Product) GetByCode(code string, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := db.Products(db.ProductWhere.Code.EQ(code)).One(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetByRegisterID gets a Product by their register id
func (s *Product) GetByRegisterID(id uuid.UUID, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := db.Products(db.ProductWhere.RegisterID.EQ(id.String())).One(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// Get product
func (s *Product) Get(id uuid.UUID, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := db.FindProduct(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// All products
func (s *Product) All(txes ...*sqlx.Tx) (db.ProductSlice, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := db.Products().All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// SearchSelect searchs/selects Products
func (s *Product) SearchSelect(
	search graphql.SearchFilter,
	limit int,
	offset int,
	cartonID null.String,
	orderID null.String,
	skuID null.String,
	distributorID null.String,
	contractID null.String,
	trackActionID null.String,
) (int64, []*db.Product, error) {
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
						db.ProductColumns.Code,
						db.ProductColumns.Description,
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
			queries = append(queries, db.ProductWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.ProductWhere.Archived.EQ(true))
		case graphql.FilterOptionProductWithoutCarton:
			queries = append(queries, db.ProductWhere.CartonID.EQ(null.StringFromPtr(nil)))
		case graphql.FilterOptionProductWithoutOrder:
			queries = append(queries, db.ProductWhere.OrderID.EQ(null.StringFromPtr(nil)))
		case graphql.FilterOptionProductWithoutSku:
			queries = append(queries, db.ProductWhere.SkuID.EQ(null.StringFromPtr(nil)))
		}
	}

	if cartonID.Valid {
		queries = append(queries, db.ProductWhere.CartonID.EQ(cartonID))
	}
	if orderID.Valid {
		queries = append(queries, db.ProductWhere.OrderID.EQ(orderID))
	}
	if skuID.Valid {
		queries = append(queries, db.ProductWhere.SkuID.EQ(skuID))
	}
	if distributorID.Valid {
		queries = append(queries, db.ProductWhere.DistributorID.EQ(distributorID))
	}
	if contractID.Valid {
		queries = append(queries, db.ProductWhere.ContractID.EQ(contractID))
	}
	if trackActionID.Valid {
		// using combination of product_latest_transactions view and the inner join lateral
		// to complete this complex search
		// which is to search by the product's last track action
		// it will not search product's previous track action
		queries = append(
			queries,
			qm.InnerJoin(
				`LATERAL
				(
					SELECT product_id, track_action_id FROM product_latest_transactions
					WHERE track_action_id = ?
				) t1 ON products.id = t1.product_id
				`,
				trackActionID.String,
			),
		)
	}

	// Get Count
	count, err := db.Products(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.ProductColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.ProductColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.ProductColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.ProductColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Products(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany products
func (s *Product) GetMany(keys []string) (db.ProductSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Products(db.ProductWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Product{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Product{}
	for _, key := range keys {
		var row *db.Product
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

// GetManyByCartonID product
func (s *Product) GetManyByCartonID(cartonID uuid.UUID, txes ...*sqlx.Tx) (db.ProductSlice, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := db.Products(db.ProductWhere.CartonID.EQ(null.StringFrom(cartonID.String()))).All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetManyByPalletID product
func (s *Product) GetManyByPalletID(palletID uuid.UUID, txes ...*sqlx.Tx) (db.ProductSlice, error) {
	conn := prepConn(s.Conn, txes...)

	cartons, err := db.Cartons(db.CartonWhere.PalletID.EQ(null.StringFrom(palletID.String()))).All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	products := db.ProductSlice{}
	for _, carton := range cartons {
		p, err := db.Products(db.ProductWhere.CartonID.EQ(null.StringFrom(carton.ID))).All(conn)
		if err != nil {
			return nil, terror.New(err, "")
		}
		products = append(products, p...)
	}

	return products, nil
}

// GetManyByContainerID product
func (s *Product) GetManyByContainerID(containerID uuid.UUID, txes ...*sqlx.Tx) (db.ProductSlice, error) {
	conn := prepConn(s.Conn, txes...)

	pallets, err := db.Pallets(db.PalletWhere.ContainerID.EQ(null.StringFrom(containerID.String()))).All(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}

	products := db.ProductSlice{}
	for _, pallet := range pallets {
		cartons, err := db.Cartons(db.CartonWhere.PalletID.EQ(null.StringFrom(pallet.ID))).All(conn)
		if err != nil {
			return nil, terror.New(err, "")
		}

		for _, carton := range cartons {
			p, err := db.Products(db.ProductWhere.CartonID.EQ(null.StringFrom(carton.ID))).All(conn)
			if err != nil {
				return nil, terror.New(err, "")
			}
			products = append(products, p...)
		}
	}

	return products, nil
}

// Insert product
func (s *Product) Insert(record *db.Product, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	err := record.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update product
func (s *Product) Update(record *db.Product, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	record.UpdatedAt = time.Now()
	_, err := record.Update(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive products + remove them from their cartons
func (s *Product) Archive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	u, err := db.FindProduct(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.CartonID = null.StringFromPtr(nil)
	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(conn, boil.Whitelist(db.ProductColumns.Archived, db.ProductColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive products
func (s *Product) Unarchive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Product, error) {
	conn := prepConn(s.Conn, txes...)

	u, err := db.FindProduct(conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(conn, boil.Whitelist(db.ProductColumns.Archived, db.ProductColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of products
func (s *Product) Count(txes ...*sqlx.Tx) (int64, error) {
	conn := prepConn(s.Conn, txes...)

	i, err := db.Products().Count(conn)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return i, nil
}

// Registered checks if the product has been registered or not
func (s *Product) Registered(record *db.Product, txes ...*sqlx.Tx) (bool, error) {
	conn := prepConn(s.Conn, txes...)

	dat, err := record.UserLoyaltyActivities().Exists(conn)
	if err != nil {
		return false, terror.New(err, "")
	}
	return dat, nil
}

// RegisteredBy returns the user that registered the product
func (s *Product) RegisteredBy(record *db.Product, txes ...*sqlx.Tx) (*db.User, error) {
	conn := prepConn(s.Conn, txes...)

	exists, err := record.UserLoyaltyActivities(qm.Load(db.UserLoyaltyActivityRels.User)).Exists(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	if !exists {
		return nil, nil
	}
	activity, err := record.UserLoyaltyActivities(qm.Load(db.UserLoyaltyActivityRels.User)).One(conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return activity.R.User, nil
}

// LatestTrackAction returns the latest transaction track action name
func (s *Product) LatestTrackAction(record *db.Product, txes ...*sqlx.Tx) (*graphql.LatestTransactionInfo, error) {
	conn := prepConn(s.Conn, txes...)

	transaction, err := db.Transactions(
		db.TransactionWhere.ProductID.EQ(null.StringFrom(record.ID)),
		qm.OrderBy(
			fmt.Sprintf("COALESCE(%s,%s) DESC",
				db.TransactionColumns.ScannedAt,
				db.TransactionColumns.CreatedAt,
			),
		),
		qm.Load(
			db.TransactionRels.TrackAction,
		),
	).One(conn)
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
