package util

import (
	"os"
	"strconv"
)

// GetEnvDefault returns the value of the environment variable, or a default
// value if the environment variable is not defined or is an empty string
func GetEnvDefault(envVar, defaultValue string) string {
	if v, ok := os.LookupEnv(envVar); ok && len(v) > 0 {
		return v
	}
	return defaultValue
}

// GetEnvBoolDefault is similar to GetEnvDefault, but with booleans instead of strings
func GetEnvBoolDefault(envVar string, defaultValue bool) bool {
	val := GetEnvDefault(envVar, strconv.FormatBool(defaultValue))
	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}
	return false
}
