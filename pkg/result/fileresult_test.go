package result

import (
	"sort"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFileResult_String(t *testing.T) {
	r := rule.NewTestRule()
	rs := FindResults(r, "my/file", "this has the term testrule", 1)
	fr := FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file\n    my/file:1:18-my/file:1:26 error      `testrule` may be insensitive, use `better-rule` instead", fr.String())

	rs = FindResults(r, "my/file", "this has no rule violations", 1)
	fr = FileResults{Filename: "my/file", Results: rs}
	assert.Equal(t, "my/file", fr.String())
}

func TestFileResult_Sort(t *testing.T) {
	r := rule.NewTestRule()
	rs1 := FindResults(r, "my/file", "this has a few testrule test-rule testrule", 1)
	rs2 := FindResults(r, "my/file", "this testrule has a few test-rule testrule", 2)

	rs := append(rs2, rs1...)

	fr := FileResults{Filename: "my/file", Results: rs}

	assert.False(t, sort.IsSorted(fr))
	sort.Sort(fr)
	assert.True(t, sort.IsSorted(fr))

	assert.EqualValues(t, fr.Results[0].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[0].GetStartPosition().Column, 15)
	assert.EqualValues(t, fr.Results[1].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[1].GetStartPosition().Column, 24)
	assert.EqualValues(t, fr.Results[2].GetStartPosition().Line, 1)
	assert.EqualValues(t, fr.Results[2].GetStartPosition().Column, 34)

	assert.EqualValues(t, fr.Results[3].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[3].GetStartPosition().Column, 5)
	assert.EqualValues(t, fr.Results[4].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[4].GetStartPosition().Column, 24)
	assert.EqualValues(t, fr.Results[5].GetStartPosition().Line, 2)
	assert.EqualValues(t, fr.Results[5].GetStartPosition().Column, 34)
}
