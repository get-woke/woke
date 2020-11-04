package parser

import (
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
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

func getTestPrinterResults(p *Parser) []*result.FileResults {
	_printer := p.Printer
	_testPrinter := _printer.(*testPrinter)
	return _testPrinter.results
}

func testParser() *Parser {
	return NewParser(testRules(), ignore.NewIgnore([]string{}), new(testPrinter))
}

func testRules() []*rule.Rule {
	testRule := rule.NewTestRule()
	return []*rule.Rule{&testRule}
}

func parsePathTests(t *testing.T) {
	t.Run("violation", func(t *testing.T) {
		f, err := newFile(t, "i have a test-rule")
		assert.NoError(t, err)

		p := testParser()
		testRule := p.Rules[0]
		violations := p.ParsePaths(f.Name())

		assert.Len(t, getTestPrinterResults(p), 1)
		assert.Equal(t, len(getTestPrinterResults(p)), violations)

		filename := filepath.ToSlash(f.Name())
		expected := result.FileResults{
			Filename: filename,
			Results: []result.Result{
				result.LineResult{
					Rule:      testRule,
					Violation: "test-rule",
					Line:      "i have a test-rule",
					StartPosition: &token.Position{
						Filename: filename,
						Offset:   0,
						Line:     1,
						Column:   9,
					},
					EndPosition: &token.Position{
						Filename: filename,
						Offset:   0,
						Line:     1,
						Column:   18,
					},
				},
			},
		}
		assert.EqualValues(t, &expected, getTestPrinterResults(p)[0])
	})

	t.Run("no violations", func(t *testing.T) {
		f, err := newFile(t, "i have a no violations\n")
		assert.NoError(t, err)

		p := testParser()
		violations := p.ParsePaths(f.Name())

		assert.NoError(t, err)
		assert.Len(t, getTestPrinterResults(p), 0)
		assert.Equal(t, len(getTestPrinterResults(p)), violations)
	})
	t.Run("IsTextFileFromFilename failure", func(t *testing.T) {
		f, err := newFile(t, "")
		assert.NoError(t, err)

		p := testParser()
		violations := p.ParsePaths(f.Name())
		assert.NoError(t, err)
		assert.Len(t, getTestPrinterResults(p), 0)
		assert.Equal(t, len(getTestPrinterResults(p)), violations)
	})

	t.Run("multiple paths", func(t *testing.T) {
		f1, err := newFile(t, "i have a test-rule\n")
		assert.NoError(t, err)
		f2, err := newFile(t, "i have a no violations\n")
		assert.NoError(t, err)

		// Test with multiple paths supplied
		p := testParser()
		violations := p.ParsePaths(f1.Name(), f2.Name())
		assert.NoError(t, err)

		assert.Len(t, getTestPrinterResults(p), 1)
		assert.Equal(t, len(getTestPrinterResults(p)), violations)
	})

	t.Run("ignored", func(t *testing.T) {
		f, err := newFile(t, "i have a test-rule violation, but am ignored\n")
		assert.NoError(t, err)

		p := testParser()
		p.Ignorer = ignore.NewIgnore([]string{filepath.ToSlash(f.Name())})

		violations := p.ParsePaths(f.Name())
		assert.NoError(t, err)
		assert.Len(t, getTestPrinterResults(p), 0)
		assert.Equal(t, len(getTestPrinterResults(p)), violations)
	})

	t.Run("default path", func(t *testing.T) {
		// Test default path (which would run tests against the parser package)
		p := testParser()
		violations := p.ParsePaths()

		assert.Equal(t, len(getTestPrinterResults(p)), violations)
		assert.Greater(t, len(getTestPrinterResults(p)), 0)
	})

	t.Run("stdin", func(t *testing.T) {
		err := writeToStdin(t, "i have a test-rule here\n", func() {
			p := testParser()
			testRule := p.Rules[0]
			violations := p.ParsePaths(os.Stdin.Name())
			assert.Len(t, getTestPrinterResults(p), 1)
			assert.Equal(t, len(getTestPrinterResults(p)), violations)

			filename := filepath.ToSlash(os.Stdin.Name())
			expected := result.FileResults{
				Filename: filename,
				Results: []result.Result{
					result.LineResult{
						Rule:      testRule,
						Violation: "test-rule",
						Line:      "i have a test-rule here",
						StartPosition: &token.Position{
							Filename: filename,
							Offset:   0,
							Line:     1,
							Column:   9,
						},
						EndPosition: &token.Position{
							Filename: filename,
							Offset:   0,
							Line:     1,
							Column:   18,
						},
					},
				},
			}
			assert.EqualValues(t, &expected, getTestPrinterResults(p)[0])
		})
		assert.NoError(t, err)
	})
}

func TestParser_ParsePaths(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("WORKER_POOL_COUNT")
	})
	os.Unsetenv("WORKER_POOL_COUNT")
	parsePathTests(t)

	os.Setenv("WORKER_POOL_COUNT", "10")
	parsePathTests(t)
}

func writeToStdin(t *testing.T, text string, f func()) error {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "")
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
		_, _ = tmpFile.WriteString("this whitelist, he put in man hours to sanity-check the master/slave dummy-value. we can do better.\n") // wokeignore:rule=whitelist,man-hours,sanity,master-slave,slave,dummy
	}
	tmpFile.Close()

	for i := 0; i < b.N; i++ {
		p := testParser()
		violations := p.ParsePaths(tmpFile.Name())
		assert.Equal(b, 1, violations)
	}
}

func BenchmarkParsePathsRoot(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	for i := 0; i < b.N; i++ {
		assert.NotPanics(b, func() {
			p := testParser()
			// Unknown how many violations this will return since it's parsing the whole repo
			// there's no way to know for sure at any given time, so just check that it doesn't panic
			_ = p.ParsePaths("../..")
		})
	}
}
