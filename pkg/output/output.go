package output

import "github.com/mattn/go-colorable"

// These are the same as https://pkg.go.dev/github.com/fatih/color@v1.9.0#pkg-variables,
// but a central place to access them within this codebase.
// colorable package enables color support on Windows
var (
	// Stdout defines the standard output of the print functions. By default
	// os.Stdout is used.
	Stdout = colorable.NewColorableStdout()

	// Stderr defines a color supporting writer for os.Stderr.
	Stderr = colorable.NewColorableStderr()
)
