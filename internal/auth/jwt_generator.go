package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/repository"
	"github.com/reugn/auth-server/internal/util/env"
)

const (
	envTokenExpireAfterMillis = "AUTH_SERVER_ACCESS_TOKEN_EXPIRATION_MILLIS"
)

// JWTGenerator generates an AccessToken.
type JWTGenerator struct {
	keys             *Keys
	signingMethod    jwt.SigningMethod
	tokenExpireAfter time.Duration
}

// NewJWTGenerator returns a new instance of JWTGenerator.
func NewJWTGenerator(keys *Keys, signingMethod jwt.SigningMethod) *JWTGenerator {
	tokenExpireAfter := time.Hour // default 1 hour
	env.ReadTime(&tokenExpireAfter, envTokenExpireAfterMillis, time.Millisecond)
	return &JWTGenerator{
		keys:             keys,
		signingMethod:    signingMethod,
		tokenExpireAfter: tokenExpireAfter,
	}
}

// Generate generates an AccessToken using the username and role claims.
func (gen *JWTGenerator) Generate(username string, role repository.UserRole) (*AccessToken, error) {
	token := jwt.New(gen.signingMethod)
	claims := Claims{}

	// set custom claims
	claims.Username = username
	claims.Role = role

	// set standard claims
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)
	if gen.tokenExpireAfter > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(gen.tokenExpireAfter))
	}

	token.Claims = &claims
	signed, err := token.SignedString(gen.keys.privateKey)
	if err != nil {
		return nil, err
	}

	// create an access token
	accessToken := &AccessToken{
		signed,
		BearerToken.ToString(),
		gen.tokenExpireAfter.Milliseconds(),
	}

	return accessToken, nil
}
