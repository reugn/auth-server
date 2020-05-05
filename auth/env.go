package auth

import (
	"os"
	"strconv"
	"time"
)

type environmentVariables struct {
	expireAfter    time.Duration
	privateKeyPath string
	publicKeyPath  string
}

var env environmentVariables = environmentVariables{time.Hour, "", ""}

func init() {
	// read environment variables
	tokenExpirationMilis, ok := os.LookupEnv("AUTH_SERVER_ACCESS_TOKEN_EXPIRATION_MILLIS")
	if ok {
		expireAfter, err := strconv.Atoi(tokenExpirationMilis)
		if err == nil {
			env.expireAfter = time.Duration(expireAfter) * time.Millisecond
		}
	}

	privateKeyPath, ok := os.LookupEnv("AUTH_SERVER_PRIVATE_KEY_PATH")
	if ok {
		env.privateKeyPath = privateKeyPath
	} else {
		env.privateKeyPath = "secrets/privkey.pem"
	}

	publicKeyPath, ok := os.LookupEnv("AUTH_SERVER_PUBLIC_KEY_PATH")
	if ok {
		env.publicKeyPath = publicKeyPath
	} else {
		env.publicKeyPath = "secrets/cert.pem"
	}
}
