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

	expected := `{"Filename":"foo.txt","Results":[{"Rule":{"Name":"blacklist","Terms":["blacklist","black-list","blacklisted","black-listed"],"Alternatives":["denylist","blocklist"],"Note":"","Severity":"warning","Options":{"WordBoundary":false}},"Violation":"blacklist","Line":"this blacklist must change","StartPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":6},"EndPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":15},"Reason":"` + "`blacklist`" + ` may be insensitive, use ` + "`denylist`" + `, ` + "`blocklist`" + ` instead"}]}` + "\n"
	assert.Equal(t, expected, got)
}
