package repository

// User model for the internal Local repository use.
type User struct {
	password string
	role     UserRole
}

// Local implements the Repository interface backed by in-memory permission details.
// Use it for test purposes only.
type Local struct {
	Users       map[string]User
	Permissions map[UserRole][]RequestDetails
}

// NewLocalRepo returns a new instance of the Local repo.
func NewLocalRepo() *Local {
	users := make(map[string]User)
	users["admin"] = User{"1234", 1}

	perms := make(map[UserRole][]RequestDetails)
	requestDetails := []RequestDetails{
		{
			"GET",
			"/dashboard",
		},
		{
			"POST",
			"/auth",
		},
		{
			"GET",
			"/health",
		},
	}
	perms[1] = requestDetails
	return &Local{users, perms}
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
func (local *Local) AuthenticateBasic(username string, password string) *UserDetails {
	if user, ok := local.Users[username]; ok {
		if user.password == password {
			return &UserDetails{
				username,
				user.role,
			}
		}
	}
	return nil
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (local *Local) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	if perms, ok := local.Permissions[userRole]; ok {
		if containsRequestDetails(perms, request) {
			return true
		}
	}
	return false
}

func containsRequestDetails(details []RequestDetails, rd RequestDetails) bool {
	for _, detail := range details {
		if detail == rd {
			return true
		}
	}
	return false
}
