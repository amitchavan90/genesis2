package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
)

////////////////
//  Resolver  //
////////////////

// Container resolver
func (r *Resolver) Container() graphql.ContainerResolver {
	return &containerResolver{r}
}

type containerResolver struct{ *Resolver }

func (r *containerResolver) PalletCount(ctx context.Context, obj *db.Container) (int, error) {
	count, err := r.ContainerStore.PalletCount(obj)
	if err != nil {
		return 0, terror.New(err, "get container pallet count")
	}
	return int(count), nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Container(ctx context.Context, code string) (*db.Container, error) {
	container, err := r.ContainerStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get container")
	}

	return container, nil
}

func (r *queryResolver) Containers(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.ContainerResult, error) {
	total, containers, err := r.ContainerStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list container")
	}

	result := &graphql.ContainerResult{
		Containers: containers,
		Total:      int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// ContainerCreate creates an container and returns the spreadsheet download link
func (r *mutationResolver) ContainerCreate(ctx context.Context, input graphql.CreateContainer) (string, error) {
	if input.Quantity <= 0 || input.Quantity > 10000 {
		return "", terror.New(fmt.Errorf("invalid container quantity (%d)", input.Quantity), "")
	}

	// Get Container count (for Container Code)
	count, err := r.ContainerStore.Count()
	if err != nil {
		return "", terror.New(err, "create container")
	}

	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return "", terror.New(terror.ErrBadContext, "")
	}

	// Create Containers
	startCode := ""
	endCode := ""

	for i := 0; i < input.Quantity; i++ {
		code := fmt.Sprintf("CON%05d", count)
		u := &db.Container{
			Description: input.Description,
			Code:        code,
			CreatedByID: userID.String(),
		}

		_, err := r.ContainerStore.Insert(u)
		if err != nil {
			return "", terror.New(err, "create container")
		}

		if i == 0 {
			startCode = code
		} else if i == input.Quantity-1 {
			endCode = code
		}

		count++
	}

	r.RecordUserActivity(ctx, "Created Containers", graphql.ObjectTypeContainer, nil, nil)

	spreadSheetLink := fmt.Sprintf("%ssheet?type=container&from=%s&to=%s", r.Config.API.BlobBaseURL, startCode, endCode)
	return spreadSheetLink, nil
}

// ContainerUpdate updates a container
func (r *mutationResolver) ContainerUpdate(ctx context.Context, id string, input graphql.UpdateContainer) (*db.Container, error) {
	// Get Container
	containerUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.ContainerStore.Get(containerUUID)
	if err != nil {
		return nil, terror.New(err, "update container")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived container")
	}

	// Update Container
	if input.Description != nil {
		u.Description = input.Description.String
	}

	updated, err := r.ContainerStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update container")
	}

	r.RecordUserActivity(ctx, "Updated Container", graphql.ObjectTypeContainer, &updated.ID, &updated.Code)

	return updated, nil
}

// ContainerArchive archives an container
func (r *mutationResolver) ContainerArchive(ctx context.Context, id string) (*db.Container, error) {
	// Get Container
	containerUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Container
	updated, err := r.ContainerStore.Archive(containerUUID)
	if err != nil {
		return nil, terror.New(err, "update container")
	}

	r.RecordUserActivity(ctx, "Archived Container", graphql.ObjectTypeContainer, &updated.ID, &updated.Code)

	return updated, nil
}

// ContainerUnarchive unarchives an container
func (r *mutationResolver) ContainerUnarchive(ctx context.Context, id string) (*db.Container, error) {
	// Get Container
	containerUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Container
	updated, err := r.ContainerStore.Unarchive(containerUUID)
	if err != nil {
		return nil, terror.New(err, "update container")
	}

	r.RecordUserActivity(ctx, "Unarchived Container", graphql.ObjectTypeContainer, &updated.ID, &updated.Code)

	return updated, nil
}

// ContainerBatchAction attempts to do an action of each container
func (r *mutationResolver) ContainerBatchAction(ctx context.Context, ids []string, action graphql.Action, value *graphql.BatchActionInput) (bool, error) {

	r.RecordUserActivity(ctx, "Batch Action: "+action.String(), graphql.ObjectTypeContainer, nil, nil)

	for _, id := range ids {
		switch action {
		case graphql.ActionArchive:
			_, err := r.ContainerArchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "archive container")
			}
		case graphql.ActionUnarchive:
			_, err := r.ContainerUnarchive(ctx, id)
			if err != nil {
				return false, terror.New(err, "unarchive container")
			}

		case graphql.ActionSetBonusLoyaltyPoints:
			containerUUID, err := uuid.FromString(id)
			if err != nil {
				return false, terror.New(err, "")
			}

			products, err := r.ProductStore.GetManyByContainerID(containerUUID)
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
					return false, terror.New(err, "change product bonus loyalty points (by container)")
				}
			}

		}
	}

	return true, nil
}
