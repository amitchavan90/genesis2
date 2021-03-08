package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"time"

	"github.com/ninja-software/terror"

	"github.com/volatiletech/null"

	"github.com/gofrs/uuid"
)

////////////////
//  Resolver  //
////////////////

// Carton resolver
func (r *Resolver) Carton() graphql.CartonResolver {
	return &cartonResolver{r}
}

type cartonResolver struct{ *Resolver }

func (r *cartonResolver) Pallet(ctx context.Context, obj *db.Carton) (*db.Pallet, error) {
	if !obj.PalletID.Valid {
		return nil, nil
	}

	id, err := uuid.FromString(obj.PalletID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := PalletLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get pallet")
	}
	return result, nil
}

func (r *cartonResolver) ProductCount(ctx context.Context, obj *db.Carton) (int, error) {
	count, err := r.CartonStore.ProductCount(obj)
	if err != nil {
		return 0, terror.New(err, "get carton product count")
	}
	return int(count), nil
}

func (r *cartonResolver) LatestTrackAction(ctx context.Context, obj *db.Carton) (*graphql.LatestTransactionInfo, error) {
	action, err := r.CartonStore.LatestTrackAction(obj)
	if err != nil {
		return nil, terror.New(err, "get carton latest track action")
	}
	return action, nil
}

