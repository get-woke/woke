package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetEnvDefault returns the value of the environment variable, or a default
// value if the environment variable is not defined or is an empty string
func TestGetEnvDefault(t *testing.T) {
	val := GetEnvDefault("MY_ENV", "default_value")
	assert.Equal(t, "default_value", val)

	os.Setenv("MY_ENV", "defined_value")
	val = GetEnvDefault("MY_ENV", "default_value")
	assert.Equal(t, "defined_value", val)

	os.Unsetenv("MY_ENV")
}

func TestGetEnvBoolDefault(t *testing.T) {
	val := GetEnvBoolDefault("MY_BOOL_ENV", true)
	assert.Equal(t, true, val)

	os.Setenv("MY_BOOL_ENV", "true")
	val = GetEnvBoolDefault("MY_BOOL_ENV", false)
	assert.Equal(t, true, val)

	os.Unsetenv("MY_BOOL_ENV")

	os.Setenv("MY_BOOL_ENV", "notABool")
	val = GetEnvBoolDefault("MY_BOOL_ENV", true)
	assert.Equal(t, true, val)

	os.Unsetenv("MY_BOOL_ENV")
}
