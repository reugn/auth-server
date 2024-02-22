package config

import (
	"errors"
)

// Secret holds the configuration for secret keys.
type Secret struct {
	// Private denotes the path to the private key.
	Private string `yaml:"private-path,omitempty" json:"private-path,omitempty"`
	// Public denotes the path to the public key.
	Public string `yaml:"public-path,omitempty" json:"public-path,omitempty"`
}

// NewSecretDefault returns a new Secret with default values.
func NewSecretDefault() *Secret {
	return &Secret{
		Private: "secrets/privkey.pem",
		Public:  "secrets/cert.pem",
	}
}

// validate validates the Secret configuration properties.
func (s *Secret) validate() error {
	if s == nil {
		return errors.New("secret config is nil")
	}
	if s.Private == "" {
		return errors.New("private key path is not specified")
	}
	if s.Public == "" {
		return errors.New("public key path is not specified")
	}
	return nil
}
