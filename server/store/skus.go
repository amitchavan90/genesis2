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

// SKUFactory creates skus
func SKUFactory() *db.StockKeepingUnit {
	u := &db.StockKeepingUnit{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Name:        faker.Lorem().Word(),
		Code:        faker.Numerify("L280####"),
		Description: faker.Lorem().Sentence(faker.Number().NumberInt(1)),
	}
	return u
}

// NewSKUStore handle sku methods
func NewSKUStore(conn *sqlx.DB) *SKU {
	ts := &SKU{conn}
	return ts
}

// SKU for persistence
type SKU struct {
	Conn *sqlx.DB
}

// GetByCode gets an SKU by their code
func (s *SKU) GetByCode(code string) (*db.StockKeepingUnit, error) {
	return db.StockKeepingUnits(db.StockKeepingUnitWhere.Code.EQ(code)).One(s.Conn)
}

// Get skus
func (s *SKU) Get(id uuid.UUID) (*db.StockKeepingUnit, error) {
	return db.FindStockKeepingUnit(s.Conn, id.String())
}

// All skus
func (s *SKU) All() (db.StockKeepingUnitSlice, error) {
	return db.StockKeepingUnits().All(s.Conn)
}

// SearchSelect searchs/selects SKUs
func (s *SKU) SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, db.StockKeepingUnitSlice, error) {
	queries := []qm.QueryMod{}

	// Search
	if search.Search != nil {
		searchText := strings.ToLower(strings.Trim(search.Search.String, " "))

		if len(searchText) > 0 {
			queries = append(queries,
				qm.Where(
					fmt.Sprintf("(LOWER(%s) LIKE ?) OR (LOWER(%s) LIKE ?)",
						db.StockKeepingUnitColumns.Name,
						db.StockKeepingUnitColumns.Code,
					),
					"%"+searchText+"%", "%"+searchText+"%",
				))
		}
	}

	// Filter
	if search.Filter != nil {
		switch *search.Filter {
		case graphql.FilterOptionActive:
			queries = append(queries, db.StockKeepingUnitWhere.Archived.EQ(false))
		case graphql.FilterOptionArchived:
			queries = append(queries, db.StockKeepingUnitWhere.Archived.EQ(true))
		}
	}

	// Get Count
	count, err := db.StockKeepingUnits(queries...).Count(s.Conn)
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
			queries = append(queries, qm.OrderBy(db.StockKeepingUnitColumns.CreatedAt+sortDir))
		case graphql.SortByOptionDateUpdated:
			queries = append(queries, qm.OrderBy(db.StockKeepingUnitColumns.UpdatedAt+sortDir))
		case graphql.SortByOptionAlphabetical:
			queries = append(queries, qm.OrderBy(db.StockKeepingUnitColumns.Name+sortDir))
		}
	} else {
		queries = append(queries, qm.OrderBy(db.StockKeepingUnitColumns.CreatedAt+sortDir))
	}

	// Get list
	queries = append(queries, qm.Limit(limit), qm.Offset(offset))
	records, err := db.StockKeepingUnits(queries...).All(s.Conn)
	if err != nil {
		return count, nil, terror.New(err, "")
	}

	return count, records, nil
}

