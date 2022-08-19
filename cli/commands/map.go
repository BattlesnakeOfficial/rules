package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func NewMapCommand() *cobra.Command {

	var mapCmd = &cobra.Command{
		Use:   "map",
		Short: "Display map information",
		Long:  "Display map information",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return mapCmd
}
