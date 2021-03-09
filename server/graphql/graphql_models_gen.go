// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"genesis/db"
	"io"
	"strconv"
	"time"

	"github.com/volatiletech/null"
)

type BatchActionInput struct {
	Str      *null.String `json:"str"`
	No       *null.Int    `json:"no"`
	DateTime *null.Time   `json:"dateTime"`
	Bool     *null.Bool   `json:"bool"`
}

type CartonResult struct {
	Cartons []*db.Carton `json:"cartons"`
	Total   int          `json:"total"`
}

type ConsumersResult struct {
	Consumers []*db.User `json:"consumers"`
	Total     int        `json:"total"`
}

type ContainerResult struct {
	Containers []*db.Container `json:"containers"`
	Total      int             `json:"total"`
}

type ContractResult struct {
	Contracts []*db.Contract `json:"contracts"`
	Total     int            `json:"total"`
	CreatedAt time.Time      `json:"createdAt"`
}

type CreateCarton struct {
	PalletID    *null.String `json:"palletID"`
	Quantity    int          `json:"quantity"`
	Description string       `json:"description"`
}

type CreateContainer struct {
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
}

type CreateOrder struct {
	ContractID *null.String `json:"contractID"`
	SkuID      *null.String `json:"skuID"`
	Quantity   int          `json:"quantity"`
}

type CreatePallet struct {
	ContainerID *null.String `json:"containerID"`
	Quantity    int          `json:"quantity"`
	Description string       `json:"description"`
}

type DistributorResult struct {
	Distributors []*db.Distributor `json:"distributors"`
	Total        int               `json:"total"`
	CreatedAt    time.Time         `json:"createdAt"`
}

type GetObjectResponse struct {
	Product   *db.Product   `json:"product"`
	Carton    *db.Carton    `json:"carton"`
	Pallet    *db.Pallet    `json:"pallet"`
	Container *db.Container `json:"container"`
}

type GetObjectsRequest struct {
	ProductIDs   []string `json:"productIDs"`
	CartonIDs    []string `json:"cartonIDs"`
	PalletIDs    []string `json:"palletIDs"`
	ContainerIDs []string `json:"containerIDs"`
}

type GetObjectsResponse struct {
	Products   []*db.Product   `json:"products"`
	Cartons    []*db.Carton    `json:"cartons"`
	Pallets    []*db.Pallet    `json:"pallets"`
	Containers []*db.Container `json:"containers"`
}

type LatestTransactionInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderResult struct {
	Orders []*db.Order `json:"orders"`
	Total  int         `json:"total"`
}

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
}

type PalletResult struct {
	Pallets []*db.Pallet `json:"pallets"`
	Total   int          `json:"total"`
}

type ProductResult struct {
	Products []*db.Product `json:"products"`
	Total    int           `json:"total"`
}

type RecordTransactionInput struct {
	TrackActionCode     string       `json:"trackActionCode"`
	ProductIDs          []string     `json:"productIDs"`
	CartonIDs           []string     `json:"cartonIDs"`
	PalletIDs           []string     `json:"palletIDs"`
	ContainerIDs        []string     `json:"containerIDs"`
	ProductScanTimes    []*time.Time `json:"productScanTimes"`
	CartonScanTimes     []*time.Time `json:"cartonScanTimes"`
	PalletScanTimes     []*time.Time `json:"palletScanTimes"`
	ContainerScanTimes  []*time.Time `json:"containerScanTimes"`
	CartonPhotoBlobIDs  []string     `json:"cartonPhotoBlobIDs"`
	ProductPhotoBlobIDs []string     `json:"productPhotoBlobIDs"`
	Memo                *null.String `json:"memo"`
	LocationGeohash     *null.String `json:"locationGeohash"`
	LocationName        *null.String `json:"locationName"`
}

type ReferralsResult struct {
	Referrals []*db.Referral `json:"referrals"`
	Total     int            `json:"total"`
}

type RequestToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RolesResult struct {
	Roles []*db.Role `json:"roles"`
	Total int        `json:"total"`
}

type SKUClone struct {
	Sku   *db.StockKeepingUnit `json:"sku"`
	Depth int                  `json:"depth"`
}

type SKUContentInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SKUResult struct {
	Skus  []*db.StockKeepingUnit `json:"skus"`
	Total int                    `json:"total"`
}

