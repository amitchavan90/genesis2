package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

////////////////
//  Resolver  //
////////////////

// Product resolver
func (r *Resolver) Product() graphql.ProductResolver {
	return &productResolver{r}
}

type productResolver struct{ *Resolver }

func (r *productResolver) Carton(ctx context.Context, obj *db.Product) (*db.Carton, error) {
	if !obj.CartonID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.CartonID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := CartonLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get carton")
	}
	return result, nil
}

func (r *productResolver) Order(ctx context.Context, obj *db.Product) (*db.Order, error) {
	if !obj.OrderID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.OrderID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := OrderLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get order")
	}
	return result, nil
}

func (r *productResolver) Contract(ctx context.Context, obj *db.Product) (*db.Contract, error) {
	if !obj.ContractID.Valid {
		return nil, nil
	}

	uuid, err := uuid.FromString(obj.ContractID.String)
	if err != nil {
		return nil, terror.New(err, "get product contract")
	}
	result, err := ContractLoaderFromContext(ctx, uuid)
	if err != nil {
		return nil, terror.New(err, "get product contract")
	}
	return result, nil
}

func (r *productResolver) Sku(ctx context.Context, obj *db.Product) (*db.StockKeepingUnit, error) {
	if !obj.SkuID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.SkuID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := SKULoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

func (r *productResolver) Distributor(ctx context.Context, obj *db.Product) (*db.Distributor, error) {
	if !obj.DistributorID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.DistributorID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := DistributorLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get distributor")
	}
	return result, nil
}

func (r *productResolver) Registered(ctx context.Context, obj *db.Product) (bool, error) {
	registered, err := r.ProductStore.Registered(obj)
	if err != nil {
		return false, terror.New(err, "check product registered")
	}
	return registered, nil
}

func (r *productResolver) RegisteredBy(ctx context.Context, obj *db.Product) (*db.User, error) {
	registeredBy, err := r.ProductStore.RegisteredBy(obj)
	if err != nil {
		return nil, terror.New(err, "get product registered by")
	}
	return registeredBy, nil
}

