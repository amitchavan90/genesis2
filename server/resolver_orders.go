package genesis

import (
	"context"
	"encoding/json"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"genesis/store"
	"time"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

////////////////
//  Resolver  //
////////////////

// Order resolver
func (r *Resolver) Order() graphql.OrderResolver {
	return &orderResolver{r}
}

type orderResolver struct{ *Resolver }

func (r *orderResolver) ProductCount(ctx context.Context, obj *db.Order) (int, error) {
	count, err := r.OrderStore.ProductCount(obj)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return int(count), nil
}

func (r *orderResolver) Sku(ctx context.Context, obj *db.Order) (*db.StockKeepingUnit, error) {
	skuID, err := r.OrderStore.SkuID(obj)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	if skuID == nil {
		return nil, nil
	}
	skuUUID, err := uuid.FromString(*skuID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := SKULoaderFromContext(ctx, skuUUID)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Order(ctx context.Context, code string) (*db.Order, error) {
	order, err := r.OrderStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get order")
	}

	return order, nil
}

func (r *queryResolver) Orders(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.OrderResult, error) {
	total, orders, err := r.OrderStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list order")
	}

	result := &graphql.OrderResult{
		Orders: orders,
		Total:  int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// OrderCreate creates an order
func (r *mutationResolver) OrderCreate(ctx context.Context, input graphql.CreateOrder) (*db.Order, error) {
	if input.Quantity <= 0 || input.Quantity > 10000 {
		return nil, terror.New(fmt.Errorf("invalid order quantity (%d)", input.Quantity), "")
	}

	// Get SKU
	sku := &db.StockKeepingUnit{}
	sku = nil
	skuID := null.StringFromPtr(nil)
	if input.SkuID != nil && input.SkuID.Valid {
		skuUUID, err := uuid.FromString(input.SkuID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}

		dat, err := r.SKUStore.Get(skuUUID)
		if err != nil {
			return nil, terror.New(err, "get sku")
		}
		sku = dat
		skuID = null.StringFrom(sku.ID)

		input.IsAppBound = sku.IsAppBound
	}

	// Get Order count (for Order Code)
	count, err := r.OrderStore.Count()
	if err != nil {
		return nil, terror.New(err, "create order")
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

	// Setup Order
	u := &db.Order{
		Code:        fmt.Sprintf("N%05d", count+1),
		IsAppBound:  input.IsAppBound,
		CreatedByID: user.ID,
	}

	var created *db.Order

	// Get Product count (for Product Code)
	productCount, err := r.ProductStore.Count()
	if err != nil {
		return nil, terror.New(err, "create product for order")
	}

	// start transaction
	// all tx are the same, so just using any store's begintx to avoid reimplement (being lazy)
	tx, err := r.ManifestStore.BeginTransaction()
	if err != nil {
		return nil, terror.New(err, "begin db tx")
	}
	defer tx.Rollback()

	// Create order
	created, err = r.OrderStore.Insert(u, tx)
	if err != nil {
		return nil, terror.New(err, "create order")
	}

	// Create transaction (contract creation)
	action, err := r.TrackActionStore.GetByCode(trackContractCreated)
	if err != nil {
		return nil, terror.New(err, "invalid track action (contract creation)")
	}
	// manifest json, to be published publicly
	now := time.Now()
	mj := store.ManifestLineJSON{
		Time:            &now,
		TractActionName: &action.Name,
		EntityName:      &createdByName,
	}
	mjs, err := json.Marshal(mj)
	if err != nil {
		return nil, terror.New(err, "json marshal fail")
	}
	t := &db.Transaction{
		TrackActionID:    action.ID,
		CreatedByID:      null.StringFrom(user.ID),
		CreatedByName:    createdByName,
		ManifestLineJSON: null.StringFrom(string(mjs)),
	}
	if !action.Blockchain {
		t.TransactionHash = null.StringFrom("-")
	}
	t, err = r.TransactionStore.Insert(t, tx)
	if err != nil {
		return nil, terror.New(err, "create product for order")
	}

	// Create products
	for i := 0; i < input.Quantity; i++ {
		productCode := fmt.Sprintf("P%05d", productCount+1)

		// Create Product
		p := &db.Product{
			Code:        productCode,
			CreatedByID: user.ID,
			OrderID:     null.StringFrom(created.ID),
			SkuID:       skuID,
		}

		if input.ContractID != nil {
			p.ContractID = *input.ContractID
		}

		if sku != nil {
			p.SkuID = null.StringFrom(sku.ID)
			p.Description = sku.Description
			p.IsBeef = sku.IsBeef
			p.IsAppBound = sku.IsAppBound
			p.IsPointBound = sku.IsPointBound
			p.LoyaltyPoints = sku.LoyaltyPoints
		}

		_, err = r.ProductStore.Insert(p, tx)
		if err != nil {
			return nil, terror.New(err, "create product for order")
		}

		// Attach transaction
		_, err = r.TransactionStore.InsertByProduct(p, action, user, createdByName, nil, tx)
		if err != nil {
			return nil, terror.New(err, "create product for order")
		}

		productCount++
	}

	// commit to db
	err = tx.Commit()
	if err != nil {
		return nil, terror.New(err, "commit create order")
	}

	r.RecordUserActivity(ctx, "Created Order", graphql.ObjectTypeOrder, &created.ID, &created.Code)

	return created, nil
}

// OrderUpdate updates a order
func (r *mutationResolver) OrderUpdate(ctx context.Context, id string, input graphql.UpdateOrder) (*db.Order, error) {
	// Get Order
	orderUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.OrderStore.Get(orderUUID)
	if err != nil {
		return nil, terror.New(err, "update order")
	}

	// NOTE: Orders don't have any columns that can be updated
	// 		 This method is being left here just-in-case it does in the future

	// // Check archived state
	// if u.Archived {
	// 	return nil, terror.New(ErrArchived, "cannot update archived order")
	// }

	// // Update Order
	// updated, err := r.OrderStore.Update(u)
	// if err != nil {
	// 	return nil, terror.New(err, "update order")
	// }

	// r.RecordUserActivity(ctx, "Updated Order", graphql.ObjectTypeOrder, &updated.ID, &updated.Code)

	return u, nil
}

// OrderArchive archives an order
func (r *mutationResolver) OrderArchive(ctx context.Context, id string) (*db.Order, error) {
	// Get Order
	orderUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Order
	updated, err := r.OrderStore.Archive(orderUUID)
	if err != nil {
		return nil, terror.New(err, "update order")
	}

	r.RecordUserActivity(ctx, "Archived Order", graphql.ObjectTypeOrder, &updated.ID, &updated.Code)

	return updated, nil
}

// OrderUnarchive unarchives an order
func (r *mutationResolver) OrderUnarchive(ctx context.Context, id string) (*db.Order, error) {
	// Get Order
	orderUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Order
	updated, err := r.OrderStore.Unarchive(orderUUID)
	if err != nil {
		return nil, terror.New(err, "update order")
	}

	r.RecordUserActivity(ctx, "Unarchived Order", graphql.ObjectTypeOrder, &updated.ID, &updated.Code)

	return updated, nil
}

// OrderBatchAction attempts to do an action of each order
func (r *mutationResolver) OrderBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypeOrder, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.OrderArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive order")
			}
		case graphql.ActionUnarchive:
			_, err := r.OrderUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive order")
			}

		}
	}

	return true, nil
}
