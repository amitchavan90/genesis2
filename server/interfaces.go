package genesis

import (
	"context"
	"database/sql"
	"genesis/db"
	"genesis/email"
	"genesis/graphql"
	"genesis/sms"
	"genesis/store"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"
)

var _ UserStorer = &store.Users{}
var _ ProspectStorer = &store.Prospects{}
var _ Notifier = &email.Console{}
var _ Notifier = &email.Mailer{}
var _ TokenStorer = &store.Tokens{}
var _ AuthProvider = &Auther{}
var _ Messenger = &sms.Mock{}
var _ Messenger = &sms.Client{}

// TODO change all *sql.Tx to *sqlx.Tx
// TODO remove handleTransactions() as it can only handle transaction in same store, not cross stores. Use prepConn() instead

// TokenStorer collects all token methods
type TokenStorer interface {
	Get(id string) (*db.IssuedToken, error)
	Insert(t *db.IssuedToken) (*db.IssuedToken, error)
	Update(t *db.IssuedToken) (*db.IssuedToken, error)
	GetAllByUser(userID string) ([]*db.IssuedToken, error)
	GetAllExpired() ([]*db.IssuedToken, error)
	Delete(t *db.IssuedToken) error
	Blacklist() (store.Blacklist, error)
}

// Notifier is used for notifying the user of things
type Notifier interface {
	ReceivedSignup(email string) error
	ForgotPassword(email string) error
}

// ProspectStorer collects all prospect methods
type ProspectStorer interface {
	Get(id uuid.UUID, tx ...*sql.Tx) (*db.Prospect, error)
	Start(email string, tx ...*sql.Tx) (*db.Prospect, error)
	Update(p *db.Prospect, tx ...*sql.Tx) (*db.Prospect, error)
	Finish(id uuid.UUID, tx ...*sql.Tx) (*db.Prospect, error)
}

// OrganisationStorer collects all todo methods
type OrganisationStorer interface {
	All() (db.OrganisationSlice, error)
	Get(id uuid.UUID) (*db.Organisation, error)
	GetMany(keys []string) (db.OrganisationSlice, []error)
	Insert(record *db.Organisation, txes ...*sql.Tx) (*db.Organisation, error)
	Update(record *db.Organisation, txes ...*sql.Tx) (*db.Organisation, error)
}

// BlobStorer collects all blob methods
type BlobStorer interface {
	Get(id uuid.UUID) (*db.Blob, error)
	GetMany(keys []string) (db.BlobSlice, []error)
	All() (db.BlobSlice, error)
	Insert(record *db.Blob, txes ...*sql.Tx) (*db.Blob, error)
	Update(t *db.Blob, txes ...*sql.Tx) (*db.Blob, error)
	Delete(blob *db.Blob, tx ...*sql.Tx) error
	Exists(id string) (bool, error)
}

// SKUStorer collects all sku methods
type SKUStorer interface {
	GetByCode(code string) (*db.StockKeepingUnit, error)
	Get(id uuid.UUID) (*db.StockKeepingUnit, error)
	GetMany(keys []string) (db.StockKeepingUnitSlice, []error)
	GetCategories(skuID string, txes ...*sql.Tx) (db.CategorySlice, error)
	GetProductCategories(skuID string, txes ...*sql.Tx) (db.ProductCategorySlice, error)
	All() (db.StockKeepingUnitSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, db.StockKeepingUnitSlice, error)
	Insert(u *db.StockKeepingUnit, tx ...*sql.Tx) (*db.StockKeepingUnit, error)
	InsertCategory(cat *db.Category, txes ...*sql.Tx) (*db.Category, error)
	InsertProductCategory(pcat *db.ProductCategory, txes ...*sql.Tx) (*db.ProductCategory, error)
	Update(u *db.StockKeepingUnit, tx ...*sql.Tx) (*db.StockKeepingUnit, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.StockKeepingUnit, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.StockKeepingUnit, error)
	Count() (int64, error)
	ProductCount(record *db.StockKeepingUnit) (int64, error)

	HasClones(record *db.StockKeepingUnit) (bool, error)
	GetClones(id string) (db.StockKeepingUnitSlice, error)

	GetContent(sku *db.StockKeepingUnit, contentType string) (db.StockKeepingUnitContentSlice, error)
	UpdateContent(sku *db.StockKeepingUnit, input []*graphql.SKUContentInput, contentType string) error
	GetPhotos(sku *db.StockKeepingUnit) (db.StockKeepingUnitPhotoSlice, error)
	UpdatePhotos(sku *db.StockKeepingUnit, blobIDs []string) error
}

