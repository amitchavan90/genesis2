package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
)

////////////////
//  Resolver  //
////////////////

// UserLoyaltyActivity resolver
func (r *Resolver) UserLoyaltyActivity() graphql.UserLoyaltyActivityResolver {
	return &userLoyaltyActivityResolver{r}
}

type userLoyaltyActivityResolver struct{ *Resolver }

func (r *userLoyaltyActivityResolver) User(ctx context.Context, obj *db.UserLoyaltyActivity) (*db.User, error) {
	userUUID, err := uuid.FromString(obj.UserID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := UserLoaderFromContext(ctx, userUUID)
	if err != nil {
		return nil, terror.New(err, "get user")
	}
	return result, nil
}

func (r *userLoyaltyActivityResolver) Product(ctx context.Context, obj *db.UserLoyaltyActivity) (*db.Product, error) {
	if !obj.ProductID.Valid {
		return nil, nil
	}

	productUUID, err := uuid.FromString(obj.ProductID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := ProductLoaderFromContext(ctx, productUUID)
	if err != nil {
		return nil, terror.New(err, "get product")
	}
	return result, nil
}

func (r *userResolver) LoyaltyPoints(ctx context.Context, obj *db.User) (int, error) {
	userUUID, err := uuid.FromString(obj.ID)
	if err != nil {
		return 0, terror.ErrParse
	}

	points, err := r.LoyaltyStore.TotalPoints(ctx, userUUID)
	if err != nil {
		return 0, terror.New(err, "get loyalty points")
	}

	return points, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) GetLoyaltyActivity(ctx context.Context, userID string) ([]*db.UserLoyaltyActivity, error) {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	activity, err := r.LoyaltyStore.UserActivity(userUUID)
	if err != nil {
		return nil, terror.New(err, "get loyalty activity")
	}

	return activity, nil
}
