package main

import (
	"crypto"
	"flag"

	"github.com/dgrijalva/jwt-go"
	"github.com/reugn/auth-server/auth"
	"github.com/reugn/auth-server/proxy"
	"github.com/reugn/auth-server/repository"
	"github.com/reugn/auth-server/utils"
)

var (
	serverHostParam = flag.String("host", "0.0.0.0", "Server host")
	serverPortParam = flag.Int("port", 8081, "Server port")
	algoParam       = flag.String("algo", "RS256", "JWT signing algorithm")
	proxyParam      = flag.String("proxy", "simple", "Proxy provider")
	repoParam       = flag.String("repo", "local", "Repository provider")
)

func main() {
	flag.Parse()
	parseAlgo()

	// load ssl keys
	keys, err := auth.NewKeys()
	utils.Check(err)

	// start http server
	server := NewHTTPServer(*serverHostParam, *serverPortParam, parseRepo(), keys)
	server.start()
}

func parseAlgo() *jwt.SigningMethodRSA {
	var signingMethodRSA *jwt.SigningMethodRSA
	switch *algoParam {
	case "RS256":
		signingMethodRSA = &jwt.SigningMethodRSA{
			Name: "RS256",
			Hash: crypto.SHA256,
		}
	case "RS384":
		signingMethodRSA = &jwt.SigningMethodRSA{
			Name: "RS384",
			Hash: crypto.SHA384,
		}
	case "RS512":
		signingMethodRSA = &jwt.SigningMethodRSA{
			Name: "RS512",
			Hash: crypto.SHA512,
		}
	default:
		panic("Invalid signing method")
	}
	return signingMethodRSA
}

func parseProxy() proxy.RequestParser {
	var parser proxy.RequestParser
	switch *proxyParam {
	case "simple":
		parser = proxy.NewSimpleParser()
	case "traefik":
		parser = proxy.NewTraefikParser()
	default:
		panic("Invalid proxy provider")
	}
	return parser
}

func parseRepo() repository.Repository {
	var repo repository.Repository
	var err error
	switch *repoParam {
	case "local":
		repo = repository.NewLocalRepo()
	case "aerospike":
		repo, err = repository.NewAerospikeRepositoryFromEnv()
		utils.Check(err)
	default:
		panic("Invalid repository provider")
	}
	return repo
}