// OrderStorer collects all order methods
type OrderStorer interface {
	GetByCode(code string) (*db.Order, error)
	Get(id uuid.UUID) (*db.Order, error)
	GetMany(keys []string) (db.OrderSlice, []error)
	All() (db.OrderSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Order, error)
	Insert(u *db.Order, tx ...*sqlx.Tx) (*db.Order, error)
	Update(u *db.Order, tx ...*sql.Tx) (*db.Order, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Order, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Order, error)
	Count() (int64, error)
	ProductCount(record *db.Order) (int64, error)
	SkuID(record *db.Order) (*string, error)

	Products(record *db.Order) (db.ProductSlice, error)
}

// ContainerStorer collects all container methods
type ContainerStorer interface {
	GetByCode(code string) (*db.Container, error)
	Get(id uuid.UUID) (*db.Container, error)
	GetMany(keys []string) (db.ContainerSlice, []error)
	GetRange(from string, to string) (db.ContainerSlice, error)
	All() (db.ContainerSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Container, error)
	Insert(u *db.Container, tx ...*sql.Tx) (*db.Container, error)
	Update(u *db.Container, tx ...*sql.Tx) (*db.Container, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Container, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Container, error)
	Count() (int64, error)
	PalletCount(record *db.Container) (int64, error)
}

// PalletStorer collects all pallet methods
type PalletStorer interface {
	GetByCode(code string) (*db.Pallet, error)
	Get(id uuid.UUID) (*db.Pallet, error)
	GetMany(keys []string) (db.PalletSlice, []error)
	GetRange(from string, to string) (db.PalletSlice, error)
	All() (db.PalletSlice, error)
	SearchSelect(
		search graphql.SearchFilter,
		limit int,
		offset int,
		containerID null.String,
		trackActionID null.String,
	) (int64, []*db.Pallet, error)
	Insert(u *db.Pallet, tx ...*sql.Tx) (*db.Pallet, error)
	Update(u *db.Pallet, tx ...*sql.Tx) (*db.Pallet, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Pallet, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Pallet, error)
	Count() (int64, error)
	CartonCount(record *db.Pallet) (int64, error)
	LatestTrackAction(record *db.Pallet) (*graphql.LatestTransactionInfo, error)
}

// CartonStorer collects all carton methods
type CartonStorer interface {
	Get(id uuid.UUID) (*db.Carton, error)
	GetByCode(code string) (*db.Carton, error)
	GetMany(keys []string) (db.CartonSlice, []error)
	GetManyByPalletID(palletID uuid.UUID) (db.CartonSlice, error)
	GetManyByContainerID(containerID uuid.UUID) (db.CartonSlice, error)
	GetRange(from string, to string) (db.CartonSlice, error)
	All() (db.CartonSlice, error)
	SearchSelect(
		search graphql.SearchFilter,
		limit int,
		offset int,
		containerID null.String,
		trackActionID null.String,
	) (int64, []*db.Carton, error)
	Insert(u *db.Carton, tx ...*sql.Tx) (*db.Carton, error)
	Update(u *db.Carton, tx ...*sql.Tx) (*db.Carton, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Carton, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Carton, error)
	Count() (int64, error)
	ProductCount(record *db.Carton) (int64, error)
	LatestTrackAction(record *db.Carton) (*graphql.LatestTransactionInfo, error)
	OrderID(record *db.Carton) (*string, error)
	DistributorID(record *db.Carton) (*string, error)
	SkuID(record *db.Carton) (*string, error)
}

