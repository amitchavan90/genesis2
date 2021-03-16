package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/volatiletech/null"

	"github.com/gofrs/uuid"
)

////////////////
//  Resolver  //
////////////////

// SKU resolver
func (r *Resolver) SKU() graphql.SKUResolver {
	return &skuResolver{r}
}

type skuResolver struct{ *Resolver }

func (r *skuResolver) MasterPlan(ctx context.Context, sku *db.StockKeepingUnit) (*db.Blob, error) {
	if !sku.MasterPlanBlobID.Valid {
		return nil, nil
	}
	blobUUID, err := uuid.FromString(sku.MasterPlanBlobID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	return r.BlobStore.Get(blobUUID)
}
func (r *skuResolver) Video(ctx context.Context, sku *db.StockKeepingUnit) (*db.Blob, error) {
	if !sku.VideoBlobID.Valid {
		return nil, nil
	}
	blobUUID, err := uuid.FromString(sku.VideoBlobID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	return r.BlobStore.Get(blobUUID)
}
func (r *skuResolver) Urls(ctx context.Context, sku *db.StockKeepingUnit) ([]*db.StockKeepingUnitContent, error) {
	content, err := r.SKUStore.GetContent(sku, db.ContentTypeURL)
	if err != nil {
		return nil, terror.New(err, "get sku urls")
	}
	return content, nil
}
func (r *skuResolver) ProductInfo(ctx context.Context, sku *db.StockKeepingUnit) ([]*db.StockKeepingUnitContent, error) {
	content, err := r.SKUStore.GetContent(sku, db.ContentTypeINFO)
	if err != nil {
		return nil, terror.New(err, "get sku product info")
	}
	return content, nil
}
func (r *skuResolver) Photos(ctx context.Context, sku *db.StockKeepingUnit) ([]*db.Blob, error) {
	content, err := r.SKUStore.GetPhotos(sku)
	if err != nil {
		return nil, terror.New(err, "get sku photos")
	}

	photos := []*db.Blob{}
	for _, photo := range content {
		photos = append(photos, photo.R.Photo)
	}

	return photos, nil
}

func (r *skuResolver) ProductCount(ctx context.Context, obj *db.StockKeepingUnit) (int, error) {
	count, err := r.SKUStore.ProductCount(obj)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return int(count), nil
}

func (r *skuResolver) HasClones(ctx context.Context, obj *db.StockKeepingUnit) (bool, error) {
	hasClones, err := r.SKUStore.HasClones(obj)
	if err != nil {
		return false, terror.New(err, "")
	}
	return hasClones, nil
}

func (r *skuResolver) Categories(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.Category, error) {
	dat, err := r.SKUStore.GetCategories(obj.ID)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

func (r *skuResolver) ProductCategories(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.ProductCategory, error) {
	dat, err := r.SKUStore.GetProductCategories(obj.ID)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return dat, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Sku(ctx context.Context, code string) (*db.StockKeepingUnit, error) {
	sku, err := r.SKUStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return sku, nil
}
func (r *queryResolver) SkuByID(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	skuUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(err, "get sku: invalid uuid")
	}
	sku, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return sku, nil
}

func (r *queryResolver) Skus(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.SKUResult, error) {
	total, skus, err := r.SKUStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list sku")
	}

	result := &graphql.SKUResult{
		Skus:  skus,
		Total: int(total),
	}

	return result, nil
}

func (r *queryResolver) SkuCloneTree(ctx context.Context, id string) ([]*graphql.SKUClone, error) {
	cloneTree := []*graphql.SKUClone{}

	// get sku
	skuUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(err, "get sku: invalid uuid")
	}
	sku, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}

	// get clones
	cloneTree, err = r.SkuCloneTreeGet(ctx, sku, cloneTree, 0)
	if err != nil {
		return nil, terror.New(err, "")
	}

	return cloneTree, nil
}

