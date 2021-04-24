package genesis

import (
	"context"
	"fmt"
	"genesis/db"
	"genesis/graphql"
	"log"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
)

////////////////
//  Resolver  //
////////////////

// UserPurchaseActivity resolver
func (r *Resolver) UserPurchaseActivity() graphql.UserPurchaseActivityResolver {
	return &userPurchaseActivityResolver{r}
}

type userPurchaseActivityResolver struct{ *Resolver }

func (r *userPurchaseActivityResolver) Product(ctx context.Context, obj *db.UserPurchaseActivity) (*db.Product, error) {
	result, err := r.UserPurchaseActivityStore.GetProduct(obj.ProductID.String)
	if err != nil {
		return nil, terror.New(err, "get product")
	}
	return result, nil
}

func (r *userPurchaseActivityResolver) User(ctx context.Context, obj *db.UserPurchaseActivity) (*db.User, error) {
	result, err := r.UserPurchaseActivityStore.GetUser(obj.UserID)
	if err != nil {
		return nil, terror.New(err, "get user")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserPurchaseActivity(ctx context.Context, id *string) (*db.UserPurchaseActivity, error) {
	purchaseUUID, err := uuid.FromString(*id)
	purchase, err := r.UserPurchaseActivityStore.Get(purchaseUUID)
	if err != nil {
		return nil, terror.New(err, "get purchase")
	}
	return purchase, nil
}

func (r *queryResolver) UserPurchaseActivities(ctx context.Context, search graphql.SearchFilter, limit int, offset int, userID *string) (*graphql.UserPurchaseActivityResult, error) {
	total, purchases, err := r.UserPurchaseActivityStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list purchase")
	}

	result := &graphql.UserPurchaseActivityResult{
		UserPurchaseActivities: purchases,
		Total:                  int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// UserPurchaseActivityCreate creates an purchase
func (r *mutationResolver) UserPurchaseActivityCreate(ctx context.Context, input graphql.UpdateUserPurchaseActivity) (*db.UserPurchaseActivity, error) {
	// Get UserPurchaseActivity count (for UserPurchaseActivity Code)
	count, err := r.UserPurchaseActivityStore.Count()
	if err != nil {
		return nil, terror.New(err, "create purchase: Error while fetching purchase count from db")
	}

	// Create UserPurchaseActivity
	t := &db.UserPurchaseActivity{
		Code: fmt.Sprintf("P%05d", count+1),
	}

	// get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create userTask: Error while fetching user")
	}

	// Get user object
	user, err := r.UserStore.Get(userID)
	if err != nil {
		return nil, terror.New(err, "Error while getting user")
	}

	purchaseID, _ := uuid.NewV4()
	t.ID = purchaseID.String()

	if input.ProductID == nil {
		return nil, terror.New(terror.ErrParse, "create purchase: product id is not provided")
	}

	productUUID, _ := uuid.FromString(input.ProductID.String)

	// Get product
	product, err := r.ProductStore.Get(productUUID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create purchase: no product found with given id")
	}

	if !product.SkuID.Valid {
		return nil, terror.New(terror.ErrParse, "create purchase: product is not bound to any sku")
	}
	if product.IsPointBound {
		return nil, terror.New(terror.ErrParse, "create purchase: product is point bound")
	}
	if product.IsClosed {
		return nil, terror.New(terror.ErrParse, "create purchase: product is closed")
	}

	// Get sku
	skuUUID, _ := uuid.FromString(product.SkuID.String)
	sku, err := r.SKUStore.Get(skuUUID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create purchase: no sku found with given id")
	}

	if sku.IsPointBound {
		// if user.WalletPoints < sku.PurchasePoints {
		// 	return nil, terror.New(terror.ErrParse, "create purchase: not enough points in the wallet")
		// }
		// user.WalletPoints -= sku.PurchasePoints
		return nil, terror.New(terror.ErrParse, "create purchase: product is point bound")
	}

	// Add loyalty points to user's wallet
	user.WalletPoints += sku.LoyaltyPoints

	t.UserID = userID.String()
	t.ProductID = *input.ProductID
	t.LoyaltyPoints = sku.LoyaltyPoints

	created, err := r.UserPurchaseActivityStore.Insert(t)
	if err != nil {
		return nil, terror.New(err, "create user purchase")
	}

	// Update product
	product.IsClosed = true
	_, err = r.ProductStore.Update(product)
	if err != nil {
		return nil, terror.New(err, "Error while updating user")
	}

	// Update user
	_, err = r.UserStore.Update(user)
	if err != nil {
		return nil, terror.New(err, "Error while updating user")
	}

	// Get referral if exists
	referral, err := r.ReferralStore.GetByUserID(user.ID)
	if err != nil {
		referral = nil
	}

	if referral != nil && !referral.IsRedemmed {
		// Update wallet points of the referee with 10 points
		refID, _ := uuid.FromString(referral.ReferredByID.String)
		referee, err := r.UserStore.Get(refID)
		if err != nil {
			log.Println("Referee not found in users db")
		}
		referee.WalletPoints += 10

		// Update referee
		_, err = r.UserStore.Update(referee)
		if err != nil {
			fmt.Println("Error while updating referee")
		}

		referral.IsRedemmed = true
		_, err = r.ReferralStore.Update(referral)
		if err != nil {
			fmt.Println("Error while updating referral")
		}

		// Update WalletTransactions
		wtID, _ := uuid.NewV4()
		wt := &db.WalletTransaction{
			ID:            wtID.String(),
			UserID:        referee.ID,
			LoyaltyPoints: 10,
			IsCredit:      true,
			Message:       fmt.Sprintf("Loyalty points awarded by referral %v %v", user.FirstName.String, user.LastName.String),
		}

		_, err = r.UserStore.InsertWalletTransaction(wt)
		if err != nil {
			return nil, terror.New(err, "create wallet transaction")
		}
	}

	// Insert wallet transaction
	wtID, _ := uuid.NewV4()
	wt := &db.WalletTransaction{
		ID:     wtID.String(),
		UserID: user.ID,
	}

	// if sku.IsPointBound {
	// 	wt.LoyaltyPoints = sku.PurchasePoints
	// 	wt.IsCredit = false
	// 	wt.Message = fmt.Sprintf("Loyalty points deducted by purchase of product %v, %v", product.Code, sku.Name)
	// }

	wt.LoyaltyPoints = sku.LoyaltyPoints
	wt.IsCredit = true
	wt.Message = fmt.Sprintf("Loyalty points awarded by purchase of product %v %v", product.Code, sku.Name)

	_, err = r.UserStore.InsertWalletTransaction(wt)
	if err != nil {
		return nil, terror.New(err, "create wallet transaction")
	}

	// r.RecordUserActivity(ctx, "Created User Purchase", graphql.ObjectTypeSku, &created.ID, &created.Code)

	return created, nil
}

// UserPurchaseActivityUpdate updates a purchase
func (r *mutationResolver) UserPurchaseActivityUpdate(ctx context.Context, id string, input graphql.UpdateUserPurchaseActivity) (*db.UserPurchaseActivity, error) {
	// Get UserPurchaseActivity
	purchaseUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	t, err := r.UserPurchaseActivityStore.Get(purchaseUUID)
	if err != nil {
		return nil, terror.New(err, "update purchase")
	}

	updated, err := r.UserPurchaseActivityStore.Update(t)
	if err != nil {
		return nil, terror.New(err, "update purchase")
	}

	return updated, nil
}