// ProductStorer collects all product methods
type ProductStorer interface {
	GetByCode(code string, tx ...*sqlx.Tx) (*db.Product, error)
	GetByRegisterID(id uuid.UUID, tx ...*sqlx.Tx) (*db.Product, error)
	Get(id uuid.UUID, tx ...*sqlx.Tx) (*db.Product, error)
	GetMany(keys []string) (db.ProductSlice, []error)
	GetManyByCartonID(cartonID uuid.UUID, tx ...*sqlx.Tx) (db.ProductSlice, error)
	GetManyByPalletID(palletID uuid.UUID, tx ...*sqlx.Tx) (db.ProductSlice, error)
	GetManyByContainerID(containerID uuid.UUID, tx ...*sqlx.Tx) (db.ProductSlice, error)
	All(tx ...*sqlx.Tx) (db.ProductSlice, error)
	SearchSelect(
		search graphql.SearchFilter,
		limit int,
		offset int,
		cartonID null.String,
		orderID null.String,
		skuID null.String,
		distributorID null.String,
		contractID null.String,
		trackActionID null.String,
	) (int64, []*db.Product, error)
	Insert(u *db.Product, tx ...*sqlx.Tx) (*db.Product, error)
	Update(u *db.Product, tx ...*sqlx.Tx) (*db.Product, error)
	Archive(id uuid.UUID, tx ...*sqlx.Tx) (*db.Product, error)
	Unarchive(id uuid.UUID, tx ...*sqlx.Tx) (*db.Product, error)
	Count(tx ...*sqlx.Tx) (int64, error)
	Registered(record *db.Product, tx ...*sqlx.Tx) (bool, error)
	RegisteredBy(record *db.Product, tx ...*sqlx.Tx) (*db.User, error)
	LatestTrackAction(record *db.Product, tx ...*sqlx.Tx) (*graphql.LatestTransactionInfo, error)
}

// UserStorer collects all user methods
type UserStorer interface {
	BeginTransaction() (*sql.Tx, error)
	GetByVerifyToken(token string, txes ...*sql.Tx) (*db.User, error)
	GetByResetToken(token string, txes ...*sql.Tx) (*db.User, error)
	GetByEmail(email string, txes ...*sql.Tx) (*db.User, error)
	GetByReferralCode(referralCode string, txes ...*sql.Tx) (*db.User, error)
	GetReferrals(refByID string, txes ...*sql.Tx) (db.ReferralSlice, error)
	GetByWechatID(wechatID string, txes ...*sql.Tx) (*db.User, error)
	Get(id uuid.UUID, txes ...*sql.Tx) (*db.User, error)
	GetWithPermissions(id uuid.UUID, txes ...*sql.Tx) (*db.User, error)
	GetMany(keys []string, txes ...*sql.Tx) (db.UserSlice, []error)
	All(txes ...*sql.Tx) (db.UserSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int, consumer bool) (int64, []*db.User, error)
	GetByOrganisation(orgID uuid.UUID, txes ...*sql.Tx) (db.UserSlice, error)
	Insert(u *db.User, tx ...*sql.Tx) (*db.User, error)
	InsertUserAndReferral(u *db.User, referredByID string, tx ...*sql.Tx) (*db.User, error)
	Update(u *db.User, tx ...*sql.Tx) (*db.User, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.User, error)
}

// ReferralStorer collects all refereal methods
type ReferralStorer interface {
	BeginTransaction() (*sql.Tx, error)
	// GetByReferee(email string, txes ...*sql.Tx) (*db.Referral, error)
	// GetByReferredByID(refByID null.String, txes ...*sql.Tx) (*db.Referral, error)
	GetByUserID(userID string, txes ...*sql.Tx) (*db.Referral, error)
	Get(id uuid.UUID, txes ...*sql.Tx) (*db.Referral, error)
	GetMany(keys []string, txes ...*sql.Tx) (db.ReferralSlice, []error)
	All(txes ...*sql.Tx) (db.ReferralSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Referral, error)
	Insert(u *db.Referral, tx ...*sql.Tx) (*db.Referral, error)
	Update(u *db.Referral, tx ...*sql.Tx) (*db.Referral, error)
}

