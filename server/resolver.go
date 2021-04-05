package genesis

import (
	"context"
	"genesis/blockchain"
	"genesis/canlog"
	"genesis/config"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
	"github.com/mailgun/mailgun-go/v3"
)

// Resolver for graphql
type Resolver struct {
	Auther                    AuthProvider
	Config                    *config.PlatformConfig
	Mailer                    *mailgun.MailgunImpl
	SKUStore                  SKUStorer
	OrderStore                OrderStorer
	ContainerStore            ContainerStorer
	PalletStore               PalletStorer
	CartonStore               CartonStorer
	ProductStore              ProductStorer
	OrganisationStore         OrganisationStorer
	UserStore                 UserStorer
	ReferralStore             ReferralStorer
	TaskStore                 TaskStorer
	UserTaskStore             UserTaskStorer
	UserSubtaskStore          UserSubtaskStorer
	RoleStore                 RoleStorer
	BlobStore                 BlobStorer
	ContractStore             ContractStorer
	DistributorStore          DistributorStorer
	TransactionStore          TransactionStorer
	ManifestStore             ManifestStorer
	TrackActionStore          TrackActionStorer
	UserActivityStore         UserActivityStorer
	UserPurchaseActivityStore UserPurchaseActivityStorer
	WalletTransactionStore    WalletTransactionStorer
	LoyaltyStore              LoyaltyStorer
	Blacklister               BlacklistProvider
	SmsMessenger              Messenger
	Blk                       *blockchain.Service
	SystemTicker              *SystemTicker
}

/////////////////
//   Queries   //
/////////////////

// Query resolver
func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	u, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return nil, terror.New(err, "user from context")
	}
	return u, nil
}

func (r *queryResolver) GetObject(ctx context.Context, id string) (*graphql.GetObjectResponse, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	response := &graphql.GetObjectResponse{}

	product, err := r.ProductStore.Get(uuid)
	if err == nil {
		response.Product = product
		return response, nil
	}

	carton, err := r.CartonStore.Get(uuid)
	if err == nil {
		response.Carton = carton
		return response, nil
	}

	pallet, err := r.PalletStore.Get(uuid)
	if err == nil {
		response.Pallet = pallet
		return response, nil
	}

	container, err := r.ContainerStore.Get(uuid)
	if err == nil {
		response.Container = container
		return response, nil
	}

	return response, nil
}

func (r *queryResolver) GetObjects(ctx context.Context, input graphql.GetObjectsRequest) (*graphql.GetObjectsResponse, error) {
	response := &graphql.GetObjectsResponse{
		Products:   []*db.Product{},
		Cartons:    []*db.Carton{},
		Pallets:    []*db.Pallet{},
		Containers: []*db.Container{},
	}

	if len(input.ProductIDs) > 0 {
		products, err := r.ProductStore.GetMany(input.ProductIDs)
		if err != nil {
			return nil, terror.New(err[0], "get objects")
		}
		response.Products = products
	}

	if len(input.CartonIDs) > 0 {
		cartons, err := r.CartonStore.GetMany(input.CartonIDs)
		if err != nil {
			return nil, terror.New(err[0], "get objects")
		}
		response.Cartons = cartons
	}

	if len(input.PalletIDs) > 0 {
		pallets, err := r.PalletStore.GetMany(input.PalletIDs)
		if err != nil {
			return nil, terror.New(err[0], "get objects")
		}
		response.Pallets = pallets
	}

	if len(input.ContainerIDs) > 0 {
		containers, err := r.ContainerStore.GetMany(input.ContainerIDs)
		if err != nil {
			return nil, terror.New(err[0], "get objects")
		}
		response.Containers = containers
	}

	return response, nil
}

///////////////
// Mutations //
///////////////

// Mutation resolver
func (r *Resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RequestToken(ctx context.Context, input *graphql.RequestToken) (string, error) {
	canlog.Set(ctx, "email", input.Email)
	err := r.Auther.ValidatePassword(ctx, input.Email, input.Password)
	if err != nil {
		return "", terror.New(terror.ErrBadCredentials, "")
	}
	u, err := r.UserStore.GetByEmail(input.Email)
	if err != nil {
		return "", terror.New(terror.ErrBadCredentials, "")
	}
	token, err := r.Auther.GenerateJWT(ctx, u.Email.String, u.ID, u.RoleID, "", nil)
	if err != nil {
		return "", terror.New(terror.ErrBadCredentials, "")
	}
	return token, nil
}

func (r *mutationResolver) DeploySmartContract(ctx context.Context) (*db.Setting, error) {
	settings, err := r.TransactionStore.GetSettings()
	if err != nil {
		return nil, terror.New(err, "get settings (3)")
	}
	if settings.SmartContractAddress != "" {
		return nil, terror.New(blockchain.ErrSmartContractAlreadyDeployed, "")
	}

	address, _, _, err := r.Blk.DeploySmartContract(ctx, "Genesis Smart Contract")
	if err != nil {
		return nil, terror.New(err, "deploy smart contract")
	}

	settings.SmartContractAddress = address.Hex()
	_, err = r.TransactionStore.UpdateSettings(settings)
	if err != nil {
		return nil, terror.New(err, "update settings")
	}

	// Track user activity
	r.RecordUserActivity(ctx, "Deployed Smart Contract", graphql.ObjectTypeBlockchain, nil, nil)

	return settings, nil
}

func (r *queryResolver) GetTickerInfo(ctx context.Context) (*graphql.TickerInfo, error) {
	info := &graphql.TickerInfo{
		LastTick:     r.SystemTicker.LastTick,
		TickInterval: int(r.SystemTicker.TickInterval.Hours()),
	}
	return info, nil
}
