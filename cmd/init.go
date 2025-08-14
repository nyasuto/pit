package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize repository",
		Args:  cobra.RangeArgs(2, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			itemOne, err := strconv.Atoi(args[0])

			if err != nil {
				log.Fatal(err)
				return nil
			}

			itemTwo, err := strconv.Atoi(args[1])

			if err != nil {
				log.Fatal(err)
				return nil
			}

			fmt.Println(itemOne + itemTwo)

			return nil
		},
	}

	return cmd
}
