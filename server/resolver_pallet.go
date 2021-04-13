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

// Pallet resolver
func (r *Resolver) Pallet() graphql.PalletResolver {
	return &palletResolver{r}
}

type palletResolver struct{ *Resolver }

func (r *palletResolver) Container(ctx context.Context, obj *db.Pallet) (*db.Container, error) {
	if !obj.ContainerID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.ContainerID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := ContainerLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get container")
	}
	return result, nil
}

func (r *palletResolver) CartonCount(ctx context.Context, obj *db.Pallet) (int, error) {
	count, err := r.PalletStore.CartonCount(obj)
	if err != nil {
		return 0, terror.New(err, "")
	}
	return int(count), nil
}

func (r *palletResolver) LatestTrackAction(ctx context.Context, obj *db.Pallet) (*graphql.LatestTransactionInfo, error) {
	action, err := r.PalletStore.LatestTrackAction(obj)
	if err != nil {
		return nil, terror.New(err, "get pallet latest track action")
	}
	return action, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Pallet(ctx context.Context, code string) (*db.Pallet, error) {
	pallet, err := r.PalletStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get pallet")
	}

	return pallet, nil
}

func (r *queryResolver) Pallets(
	ctx context.Context,
	search graphql.SearchFilter,
	limit int,
	offset int,
	containerID *string,
	trackActionID *string,
) (*graphql.PalletResult, error) {
	total, pallets, err := r.PalletStore.SearchSelect(
		search,
		limit,
		offset,
		null.StringFromPtr(containerID),
		null.StringFromPtr(trackActionID),
	)
	if err != nil {
		return nil, terror.New(err, "list pallet")
	}

	result := &graphql.PalletResult{
		Pallets: pallets,
		Total:   int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// PalletCreate creates an pallet and returns the spreadsheet download link
func (r *mutationResolver) PalletCreate(ctx context.Context, input graphql.CreatePallet) (string, error) {
	if input.Quantity <= 0 || input.Quantity > 10000 {
		return "", terror.New(fmt.Errorf("invalid pallet quantity (%d)", input.Quantity), "")
	}

	containerID := null.StringFromPtr(nil)
	if input.ContainerID != nil && len(input.ContainerID.String) > 1 {
		containerID = *input.ContainerID
	}

	// Get Pallet count (for Pallet Code)
	count, err := r.PalletStore.Count()
	if err != nil {
		return "", terror.New(err, "create pallet")
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return "", terror.New(terror.ErrBadContext, "")
	}

	// Create Pallets
	startCode := ""
	endCode := ""

	for i := 0; i < input.Quantity; i++ {
		code := fmt.Sprintf("PAL%05d", count)
		u := &db.Pallet{
			Code:        code,
			CreatedByID: userID.String(),
			ContainerID: containerID,
			Description: input.Description,
		}

		_, err := r.PalletStore.Insert(u)
		if err != nil {
			return "", terror.New(err, "create pallet")
		}

		if i == 0 {
			startCode = code
		} else if i == input.Quantity-1 {
			endCode = code
		}

		count++
	}

	r.RecordUserActivity(ctx, "Created Pallets", graphql.ObjectTypePallet, nil, nil)

	spreadSheetLink := fmt.Sprintf("%ssheet?type=pallet&from=%s&to=%s", r.Config.API.BlobBaseURL, startCode, endCode)
	return spreadSheetLink, nil
}

// PalletUpdate updates a pallet
func (r *mutationResolver) PalletUpdate(ctx context.Context, id string, input graphql.UpdatePallet) (*db.Pallet, error) {
	// Get Pallet
	palletUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.PalletStore.Get(palletUUID)
	if err != nil {
		return nil, terror.New(err, "update pallet")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived pallet")
	}

	// Update Pallet
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.ContainerID != nil {
		if len(input.ContainerID.String) <= 1 {
			if u.ContainerID.Valid {
				// Track Container Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackRemovedFromContainer,
					PalletIDs:       []string{u.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: removed from container")
				}
			}

			u.ContainerID = null.StringFromPtr(nil)
		} else {
			if u.ContainerID.String != input.ContainerID.String {
				// Track Container Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackMovedToContainer,
					PalletIDs:       []string{u.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: moved to container")
				}
			}

			u.ContainerID = *input.ContainerID
		}
	}

	updated, err := r.PalletStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update pallet")
	}

	r.RecordUserActivity(ctx, "Updated Pallet", graphql.ObjectTypePallet, &updated.ID, &updated.Code)

	return updated, nil
}

// PalletArchive archives an pallet
func (r *mutationResolver) PalletArchive(ctx context.Context, id string) (*db.Pallet, error) {
	// Get Pallet
	palletUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Pallet
	updated, err := r.PalletStore.Archive(palletUUID)
	if err != nil {
		return nil, terror.New(err, "update pallet")
	}

	r.RecordUserActivity(ctx, "Archived Pallet", graphql.ObjectTypePallet, &updated.ID, &updated.Code)

	return updated, nil
}

// PalletUnarchive unarchives an pallet
func (r *mutationResolver) PalletUnarchive(ctx context.Context, id string) (*db.Pallet, error) {
	// Get Pallet
	palletUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Pallet
	updated, err := r.PalletStore.Unarchive(palletUUID)
	if err != nil {
		return nil, terror.New(err, "update pallet")
	}

	r.RecordUserActivity(ctx, "Unarchived Pallet", graphql.ObjectTypePallet, &updated.ID, &updated.Code)

	return updated, nil
}

// PalletBatchAction attempts to do an action of each pallet
func (r *mutationResolver) PalletBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypePallet, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.PalletArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive pallet")
			}
		case graphql.ActionUnarchive:
			_, err := r.PalletUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive pallet")
			}

		case graphql.ActionDetachFromContainer:
			newID := null.StringFrom("-")
			input := graphql.UpdatePallet{ContainerID: &newID}
			_, err := r.PalletUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "detach pallet from container")
			}
		case graphql.ActionSetContainer:
			input := graphql.UpdatePallet{ContainerID: value.Str}
			_, err := r.PalletUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "move pallet to container")
			}

		case graphql.ActionSetBonusLoyaltyPoints:
			palletUUID, err := uuid.FromString(id)
			if err != nil {
				return false, terror.New(err, "")
			}

			products, err := r.ProductStore.GetManyByPalletID(palletUUID)
			if err != nil {
				return false, terror.New(err, "")
			}

			for _, product := range products {
				input := graphql.UpdateProduct{
					LoyaltyPoints:       value.No,
					LoyaltyPointsExpire: value.DateTime,
				}
				_, err := r.ProductUpdate(ctx, product.ID, input)
				if err != nil {
					return false, terror.New(err, "change product bonus loyalty points (by pallet)")
				}
			}

		}
	}

	return true, nil
}
