package result

import (
	"fmt"
	"testing"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestFindResults(t *testing.T) {
	rs := FindResults(&rule.WhitelistRule, "my/file", "this has the term whitelist", 1)
	assert.Len(t, rs, 1)
	assert.Equal(t, rule.WhitelistRule.Reason("whitelist"), rs[0].Reason())
	assert.Equal(t, fmt.Sprintf("    my/file:1:18-my/file:1:27 warning    %s", rs[0].Reason()), rs[0].String())

	rs = FindResults(&rule.WhitelistRule, "my/file", "this has no rule violations", 1)
	assert.Len(t, rs, 0)

	rs = FindResults(&rule.WhitelistRule, "my/file", "this has the term whitelist #wokeignore:rule=whitelist", 1)
	assert.Len(t, rs, 0)
}
