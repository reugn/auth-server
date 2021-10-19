package repository

import (
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/vault/api"
)

type vaultEnv struct {
	vaultAddr              string
	vaultToken             string
	basicAuthKeyPrefix     string
	authorizationKeyPrefix string
}

// VaultRepository implements the Repository interface backed by HashiCorp Vault.
type VaultRepository struct {
	client *api.Client
	env    vaultEnv
}

func getVaultEnv() vaultEnv {
	// set defaults
	env := vaultEnv{
		vaultAddr:              "localhost:8200",
		basicAuthKeyPrefix:     "secret/basic",
		authorizationKeyPrefix: "secret/authorization",
	}

	vaultAddr, ok := os.LookupEnv("AUTH_SERVER_VAULT_ADDR")
	if ok {
		env.vaultAddr = vaultAddr
	}
	vaultToken, ok := os.LookupEnv("AUTH_SERVER_VAULT_TOKEN")
	if ok {
		env.vaultToken = vaultToken
	}
	basicKey, ok := os.LookupEnv("AUTH_SERVER_VAULT_BASIC_KEY")
	if ok {
		env.basicAuthKeyPrefix = basicKey
	}
	authKey, ok := os.LookupEnv("AUTH_SERVER_VAULT_AUTHORIZATION_KEY")
	if ok {
		env.authorizationKeyPrefix = authKey
	}

	return env
}

// NewVaultRepositoryFromEnv returns a new instance of VaultRepository using env configuration.
func NewVaultRepositoryFromEnv() (*VaultRepository, error) {
	env := getVaultEnv()
	config := &api.Config{
		Address: env.vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(env.vaultToken)

	return &VaultRepository{
		client: client,
		env:    env,
	}, nil
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
// Uses the bcrypt password-hashing function to validate the password.
func (vr *VaultRepository) AuthenticateBasic(username string, password string) *UserDetails {
	secret, err := vr.client.Logical().Read(vr.env.basicAuthKeyPrefix + "/" + username)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if hashed, ok := secret.Data["password"].(string); ok {
		if !pwdMatch(hashed, password) {
			return nil
		}
	} else {
		return nil
	}

	return &UserDetails{
		UserName: username,
		UserRole: secret.Data["role"].(UserRole),
	}
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (vr *VaultRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	secret, err := vr.client.Logical().Read(vr.env.authorizationKeyPrefix + "/" + strconv.Itoa(int(userRole)))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	scopes, ok := secret.Data["scopes"].([]map[string]string)
	if !ok {
		log.Printf("VaultRepository: error on reading scopes for %d", userRole)
		return false
	}

	return isAuthorizedRequest(scopes, request)
}
