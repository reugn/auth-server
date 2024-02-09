package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/auth"
	"github.com/reugn/auth-server/proxy"
	"github.com/reugn/auth-server/repository"
	"github.com/reugn/auth-server/util"
)

const authServerVersion = "0.3.1"

var (
	serverHostParam = flag.String("host", "0.0.0.0", "Server host")
	serverPortParam = flag.Int("port", 8081, "Server port")
	algoParam       = flag.String("algo", "RS256", "JWT signing algorithm")
	proxyParam      = flag.String("proxy", "simple", "Proxy provider")
	repoParam       = flag.String("repo", "local", "Repository provider")
)

func main() {
	flag.Parse()

	// load ssl keys
	keys, err := auth.NewKeys()
	util.Check(err)

	// start http server
	server := NewHTTPServer(*serverHostParam, *serverPortParam, keys)
	server.start()
}

func parseAlgo() *jwt.SigningMethodRSA {
	var signingMethodRSA *jwt.SigningMethodRSA
	switch strings.ToUpper(*algoParam) {
	case "RS256":
		signingMethodRSA = jwt.SigningMethodRS256
	case "RS384":
		signingMethodRSA = jwt.SigningMethodRS384
	case "RS512":
		signingMethodRSA = jwt.SigningMethodRS512
	default:
		panic(fmt.Sprintf("Unsupported signing method: %s", *algoParam))
	}
	return signingMethodRSA
}

func parseProxy() proxy.RequestParser {
	var parser proxy.RequestParser
	switch strings.ToLower(*proxyParam) {
	case "simple":
		parser = proxy.NewSimpleParser()
	case "traefik":
		parser = proxy.NewTraefikParser()
	default:
		panic(fmt.Sprintf("Unsupported proxy provider: %s", *proxyParam))
	}
	return parser
}

func parseRepository() repository.Repository {
	var repo repository.Repository
	var err error
	switch strings.ToLower(*repoParam) {
	case "local":
		repo, err = repository.NewLocal()
		util.Check(err)
	case "aerospike":
		repo, err = repository.NewAerospike()
		util.Check(err)
	case "vault":
		repo, err = repository.NewVault()
		util.Check(err)
	default:
		panic(fmt.Sprintf("Unsupported storage provider: %s", *repoParam))
	}
	return repo
}
