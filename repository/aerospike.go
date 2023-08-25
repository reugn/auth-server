package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"

	as "github.com/aerospike/aerospike-client-go/v6"
)

type aerospikeEnv struct {
	hostname         string
	port             int
	namespase        string
	setName          string
	basicAuthKey     string
	authorizationKey string
}

// AerospikeRepository implements the Repository interface backed by Aerospike Database.
type AerospikeRepository struct {
	client  *as.Client
	env     aerospikeEnv
	baseKey *as.Key
	authKey *as.Key
}

func getAerospikeEnv() aerospikeEnv {
	// set defaults
	env := aerospikeEnv{
		hostname:         "localhost",
		port:             3000,
		namespase:        "test",
		setName:          "auth",
		basicAuthKey:     "basic",
		authorizationKey: "authorization",
	}

	hostname, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_HOST")
	if ok {
		env.hostname = hostname
	}
	port, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_PORT")
	if ok {
		iport, err := strconv.Atoi(port)
		if err == nil {
			env.port = iport
		}
	}
	namespace, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_NAMESPACE")
	if ok {
		env.namespase = namespace
	}
	setName, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_SETNAME")
	if ok {
		env.setName = setName
	}
	basicKey, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_BASIC_KEY")
	if ok {
		env.basicAuthKey = basicKey
	}
	authKey, ok := os.LookupEnv("AUTH_SERVER_AEROSPIKE_AUTHORIZATION_KEY")
	if ok {
		env.authorizationKey = authKey
	}

	return env
}

// NewAerospikeRepositoryFromEnv returns a new instance of AerospikeRepository using env configuration.
func NewAerospikeRepositoryFromEnv() (*AerospikeRepository, error) {
	env := getAerospikeEnv()
	client, err := as.NewClient(env.hostname, env.port)
	if err != nil {
		return nil, err
	}

	baseKey, err := as.NewKey(env.namespase, env.setName, env.basicAuthKey)
	if err != nil {
		return nil, err
	}

	authKey, err := as.NewKey(env.namespase, env.setName, env.authorizationKey)
	if err != nil {
		return nil, err
	}

	return &AerospikeRepository{
		client:  client,
		env:     env,
		baseKey: baseKey,
		authKey: authKey,
	}, nil
}

// AuthenticateBasic validates the basic username and password before issuing a JWT.
// Uses the bcrypt password-hashing function to validate the password.
func (aero *AerospikeRepository) AuthenticateBasic(username string, password string) *UserDetails {
	record, err := aero.client.Get(nil, aero.baseKey, username)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Bin(user1: {username: user1, password: sha256, role: 1})
	userBin := record.Bins[username].(map[string]interface{})
	if hashed, ok := userBin["password"].(string); ok {
		if !pwdMatch(hashed, password) {
			return nil
		}
	} else {
		return nil
	}

	return &UserDetails{
		UserName: username,
		UserRole: userBin["role"].(UserRole),
	}
}

// AuthorizeRequest checks if the role has permissions to access the endpoint.
func (aero *AerospikeRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	record, err := aero.client.Get(nil, aero.authKey, fmt.Sprint(userRole))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// Bin(1: [{method: GET, uri: /health}])
	scopes := record.Bins[fmt.Sprint(userRole)].([]map[string]string)

	return isAuthorizedRequest(scopes, request)
}
