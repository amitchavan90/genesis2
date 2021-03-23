package genesis

import (
	"context"
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
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

func (r *userPurchaseActivityResolver) User(ctx context.Context, obj *db.UserPurchaseActivity) (*db.User, error) {
	userUUID, _ := uuid.FromString(obj.ID)
	result, err := r.UserStore.Get(userUUID)
	if err != nil {
		return nil, terror.New(err, "get subtasks")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserPurchaseActivity(ctx context.Context, id *string) (*db.UserPurchaseActivity, error) {
	taskUUID, err := uuid.FromString(*id)
	task, err := r.UserPurchaseActivityStore.Get(taskUUID)
	if err != nil {
		return nil, terror.New(err, "get task")
	}
	return task, nil
}

func (r *queryResolver) UserPurchaseActivities(ctx context.Context, search graphql.SearchFilter, limit int, offset int, userID *string) (*graphql.UserPurchaseActivityResult, error) {
	total, tasks, err := r.UserPurchaseActivityStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list task")
	}

	result := &graphql.UserPurchaseActivityResult{
		UserPurchaseActivities: tasks,
		Total:                  int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// UserPurchaseActivityCreate creates an task
func (r *mutationResolver) UserPurchaseActivityCreate(ctx context.Context, input graphql.UpdateUserPurchaseActivity) (*db.UserPurchaseActivity, error) {
	// Create UserPurchaseActivity
	t := &db.UserPurchaseActivity{}

	purchaseID, _ := uuid.NewV4()
	t.ID = purchaseID.String()

	if input.ProductID == nil {
		return nil, terror.New(terror.ErrParse, "create purchase: product id is not provided")
	}

	productUUID, _ := uuid.FromString(input.ProductID.String)

	// Get product
	product, err := r.ProductStore.Get(productUUID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create task: no product found with given id")
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

// UserPurchaseActivityUpdate updates a task
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
