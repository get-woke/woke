package result

import (
	"sort"
	"testing"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestFileResult_String(t *testing.T) {
	rs := FindResults(&rule.WhitelistRule, "my/file", "this has the term whitelist", 1)
	fr := FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file\n    my/file:1:18-my/file:1:27 warning    `whitelist` may be insensitive, use `allowlist` instead", fr.String())

	rs = FindResults(&rule.WhitelistRule, "my/file", "this has no rule violations", 1)
	fr = FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file", fr.String())
}

func TestFileResult_Sort(t *testing.T) {
	rs1 := FindResults(&rule.WhitelistRule, "my/file", "this has a few whitelist white-list whitelist", 1)
	rs2 := FindResults(&rule.WhitelistRule, "my/file", "this whitelist has a few white-list whitelist", 2)

	rs := append(rs2, rs1...)

	fr := FileResults{Filename: "my/file", Results: rs}

	assert.False(t, sort.IsSorted(fr))
	sort.Sort(fr)
	assert.True(t, sort.IsSorted(fr))

	assert.EqualValues(t, fr.Results[0].StartPosition.Line, 1)
	assert.EqualValues(t, fr.Results[0].StartPosition.Column, 15)
	assert.EqualValues(t, fr.Results[1].StartPosition.Line, 1)
	assert.EqualValues(t, fr.Results[1].StartPosition.Column, 25)
	assert.EqualValues(t, fr.Results[2].StartPosition.Line, 1)
	assert.EqualValues(t, fr.Results[2].StartPosition.Column, 36)

	assert.EqualValues(t, fr.Results[3].StartPosition.Line, 2)
	assert.EqualValues(t, fr.Results[3].StartPosition.Column, 5)
	assert.EqualValues(t, fr.Results[4].StartPosition.Line, 2)
	assert.EqualValues(t, fr.Results[4].StartPosition.Column, 25)
	assert.EqualValues(t, fr.Results[5].StartPosition.Line, 2)
	assert.EqualValues(t, fr.Results[5].StartPosition.Column, 36)

}