type SearchFilter struct {
	Search  *null.String  `json:"search"`
	Filter  *FilterOption `json:"filter"`
	SortBy  *SortByOption `json:"sortBy"`
	SortDir *SortDir      `json:"sortDir"`
}

type TasksResult struct {
	Tasks []*db.Task `json:"tasks"`
	Total int        `json:"total"`
}

type TickerInfo struct {
	LastTick     time.Time `json:"lastTick"`
	TickInterval int       `json:"tickInterval"`
}

type TrackActionResult struct {
	TrackActions []*db.TrackAction `json:"trackActions"`
	Total        int               `json:"total"`
}

type TransactionPhotos struct {
	CartonPhoto  *db.Blob `json:"cartonPhoto"`
	ProductPhoto *db.Blob `json:"productPhoto"`
}

type TransactionsResult struct {
	Transactions []*db.Transaction `json:"transactions"`
	Total        int               `json:"total"`
}

type UpdateCarton struct {
	Code        *null.String `json:"code"`
	Weight      *null.String `json:"weight"`
	PalletID    *null.String `json:"palletID"`
	ProcessedAt *null.Time   `json:"processedAt"`
	Description *null.String `json:"description"`
	MeatType    *null.String `json:"meatType"`
}

type UpdateCategory struct {
	Name string `json:"name"`
}

type UpdateContainer struct {
	Code        *null.String `json:"code"`
	Description *null.String `json:"description"`
}

type UpdateContract struct {
	Name         *null.String `json:"name"`
	Description  *null.String `json:"description"`
	SupplierName *null.String `json:"supplierName"`
	DateSigned   *null.Time   `json:"dateSigned"`
}

type UpdateDistributor struct {
	Name *null.String `json:"name"`
	Code *null.String `json:"code"`
}

type UpdateLoyaltyPoints struct {
	ProductIDs          []string     `json:"productIDs"`
	CartonID            *null.String `json:"cartonID"`
	PalletID            *null.String `json:"palletID"`
	ContainerID         *null.String `json:"containerID"`
	LoyaltyPoints       int          `json:"loyaltyPoints"`
	LoyaltyPointsExpire time.Time    `json:"loyaltyPointsExpire"`
}

type UpdateOrder struct {
	Code *null.String `json:"code"`
}

type UpdatePallet struct {
	Code        *null.String `json:"code"`
	Description *null.String `json:"description"`
	ContainerID *null.String `json:"containerID"`
}

type UpdateProduct struct {
	Code                 *null.String `json:"code"`
	CartonID             *null.String `json:"cartonID"`
	OrderID              *null.String `json:"orderID"`
	SkuID                *null.String `json:"skuID"`
	ContractID           *null.String `json:"contractID"`
	DistributorID        *null.String `json:"distributorID"`
	LoyaltyPoints        *null.Int    `json:"loyaltyPoints"`
	LoyaltyPointsExpire  *null.Time   `json:"loyaltyPointsExpire"`
	InheritCartonHistory *null.Bool   `json:"inheritCartonHistory"`
	Description          *null.String `json:"description"`
}

type UpdateProductCategory struct {
	Name string `json:"name"`
}

type UpdateRole struct {
	Name           *null.String `json:"name"`
	Permissions    []string     `json:"permissions"`
	TrackActionIDs []string     `json:"trackActionIDs"`
}

type UpdateSku struct {
	Name              *null.String             `json:"name"`
	Code              *null.String             `json:"code"`
	Description       *null.String             `json:"description"`
	Weight            *null.Int                `json:"weight"`
	WeightUnit        *null.String             `json:"weightUnit"`
	Price             *null.Int                `json:"price"`
	Currency          *null.String             `json:"currency"`
	IsBeef            *null.Bool               `json:"isBeef"`
	IsRetailSku       *null.Bool               `json:"isRetailSku"`
	IsPointSku        *null.Bool               `json:"isPointSku"`
	IsAppSku          *null.Bool               `json:"isAppSku"`
	IsMiniappSku      *null.Bool               `json:"isMiniappSku"`
	LoyaltyPoints     *null.Int                `json:"loyaltyPoints"`
	MasterPlanBlobID  *null.String             `json:"masterPlanBlobID"`
	VideoBlobID       *null.String             `json:"videoBlobID"`
	Urls              []*SKUContentInput       `json:"urls"`
	ProductInfo       []*SKUContentInput       `json:"productInfo"`
	PhotoBlobIDs      []string                 `json:"photoBlobIDs"`
	Categories        []*UpdateCategory        `json:"categories"`
	ProductCategories []*UpdateProductCategory `json:"productCategories"`
	CloneParentID     *null.String             `json:"cloneParentID"`
}