func (r *cartonResolver) Order(ctx context.Context, obj *db.Carton) (*db.Order, error) {
	orderID, err := r.CartonStore.OrderID(obj)
	if err != nil {
		return nil, terror.New(err, "get order")
	}
	if orderID == nil {
		return nil, nil
	}
	orderUUID, err := uuid.FromString(*orderID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := OrderLoaderFromContext(ctx, orderUUID)
	if err != nil {
		return nil, terror.New(err, "get order")
	}
	return result, nil
}

func (r *cartonResolver) Distributor(ctx context.Context, obj *db.Carton) (*db.Distributor, error) {
	distributorID, err := r.CartonStore.DistributorID(obj)
	if err != nil {
		return nil, terror.New(err, "get distributor")
	}
	if distributorID == nil {
		return nil, nil
	}
	distributorUUID, err := uuid.FromString(*distributorID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := DistributorLoaderFromContext(ctx, distributorUUID)
	if err != nil {
		return nil, terror.New(err, "get distributor")
	}
	return result, nil
}

func (r *cartonResolver) Sku(ctx context.Context, obj *db.Carton) (*db.StockKeepingUnit, error) {
	skuID, err := r.CartonStore.SkuID(obj)
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

func (r *cartonResolver) Transactions(ctx context.Context, obj *db.Carton) ([]*db.Transaction, error) {
	// Get Transactions
	cartonUUID, err := uuid.FromString(obj.ID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	transactions, err := r.TransactionStore.GetByCartonID(cartonUUID)
	if err != nil {
		return nil, terror.New(err, "get carton transactions")
	}

	// Omit private transactions if not logged in
	userID, _ := r.Auther.UserIDFromContext(ctx)

	// skip repeated track action (same carton but diff products)
	foundTrackActionIDs := make(map[string]bool)

	result := db.TransactionSlice{}
	for _, t := range transactions {
		if userID == uuid.Nil && t.R.TrackAction.Private {
			continue
		}

		// prevent products being updated multiple times in one action
		_, found := foundTrackActionIDs[t.TrackActionID]
		if found {
			continue
		}
		foundTrackActionIDs[t.TrackActionID] = true

		// public only see little detail
		if userID == uuid.Nil {
			result = append(
				result,
				&db.Transaction{
					ID:            t.ID,
					TrackActionID: t.TrackActionID,
					ScannedAt:     t.ScannedAt,
					CartonID:      t.CartonID,
				},
			)
			continue
		}

		result = append(result, t)
	}

	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Carton(ctx context.Context, code string) (*db.Carton, error) {
	carton, err := r.CartonStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get carton")
	}

	return carton, nil
}

func (r *queryResolver) Cartons(
	ctx context.Context,
	search graphql.SearchFilter,
	limit int,
	offset int,
	containerID *string,
	trackActionID *string,
) (*graphql.CartonResult, error) {
	total, cartons, err := r.CartonStore.SearchSelect(
		search,
		limit,
		offset,
		null.StringFromPtr(containerID),
		null.StringFromPtr(trackActionID),
	)
	if err != nil {
		return nil, terror.New(err, "list carton")
	}

	result := &graphql.CartonResult{
		Cartons: cartons,
		Total:   int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// CartonCreate creates an carton and returns the spreadsheet download link
func (r *mutationResolver) CartonCreate(ctx context.Context, input graphql.CreateCarton) (string, error) {
	if input.Quantity <= 0 || input.Quantity > 10000 {
		return "", terror.New(fmt.Errorf("invalid carton quantity (%d)", input.Quantity), "")
	}

	palletID := null.StringFromPtr(nil)
	if input.PalletID != nil && len(input.PalletID.String) > 1 {
		palletID = *input.PalletID
	}

	// Get Carton count (for Carton Code)
	count, err := r.CartonStore.Count()
	if err != nil {
		return "", terror.New(err, "create carton")
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return "", terror.New(terror.ErrBadContext, "")
	}

	// Create Cartons
	ssLink := helpers.GetCartonSpreadsheetLink(fmt.Sprintf("CAR%05d", count), input.Quantity, int(count))

	for i := 0; i < input.Quantity; i++ {
		code := fmt.Sprintf("CAR%05d", count)
		u := &db.Carton{
			Code:            code,
			CreatedByID:     userID.String(),
			PalletID:        palletID,
			Description:     input.Description,
			SpreadsheetLink: ssLink,
		}

		_, err := r.CartonStore.Insert(u)
		if err != nil {
			return "", terror.New(err, "create carton")
		}

		count++
	}

	r.RecordUserActivity(ctx, "Created Cartons", graphql.ObjectTypeContainer, nil, nil)

	return ssLink, nil
}

// CartonUpdate updates a carton
func (r *mutationResolver) CartonUpdate(ctx context.Context, id string, input graphql.UpdateCarton) (*db.Carton, error) {
	// Get Carton
	cartonUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.CartonStore.Get(cartonUUID)
	if err != nil {
		return nil, terror.New(err, "update carton")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived carton")
	}

	// Update Carton
	if input.Weight != nil {
		u.Weight = input.Weight.String
	}
	if input.ProcessedAt != nil && !input.ProcessedAt.Time.After(time.Now()) {
		u.ProcessedAt = *input.ProcessedAt
	}
	if input.Description != nil {
		u.Description = input.Description.String
	}
	if input.MeatType != nil {
		u.MeatType = input.MeatType.String
	}
	if input.PalletID != nil {
		if len(input.PalletID.String) <= 1 {
			if u.PalletID.Valid {
				// Track Pallet Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackRemovedFromPallet,
					CartonIDs:       []string{u.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: removed from pallet")
				}
			}

			u.PalletID = null.StringFromPtr(nil)
		} else {
			if u.PalletID.String != input.PalletID.String {
				// Track Pallet Move
				_, err = r.RecordTransaction(ctx, graphql.RecordTransactionInput{
					TrackActionCode: trackMovedToPallet,
					CartonIDs:       []string{u.ID},
				})
				if err != nil {
					fmt.Printf("Error on track action: moved to pallet")
				}
			}

			u.PalletID = *input.PalletID
		}
	}

	updated, err := r.CartonStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update carton")
	}

	r.RecordUserActivity(ctx, "Updated Carton", graphql.ObjectTypeCarton, &updated.ID, &updated.Code)

	return updated, nil
}

// CartonArchive archives an carton
func (r *mutationResolver) CartonArchive(ctx context.Context, id string) (*db.Carton, error) {
	// Get Carton
	cartonUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Carton
	updated, err := r.CartonStore.Archive(cartonUUID)
	if err != nil {
		return nil, terror.New(err, "update carton")
	}

	r.RecordUserActivity(ctx, "Archived Carton", graphql.ObjectTypeCarton, &updated.ID, &updated.Code)

	return updated, nil
}

// CartonUnarchive unarchives an carton
func (r *mutationResolver) CartonUnarchive(ctx context.Context, id string) (*db.Carton, error) {
	// Get Carton
	cartonUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Carton
	updated, err := r.CartonStore.Unarchive(cartonUUID)
	if err != nil {
		return nil, terror.New(err, "update carton")
	}

	r.RecordUserActivity(ctx, "Unarchived Carton", graphql.ObjectTypeCarton, &updated.ID, &updated.Code)

	return updated, nil
}

// CartonBatchAction attempts to do an action of each carton
func (r *mutationResolver) CartonBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypeCarton, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.CartonArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive carton")
			}
		case graphql.ActionUnarchive:
			_, err := r.CartonUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive carton")
			}

		case graphql.ActionDetachFromPallet:
			newID := null.StringFrom("-")
			input := graphql.UpdateCarton{PalletID: &newID}
			_, err := r.CartonUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "detach carton from pallet")
			}
		case graphql.ActionSetPallet:
			input := graphql.UpdateCarton{PalletID: value.Str}
			_, err := r.CartonUpdate(ctx, id, input)
			if err != nil {
				return false, terror.New(err, "move carton to pallet")
			}

		case graphql.ActionDetachFromDistributor:
			newID := null.StringFrom("-")
			err := r.cartonSetDistributor(ctx, id, &newID)
			if err != nil {
				return false, terror.New(err, "")
			}
		case graphql.ActionSetDistributor:
			err := r.cartonSetDistributor(ctx, id, value.Str)
			if err != nil {
				return false, terror.New(err, "")
			}

		case graphql.ActionDetachFromContract:
			newID := null.StringFrom("-")
			err := r.cartonSetContract(ctx, id, &newID)
			if err != nil {
				return false, terror.New(err, "")
			}
		case graphql.ActionSetContract:
			err := r.cartonSetContract(ctx, id, value.Str)
			if err != nil {
				return false, terror.New(err, "")
			}

		case graphql.ActionSetBonusLoyaltyPoints:
			cartonUUID, err := uuid.FromString(id)
			if err != nil {
				return false, terror.New(err, "")
			}

			products, err := r.ProductStore.GetManyByCartonID(cartonUUID)
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
					return false, terror.New(err, "change product bonus loyalty points (by carton)")
				}
			}

		}
	}

	return true, nil
}

