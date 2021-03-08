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

// OrderFactory creates orders
func OrderFactory() *db.Order {
	u := &db.Order{
		ID:   uuid.Must(uuid.NewV4()).String(),
		Code: faker.Numerify("N#####"),
	}
	return u
}

// NewOrderStore handle order methods
func NewOrderStore(conn *sqlx.DB) *Order {
	ts := &Order{conn}
	return ts
}

// Order for persistence
type Order struct {
	Conn *sqlx.DB
}

// GetByCode gets a Order by their code
func (s *Order) GetByCode(code string) (*db.Order, error) {
	return db.Orders(db.OrderWhere.Code.EQ(code)).One(s.Conn)
}

// Get order
func (s *Order) Get(id uuid.UUID) (*db.Order, error) {
	return db.FindOrder(s.Conn, id.String())
}

// All orders
func (s *Order) All() (db.OrderSlice, error) {
	return db.Orders().All(s.Conn)
}

// SearchSelect searchs/selects Orders
func (s *Order) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Order, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("LOWER(%s) LIKE ?", db.OrderColumns.Code),
					"%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.OrderWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.OrderWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.Orders(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.OrderColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.OrderColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.OrderColumns.Code+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.OrderColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.Orders(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany orders
func (s *Order) GetMany(keys []string) (db.OrderSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.Orders(db.OrderWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.Order{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.Order{}
	for _, key := range keys {
		var row *db.Order
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

// Insert order
func (s *Order) Insert(record *db.Order, txes ...*sqlx.Tx) (*db.Order, error) {
	conn := prepConn(s.Conn, txes...)

	err := record.Insert(conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Update order
func (s *Order) Update(record *db.Order, txes ...*sql.Tx) (*db.Order, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive orders
func (s *Order) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Order, error) {
	u, err := db.FindOrder(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.OrderColumns.Archived, db.OrderColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive orders
func (s *Order) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Order, error) {
	u, err := db.FindOrder(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.OrderColumns.Archived, db.OrderColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of orders
func (s *Order) Count() (int64, error) {
	return db.Orders().Count(s.Conn)
}

// ProductCount returns the count of products in the order
func (s *Order) ProductCount(record *db.Order) (int64, error) {
	return db.Products(db.ProductWhere.OrderID.EQ(null.StringFrom(record.ID))).Count(s.Conn)
}

// SkuID returns the sku id related to the order based on it's products (returns nil if not all products are from the same sku)
func (s *Order) SkuID(record *db.Order) (*string, error) {
	products, err := db.Products(
		db.ProductWhere.OrderID.EQ(null.StringFrom(record.ID)),
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

// Products returns the products in an order (w/ sku info)
func (s *Order) Products(record *db.Order) (db.ProductSlice, error) {
	return db.Products(
		db.ProductWhere.OrderID.EQ(null.StringFrom(record.ID)),
		qm.Load(
			db.ProductRels.Sku,
			qm.Select(
				db.StockKeepingUnitColumns.ID,
				db.StockKeepingUnitColumns.Code,
				db.StockKeepingUnitColumns.Name,
			),
		),
	).All(s.Conn)
}
