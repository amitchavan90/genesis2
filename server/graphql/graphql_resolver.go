package graphql

import (
	"context"
	"genesis/db"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Blob() BlobResolver {
	return &blobResolver{r}
}
func (r *Resolver) Carton() CartonResolver {
	return &cartonResolver{r}
}
func (r *Resolver) Container() ContainerResolver {
	return &containerResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Order() OrderResolver {
	return &orderResolver{r}
}
func (r *Resolver) Organisation() OrganisationResolver {
	return &organisationResolver{r}
}
func (r *Resolver) Pallet() PalletResolver {
	return &palletResolver{r}
}
func (r *Resolver) Product() ProductResolver {
	return &productResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Referral() ReferralResolver {
	return &referralResolver{r}
}
func (r *Resolver) Role() RoleResolver {
	return &roleResolver{r}
}
func (r *Resolver) SKU() SKUResolver {
	return &skuResolver{r}
}
func (r *Resolver) Settings() SettingsResolver {
	return &settingsResolver{r}
}
func (r *Resolver) Task() TaskResolver {
	return &taskResolver{r}
}
func (r *Resolver) UserTask() UserTaskResolver {
	return &userTaskResolver{r}
}
func (r *Resolver) TrackAction() TrackActionResolver {
	return &trackActionResolver{r}
}
func (r *Resolver) Transaction() TransactionResolver {
	return &transactionResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (r *Resolver) UserActivity() UserActivityResolver {
	return &userActivityResolver{r}
}
func (r *Resolver) UserLoyaltyActivity() UserLoyaltyActivityResolver {
	return &userLoyaltyActivityResolver{r}
}

type blobResolver struct{ *Resolver }

func (r *blobResolver) FileURL(ctx context.Context, obj *db.Blob) (*string, error) {
	panic("not implemented")
}

type cartonResolver struct{ *Resolver }

func (r *cartonResolver) Pallet(ctx context.Context, obj *db.Carton) (*db.Pallet, error) {
	panic("not implemented")
}
func (r *cartonResolver) Sku(ctx context.Context, obj *db.Carton) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *cartonResolver) Order(ctx context.Context, obj *db.Carton) (*db.Order, error) {
	panic("not implemented")
}
func (r *cartonResolver) Distributor(ctx context.Context, obj *db.Carton) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *cartonResolver) Transactions(ctx context.Context, obj *db.Carton) ([]*db.Transaction, error) {
	panic("not implemented")
}
func (r *cartonResolver) LatestTrackAction(ctx context.Context, obj *db.Carton) (*LatestTransactionInfo, error) {
	panic("not implemented")
}
func (r *cartonResolver) ProductCount(ctx context.Context, obj *db.Carton) (int, error) {
	panic("not implemented")
}

type containerResolver struct{ *Resolver }

func (r *containerResolver) PalletCount(ctx context.Context, obj *db.Container) (int, error) {
	panic("not implemented")
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RequestToken(ctx context.Context, input *RequestToken) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*db.Blob, error) {
	panic("not implemented")
}
func (r *mutationResolver) FileUploadMultiple(ctx context.Context, files []*graphql.Upload) ([]*db.Blob, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeploySmartContract(ctx context.Context) (*db.Setting, error) {
	panic("not implemented")
}
func (r *mutationResolver) RoleCreate(ctx context.Context, input UpdateRole) (*db.Role, error) {
	panic("not implemented")
}
func (r *mutationResolver) RoleUpdate(ctx context.Context, id string, input UpdateRole) (*db.Role, error) {
	panic("not implemented")
}
func (r *mutationResolver) RoleArchive(ctx context.Context, id string) (*db.Role, error) {
	panic("not implemented")
}
func (r *mutationResolver) RoleUnarchive(ctx context.Context, id string) (*db.Role, error) {
	panic("not implemented")
}
func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, password string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ChangeDetails(ctx context.Context, input UpdateUser) (*db.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserCreate(ctx context.Context, input UpdateUser) (*db.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserUpdate(ctx context.Context, id string, input UpdateUser) (*db.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) ForgotPassword(ctx context.Context, email string, viaSms *bool) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResetPassword(ctx context.Context, token string, password string, email *null.String) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResendEmailVerification(ctx context.Context, email string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserArchive(ctx context.Context, id string) (*db.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserUnarchive(ctx context.Context, id string) (*db.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) TaskCreate(ctx context.Context, input UpdateTask) (*db.Task, error) {
	panic("not implemented")
}
func (r *mutationResolver) TaskUpdate(ctx context.Context, id string, input UpdateTask) (*db.Task, error) {
	panic("not implemented")
}
func (r *mutationResolver) TaskArchive(ctx context.Context, id string) (*db.Task, error) {
	panic("not implemented")
}
func (r *mutationResolver) TaskUnarchive(ctx context.Context, id string) (*db.Task, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserTaskCreate(ctx context.Context, input UpdateUserTask) (*db.UserTask, error) {
	panic("not implemented")
}
func (r *mutationResolver) UserTaskUpdate(ctx context.Context, id string, input UpdateUserTask) (*db.UserTask, error) {
	panic("not implemented")
}
func (r *mutationResolver) SkuCreate(ctx context.Context, input UpdateSku) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *mutationResolver) SkuUpdate(ctx context.Context, id string, input UpdateSku) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *mutationResolver) SkuArchive(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *mutationResolver) SkuUnarchive(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *mutationResolver) SkuBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) OrderCreate(ctx context.Context, input CreateOrder) (*db.Order, error) {
	panic("not implemented")
}
func (r *mutationResolver) OrderUpdate(ctx context.Context, id string, input UpdateOrder) (*db.Order, error) {
	panic("not implemented")
}
func (r *mutationResolver) OrderArchive(ctx context.Context, id string) (*db.Order, error) {
	panic("not implemented")
}
func (r *mutationResolver) OrderUnarchive(ctx context.Context, id string) (*db.Order, error) {
	panic("not implemented")
}
func (r *mutationResolver) OrderBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContainerCreate(ctx context.Context, input CreateContainer) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContainerUpdate(ctx context.Context, id string, input UpdateContainer) (*db.Container, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContainerArchive(ctx context.Context, id string) (*db.Container, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContainerUnarchive(ctx context.Context, id string) (*db.Container, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContainerBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) PalletCreate(ctx context.Context, input CreatePallet) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) PalletUpdate(ctx context.Context, id string, input UpdatePallet) (*db.Pallet, error) {
	panic("not implemented")
}
func (r *mutationResolver) PalletArchive(ctx context.Context, id string) (*db.Pallet, error) {
	panic("not implemented")
}
func (r *mutationResolver) PalletUnarchive(ctx context.Context, id string) (*db.Pallet, error) {
	panic("not implemented")
}
func (r *mutationResolver) PalletBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) CartonCreate(ctx context.Context, input CreateCarton) (string, error) {
	panic("not implemented")
}
func (r *mutationResolver) CartonUpdate(ctx context.Context, id string, input UpdateCarton) (*db.Carton, error) {
	panic("not implemented")
}
func (r *mutationResolver) CartonArchive(ctx context.Context, id string) (*db.Carton, error) {
	panic("not implemented")
}
func (r *mutationResolver) CartonUnarchive(ctx context.Context, id string) (*db.Carton, error) {
	panic("not implemented")
}
func (r *mutationResolver) CartonBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ProductCreate(ctx context.Context, input UpdateProduct) (*db.Product, error) {
	panic("not implemented")
}
func (r *mutationResolver) ProductUpdate(ctx context.Context, id string, input UpdateProduct) (*db.Product, error) {
	panic("not implemented")
}
func (r *mutationResolver) ProductArchive(ctx context.Context, id string) (*db.Product, error) {
	panic("not implemented")
}
func (r *mutationResolver) ProductUnarchive(ctx context.Context, id string) (*db.Product, error) {
	panic("not implemented")
}
func (r *mutationResolver) ProductBatchAction(ctx context.Context, ids []string, action Action, value *BatchActionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) TrackActionCreate(ctx context.Context, input UpdateTrackAction) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *mutationResolver) TrackActionUpdate(ctx context.Context, id string, input UpdateTrackAction) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *mutationResolver) TrackActionArchive(ctx context.Context, id string) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *mutationResolver) TrackActionUnarchive(ctx context.Context, id string) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *mutationResolver) RecordTransaction(ctx context.Context, input RecordTransactionInput) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) FlushPendingTransactions(ctx context.Context) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContractCreate(ctx context.Context, input UpdateContract) (*db.Contract, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContractUpdate(ctx context.Context, id string, input UpdateContract) (*db.Contract, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContractArchive(ctx context.Context, id string) (*db.Contract, error) {
	panic("not implemented")
}
func (r *mutationResolver) ContractUnarchive(ctx context.Context, id string) (*db.Contract, error) {
	panic("not implemented")
}
func (r *mutationResolver) DistributorCreate(ctx context.Context, input UpdateDistributor) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *mutationResolver) DistributorUpdate(ctx context.Context, id string, input UpdateDistributor) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *mutationResolver) DistributorArchive(ctx context.Context, id string) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *mutationResolver) DistributorUnarchive(ctx context.Context, id string) (*db.Distributor, error) {
	panic("not implemented")
}

