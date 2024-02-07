package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/reugn/auth-server/auth"
	"github.com/reugn/auth-server/proxy"
	"github.com/reugn/auth-server/repository"
	"github.com/reugn/auth-server/utils"
)

var rateLimiter = NewIPRateLimiter(1, 10)

var ipsWhiteList = map[string]struct{}{
	"127.0.0.1": {},
}

// HTTPServer is the authentication HTTP server wrapper.
type HTTPServer struct {
	addr         string
	parser       proxy.RequestParser
	repo         repository.Repository
	jwtGenerator *auth.JWTGenerator
	jwtValidtor  *auth.JWTValidator
}

// NewHTTPServer returns a new instance of HTTPServer.
func NewHTTPServer(host string, port int, keys *auth.Keys) *HTTPServer {
	addr := host + ":" + strconv.Itoa(port)
	repository := parseRepo()
	generator := auth.NewJWTGenerator(keys, parseAlgo())
	validator := auth.NewJWTValidator(keys, repository)

	return &HTTPServer{
		addr:         addr,
		parser:       parseProxy(),
		repo:         repository,
		jwtGenerator: generator,
		jwtValidtor:  validator,
	}
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, ok := ipsWhiteList[ip]
		if !ok {
			limiter := rateLimiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (ws *HTTPServer) start() {
	mux := http.NewServeMux()

	// root route
	mux.HandleFunc("/", rootActionHandler)

	// health route
	mux.HandleFunc("/health", healthActionHandler)

	// readiness route
	mux.HandleFunc("/ready", readyActionHandler)

	// version route
	mux.HandleFunc("/version", versionActionHandler)

	// token issuing route
	// uses basic authentication
	mux.HandleFunc("/token", ws.tokenActionHandler)

	// authorization route
	// validates bearer JWT
	mux.HandleFunc("/auth", ws.authActionHandler)

	err := http.ListenAndServe(ws.addr, rateLimiterMiddleware(mux))
	utils.Check(err)
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

func versionActionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, authServerVersion)
}

func (ws *HTTPServer) tokenActionHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDetails := ws.repo.AuthenticateBasic(user, pass)
	if userDetails == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := ws.jwtGenerator.Generate(userDetails.UserName, userDetails.UserRole)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", token.Marshal())
}

func (ws *HTTPServer) authActionHandler(w http.ResponseWriter, r *http.Request) {
	requestDetails := ws.parser.ParseRequestDetails(r)
	auth := ws.parser.ParseAuthorizationToken(r)

	if !ws.jwtValidtor.Authorize(auth, requestDetails) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
