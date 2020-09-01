package result

import (
	"testing"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestFileResult_String(t *testing.T) {
	rs := FindResults(&rule.WhitelistRule, "my/file", "this has the term whitelist", 1)
	fr := FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file\n    my/file:1:18-my/file:1:27 warn       `whitelist` may be insensitive, use `allowlist` instead", fr.String())

	rs = FindResults(&rule.WhitelistRule, "my/file", "this has no rule violations", 1)
	fr = FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file", fr.String())

}
