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
	"runtime"
	"strings"
	"time"

	"github.com/get-woke/woke/pkg/config"
	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/output"
	"github.com/get-woke/woke/pkg/parser"
	"github.com/get-woke/woke/pkg/printer"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// flags
	exitOneOnFailure bool
	cfgFile          string
	debug            bool
	stdin            bool
	outputName       string
	noIgnore         bool

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
	RunE: rootRunE,
}

func rootRunE(cmd *cobra.Command, args []string) error {
	setDebugLogLevel()
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Debug().Msg(getVersion("default"))

	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Msg("woke completed")
	}()

	cfg, err := config.NewConfig(viper.ConfigFileUsed())
	if err != nil {
		return err
	}

	var ignorer *ignore.Ignore
	if !noIgnore {
		ignorer = ignore.NewIgnore(cfg.IgnoreFiles)
	}

	p := parser.NewParser(cfg.Rules, ignorer)

	print, err := printer.NewPrinter(outputName)
	if err != nil {
		return err
	}

	violations := p.ParsePaths(print, parseArgs(args)...)

	if exitOneOnFailure && violations > 0 {
		// We intentionally return an error if exitOneOnFailure is true, but don't want to show usage
		cmd.SilenceUsage = true
		err = fmt.Errorf("files with violations: %d", violations)
	}

	if violations == 0 {
		fmt.Fprintln(output.Stdout, "No violations found. Stay woke \u270a")
	}

	return err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Version = getVersion("short")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is .woke.yaml in current directory, or $HOME)")
	rootCmd.PersistentFlags().BoolVar(&exitOneOnFailure, "exit-1-on-failure", false, "Exit with exit code 1 on failures")
	rootCmd.PersistentFlags().BoolVar(&stdin, "stdin", false, "Read from stdin")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&noIgnore, "no-ignore", false, "Files matching entries in .gitignore/.wokeignore are parsed")
	rootCmd.PersistentFlags().StringVarP(&outputName, "output", "o", printer.OutFormatText, fmt.Sprintf("Output type [%s]", printer.OutFormatsString))
}

func parseArgs(args []string) []string {
	if len(args) == 0 {
		args = parser.DefaultPath
	}

	if stdin {
		args = []string{os.Stdin.Name()}
	}

	return args
}

func setDebugLogLevel() {
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

func initConfig() {
	// Require yaml for now, since the unmarshaling of the config only
	// supports yaml. The config loading will need to be refactored to
	// support viper unmarshaling.
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in working directory, then home directory with name ".woke.[yml|yaml]"
		viper.SetConfigName(".woke")
		viper.AddConfigPath(".")

		// Find home directory.
		if home, err := homedir.Dir(); err == nil {
			viper.AddConfigPath(home)
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}
