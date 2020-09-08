package parser

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type testPrinter struct {
	results []*result.FileResults
}

func (p *testPrinter) Print(r *result.FileResults) error {
	p.results = append(p.results, r)
	return nil
}

func testParser() *Parser {
	return NewParser(rule.DefaultRules, ignore.NewIgnore([]string{}, []string{}))
}

func parsePathTests(t *testing.T) {
	f1, err := newFile(t, "i have a whitelist\n")
	assert.NoError(t, err)
	defer os.Remove(f1.Name())
	pr := new(testPrinter)
	p := testParser()
	violations := p.ParsePaths(pr, f1.Name())

	// assert.Len(t, pr.results, 1)
	assert.Equal(t, violations, 1)
	expected := result.FileResults{
		Filename: f1.Name(),
		Results: []result.Result{
			{
				Rule:      &rule.WhitelistRule,
				Violation: "whitelist",
				Line:      "i have a whitelist",
				StartPosition: &token.Position{
					Filename: f1.Name(),
					Offset:   0,
					Line:     1,
					Column:   9,
				},
				EndPosition: &token.Position{
					Filename: f1.Name(),
					Offset:   0,
					Line:     1,
					Column:   18,
				},
			},
		},
	}
	assert.EqualValues(t, &expected, pr.results[0])

	f2, err := newFile(t, "i have a no violations\n")
	assert.NoError(t, err)
	defer os.Remove(f2.Name())
	p = testParser()
	pr = new(testPrinter)
	violations = p.ParsePaths(pr, f2.Name())
	assert.NoError(t, err)
	assert.Len(t, pr.results, 0)
	assert.Equal(t, violations, 0)

	// Test for IsTextFileFromFilename failure
	f3, err := newFile(t, "")
	assert.NoError(t, err)
	defer os.Remove(f3.Name())

	p = testParser()
	pr = new(testPrinter)
	violations = p.ParsePaths(pr, f3.Name())
	assert.NoError(t, err)
	assert.Equal(t, violations, 0)
	assert.Len(t, pr.results, 0)

	// Test with multiple paths supplied
	// Disabled since this functionality is currently broken
	p = testParser()
	pr = new(testPrinter)
	violations = p.ParsePaths(pr, f1.Name(), f2.Name())
	assert.NoError(t, err)
	fmt.Println(pr.results)
	assert.Equal(t, violations, 1)
	assert.Len(t, pr.results, 1)

	// Test ignored file
	f4, err := newFile(t, "i have a whitelist violation, but am ignored\n")
	assert.NoError(t, err)
	defer os.Remove(f4.Name())

	p = testParser()
	p.Ignorer = ignore.NewIgnore([]string{f4.Name()}, []string{})
	pr = new(testPrinter)

	violations = p.ParsePaths(pr, f4.Name())
	assert.NoError(t, err)
	assert.Len(t, pr.results, 0)
	assert.Equal(t, violations, len(pr.results))

	// Test default path (which would run tests against the parser package)
	p = testParser()
	pr = new(testPrinter)
	violations = p.ParsePaths(pr)
	assert.NoError(t, err)
	fmt.Println(pr.results)
	assert.Equal(t, violations, len(pr.results))
	assert.Greater(t, len(pr.results), 0)

	// Stdin
	err = writeToStdin(t, "i have a whitelist here\n", func() {
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, os.Stdin.Name())
		assert.NoError(t, err)
		assert.Len(t, pr.results, 1)
		assert.Equal(t, violations, 1)
		expected := result.FileResults{
			Filename: os.Stdin.Name(),
			Results: []result.Result{
				{
					Rule:      &rule.WhitelistRule,
					Violation: "whitelist",
					Line:      "i have a whitelist here",
					StartPosition: &token.Position{
						Filename: os.Stdin.Name(),
						Offset:   0,
						Line:     1,
						Column:   9,
					},
					EndPosition: &token.Position{
						Filename: os.Stdin.Name(),
						Offset:   0,
						Line:     1,
						Column:   18,
					},
				},
			},
		}
		assert.EqualValues(t, &expected, pr.results[0])
	})
	assert.NoError(t, err)
}

func TestParser_ParsePaths(t *testing.T) {
	os.Unsetenv("WORKER_POOL_COUNT")
	parsePathTests(t)

	os.Setenv("WORKER_POOL_COUNT", "10")
	defer os.Unsetenv("WORKER_POOL_COUNT")
	parsePathTests(t)
}

func writeToStdin(t *testing.T, text string, f func()) error {
	tmpfile, err := ioutil.TempFile(t.TempDir(), "")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(text)); err != nil {
		return err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return err
	}

	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()
	os.Stdin = tmpfile
	f()
	return tmpfile.Close()
}

func BenchmarkParsePaths(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
	tmpFile, err := ioutil.TempFile(os.TempDir(), "")
	assert.NoError(b, err)

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())

	for i := 0; i < 100; i++ {
		_, _ = tmpFile.WriteString("this whitelist, he put in man hours to sanity-check the master/slave dummy-value. we can do better.\n")
	}
	tmpFile.Close()

	for i := 0; i < b.N; i++ {
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, tmpFile.Name())
		assert.Equal(b, violations, 6)
	}
}

func BenchmarkParsePathsRoot(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	for i := 0; i < b.N; i++ {
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, "../..")
		assert.Equal(b, violations, 6)
	}
}
