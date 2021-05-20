package result

import (
	"sort"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFileResult_String(t *testing.T) {
	rs := FindResults(&rule.TestRule, "my/file", "this has the term whitelist", 1)
	fr := FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file\n    my/file:1:18-my/file:1:27 warning    `whitelist` may be insensitive, use `allowlist` instead", fr.String())

	rs = FindResults(&rule.TestRule, "my/file", "this has no rule violations", 1)
	fr = FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file", fr.String())
}

func TestFileResult_Sort(t *testing.T) {
	rs1 := FindResults(&rule.TestRule, "my/file", "this has a few whitelist white-list whitelist", 1)
	rs2 := FindResults(&rule.TestRule, "my/file", "this whitelist has a few white-list whitelist", 2)

	fr := FileResults{Filename: "my/file", Results: append(rs2, rs1...)}

	assert.False(t, sort.IsSorted(fr))
	sort.Sort(fr)
	assert.True(t, sort.IsSorted(fr))

	assert.EqualValues(t, fr.Results[0].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[0].GetStartPosition().Column, 15)
	assert.EqualValues(t, fr.Results[1].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[1].GetStartPosition().Column, 25)
	assert.EqualValues(t, fr.Results[2].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[2].GetStartPosition().Column, 36)

	assert.EqualValues(t, fr.Results[3].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[3].GetStartPosition().Column, 5)
	assert.EqualValues(t, fr.Results[4].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[4].GetStartPosition().Column, 25)
	assert.EqualValues(t, fr.Results[5].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[5].GetStartPosition().Column, 36)
}