type UpdateSubtask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTask struct {
	Title             string           `json:"title"`
	Description       string           `json:"description"`
	LoyaltyPoints     int              `json:"loyaltyPoints"`
	IsTimeBound       bool             `json:"isTimeBound"`
	IsPeopleBound     bool             `json:"isPeopleBound"`
	IsProductRelevant bool             `json:"isProductRelevant"`
	IsFinal           *bool            `json:"isFinal"`
	FinishDate        *time.Time       `json:"finishDate"`
	MaximumPeople     int              `json:"maximumPeople"`
	SkuID             *null.String     `json:"skuID"`
	Subtasks          []*UpdateSubtask `json:"subtasks"`
}

type UpdateTrackAction struct {
	Name          *null.String `json:"name"`
	RequirePhotos []bool       `json:"requirePhotos"`
	NameChinese   *null.String `json:"nameChinese"`
	Private       *null.Bool   `json:"private"`
	Blockchain    *null.Bool   `json:"blockchain"`
}

type UpdateUser struct {
	Email          *null.String `json:"email"`
	FirstName      *null.String `json:"firstName"`
	LastName       *null.String `json:"lastName"`
	RoleID         *null.String `json:"roleID"`
	Password       *null.String `json:"password"`
	AffiliateOrg   *null.String `json:"affiliateOrg"`
	MobilePhone    *null.String `json:"mobilePhone"`
	ReferredByCode *null.String `json:"referredByCode"`
}

type UserActivityResult struct {
	UserActivities []*db.UserActivity `json:"userActivities"`
	Total          int                `json:"total"`
}

type UsersResult struct {
	Users []*db.User `json:"users"`
	Total int        `json:"total"`
}

type Action string

const (
	ActionArchive               Action = "Archive"
	ActionUnarchive             Action = "Unarchive"
	ActionSetSku                Action = "SetSKU"
	ActionSetOrder              Action = "SetOrder"
	ActionSetDistributor        Action = "SetDistributor"
	ActionSetContract           Action = "SetContract"
	ActionSetCarton             Action = "SetCarton"
	ActionSetPallet             Action = "SetPallet"
	ActionSetContainer          Action = "SetContainer"
	ActionDetachFromSku         Action = "DetachFromSKU"
	ActionDetachFromOrder       Action = "DetachFromOrder"
	ActionDetachFromDistributor Action = "DetachFromDistributor"
	ActionDetachFromContract    Action = "DetachFromContract"
	ActionDetachFromCarton      Action = "DetachFromCarton"
	ActionDetachFromPallet      Action = "DetachFromPallet"
	ActionDetachFromContainer   Action = "DetachFromContainer"
	ActionSetBonusLoyaltyPoints Action = "SetBonusLoyaltyPoints"
	ActionInheritCartonHistory  Action = "InheritCartonHistory"
)

var AllAction = []Action{
	ActionArchive,
	ActionUnarchive,
	ActionSetSku,
	ActionSetOrder,
	ActionSetDistributor,
	ActionSetContract,
	ActionSetCarton,
	ActionSetPallet,
	ActionSetContainer,
	ActionDetachFromSku,
	ActionDetachFromOrder,
	ActionDetachFromDistributor,
	ActionDetachFromContract,
	ActionDetachFromCarton,
	ActionDetachFromPallet,
	ActionDetachFromContainer,
	ActionSetBonusLoyaltyPoints,
	ActionInheritCartonHistory,
}

func (e Action) IsValid() bool {
	switch e {
	case ActionArchive, ActionUnarchive, ActionSetSku, ActionSetOrder, ActionSetDistributor, ActionSetContract, ActionSetCarton, ActionSetPallet, ActionSetContainer, ActionDetachFromSku, ActionDetachFromOrder, ActionDetachFromDistributor, ActionDetachFromContract, ActionDetachFromCarton, ActionDetachFromPallet, ActionDetachFromContainer, ActionSetBonusLoyaltyPoints, ActionInheritCartonHistory:
		return true
	}
	return false
}

func (e Action) String() string {
	return string(e)
}

