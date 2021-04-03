package genesis

import (
	"context"
	"fmt"
	"genesis/crypto"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"strings"

	"github.com/ninja-software/terror"

	"github.com/volatiletech/null"

	"github.com/gofrs/uuid"
)

////////////////
//  Resolver  //
////////////////

// User resolver
func (r *Resolver) User() graphql.UserResolver {
	return &userResolver{r}
}

type userResolver struct{ *Resolver }

func (r *userResolver) Role(ctx context.Context, obj *db.User) (*db.Role, error) {
	id, err := uuid.FromString(obj.RoleID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := RoleLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get role")
	}
	return result, nil
}

func (r *userResolver) Organisation(ctx context.Context, obj *db.User) (*db.Organisation, error) {
	if !obj.OrganisationID.Valid {
		return nil, nil
	}
	id, err := uuid.FromString(obj.OrganisationID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	result, err := OrganisationLoaderFromContext(ctx, id)
	if err != nil {
		return nil, terror.New(err, "get organisation")
	}
	return result, nil
}

func (r *userResolver) Referrals(ctx context.Context, obj *db.User) ([]*db.Referral, error) {
	result, err := r.UserStore.GetReferrals(obj.ID)
	if err != nil {
		return nil, terror.New(err, "get referrals")
	}
	return result, nil
}

func (r *userResolver) WalletHistory(ctx context.Context, obj *db.User) ([]*db.WalletHistory, error) {
	result, err := r.UserStore.GetWalletHistory(obj.ID)
	if err != nil {
		return nil, terror.New(err, "get referrals")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) User(ctx context.Context, email *string, wechatID *string) (*db.User, error) {
	if email != nil {
		user, err := r.UserStore.GetByEmail(*email)
		if err != nil {
			return nil, terror.New(err, "get user")
		}
		return user, nil
	}
	if wechatID != nil {
		user, err := r.UserStore.GetByWechatID(*wechatID)
		if err != nil {
			return nil, terror.New(err, "get user")
		}
		return user, nil
	}

	return nil, terror.New(fmt.Errorf("get user: no argument provided"), "")
}

func (r *queryResolver) Users(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.UsersResult, error) {
	total, users, err := r.UserStore.SearchSelect(search, limit, offset, false)
	if err != nil {
		return nil, terror.New(err, "list user")
	}

	result := &graphql.UsersResult{
		Users: users,
		Total: int(total),
	}

	return result, nil
}

func (r *queryResolver) Consumers(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.ConsumersResult, error) {
	total, users, err := r.UserStore.SearchSelect(search, limit, offset, true)
	if err != nil {
		return nil, terror.New(err, "list consumers")
	}

	result := &graphql.ConsumersResult{
		Consumers: users,
		Total:     int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// ChangeDetails of a user
func (r *mutationResolver) ChangeDetails(ctx context.Context, input graphql.UpdateUser) (*db.User, error) {
	u, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	if input.FirstName != nil {
		u.FirstName = null.StringFrom(input.FirstName.String)
	}
	if input.LastName != nil {
		u.LastName = null.StringFrom(input.LastName.String)
	}
	if input.Email != nil {
		u.Email = null.StringFrom(strings.ToLower(input.Email.String))
	}
	if input.MobilePhone != nil {
		u.MobilePhone = null.StringFrom(input.MobilePhone.String)
	}
	if input.AffiliateOrg != nil {
		u.AffiliateOrg = *input.AffiliateOrg
	}

	_, err = r.UserStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	r.RecordUserActivity(ctx, "Updated User", graphql.ObjectTypeSelf, nil, nil)

	return u, nil
}

// UserCreate creates a user
func (r *mutationResolver) UserCreate(ctx context.Context, input graphql.UpdateUser) (*db.User, error) {
	u := &db.User{}
	referee := &db.User{}

	if input.ReferredByCode != nil {
		dat, err := r.UserStore.GetByReferralCode(input.ReferredByCode.String)
		if err != nil {
			return nil, terror.New(err, "No referee found")
		}
		referee = dat
	}

	if input.FirstName != nil {
		u.FirstName = null.StringFrom(input.FirstName.String)
	}
	if input.LastName != nil {
		u.LastName = null.StringFrom(input.LastName.String)
	}
	if input.Email != nil {
		u.Email = null.StringFrom(strings.ToLower(input.Email.String))
	}
	if input.MobilePhone != nil {
		u.MobilePhone = null.StringFrom(input.MobilePhone.String)
	}
	if input.RoleID != nil {
		u.RoleID = input.RoleID.String
	} else {
		// RoleID is required when creating a new user
		return nil, terror.New(fmt.Errorf("No role provided"), "")
	}
	if input.AffiliateOrg != nil {
		u.AffiliateOrg = *input.AffiliateOrg
	}

	// Set/Generate password
	password := ""
	if input.Password != nil {
		password = input.Password.String
	} else {
		g, err := uuid.NewV4()
		if err != nil {
			return nil, terror.New(err, "failed to generate password")
		}
		password = g.String()
		fmt.Println("generated", password)
	}
	err := r.CheckPasswordStrength(ctx, password)
	if err != nil {
		return nil, terror.New(err, "")
	}
	hashed := crypto.HashPassword(password)
	u.PasswordHash = hashed
	if r.Config.Skip.UserVerification {
		u.Verified = true
	}

	// Set/Genetate referral code
	u.ReferralCode = null.StringFrom(helpers.GenerateID(7))
	// referredByID := referee.ID

	// Create user
	created, err := r.UserStore.Insert(u)
	if err != nil {
		return nil, terror.New(err, "create user")
	}

	// Get Referral count
	count, err := r.ReferralStore.Count()
	if err != nil {
		return nil, terror.New(err, "create referral: Error while fetching referral count from db")
	}

	// Define Referrel model
	ref := &db.Referral{
		Code: fmt.Sprintf("T%05d", count+1),
	}

	refID, _ := uuid.NewV4()
	ref.ID = refID.String()
	ref.UserID = created.ID
	ref.ReferredByID = null.StringFrom(referee.ID)
	ref.IsRedemmed = false

	// Create referral
	_, err = r.ReferralStore.Insert(ref)
	if err != nil {
		return nil, terror.New(err, "create user")
	}

	r.RecordUserActivity(ctx, "Created User", graphql.ObjectTypeUser, &created.ID, &created.Email.String)

	return created, nil
}

// UserUpdate updates a user
func (r *mutationResolver) UserUpdate(ctx context.Context, id string, input graphql.UpdateUser) (*db.User, error) {
	// Get current user
	me, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Get User
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.UserStore.Get(userUUID)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	// Check archived state
	if u.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived user")
	}

	// Prevent updating a user who has a higher tier role than yourself
	if me.R.Role.Tier >= u.R.Role.Tier {
		return nil, terror.New(terror.ErrUnauthorized, "")
	}

	// Prevent using blank email address
	if input.Email.String == "" {
		return nil, terror.New(terror.ErrInvalidInput, "cannot use blank email address")
	}

	// Update User
	if input.FirstName != nil {
		u.FirstName = null.StringFrom(input.FirstName.String)
	}
	if input.LastName != nil {
		u.LastName = null.StringFrom(input.LastName.String)
	}
	if input.Email != nil {
		u.Email = null.StringFrom(strings.ToLower(input.Email.String))
	}
	if input.MobilePhone != nil {
		u.MobilePhone = null.StringFrom(input.MobilePhone.String)
	}
	if input.RoleID != nil && input.RoleID.String != u.RoleID {
		// Prevent giving user a role that's greater than your own
		settingRoleUUID, err := uuid.FromString(input.RoleID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		settingRole, err := r.RoleStore.Get(settingRoleUUID)
		if err != nil {
			return nil, terror.New(err, "update user")
		}

		if me.R.Role.Tier > settingRole.Tier {
			return nil, terror.New(terror.ErrUnauthorized, "")
		}

		u.RoleID = input.RoleID.String
	}
	if input.AffiliateOrg != nil {
		u.AffiliateOrg = *input.AffiliateOrg
	}
	if input.Password != nil {
		err := r.CheckPasswordStrength(ctx, input.Password.String)
		if err != nil {
			return nil, terror.New(err, "")
		}
		hashed := crypto.HashPassword(input.Password.String)
		u.PasswordHash = hashed
	}

	updated, err := r.UserStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	r.RecordUserActivity(ctx, "Updated User", graphql.ObjectTypeUser, &updated.ID, &updated.Email.String)

	return updated, nil
}

// UserArchive archives a user
func (r *mutationResolver) UserArchive(ctx context.Context, id string) (*db.User, error) {
	// Get User
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive User
	updated, err := r.UserStore.Archive(userUUID)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	r.RecordUserActivity(ctx, "Archived User", graphql.ObjectTypeUser, &updated.ID, &updated.ID)

	return updated, nil
}

// UserUnarchive archives a user
func (r *mutationResolver) UserUnarchive(ctx context.Context, id string) (*db.User, error) {
	// Get User
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive User
	updated, err := r.UserStore.Unarchive(userUUID)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	r.RecordUserActivity(ctx, "Archived User", graphql.ObjectTypeUser, &updated.ID, &updated.ID)

	return updated, nil
}
