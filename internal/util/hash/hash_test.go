package hash_test

import (
	"testing"

	"github.com/reugn/auth-server/internal/util/hash"
)

func TestSha256(t *testing.T) {
	if hash.Sha256("1234") != "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4" {
		t.Fatal("Sha256")
	}
}
