package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/reugn/auth-server/repository"
)

// JWTGenerator generates an AccessToken
type JWTGenerator struct {
	Keys          *Keys
	SigningMethod jwt.SigningMethod
}

// NewJWTGenerator returns a new instance of JWTGenerator
func NewJWTGenerator(keys *Keys, signingMethod jwt.SigningMethod) *JWTGenerator {
	return &JWTGenerator{keys, signingMethod}
}

// Generate generates an AccessToken using username and role claims
func (gen *JWTGenerator) Generate(username string, role repository.UserRole) (*AccessToken, error) {
	token := jwt.New(gen.SigningMethod)
	claims := Claims{}

	// set custom claims
	claims.Username = username
	claims.Role = role

	// set standard claims
	now := time.Now()
	claims.IssuedAt = now.Unix()
	if env.expireAfter > 0 {
		claims.ExpiresAt = now.Add(env.expireAfter).Unix()
	}

	token.Claims = &claims
	signed, err := token.SignedString(gen.Keys.PrivateKey)
	if err != nil {
		return nil, err
	}

	// create an access token
	accessToken := &AccessToken{
		signed,
		BearerToken.ToString(),
		env.expireAfter.Milliseconds(),
	}

	return accessToken, nil
}
