package server

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
	"skybluetrades.net/work-planning-demo/model"
)

// NewAuthenticator creates a new JWT-based authenticator to use in
// the Echo middleware. This checks a bearer token on each request,
// and if that's successful, it extracts token claims into the Echo
// request context for later use.
func NewAuthenticator(cfg *Config) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// Not sure when this might go wrong, but the oapi-codegen
		// examples check it!
		if input.SecuritySchemeName != "BearerAuth" {
			return errors.New("incorrect authentication security scheme")
		}

		// Extract the access token from the Authorization header.
		token, err := parseAuthHeader(input.RequestValidationInput.Request)
		if err != nil {
			return err
		}

		// Parse and validate the access token.
		claims, err := ValidateToken(token, cfg.AuthKey)
		if err != nil {
			return err
		}

		// Check claims for particular scopes from the OpenAPI spec (here,
		// we just have an "admin" scope).
		for _, scope := range input.Scopes {
			if scope == "admin" && !claims.IsAdmin {
				return errors.New("admin action not permitted")
			}
		}

		// Save the access token claims for later processing.
		echoCtx := middleware.GetEchoContext(ctx)
		echoCtx.Set("claims", claims)

		return nil
	}
}

// JWTClaim is the claim structure for JWT access tokens.
type JWTClaim struct {
	ID      model.WorkerID `json:"id"`
	IsAdmin bool           `json:"is_admin"`
	jwt.StandardClaims
}

// JWTRefreshClaim is the claim structure for JWT refresh tokens.
type JWTRefreshClaim struct {
	ID model.WorkerID `json:"id"`
	jwt.StandardClaims
}

// GenerateTokens creates access and refresh tokens for a given
// Worker.
func GenerateTokens(worker *model.Worker, cfg *Config) (string, string, error) {
	// Create the access claims for the worker. This stores the user ID
	// and whether the user is an admin.
	claims := &JWTClaim{
		ID:      worker.ID,
		IsAdmin: worker.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(cfg.AccessTokenLease) * time.Second).Unix(),
		},
	}

	// Create and sign the access token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.AuthKey))
	if err != nil {
		return "", "", err
	}

	// Create the refresh claims for the worker. This stores just the
	// user ID. When a token refresh is requested, the user is looked up
	// in the database, and tokens are regenerated — this means that any
	// changes to user characteristics that affect access token claims
	// become active at the next token refresh.
	refreshClaims := &JWTRefreshClaim{
		ID: worker.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(cfg.RefreshTokenLease) * time.Second).Unix(),
		},
	}

	// Create and sign the refresh token.
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString([]byte(cfg.AuthKey))
	if err != nil {
		return "", "", err
	}

	return t, rt, nil
}

// Extract JWT from request headers.
func parseAuthHeader(r *http.Request) (string, error) {
	// Get the Authorization header, which should be of the form "Bearer
	// <access_token>".
	auths, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New("no Authorization header")
	}

	// Check "Bearer " prefix.
	if !strings.HasPrefix(auths[0], "Bearer ") {
		return "", errors.New("invalid Authorization header")
	}

	// Split off JWT.
	return auths[0][7:], nil
}

// ValidateToken validates an access token.
func ValidateToken(signedToken string, secretKey string) (*JWTClaim, error) {
	// Parse the access token claims.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	// The Valid flag on the token will be false if the token has
	// expired.
	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	// Extract the claims, which is the thing we care about here.
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	return claims, nil
}

// ValidateRefreshToken validates an access token.
func ValidateRefreshToken(signedToken string, secretKey string) (*JWTRefreshClaim, error) {
	// Parse the refresh token claims.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTRefreshClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	// The Valid flag on the token will be false if the token has
	// expired.
	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	// Extract the claims, which is the thing we care about here.
	claims, ok := token.Claims.(*JWTRefreshClaim)
	if !ok {
		return nil, errors.New("couldn't parse refresh claims")
	}
	return claims, nil
}