func (e *Action) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Action(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Action", str)
	}
	return nil
}

func (e Action) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FilterOption string

const (
	FilterOptionAll                    FilterOption = "All"
	FilterOptionActive                 FilterOption = "Active"
	FilterOptionArchived               FilterOption = "Archived"
	FilterOptionProductWithoutOrder    FilterOption = "ProductWithoutOrder"
	FilterOptionProductWithoutCarton   FilterOption = "ProductWithoutCarton"
	FilterOptionProductWithoutSku      FilterOption = "ProductWithoutSKU"
	FilterOptionCartonWithoutPallet    FilterOption = "CartonWithoutPallet"
	FilterOptionPalletWithoutContainer FilterOption = "PalletWithoutContainer"
	FilterOptionSystem                 FilterOption = "System"
	FilterOptionBlockchain             FilterOption = "Blockchain"
	FilterOptionPending                FilterOption = "Pending"
)

var AllFilterOption = []FilterOption{
	FilterOptionAll,
	FilterOptionActive,
	FilterOptionArchived,
	FilterOptionProductWithoutOrder,
	FilterOptionProductWithoutCarton,
	FilterOptionProductWithoutSku,
	FilterOptionCartonWithoutPallet,
	FilterOptionPalletWithoutContainer,
	FilterOptionSystem,
	FilterOptionBlockchain,
	FilterOptionPending,
}

func (e FilterOption) IsValid() bool {
	switch e {
	case FilterOptionAll, FilterOptionActive, FilterOptionArchived, FilterOptionProductWithoutOrder, FilterOptionProductWithoutCarton, FilterOptionProductWithoutSku, FilterOptionCartonWithoutPallet, FilterOptionPalletWithoutContainer, FilterOptionSystem, FilterOptionBlockchain, FilterOptionPending:
		return true
	}
	return false
}

func (e FilterOption) String() string {
	return string(e)
}

func (e *FilterOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterOption", str)
	}
	return nil
}

func (e FilterOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ObjectType string

const (
	ObjectTypeSelf        ObjectType = "Self"
	ObjectTypeUser        ObjectType = "User"
	ObjectTypeRole        ObjectType = "Role"
	ObjectTypeSku         ObjectType = "Sku"
	ObjectTypeContract    ObjectType = "Contract"
	ObjectTypeOrder       ObjectType = "Order"
	ObjectTypeContainer   ObjectType = "Container"
	ObjectTypePallet      ObjectType = "Pallet"
	ObjectTypeCarton      ObjectType = "Carton"
	ObjectTypeDistributor ObjectType = "Distributor"
	ObjectTypeProduct     ObjectType = "Product"
	ObjectTypeTrackAction ObjectType = "TrackAction"
	ObjectTypeBlob        ObjectType = "Blob"
	ObjectTypeBlockchain  ObjectType = "Blockchain"
)

var AllObjectType = []ObjectType{
	ObjectTypeSelf,
	ObjectTypeUser,
	ObjectTypeRole,
	ObjectTypeSku,
	ObjectTypeContract,
	ObjectTypeOrder,
	ObjectTypeContainer,
	ObjectTypePallet,
	ObjectTypeCarton,
	ObjectTypeDistributor,
	ObjectTypeProduct,
	ObjectTypeTrackAction,
	ObjectTypeBlob,
	ObjectTypeBlockchain,
}

func (e ObjectType) IsValid() bool {
	switch e {
	case ObjectTypeSelf, ObjectTypeUser, ObjectTypeRole, ObjectTypeSku, ObjectTypeContract, ObjectTypeOrder, ObjectTypeContainer, ObjectTypePallet, ObjectTypeCarton, ObjectTypeDistributor, ObjectTypeProduct, ObjectTypeTrackAction, ObjectTypeBlob, ObjectTypeBlockchain:
		return true
	}
	return false
}

func (e ObjectType) String() string {
	return string(e)
}

func (e *ObjectType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObjectType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObjectType", str)
	}
	return nil
}

