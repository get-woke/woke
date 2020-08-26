// Package cmd ...
/*
Copyright © 2020 Caitlin Elfring <celfring@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caitlinelfring/woke/pkg/config"
	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/caitlinelfring/woke/pkg/parser"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const defaultGlob = "**"

var (
	// flags
	exitOneOnFailure bool
	ruleConfig       string
	debug            bool
	stdin            bool

	// Version is populated by goreleaser during build
	// Version...
	Version = "main"
	// Commit ...
	Commit = "000000"
	// Date ...
	Date = "today"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "woke [globs ...]",
	Short: "Check for usage of non-inclusive language in your code and provide alternatives",
	Long: `
woke is a linter that will check your source code for usage of non-inclusive
language and provide suggestions for alternatives. Rules can be customized
to suit your needs.

Provide a list file globs for files you'd like to check.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		setLogLevel()
		log.Debug().Msg(getVersion("default"))

		start := time.Now()
		defer func() {
			log.Debug().
				Dur("durationMS", time.Now().Sub(start)).
				Msg("woke completed")
		}()

		c, err := config.NewConfig(ruleConfig)
		if err != nil {
			return err
		}
		p := parser.Parser{Rules: c.Rules}

		// var results result.FileResults
		var ignorer *ignore.Ignore
		var fileArgs []string
		if stdin {
			fileArgs = []string{"/dev/stdin"}
		} else {
			ignorer, _ = ignore.NewIgnore(c.IgnoreFiles)
			fileArgs = args
		}
		results := p.ParseFiles(fileArgs, ignorer)
		for _, fr := range results {
			cmd.Println(fr.String())
		}

		// cmd.Println(results.OutputString(output))

		if len(results) == 0 {
			log.Info().Msg("👏 Great work using inclusive language in your code! Stay woke! 🙌")
		}

		if len(results) > 0 && exitOneOnFailure {
			// We intentionally return an error if exitOneOnFailure is true, but don't want to show usage
			cmd.SilenceUsage = true
			return fmt.Errorf("Total failures: %d", len(results))
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = getVersion("short")

	rootCmd.PersistentFlags().StringVarP(&ruleConfig, "rule-config", "r", "", "YAML file with list of rules")
	rootCmd.PersistentFlags().BoolVar(&exitOneOnFailure, "exit-1-on-failure", false, "Exit with exit code 1 on failures. Otherwise, will always exit 0 if any failures occur")
	rootCmd.PersistentFlags().BoolVar(&stdin, "stdin", false, "Read from stdin")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
}

func setLogLevel() {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func getVersion(t string) string {
	switch strings.ToLower(t) {
	case "short":
		return Version
	default:
		return fmt.Sprintf("woke version %s built from %s on %s", Version, Commit, Date)
	}
}
