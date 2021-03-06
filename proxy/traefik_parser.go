package proxy

import (
	"net/http"
	"strings"

	"github.com/reugn/auth-server/repository"
)

// TraefikParser implementation of RequestParser interface
type TraefikParser struct{}

// NewTraefikParser returns a new instance of TraefikParser
func NewTraefikParser() *TraefikParser {
	return &TraefikParser{}
}

// ParseAuthorizationToken parses and returns an Authorization Bearer token from the original request
func (tp *TraefikParser) ParseAuthorizationToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return authHeader
	}
	splitToken := strings.Split(authHeader, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}

// ParseRequestDetails parses and returns a RequestDetails from the original request
func (tp *TraefikParser) ParseRequestDetails(r *http.Request) *repository.RequestDetails {
	return &repository.RequestDetails{
		Method: r.Header.Get("X-Forwarded-Method"),
		URI:    r.Header.Get("X-Forwarded-Uri"),
	}
}
