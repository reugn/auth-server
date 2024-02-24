package auth

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/config"
)

// Keys represents a container for the private and public keys.
type Keys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewKeys returns a new instance of Keys.
func NewKeys(config *config.Secret) (*Keys, error) {
	return NewKeysFromFile(config.Private, config.Public)
}

// NewKeysFromFile creates and returns a new instance of Keys from the files
// containing the secrets information.
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
