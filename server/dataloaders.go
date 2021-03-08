package genesis

import (
	"genesis/dataloaders"
	"genesis/db"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"

	"context"
	"net/http"
	"time"
)

// ContextKey holds a custom String func for uniqueness
type ContextKey string

func (k ContextKey) String() string {
	return "dataloader_" + string(k)
}

// UserLoaderKey declares a statically typed key for context reference in other packages
const UserLoaderKey ContextKey = "user_loader"

// ReferralLoaderKey declares a statically typed key for context reference in other packages
const ReferralLoaderKey ContextKey = "referral_loader"

// TaskLoaderKey declares a statically typed key for context reference in other packages
const TaskLoaderKey ContextKey = "task_loader"

// RoleLoaderKey declares a statically typed key for context reference in other packages
const RoleLoaderKey ContextKey = "role_loader"

// OrganisationLoaderKey declares a statically typed key for context reference in other packages
const OrganisationLoaderKey ContextKey = "organisation_loader"

// OrganisationUsersLoaderKey declares a statically typed key for context reference in other packages
const OrganisationUsersLoaderKey ContextKey = "organisation_users_loader"

// SKULoaderKey declares a statically typed key for context reference in other packages
const SKULoaderKey ContextKey = "sku_loader"

// OrderLoaderKey declares a statically typed key for context reference in other packages
const OrderLoaderKey ContextKey = "order_loader"

// DistributorLoaderKey declares a statically typed key for context reference in other packages
const DistributorLoaderKey ContextKey = "distributor_loader"

// ContainerLoaderKey declares a statically typed key for context reference in other packages
const ContainerLoaderKey ContextKey = "container_loader"

// PalletLoaderKey declares a statically typed key for context reference in other packages
const PalletLoaderKey ContextKey = "pallet_loader"

// CartonLoaderKey declares a statically typed key for context reference in other packages
const CartonLoaderKey ContextKey = "carton_loader"

// ProductLoaderKey declares a statically typed key for context reference in other packages
const ProductLoaderKey ContextKey = "product_loader"

// ContractLoaderKey declares a statically typed key for context reference in other packages
const ContractLoaderKey ContextKey = "contract_loader"

// TransactionLoaderKey declares a statically typed key for context reference in other packages
const TransactionLoaderKey ContextKey = "transaction_loader"

// ManifestLoaderKey declares a statically typed key for context reference in other packages
const ManifestLoaderKey ContextKey = "manifest_loader"

// TrackActionLoaderKey declares a statically typed key for context reference in other packages
const TrackActionLoaderKey ContextKey = "track_action_loader"

// UserLoaderFromContext runs the dataloader inside the context
func UserLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return ctx.Value(UserLoaderKey).(*dataloaders.UserLoader).Load(id.String())
}

// ReferralLoaderFromContext runs the dataloader inside the context
func ReferralLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Referral, error) {
	return ctx.Value(ReferralLoaderKey).(*dataloaders.ReferralLoader).Load(id.String())
}

// TaskLoaderFromContext runs the dataloader inside the context
func TaskLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Task, error) {
	return ctx.Value(TaskLoaderKey).(*dataloaders.TaskLoader).Load(id.String())
}

// RoleLoaderFromContext runs the dataloader inside the context
func RoleLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Role, error) {
	return ctx.Value(RoleLoaderKey).(*dataloaders.RoleLoader).Load(id.String())
}

// OrganisationUsersLoaderFromContext runs the dataloader inside the context
func OrganisationUsersLoaderFromContext(ctx context.Context, id uuid.UUID) ([]*db.User, error) {
	return ctx.Value(OrganisationUsersLoaderKey).(*dataloaders.OrganisationUsersLoader).Load(id.String())
}

// OrganisationLoaderFromContext runs the dataloader inside the context
func OrganisationLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Organisation, error) {
	return ctx.Value(OrganisationLoaderKey).(*dataloaders.OrganisationLoader).Load(id.String())
}

// SKULoaderFromContext runs the dataloader inside the context
func SKULoaderFromContext(ctx context.Context, id uuid.UUID) (*db.StockKeepingUnit, error) {
	return ctx.Value(SKULoaderKey).(*dataloaders.SKULoader).Load(id.String())
}

