package repository

// UserRole model
type UserRole int

// UserDetails model
type UserDetails struct {
	UserName string
	UserRole UserRole
}

// RequestDetails model
type RequestDetails struct {
	Method string
	URI    string
}

// Repository is the interface to a custom authentication/authorization backend facade
type Repository interface {

	// AuthenticateBasic validates basic username and password before issuing a JWT
	AuthenticateBasic(username string, password string) *UserDetails

	// AuthorizeRequest checks if the role has permissions to access the endpoint
	AuthorizeRequest(userRole UserRole, request RequestDetails) bool
}
