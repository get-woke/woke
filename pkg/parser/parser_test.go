package parser

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
)

type testPrinter struct {
	results []*result.FileResults
}

// Print doesn't actually write anything, just stores the results in memory so they can be read later
func (p *testPrinter) Print(r *result.FileResults) error {
	p.results = append(p.results, r)
	return nil
}

func (p *testPrinter) Start() {
}

func (p *testPrinter) End() {
}

func (p *testPrinter) PrintSuccessExitMessage() bool {
	return true
}

func testParser() (parser *Parser, err error) {
	r := rule.TestRule
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	fs := osfs.New(cwd)
	ignorer, err := ignore.NewIgnore(fs, []string{})
	if err != nil {
		return
	}
	parser = NewParser([]*rule.Rule{&r}, ignorer)
	return
}

func parsePathTests(t *testing.T) {
	t.Run("finding", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist")
		assert.NoError(t, err)

		pr := new(testPrinter)
		p, err := testParser()
		assert.NoError(t, err)
		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), findings)

		filename := filepath.ToSlash(f.Name())
		expected := result.FileResults{
			Filename: filename,
			Results: []result.Result{
				result.LineResult{
					Rule:    &rule.TestRule,
					Finding: "whitelist",
					Line:    "i have a whitelist",
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
		assert.EqualValues(t, &expected, pr.results[0])
	})

	t.Run("no findings", func(t *testing.T) {
		f, err := newFile(t, "i have a no findings\n")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		pr := new(testPrinter)
		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("finding in filename - empty file", func(t *testing.T) {
		f, err := newFileWithPrefix(t, "whitelist", "")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		pr := new(testPrinter)
		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("IsTextFileFromFilename failure", func(t *testing.T) {
		f, err := newFile(t, "")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		pr := new(testPrinter)
		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("multiple paths", func(t *testing.T) {
		f1, err := newFile(t, "i have a whitelist\n")
		assert.NoError(t, err)
		f2, err := newFile(t, "i have a no findings\n")
		assert.NoError(t, err)

		// Test with multiple paths supplied
		p, err := testParser()
		assert.NoError(t, err)
		pr := new(testPrinter)
		findings := p.ParsePaths(pr, f1.Name(), f2.Name())
		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("ignored", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist finding, but am ignored\n")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		cwd, err := os.Getwd()
		assert.NoError(t, err)
		fs := osfs.New(cwd)
		ignorer, err := ignore.NewIgnore(fs, []string{filepath.ToSlash(f.Name())})
		assert.NoError(t, err)
		p.Ignorer = ignorer
		pr := new(testPrinter)

		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("ignored inline", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist finding, but am ignored # wokeignore:rule=whitelist\n")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		pr := new(testPrinter)

		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 0)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("ignored inline with no ignorer", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist finding, but am ignored # wokeignore:rule=whitelist\n")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		p.Ignorer = nil
		pr := new(testPrinter)

		findings := p.ParsePaths(pr, f.Name())
		assert.Len(t, pr.results, 1)
		assert.Equal(t, len(pr.results), findings)
	})

	t.Run("default path", func(t *testing.T) {
		// Test default path (which would run tests against the parser package)
		p, err := testParser()
		assert.NoError(t, err)
		cwd, err := os.Getwd()
		assert.NoError(t, err)
		fs := osfs.New(cwd)
		ignorer, err := ignore.NewIgnore(fs, []string{"*_test.go"})
		assert.NoError(t, err)
		p.Ignorer = ignorer
		pr := new(testPrinter)
		findings := p.ParsePaths(pr)

		assert.Equal(t, len(pr.results), findings)
		assert.Equal(t, len(pr.results), 0)
	})

	t.Run("stdin", func(t *testing.T) {
		err := writeToStdin(t, "i have a whitelist here\n", func() {
			p, err := testParser()
			assert.NoError(t, err)
			pr := new(testPrinter)
			findings := p.ParsePaths(pr, os.Stdin.Name())
			assert.Len(t, pr.results, 1)
			assert.Equal(t, len(pr.results), findings)

			filename := filepath.ToSlash(os.Stdin.Name())
			expected := result.FileResults{
				Filename: filename,
				Results: []result.Result{
					result.LineResult{
						Rule:    &rule.TestRule,
						Finding: "whitelist",
						Line:    "i have a whitelist here",
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
			assert.EqualValues(t, &expected, pr.results[0])
		})
		assert.NoError(t, err)
	})

	t.Run("note is not included in output message", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist")
		assert.NoError(t, err)
		const TestNote = "TEST NOTE"
		p, err := testParser()
		assert.NoError(t, err)
		p.Rules[0].Note = TestNote
		p.Rules[0].Options.IncludeNote = nil
		pr := new(testPrinter)
		p.ParsePaths(pr, f.Name())

		assert.NotContains(t, pr.results[0].Results[0].Reason(), TestNote)
	})

	t.Run("note is included in output message", func(t *testing.T) {
		f, err := newFile(t, "i have a whitelist")
		assert.NoError(t, err)
		const TestNote = "TEST NOTE"
		includeNote := true
		p, err := testParser()
		assert.NoError(t, err)
		p.Rules[0].Note = TestNote
		p.Rules[0].Options.IncludeNote = &includeNote
		// Test IncludeNote flag doesn't get overridden with SetIncludeNote method
		p.Rules[0].SetIncludeNote(false)
		pr := new(testPrinter)
		p.ParsePaths(pr, f.Name())

		assert.Contains(t, pr.results[0].Results[0].Reason(), TestNote)
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
	tmpfile, err := os.CreateTemp(os.TempDir(), "")
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
	tmpFile, err := os.CreateTemp(b.TempDir(), "")
	assert.NoError(b, err)

	for i := 0; i < 100; i++ {
		_, _ = tmpFile.WriteString("this whitelist, he put in man hours to sanity-check the master/slave dummy-value. we can do better.\n")
	}
	tmpFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p, err := testParser()
		assert.NoError(b, err)
		pr := new(testPrinter)
		findings := p.ParsePaths(pr, tmpFile.Name())
		assert.Equal(b, 1, findings)
	}
}

func BenchmarkParsePathsRoot(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	for i := 0; i < b.N; i++ {
		assert.NotPanics(b, func() {
			p, err := testParser()
			assert.NoError(b, err)
			pr := new(testPrinter)
			// Unknown how many findings this will return since it's parsing the whole repo
			// there's no way to know for sure at any given time, so just check that it doesn't panic
			_ = p.ParsePaths(pr, "../..")
		})
	}
}
