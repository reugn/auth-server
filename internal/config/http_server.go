package config

import (
	"errors"
	"fmt"
)

// HTTP contains HTTP server configuration properties.
type HTTP struct {
	// The address to listen on.
	Host string `yaml:"host,omitempty" json:"host,omitempty"`
	// The port to listen on.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`
	// Rate limiter configuration.
	Rate RateLimiter `yaml:"rate,omitempty" json:"rate,omitempty"`
}

// RateLimiter contains rate limiter configuration properties.
type RateLimiter struct {
	// Rate limiter tokens per second threshold.
	Tps int `yaml:"tps,omitempty" json:"tps,omitempty"`
	// Rate limiter token bucket size (bursts threshold).
	Size int `yaml:"size,omitempty" json:"size,omitempty"`
	// A list of IP addresses to exclude from rate limiting.
	WhiteList []string `yaml:"white-list,omitempty" json:"white-list,omitempty"`
}

func (c *RateLimiter) validate() error {
	if c == nil {
		return errors.New("rate limiter config is nil")
	}
	if c.Tps < 1 {
		return fmt.Errorf("invalid rate tps: %d", c.Tps)
	}
	if c.Size < 1 {
		return fmt.Errorf("invalid rate size: %d", c.Size)
	}
	return nil
}

// NewHTTPDefault returns a new HTTP config with default values.
func NewHTTPDefault() *HTTP {
	return &HTTP{
		Host: "0.0.0.0",
		Port: 8080,
		Rate: RateLimiter{
			Tps:  1024,
			Size: 1024,
		},
	}
}

// validate validates the HTTP configuration.
func (c *HTTP) validate() error {
	if c == nil {
		return errors.New("http config is nil")
	}
	if c.Host == "" {
		return errors.New("host is not specified")
	}
	if c.Port < 1 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if err := c.Rate.validate(); err != nil {
		return err
	}
	return nil
}