func (e ObjectType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Perm string

const (
	PermUserList                       Perm = "UserList"
	PermUserCreate                     Perm = "UserCreate"
	PermUserRead                       Perm = "UserRead"
	PermUserUpdate                     Perm = "UserUpdate"
	PermUserArchive                    Perm = "UserArchive"
	PermUserUnarchive                  Perm = "UserUnarchive"
	PermReferralList                   Perm = "ReferralList"
	PermReferralRead                   Perm = "ReferralRead"
	PermTaskList                       Perm = "TaskList"
	PermTaskCreate                     Perm = "TaskCreate"
	PermTaskRead                       Perm = "TaskRead"
	PermTaskUpdate                     Perm = "TaskUpdate"
	PermTaskArchive                    Perm = "TaskArchive"
	PermTaskUnarchive                  Perm = "TaskUnarchive"
	PermUserTaskList                   Perm = "UserTaskList"
	PermUserTaskCreate                 Perm = "UserTaskCreate"
	PermUserTaskRead                   Perm = "UserTaskRead"
	PermUserTaskUpdate                 Perm = "UserTaskUpdate"
	PermUserTaskArchive                Perm = "UserTaskArchive"
	PermUserTaskUnarchive              Perm = "UserTaskUnarchive"
	PermOrganisationList               Perm = "OrganisationList"
	PermOrganisationCreate             Perm = "OrganisationCreate"
	PermOrganisationRead               Perm = "OrganisationRead"
	PermOrganisationUpdate             Perm = "OrganisationUpdate"
	PermOrganisationArchive            Perm = "OrganisationArchive"
	PermOrganisationUnarchive          Perm = "OrganisationUnarchive"
	PermRoleList                       Perm = "RoleList"
	PermRoleCreate                     Perm = "RoleCreate"
	PermRoleRead                       Perm = "RoleRead"
	PermRoleUpdate                     Perm = "RoleUpdate"
	PermRoleArchive                    Perm = "RoleArchive"
	PermRoleUnarchive                  Perm = "RoleUnarchive"
	PermSKUList                        Perm = "SKUList"
	PermSKUCreate                      Perm = "SKUCreate"
	PermSKURead                        Perm = "SKURead"
	PermSKUUpdate                      Perm = "SKUUpdate"
	PermSKUArchive                     Perm = "SKUArchive"
	PermSKUUnarchive                   Perm = "SKUUnarchive"
	PermCategoryList                   Perm = "CategoryList"
	PermCategoryCreate                 Perm = "CategoryCreate"
	PermCategoryRead                   Perm = "CategoryRead"
	PermCategoryUpdate                 Perm = "CategoryUpdate"
	PermCategoryArchive                Perm = "CategoryArchive"
	PermCategoryUnarchive              Perm = "CategoryUnarchive"
	PermProductCategoryList            Perm = "ProductCategoryList"
	PermProductCategoryCreate          Perm = "ProductCategoryCreate"
	PermProductCategoryRead            Perm = "ProductCategoryRead"
	PermProductCategoryUpdate          Perm = "ProductCategoryUpdate"
	PermProductCategoryArchive         Perm = "ProductCategoryArchive"
	PermProductCategoryUnarchive       Perm = "ProductCategoryUnarchive"
	PermContainerList                  Perm = "ContainerList"
	PermContainerRead                  Perm = "ContainerRead"
	PermContainerCreate                Perm = "ContainerCreate"
	PermContainerUpdate                Perm = "ContainerUpdate"
	PermContainerArchive               Perm = "ContainerArchive"
	PermContainerUnarchive             Perm = "ContainerUnarchive"
	PermPalletList                     Perm = "PalletList"
	PermPalletRead                     Perm = "PalletRead"
	PermPalletCreate                   Perm = "PalletCreate"
	PermPalletUpdate                   Perm = "PalletUpdate"
	PermPalletArchive                  Perm = "PalletArchive"
	PermPalletUnarchive                Perm = "PalletUnarchive"
	PermCartonList                     Perm = "CartonList"
	PermCartonRead                     Perm = "CartonRead"
	PermCartonCreate                   Perm = "CartonCreate"
	PermCartonUpdate                   Perm = "CartonUpdate"
	PermCartonArchive                  Perm = "CartonArchive"
	PermCartonUnarchive                Perm = "CartonUnarchive"
	PermProductList                    Perm = "ProductList"
	PermProductRead                    Perm = "ProductRead"
	PermProductCreate                  Perm = "ProductCreate"
	PermProductUpdate                  Perm = "ProductUpdate"
	PermProductArchive                 Perm = "ProductArchive"
	PermProductUnarchive               Perm = "ProductUnarchive"
	PermOrderList                      Perm = "OrderList"
	PermOrderRead                      Perm = "OrderRead"
	PermOrderCreate                    Perm = "OrderCreate"
	PermOrderUpdate                    Perm = "OrderUpdate"
	PermOrderArchive                   Perm = "OrderArchive"
	PermOrderUnarchive                 Perm = "OrderUnarchive"
	PermTrackActionList                Perm = "TrackActionList"
	PermTrackActionRead                Perm = "TrackActionRead"
	PermTrackActionCreate              Perm = "TrackActionCreate"
	PermTrackActionUpdate              Perm = "TrackActionUpdate"
	PermTrackActionArchive             Perm = "TrackActionArchive"
	PermTrackActionUnarchive           Perm = "TrackActionUnarchive"
	PermContractList                   Perm = "ContractList"
	PermContractRead                   Perm = "ContractRead"
	PermContractCreate                 Perm = "ContractCreate"
	PermContractUpdate                 Perm = "ContractUpdate"
	PermContractArchive                Perm = "ContractArchive"
	PermContractUnarchive              Perm = "ContractUnarchive"
	PermDistributorList                Perm = "DistributorList"
	PermDistributorRead                Perm = "DistributorRead"
	PermDistributorCreate              Perm = "DistributorCreate"
	PermDistributorUpdate              Perm = "DistributorUpdate"
	PermDistributorArchive             Perm = "DistributorArchive"
	PermDistributorUnarchive           Perm = "DistributorUnarchive"
	PermActivityListBlockchainActivity Perm = "ActivityListBlockchainActivity"
	PermActivityListUserActivity       Perm = "ActivityListUserActivity"
	PermUseAdvancedMode                Perm = "UseAdvancedMode"
	PermUseAdminPortal                 Perm = "UseAdminPortal"
)

var AllPerm = []Perm{
	PermUserList,
	PermUserCreate,
	PermUserRead,
	PermUserUpdate,
	PermUserArchive,
	PermUserUnarchive,
	PermReferralList,
	PermReferralRead,
	PermTaskList,
	PermTaskCreate,
	PermTaskRead,
	PermTaskUpdate,
	PermTaskArchive,
	PermTaskUnarchive,
	PermUserTaskList,
	PermUserTaskCreate,
	PermUserTaskRead,
	PermUserTaskUpdate,
	PermUserTaskArchive,
	PermUserTaskUnarchive,
	PermOrganisationList,
	PermOrganisationCreate,
	PermOrganisationRead,
	PermOrganisationUpdate,
	PermOrganisationArchive,
	PermOrganisationUnarchive,
	PermRoleList,
	PermRoleCreate,
	PermRoleRead,
	PermRoleUpdate,
	PermRoleArchive,
	PermRoleUnarchive,
	PermSKUList,
	PermSKUCreate,
	PermSKURead,
	PermSKUUpdate,
	PermSKUArchive,
	PermSKUUnarchive,
	PermCategoryList,
	PermCategoryCreate,
	PermCategoryRead,
	PermCategoryUpdate,
	PermCategoryArchive,
	PermCategoryUnarchive,
	PermProductCategoryList,
	PermProductCategoryCreate,
	PermProductCategoryRead,
	PermProductCategoryUpdate,
	PermProductCategoryArchive,
	PermProductCategoryUnarchive,
	PermContainerList,
	PermContainerRead,
	PermContainerCreate,
	PermContainerUpdate,
	PermContainerArchive,
	PermContainerUnarchive,
	PermPalletList,
	PermPalletRead,
	PermPalletCreate,
	PermPalletUpdate,
	PermPalletArchive,
	PermPalletUnarchive,
	PermCartonList,
	PermCartonRead,
	PermCartonCreate,
	PermCartonUpdate,
	PermCartonArchive,
	PermCartonUnarchive,
	PermProductList,
	PermProductRead,
	PermProductCreate,
	PermProductUpdate,
	PermProductArchive,
	PermProductUnarchive,
	PermOrderList,
	PermOrderRead,
	PermOrderCreate,
	PermOrderUpdate,
	PermOrderArchive,
	PermOrderUnarchive,
	PermTrackActionList,
	PermTrackActionRead,
	PermTrackActionCreate,
	PermTrackActionUpdate,
	PermTrackActionArchive,
	PermTrackActionUnarchive,
	PermContractList,
	PermContractRead,
	PermContractCreate,
	PermContractUpdate,
	PermContractArchive,
	PermContractUnarchive,
	PermDistributorList,
	PermDistributorRead,
	PermDistributorCreate,
	PermDistributorUpdate,
	PermDistributorArchive,
	PermDistributorUnarchive,
	PermActivityListBlockchainActivity,
	PermActivityListUserActivity,
	PermUseAdvancedMode,
	PermUseAdminPortal,
}

func (e Perm) IsValid() bool {
	switch e {
	case PermUserList, PermUserCreate, PermUserRead, PermUserUpdate, PermUserArchive, PermUserUnarchive, PermReferralList, PermReferralRead, PermTaskList, PermTaskCreate, PermTaskRead, PermTaskUpdate, PermTaskArchive, PermTaskUnarchive, PermUserTaskList, PermUserTaskCreate, PermUserTaskRead, PermUserTaskUpdate, PermUserTaskArchive, PermUserTaskUnarchive, PermOrganisationList, PermOrganisationCreate, PermOrganisationRead, PermOrganisationUpdate, PermOrganisationArchive, PermOrganisationUnarchive, PermRoleList, PermRoleCreate, PermRoleRead, PermRoleUpdate, PermRoleArchive, PermRoleUnarchive, PermSKUList, PermSKUCreate, PermSKURead, PermSKUUpdate, PermSKUArchive, PermSKUUnarchive, PermCategoryList, PermCategoryCreate, PermCategoryRead, PermCategoryUpdate, PermCategoryArchive, PermCategoryUnarchive, PermProductCategoryList, PermProductCategoryCreate, PermProductCategoryRead, PermProductCategoryUpdate, PermProductCategoryArchive, PermProductCategoryUnarchive, PermContainerList, PermContainerRead, PermContainerCreate, PermContainerUpdate, PermContainerArchive, PermContainerUnarchive, PermPalletList, PermPalletRead, PermPalletCreate, PermPalletUpdate, PermPalletArchive, PermPalletUnarchive, PermCartonList, PermCartonRead, PermCartonCreate, PermCartonUpdate, PermCartonArchive, PermCartonUnarchive, PermProductList, PermProductRead, PermProductCreate, PermProductUpdate, PermProductArchive, PermProductUnarchive, PermOrderList, PermOrderRead, PermOrderCreate, PermOrderUpdate, PermOrderArchive, PermOrderUnarchive, PermTrackActionList, PermTrackActionRead, PermTrackActionCreate, PermTrackActionUpdate, PermTrackActionArchive, PermTrackActionUnarchive, PermContractList, PermContractRead, PermContractCreate, PermContractUpdate, PermContractArchive, PermContractUnarchive, PermDistributorList, PermDistributorRead, PermDistributorCreate, PermDistributorUpdate, PermDistributorArchive, PermDistributorUnarchive, PermActivityListBlockchainActivity, PermActivityListUserActivity, PermUseAdvancedMode, PermUseAdminPortal:
		return true
	}
	return false
}

func (e Perm) String() string {
	return string(e)
}

func (e *Perm) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Perm(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Perm", str)
	}
	return nil
}

