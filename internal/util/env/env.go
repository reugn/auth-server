package env

import (
	"os"
	"strconv"
	"time"
)

// ReadString retrieves the string value of the environment variable named
// by the key.
func ReadString(value *string, key string) {
	envValue, ok := os.LookupEnv(key)
	if ok {
		*value = envValue
	}
}

// ReadInt retrieves the integer value of the environment variable named
// by the key.
func ReadInt(value *int, key string) {
	envValue, ok := os.LookupEnv(key)
	if ok {
		intValue, err := strconv.Atoi(envValue)
		if err == nil {
			*value = intValue
		}
	}
}

// ReadTime retrieves the time value of the environment variable named
// by the key.
func ReadTime(value *time.Duration, key string, timeUnit time.Duration) {
	envValue, ok := os.LookupEnv(key)
	if ok {
		intValue, err := strconv.Atoi(envValue)
		if err == nil {
			*value = time.Duration(intValue) * timeUnit
		}
	}
}
