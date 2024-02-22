package repository

import (
	"os"
	"testing"
)

func Test_getAerospikeConfig(t *testing.T) {
	os.Setenv(envAerospikeHost, "127.0.0.1")
	os.Setenv(envAerospikePort, "3300")
	os.Setenv(envAerospikeNamespace, "test1")
	os.Setenv(envAerospikeSet, "set1")
	os.Setenv(envAerospikeBasicKey, "basic1")
	os.Setenv(envAerospikeAuthKey, "authorization1")

	config := getAerospikeConfig()
	if config.hostname != "127.0.0.1" {
		t.Fail()
	}
	if config.port != 3300 {
		t.Fail()
	}
	if config.namespase != "test1" {
		t.Fail()
	}
	if config.setName != "set1" {
		t.Fail()
	}
	if config.basicAuthKey != "basic1" {
		t.Fail()
	}
	if config.authorizationKey != "authorization1" {
		t.Fail()
	}
}
