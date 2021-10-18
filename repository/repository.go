package repository

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// UserRole represents a user role.
type UserRole int

// UserDetails represents user details.
type UserDetails struct {
	UserName string
	UserRole UserRole
}

// RequestDetails represents request details.
type RequestDetails struct {
	Method string
	URI    string
}

// Repository is the interface to a custom authentication/authorization backend facade.
type Repository interface {

	// AuthenticateBasic validates the basic username and password before issuing a JWT.
	AuthenticateBasic(username string, password string) *UserDetails

	// AuthorizeRequest checks if the role has permissions to access the endpoint.
	AuthorizeRequest(userRole UserRole, request RequestDetails) bool
}

func isAuthorizedRequest(scopes []map[string]string, request RequestDetails) bool {
	for _, scope := range scopes {
		if (scope["method"] == "*" || scope["method"] == request.Method) &&
			(scope["uri"] == "*" || strings.HasPrefix(request.URI, scope["uri"])) {
			return true
		}
	}
	return false
}

func hashAndSalt(pwd string) ([]byte, error) {
	bytePwd := []byte(pwd)

	// Use bcrypt.GenerateFromPassword to hash and salt the password.
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func pwdMatch(stored string, provided []byte) bool {
	byteStored := []byte(stored)

	err := bcrypt.CompareHashAndPassword(byteStored, provided)
	if err != nil {
		return false
	}

	return true
}
