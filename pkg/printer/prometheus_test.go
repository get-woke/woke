/**
 * Copyright 2022 Cisco and its affiliates
 * All rights reserved.
**/

package printer

import (
	"bytes"
	"fmt"
	config "github.com/get-woke/woke/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrometheus_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewPrometheus(buf, new(config.Config))
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()
	fmt.Printf("buf: %s", buf)
	fmt.Printf("res: %s", res.Results[0].String())
	expected := fmt.Sprintf("woke_result{file=\"foo.txt:1:6-15\", term=\"%s\"} 1 \n", res.Results[0].GetRuleName())
	assert.Equal(t, expected, got)
}

func TestPrometheus_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewPrometheus(buf, new(config.Config))
	p.Start()
	got := buf.String()
	assert.Equal(t, ``, got)
}

func TestPrometheus_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewPrometheus(buf, new(config.Config))
	p.End()
	got := buf.String()
	assert.Equal(t, ``, got)
}

func TestPrometheus_PrintSuccessExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewPrometheus(buf, new(config.Config))
	assert.Equal(t, true, p.PrintSuccessExitMessage())
}
