package main

import (
	"fmt"
	"os"

	"github.com/BrandonThomas84/code-review-automation/internal/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

