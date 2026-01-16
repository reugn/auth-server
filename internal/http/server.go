package http

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/reugn/auth-server/internal/auth"
	"github.com/reugn/auth-server/internal/config"
	"github.com/reugn/auth-server/internal/proxy"
	"github.com/reugn/auth-server/internal/repository"
	"github.com/reugn/auth-server/internal/util/httputil"
	"golang.org/x/time/rate"
)

// Server represents the entry point to interact with the service via HTTP requests.
type Server struct {
	address      string
	version      string
	parser       proxy.RequestParser
	repository   repository.Repository
	rateLimiter  *IPRateLimiter
	ipWhiteList  *IPWhiteList
	jwtGenerator *auth.JWTGenerator
	jwtValidator *auth.JWTValidator
}

// NewServer returns a new instance of Server.
// It initializes the server with the given version, keys, and configuration.
func NewServer(version string, keys *auth.Keys, config *config.Service) (*Server, error) {
	address := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)
	repo, err := config.Repository()
	if err != nil {
		return nil, err
	}

	// get the signing method
	signingMethod, err := config.SigningMethodRSA()
	if err != nil {
		return nil, err
	}

	// get the IP white list
	ipWhiteList, err := NewIPWhiteList(config.HTTP.Rate.WhiteList)
	if err != nil {
		return nil, err
	}

	return &Server{
		address:      address,
		version:      version,
		parser:       config.RequestParser(),
		repository:   repo,
		rateLimiter:  NewIPRateLimiter(rate.Limit(config.HTTP.Rate.Tps), config.HTTP.Rate.Size),
		ipWhiteList:  ipWhiteList,
		jwtGenerator: auth.NewJWTGenerator(keys, signingMethod),
		jwtValidator: auth.NewJWTValidator(keys, repo),
	}, nil
}

func (ws *Server) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the IP from the remote address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		// check if the IP is in the white list
		if !ws.ipWhiteList.isAllowed(ip) {
			limiter := ws.rateLimiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests),
					http.StatusTooManyRequests)
				return
			}
		}

		// serve the request
		next.ServeHTTP(w, r)
	})
}

// Start initiates the HTTP server.
func (ws *Server) Start() error {
	mux := http.NewServeMux()

	// root route
	mux.HandleFunc("/", rootActionHandler)

	// health route
	mux.HandleFunc("/health", healthActionHandler)

	// readiness route
	mux.HandleFunc("/ready", readyActionHandler)

	// version route
	mux.HandleFunc("/version", ws.versionActionHandler)

	// token issuing route, requires basic authentication
	mux.HandleFunc("/token", ws.tokenActionHandler)

	// authorization route, requires a JSON Web Token
	mux.HandleFunc("/auth", ws.authActionHandler)
	// prefix match for Envoy ext_authz (handles /auth/health, /auth/ready, etc.)
	mux.HandleFunc("/auth/", ws.authActionHandler)

	//nolint:gosec
	return http.ListenAndServe(ws.address, ws.rateLimiterMiddleware(mux))
}

func rootActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, "")
}

func healthActionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func readyActionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func (ws *Server) versionActionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, ws.version)
}

func (ws *Server) tokenActionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Token generation request")
	// get the basic authentication credentials
	user, pass, ok := r.BasicAuth()
	if !ok {
		slog.Info("Failed to get basic authentication credentials", httputil.LogAttrs(r)...)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// authenticate the user
	userDetails := ws.repository.AuthenticateBasic(user, pass)
	if userDetails == nil {
		slog.Info("Failed to authenticate user", httputil.LogAttrs(r,
			slog.String("username", user))...)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	slog.Debug("User authenticated", httputil.LogAttrs(r,
		slog.String("username", user), slog.String("role", string(userDetails.UserRole)))...)

	// generate the access token
	accessToken, err := ws.jwtGenerator.Generate(userDetails.UserName, userDetails.UserRole)
	if err != nil {
		slog.Error("Failed to generate access token", httputil.LogAttrs(r,
			slog.String("username", user), slog.Any("error", err))...)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// marshal the access token
	marshalled, err := accessToken.Marshal()
	if err != nil {
		slog.Error("Failed to marshal access token", httputil.LogAttrs(r,
			slog.String("username", user), slog.Any("error", err))...)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write the access token to the response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, marshalled)
}

func (ws *Server) authActionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Token authorization request")
	requestDetails := ws.parser.ParseRequestDetails(r)
	authToken := ws.parser.ParseAuthorizationToken(r)

	// authorize the request using the JWT validator and write the response
	if !ws.jwtValidator.Authorize(authToken, requestDetails) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// success response
	w.WriteHeader(http.StatusOK)
}