// GetMany skus
func (s *SKU) GetMany(keys []string) (db.StockKeepingUnitSlice, []error) {
	if len(keys) == 0 {
		return nil, []error{errors.New("no keys provided")}
	}
	records, err := db.StockKeepingUnits(db.StockKeepingUnitWhere.ID.IN(keys)).All(s.Conn)
	if errors.Is(err, sql.ErrNoRows) {
		return []*db.StockKeepingUnit{}, nil
	}
	if err != nil {
		return nil, []error{err}
	}

	result := []*db.StockKeepingUnit{}
	for _, key := range keys {
		var row *db.StockKeepingUnit
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

// Insert skus
func (s *SKU) Insert(record *db.StockKeepingUnit, txes ...*sql.Tx) (*db.StockKeepingUnit, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// InsertCategory skus
func (s *SKU) InsertCategory(cat *db.Category, txes ...*sql.Tx) (*db.Category, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return cat.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return cat, nil
}

// InsertProductCategory skus
func (s *SKU) InsertProductCategory(pcat *db.ProductCategory, txes ...*sql.Tx) (*db.ProductCategory, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return pcat.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return pcat, nil
}

// Update skus
func (s *SKU) Update(record *db.StockKeepingUnit, txes ...*sql.Tx) (*db.StockKeepingUnit, error) {
	record.UpdatedAt = time.Now()
	_, err := record.Update(s.Conn, boil.Infer())
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// Archive will archive skus
func (s *SKU) Archive(id uuid.UUID, txes ...*sql.Tx) (*db.StockKeepingUnit, error) {
	u, err := db.FindStockKeepingUnit(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if u.Archived {
		return u, nil
	}

	u.Archived = true
	u.ArchivedAt = null.TimeFrom(time.Now())
	_, err = u.Update(s.Conn, boil.Whitelist(db.StockKeepingUnitColumns.Archived, db.StockKeepingUnitColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Unarchive will unarchive skus
func (s *SKU) Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.StockKeepingUnit, error) {
	u, err := db.FindStockKeepingUnit(s.Conn, id.String())
	if err != nil {
		return nil, terror.New(err, "")
	}

	if !u.Archived {
		return u, nil
	}

	u.Archived = false
	u.ArchivedAt = null.TimeFromPtr(nil)
	_, err = u.Update(s.Conn, boil.Whitelist(db.StockKeepingUnitColumns.Archived, db.StockKeepingUnitColumns.ArchivedAt))
	if err != nil {
		return nil, terror.New(err, "")
	}
	return u, nil
}

// Count gives the amount of skus
func (s *SKU) Count() (int64, error) {
	return db.StockKeepingUnits().Count(s.Conn)
}

// ProductCount returns the count of products in the order
func (s *SKU) ProductCount(record *db.StockKeepingUnit) (int64, error) {
	return db.Products(db.ProductWhere.SkuID.EQ(null.StringFrom(record.ID))).Count(s.Conn)
}

// HasClones checks if sku has any clones
func (s *SKU) HasClones(record *db.StockKeepingUnit) (bool, error) {
	return record.CloneParentStockKeepingUnits().Exists(s.Conn)
}

// GetClones returns all clones of a sku
func (s *SKU) GetClones(id string) (db.StockKeepingUnitSlice, error) {
	return db.StockKeepingUnits(db.StockKeepingUnitWhere.CloneParentID.EQ(null.StringFrom(id))).All(s.Conn)
}

// GetContent returns SKU content (string pairs)
func (s *SKU) GetContent(sku *db.StockKeepingUnit, contentType string) (db.StockKeepingUnitContentSlice, error) {
	return sku.SkuStockKeepingUnitContents(db.StockKeepingUnitContentWhere.ContentType.EQ(contentType)).All(s.Conn)
}

// GetCategories skus by skuID
func (s *SKU) GetCategories(skuID string, txes ...*sql.Tx) (db.CategorySlice, error) {
	dat, err := db.Categories(db.CategoryWhere.SkuID.EQ(skuID)).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// GetProductCategories skus by skuID
func (s *SKU) GetProductCategories(skuID string, txes ...*sql.Tx) (db.ProductCategorySlice, error) {
	dat, err := db.ProductCategories(db.ProductCategoryWhere.SkuID.EQ(skuID)).All(s.Conn)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

// UpdateContent clears and adds new SKU content (string pairs)
func (s *SKU) UpdateContent(sku *db.StockKeepingUnit, input []*graphql.SKUContentInput, contentType string) error {
	// delete old content
	_, err := sku.SkuStockKeepingUnitContents(db.StockKeepingUnitContentWhere.ContentType.EQ(contentType)).DeleteAll(s.Conn)
	if err != nil {
		return nil
	}

	// add new content
	content := []*db.StockKeepingUnitContent{}
	for _, c := range input {
		content = append(content,
			&db.StockKeepingUnitContent{
				Title:       c.Title,
				Content:     c.Content,
				ContentType: contentType,
			})
	}
	return sku.AddSkuStockKeepingUnitContents(s.Conn, true, content...)
}

// GetPhotos return SKU photos
func (s *SKU) GetPhotos(sku *db.StockKeepingUnit) (db.StockKeepingUnitPhotoSlice, error) {
	return sku.SkuStockKeepingUnitPhotos(
		qm.Load(db.StockKeepingUnitPhotoRels.Photo),
		qm.OrderBy(db.StockKeepingUnitPhotoColumns.SortIndex),
	).All(s.Conn)
}

// UpdatePhotos removes/adds photos to match input
func (s *SKU) UpdatePhotos(sku *db.StockKeepingUnit, blobIDs []string) error {
	// deleted last photos?
	if len(blobIDs) == 1 && blobIDs[0] == "" {
		blobIDs = []string{}
	}

	// get current photos
	photos, err := sku.SkuStockKeepingUnitPhotos(
		qm.Load(db.StockKeepingUnitPhotoRels.Photo),
		qm.OrderBy(db.StockKeepingUnitPhotoColumns.SortIndex),
	).All(s.Conn)
	if err != nil {
		return err
	}

	// delete photos?
	for _, photo := range photos {
		removePhoto := true
		for _, fileID := range blobIDs {
			if photo.PhotoID == fileID {
				removePhoto = false
				break
			}
		}
		if !removePhoto {
			continue
		}

		blobID := photo.PhotoID

		_, err = photo.Delete(s.Conn)
		if err != nil {
			return err
		}

		// delete blob
		_, err := db.Blobs(db.BlobWhere.ID.EQ(blobID)).DeleteAll(s.Conn)
		if err != nil {
			return err
		}
	}

	// add new photos?
	for _, fileID := range blobIDs {
		newPhoto := true
		for _, photo := range photos {
			if photo.PhotoID == fileID {
				newPhoto = false
				break
			}
		}
		if !newPhoto {
			continue
		}

		photo := &db.StockKeepingUnitPhoto{
			SkuID:   sku.ID,
			PhotoID: fileID,
		}

		err := photo.Insert(s.Conn, boil.Infer())
		if err != nil {
			return err
		}
	}

	// update sort indices
	photos, err = sku.SkuStockKeepingUnitPhotos().All(s.Conn)
	if err != nil {
		return err
	}

	for i, fileID := range blobIDs {
		for _, photo := range photos {
			if photo.PhotoID == fileID {
				if photo.SortIndex == i {
					break
				}

				photo.SortIndex = i
				_, err := photo.Update(s.Conn, boil.Whitelist(db.StockKeepingUnitPhotoColumns.SortIndex))
				if err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}
