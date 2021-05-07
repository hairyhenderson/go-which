// CLI for the which command.
//
// See https://github.com/hairyhenderson/go-which for details.
package main

import (
	"os"

	"github.com/hairyhenderson/go-which/internal/cmd"
)

func main() {
	success := cmd.Run(os.Args, os.Stdout, os.Stderr)

	if !success {
		os.Exit(1)
	}
}
