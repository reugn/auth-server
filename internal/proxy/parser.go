package proxy

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/reugn/auth-server/internal/repository"
)

// bearerPrefix is the prefix for the Authorization header.
const bearerPrefix = "Bearer "

// RequestParser represents a request parser. It parses the authorization token and the request
// details from the original request.
type RequestParser interface {
	// ParseAuthorizationToken parses and returns an Authorization token from the original request.
	ParseAuthorizationToken(r *http.Request) string
	// ParseRequestDetails parses and returns a RequestDetails from the original request.
	ParseRequestDetails(r *http.Request) *repository.RequestDetails
}

// HeaderConfig defines which headers to check for method and URI.
type HeaderConfig struct {
	// MethodHeaders are headers to check for the original HTTP method (in priority order).
	MethodHeaders []string
	// URIHeaders are headers to check for the original URI (in priority order).
	URIHeaders []string
	// PrefixHeader is an optional header for path prefix reconstruction.
	PrefixHeader string
}

// DefaultHeaderConfigs provides pre-configured header settings for common proxies.
var DefaultHeaderConfigs = map[string]HeaderConfig{
	"direct": {
		MethodHeaders: nil,
		URIHeaders:    nil,
		PrefixHeader:  "",
	},
	"nginx": {
		MethodHeaders: []string{"X-Forwarded-Method"},
		URIHeaders:    []string{"X-Forwarded-Uri"},
		PrefixHeader:  "",
	},
	"traefik": {
		MethodHeaders: []string{"X-Forwarded-Method"},
		URIHeaders:    []string{"X-Forwarded-Uri"},
		PrefixHeader:  "X-Forwarded-Prefix",
	},
	"envoy": {
		MethodHeaders: []string{"X-Original-Method", "X-Envoy-Original-Method", "X-Forwarded-Method"},
		URIHeaders:    []string{"X-Original-Uri", "X-Envoy-Original-Path", "X-Forwarded-Uri"},
		PrefixHeader:  "",
	},
	"haproxy": {
		MethodHeaders: []string{"X-Forwarded-Method"},
		URIHeaders:    []string{"X-Forwarded-Uri", "X-Original-URI"},
		PrefixHeader:  "",
	},
	"kong": {
		MethodHeaders: []string{"X-Forwarded-Method"},
		URIHeaders:    []string{"X-Forwarded-Path", "X-Forwarded-Uri"},
		PrefixHeader:  "X-Forwarded-Prefix",
	},
}

// Parser is a configurable request parser that works with any proxy.
type Parser struct {
	config HeaderConfig
}

var _ RequestParser = (*Parser)(nil)

// NewParser returns a Parser configured for a known proxy.
// Falls back to direct (no forwarded headers) if the proxy is unknown.
func NewParser(proxyName string) *Parser {
	config, ok := DefaultHeaderConfigs[strings.ToLower(proxyName)]
	if !ok {
		slog.Warn("Unknown proxy provider, using direct mode", "proxy", proxyName)
		config = DefaultHeaderConfigs["direct"]
	}
	return &Parser{config: config}
}

// NewParserWithConfig returns a new Parser with a custom header configuration.
func NewParserWithConfig(config HeaderConfig) *Parser {
	return &Parser{config: config}
}

// ParseAuthorizationToken parses and returns an Authorization Bearer token from the request.
// It returns an empty string if the authorization header is not present or if the authorization
// header is not in the correct format (Bearer <token>).
func (p *Parser) ParseAuthorizationToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		slog.Debug("Invalid Authorization header format")
		return ""
	}

	return strings.TrimSpace(authHeader[len(bearerPrefix):])
}

// ParseRequestDetails parses and returns RequestDetails from the request.
func (p *Parser) ParseRequestDetails(r *http.Request) *repository.RequestDetails {
	return &repository.RequestDetails{
		Method: p.parseMethod(r),
		URI:    p.parseURI(r),
	}
}

// parseMethod extracts the HTTP method from forwarded headers or the request itself.
func (p *Parser) parseMethod(r *http.Request) string {
	if len(p.config.MethodHeaders) > 0 {
		if method := getFirstHeader(r, p.config.MethodHeaders...); method != "" {
			return method
		}
	}
	return r.Method
}

// parseURI extracts the URI from forwarded headers or the request itself.
func (p *Parser) parseURI(r *http.Request) string {
	// First, try configured URI headers
	if len(p.config.URIHeaders) > 0 {
		if uri := getFirstHeader(r, p.config.URIHeaders...); uri != "" {
			return uri
		}
	}

	// Handle prefix reconstruction if configured
	if p.config.PrefixHeader != "" {
		if prefix := r.Header.Get(p.config.PrefixHeader); prefix != "" {
			return reconstructURI(prefix, r.URL.RequestURI())
		}
	}

	return r.URL.RequestURI()
}

// getFirstHeader returns the first non-empty header value from the provided headers.
func getFirstHeader(r *http.Request, headers ...string) string {
	for _, header := range headers {
		if value := r.Header.Get(header); value != "" {
			return value
		}
	}
	return ""
}

// reconstructURI combines a prefix with the current path.
func reconstructURI(prefix, currentPath string) string {
	prefix = strings.TrimSuffix(prefix, "/")
	currentPath = strings.TrimPrefix(currentPath, "/")

	if prefix == "" {
		return "/" + currentPath
	}
	if currentPath == "" {
		return prefix
	}
	return prefix + "/" + currentPath
}