// OrderLoaderFromContext runs the dataloader inside the context
func OrderLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Order, error) {
	return ctx.Value(OrderLoaderKey).(*dataloaders.OrderLoader).Load(id.String())
}

// DistributorLoaderFromContext runs the dataloader inside the context
func DistributorLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Distributor, error) {
	return ctx.Value(DistributorLoaderKey).(*dataloaders.DistributorLoader).Load(id.String())
}

// ContainerLoaderFromContext runs the dataloader inside the context
func ContainerLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Container, error) {
	return ctx.Value(ContainerLoaderKey).(*dataloaders.ContainerLoader).Load(id.String())
}

// PalletLoaderFromContext runs the dataloader inside the context
func PalletLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Pallet, error) {
	return ctx.Value(PalletLoaderKey).(*dataloaders.PalletLoader).Load(id.String())
}

// CartonLoaderFromContext runs the dataloader inside the context
func CartonLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Carton, error) {
	return ctx.Value(CartonLoaderKey).(*dataloaders.CartonLoader).Load(id.String())
}

// ProductLoaderFromContext runs the dataloader inside the context
func ProductLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Product, error) {
	return ctx.Value(ProductLoaderKey).(*dataloaders.ProductLoader).Load(id.String())
}

// ContractLoaderFromContext runs the dataloader inside the context
func ContractLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Contract, error) {
	return ctx.Value(ContractLoaderKey).(*dataloaders.ContractLoader).Load(id.String())
}

// TransactionLoaderFromContext runs the dataloader inside the context
func TransactionLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	return ctx.Value(TransactionLoaderKey).(*dataloaders.TransactionLoader).Load(id.String())
}

// TrackActionLoaderFromContext runs the dataloader inside the context
func TrackActionLoaderFromContext(ctx context.Context, id uuid.UUID) (*db.TrackAction, error) {
	return ctx.Value(TrackActionLoaderKey).(*dataloaders.TrackActionLoader).Load(id.String())
}

