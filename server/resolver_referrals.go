package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"
)

////////////////
//  Resolver  //
////////////////

// Referral resolver
func (r *Resolver) Referral() graphql.ReferralResolver {
	return &referralResolver{r}
}

type referralResolver struct{ *Resolver }

func (r *referralResolver) ReferredByID(ctx context.Context, obj *db.Referral) (string, error) {
	ref, err := r.ReferralStore.GetByUserID(obj.UserID)
	if err != nil {
		return "", terror.New(err, "")
	}
	return ref.ReferredByID.String, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Referral(ctx context.Context, userID *string) (*db.Referral, error) {
	referral, err := r.ReferralStore.GetByUserID(*userID)
	if err != nil {
		return nil, terror.New(err, "get referral")
	}
	return referral, nil
}

func (r *queryResolver) Referrals(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.ReferralsResult, error) {
	total, referrals, err := r.ReferralStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list referral")
	}

	result := &graphql.ReferralsResult{
		Referrals: referrals,
		Total:     int(total),
	}

	return result, nil
}
