package repository

import (
	"log/slog"

	as "github.com/aerospike/aerospike-client-go/v7"
	"github.com/reugn/auth-server/internal/util/env"
)

// Environment variables to configure AerospikeRepository.
const (
	envAerospikeHost      = "AUTH_SERVER_AEROSPIKE_HOST"
	envAerospikePort      = "AUTH_SERVER_AEROSPIKE_PORT"
	envAerospikeNamespace = "AUTH_SERVER_AEROSPIKE_NAMESPACE"
	envAerospikeSet       = "AUTH_SERVER_AEROSPIKE_SETNAME"
	envAerospikeBasicKey  = "AUTH_SERVER_AEROSPIKE_BASIC_KEY"
	envAerospikeAuthKey   = "AUTH_SERVER_AEROSPIKE_AUTHORIZATION_KEY"
)

// aerospikeConfig contains AerospikeRepository configuration properties.
type aerospikeConfig struct {
	hostname         string
	port             int
	namespase        string
	setName          string
	basicAuthKey     string
	authorizationKey string
}

// AerospikeRepository implements the Repository interface using Aerospike Database
// as the storage backend.
type AerospikeRepository struct {
	client  *as.Client
	config  aerospikeConfig
	baseKey *as.Key
	authKey *as.Key
}

var _ Repository = (*AerospikeRepository)(nil)

func getAerospikeConfig() aerospikeConfig {
	// set defaults
	config := aerospikeConfig{
		hostname:         "localhost",
		port:             3000,
		namespase:        "test",
		setName:          "auth",
		basicAuthKey:     "basic",
		authorizationKey: "authorization",
	}

	// read configuration from environment variables
	env.ReadString(&config.hostname, envAerospikeHost)
	env.ReadInt(&config.port, envAerospikePort)
	env.ReadString(&config.namespase, envAerospikeNamespace)
	env.ReadString(&config.setName, envAerospikeSet)
	env.ReadString(&config.basicAuthKey, envAerospikeBasicKey)
	env.ReadString(&config.authorizationKey, envAerospikeAuthKey)

	return config
}

// NewAerospike returns a new AerospikeRepository using environment variables for configuration.
func NewAerospike() (*AerospikeRepository, error) {
	config := getAerospikeConfig() // read configuration
	client, err := as.NewClient(config.hostname, config.port)
	if err != nil {
		return nil, err
	}
	baseKey, err := as.NewKey(config.namespase, config.setName, config.basicAuthKey)
	if err != nil {
		return nil, err
	}
	authKey, err := as.NewKey(config.namespase, config.setName, config.authorizationKey)
	if err != nil {
		return nil, err
	}

	return &AerospikeRepository{
		client:  client,
		config:  config,
		baseKey: baseKey,
		authKey: authKey,
	}, nil
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
// It uses the bcrypt password-hashing function to validate the password.
func (aero *AerospikeRepository) AuthenticateBasic(username string, password string) *UserDetails {
	record, err := aero.client.Get(nil, aero.baseKey, username)
	if err != nil {
		slog.Error("Failed to fetch record", "key", aero.baseKey, "err", err)
		return nil
	}

	// Bin(user1: {username: user1, password: sha256, role: admin})
	userBin := record.Bins[username].(map[string]interface{})
	hashed, ok := userBin["password"].(string)
	if !ok || !pwdMatch(hashed, password) {
		slog.Debug("Failed to authenticate", "user", username)
		return nil
	}

	return &UserDetails{
		UserName: username,
		UserRole: userBin["role"].(UserRole),
	}
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (aero *AerospikeRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	record, err := aero.client.Get(nil, aero.authKey, string(userRole))
	if err != nil {
		slog.Error("Failed to fetch record", "key", aero.authKey, "err", err)
		return false
	}
	// Bin(admin: [{method: GET, uri: /health}])
	scopes := record.Bins[string(userRole)].([]map[string]string)

	return isAuthorizedRequest(scopes, request)
}
