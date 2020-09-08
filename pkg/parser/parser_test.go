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
	t.Run("violation", func(t *testing.T) {
		f1, err := newFile(t, "i have a whitelist\n")
		assert.NoError(t, err)
		defer os.Remove(f1.Name())
		pr := new(testPrinter)
		p := testParser()
		violations := p.ParsePaths(pr, f1.Name())

		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), violations)
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
	})

	t.Run("no violations", func(t *testing.T) {
		f, err := newFile(t, "i have a no violations\n")
		assert.NoError(t, err)

		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, f.Name())

		assert.NoError(t, err)
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), violations)
	})
	t.Run("IsTextFileFromFilename failure", func(t *testing.T) {
		f, err := newFile(t, "")
		assert.NoError(t, err)
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, f.Name())
		assert.NoError(t, err)
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), violations)
	})

	t.Run("multiple paths", func(t *testing.T) {
		f1, err := newFile(t, "i have a whitelist\n")
		assert.NoError(t, err)
		f2, err := newFile(t, "i have a no violations\n")
		assert.NoError(t, err)

		// Test with multiple paths supplied
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, f1.Name(), f2.Name())
		assert.NoError(t, err)
		fmt.Println(pr.results)
		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), violations)
	})

	t.Run("ignored", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist violation, but am ignored\n")
		assert.NoError(t, err)

		p := testParser()
		p.Ignorer = ignore.NewIgnore([]string{f.Name()}, []string{})
		pr := new(testPrinter)

		violations := p.ParsePaths(pr, f.Name())
		assert.NoError(t, err)
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), violations)
	})

	t.Run("default path", func(t *testing.T) {
		// Test default path (which would run tests against the parser package)
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr)

		assert.Equal(t, len(pr.results), violations)
		assert.Greater(t, len(pr.results), 0)
	})

	t.Run("stdin", func(t *testing.T) {
		err := writeToStdin(t, "i have a whitelist here\n", func() {
			p := testParser()
			pr := new(testPrinter)
			violations := p.ParsePaths(pr, os.Stdin.Name())
			assert.Len(t, pr.results, 1)
			assert.Equal(t, len(pr.results), violations)

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
	})
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
	// TODO: Use b.TempDir() instead of os.TempDir()
	// Fix in go 1.16: https://github.com/golang/go/issues/41062
	tmpFile, err := ioutil.TempFile(os.TempDir(), "")
	assert.NoError(b, err)

	// Remember to clean up the file afterwards
	// TODO: Can be removed once b.TempDir() is used above, since the testing package
	// cleans up directories for us
	defer os.Remove(tmpFile.Name())

	for i := 0; i < b.N; i++ {
		_, _ = tmpFile.WriteString("this whitelist, he put in man hours to sanity-check the master/slave dummy-value. we can do better.\n")
	}
	tmpFile.Close()

	for i := 0; i < b.N; i++ {
		p := testParser()
		pr := new(testPrinter)
		violations := p.ParsePaths(pr, tmpFile.Name())
		assert.Equal(b, 1, violations)
	}
}

func BenchmarkParsePathsRoot(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	for i := 0; i < b.N; i++ {
		assert.NotPanics(b, func() {
			p := testParser()
			pr := new(testPrinter)
			// Unknown how many violations this will return since it's parsing the whole repo
			// there's no way to know for sure at any given time, so just check that it doesn't panic
			_ = p.ParsePaths(pr, "../..")
		})
	}
}
