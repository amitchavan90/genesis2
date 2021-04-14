package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
)

///////////////
//   Query   //
///////////////

func (r *queryResolver) Distributor(ctx context.Context, code string) (*db.Distributor, error) {
	distributor, err := r.DistributorStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get distributor")
	}

	return distributor, nil
}

func (r *queryResolver) Distributors(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.DistributorResult, error) {
	total, distributors, err := r.DistributorStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list distributor")
	}

	result := &graphql.DistributorResult{
		Distributors: distributors,
		Total:        int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// DistributorCreate creates an distributor and returns the spreadsheet download link
func (r *mutationResolver) DistributorCreate(ctx context.Context, input graphql.UpdateDistributor) (*db.Distributor, error) {
	if input.Name == nil || input.Code == nil {
		return nil, terror.ErrInvalidInput
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	// Create Distributor
	u := &db.Distributor{
		Name:        input.Name.String,
		Code:        input.Code.String,
		CreatedByID: userID.String(),
	}

	created, err := r.DistributorStore.Insert(u)
	if err != nil {
		return nil, terror.New(err, "create distributor")
	}

	r.RecordUserActivity(ctx, "Created Distributor", graphql.ObjectTypeDistributor, &created.ID, &created.Code)

	return created, nil
}

// DistributorUpdate updates a distributor
func (r *mutationResolver) DistributorUpdate(ctx context.Context, id string, input graphql.UpdateDistributor) (*db.Distributor, error) {
	// Get Distributor
	distributorUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.DistributorStore.Get(distributorUUID)
	if err != nil {
		return nil, terror.New(err, "update distributor")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived distributor")
	}

	// Update Distributor
	if input.Name != nil {
		u.Name = input.Name.String
	}

	updated, err := r.DistributorStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update distributor")
	}

	r.RecordUserActivity(ctx, "Updated Distributor", graphql.ObjectTypeDistributor, &updated.ID, &updated.Code)

	return updated, nil
}

// DistributorArchive archives an distributor
func (r *mutationResolver) DistributorArchive(ctx context.Context, id string) (*db.Distributor, error) {
	// Get Distributor
	distributorUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Distributor
	updated, err := r.DistributorStore.Archive(distributorUUID)
	if err != nil {
		return nil, terror.New(err, "update distributor")
	}

	r.RecordUserActivity(ctx, "Archived Distributor", graphql.ObjectTypeDistributor, &updated.ID, &updated.Code)

	return updated, nil
}

// DistributorUnarchive unarchives an distributor
func (r *mutationResolver) DistributorUnarchive(ctx context.Context, id string) (*db.Distributor, error) {
	// Get Distributor
	distributorUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Distributor
	updated, err := r.DistributorStore.Unarchive(distributorUUID)
	if err != nil {
		return nil, terror.New(err, "update distributor")
	}

	r.RecordUserActivity(ctx, "Unarchived Distributor", graphql.ObjectTypeDistributor, &updated.ID, &updated.Code)

	return updated, nil
}
