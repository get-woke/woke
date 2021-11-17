package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValidURL(t *testing.T) {
	t.Run("valid-url-test1", func(t *testing.T) {
		boolResponse := isValidURL("https://raw.githubusercontent.com/get-woke/woke/main/example.yaml")
		assert.True(t, boolResponse)
	})

	t.Run("invalid-url-test1", func(t *testing.T) {
		boolResponse := isValidURL("Users/Document/test.yaml")
		assert.False(t, boolResponse)
	})

	t.Run("invalid-url-test2", func(t *testing.T) {
		boolResponse := isValidURL("/Users/Document/test.yaml")
		assert.False(t, boolResponse)
	})

	t.Run("invalid-url-test3", func(t *testing.T) {
		boolResponse := isValidURL("C:User\testpath\test.yaml")
		assert.False(t, boolResponse)
	})

	t.Run("invalid-url-test4", func(t *testing.T) {
		boolResponse := isValidURL("C:\\directory.com\test.yaml")
		assert.False(t, boolResponse)
	})
}
