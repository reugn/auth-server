package repository

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/reugn/auth-server/internal/util/env"
	"gopkg.in/yaml.v3"
)

const (
	EnvLocalConfigPath     = "AUTH_SERVER_LOCAL_CONFIG_PATH"
	DefaultLocalConfigPath = "../../config/local_repository_config.yml"
)

// AuthDetails contains authentication details for the user.
type AuthDetails struct {
	Password string   `yaml:"password"`
	Role     UserRole `yaml:"role"`
}

// Local implements the Repository interface by loading authentication details from
// a local configuration file.
type Local struct {
	Users map[string]AuthDetails        `yaml:"users"`
	Roles map[UserRole][]RequestDetails `yaml:"roles"`
}

var _ Repository = (*Local)(nil)

// NewLocal returns a new Local repository using an environment variable to
// read a custom path to the configuration file.
func NewLocal() (*Local, error) {
	configPath := DefaultLocalConfigPath
	env.ReadString(&configPath, EnvLocalConfigPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", configPath, err)
	}

	localRepository := &Local{}
	if err = yaml.Unmarshal(data, localRepository); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration data: %w", err)
	}

	return localRepository, nil
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
func (local *Local) AuthenticateBasic(username string, password string) *UserDetails {
	if authDetails, ok := local.Users[username]; ok {
		if authDetails.Password == password {
			return &UserDetails{
				UserName: username,
				UserRole: authDetails.Role,
			}
		}
	}
	return nil
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (local *Local) AuthorizeRequest(userRole UserRole, requestDetails RequestDetails) bool {
	if permissions, ok := local.Roles[userRole]; ok {
		if slices.Contains(permissions, requestDetails) {
			slog.Debug("Request authorized", "request", requestDetails)
			return true
		}
	}
	slog.Debug("Authorization failed for the request", "request", requestDetails)
	return false
}
