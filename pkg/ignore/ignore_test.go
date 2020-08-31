package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreMatch(t *testing.T) {
	i, err := NewIgnore([]string{"my/files/*"})
	assert.NoError(t, err)
	assert.NotNil(t, i)

	assert.False(t, i.Match("not/foo"))
	assert.True(t, i.Match("my/files/file1"))
}

func TestReadIgnoreFIle(t *testing.T) {
	lines := readIgnoreFile("testdata/.gitignore")

	assert.Equal(t, []string{"*.DS_Store"}, lines)
}
