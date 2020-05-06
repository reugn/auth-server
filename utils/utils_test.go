package utils_test

import (
	"reflect"
	"testing"

	"github.com/reugn/auth-server/utils"
)

func TestSha256(t *testing.T) {
	assertEqual(t, "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4", utils.Sha256("1234"))
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%v != %v", a, b)
	}
}
