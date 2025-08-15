package cmd

import (
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
	Short: "A tiny, educational Git implementation in Go",
	Long: `Pit is an educational Git implementation that helps you understand
how Git works internally. It implements Git's core functionality
with a focus on simplicity and learning.

Pit creates '.pit' repositories instead of '.git' to avoid conflicts
with real Git repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
