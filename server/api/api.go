package api

import (
	"errors"
	"genesis"
	"io"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/mailgun/mailgun-go/v3"

	"context"
	"fmt"
	"genesis/blockchain"
	"genesis/canlog"
	"genesis/config"
	"genesis/graphql"
	"genesis/restpkg/routes"
	"net/http"
	"time"

	"github.com/ninja-software/terror"

	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/gqlerror"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

// NewAPIController returns the public router layer
func NewAPIController(
	log *zap.SugaredLogger,
	config *config.PlatformConfig,
	skuStore genesis.SKUStorer,
	orderStore genesis.OrderStorer,
	containerStore genesis.ContainerStorer,
	palletStore genesis.PalletStorer,
	cartonStore genesis.CartonStorer,
	productStore genesis.ProductStorer,
	organisationStore genesis.OrganisationStorer,
	userStore genesis.UserStorer,
	referralStore genesis.ReferralStorer,
	taskStore genesis.TaskStorer,
	userTaskStore genesis.UserTaskStorer,
	roleStore genesis.RoleStorer,
	blobStore genesis.BlobStorer,
	loyaltyStore genesis.LoyaltyStorer,
	contractStore genesis.ContractStorer,
	distributorStore genesis.DistributorStorer,
	transactionStore genesis.TransactionStorer,
	manifestStore genesis.ManifestStorer,
	trackActionStore genesis.TrackActionStorer,
	userActivityStore genesis.UserActivityStorer,
	userPurchaseActivityStore genesis.UserPurchaseActivityStorer,
	blacklistProvider genesis.BlacklistProvider,
	tks genesis.TokenStorer,
	jwtSecret string,
	auther *genesis.Auther,
	mailer *mailgun.MailgunImpl,
	smsMessenger genesis.Messenger,
	blk *blockchain.Service,
	systemTicker *genesis.SystemTicker,
) http.Handler {
	authentication := NewAuthRoutes(auther, config.UserAuth, userStore, roleStore, organisationStore, referralStore, userActivityStore)

	websocketUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r := chi.NewRouter()
	// r.Use(debugClientInput())
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(auther.VerifyMiddleware())
	r.Use(canonicalLogger(log.Desugar()))
	r.Use(genesis.DataloaderMiddleware(
		skuStore,
		orderStore,
		distributorStore,
		containerStore,
		palletStore,
		cartonStore,
		productStore,
		organisationStore,
		userStore,
		referralStore,
		taskStore,
		userTaskStore,
		roleStore,
		userPurchaseActivityStore,
		contractStore,
		transactionStore,
		manifestStore,
		trackActionStore,
	))
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Mount("/api/auth", authentication)
	r.Route("/api/gql", func(r chi.Router) {
		r.Handle("/", handler.Playground("GraphQL playground", "/api/gql/query"))
		r.Handle("/query", handler.GraphQL(
			graphql.NewExecutableSchema(
				graphql.Config{
					Resolvers: &genesis.Resolver{
						Auther:                    auther,
						Config:                    config,
						Mailer:                    mailer,
						SKUStore:                  skuStore,
						OrderStore:                orderStore,
						ContainerStore:            containerStore,
						PalletStore:               palletStore,
						CartonStore:               cartonStore,
						ProductStore:              productStore,
						OrganisationStore:         organisationStore,
						UserStore:                 userStore,
						ReferralStore:             referralStore,
						TaskStore:                 taskStore,
						UserTaskStore:             userTaskStore,
						RoleStore:                 roleStore,
						BlobStore:                 blobStore,
						LoyaltyStore:              loyaltyStore,
						ContractStore:             contractStore,
						DistributorStore:          distributorStore,
						TransactionStore:          transactionStore,
						ManifestStore:             manifestStore,
						TrackActionStore:          trackActionStore,
						UserActivityStore:         userActivityStore,
						UserPurchaseActivityStore: userPurchaseActivityStore,
						Blacklister:               blacklistProvider,
						SmsMessenger:              smsMessenger,
						Blk:                       blk,
						SystemTicker:              systemTicker,
					},
					Directives: genesis.NewDirectiveRoot(roleStore),
				},
			),
			handler.ErrorPresenter(
				func(ctx context.Context, e error) *gqlerror.Error {
					// non-panic error recovery
					echoMsg := terror.Echo(e)

					if errors.Is(e, terror.ErrBadContext) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("There was a problem reading your credentials. Please sign in and try again."))
					}
					if errors.Is(e, terror.ErrBadClaims) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("There was a problem reading your credentials. Please sign in and try again."))
					}
					if errors.Is(e, terror.ErrBlacklisted) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Login token is no longer valid. Please login again."))
					}
					if errors.Is(e, terror.ErrUnauthorized) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("You are not authorized to do this action."))
					}
					if errors.Is(e, terror.ErrBadCredentials) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Please check your username or password and try again."))
					}
					if errors.Is(e, terror.ErrNotImplemented) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("This functionality is not yet available."))
					}
					if errors.Is(e, blockchain.ErrBlockchainConnectionIssue) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Unable to connect to smart contact."))
					}
					if errors.Is(e, blockchain.ErrBlockchainOutOfGas) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Blockchain smart contract out of gas."))
					}
					if errors.Is(e, blockchain.ErrSmartContractAlreadyDeployed) {
						canlog.SetErr(ctx, e)
						canlog.LogError(ctx, "request")
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Blockchain smart contract already deployed."))
					}
					if errors.Is(e, terror.ErrAuthWrongPassword) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Password did not match existing one."))
					}
					if errors.Is(e, genesis.ErrPasswordShort) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Password too short."))
					}
					if errors.Is(e, genesis.ErrPasswordCommon) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Password too common."))
					}
					if errors.Is(e, genesis.ErrEmailInvalid) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Invalid email address."))
					}
					if errors.Is(e, genesis.ErrMobileNotSet) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("No mobile number found."))
					}
					if errors.Is(e, genesis.ErrTokenInvalid) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Token is invalid."))
					}
					if errors.Is(e, genesis.ErrTokenExpired) {
						return gqlgraphql.DefaultErrorPresenter(ctx, fmt.Errorf("Token is expired."))
					}

					var bErr *terror.Error
					if errors.As(e, &bErr) {
						user, _ := auther.UserFromContext(ctx)
						// reconstruct since we dont have http.Request
						req := &http.Request{
							RemoteAddr: canlog.Get(ctx, "ip").(string),
							Method:     canlog.Get(ctx, "req_method").(string),
							RequestURI: canlog.Get(ctx, "req_uri").(string),
							Header:     canlog.Get(ctx, "req_header").(http.Header),
							Body:       canlog.Get(ctx, "req_body").(io.ReadCloser),
						}
						genesis.SentrySend(ctx, user, req, e, echoMsg)

						// generic error
						// e = errors.New("There was a problem with the server.   Please try again later")

						// extra spacing for the Please to indicate us it is errors.As route
						return gqlgraphql.DefaultErrorPresenter(ctx, e)
					}

					return gqlgraphql.DefaultErrorPresenter(ctx, e)
				},
			),
			handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
				// panic recovery

				// construct error to terror.Error then pass to .ErrorPresenter
				return terror.NewPanic(err.(error))
			}),
			handler.WebsocketKeepAliveDuration(5*time.Second),
			handler.WebsocketUpgrader(websocketUpgrader),
		),
		)
	})
	r.Get(config.API.BlobBaseURL, downloadHandler(blobStore))
	r.Get(fmt.Sprintf("%sorder", config.API.BlobBaseURL), genesis.DownloadSpreadSheetOrder(orderStore, auther, config.API))
	r.Get(fmt.Sprintf("%ssheet", config.API.BlobBaseURL), genesis.DownloadSpreadSheet(cartonStore, palletStore, containerStore, auther, config.API))

	r.Get("/api/steak/view", genesis.SteakView(
		productStore,
		distributorStore,
	))
	r.Get("/api/steak/detail", genesis.SteakDetail(
		productStore,
		skuStore,
		loyaltyStore,
		distributorStore,
	))
	r.Post("/api/steak/close", genesis.SteakClose(
		userStore,
		productStore,
		roleStore,
		organisationStore,
		loyaltyStore,
		skuStore,
		transactionStore,
		trackActionStore,
		blk,
	))
	r.Get("/api/steak/final", genesis.SteakFinal(
		productStore,
		loyaltyStore,
		skuStore,
	))
	r.Get("/api/manifests", genesis.ProofOfSteakManifestBasicAll(
		config.Blockchain.EtherscanHost,
		transactionStore,
		manifestStore,
		trackActionStore,
		blk,
	))
	r.Get("/api/manifest/mr/{merkleRoot}", genesis.ProofOfSteakManifestByMerkleRoot(
		config.Blockchain.EtherscanHost,
		transactionStore,
		manifestStore,
		trackActionStore,
		blk,
	))
	r.Get("/api/manifest/tx/{txID}", genesis.ProofOfSteakManifestByHash(
		config.Blockchain.EtherscanHost,
		transactionStore,
		manifestStore,
		trackActionStore,
		blk,
	))
	r.Get("/api/manifest/line/{lineHash}", genesis.ProofOfSteakManifestByLine(
		config.Blockchain.EtherscanHost,
		transactionStore,
		manifestStore,
		trackActionStore,
		blk,
	))
	// Map rest urls
	routes.MapUrls(r)

	return r
}
