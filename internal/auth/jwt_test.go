package auth

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reugn/auth-server/internal/repository"
)

func TestJWT_Authorize(t *testing.T) {
	os.Setenv(repository.EnvLocalConfigPath, repository.DefaultLocalConfigPath)
	repo, err := repository.NewLocal()
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv(envPrivateKeyPath, defaultPrivateKeyPath)
	os.Setenv(envPublicKeyPath, defaultPublicKeyPath)
	keys, err := NewKeys()
	if err != nil {
		t.Fatal(err)
	}
	tokenGenerator := NewJWTGenerator(keys, jwt.SigningMethodRS256)
	tokenValidator := NewJWTValidator(keys, repo)

	tests := []struct {
		name       string
		username   string
		password   string
		request    repository.RequestDetails
		authorized bool
	}{
		{
			"configured-uri",
			"admin",
			"1234",
			repository.RequestDetails{
				Method: "GET",
				URI:    "/health",
			},
			true,
		},
		{
			"unknown-uri",
			"admin",
			"1234",
			repository.RequestDetails{
				Method: "GET",
				URI:    "/health2",
			},
			false,
		},
		{
			"invalid-user",
			"admin2",
			"1234",
			repository.RequestDetails{
				Method: "GET",
				URI:    "/health",
			},
			false,
		},
		{
			"invalid-password",
			"admin",
			"1111",
			repository.RequestDetails{
				Method: "GET",
				URI:    "/health",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userDetails := repo.AuthenticateBasic(tt.username, tt.password)
			if userDetails == nil {
				if tt.authorized {
					t.Fatal("authentication failed")
				} else {
					return
				}
			}
			token, err := tokenGenerator.Generate(tt.username, userDetails.UserRole)
			if err != nil {
				t.Fatal(err)
			}
			authorized := tokenValidator.Authorize(token.Token, &tt.request)
			if authorized != tt.authorized {
				t.Fatal("authorization result mismatch")
			}
		})
	}
}
