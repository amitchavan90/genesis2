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

////////////////
//  Resolver  //
////////////////

// TrackAction resolver
func (r *Resolver) TrackAction() graphql.TrackActionResolver {
	return &trackActionResolver{r}
}

type trackActionResolver struct{ *Resolver }

func (r *trackActionResolver) RequirePhotos(ctx context.Context, obj *db.TrackAction) ([]bool, error) {
	return obj.RequirePhotos, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) TrackAction(ctx context.Context, code string) (*db.TrackAction, error) {
	trackAction, err := r.TrackActionStore.GetByCode(code)
	if err != nil {
		return nil, terror.New(err, "get track action")
	}
	return trackAction, nil
}

func (r *queryResolver) TrackActions(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.TrackActionResult, error) {
	total, trackActions, err := r.TrackActionStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list track actions")
	}

	result := &graphql.TrackActionResult{
		TrackActions: trackActions,
		Total:        int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) TrackActionCreate(ctx context.Context, input graphql.UpdateTrackAction) (*db.TrackAction, error) {
	// Get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	// Get TrackAction count
	count, err := r.TrackActionStore.Count()
	if err != nil {
		return nil, terror.New(err, "create track action")
	}

	// Check name
	if input.Name == nil || !input.Name.Valid {
		return nil, terror.New(err, "track action requires a name")
	}

	nameChinese := ""
	if input.NameChinese != nil {
		nameChinese = input.NameChinese.String
	}

	private := false
	if input.Private != nil {
		private = input.Private.Bool
	}

	blockchain := true
	if input.Blockchain != nil {
		blockchain = input.Blockchain.Bool
	}

	requirePhotos := []bool{true, true}
	if input.RequirePhotos != nil && len(input.RequirePhotos) == 2 {
		requirePhotos = input.RequirePhotos
	}

	// Create TrackAction
	trackAction := &db.TrackAction{
		Code:          fmt.Sprintf("TRACK%03d", count),
		Name:          input.Name.String,
		NameChinese:   nameChinese,
		Private:       private,
		Blockchain:    blockchain,
		CreatedByID:   null.StringFrom(userID.String()),
		RequirePhotos: requirePhotos,
	}

	created, err := r.TrackActionStore.Insert(trackAction)
	if err != nil {
		return nil, terror.New(err, "create track action")
	}

	r.RecordUserActivity(ctx, "Created Track Action", graphql.ObjectTypeTrackAction, &created.ID, &created.Code)

	return created, nil
}

func (r *mutationResolver) TrackActionUpdate(ctx context.Context, id string, input graphql.UpdateTrackAction) (*db.TrackAction, error) {
	// Get TrackAction
	trackActionUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	trackAction, err := r.TrackActionStore.Get(trackActionUUID)
	if err != nil {
		return nil, terror.New(err, "update track action")
	}

	// Check archived state
	if trackAction.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived track action")
	}

	if input.Name != nil {
		trackAction.Name = input.Name.String
	}
	if input.NameChinese != nil {
		trackAction.NameChinese = input.NameChinese.String
	}
	if input.Private != nil {
		trackAction.Private = input.Private.Bool
	}
	if input.Blockchain != nil {
		trackAction.Blockchain = input.Blockchain.Bool
	}
	if input.RequirePhotos != nil && len(input.RequirePhotos) == 2 {
		trackAction.RequirePhotos = input.RequirePhotos
	}

	updated, err := r.TrackActionStore.Update(trackAction)
	if err != nil {
		return nil, terror.New(err, "update track action")
	}

	r.RecordUserActivity(ctx, "Updated Track Action", graphql.ObjectTypeTrackAction, &updated.ID, &updated.Code)

	return updated, nil
}

func (r *mutationResolver) TrackActionArchive(ctx context.Context, id string) (*db.TrackAction, error) {
	// Get TrackAction
	trackActionUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive TrackAction
	updated, err := r.TrackActionStore.Archive(trackActionUUID)
	if err != nil {
		return nil, terror.New(err, "archive track action")
	}

	r.RecordUserActivity(ctx, "Archived Track Action", graphql.ObjectTypeTrackAction, &updated.ID, &updated.Code)

	return updated, nil
}

func (r *mutationResolver) TrackActionUnarchive(ctx context.Context, id string) (*db.TrackAction, error) {
	// Get TrackAction
	trackActionUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive TrackAction
	updated, err := r.TrackActionStore.Unarchive(trackActionUUID)
	if err != nil {
		return nil, terror.New(err, "unarchive track action")
	}

	r.RecordUserActivity(ctx, "Unarchived Track Action", graphql.ObjectTypeTrackAction, &updated.ID, &updated.Code)

	return updated, nil
}
