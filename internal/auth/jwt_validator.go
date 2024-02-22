package auth

import (
	"encoding/json"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/repository"
)

// JWTValidator validates and authorizes an AccessToken.
type JWTValidator struct {
	keys    *Keys
	backend repository.Repository
}

// NewJWTValidator returns a new instance of JWTValidator.
func NewJWTValidator(keys *Keys, backend repository.Repository) *JWTValidator {
	return &JWTValidator{
		keys:    keys,
		backend: backend,
	}
}

// validate validates the AccessToken.
func (v *JWTValidator) validate(jtwToken string) (*Claims, error) {
	token, err := jwt.Parse(jtwToken, func(_ *jwt.Token) (interface{}, error) {
		return v.keys.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return v.validateClaims(token)
}

func (v *JWTValidator) validateClaims(token *jwt.Token) (*Claims, error) {
	claims, err := getClaims(token)
	if err != nil {
		return nil, err
	}

	// validate expiration
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

func getClaims(token *jwt.Token) (*Claims, error) {
	mapClaims := token.Claims.(jwt.MapClaims)
	jsonClaims, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}

	claims := Claims{}
	err = json.Unmarshal(jsonClaims, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

// Authorize validates the token and authorizes the actual request.
func (v *JWTValidator) Authorize(token string, request *repository.RequestDetails) bool {
	claims, err := v.validate(token)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return v.backend.AuthorizeRequest(claims.Role, *request)
}
