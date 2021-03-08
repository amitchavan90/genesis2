package genesis

import (
	"context"
	"encoding/base64"
	"fmt"
	"genesis/canlog"
	"genesis/db"
	"genesis/graphql"
	"math/rand"
	"net/http"
	"time"

	"github.com/ninja-software/terror"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofrs/uuid"

	"github.com/go-chi/jwtauth"
)

// AuthProvider contains methods for authentication
type AuthProvider interface {
	// Get current user and related role w/ permissions
	UserFromContext(ctx context.Context) (*db.User, error)
	// Get current user id
	UserIDFromContext(ctx context.Context) (uuid.UUID, error)
	// Get current user's role w/ track actions if incTrackActions = true
	RoleFromContext(ctx context.Context, incTrackActions bool) (*db.Role, error)
	// Get current user's role id
	RoleIDFromContext(ctx context.Context) (uuid.UUID, error)
	GenerateJWT(ctx context.Context, user string, userID string, roleID string, userAgent string, expiration *time.Time) (string, error)
	ValidatePassword(ctx context.Context, email string, password string) error
}

// Auther to handle JWT authentication
type Auther struct {
	TokenExpirationDays int
	TokenAuth           *jwtauth.JWTAuth
	Blacklister         BlacklistProvider
	TokenStore          TokenStorer
	UserStore           UserStorer
	RoleStore           RoleStorer
}

// NewAuther for JWT and blacklisting
func NewAuther(jwtsecret string, userStore UserStorer, blacklister BlacklistProvider, tokenStore TokenStorer, roleStore RoleStorer, tokenExpirationDays int) *Auther {
	result := &Auther{
		TokenAuth:           jwtauth.New("HS256", []byte(jwtsecret), []byte(jwtsecret)),
		Blacklister:         blacklister,
		TokenStore:          tokenStore,
		UserStore:           userStore,
		RoleStore:           roleStore,
		TokenExpirationDays: tokenExpirationDays,
	}
	return result
}

// ClaimsFromContext a map of all claims in JWT
func ClaimsFromContext(ctx context.Context) (map[string]string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	result := map[string]string{}
	for k, v := range claims {
		val, ok := v.(string)
		if !ok {
			continue
		}
		result[k] = val
	}

	if err != nil {
		return result, terror.New(err, "")
	}

	return result, nil
}

// ClaimKey is a type used to set values in the JWT
type ClaimKey string

// ClaimUserName JWT key value
const ClaimUserName ClaimKey = "username"

// ClaimUserID JWT key value
const ClaimUserID ClaimKey = "uid"

// ClaimRoleID JWT key value
const ClaimRoleID ClaimKey = "roleId"

// ClaimTokenID JWT key value
const ClaimTokenID ClaimKey = "tokenId"

// ClaimExistsInContext returns a specific claim value from the JWT
func ClaimExistsInContext(ctx context.Context, key ClaimKey) bool {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return false
	}
	_, ok := claims.Get(string(key))
	if !ok {
		return false
	}
	return true
}

// ClaimValueFromContext returns a specific claim value from the JWT
func ClaimValueFromContext(ctx context.Context, key ClaimKey) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", terror.New(terror.ErrBadContext, "")
	}
	valI, ok := claims.Get(string(key))
	if !ok {
		canlog.Set(ctx, "claimkey", key)
		return "", terror.New(terror.ErrBadClaims, "")
	}
	val, ok := valI.(string)
	if !ok {
		return "", terror.New(terror.ErrTypeCast, "")
	}
	return val, nil
}

// UserFromContext grabs the user from the context if a JWT is inside (includes role w/ permissions)
func (a *Auther) UserFromContext(ctx context.Context) (*db.User, error) {
	id, err := a.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(err, "")
	}

	u, err := a.UserStore.GetWithPermissions(id)
	if err != nil {
		return nil, terror.New(err, "get user with permissions")
	}

	return u, nil
}

// UserIDFromContext will validate JWT and return user ID
func (a *Auther) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	jwtToken, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, terror.New(terror.ErrBadContext, "")
	}

	if jwtauth.IsExpired(jwtToken) {
		return uuid.Nil, terror.New(fmt.Errorf("token has expired"), "")
	}

	userIDI, ok := claims.Get(string(ClaimUserID))
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimUserID)
		return uuid.Nil, terror.New(terror.ErrBadClaims, "")
	}

	tokenID, ok := claims.Get("tokenId")
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimTokenID)
		return uuid.Nil, terror.New(terror.ErrBadClaims, "")
	}

	blacklisted := a.Blacklister.OnList(tokenID.(string))
	if blacklisted {
		return uuid.Nil, terror.New(terror.ErrBlacklisted, "")
	}

	userIDStr, ok := userIDI.(string)
	if !ok {
		return uuid.Nil, terror.New(terror.ErrTypeCast, "")
	}
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return uuid.Nil, terror.New(terror.ErrParse, "")
	}

	return userID, nil
}

