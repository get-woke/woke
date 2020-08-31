package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownCodify(t *testing.T) {
	assert.Equal(t, "`mystring`", MarkdownCodify("mystring"))
}
