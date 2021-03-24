package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
)

////////////////
//  Resolver  //
////////////////

// UserPurchaseActivity resolver
func (r *Resolver) UserPurchaseActivity() graphql.UserPurchaseActivityResolver {
	return &userPurchaseActivityResolver{r}
}

type userPurchaseActivityResolver struct{ *Resolver }

func (r *userPurchaseActivityResolver) Product(ctx context.Context, obj *db.UserPurchaseActivity) (*db.Product, error) {
	result, err := r.UserPurchaseActivityStore.GetProduct(obj.ProductID.String)
	if err != nil {
		return nil, terror.New(err, "get product")
	}
	return result, nil
}

func (r *userPurchaseActivityResolver) User(ctx context.Context, obj *db.UserPurchaseActivity) (*db.User, error) {
	result, err := r.UserPurchaseActivityStore.GetUser(obj.UserID)
	if err != nil {
		return nil, terror.New(err, "get user")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserPurchaseActivity(ctx context.Context, id *string) (*db.UserPurchaseActivity, error) {
	purchaseUUID, err := uuid.FromString(*id)
	purchase, err := r.UserPurchaseActivityStore.Get(purchaseUUID)
	if err != nil {
		return nil, terror.New(err, "get purchase")
	}
	return purchase, nil
}

func (r *queryResolver) UserPurchaseActivities(ctx context.Context, search graphql.SearchFilter, limit int, offset int, userID *string) (*graphql.UserPurchaseActivityResult, error) {
	total, purchases, err := r.UserPurchaseActivityStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list purchase")
	}

	result := &graphql.UserPurchaseActivityResult{
		UserPurchaseActivities: purchases,
		Total:                  int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// UserPurchaseActivityCreate creates an purchase
func (r *mutationResolver) UserPurchaseActivityCreate(ctx context.Context, input graphql.UpdateUserPurchaseActivity) (*db.UserPurchaseActivity, error) {
	// Get UserPurchaseActivity count (for UserPurchaseActivity Code)
	count, err := r.UserPurchaseActivityStore.Count()
	if err != nil {
		return nil, terror.New(err, "create purchase: Error while fetching purchase count from db")
	}

	// Create UserPurchaseActivity
	t := &db.UserPurchaseActivity{
		Code: fmt.Sprintf("P%05d", count),
	}

	purchaseID, _ := uuid.NewV4()
	t.ID = purchaseID.String()

	if input.ProductID == nil {
		return nil, terror.New(terror.ErrParse, "create purchase: product id is not provided")
	}

	productUUID, _ := uuid.FromString(input.ProductID.String)

	// Get product
	product, err := r.ProductStore.Get(productUUID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create purchase: no product found with given id")
	}

	// get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create userTask: Error while fetching user")
	}

	t.UserID = userID.String()
	t.ProductID = *input.ProductID
	t.LoyaltyPoints = product.LoyaltyPoints

	created, err := r.UserPurchaseActivityStore.Insert(t)
	if err != nil {
		return nil, terror.New(err, "create user purchase")
	}

	// r.RecordUserActivity(ctx, "Created USer Purchase", graphql.ObjectTypeSku, &created.ID, &created.Code)

	return created, nil
}

// UserPurchaseActivityUpdate updates a purchase
func (r *mutationResolver) UserPurchaseActivityUpdate(ctx context.Context, id string, input graphql.UpdateUserPurchaseActivity) (*db.UserPurchaseActivity, error) {
	// Get UserPurchaseActivity
	purchaseUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	t, err := r.UserPurchaseActivityStore.Get(purchaseUUID)
	if err != nil {
		return nil, terror.New(err, "update purchase")
	}

	updated, err := r.UserPurchaseActivityStore.Update(t)
	if err != nil {
		return nil, terror.New(err, "update purchase")
	}

	return updated, nil
}
