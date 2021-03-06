package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"genesis"
	"io"
	"strings"

	gqlgraphql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gofrs/uuid"
	"github.com/mailgun/mailgun-go/v3"

	"context"
	"fmt"
	"genesis/blockchain"
	"genesis/canlog"
	"genesis/config"
	"genesis/db"
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
	roleStore genesis.RoleStorer,
	blobStore genesis.BlobStorer,
	loyaltyStore genesis.LoyaltyStorer,
	contractStore genesis.ContractStorer,
	distributorStore genesis.DistributorStorer,
	transactionStore genesis.TransactionStorer,
	manifestStore genesis.ManifestStorer,
	trackActionStore genesis.TrackActionStorer,
	userActivityStore genesis.UserActivityStorer,
	blacklistProvider genesis.BlacklistProvider,
	tks genesis.TokenStorer,
	jwtSecret string,
	auther *genesis.Auther,
	mailer *mailgun.MailgunImpl,
	smsMessenger genesis.Messenger,
	blk *blockchain.Service,
	systemTicker *genesis.SystemTicker,
) http.Handler {
	authentication := NewAuthRoutes(auther, config.UserAuth, userStore, userActivityStore)

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
		roleStore,
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
						Auther:            auther,
						Config:            config,
						Mailer:            mailer,
						SKUStore:          skuStore,
						OrderStore:        orderStore,
						ContainerStore:    containerStore,
						PalletStore:       palletStore,
						CartonStore:       cartonStore,
						ProductStore:      productStore,
						OrganisationStore: organisationStore,
						UserStore:         userStore,
						ReferralStore:     referralStore,
						TaskStore:         taskStore,
						RoleStore:         roleStore,
						BlobStore:         blobStore,
						LoyaltyStore:      loyaltyStore,
						ContractStore:     contractStore,
						DistributorStore:  distributorStore,
						TransactionStore:  transactionStore,
						ManifestStore:     manifestStore,
						TrackActionStore:  trackActionStore,
						UserActivityStore: userActivityStore,
						Blacklister:       blacklistProvider,
						SmsMessenger:      smsMessenger,
						Blk:               blk,
						SystemTicker:      systemTicker,
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

// AuthController contains handlers involving authentication
type AuthController struct {
	UserStore         genesis.UserStorer
	UserActivityStore genesis.UserActivityStorer
	auther            *genesis.Auther
	authConfig        *config.UserAuth
	cookieDefaults    CookieSettings
}

// CookieSettings are the default values used to set cookies
type CookieSettings struct {
	SameSite http.SameSite
	HTTPOnly bool
	Secure   bool
	Path     string
}

// NewAuthRoutes returns a router for use in authentication
func NewAuthRoutes(auther *genesis.Auther, authConfig *config.UserAuth, userStore genesis.UserStorer, userActivityStore genesis.UserActivityStorer) chi.Router {
	cookieDefaults := CookieSettings{
		SameSite: http.SameSiteNoneMode,
		HTTPOnly: true,
		Path:     "/",
	}
	c := &AuthController{
		userStore,
		userActivityStore,
		auther,
		authConfig,
		cookieDefaults,
	}
	r := chi.NewRouter()
	r.Post("/login", c.login())
	r.Post("/logout", c.logout())
	r.Post("/verify_account", c.verifyAccount())
	return r
}

// LoginRequest structs for the HTTP request/response cycle
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse structs for the HTTP request/response cycle
type LoginResponse struct {
	Verified bool `json:"verified"`
	Success  bool `json:"success"`
}

// writeError writes a http errors to the ResponseWriter
func writeError(w http.ResponseWriter, err error, message string, code int) {
	http.Error(w, fmt.Sprintf(`{"error":"%s","message":"%s"}`, err.Error(), message), code)
}

// login logs a user in
func (c *AuthController) login() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		req := &LoginRequest{}

		failedMsg := "Login failed, please try again."

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		email := strings.ToLower(req.Email)

		user, err := c.UserStore.GetByEmail(email)
		if err != nil {
			failedMsg := "Invalid user, please try again."
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		userUUID, err := uuid.FromString(user.ID)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		u, err := c.UserStore.GetWithPermissions(userUUID)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		// doesn't have permission to use admin portal
		if !hasPerm(graphql.PermUseAdminPortal.String(), u.R.Role.Permissions) {
			failedMsg := "Must be an admin"
			writeError(w, errors.New(failedMsg), failedMsg, http.StatusUnauthorized)
			return
		}

		err = c.auther.ValidatePassword(r.Context(), email, req.Password)
		if err != nil {
			failedMsg := "Invalid Password"
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(time.Duration(c.authConfig.TokenExpiryDays) * time.Hour * 24)
		jwt, err := c.auther.GenerateJWT(r.Context(), user.Email.String, user.ID, user.RoleID, r.UserAgent(), &expiration)
		if err != nil {
			failedMsg := "jwt expired"
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{Name: "jwt", Value: jwt, Expires: expiration, HttpOnly: c.cookieDefaults.HTTPOnly, Path: c.cookieDefaults.Path, SameSite: c.cookieDefaults.SameSite, Secure: c.cookieDefaults.Secure}
		http.SetCookie(w, &cookie)

		// record user activity
		_, err = c.UserActivityStore.Insert(&db.UserActivity{
			UserID:     user.ID,
			Action:     "Sign in",
			ObjectType: graphql.ObjectTypeSelf.String(),
		})
		if err != nil {
			fmt.Println("update user activity: %w", err)
		}

		w.WriteHeader(http.StatusOK)
	}
	return fn
}

func hasPerm(perm string, perms []string) bool {
	for _, p := range perms {
		if p == perm {
			return true
		}
	}

	return false
}

// logout
func (c *AuthController) logout() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, err := c.auther.UserIDFromContext(r.Context())
		if err != nil {
			fmt.Println("get user on logout: %w", err)
		}

		cookie := http.Cookie{Name: "jwt", Value: "", Expires: time.Unix(0, 0), HttpOnly: c.cookieDefaults.HTTPOnly, Path: c.cookieDefaults.Path, SameSite: c.cookieDefaults.SameSite, Secure: c.cookieDefaults.Secure}
		http.SetCookie(w, &cookie)

		// record user activity
		_, err = c.UserActivityStore.Insert(&db.UserActivity{
			UserID:     userID.String(),
			Action:     "Sign out",
			ObjectType: graphql.ObjectTypeSelf.String(),
		})
		if err != nil {
			fmt.Println("update user activity: %w", err)
		}

		w.WriteHeader(http.StatusOK)
	}
	return fn
}