func (r *queryResolver) SkuCloneTreeGet(ctx context.Context, sku *db.StockKeepingUnit, cloneTree []*graphql.SKUClone, depth int) ([]*graphql.SKUClone, error) {
	if !sku.CloneParentID.Valid {
		// original sku - add
		cloneTree = append(cloneTree,
			&graphql.SKUClone{
				Sku:   sku,
				Depth: depth,
			},
		)

		// add clone children
		cloneTree, err := r.getSkuCloneChildren(ctx, sku, cloneTree, depth)
		if err != nil {
			return cloneTree, terror.New(err, "get sku clone children")
		}

		return cloneTree, nil
	}

	// get clone parent
	skuUUID, err := uuid.FromString(sku.CloneParentID.String)
	if err != nil {
		return nil, terror.New(err, "get sku clone parent: invalid uuid")
	}
	parent, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return cloneTree, terror.New(err, "get sku clone parent")
	}

	cloneTree, err = r.SkuCloneTreeGet(ctx, parent, cloneTree, depth-1)
	if err != nil {
		return cloneTree, terror.New(err, "")
	}

	return cloneTree, nil
}
func (r *queryResolver) getSkuCloneChildren(ctx context.Context, sku *db.StockKeepingUnit, cloneTree []*graphql.SKUClone, depth int) ([]*graphql.SKUClone, error) {
	hasChildren, err := r.SKUStore.HasClones(sku)
	if err != nil {
		return cloneTree, terror.New(err, "")
	}
	if !hasChildren {
		return cloneTree, nil
	}

	children, err := r.SKUStore.GetClones(sku.ID)
	if err != nil {
		return cloneTree, terror.New(err, "")
	}

	for _, child := range children {
		// add clone child
		cloneTree = append(cloneTree,
			&graphql.SKUClone{
				Sku:   child,
				Depth: depth + 1,
			},
		)

		// get chilren's children
		cloneTree, err = r.getSkuCloneChildren(ctx, child, cloneTree, depth+1)
		if err != nil {
			return cloneTree, terror.New(err, "")
		}
	}

	return cloneTree, nil
}

///////////////
// Mutations //
///////////////

