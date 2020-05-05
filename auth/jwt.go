package auth

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/reugn/auth-server/repository"
)

// TokenType enum
type TokenType int

const (
	// A BearerToken is an opaque string, not intended to have any meaning to clients using it.
	// Some servers will issue tokens that are a short `string` of hexadecimal characters,
	// while others may use structured tokens such as JSON Web Tokens.
	BearerToken TokenType = iota

	// A BasicToken is a string where credentials is the base64 encoding of id and
	// password joined by a single colon :
	BasicToken
)

// ToString converts TokenType to a string
func (t TokenType) ToString() string {
	return [...]string{"Bearer", "Basic"}[t]
}

// Claims is a custom JWT claims model
type Claims struct {
	jwt.StandardClaims
	Username string              `json:"user"`
	Role     repository.UserRole `json:"role"`
}

// AccessToken model
type AccessToken struct {
	Token   string `json:"access_token"`
	Type    string `json:"token_type"`
	Expires int64  `json:"expires_in"`
}

// Marshal AccessToken to a json string
func (t *AccessToken) Marshal() string {
	jsonByteArray, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(jsonByteArray)
}
