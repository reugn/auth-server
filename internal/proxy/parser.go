package proxy

import (
	"net/http"

	"github.com/reugn/auth-server/internal/repository"
)

// RequestParser represents a request parser. It parses the authorization token and the request
// details from the original request.
type RequestParser interface {

	// ParseAuthorizationToken parses and returns an Authorization token from the original request.
	ParseAuthorizationToken(r *http.Request) string

	// ParseRequestDetails parses and returns a RequestDetails from the original request.
	ParseRequestDetails(r *http.Request) *repository.RequestDetails
}