// SkuCreate creates an sku
func (r *mutationResolver) SkuCreate(ctx context.Context, input graphql.UpdateSku) (*db.StockKeepingUnit, error) {
	// Get SKU count (for SKU Code)
	count, err := r.SKUStore.Count()
	if err != nil {
		return nil, terror.New(err, "create sku")
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	// Create SKU
	u := &db.StockKeepingUnit{
		Code:        fmt.Sprintf("L28%05d", count),
		CreatedByID: userID.String(),
	}

	if input.Name != nil {
		u.Name = input.Name.String
	}
	if input.Brand != nil {
		u.Brand = input.Brand.String
	}
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.Ingredients != nil {
		u.Ingredients = input.Ingredients.String
	}
	if input.LoyaltyPoints != nil {
		u.LoyaltyPoints = input.LoyaltyPoints.Int
	}
	if input.MasterPlanBlobID != nil {
		if input.MasterPlanBlobID.String == "-" {
			u.MasterPlanBlobID = null.StringFromPtr(nil)
		} else {
			u.MasterPlanBlobID = *input.MasterPlanBlobID
		}
	}
	if input.VideoBlobID != nil {
		if input.VideoBlobID.String == "-" {
			u.VideoBlobID = null.StringFromPtr(nil)
		} else {
			u.VideoBlobID = *input.VideoBlobID
		}
	}
	if input.CloneParentID != nil {
		u.CloneParentID = *input.CloneParentID
	}

	u.IsBeef = input.IsBeef.Bool
	u.IsPointSku = input.IsPointSku.Bool
	u.IsAppSku = input.IsAppSku.Bool
	u.Weight = input.Weight.Int
	u.WeightUnit = input.WeightUnit.String
	u.Currency = input.Currency.String
	u.Price = input.Price.Int
	u.PurchasePoints = input.PurchasePoints.Int
	u.LoyaltyPoints = input.LoyaltyPoints.Int

	created, err := r.SKUStore.Insert(u)
	if err != nil {
		return nil, terror.New(err, "create sku")
	}

	// Add Content
	if input.Urls != nil && len(input.Urls) > 0 {
		err := r.SKUStore.UpdateContent(created, input.Urls, db.ContentTypeURL)
		if err != nil {
			return nil, terror.New(err, "create sku")
		}
	}
	if input.ProductInfo != nil && len(input.ProductInfo) > 0 {
		err := r.SKUStore.UpdateContent(created, input.ProductInfo, db.ContentTypeINFO)
		if err != nil {
			return nil, terror.New(err, "create sku")
		}
	}
	if input.PhotoBlobIDs != nil && len(input.PhotoBlobIDs) > 0 {
		err := r.SKUStore.UpdatePhotos(created, input.PhotoBlobIDs)
		if err != nil {
			return nil, terror.New(err, "create sku")
		}
	}

	// Add categories
	if len(input.Categories) >= 0 {
		for i := range input.Categories {
			cat := &db.Category{}
			id, _ := uuid.NewV4()
			cat.ID = id.String()
			cat.SkuID = created.ID
			cat.Name = input.Categories[i].Name
			_, err = r.SKUStore.InsertCategory(cat)
			if err != nil {
				return nil, terror.New(err, "create sku")
			}
		}
	}

	if len(input.ProductCategories) >= 0 {
		for i := range input.ProductCategories {
			pcat := &db.ProductCategory{}
			id, _ := uuid.NewV4()
			pcat.ID = id.String()
			pcat.SkuID = created.ID
			pcat.Name = input.ProductCategories[i].Name
			_, err = r.SKUStore.InsertProductCategory(pcat)
			if err != nil {
				return nil, terror.New(err, "create sku")
			}
		}
	}

	r.RecordUserActivity(ctx, "Created SKU", graphql.ObjectTypeSku, &created.ID, &created.Code)

	return created, nil
}

// SkuUpdate updates a sku
func (r *mutationResolver) SkuUpdate(ctx context.Context, id string, input graphql.UpdateSku) (*db.StockKeepingUnit, error) {
	// Get SKU
	skuUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return nil, terror.New(err, "update sku")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived sku")
	}

	// Update SKU
	if input.Name != nil {
		u.Name = input.Name.String
	}
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.LoyaltyPoints != nil {
		u.LoyaltyPoints = input.LoyaltyPoints.Int
	}
	if input.MasterPlanBlobID != nil {
		if input.MasterPlanBlobID.String == "-" {
			u.MasterPlanBlobID = null.StringFromPtr(nil)
		} else {
			u.MasterPlanBlobID = *input.MasterPlanBlobID
		}
	}
	if input.VideoBlobID != nil {
		if input.VideoBlobID.String == "-" {
			u.VideoBlobID = null.StringFromPtr(nil)
		} else {
			u.VideoBlobID = *input.VideoBlobID
		}
	}

	updated, err := r.SKUStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update sku")
	}

	// Add Content
	if input.Urls != nil && len(input.Urls) > 0 {
		err := r.SKUStore.UpdateContent(updated, input.Urls, db.ContentTypeURL)
		if err != nil {
			return nil, terror.New(err, "update sku")
		}
	}
	if input.ProductInfo != nil && len(input.ProductInfo) > 0 {
		err := r.SKUStore.UpdateContent(updated, input.ProductInfo, db.ContentTypeINFO)
		if err != nil {
			return nil, terror.New(err, "update sku")
		}
	}
	if input.PhotoBlobIDs != nil && len(input.PhotoBlobIDs) > 0 {
		err := r.SKUStore.UpdatePhotos(updated, input.PhotoBlobIDs)
		if err != nil {
			return nil, terror.New(err, "update sku")
		}
	}

	r.RecordUserActivity(ctx, "Updated SKU", graphql.ObjectTypeSku, &updated.ID, &updated.Code)

	return updated, nil
}

// SkuArchive archives an sku
func (r *mutationResolver) SkuArchive(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	// Get SKU
	skuUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive SKU
	updated, err := r.SKUStore.Archive(skuUUID)
	if err != nil {
		return nil, terror.New(err, "update sku")
	}

	r.RecordUserActivity(ctx, "Archived SKU", graphql.ObjectTypeSku, &updated.ID, &updated.Code)

	return updated, nil
}

// SkuUnarchive unarchives an sku
func (r *mutationResolver) SkuUnarchive(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	// Get SKU
	skuUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive SKU
	updated, err := r.SKUStore.Unarchive(skuUUID)
	if err != nil {
		return nil, terror.New(err, "update sku")
	}

	r.RecordUserActivity(ctx, "Unarchived SKU", graphql.ObjectTypeSku, &updated.ID, &updated.Code)

	return updated, nil
}

// SkuBatchAction attempts to do an action of each sku
func (r *mutationResolver) SkuBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypeSku, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.SkuArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive sku")
			}
		case graphql.ActionUnarchive:
			_, err := r.SkuUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive sku")
			}

		}
	}

	return true, nil
}
