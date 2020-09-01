package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_Print(t *testing.T) {
	p := NewText(true)
	got := captureOutput(func() {
		assert.NoError(t, p.Print(generateFileResult()))
	})
	expected := "foo.txt\n\t5:3-5:12       warn                 `blacklist` may be insensitive, use `blocklist` instead\n\n"
	assert.Equal(t, expected, got)
}
