package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/repository"
)

// JWTValidator validates and authorizes an AccessToken. It uses the keys to validate
// the token and the backend to authorize the request.
type JWTValidator struct {
	keys    *Keys
	backend repository.Repository
}

// NewJWTValidator returns a new JWTValidator.
func NewJWTValidator(keys *Keys, backend repository.Repository) *JWTValidator {
	return &JWTValidator{
		keys:    keys,
		backend: backend,
	}
}

// validate validates the AccessToken. It returns the claims if the token is valid,
// an error otherwise.
func (v *JWTValidator) validate(jwtToken string) (*Claims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return v.keys.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return v.validateClaims(token)
}

// validateClaims validates the claims of the token.
func (v *JWTValidator) validateClaims(token *jwt.Token) (*Claims, error) {
	claims, err := getClaims(token)
	if err != nil {
		return nil, err
	}

	// validate expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		slog.Debug("Token expired")
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

// getClaims gets the claims from the token.
func getClaims(token *jwt.Token) (*Claims, error) {
	mapClaims := token.Claims.(jwt.MapClaims)
	jsonClaims, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	claims := Claims{}
	err = json.Unmarshal(jsonClaims, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &claims, nil
}

// Authorize validates the token and authorizes the actual request. It returns true
// if the request is authorized, false otherwise.
func (v *JWTValidator) Authorize(token string, request *repository.RequestDetails) bool {
	claims, err := v.validate(token)
	if err != nil {
		slog.Debug("Failed to authorize token", "err", err)
		return false
	}

	return v.backend.AuthorizeRequest(claims.Role, *request)
}
