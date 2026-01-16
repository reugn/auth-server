package repository

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// UserRole represents a user role.
type UserRole string

// UserDetails represents user details.
type UserDetails struct {
	UserName string
	UserRole UserRole
}

// RequestDetails represents request details.
type RequestDetails struct {
	Method string `yaml:"method"`
	URI    string `yaml:"uri"`
}

// String returns a string representation of the RequestDetails.
func (r RequestDetails) String() string {
	return fmt.Sprintf("%s %s", r.Method, r.URI)
}

// A Repository acts as a gateway to the authentication and authorization
// operations, facilitating secure access to resources.
type Repository interface {

	// AuthenticateBasic validates the basic username and password before issuing a JWT.
	AuthenticateBasic(username string, password string) *UserDetails

	// AuthorizeRequest checks if the role has permissions to access the endpoint.
	AuthorizeRequest(userRole UserRole, request RequestDetails) bool
}

// isAuthorizedRequest checks if the role has permissions to access the endpoint.
func isAuthorizedRequest(scopes []map[string]string, request RequestDetails) bool {
	for _, scope := range scopes {
		if (scope["method"] == "*" || scope["method"] == request.Method) &&
			(scope["uri"] == "*" || strings.HasPrefix(request.URI, scope["uri"])) {
			slog.Debug("Request authorized", "request", request)
			return true
		}
	}
	slog.Debug("Authorization failed for the request", "request", request)
	return false
}

// HashAndSalt hashes and salts the password using the bcrypt algorithm.
func HashAndSalt(pwd string) ([]byte, error) {
	bytePwd := []byte(pwd)

	// use bcrypt.GenerateFromPassword to hash and salt the password
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("failed to bcrypt hash password: %w", err)
	}

	return hash, nil
}

// pwdMatch checks if the hashed password matches the plain password.
func pwdMatch(hashed string, plain string) bool {
	hashedBytes := []byte(hashed)
	plainBytes := []byte(plain)

	err := bcrypt.CompareHashAndPassword(hashedBytes, plainBytes)
	return err == nil
}
