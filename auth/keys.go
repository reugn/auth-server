package auth

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/util/env"
)

const (
	envPrivateKeyPath = "AUTH_SERVER_PRIVATE_KEY_PATH"
	envPublicKeyPath  = "AUTH_SERVER_PUBLIC_KEY_PATH"

	defaultPrivateKeyPath = "secrets/privkey.pem"
	defaultPublicKeyPath  = "secrets/cert.pem"
)

// Keys represents a container for the private and public keys.
type Keys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewKeys returns a new instance of Keys.
func NewKeys() (*Keys, error) {
	privateKeyPath := defaultPrivateKeyPath
	env.ReadString(&privateKeyPath, envPrivateKeyPath)

	publicKeyPath := defaultPublicKeyPath
	env.ReadString(&publicKeyPath, envPublicKeyPath)

	return NewKeysFromFile(privateKeyPath, publicKeyPath)
}

// NewKeysFromFile creates and returns a new instance of Keys from the files.
func NewKeysFromFile(privateKeyPath string, publicKeyPath string) (*Keys, error) {
	priv, err := parsePrivateKey(&privateKeyPath, nil)
	if err != nil {
		return nil, err
	}

	pub, err := parsePublicKey(&publicKeyPath, nil)
	if err != nil {
		return nil, err
	}

	return &Keys{priv, pub}, nil
}

// NewKeysFromPem creates and returns a new instance of Keys from the pem byte arrays.
func NewKeysFromPem(privatePem []byte, publicPem []byte) (*Keys, error) {
	priv, err := parsePrivateKey(nil, privatePem)
	if err != nil {
		return nil, err
	}

	pub, err := parsePublicKey(nil, publicPem)
	if err != nil {
		return nil, err
	}

	return &Keys{priv, pub}, nil
}

func parsePrivateKey(privateKeyPath *string, pem []byte) (*rsa.PrivateKey, error) {
	if privateKeyPath != nil {
		pem, err := os.ReadFile(*privateKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPrivateKeyFromPEM(pem)
	} else if pem != nil {
		return jwt.ParseRSAPrivateKeyFromPEM(pem)
	}
	return nil, errors.New("parsePrivateKey nil parameters")
}

func parsePublicKey(publicKeyPath *string, pem []byte) (*rsa.PublicKey, error) {
	if publicKeyPath != nil {
		pem, err := os.ReadFile(*publicKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	} else if pem != nil {
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	}
	return nil, errors.New("parsePublicKey nil parameters")
}
