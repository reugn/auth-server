package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

// Keys container for private and public keys
type Keys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewKeys returns a new instance of Keys
func NewKeys() (*Keys, error) {
	return NewKeysFromFile(env.privateKeyPath, env.publicKeyPath)
}

// NewKeysFromFile returns a new instance of Keys using file pathes
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

// NewKeysFromPem returns a new instance of Keys using pem byte arrays
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
		pem, err := ioutil.ReadFile(*privateKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPrivateKeyFromPEM(pem)
	} else if pem != nil {
		return jwt.ParseRSAPrivateKeyFromPEM(pem)
	} else {
		return nil, errors.New("parsePrivateKey nil parameters")
	}
}

func parsePublicKey(publicKeyPath *string, pem []byte) (*rsa.PublicKey, error) {
	if publicKeyPath != nil {
		pem, err := ioutil.ReadFile(*publicKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	} else if pem != nil {
		return jwt.ParseRSAPublicKeyFromPEM(pem)
	} else {
		return nil, errors.New("parsePublicKey nil parameters")
	}
}