func (e Perm) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortByOption string

const (
	SortByOptionDateCreated  SortByOption = "DateCreated"
	SortByOptionDateUpdated  SortByOption = "DateUpdated"
	SortByOptionAlphabetical SortByOption = "Alphabetical"
)

var AllSortByOption = []SortByOption{
	SortByOptionDateCreated,
	SortByOptionDateUpdated,
	SortByOptionAlphabetical,
}

func (e SortByOption) IsValid() bool {
	switch e {
	case SortByOptionDateCreated, SortByOptionDateUpdated, SortByOptionAlphabetical:
		return true
	}
	return false
}

func (e SortByOption) String() string {
	return string(e)
}

func (e *SortByOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortByOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortByOption", str)
	}
	return nil
}

func (e SortByOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDir string

const (
	SortDirAscending  SortDir = "Ascending"
	SortDirDescending SortDir = "Descending"
)

var AllSortDir = []SortDir{
	SortDirAscending,
	SortDirDescending,
}

func (e SortDir) IsValid() bool {
	switch e {
	case SortDirAscending, SortDirDescending:
		return true
	}
	return false
}

func (e SortDir) String() string {
	return string(e)
}

func (e *SortDir) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDir(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDir", str)
	}
	return nil
}

func (e SortDir) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
