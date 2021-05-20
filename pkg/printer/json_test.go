package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON_Print(t *testing.T) {
	p := NewJSON()
	res := generateFileResult()
	buf := new(bytes.Buffer)
	assert.NoError(t, p.Print(buf, res))
	got := buf.String()

	expected := `{"Filename":"foo.txt","Results":[{"Rule":{"Name":"whitelist","Terms":["whitelist","white-list","whitelisted","white-listed"],"Alternatives":["allowlist"],"Note":"","Severity":"warning","Options":{"WordBoundary":false,"AddNoteToMessage":null}},"Violation":"whitelist","Line":"this whitelist must change","StartPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":6},"EndPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":15},"Reason":"` + "`whitelist`" + ` may be insensitive, use ` + "`allowlist`" + ` instead"}]}` + "\n"
	assert.Equal(t, expected, got)
}
