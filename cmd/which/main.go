/*
CLI for the which command
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	which "github.com/hairyhenderson/go-which"
	"github.com/hairyhenderson/go-which/internal/version"
)

func main() {
	var (
		all    bool
		silent bool
	)

	flag.BoolVar(&all, "a", false, "List all instances of executables found (instead of just the first).")
	flag.BoolVar(&silent, "s", false, "No output, just return 0 if all executables are found, or 1 if some were not found.")
	ver := flag.Bool("v", false, "Print the version")
	flag.Parse()

	if *ver {
		fmt.Printf("%s (go-which) version %s (%s)\n", os.Args[0], version.Version, version.GitCommit)
		return
	}

	programs := flag.Args()
	if len(programs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if all {
		found := which.All(programs...)
		if len(found) == 0 {
			os.Exit(1)
		}
		fmt.Println(strings.Join(found, "\n"))
		return
	}

	if silent {
		found := which.Found(programs...)
		if found {
			return
		}
		os.Exit(1)
	}

	found := which.Which(programs...)
	if found == "" {
		os.Exit(1)
	}
	fmt.Println(found)
}
