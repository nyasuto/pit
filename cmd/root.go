package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)
func init() {
	cobra.OnInitialize()
	RootCmd.AddCommand(
		initCmd(),
	)
}

// RootCmd is root command
var RootCmd = &cobra.Command{
	Use:   "pit",
	Short: "petit git",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root command")
	},
}
