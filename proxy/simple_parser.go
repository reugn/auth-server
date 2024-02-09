package proxy

import (
	"net/http"
	"strings"

	"github.com/reugn/auth-server/repository"
)

// SimpleParser implements the RequestParser interface.
type SimpleParser struct{}

var _ RequestParser = (*SimpleParser)(nil)

// NewSimpleParser returns a new SimpleParser.
func NewSimpleParser() *SimpleParser {
	return &SimpleParser{}
}

// ParseAuthorizationToken parses and returns an Authorization Bearer token from the original request.
func (sp *SimpleParser) ParseAuthorizationToken(r *http.Request) string {
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

// ParseRequestDetails parses and returns a RequestDetails from the original request.
func (sp *SimpleParser) ParseRequestDetails(r *http.Request) *repository.RequestDetails {
	return &repository.RequestDetails{
		Method: r.Method,
		URI:    r.URL.RequestURI(),
	}
}
