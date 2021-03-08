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

// RecordUserActivity adds a UserActivity to the db
func (r *Resolver) RecordUserActivity(
	ctx context.Context,
	action string,
	objectType graphql.ObjectType,
	objectID *string,
	objectCode *string,
) {
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		fmt.Println("sqlboiler hook - get user from context: %w", err)
		return
	}

	_, err = r.UserActivityStore.Insert(
		&db.UserActivity{
			UserID:     userID.String(),
			Action:     action,
			ObjectType: objectType.String(),
			ObjectID:   null.StringFromPtr(objectID),
			ObjectCode: null.StringFromPtr(objectCode),
		},
	)
	if err != nil {
		fmt.Println("update user activity: %w", err)
		return
	}
}

////////////////
//  Resolver  //
////////////////

// UserActivity resolver
func (r *Resolver) UserActivity() graphql.UserActivityResolver {
	return &userActivityResolver{r}
}

type userActivityResolver struct{ *Resolver }

func (r *userActivityResolver) User(ctx context.Context, obj *db.UserActivity) (*db.User, error) {
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
func (r *userActivityResolver) ObjectType(ctx context.Context, obj *db.UserActivity) (graphql.ObjectType, error) {
	return (graphql.ObjectType)(obj.ObjectType), nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserActivities(ctx context.Context, search graphql.SearchFilter, limit int, offset int, userID *string) (*graphql.UserActivityResult, error) {
	total, activities, err := r.UserActivityStore.SearchSelect(
		search,
		limit,
		offset,
		null.StringFromPtr(userID),
	)
	if err != nil {
		return nil, terror.New(err, "list user activity")
	}

	result := &graphql.UserActivityResult{
		UserActivities: activities,
		Total:          int(total),
	}

	return result, nil
}
