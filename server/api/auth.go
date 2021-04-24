package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"genesis"
	"genesis/config"
	"genesis/crypto"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
)

// AuthController contains handlers involving authentication
type AuthController struct {
	UserStore         genesis.UserStorer
	RoleStore         genesis.RoleStorer
	OrganisationStore genesis.OrganisationStorer
	ReferralStore     genesis.ReferralStorer
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
func NewAuthRoutes(auther *genesis.Auther, authConfig *config.UserAuth, userStore genesis.UserStorer, roleStore genesis.RoleStorer, organisationStore genesis.OrganisationStorer, referralStore genesis.ReferralStorer, userActivityStore genesis.UserActivityStorer) chi.Router {
	cookieDefaults := CookieSettings{
		SameSite: http.SameSiteNoneMode,
		HTTPOnly: true,
		Path:     "/",
	}
	c := &AuthController{
		userStore,
		roleStore,
		organisationStore,
		referralStore,
		userActivityStore,
		auther,
		authConfig,
		cookieDefaults,
	}
	r := chi.NewRouter()
	r.Get("/roles", c.roleList())
	r.Post("/register", c.register())
	r.Post("/login", c.login())
	r.Post("/consumer/login", c.consumerLogin())
	r.Post("/logout", c.logout())
	r.Post("/verify_account", c.verifyAccount())

	// Social Login
	r.Post("/social_login", c.socialLogin())
	r.Get("/google_auth/home", c.handleGoogleMain())
	r.Get("/google_auth/login", c.handleGoogleLogin())
	r.Get("/google_auth/callback", c.handleGoogleCallback())
	r.Get("/facebook_auth/home", c.handleFacebookMain())
	r.Get("/facebook_auth/login", c.handleFacebookLogin())
	r.Get("/facebook_auth/callback", c.handleFacebookCallback())
	return r
}

// RegisterRequest structs for the HTTP request/response cycle
type RegisterRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	MobilePhone    string `json:"mobilePhone"`
	Password       string `json:"password"`
	RoleID         string `json:"roleID"`
	AffilatedOrg   string `json:"affiliateOrg"`
	ReferredByCode string `json:"referredByCode"`
}

type RegisterResponse struct {
	FirstName    null.String `json:"firstName"`
	LastName     null.String `json:"lastName"`
	Email        null.String `json:"email"`
	MobilePhone  null.String `json:"mobilePhone"`
	RoleID       string      `json:"roleID"`
	AffilatedOrg null.String `json:"affiliateOrg"`
	ReferralCode null.String `json:"referralCode"`
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

// RestResponse handles the http status and renders body in JSON
func RestResponse(w http.ResponseWriter, r *http.Request, status int, body interface{}) {
	render.Status(r, status)
	render.JSON(w, r, body)
}

// list roles
func (c *AuthController) roleList() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		roles, err := c.RoleStore.All()
		if err != nil {
			failedMsg := "Failed to load role list, please try again."
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		RestResponse(w, r, http.StatusOK, roles)
	}

	return fn
}

