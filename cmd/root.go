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
	"path/filepath"

	"github.com/caitlinelfring/woke/pkg/config"
	"github.com/caitlinelfring/woke/pkg/parser"
	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/spf13/cobra"
)

var FileGlob string
var ExitOneOnFailure bool
var RuleConfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "woke",
	Short: "Check for usage of non-inclusive language in your code and provide alternatives",
	Long: `woke is a linter that will check your source code for usage of non-inclusive
	language and provide suggestions for alternatives. Rules can be customized to suit your needs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		c, err := config.NewConfig(RuleConfig)
		if err != nil {
			return err
		}
		p := parser.Parser{Rules: c.Rules}
		files, err := filepath.Glob(FileGlob)
		if err != nil {
			return err
		}

		allResults := []*rule.Result{}
		for _, f := range files {
			// skip rules config, which will always produce failures
			if f == RuleConfig {
				continue
			}
			results, err := p.Parse(f)
			if err != nil {
				fmt.Printf("parser failed on %s: %s\n", f, err)
			}
			allResults = append(allResults, results...)
		}

		for _, r := range allResults {
			fmt.Println(r)
		}
		if len(allResults) > 0 && ExitOneOnFailure {
			return fmt.Errorf("Total failures: %d", len(allResults))
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&FileGlob, "file-glob", "f", "**", "File glob of files to check for inclusive language")
	rootCmd.PersistentFlags().StringVarP(&RuleConfig, "rule-config", "r", "default.yaml", "YAML file with list of rules")
	rootCmd.PersistentFlags().BoolVar(&ExitOneOnFailure, "exit-1-on-failure", false, "Exit with exit code 1 on failures. Otherwise, will always exit 0 if any failures occur")
}
