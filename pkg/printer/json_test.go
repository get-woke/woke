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
	expected := `{"Filename":"foo.txt","Results":[{"Rule":{"Name":"blacklist","Terms":["blacklist","black-list","blacklisted","black-listed"],"Alternatives":["denylist","blocklist"],"Note":"","Severity":"warning"},"Violation":"blacklist","Line":"this blacklist must change","StartPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":6},"EndPosition":{"Filename":"foo.txt","Offset":0,"Line":1,"Column":15},"Reason":"` + "`blacklist`" + ` may be insensitive, use ` + "`denylist`" + `, ` + "`blocklist`" + ` instead"}]}` + "\n"
	assert.Equal(t, expected, got)
}