// WithDataloaders returns a new context that contains dataloaders
func WithDataloaders(
	ctx context.Context,
	SKUStorer SKUStorer,
	OrderStorer OrderStorer,
	DistributorStorer DistributorStorer,
	ContainerStorer ContainerStorer,
	PalletStorer PalletStorer,
	CartonStorer CartonStorer,
	ProductStorer ProductStorer,
	OrganisationStorer OrganisationStorer,
	UserStorer UserStorer,
	referralStore ReferralStorer,
	taskStore TaskStorer,
	RoleStorer RoleStorer,
	ContractStorer ContractStorer,
	TransactionStorer TransactionStorer,
	ManifestStorer ManifestStorer,
	TrackActionStorer TrackActionStorer,
) context.Context {
	userloader := dataloaders.NewUserLoader(
		dataloaders.UserLoaderConfig{
			Fetch: func(ids []string) ([]*db.User, []error) {
				return UserStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	roleLoader := dataloaders.NewRoleLoader(
		dataloaders.RoleLoaderConfig{
			Fetch: func(ids []string) ([]*db.Role, []error) {
				return RoleStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	organisationLoader := dataloaders.NewOrganisationLoader(
		dataloaders.OrganisationLoaderConfig{
			Fetch: func(ids []string) ([]*db.Organisation, []error) {
				return OrganisationStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)
	organisationUsersLoader := dataloaders.NewOrganisationUsersLoader(
		dataloaders.OrganisationUsersLoaderConfig{
			Fetch: func(ids []string) ([][]*db.User, []error) {
				result := [][]*db.User{}
				for _, id := range ids {
					userID, err := uuid.FromString(id)
					if err != nil {
						return nil, []error{terror.ErrParse}
					}
					records, err := UserStorer.GetByOrganisation(userID)
					if err != nil {
						return nil, []error{terror.ErrDataloader}
					}
					result = append(result, records)
				}
				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	skuLoader := dataloaders.NewSKULoader(
		dataloaders.SKULoaderConfig{
			Fetch: func(ids []string) ([]*db.StockKeepingUnit, []error) {
				return SKUStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	orderLoader := dataloaders.NewOrderLoader(
		dataloaders.OrderLoaderConfig{
			Fetch: func(ids []string) ([]*db.Order, []error) {
				return OrderStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	distributorLoader := dataloaders.NewDistributorLoader(
		dataloaders.DistributorLoaderConfig{
			Fetch: func(ids []string) ([]*db.Distributor, []error) {
				return DistributorStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	containerLoader := dataloaders.NewContainerLoader(
		dataloaders.ContainerLoaderConfig{
			Fetch: func(ids []string) ([]*db.Container, []error) {
				return ContainerStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	palletLoader := dataloaders.NewPalletLoader(
		dataloaders.PalletLoaderConfig{
			Fetch: func(ids []string) ([]*db.Pallet, []error) {
				return PalletStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	cartonLoader := dataloaders.NewCartonLoader(
		dataloaders.CartonLoaderConfig{
			Fetch: func(ids []string) ([]*db.Carton, []error) {
				return CartonStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	productLoader := dataloaders.NewProductLoader(
		dataloaders.ProductLoaderConfig{
			Fetch: func(ids []string) ([]*db.Product, []error) {
				return ProductStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	contractLoader := dataloaders.NewContractLoader(
		dataloaders.ContractLoaderConfig{
			Fetch: func(ids []string) ([]*db.Contract, []error) {
				return ContractStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	TransactionLoader := dataloaders.NewTransactionLoader(
		dataloaders.TransactionLoaderConfig{
			Fetch: func(ids []string) ([]*db.Transaction, []error) {
				return TransactionStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	manifestLoader := dataloaders.NewManifestLoader(
		dataloaders.ManifestLoaderConfig{
			Fetch: func(ids []string) ([]*db.Manifest, []error) {
				return ManifestStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	trackActionLoader := dataloaders.NewTrackActionLoader(
		dataloaders.TrackActionLoaderConfig{
			Fetch: func(ids []string) ([]*db.TrackAction, []error) {
				return TrackActionStorer.GetMany(ids)
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	ctx = context.WithValue(ctx, UserLoaderKey, userloader)
	ctx = context.WithValue(ctx, RoleLoaderKey, roleLoader)
	ctx = context.WithValue(ctx, OrganisationLoaderKey, organisationLoader)
	ctx = context.WithValue(ctx, OrganisationUsersLoaderKey, organisationUsersLoader)
	ctx = context.WithValue(ctx, SKULoaderKey, skuLoader)
	ctx = context.WithValue(ctx, OrderLoaderKey, orderLoader)
	ctx = context.WithValue(ctx, DistributorLoaderKey, distributorLoader)
	ctx = context.WithValue(ctx, ContainerLoaderKey, containerLoader)
	ctx = context.WithValue(ctx, PalletLoaderKey, palletLoader)
	ctx = context.WithValue(ctx, CartonLoaderKey, cartonLoader)
	ctx = context.WithValue(ctx, ProductLoaderKey, productLoader)
	ctx = context.WithValue(ctx, ContractLoaderKey, contractLoader)
	ctx = context.WithValue(ctx, TransactionLoaderKey, TransactionLoader)
	ctx = context.WithValue(ctx, ManifestLoaderKey, manifestLoader)
	ctx = context.WithValue(ctx, TrackActionLoaderKey, trackActionLoader)
	return ctx
}

// DataloaderMiddleware runs before each API call and loads the dataloaders into context
func DataloaderMiddleware(
	skuStorer SKUStorer,
	orderStorer OrderStorer,
	distributorStorer DistributorStorer,
	containerStorer ContainerStorer,
	palletStorer PalletStorer,
	cartonStorer CartonStorer,
	productStorer ProductStorer,
	organisationStore OrganisationStorer,
	userStore UserStorer,
	referralStore ReferralStorer,
	taskStore TaskStorer,
	roleStore RoleStorer,
	contractStore ContractStorer,
	transactionStorer TransactionStorer,
	manifestStorer ManifestStorer,
	trackActionStore TrackActionStorer,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithDataloaders(
				r.Context(),
				skuStorer,
				orderStorer,
				distributorStorer,
				containerStorer,
				palletStorer,
				cartonStorer,
				productStorer,
				organisationStore,
				userStore,
				referralStore,
				taskStore,
				roleStore,
				contractStore,
				transactionStorer,
				manifestStorer,
				trackActionStore,
			))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
