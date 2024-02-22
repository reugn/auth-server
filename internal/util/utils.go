package util

import (
	"crypto/sha256"
	"fmt"
)

// Sha256 returns the Sha256 hexadecimal representation of the string.
func Sha256(str string) string {
	sha256pwd := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sha256pwd)
}

// Check panics if the error is not nil.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
