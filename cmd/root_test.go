package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// run profiling with
// go test -v -cpuprofile cpu.prof -memprofile mem.prof -bench=. ./cmd
// memory:
//    go tool pprof mem.prof
// cpu:
//    go tool pprof cpu.prof
func BenchmarkExecute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert.NoError(b, TestExecute())
	}
}
