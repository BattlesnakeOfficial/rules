package commands

import (
	"github.com/spf13/cobra"
)

func NewMapCommand() *cobra.Command {

	var mapCmd = &cobra.Command{
		Use:   "map",
		Short: "Display map information",
		Long:  "Display map information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return mapCmd
}
