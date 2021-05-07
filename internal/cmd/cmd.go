package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/hairyhenderson/go-which"
	"github.com/hairyhenderson/go-which/internal/version"
)

func Run(args []string, stdOut, stdErr io.Writer) bool {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(stdErr)

	var all, silent, ver bool

	fs.BoolVar(&all, "a", false, "List all instances of executables found (instead of just the first).")
	fs.BoolVar(&silent, "s", false, "No output, just return 0 if all executables are found, or 1 if some were not found.")
	fs.BoolVar(&ver, "v", false, "Print the version")

	err := fs.Parse(args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return false
	}

	if err != nil {
		fs.Usage()
		fmt.Fprintf(stdErr, "error: %v\n", err)

		return false
	}

	if ver {
		fmt.Fprintf(stdErr, "%s (go-which) version %s (%s)\n", args[0], version.Version, version.GitCommit)

		return true
	}

	programs := fs.Args()
	if len(programs) == 0 {
		fs.Usage()

		return false
	}

	if all {
		found := which.All(programs...)

		if len(found) == 0 {
			return false
		}

		fmt.Fprintln(stdOut, strings.Join(found, "\n"))

		return true
	}

	if silent {
		return which.Found(programs...)
	}

	found := which.Which(programs...)
	if found == "" {
		return false
	}

	fmt.Fprintf(stdOut, "%s\n", found)

	return true
}