// register creates a user
func (c *AuthController) register() func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		input := &RegisterRequest{}
		user := &db.User{}

		failedMsg := "Register failed, please try again."

		err := json.NewDecoder(r.Body).Decode(input)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if input.Email == "" {
			RestResponse(w, r, http.StatusBadRequest, "Email is required")
			return
		}
		if input.FirstName == "" {
			RestResponse(w, r, http.StatusBadRequest, "First Name is required")
			return
		}
		if input.LastName == "" {
			RestResponse(w, r, http.StatusBadRequest, "Last Name is required")
			return
		}
		if input.MobilePhone == "" {
			RestResponse(w, r, http.StatusBadRequest, "Mobile number is required")
			return
		}
		if input.RoleID == "" {
			RestResponse(w, r, http.StatusBadRequest, "Role ID is required")
			return
		}
		if input.Password == "" {
			RestResponse(w, r, http.StatusBadRequest, "Password is required")
			return
		}

		email := strings.ToLower(input.Email)

		// Verify email
		u, _ := c.UserStore.GetByEmail(email)
		if u != nil {
			failedMsg := "Email already registered, please try again."
			RestResponse(w, r, http.StatusBadRequest, failedMsg)
			return
		}

		// Verify email
		d, _ := c.UserStore.GetByMobilePhone(input.MobilePhone)
		if d != nil {
			failedMsg := "Mobile number already registered, please try again."
			RestResponse(w, r, http.StatusBadRequest, failedMsg)
			return
		}

		// Verify role
		roleUUID, err := uuid.FromString(input.RoleID)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}
		_, err = c.auther.RoleStore.Get(roleUUID)
		if err != nil {
			failedMsg := "Role does not exist. Please verify role ID"
			RestResponse(w, r, http.StatusBadRequest, failedMsg)
			return
		}

		// Verify affilited org
		if input.AffilatedOrg != "" {
			orgUUID, err := uuid.FromString(input.RoleID)
			if err != nil {
				writeError(w, err, failedMsg, http.StatusInternalServerError)
				return
			}
			_, err = c.OrganisationStore.Get(orgUUID)
			if err != nil {
				failedMsg := "Organisation does not exist. Please verify organisation ID"
				RestResponse(w, r, http.StatusBadRequest, failedMsg)
				return
			}
			user.AffiliateOrg = null.StringFrom(input.AffilatedOrg)
		}

		// Set/Generate password
		password := ""
		if input.Password != "" {
			password = input.Password
		} else {
			g, err := uuid.NewV4()
			if err != nil {
				RestResponse(w, r, http.StatusInternalServerError, failedMsg)
				return
			}
			password = g.String()
		}
		hashed := crypto.HashPassword(password)
		input.Password = hashed

		user.FirstName = null.StringFrom(input.FirstName)
		user.LastName = null.StringFrom(input.LastName)
		user.Email = null.StringFrom(input.Email)
		user.MobilePhone = null.StringFrom(input.MobilePhone)
		user.RoleID = input.RoleID
		user.PasswordHash = hashed
		user.ReferralCode = null.StringFrom(helpers.GenerateID(7))

		created, err := c.UserStore.Insert(user)
		if err != nil {
			RestResponse(w, r, http.StatusInternalServerError, failedMsg)
			return
		}

		// Verify referee by referral code
		if input.ReferredByCode != "" {
			referee, err := c.UserStore.GetByReferralCode(input.ReferredByCode)
			if err != nil {
				failedMsg := "Referee does not exist. Please verify referral code"
				RestResponse(w, r, http.StatusBadRequest, failedMsg)
				return
			}

			// Get Referral count
			count, err := c.ReferralStore.Count()
			if err != nil {
				failedMsg := "create referral: Error while fetching referral count from db"
				RestResponse(w, r, http.StatusBadRequest, failedMsg)
				return
			}

			refID, _ := uuid.NewV4()
			referral := &db.Referral{
				Code:         fmt.Sprintf("R%05d", count+1),
				ID:           refID.String(),
				UserID:       created.ID,
				ReferredByID: null.StringFrom(referee.ID),
				IsRedemmed:   false,
			}

			_, err = c.ReferralStore.Insert(referral)
			if err != nil {
				failedMsg := "Error while creating referral"
				RestResponse(w, r, http.StatusInternalServerError, failedMsg)
				return
			}
		}

		// Gnerate token
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

		response := RegisterResponse{}
		response.FirstName = created.FirstName
		response.LastName = created.LastName
		response.Email = created.Email
		response.MobilePhone = created.MobilePhone
		response.RoleID = created.RoleID
		response.ReferralCode = created.ReferralCode

		RestResponse(w, r, http.StatusOK, response)
	}

	return fn
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

// consumerLogin logs a consumer in
func (c *AuthController) consumerLogin() func(w http.ResponseWriter, r *http.Request) {
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

		// consumer role
		consumerRole := "CONSUMER"
		roleUUID, err := uuid.FromString(user.RoleID)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		role, err := c.RoleStore.Get(roleUUID)
		if err != nil {
			writeError(w, err, failedMsg, http.StatusInternalServerError)
			return
		}

		if role.Name != consumerRole {
			failedMsg := "Invalid consumer, please try again."
			err := terror.New(nil, failedMsg)
			writeError(w, err, failedMsg, http.StatusInternalServerError)
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

// socialLogin logs a user in
func (c *AuthController) socialLogin() func(w http.ResponseWriter, r *http.Request) {
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
