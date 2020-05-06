package repository

import (
	"log"
	"os"
	"strconv"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/reugn/auth-server/utils"
)

type aerospikeEnv struct {
	hostname         string
	port             int
	namespase        string
	setName          string
	basicAuthKey     string
	authorizationKey string
}

// Aerospike repository implementation
type AerospikeRepository struct {
	client  *as.Client
	env     aerospikeEnv
	baseKey *as.Key
	authKey *as.Key
}

func getAerospikeEnv() aerospikeEnv {
	// set defaults
	env := aerospikeEnv{"localhost", 3000, "auth", "", "basic", "authorization"}

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

// NewAerospikeRepositoryFromEnv returns a new instance of AerospikeRepository using env configuration
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

// AuthenticateBasic validates basic username and password before issuing a JWT
func (aero *AerospikeRepository) AuthenticateBasic(username string, password string) *UserDetails {
	record, err := aero.client.Get(nil, aero.baseKey, username)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	sha256pwd := utils.Sha256(password)
	// Bin(user1: {username: user1, password: sha256, role: 1})
	userBin := record.Bins[username].(map[string]interface{})
	if pass, ok := userBin["password"]; ok {
		if pass.(string) != sha256pwd {
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

// AuthorizeRequest checks if the role has permissions to access the endpoint
func (aero *AerospikeRepository) AuthorizeRequest(userRole UserRole, request RequestDetails) bool {
	record, err := aero.client.Get(nil, aero.baseKey, string(userRole))
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	// Bin(1: [{method: GET, uri: /health}])
	permBin := record.Bins[string(userRole)].([]map[string]string)

	return isAuthorizedRequest(permBin, request)
}

func isAuthorizedRequest(perms []map[string]string, request RequestDetails) bool {
	for _, details := range perms {
		if (details["method"] == "*" || details["method"] == request.Method) &&
			(details["uri"] == "*" || details["uri"] == request.URI) {
			return true
		}
	}
	return false
}
