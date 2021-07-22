package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownCodify(t *testing.T) {
	assert.Equal(t, "`mystring`", MarkdownCodify("mystring"))
}

func TestInSlice(t *testing.T) {
	tests := []struct {
		s         string
		sl        []string
		assertion assert.BoolAssertionFunc
	}{
		{"foo", []string{"foo", "bar"}, assert.True},
		{"bar", []string{"foo", "baz"}, assert.False},
		{"", []string{"", "baz"}, assert.True},
		{"", []string{"foo", "baz"}, assert.False},
		{"baz", []string{}, assert.False},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-%d", tt.s, i), func(t *testing.T) {
			tt.assertion(t, InSlice(tt.s, tt.sl))
		})
	}
}

func TestContainsAlphanumeric(t *testing.T) {
	tests := []struct {
		s         string
		assertion assert.BoolAssertionFunc
	}{
		{"foo", assert.True},
		{"bar123", assert.True},
		{"", assert.False},
		{" ", assert.False},
		{"123", assert.True},
		{"<-- -->", assert.False},
		{"#", assert.False},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s-%d", tt.s, i), func(t *testing.T) {
			tt.assertion(t, ContainsAlphanums(tt.s))
		})
	}
}
