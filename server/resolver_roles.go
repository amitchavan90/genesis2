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

// Role resolver
func (r *Resolver) Role() graphql.RoleResolver {
	return &roleResolver{r}
}

func (r *roleResolver) Permissions(ctx context.Context, obj *db.Role) ([]graphql.Perm, error) {
	result := []graphql.Perm{}
	for _, p := range obj.Permissions {
		result = append(result, (graphql.Perm)(p))
	}
	return result, nil
}

func (r *roleResolver) TrackActions(ctx context.Context, obj *db.Role) ([]*db.TrackAction, error) {
	result, err := r.RoleStore.GetTrackActions(obj)
	if err != nil {
		return nil, terror.New(err, "get role track actions")
	}

	return result, nil
}

type roleResolver struct{ *Resolver }

///////////////
//   Query   //
///////////////

func (r *queryResolver) Role(ctx context.Context, name string) (*db.Role, error) {
	role, err := r.RoleStore.GetByName(name)
	if err != nil {
		return nil, terror.New(err, "get role")
	}

	return role, nil
}

func (r *queryResolver) Roles(ctx context.Context, search graphql.SearchFilter, limit int, offset int, excludeSuper bool) (*graphql.RolesResult, error) {
	// Only let other SUPERADMINS select SUPERADMIN when creating/updating users
	if excludeSuper {
		// Get user
		role, err := r.Auther.RoleFromContext(ctx, false)
		if err != nil {
			return nil, terror.New(terror.ErrBadContext, "")
		}

		if role.Name == "SUPERADMIN" {
			excludeSuper = false
		}
	}

	total, roles, err := r.RoleStore.SearchSelect(search, limit, offset, excludeSuper)
	if err != nil {
		return nil, terror.New(err, "list role")
	}

	result := &graphql.RolesResult{
		Roles: roles,
		Total: int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// RoleCreate creates a role
func (r *mutationResolver) RoleCreate(ctx context.Context, input graphql.UpdateRole) (*db.Role, error) {
	// Create role
	role := &db.Role{}

	if input.Name != nil {
		role.Name = input.Name.String
	}
	if input.Permissions != nil {
		for _, p := range input.Permissions {
			validPerm := false
			for _, vp := range graphql.AllPerm {
				if p == vp.String() {
					validPerm = true
					break
				}
			}
			if !validPerm {
				return nil, terror.New(fmt.Errorf("invalid permission"), "")
			}
		}
		role.Permissions = input.Permissions

		// Prevent making roles equal to SUPERADMIN
		if len(role.Permissions) >= len(graphql.AllPerm) {
			return nil, terror.New(terror.ErrUnauthorized, "")
		}
	}

	// get track actions
	var errs []error
	trackActions := db.TrackActionSlice{}
	if input.TrackActionIDs != nil && len(input.TrackActionIDs) > 0 {
		trackActions, errs = r.TrackActionStore.GetMany(input.TrackActionIDs)
		if errs != nil {
			return nil, terror.New(errs[0], "create role")
		}
	}

	// create
	created, err := r.RoleStore.Insert(role)
	if err != nil {
		return nil, terror.New(err, "create role")
	}

	// set track actions
	if len(trackActions) > 0 {
		err = r.RoleStore.SetTrackActions(role, trackActions)
		if err != nil {
			return nil, terror.New(err, "set role track actions")
		}
	}

	r.RecordUserActivity(ctx, "Created Role", graphql.ObjectTypeRole, &created.ID, &created.Name)

	return created, nil
}

func (r *mutationResolver) RoleUpdate(ctx context.Context, id string, input graphql.UpdateRole) (*db.Role, error) {
	// Get role
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	role, err := r.RoleStore.Get(roleUUID)
	if err != nil {
		return nil, terror.New(err, "update role")
	}

	// Check archived state
	if role.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived role")
	}

	// get track actions
	var errs []error
	trackActions := db.TrackActionSlice{}
	if input.TrackActionIDs != nil {
		trackActions, errs = r.TrackActionStore.GetMany(input.TrackActionIDs)
		if errs != nil {
			return nil, terror.New(errs[0], "update role")
		}
	}

	// update role
	if input.Name != nil {
		role.Name = input.Name.String
	}
	if input.Permissions != nil {
		for _, p := range input.Permissions {
			validPerm := false
			for _, vp := range graphql.AllPerm {
				if p == vp.String() {
					validPerm = true
					break
				}
			}
			if !validPerm {
				return nil, terror.New(err, "invalid permission")
			}
		}
		role.Permissions = input.Permissions

		// Prevent making roles equal to SUPERADMIN
		if len(role.Permissions) >= len(graphql.AllPerm) {
			return nil, terror.New(terror.ErrUnauthorized, "")
		}
	}

	updated, err := r.RoleStore.Update(role)
	if err != nil {
		return nil, terror.New(err, "update role")
	}

	// set track actions
	err = r.RoleStore.SetTrackActions(role, trackActions)
	if err != nil {
		return nil, terror.New(err, "update role track actions")
	}

	r.RecordUserActivity(ctx, "Updated Role", graphql.ObjectTypeRole, &updated.ID, &updated.Name)

	return updated, nil
}

func (r *mutationResolver) RoleArchive(ctx context.Context, id string) (*db.Role, error) {
	// Get Role
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Role
	updated, err := r.RoleStore.Archive(roleUUID)
	if err != nil {
		return nil, terror.New(err, "update role")
	}

	r.RecordUserActivity(ctx, "Archived Role", graphql.ObjectTypeRole, &updated.ID, &updated.Name)

	return updated, nil
}
func (r *mutationResolver) RoleUnarchive(ctx context.Context, id string) (*db.Role, error) {
	// Get Role
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Role
	updated, err := r.RoleStore.Unarchive(roleUUID)
	if err != nil {
		return nil, terror.New(err, "update role")
	}

	r.RecordUserActivity(ctx, "Unarchived Role", graphql.ObjectTypeRole, &updated.ID, &updated.Name)

	return updated, nil
}
