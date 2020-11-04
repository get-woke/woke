package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON_Print(t *testing.T) {
	p := NewJSON()
	res := generateFileResult()
	got := captureOutput(func() {
		assert.NoError(t, p.Print(res))
	})
	expected := `{"Filename":"foo.txt","Results":[{"Rule":{"Name":"test-rule","Terms":["testrule","test-rule"],"Alternatives":["better-rule"],"Note":"","Severity":"error","Options":{"WordBoundary":false}},"Violation":"testrule","Line":"this testrule must change","StartPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":6},"EndPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":15},"Reason":"` + "`testrule`" + ` may be insensitive, use ` + "`better-rule`" + ` instead"}]}` + "\n"
	assert.Equal(t, expected, got)
}
