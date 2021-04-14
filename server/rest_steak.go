package genesis

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"genesis/blockchain"
	"genesis/db"
	"genesis/helpers"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
)

// helper to create string pointer
func stringPointer(str string) *string {
	return &str
}

// restWriteError writes a http errors to the ResponseWriter
func restWriteError(ctx context.Context, w http.ResponseWriter, httpStatus int, err error, friendlyMessage ...string) {
	terror.Echo(err)
	http.Error(w, fmt.Sprintf(`{"success":false, "message":"%s"}`, err.Error()), httpStatus)
}

// helper to write http response in json.
func restWriteJSON(ctx context.Context, w http.ResponseWriter, httpStatus int, dat interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	w.WriteHeader(httpStatus)
	err := encoder.Encode(dat)
	if err != nil {
		err = terror.New(fmt.Errorf("failed JSON encode"), "")
		restWriteError(ctx, w, http.StatusInternalServerError, err)
		return
	}
}

// SteakView redirects to view page with distributorCode added to the arguments /view?productID=UUID&distributorCode=ABCD
func SteakView(
	ProductStore ProductStorer,
	DistributorStore DistributorStorer,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Get URL query
		productID := r.URL.Query().Get("productID")
		productUUID, err := uuid.FromString(productID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid product uuid"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get product
		product, err := ProductStore.Get(productUUID)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find product"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get Distributor Code
		distributorCode := ""
		if product.DistributorID.Valid {
			distributorUUID, err := uuid.FromString(product.DistributorID.String)
			if err != nil {
				err = terror.New(fmt.Errorf("invalid distributor uuid"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}

			distributor, err := DistributorStore.Get(distributorUUID)
			if err != nil {
				err = terror.New(fmt.Errorf("failed to find distributor"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}
			distributorCode = distributor.Code
		}

		redirect := fmt.Sprintf("/view?productID=%s&distributorCode=%s", productID, distributorCode)

		fmt.Println(http.StatusFound, redirect)
		http.Redirect(w, r, redirect, http.StatusFound)
	}

	return fn
}

type steakInfoResponse struct {
	// always show
	Success bool `json:"success"`
	// if something fails
	Message *string `json:"message,omitempty"`
	// if success
	BonusPointsExpiry *time.Time `json:"bonusPointsExpiry,omitempty"`
	BonusPoints       *int       `json:"bonusPoints,omitempty"`
	BasePoints        *int       `json:"basePoints,omitempty"`
	CloseID           *string    `json:"closeID,omitempty"`
	DistributorCode   *string    `json:"distributorCode,omitempty"`
}

/*
### note

product_id, qr code on the outside package
register_id, qr code below the meat

steakID == register_id

client call them steak token
internally, we call them register id

bonus points is the product points
base points is the sku points

*/

// SteakDetail shows steak information, QR2
func SteakDetail(
	ProductStore ProductStorer,
	SKUStore SKUStorer,
	LoyaltyStore LoyaltyStorer,
	DistributorStore DistributorStorer,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// get url query
		registerID := r.URL.Query().Get("steakID")
		productID := r.URL.Query().Get("productID")

		// Get product
		if registerID == "" {
			err = terror.New(fmt.Errorf("steakID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		registerUUID, err := uuid.FromString(registerID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid steak uuid"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		product, err := ProductStore.GetByRegisterID(registerUUID)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find product"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// sanity check, making sure WeChat Mini Program scan QR1 and QR2 in correct order
		if productID == "" {
			err = terror.New(fmt.Errorf("productID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		if product.ID != productID {
			err = terror.New(fmt.Errorf("incorrect productID"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Check if product is already registered,
		// if not, generate new close id if not already existed
		registered, err := ProductStore.Registered(product)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to check product registered status"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}
		if !registered {
			uid, err := uuid.NewV4()
			if err != nil {
				err = terror.New(fmt.Errorf("error generating uuid"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}

			product.CloseRegisterID = null.StringFrom(uid.String())
			product, err = ProductStore.Update(product)
			if err != nil {
				err = terror.New(fmt.Errorf("error updating product"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}
		}

		// Get DistributorCode
		var distributorCode *string = nil
		if product.DistributorID.Valid {
			distributorUUID, err := uuid.FromString(product.DistributorID.String)
			if err != nil {
				err = terror.New(fmt.Errorf("invalid distributor uuid"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}

			distributor, err := DistributorStore.Get(distributorUUID)
			if err != nil {
				err = terror.New(fmt.Errorf("failed to find distributor"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
				return
			}
			distributorCode = &distributor.Code
		}

		// Get product SKU points
		skuPoints, err := getSKUPoints(SKUStore, product.SkuID)
		if err != nil {
			err = terror.New(fmt.Errorf("cannot find sku"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
		}

		// Response
		response := steakInfoResponse{
			Success: true,
			// these 2 must show when success
			BonusPoints:       &product.LoyaltyPoints,
			BonusPointsExpiry: &product.LoyaltyPointsExpire,
			BasePoints:        &skuPoints,
			DistributorCode:   distributorCode,
		}
		// only show if not yet registered
		if !registered {
			response.CloseID = stringPointer(product.CloseRegisterID.String)
		}
		restWriteJSON(r.Context(), w, http.StatusOK, response)
	}

	return fn
}

type steakCloseRequest struct {
	SteakID  string `json:"steakID"`
	WeChatID string `json:"weChatID"`
	CloseID  string `json:"closeID"`
}

type steakCloseResponse struct {
	// always show
	Success bool `json:"success"`
	// if something fails
	Message *string `json:"message,omitempty"`
	// if success
	RegisterID  *string `json:"steakID,omitempty"`
	UserID      *string `json:"customerID,omitempty"`
	WeChatID    *string `json:"weChatID,omitempty"`
	PointsGiven *int    `json:"pointsGiven,omitempty"`
}

// SteakClose register a product and reply points to wechat, QR2+close
func SteakClose(
	UserStore UserStorer,
	ProductStore ProductStorer,
	RoleStore RoleStorer,
	OrganisationStore OrganisationStorer,
	LoyaltyStore LoyaltyStorer,
	SKUStore SKUStorer,
	TransactionStore TransactionStorer,
	TrackActionStore TrackActionStorer,
	Blk *blockchain.Service,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		req := &steakCloseRequest{}
		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid json"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		// sanity check
		// steak id
		registerID := req.SteakID
		if registerID == "" {
			err = terror.New(fmt.Errorf("steakID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		registerUUID, err := uuid.FromString(registerID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid steakID"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		// close id
		closeID := req.CloseID
		if closeID == "" {
			err = terror.New(fmt.Errorf("closeID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		closeUUID, err := uuid.FromString(closeID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid closeID"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		// wechat
		weChatID := req.WeChatID
		// officially wechat id is 6-20 long, but we add extra incase some weird length rune length counting bug,
		// or change of official length. we make it long enough to fit, but short enough that it is less likely to be abused.
		if len(weChatID) > 36 {
			err = terror.New(fmt.Errorf("weChatID too long"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get product
		product, err := ProductStore.GetByRegisterID(registerUUID)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find product"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		if product.CloseRegisterID.String != closeUUID.String() {
			err = terror.New(fmt.Errorf("closeID not match"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get product SKU points
		skuPoints, err := getSKUPoints(SKUStore, product.SkuID)
		if err != nil {
			err = terror.New(fmt.Errorf("cannot find sku"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
		}

		// Check if product is already registered
		registered, err := ProductStore.Registered(product)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to check product registered status"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}
		if registered {
			err = terror.New(fmt.Errorf("product already registered"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get user
		var user *db.User
		if weChatID != "" {
			user, err = UserStore.GetByWechatID(weChatID)
			if err != nil {
				// Create new user
				role, err := RoleStore.GetByName("CONSUMER")
				if err != nil {
					err = terror.New(fmt.Errorf("failed to create user (1)"), "")
					restWriteError(r.Context(), w, http.StatusBadRequest, err)
					return
				}

				user = &db.User{
					WechatID: null.StringFrom(weChatID),
					RoleID:   role.ID,
				}
				user, err = UserStore.Insert(user)
				if err != nil {
					err = terror.New(fmt.Errorf("failed to create user (3)"), "")
					restWriteError(r.Context(), w, http.StatusBadRequest, err)
					return
				}
			}
		}

		// Check bonus expire
		bonus := 0
		if product.LoyaltyPoints > 0 && time.Now().Before(product.LoyaltyPointsExpire) {
			bonus = product.LoyaltyPoints
		}

		// Create product transaction
		action, err := TrackActionStore.GetByCode(trackRegistered)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid track action code"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		changedByName := ""
		if user != nil {
			changedByName = helpers.LimitName(user.WechatID.String)
		}

		_, err = TransactionStore.InsertByProduct(product, action, user, changedByName, nil)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to register product (1)"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}

		// Create loyalty activity
		pointsGiven := 0
		if user != nil {
			pointsGiven := skuPoints + bonus
			activity := &db.UserLoyaltyActivity{
				UserID:    user.ID,
				ProductID: null.StringFrom(product.ID),
				Amount:    pointsGiven,
				Bonus:     bonus,
			}
			activity, err = LoyaltyStore.Insert(activity)
			if err != nil {
				err = terror.New(fmt.Errorf("failed to register product (2)"), "")
				restWriteError(r.Context(), w, http.StatusInternalServerError, err)
				return
			}
		}

		// Archive product
		productUUID, err := uuid.FromString(product.ID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid product id"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		product, err = ProductStore.Archive(productUUID)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to archive product"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		userID := ""
		if user != nil {
			userID = user.ID
		}
		response := &steakCloseResponse{
			// success
			Success: true,
			Message: nil,
			// these 4 must show
			RegisterID:  stringPointer(registerID),
			UserID:      &userID,
			WeChatID:    stringPointer(weChatID),
			PointsGiven: &pointsGiven,
		}

		restWriteJSON(r.Context(), w, http.StatusOK, response)
	}

	return fn
}

type steakFinalTemplateData struct {
	PointsGiven int
}

// SteakFinal display a html page of thanking customer
func SteakFinal(
	ProductStore ProductStorer,
	LoyaltyStore LoyaltyStorer,
	SKUStore SKUStorer,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Detect demo mode
		isDemo := false
		if r.URL.Query().Get("demo") == "true" {
			isDemo = true
		}

		// Get product
		registerID := r.URL.Query().Get("steakID")
		if registerID == "" {
			err = terror.New(fmt.Errorf("steakID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		registerUUID, err := uuid.FromString(registerID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid steak uuid"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		product, err := ProductStore.GetByRegisterID(registerUUID)
		if err != nil {
			err = terror.New(fmt.Errorf("cannot find product"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Check if product is already registered
		registered, err := ProductStore.Registered(product)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to check product registered status"), "")
			restWriteError(r.Context(), w, http.StatusInternalServerError, err)
			return
		}
		if !registered && !isDemo {
			err = terror.New(fmt.Errorf("product not yet registered"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		productUUID, err := uuid.FromString(product.ID)
		if err != nil {
			err = terror.New(fmt.Errorf("invalid product uuid"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		pointsGiven := 0

		loyalty, err := LoyaltyStore.GetByProductID(productUUID)
		if err != nil && err != sql.ErrNoRows {
			err = terror.New(fmt.Errorf("cannot find data"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}
		if loyalty != nil {
			pointsGiven = loyalty.Amount
		}

		if isDemo {
			// get sku points
			skuPoints, err := getSKUPoints(SKUStore, product.SkuID)
			if err != nil {
				err = terror.New(fmt.Errorf("cannot find sku point"), "")
				restWriteError(r.Context(), w, http.StatusBadRequest, err)
			}

			// get product points
			var productPoints int
			if product.LoyaltyPoints > 0 && time.Now().Before(product.LoyaltyPointsExpire) {
				productPoints = product.LoyaltyPoints
			}

			// calc demo total points
			pointsGiven = skuPoints + productPoints
		}

		redirect := fmt.Sprintf("/view?productID=%s&registered=true&pointsGiven=%d", product.ID, pointsGiven)
		http.Redirect(w, r, redirect, http.StatusFound)
	}

	return fn
}

func getSKUPoints(SKUStore SKUStorer, skuID null.String) (int, error) {
	if !skuID.Valid {
		return 0, terror.New(fmt.Errorf("invalid skuID"), "")
	}

	skuUUID, err := uuid.FromString(skuID.String)
	if err != nil {
		return 0, terror.New(err, "")
	}

	sku, err := SKUStore.Get(skuUUID)
	if err != nil {
		return 0, terror.New(err, "")
	}

	return sku.LoyaltyPoints, nil
}
