package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/proxy"
	"github.com/reugn/auth-server/internal/repository"
)

const (
	signingMethodRS256 = "RS256"
	signingMethodRS384 = "RS384"
	signingMethodRS512 = "RS512"
)

var validSigningMethods = []string{signingMethodRS256, signingMethodRS384, signingMethodRS512}

// Service contains the entire service configuration.
type Service struct {
	SigningMethod      string  `yaml:"signing-method,omitempty" json:"signing-method,omitempty"`
	ProxyProvider      string  `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	RepositoryProvider string  `yaml:"repository,omitempty" json:"repository,omitempty"`
	HTTP               *HTTP   `yaml:"http,omitempty" json:"http,omitempty"`
	Logger             *Logger `yaml:"logger,omitempty" json:"logger,omitempty"`
}

// NewServiceDefault returns a new Service config with default values.
func NewServiceDefault() *Service {
	return &Service{
		SigningMethod:      signingMethodRS256,
		ProxyProvider:      "simple",
		RepositoryProvider: "local",
		HTTP:               NewHTTPDefault(),
		Logger:             NewLoggerDefault(),
	}
}

func (c *Service) SigningMethodRSA() (*jwt.SigningMethodRSA, error) {
	var signingMethodRSA *jwt.SigningMethodRSA
	switch strings.ToUpper(c.SigningMethod) {
	case signingMethodRS256:
		signingMethodRSA = jwt.SigningMethodRS256
	case signingMethodRS384:
		signingMethodRSA = jwt.SigningMethodRS384
	case signingMethodRS512:
		signingMethodRSA = jwt.SigningMethodRS512
	default:
		return nil, fmt.Errorf("unsupported signing method: %s", c.SigningMethod)
	}
	return signingMethodRSA, nil
}

func (c *Service) RequestParser() (proxy.RequestParser, error) {
	var parser proxy.RequestParser
	switch strings.ToLower(c.ProxyProvider) {
	case "simple":
		parser = proxy.NewSimpleParser()
	case "traefik":
		parser = proxy.NewTraefikParser()
	default:
		return nil, fmt.Errorf("unsupported proxy provider: %s", c.ProxyProvider)
	}
	return parser, nil
}

func (c *Service) Repository() (repository.Repository, error) {
	switch strings.ToLower(c.RepositoryProvider) {
	case "local":
		return repository.NewLocal()
	case "aerospike":
		return repository.NewAerospike()
	case "vault":
		return repository.NewVault()
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", c.RepositoryProvider)
	}
}

// Validate validates the service configuration.
func (c *Service) Validate() error {
	if c == nil {
		return errors.New("service config is nil")
	}
	if !slices.Contains(validSigningMethods, strings.ToUpper(c.SigningMethod)) {
		return fmt.Errorf("invalid signing method: %s", c.SigningMethod)
	}
	if c.ProxyProvider == "" {
		return errors.New("proxy provider is not specified")
	}
	if c.RepositoryProvider == "" {
		return errors.New("repository provider is not specified")
	}
	if err := c.HTTP.validate(); err != nil {
		return err
	}
	if err := c.Logger.validate(); err != nil {
		return err
	}
	return nil
}

// String returns string representation of the service configuration.
func (c *Service) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