func (r *productResolver) Transactions(ctx context.Context, obj *db.Product) ([]*db.Transaction, error) {
	// Get product transactions
	productUUID, err := uuid.FromString(obj.ID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	transactions, err := r.TransactionStore.GetByProductID(productUUID)
	if err != nil {
		return nil, terror.New(err, "get product transactions")
	}

	// Omit private transactions if not logged in
	userID, _ := r.Auther.UserIDFromContext(ctx)

	result := db.TransactionSlice{}
	for _, t := range transactions {
		if userID == uuid.Nil && t.R.TrackAction.Private {
			continue
		}

		// public only see little detail
		if userID == uuid.Nil {
			result = append(
				result,
				&db.Transaction{
					ID:                 t.ID,
					TrackActionID:      t.TrackActionID,
					CreatedByName:      t.CreatedByName,
					CreatedAt:          t.CreatedAt,
					ScannedAt:          t.ScannedAt,
					ProductID:          t.ProductID,
					TransactionHash:    t.TransactionHash,
					ManifestID:         t.ManifestID,
					ManifestLineJSON:   t.ManifestLineJSON,
					ManifestLineSha256: t.ManifestLineSha256,
					LocationName:       t.LocationName,
					LocationGeohash:    t.LocationGeohash,
				},
			)
			continue
		}

		result = append(result, t)
	}

	return result, nil
}

func (r *productResolver) LatestTrackAction(ctx context.Context, obj *db.Product) (*graphql.LatestTransactionInfo, error) {
	action, err := r.ProductStore.LatestTrackAction(obj)
	if err != nil {
		return nil, terror.New(err, "get product latest track action")
	}
	return action, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Product(ctx context.Context, code string) (*db.Product, error) {
	product, err := r.ProductStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get product")
	}
	return product, nil

}
func (r *queryResolver) ProductByID(ctx context.Context, id string) (*db.Product, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	product, err := ProductLoaderFromContext(ctx, uuid)
	if err != nil {
		return nil, terror.New(err, "get product")
	}
	return product, nil
}

func (r *queryResolver) Products(
	ctx context.Context,
	search graphql.SearchFilter,
	limit int,
	offset int,
	cartonID *string,
	orderID *string,
	skuID *string,
	distributorID *string,
	contractID *string,
	trackActionID *string,
) (*graphql.ProductResult, error) {
	total, products, err := r.ProductStore.SearchSelect(
		search,
		limit,
		offset,
		null.StringFromPtr(cartonID),
		null.StringFromPtr(orderID),
		null.StringFromPtr(skuID),
		null.StringFromPtr(distributorID),
		null.StringFromPtr(contractID),
		null.StringFromPtr(trackActionID),
	)
	if err != nil {
		return nil, terror.New(err, "list product")
	}

	result := &graphql.ProductResult{
		Products: products,
		Total:    int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// ProductCreate creates an product
func (r *mutationResolver) ProductCreate(ctx context.Context, input graphql.UpdateProduct) (*db.Product, error) {
	// Get Product count (for Product Code)
	count, err := r.ProductStore.Count()
	if err != nil {
		return nil, terror.New(err, "create product")
	}

	// Get user
	user, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}
	createdByName := user.AffiliateOrg.String
	if createdByName == "" {
		createdByName = user.LastName.String
	}

	// Get sku
	if input.SkuID == nil {
		return nil, terror.New(err, "product create: SKU ID is required")
	}

	skuUUID, _ := uuid.FromString(input.SkuID.String)
	sku, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return nil, terror.New(err, "product create: no sku found with give ID")
	}

	// Create Product
	p := &db.Product{
		Code:        fmt.Sprintf("P%05d", count+1),
		CreatedByID: user.ID,
		// Description: input.Description.String,
	}
	if input.Description != nil {
		p.Description = input.Description.String
	}
	if input.CartonID != nil {
		if len(input.CartonID.String) <= 1 {
			p.CartonID = null.StringFromPtr(nil)
		} else {
			p.CartonID = *input.CartonID
		}
	}
	if input.OrderID != nil {
		if len(input.OrderID.String) <= 1 {
			p.OrderID = null.StringFromPtr(nil)
		} else {
			p.OrderID = *input.OrderID
		}
	}
	if input.DistributorID != nil {
		if len(input.DistributorID.String) <= 1 {
			p.DistributorID = null.StringFromPtr(nil)
		} else {
			p.DistributorID = *input.DistributorID
		}
	}
	// if input.SkuID != nil {
	// 	if len(input.SkuID.String) <= 1 {
	// 		p.SkuID = null.StringFromPtr(nil)
	// 	} else {
	// 		p.SkuID = *input.SkuID
	// 	}
	// }
	if input.ContractID != nil {
		if len(input.ContractID.String) <= 1 {
			p.ContractID = null.StringFromPtr(nil)
		} else {
			p.ContractID = *input.ContractID
		}
	}

	if input.LoyaltyPoints != nil {
		p.LoyaltyPoints = input.LoyaltyPoints.Int
	} else {
		p.LoyaltyPoints = sku.LoyaltyPoints
	}

	if input.LoyaltyPointsExpire != nil {
		p.LoyaltyPointsExpire = input.LoyaltyPointsExpire.Time
	}

	p.IsAppBound = sku.IsAppBound
	p.IsPointBound = sku.IsPointBound

	created, err := r.ProductStore.Insert(p)
	if err != nil {
		return nil, terror.New(err, "create product")
	}

	r.RecordUserActivity(ctx, "Created Product", graphql.ObjectTypeProduct, &created.ID, &created.Code)

	// Copy carton transactions
	if p.CartonID.Valid && input.InheritCartonHistory != nil && input.InheritCartonHistory.Bool == true {
		cartonUUID, err := uuid.FromString(p.CartonID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		err = r.CopyCartonTransactionsToProduct(ctx, cartonUUID, p)
		if err != nil {
			return nil, terror.New(err, "create product")
		}
	}

	// Create product transaction (contract creation)
	action, err := r.TrackActionStore.GetByCode(trackContractCreated)
	if err != nil {
		if err != nil {
			return nil, terror.New(err, "invalid track action (contract creation)")
		}
	}

	_, err = r.TransactionStore.InsertByProduct(p, action, user, createdByName, nil)
	if err != nil {
		return nil, terror.New(err, "create product (attach to product)")
	}

	return created, nil
}

// ProductUpdate updates a product
func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input graphql.UpdateProduct) (*db.Product, error) {
	// Get Product
	productUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	p, err := r.ProductStore.Get(productUUID)
	if err != nil {
		return nil, terror.New(err, "update product")
	}

	// Check archived state
	if p.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived product")
	}

	// Update Product
	if input.Description != nil {
		p.Description = input.Description.String
	}
	if input.LoyaltyPoints != nil {
		p.LoyaltyPoints = input.LoyaltyPoints.Int
	}
	if input.LoyaltyPointsExpire != nil {
		p.LoyaltyPointsExpire = input.LoyaltyPointsExpire.Time
	}
	if input.CartonID != nil {
		if len(input.CartonID.String) <= 1 {
			if p.CartonID.Valid {
				// Track Carton Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackRemovedFromCarton,
					ProductIDs:      []string{p.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: removed from carton")
				}
			}

			p.CartonID = null.StringFromPtr(nil)
		} else {
			if p.CartonID.String != input.CartonID.String {
				// Track Carton Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackMovedToCarton,
					ProductIDs:      []string{p.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: moved to carton")
				}
				// Copy Transactions from carton
				if input.InheritCartonHistory != nil && input.InheritCartonHistory.Bool == true {
					cartonUUID, err := uuid.FromString(input.CartonID.String)
					if err != nil {
						return nil, terror.New(terror.ErrParse, "")
					}
					err = r.CopyCartonTransactionsToProduct(ctx, cartonUUID, p)
					if err != nil {
						return nil, terror.New(err, "update product")
					}
				}
			}

			p.CartonID = *input.CartonID
		}
	}
	if input.OrderID != nil {
		if len(input.OrderID.String) <= 1 {
			p.OrderID = null.StringFromPtr(nil)
		} else {
			p.OrderID = *input.OrderID
		}
	}
	if input.DistributorID != nil {
		if len(input.DistributorID.String) <= 1 {
			p.DistributorID = null.StringFromPtr(nil)
		} else {
			p.DistributorID = *input.DistributorID
		}
	}
	if input.SkuID != nil {
		if len(input.SkuID.String) <= 1 {
			p.SkuID = null.StringFromPtr(nil)
		} else {
			p.SkuID = *input.SkuID
		}
	}
	if input.ContractID != nil {
		if len(input.ContractID.String) <= 1 {
			p.ContractID = null.StringFromPtr(nil)
		} else {
			p.ContractID = *input.ContractID
		}
	}

	updated, err := r.ProductStore.Update(p)
	if err != nil {
		return nil, terror.New(err, "update product")
	}

	r.RecordUserActivity(ctx, "Updated Product", graphql.ObjectTypeProduct, &updated.ID, &updated.Code)

	return updated, nil
}

