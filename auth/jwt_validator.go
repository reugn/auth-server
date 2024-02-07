package auth

import (
	"encoding/json"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/reugn/auth-server/repository"
)

// JWTValidator validates and authorizes an AccessToken.
type JWTValidator struct {
	keys *Keys
	repo repository.Repository
}

// NewJWTValidator returns a new instance of JWTValidator.
func NewJWTValidator(key *Keys, repo repository.Repository) *JWTValidator {
	return &JWTValidator{key, repo}
}

// validate validates the AccessToken.
func (val *JWTValidator) validate(jtwToken string) (*Claims, error) {
	token, err := jwt.Parse(jtwToken, func(token *jwt.Token) (interface{}, error) {
		return val.keys.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return val.validateClaims(token)
}

func (val *JWTValidator) validateClaims(token *jwt.Token) (*Claims, error) {
	claims, err := getClaims(token)
	if err != nil {
		return nil, err
	}

	// validate expiration
	if claims.ExpiresAt < time.Now().Unix() {
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
func (val *JWTValidator) Authorize(token string, request *repository.RequestDetails) bool {
	claims, err := val.validate(token)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return val.repo.AuthorizeRequest(claims.Role, *request)
}
