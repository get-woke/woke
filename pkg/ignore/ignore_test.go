package ignore

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestIgnore_Match(t *testing.T) {
	i := NewIgnore([]string{"my/files/*"})
	assert.NotNil(t, i)

	assert.False(t, i.Match("not/foo"))
	assert.True(t, i.Match("my/files/file1"))
	assert.False(t, i.Match("my/files"))
}

func TestReadIgnoreFile(t *testing.T) {
	ignoreLines := readIgnoreFile("testdata/.gitignore")
	assert.Equal(t, []string{"*.DS_Store"}, ignoreLines)

	noIgnoreLines := readIgnoreFile(".gitignore")
	assert.Equal(t, []string{}, noIgnoreLines)
}
