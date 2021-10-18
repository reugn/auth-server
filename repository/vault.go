package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

type vaultEnv struct {
	vaultAddr        string
	vaultToken       string
	basicAuthKey     string
	authorizationKey string
}

// VaultRepository implements the Repository interface backed by HashiCorp Vault.
type VaultRepository struct {
	client *api.Client
	env    vaultEnv
}

func getVaultEnv() vaultEnv {
	// set defaults
	env := vaultEnv{
		vaultAddr:        "localhost:8200",
		basicAuthKey:     "secret/data/basic",
		authorizationKey: "secret/data/authorization",
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
		env.basicAuthKey = basicKey
	}
	authKey, ok := os.LookupEnv("AUTH_SERVER_VAULT_AUTHORIZATION_KEY")
	if ok {
		env.authorizationKey = authKey
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
	secret, err := vr.client.Logical().Read(vr.env.basicAuthKey)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	providedPwd, hashErr := hashAndSalt(password)
	if hashErr != nil {
		log.Println(err.Error())
		return nil
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Printf("VaultRepository read data error: %T %#v\n", secret.Data["data"], secret.Data["data"])
		return nil
	}

	if pass, ok := data["password"]; ok {
		if !pwdMatch(pass.(string), providedPwd) {
			return nil
		}
	} else {
		return nil
	}

	return &UserDetails{
		UserName: username,
		UserRole: data["role"].(UserRole),
	}
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (vr *VaultRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	secret, err := vr.client.Logical().Read(vr.env.authorizationKey)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Printf("VaultRepository read data error: %T %#v\n", secret.Data["data"], secret.Data["data"])
		return false
	}

	scopes := data[fmt.Sprint(userRole)].([]map[string]string)

	return isAuthorizedRequest(scopes, request)
}
