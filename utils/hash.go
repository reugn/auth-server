package utils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(str string) string {
	sha256pwd := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sha256pwd)
}