func (r *mutationResolver) cartonSetDistributor(ctx context.Context, cartonID string, distributorID *null.String) error {
	input := graphql.UpdateProduct{DistributorID: distributorID}

	cartonUUID, err := uuid.FromString(cartonID)
	if err != nil {
		return terror.New(err, "invalid carton uuid")
	}

	products, err := r.ProductStore.GetManyByCartonID(cartonUUID)
	if err != nil {
		return terror.New(err, "failed to get products by carton id")
	}

	for _, product := range products {
		_, err := r.ProductUpdate(ctx, product.ID, input)
		if err != nil {
			return terror.New(err, "set product distributor (by carton)")
		}
	}

	return nil
}

func (r *mutationResolver) cartonSetContract(ctx context.Context, cartonID string, contractID *null.String) error {
	input := graphql.UpdateProduct{ContractID: contractID}

	cartonUUID, err := uuid.FromString(cartonID)
	if err != nil {
		return terror.New(err, "invalid carton uuid")
	}

	products, err := r.ProductStore.GetManyByCartonID(cartonUUID)
	if err != nil {
		return terror.New(err, "failed to get products by carton id")
	}

	for _, product := range products {
		_, err := r.ProductUpdate(ctx, product.ID, input)
		if err != nil {
			return terror.New(err, "set product contract (by carton)")
		}
	}

	return nil
}