type orderResolver struct{ *Resolver }

func (r *orderResolver) Sku(ctx context.Context, obj *db.Order) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *orderResolver) ProductCount(ctx context.Context, obj *db.Order) (int, error) {
	panic("not implemented")
}

type organisationResolver struct{ *Resolver }

func (r *organisationResolver) Users(ctx context.Context, obj *db.Organisation) ([]*db.User, error) {
	panic("not implemented")
}

type palletResolver struct{ *Resolver }

func (r *palletResolver) Container(ctx context.Context, obj *db.Pallet) (*db.Container, error) {
	panic("not implemented")
}
func (r *palletResolver) LatestTrackAction(ctx context.Context, obj *db.Pallet) (*LatestTransactionInfo, error) {
	panic("not implemented")
}
func (r *palletResolver) CartonCount(ctx context.Context, obj *db.Pallet) (int, error) {
	panic("not implemented")
}

type productResolver struct{ *Resolver }

func (r *productResolver) Sku(ctx context.Context, obj *db.Product) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *productResolver) Carton(ctx context.Context, obj *db.Product) (*db.Carton, error) {
	panic("not implemented")
}
func (r *productResolver) Order(ctx context.Context, obj *db.Product) (*db.Order, error) {
	panic("not implemented")
}
func (r *productResolver) Contract(ctx context.Context, obj *db.Product) (*db.Contract, error) {
	panic("not implemented")
}
func (r *productResolver) Distributor(ctx context.Context, obj *db.Product) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *productResolver) Registered(ctx context.Context, obj *db.Product) (bool, error) {
	panic("not implemented")
}
func (r *productResolver) RegisteredBy(ctx context.Context, obj *db.Product) (*db.User, error) {
	panic("not implemented")
}
func (r *productResolver) Transactions(ctx context.Context, obj *db.Product) ([]*db.Transaction, error) {
	panic("not implemented")
}
func (r *productResolver) LatestTrackAction(ctx context.Context, obj *db.Product) (*LatestTransactionInfo, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Settings(ctx context.Context) (*db.Setting, error) {
	panic("not implemented")
}
func (r *queryResolver) GetTickerInfo(ctx context.Context) (*TickerInfo, error) {
	panic("not implemented")
}
func (r *queryResolver) GetObject(ctx context.Context, id string) (*GetObjectResponse, error) {
	panic("not implemented")
}
func (r *queryResolver) GetObjects(ctx context.Context, input GetObjectsRequest) (*GetObjectsResponse, error) {
	panic("not implemented")
}
func (r *queryResolver) Referrals(ctx context.Context, search SearchFilter, limit int, offset int) (*ReferralsResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Referral(ctx context.Context, userID *string) (*db.Referral, error) {
	panic("not implemented")
}
func (r *queryResolver) Tasks(ctx context.Context, search SearchFilter, limit int, offset int) (*TasksResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Task(ctx context.Context, id *string) (*db.Task, error) {
	panic("not implemented")
}
func (r *queryResolver) UserTasks(ctx context.Context, search SearchFilter, limit int, offset int) (*UserTasksResult, error) {
	panic("not implemented")
}
func (r *queryResolver) UserTask(ctx context.Context, id *string) (*db.UserTask, error) {
	panic("not implemented")
}
func (r *queryResolver) Roles(ctx context.Context, search SearchFilter, limit int, offset int, excludeSuper bool) (*RolesResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Role(ctx context.Context, name string) (*db.Role, error) {
	panic("not implemented")
}
func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Organisations(ctx context.Context) ([]*db.Organisation, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context, search SearchFilter, limit int, offset int) (*UsersResult, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, email *string, wechatID *string) (*db.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Consumers(ctx context.Context, search SearchFilter, limit int, offset int) (*ConsumersResult, error) {
	panic("not implemented")
}
func (r *queryResolver) VerifyResetToken(ctx context.Context, token string, email *null.String) (bool, error) {
	panic("not implemented")
}
func (r *queryResolver) Skus(ctx context.Context, search SearchFilter, limit int, offset int) (*SKUResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Sku(ctx context.Context, code string) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *queryResolver) SkuByID(ctx context.Context, id string) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *queryResolver) SkuCloneTree(ctx context.Context, id string) ([]*SKUClone, error) {
	panic("not implemented")
}
func (r *queryResolver) Orders(ctx context.Context, search SearchFilter, limit int, offset int) (*OrderResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Order(ctx context.Context, code string) (*db.Order, error) {
	panic("not implemented")
}
func (r *queryResolver) Containers(ctx context.Context, search SearchFilter, limit int, offset int) (*ContainerResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Container(ctx context.Context, code string) (*db.Container, error) {
	panic("not implemented")
}
func (r *queryResolver) Pallets(ctx context.Context, search SearchFilter, limit int, offset int, containerID *string, trackActionID *string) (*PalletResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Pallet(ctx context.Context, code string) (*db.Pallet, error) {
	panic("not implemented")
}
func (r *queryResolver) Cartons(ctx context.Context, search SearchFilter, limit int, offset int, palletID *string, trackActionID *string) (*CartonResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Carton(ctx context.Context, code string) (*db.Carton, error) {
	panic("not implemented")
}
func (r *queryResolver) Products(ctx context.Context, search SearchFilter, limit int, offset int, cartonID *string, orderID *string, skuID *string, distributorID *string, contractID *string, trackActionID *string) (*ProductResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Product(ctx context.Context, code string) (*db.Product, error) {
	panic("not implemented")
}
func (r *queryResolver) ProductByID(ctx context.Context, id string) (*db.Product, error) {
	panic("not implemented")
}
func (r *queryResolver) GetLoyaltyActivity(ctx context.Context, userID string) ([]*db.UserLoyaltyActivity, error) {
	panic("not implemented")
}
func (r *queryResolver) TrackActions(ctx context.Context, search SearchFilter, limit int, offset int) (*TrackActionResult, error) {
	panic("not implemented")
}
func (r *queryResolver) TrackAction(ctx context.Context, code string) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *queryResolver) Transactions(ctx context.Context, search SearchFilter, limit int, offset int, productID *string, cartonID *string, trackActionID *string) (*TransactionsResult, error) {
	panic("not implemented")
}
func (r *queryResolver) PendingTransactionsCount(ctx context.Context) (int, error) {
	panic("not implemented")
}
func (r *queryResolver) EthereumAccountAddress(ctx context.Context) (string, error) {
	panic("not implemented")
}
func (r *queryResolver) EthereumAccountBalance(ctx context.Context) (string, error) {
	panic("not implemented")
}
func (r *queryResolver) Contracts(ctx context.Context, search SearchFilter, limit int, offset int) (*ContractResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Contract(ctx context.Context, code string) (*db.Contract, error) {
	panic("not implemented")
}
func (r *queryResolver) Distributors(ctx context.Context, search SearchFilter, limit int, offset int) (*DistributorResult, error) {
	panic("not implemented")
}
func (r *queryResolver) Distributor(ctx context.Context, code string) (*db.Distributor, error) {
	panic("not implemented")
}
func (r *queryResolver) UserActivities(ctx context.Context, search SearchFilter, limit int, offset int, userID *string) (*UserActivityResult, error) {
	panic("not implemented")
}

type referralResolver struct{ *Resolver }

func (r *referralResolver) ReferredByID(ctx context.Context, obj *db.Referral) (string, error) {
	panic("not implemented")
}

type roleResolver struct{ *Resolver }

func (r *roleResolver) Permissions(ctx context.Context, obj *db.Role) ([]Perm, error) {
	panic("not implemented")
}
func (r *roleResolver) TrackActions(ctx context.Context, obj *db.Role) ([]*db.TrackAction, error) {
	panic("not implemented")
}

type skuResolver struct{ *Resolver }

func (r *skuResolver) HasClones(ctx context.Context, obj *db.StockKeepingUnit) (bool, error) {
	panic("not implemented")
}
func (r *skuResolver) MasterPlan(ctx context.Context, obj *db.StockKeepingUnit) (*db.Blob, error) {
	panic("not implemented")
}
func (r *skuResolver) Video(ctx context.Context, obj *db.StockKeepingUnit) (*db.Blob, error) {
	panic("not implemented")
}
func (r *skuResolver) Urls(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.StockKeepingUnitContent, error) {
	panic("not implemented")
}
func (r *skuResolver) ProductInfo(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.StockKeepingUnitContent, error) {
	panic("not implemented")
}
func (r *skuResolver) Photos(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.Blob, error) {
	panic("not implemented")
}
func (r *skuResolver) Categories(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.Category, error) {
	panic("not implemented")
}
func (r *skuResolver) ProductCategories(ctx context.Context, obj *db.StockKeepingUnit) ([]*db.ProductCategory, error) {
	panic("not implemented")
}
func (r *skuResolver) ProductCount(ctx context.Context, obj *db.StockKeepingUnit) (int, error) {
	panic("not implemented")
}

type settingsResolver struct{ *Resolver }

func (r *settingsResolver) ConsumerHost(ctx context.Context, obj *db.Setting) (string, error) {
	panic("not implemented")
}
func (r *settingsResolver) AdminHost(ctx context.Context, obj *db.Setting) (string, error) {
	panic("not implemented")
}
func (r *settingsResolver) EtherscanHost(ctx context.Context, obj *db.Setting) (string, error) {
	panic("not implemented")
}
func (r *settingsResolver) FieldappVersion(ctx context.Context, obj *db.Setting) (string, error) {
	panic("not implemented")
}

type taskResolver struct{ *Resolver }

func (r *taskResolver) FinishDate(ctx context.Context, obj *db.Task) (*time.Time, error) {
	panic("not implemented")
}
func (r *taskResolver) Sku(ctx context.Context, obj *db.Task) (*db.StockKeepingUnit, error) {
	panic("not implemented")
}
func (r *taskResolver) Subtasks(ctx context.Context, obj *db.Task) ([]*db.Subtask, error) {
	panic("not implemented")
}

type userTaskResolver struct{ *Resolver }

func (r *userTaskResolver) Task(ctx context.Context, obj *db.UserTask) (*db.Task, error) {
	panic("not implemented")
}
func (r *userTaskResolver) User(ctx context.Context, obj *db.UserTask) (*db.User, error) {
	panic("not implemented")
}
func (r *userTaskResolver) UserSubtasks(ctx context.Context, obj *db.UserTask) ([]*db.UserSubtask, error) {
	panic("not implemented")
}

type trackActionResolver struct{ *Resolver }

func (r *trackActionResolver) RequirePhotos(ctx context.Context, obj *db.TrackAction) ([]bool, error) {
	panic("not implemented")
}

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) TransactionPending(ctx context.Context, obj *db.Transaction) (bool, error) {
	panic("not implemented")
}
func (r *transactionResolver) Manifest(ctx context.Context, obj *db.Transaction) (*db.Manifest, error) {
	panic("not implemented")
}
func (r *transactionResolver) Action(ctx context.Context, obj *db.Transaction) (*db.TrackAction, error) {
	panic("not implemented")
}
func (r *transactionResolver) CreatedBy(ctx context.Context, obj *db.Transaction) (*db.User, error) {
	panic("not implemented")
}
func (r *transactionResolver) Carton(ctx context.Context, obj *db.Transaction) (*db.Carton, error) {
	panic("not implemented")
}
func (r *transactionResolver) Product(ctx context.Context, obj *db.Transaction) (*db.Product, error) {
	panic("not implemented")
}
func (r *transactionResolver) Photos(ctx context.Context, obj *db.Transaction) (*TransactionPhotos, error) {
	panic("not implemented")
}

type userResolver struct{ *Resolver }

func (r *userResolver) Organisation(ctx context.Context, obj *db.User) (*db.Organisation, error) {
	panic("not implemented")
}
func (r *userResolver) Role(ctx context.Context, obj *db.User) (*db.Role, error) {
	panic("not implemented")
}
func (r *userResolver) LoyaltyPoints(ctx context.Context, obj *db.User) (int, error) {
	panic("not implemented")
}
func (r *userResolver) Referrals(ctx context.Context, obj *db.User) ([]*db.Referral, error) {
	panic("not implemented")
}

type userActivityResolver struct{ *Resolver }

func (r *userActivityResolver) User(ctx context.Context, obj *db.UserActivity) (*db.User, error) {
	panic("not implemented")
}
func (r *userActivityResolver) ObjectType(ctx context.Context, obj *db.UserActivity) (ObjectType, error) {
	panic("not implemented")
}

type userLoyaltyActivityResolver struct{ *Resolver }

func (r *userLoyaltyActivityResolver) User(ctx context.Context, obj *db.UserLoyaltyActivity) (*db.User, error) {
	panic("not implemented")
}
func (r *userLoyaltyActivityResolver) Product(ctx context.Context, obj *db.UserLoyaltyActivity) (*db.Product, error) {
	panic("not implemented")
}