// verifyAccount verifies an account and logs the user in
func (c *AuthController) verifyAccount() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		type Req struct {
			Token string `json:"token"`
		}

		failedMsg := "Email verification failed, please try again."

		req := &Req{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		user, err := c.UserStore.GetByVerifyToken(req.Token)
		if err != nil && err != sql.ErrNoRows {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		if user == nil || !user.Email.Valid {
			writeError(w, err, "User with that verification code can't be found.", http.StatusInternalServerError)
			return
		}

		if user.Verified {
			err := fmt.Errorf("user already verified")
			writeError(w, err, "Verification already complete", http.StatusInternalServerError)
		}

		user.Verified = true

		user, err = c.UserStore.Update(user)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
		}

		expiration := time.Now().Add(time.Hour * time.Duration(24) * time.Duration(c.auther.TokenExpirationDays))

		// Authenticate the user so they can continue to the portal and set their password
		jwt, err := c.auther.GenerateJWT(r.Context(), user.Email.String, user.ID, user.RoleID, r.UserAgent(), &expiration)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{Name: "jwt", Value: jwt, Expires: expiration, HttpOnly: c.cookieDefaults.HTTPOnly, Path: c.cookieDefaults.Path, SameSite: c.cookieDefaults.SameSite, Secure: c.cookieDefaults.Secure}
		http.SetCookie(w, &cookie)

		// record user activity
		_, err = c.UserActivityStore.Insert(&db.UserActivity{
			UserID:     user.ID,
			Action:     "Verified Account",
			ObjectType: graphql.ObjectTypeSelf.String(),
		})
		if err != nil {
			fmt.Println("update user activity: %w", err)
		}

		w.WriteHeader(http.StatusOK)
	}
	return fn
}

func downloadHandler(BlobStore genesis.BlobStorer) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "no id provided", http.StatusBadRequest)
			return
		}

		blobUUID, err := uuid.FromString(id)
		if err != nil {
			http.Error(w, "could not parse to UUID", http.StatusBadRequest)
			return
		}
		att, err := BlobStore.Get(blobUUID)
		if err != nil {
			http.Error(w, "could not get file", http.StatusBadRequest)
			return
		}

		disposition := "attachment"
		isViewDisposition := r.URL.Query().Get("view")
		if isViewDisposition == "true" {
			disposition = "inline"
		}

		// tell the browser the returned content should be downloaded/inline
		if att.MimeType != "" && att.MimeType != "unknown" {
			w.Header().Add("Content-Type", att.MimeType)
		}
		w.Header().Add("Content-Disposition", fmt.Sprintf("%s;filename=%s", disposition, att.FileName))
		rdr := bytes.NewReader(att.File)
		http.ServeContent(w, r, att.FileName, time.Now(), rdr)

	}
	return fn
}
