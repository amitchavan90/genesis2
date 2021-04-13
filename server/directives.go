package genesis

import (
	"context"
	"errors"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	gql "github.com/99designs/gqlgen/graphql"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/e"
)

// NewDirectiveRoot handles the directives
func NewDirectiveRoot(roleStore RoleStorer) graphql.DirectiveRoot {
	return graphql.DirectiveRoot{
		HasPerm: func(ctx context.Context, obj interface{}, next gql.Resolver, p graphql.Perm) (res interface{}, err error) {
			roleID, err := ClaimValueFromContext(ctx, ClaimRoleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			roleUUID, err := uuid.FromString(roleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			role, err := roleStore.Get(roleUUID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}

			for _, userPerm := range role.Permissions {
				if userPerm == p.String() {
					return next(ctx)
				}
			}

			return nil, terror.New(terror.ErrUnauthorized, "")
		},
		HasAnyPerm: func(ctx context.Context, obj interface{}, next gql.Resolver, p []graphql.Perm) (res interface{}, err error) {
			roleID, err := ClaimValueFromContext(ctx, ClaimRoleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			roleUUID, err := uuid.FromString(roleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			role, err := roleStore.Get(roleUUID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}

			for _, userPerm := range role.Permissions {
				for _, perm := range p {
					if perm.String() == userPerm {
						return next(ctx)
					}
				}
			}

			return nil, terror.New(terror.ErrUnauthorized, "")
		},
		HasAllPerms: func(ctx context.Context, obj interface{}, next gql.Resolver, p []graphql.Perm) (res interface{}, err error) {
			roleID, err := ClaimValueFromContext(ctx, ClaimRoleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			roleUUID, err := uuid.FromString(roleID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}
			role, err := roleStore.Get(roleUUID)
			if err != nil {
				return nil, errors.New(e.ErrorMessage(err))
			}

			if len(role.Permissions) == 0 {
				return nil, terror.New(terror.ErrUnauthorized, "")
			}

			for _, perm := range p {
				hasPerm := false
				for _, userPerm := range role.Permissions {
					if userPerm == perm.String() {
						hasPerm = true
						break
					}
				}
				if !hasPerm {
					return nil, terror.New(terror.ErrUnauthorized, "")
				}
			}

			return next(ctx)
		},
	}
}
