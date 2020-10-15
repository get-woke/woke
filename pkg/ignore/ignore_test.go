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
