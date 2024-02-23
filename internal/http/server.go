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
func NewServer(version string, keys *auth.Keys, config *config.Service) (*Server, error) {
	address := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)
	repository, err := config.Repository()
	if err != nil {
		return nil, err
	}
	signingMethod, err := config.SigningMethodRSA()
	if err != nil {
		return nil, err
	}
	generator := auth.NewJWTGenerator(keys, signingMethod)
	validator := auth.NewJWTValidator(keys, repository)

	requestParser, err := config.RequestParser()
	if err != nil {
		return nil, err
	}
	ipWhiteList, err := NewIPWhiteList(config.HTTP.Rate.WhiteList)
	if err != nil {
		return nil, err
	}
	return &Server{
		address:      address,
		version:      version,
		parser:       requestParser,
		repository:   repository,
		rateLimiter:  NewIPRateLimiter(rate.Limit(config.HTTP.Rate.Tps), config.HTTP.Rate.Size),
		ipWhiteList:  ipWhiteList,
		jwtGenerator: generator,
		jwtValidator: validator,
	}, nil
}

func (ws *Server) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		if !ws.ipWhiteList.isAllowed(ip) {
			limiter := ws.rateLimiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests),
					http.StatusTooManyRequests)
				return
			}
		}
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
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDetails := ws.repository.AuthenticateBasic(user, pass)
	if userDetails == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken, err := ws.jwtGenerator.Generate(userDetails.UserName, userDetails.UserRole)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshalled, err := accessToken.Marshal()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", marshalled)
}

func (ws *Server) authActionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Token authorization request")
	requestDetails := ws.parser.ParseRequestDetails(r)
	authToken := ws.parser.ParseAuthorizationToken(r)

	if !ws.jwtValidator.Authorize(authToken, requestDetails) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