// RoleStorer collects all role methods
type RoleStorer interface {
	Get(id uuid.UUID) (*db.Role, error)
	GetWithTrackActions(id uuid.UUID) (*db.Role, error)
	GetMany(keys []string) (db.RoleSlice, []error)
	GetByName(name string) (*db.Role, error)
	GetByUser(userID uuid.UUID) (*db.Role, error)
	All() (db.RoleSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int, excludeSuper bool) (int64, []*db.Role, error)
	Insert(record *db.Role, tx ...*sql.Tx) (*db.Role, error)
	Update(record *db.Role, tx ...*sql.Tx) (*db.Role, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Role, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Role, error)
	GetTrackActions(record *db.Role) (db.TrackActionSlice, error)
	SetTrackActions(record *db.Role, actions db.TrackActionSlice) error
}

// TaskStorer collects all role methods
type TaskStorer interface {
	Get(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error)
	GetMany(keys []string, txes ...*sql.Tx) (db.TaskSlice, []error)
	GetSubtasks(taskID string, txes ...*sql.Tx) (db.SubtaskSlice, error)
	All(txes ...*sql.Tx) (db.TaskSlice, error)
	Count() (int64, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Task, error)
	Insert(record *db.Task, subtasks []db.Subtask, tx ...*sql.Tx) (*db.Task, error)
	Update(record *db.Task, tx ...*sql.Tx) (*db.Task, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Task, error)
}

// LoyaltyStorer collects all loyalty methods
type LoyaltyStorer interface {
	GetByProductID(id uuid.UUID) (*db.UserLoyaltyActivity, error)
	Insert(record *db.UserLoyaltyActivity, txes ...*sql.Tx) (*db.UserLoyaltyActivity, error)
	UserActivity(userID uuid.UUID) (db.UserLoyaltyActivitySlice, error)
	TotalPoints(ctx context.Context, userID uuid.UUID) (int, error)
}

// Messenger collects all SMS methods
type Messenger interface {
	SendMessage(mobileNumber string, content string) error
	SendToken(mobileNumber string, token string) error
}

// TrackActionStorer collects all TrackAction methods
type TrackActionStorer interface {
	Get(id uuid.UUID) (*db.TrackAction, error)
	GetMany(keys []string) (db.TrackActionSlice, []error)
	GetByCode(code string) (*db.TrackAction, error)
	All() (db.TrackActionSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.TrackAction, error)
	Insert(u *db.TrackAction, tx ...*sql.Tx) (*db.TrackAction, error)
	Update(u *db.TrackAction, tx ...*sql.Tx) (*db.TrackAction, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.TrackAction, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.TrackAction, error)
	Count() (int64, error)
}

// ContractStorer collects all Contract methods
type ContractStorer interface {
	Get(id uuid.UUID) (*db.Contract, error)
	GetMany(keys []string) (db.ContractSlice, []error)
	GetByCode(code string) (*db.Contract, error)
	All() (db.ContractSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Contract, error)
	Insert(u *db.Contract, tx ...*sql.Tx) (*db.Contract, error)
	Update(u *db.Contract, tx ...*sql.Tx) (*db.Contract, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Contract, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Contract, error)
	Count() (int64, error)
}

// DistributorStorer collects all Distributor methods
type DistributorStorer interface {
	Get(id uuid.UUID) (*db.Distributor, error)
	GetMany(keys []string) (db.DistributorSlice, []error)
	GetByCode(code string) (*db.Distributor, error)
	All() (db.DistributorSlice, error)
	SearchSelect(search graphql.SearchFilter, limit int, offset int) (int64, []*db.Distributor, error)
	Insert(u *db.Distributor, tx ...*sql.Tx) (*db.Distributor, error)
	Update(u *db.Distributor, tx ...*sql.Tx) (*db.Distributor, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Distributor, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Distributor, error)
	Count() (int64, error)
}

