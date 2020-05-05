package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/reugn/auth-server/auth"
	"github.com/reugn/auth-server/proxy"
	"github.com/reugn/auth-server/repository"
)

// HTTPServer is an auth http server wrapper
type HTTPServer struct {
	addr         string
	parser       proxy.RequestParser
	repo         repository.Repository
	jwtGenerator *auth.JWTGenerator
	jwtValidtor  *auth.JWTValidator
}

// NewHTTPServer returns a new instance of HTTPServer
func NewHTTPServer(host string, port int, repo repository.Repository, keys *auth.Keys) *HTTPServer {
	addr := host + ":" + strconv.Itoa(port)
	generator := auth.NewJWTGenerator(keys)
	validator := auth.NewJWTValidator(keys, repo)

	return &HTTPServer{
		addr,
		parseProxy(),
		repo,
		generator,
		validator,
	}
}

func (ws *HTTPServer) start() {
	// root route
	http.HandleFunc("/", rootActionHandler)

	// health route
	http.HandleFunc("/health", healthActionHandler)

	// readiness route
	http.HandleFunc("/ready", readyActionHandler)

	// token issuing route
	// uses basic authentication
	http.HandleFunc("/token", ws.tokenActionHandler)

	// authorization route
	// validates bearer JWT
	http.HandleFunc("/auth", ws.authActionHandler)

	http.ListenAndServe(ws.addr, nil)
}

func rootActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, "")
}

func healthActionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func readyActionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
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
