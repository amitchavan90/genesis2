package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"io/ioutil"

	"github.com/ninja-software/terror"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/h2non/filetype"
)

////////////////
//  Resolver  //
////////////////

// Blob resolver
func (r *Resolver) Blob() graphql.BlobResolver {
	return &blobResolver{r}
}

type blobResolver struct{ *Resolver }

func (r *blobResolver) FileURL(ctx context.Context, obj *db.Blob) (*string, error) {
	result := fmt.Sprintf("%s?id=%s", r.Config.API.BlobBaseURL, obj.ID)
	return &result, nil
}

///////////////
// Mutations //
///////////////

// UploadBlob adds a new blob to the db
func (r *mutationResolver) UploadBlob(ctx context.Context, userID string, file gql.Upload) (*db.Blob, error) {
	// Read file data
	data, err := ioutil.ReadAll(file.File)
	if err != nil {
		return nil, terror.New(err, "upload blob - read file")
	}

	// get mime type
	kind, err := filetype.Match(data)
	if err != nil {
		return nil, terror.New(err, "upload blob - get mime type")
	}

	if kind == filetype.Unknown {
		return nil, terror.New(fmt.Errorf("Image type is unknown"), "")
	}

	mimeType := kind.MIME.Value
	extension := kind.Extension

	a := &db.Blob{
		FileName:      file.Filename,
		MimeType:      mimeType,
		Extension:     extension,
		FileSizeBytes: file.Size,
		File:          data,
	}

	b, err := r.BlobStore.Insert(a)
	if err != nil {
		return nil, terror.New(err, "upload blob - insert")
	}

	fileURL := fmt.Sprintf("%s?id=%s", r.Config.API.BlobBaseURL, a.ID)
	r.RecordUserActivity(ctx, "Uploaded File", graphql.ObjectTypeBlob, &a.ID, &fileURL)

	return b, nil
}

func (r *mutationResolver) FileUpload(ctx context.Context, file gql.Upload) (*db.Blob, error) {
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	blob, err := r.UploadBlob(ctx, userID.String(), file)
	if err != nil {
		return nil, terror.New(err, "file upload")
	}

	return blob, nil
}
func (r *mutationResolver) FileUploadMultiple(ctx context.Context, files []*gql.Upload) ([]*db.Blob, error) {
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrBadContext, "")
	}

	blobs := []*db.Blob{}

	for _, file := range files {
		blob, err := r.UploadBlob(ctx, userID.String(), *file)
		if err != nil {
			return nil, terror.New(err, "file upload (multiple)")
		}

		blobs = append(blobs, blob)
	}

	return blobs, nil
}
