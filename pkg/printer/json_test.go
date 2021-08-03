package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()

	expected := `{"Filename":"foo.txt","Results":[{"Rule":{"Name":"whitelist","Terms":["whitelist","white-list","whitelisted","white-listed"],"Alternatives":["allowlist"],"Note":"","Severity":"warning","Options":{"WordBoundary":false,"WordBoundaryStart":false,"WordBoundaryEnd":false,"IncludeNote":null}},"Finding":"whitelist","Line":"this whitelist must change","StartPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":6},"EndPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":15},"Reason":"` + "`whitelist`" + ` may be insensitive, use ` + "`allowlist`" + ` instead"}]}` + "\n"
	assert.Equal(t, expected, got)
}
