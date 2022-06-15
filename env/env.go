// Package env contains shortcuts for reading and parsing environment variables.
package env

import (
	"fmt"
	"os"
	"strconv"
)

// Get retrieves the value of the environment variable named key. It returns fallback string if the variable is not present.
func Get(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

// MustGet retrieves the value of the environment variable named key. It panics if the variable is not present.
func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprint("Required environment variable not set: ", key))
	}

	return val
}

// GetBool retrieves the value of the environment variable named key as a boolean.
// If the value cannot be parsed as a boolean, the default is returned. If there is no default supplied, false is assumed.
func GetBool(key string, def ...bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}

	if len(def) != 0 {
		return def[0]
	}

	return false
}
