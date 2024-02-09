package repository

import (
	"os"
	"testing"
)

func Test_getVaultConfig(t *testing.T) {
	os.Setenv(envVaultAddr, "127.0.0.1:8200")
	os.Setenv(envVaultToken, "token1")
	os.Setenv(envVaultBasicKey, "secret/basic1")
	os.Setenv(envVaultAuthKey, "secret/authorization1")

	config := getVaultConfig()
	if config.vaultAddr != "127.0.0.1:8200" {
		t.Fail()
	}
	if config.vaultToken != "token1" {
		t.Fail()
	}
	if config.basicAuthKeyPrefix != "secret/basic1" {
		t.Fail()
	}
	if config.authorizationKeyPrefix != "secret/authorization1" {
		t.Fail()
	}
}
