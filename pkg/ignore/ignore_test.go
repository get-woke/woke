package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreMatch(t *testing.T) {
	i, err := NewIgnore("my/fles/*")
	assert.NoError(t, err)
	assert.NotNil(t, i)

	assert.False(t, i.Match("not/foo"))
	assert.True(t, i.Match("my/files/file1"))
}

func TestReadIgnoreFIle(t *testing.T) {
	lines := readIgnoreFile("testdata/.gitignore")

	assert.Equal(t, "*.DS_Store", lines)
}
