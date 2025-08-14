package main

import (
	"fmt"
	"os"

	"github.com/nyasuto/pit/cmd" // Adjust the import path to where your cmd package is located
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(-1)
	}
}
