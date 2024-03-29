package repository

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/vault/api"
	"github.com/reugn/auth-server/internal/util/env"
)

// Environment variables to configure VaultRepository.
const (
	envVaultAddr     = "AUTH_SERVER_VAULT_ADDR"
	envVaultToken    = "AUTH_SERVER_VAULT_TOKEN"
	envVaultBasicKey = "AUTH_SERVER_VAULT_BASIC_KEY"
	envVaultAuthKey  = "AUTH_SERVER_VAULT_AUTHORIZATION_KEY"
)

// vaultConfig contains VaultRepository configuration properties.
type vaultConfig struct {
	vaultAddr              string
	vaultToken             string
	basicAuthKeyPrefix     string
	authorizationKeyPrefix string
}

// VaultRepository implements the Repository interface using HashiCorp Vault
// as the storage backend.
type VaultRepository struct {
	client *api.Client
	config vaultConfig
}

var _ Repository = (*VaultRepository)(nil)

func getVaultConfig() vaultConfig {
	// set defaults
	config := vaultConfig{
		vaultAddr:              "localhost:8200",
		basicAuthKeyPrefix:     "secret/basic",
		authorizationKeyPrefix: "secret/authorization",
	}

	// read configuration from environment variables
	env.ReadString(&config.vaultAddr, envVaultAddr)
	env.ReadString(&config.vaultToken, envVaultToken)
	env.ReadString(&config.basicAuthKeyPrefix, envVaultBasicKey)
	env.ReadString(&config.authorizationKeyPrefix, envVaultAuthKey)

	return config
}

// NewVault returns a new VaultRepository using environment variables for configuration.
func NewVault() (*VaultRepository, error) {
	config := getVaultConfig() // read configuration
	apiConfig := &api.Config{
		Address: config.vaultAddr,
	}
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}
	client.SetToken(config.vaultToken)

	return &VaultRepository{
		client: client,
		config: config,
	}, nil
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
// It uses the bcrypt password-hashing function to validate the password.
func (vr *VaultRepository) AuthenticateBasic(username string, password string) *UserDetails {
	path := fmt.Sprintf("%s/%s", vr.config.basicAuthKeyPrefix, username)
	secret, err := vr.client.Logical().Read(path)
	if err != nil {
		slog.Error("Failed to read path", "path", path, "err", err)
		return nil
	}

	hashed, ok := secret.Data["password"].(string)
	if !ok || !pwdMatch(hashed, password) {
		slog.Debug("Failed to authenticate", "user", username)
		return nil
	}

	return &UserDetails{
		UserName: username,
		UserRole: secret.Data["role"].(UserRole),
	}
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (vr *VaultRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	path := fmt.Sprintf("%s/%s", vr.config.authorizationKeyPrefix, userRole)
	secret, err := vr.client.Logical().Read(path)
	if err != nil {
		slog.Error("Failed to read path", "path", path, "err", err)
		return false
	}

	scopes, ok := secret.Data["scopes"].([]map[string]string)
	if !ok {
		slog.Error("Error reading scopes", "role", userRole)
		return false
	}

	return isAuthorizedRequest(scopes, request)
}
