package printer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_Print(t *testing.T) {
	p := NewText(true)
	res := generateFileResult()
	got := captureOutput(func() {
		assert.NoError(t, p.Print(res))
	})
	expected := fmt.Sprintf("foo.txt\n\t5:3-5:12       warn                 %s\n\n", res.Results[0].Reason())
	assert.Equal(t, expected, got)
}
