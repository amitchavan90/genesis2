package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
)

///////////////
//   Query   //
///////////////

func (r *queryResolver) Contract(ctx context.Context, code string) (*db.Contract, error) {
	contract, err := r.ContractStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get contract")
	}

	return contract, nil
}

func (r *queryResolver) Contracts(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.ContractResult, error) {
	total, contracts, err := r.ContractStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list contract")
	}

	result := &graphql.ContractResult{
		Contracts: contracts,
		Total:     int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// ContractCreate creates an contract and returns the spreadsheet download link
func (r *mutationResolver) ContractCreate(ctx context.Context, input graphql.UpdateContract) (*db.Contract, error) {
	// Get Contract count (for Contract Code)
	count, err := r.ContractStore.Count()
	if err != nil {
		return nil, terror.New(err, "create contract")
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	// Create Contract
	code := fmt.Sprintf("CONTRACT%05d", count)
	u := &db.Contract{
		Code:        code,
		CreatedByID: userID.String(),
	}

	if input.Name != nil {
		u.Name = input.Name.String
	}
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.SupplierName != nil {
		u.SupplierName = input.SupplierName.String
	}
	if input.DateSigned != nil {
		u.DateSigned = *input.DateSigned
	}

	u.Latitude = input.Latitude
	u.Longitude = input.Longitude

	created, err := r.ContractStore.Insert(u)
	if err != nil {
		return nil, terror.New(err, "create contract")
	}

	r.RecordUserActivity(ctx, "Created Contract", graphql.ObjectTypeContract, &created.ID, &created.Code)

	return created, nil
}

// ContractUpdate updates a contract
func (r *mutationResolver) ContractUpdate(ctx context.Context, id string, input graphql.UpdateContract) (*db.Contract, error) {
	// Get Contract
	contractUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.ContractStore.Get(contractUUID)
	if err != nil {
		return nil, terror.New(err, "update contract")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived contract")
	}

	// Update Contract
	if input.Name != nil {
		u.Name = input.Name.String
	}
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.SupplierName != nil {
		u.SupplierName = input.SupplierName.String
	}
	if input.DateSigned != nil {
		u.DateSigned = *input.DateSigned
	}

	updated, err := r.ContractStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update contract")
	}

	r.RecordUserActivity(ctx, "Updated Contract", graphql.ObjectTypeContract, &updated.ID, &updated.Code)

	return updated, nil
}

// ContractArchive archives an contract
func (r *mutationResolver) ContractArchive(ctx context.Context, id string) (*db.Contract, error) {
	// Get Contract
	contractUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Contract
	updated, err := r.ContractStore.Archive(contractUUID)
	if err != nil {
		return nil, terror.New(err, "update contract")
	}

	r.RecordUserActivity(ctx, "Archived Contract", graphql.ObjectTypeContract, &updated.ID, &updated.Code)

	return updated, nil
}

// ContractUnarchive unarchives an contract
func (r *mutationResolver) ContractUnarchive(ctx context.Context, id string) (*db.Contract, error) {
	// Get Contract
	contractUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Contract
	updated, err := r.ContractStore.Unarchive(contractUUID)
	if err != nil {
		return nil, terror.New(err, "update contract")
	}

	r.RecordUserActivity(ctx, "Unarchived Contract", graphql.ObjectTypeContract, &updated.ID, &updated.Code)

	return updated, nil
}
