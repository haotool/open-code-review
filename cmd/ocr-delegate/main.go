package main

import (
	"fmt"
	"os"

	"github.com/alibaba/open-code-review/internal/delegatecli"
)

func main() {
	os.Exit(runMain(os.Args[1:]))
}

// runMain is the testable entry point; returns a process exit code.
func runMain(args []string) int {
	if len(args) > 0 {
		switch args[0] {
		case "--version", "-V", "version":
			printVersion()
			return 0
		}
	}

	if err := delegatecli.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