// ProductArchive archives an product
func (r *mutationResolver) ProductArchive(ctx context.Context, id string) (*db.Product, error) {
	// Get Product
	productUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Product
	updated, err := r.ProductStore.Archive(productUUID)
	if err != nil {
		return nil, terror.New(err, "archive product")
	}

	r.RecordUserActivity(ctx, "Archived Product", graphql.ObjectTypeProduct, &updated.ID, &updated.Code)

	return updated, nil
}

// ProductUnarchive unarchives an product
func (r *mutationResolver) ProductUnarchive(ctx context.Context, id string) (*db.Product, error) {
	// Get Product
	productUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Product
	updated, err := r.ProductStore.Unarchive(productUUID)
	if err != nil {
		return nil, terror.New(err, "unarchive product")
	}

	r.RecordUserActivity(ctx, "Unarchived Product", graphql.ObjectTypeProduct, &updated.ID, &updated.Code)

	return updated, nil
}

// ProductBatchAction attempts to do an action of each product
func (r *mutationResolver) ProductBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypeProduct, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.ProductArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive product")
			}
		case graphql.ActionUnarchive:
			_, err := r.ProductUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive product")
			}

		case graphql.ActionDetachFromCarton:
			newID := null.StringFrom("-")
			input := graphql.UpdateProduct{CartonID: &newID}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "detach product from carton")
			}
		case graphql.ActionSetCarton:
			input := graphql.UpdateProduct{CartonID: value.Str, InheritCartonHistory: value.Bool}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "move product to carton")
			}

		case graphql.ActionDetachFromOrder:
			newID := null.StringFrom("-")
			input := graphql.UpdateProduct{OrderID: &newID}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "detach product from order")
			}
		case graphql.ActionSetOrder:
			input := graphql.UpdateProduct{OrderID: value.Str}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "move product to order")
			}

		case graphql.ActionDetachFromSku:
			newID := null.StringFrom("-")
			input := graphql.UpdateProduct{SkuID: &newID}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "remove product sku")
			}
		case graphql.ActionSetSku:
			input := graphql.UpdateProduct{SkuID: value.Str}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "change product sku")
			}

		case graphql.ActionDetachFromDistributor:
			newID := null.StringFrom("-")
			input := graphql.UpdateProduct{DistributorID: &newID}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "remove product distributor")
			}
		case graphql.ActionSetDistributor:
			input := graphql.UpdateProduct{DistributorID: value.Str}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "change product distributor")
			}

		case graphql.ActionDetachFromContract:
			newID := null.StringFrom("-")
			input := graphql.UpdateProduct{ContractID: &newID}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "remove product contract")
			}
		case graphql.ActionSetContract:
			input := graphql.UpdateProduct{ContractID: value.Str}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "change product contract")
			}

		case graphql.ActionSetBonusLoyaltyPoints:
			input := graphql.UpdateProduct{
				LoyaltyPoints:       value.No,
				LoyaltyPointsExpire: value.DateTime,
			}
			_, err := r.ProductUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "change product bonus loyalty points")
			}

		case graphql.ActionInheritCartonHistory:
			if value.Str == nil || !value.Str.Valid {
				return false, nil
			}

			// Get Carton UUID
			cartonUUID, err := uuid.FromString(value.Str.String)
			if err != nil {
				return false, terror.New(terror.ErrParse, "")
			}

			// Get Product
			productUUID, err := uuid.FromString(id)
			if err != nil {
				return false, terror.New(terror.ErrParse, "")
			}
			p, err := r.ProductStore.Get(productUUID)
			if err != nil {
				return false, terror.New(err, "ActionInheritCartonHistory")
			}

			// Copy carton transactions to product
			err = r.CopyCartonTransactionsToProduct(ctx, cartonUUID, p)
			if err != nil {
				return false, terror.New(err, "ActionInheritCartonHistory")
			}

		}
	}

	return true, nil
}