// TransactionStorer collects all Transaction methods
type TransactionStorer interface {
	Get(id uuid.UUID) (*db.Transaction, error)
	GetMany(keys []string) (db.TransactionSlice, []error)
	GetByCartonID(id uuid.UUID) (db.TransactionSlice, error)
	GetByProductID(id uuid.UUID) (db.TransactionSlice, error)

	InsertByCarton(
		carton *db.Carton,
		trackAction *db.TrackAction,
		user *db.User,
		createdByName string,
		scannedAt time.Time,
		cartonPhotoBlobID null.String,
		productPhotoBlobID null.String,
		optTransaction *db.Transaction,
		txes ...*sqlx.Tx,
	) (*db.Transaction, error)
	InsertByProduct(
		product *db.Product,
		trackAction *db.TrackAction,
		user *db.User,
		createdByName string,
		optTransaction *db.Transaction,
		txes ...*sqlx.Tx,
	) (*db.Transaction, error)
	AttachManyToProduct(transactions db.TransactionSlice, product *db.Product, tx *sqlx.Tx) error

	All() (db.TransactionSlice, error)
	SearchSelect(
		search graphql.SearchFilter,
		limit int,
		offset int,
		productID null.String,
		cartonID null.String,
		trackActionID null.String,
	) (int64, []*db.Transaction, error)
	Insert(u *db.Transaction, tx ...*sqlx.Tx) (*db.Transaction, error)
	Update(u *db.Transaction, tx ...*sql.Tx) (*db.Transaction, error)
	Archive(id uuid.UUID, txes ...*sql.Tx) (*db.Transaction, error)
	Unarchive(id uuid.UUID, txes ...*sql.Tx) (*db.Transaction, error)

	AllPending() (db.TransactionSlice, error)
	AllPendingCount() (int64, error)

	// Settings store methods (too lazy to make new store atm)
	GetSettings() (*db.Setting, error)
	UpdateSettings(record *db.Setting) (*db.Setting, error)
}

// UserActivityStorer collects all UserActivity methods
type UserActivityStorer interface {
	Get(id uuid.UUID) (*db.UserActivity, error)
	GetMany(keys []string) (db.UserActivitySlice, []error)
	All() (db.UserActivitySlice, error)
	SearchSelect(
		search graphql.SearchFilter,
		limit int,
		offset int,
		userID null.String,
	) (int64, db.UserActivitySlice, error)
	Insert(u *db.UserActivity, tx ...*sql.Tx) (*db.UserActivity, error)
}

// ManifestStorer collects all ManifestStorer methods
type ManifestStorer interface {
	BeginTransaction() (*sqlx.Tx, error)
	Get(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error)
	GetByLineHash(lineHash string, txes ...*sqlx.Tx) (*db.Manifest, error)
	GetByBlockchainByTransactionHash(bcTxHash string, txes ...*sqlx.Tx) (*db.Manifest, error)
	GetByBlockchainByMerkleRootHash(merkleRootHash string, txes ...*sqlx.Tx) (*db.Manifest, error)
	GetMany(keys []string, txes ...*sqlx.Tx) (db.ManifestSlice, []error)
	All(txes ...*sqlx.Tx) (db.ManifestSlice, error)

	Insert(u *db.Manifest, txes ...*sqlx.Tx) (*db.Manifest, error)
	Update(u *db.Manifest, txes ...*sqlx.Tx) (*db.Manifest, error)
	Archive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error)
	Unarchive(id uuid.UUID, txes ...*sqlx.Tx) (*db.Manifest, error)

	AllUnarchived(txes ...*sqlx.Tx) (db.ManifestSlice, error)
	AllPending(txes ...*sqlx.Tx) (db.ManifestSlice, error)
	AllConfirmed(txes ...*sqlx.Tx) (db.ManifestSlice, error)
	AllUnconfirmed(txes ...*sqlx.Tx) (db.ManifestSlice, error)
	AllUnfinished(txes ...*sqlx.Tx) (db.ManifestSlice, error)
}