// RoleIDFromContext will validate JWT and return role ID
func (a *Auther) RoleIDFromContext(ctx context.Context) (uuid.UUID, error) {
	jwtToken, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, terror.New(terror.ErrBadContext, "")
	}

	if jwtauth.IsExpired(jwtToken) {
		return uuid.Nil, terror.New(fmt.Errorf("token has expired"), "")
	}

	roleIDI, ok := claims.Get(string(ClaimRoleID))
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimRoleID)
		return uuid.Nil, terror.New(terror.ErrBadClaims, "")
	}

	tokenID, ok := claims.Get("tokenId")
	if !ok {
		canlog.Set(ctx, "claimkey", ClaimTokenID)
		return uuid.Nil, terror.New(terror.ErrBadClaims, "")
	}

	blacklisted := a.Blacklister.OnList(tokenID.(string))
	if blacklisted {
		return uuid.Nil, terror.New(terror.ErrBlacklisted, "")
	}

	roleIDStr, ok := roleIDI.(string)
	if !ok {
		return uuid.Nil, terror.New(terror.ErrTypeCast, "")
	}
	roleID, err := uuid.FromString(roleIDStr)
	if err != nil {
		return uuid.Nil, terror.New(terror.ErrParse, "")
	}

	return roleID, nil
}

// RoleFromContext grabs the role  from the context if a JWT is inside
func (a *Auther) RoleFromContext(ctx context.Context, incTrackActions bool) (*db.Role, error) {
	id, err := a.RoleIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(err, "")
	}

	var r *db.Role
	if incTrackActions {
		r, err = a.RoleStore.GetWithTrackActions(id)
	} else {
		r, err = a.RoleStore.Get(id)
	}
	if err != nil {
		return nil, terror.New(err, "get user")
	}

	return r, nil
}

// GenerateJWT returns the token for client side persistence
func (a *Auther) GenerateJWT(ctx context.Context, username string, userID string, roleID string, userAgent string, expiration *time.Time) (string, error) {
	// Record token in issued token records

	var expr time.Time
	if expiration == nil {
		expr = time.Now().Add(time.Hour * time.Duration(24) * time.Duration(a.TokenExpirationDays))
	} else {
		expr = *expiration
	}

	//TODO: device currently being set by request.UserAgent() ... might need to parse to get a better device name
	newToken := &db.IssuedToken{
		UserID:       userID,
		Device:       userAgent,
		TokenExpires: expr,
	}
	token, err := a.TokenStore.Insert(newToken)
	if err != nil {
		return "", terror.New(err, "insert token")
	}

	_, tokenString, err := a.TokenAuth.Encode(jwtauth.Claims{
		string(ClaimUserName): username,
		string(ClaimUserID):   userID,
		string(ClaimRoleID):   roleID,
		string(ClaimTokenID):  token.ID,
	})
	if err != nil {
		return "", terror.New(err, "encode token")
	}
	return tokenString, nil
}

// VerifyMiddleware for authentication adds JWT to context down the HTTP chain
func (a *Auther) VerifyMiddleware() func(http.Handler) http.Handler {
	return jwtauth.Verifier(a.TokenAuth)
}

// ValidatePassword will check the login details
func (a *Auther) ValidatePassword(ctx context.Context, email string, password string) error {
	user, err := a.UserStore.GetByEmail(email)
	if err != nil {
		return terror.New(err, "get user")
	}

	storedHash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return terror.New(err, "decode hash")
	}

	err = bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err != nil {
		return terror.New(err, "compare hash")
	}

	return nil
}

// GenerateAlphanumericCode generates a simple alphanumeric code (excluding 1, i , l, o and 0)
func GenerateAlphanumericCode(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	code := ""
	for i := 0; i < length; i++ {
		code += string(chars[rand.Intn(len(chars)-1)])
	}
	return code
}

// HasPermission checks if a role use a permission
func HasPermission(role *db.Role, perm graphql.Perm) bool {
	for _, p := range role.Permissions {
		if p == perm.String() {
			return true
		}
	}
	return false
}
