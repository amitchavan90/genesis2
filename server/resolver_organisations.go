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

// Organisation resolver
func (r *Resolver) Organisation() graphql.OrganisationResolver {
	return &organisationResolver{r}
}

type organisationResolver struct{ *Resolver }

func (r *organisationResolver) Users(ctx context.Context, obj *db.Organisation) ([]*db.User, error) {
	orgID, err := uuid.FromString(obj.ID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := OrganisationUsersLoaderFromContext(ctx, orgID)
	if err != nil {
		return nil, terror.New(err, "list users from org")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Organisations(ctx context.Context) ([]*db.Organisation, error) {
	result, err := r.OrganisationStore.All()
	if err != nil {
		return nil, terror.New(err, "list organisation")
	}
	return result, nil
}
