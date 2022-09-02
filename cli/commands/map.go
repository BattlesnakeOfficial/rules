package commands

import (
	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"
)

func NewMapCommand() *cobra.Command {

	var mapCmd = &cobra.Command{
		Use:   "map",
		Short: "Display map information",
		Long:  "Display map information",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.ERROR.Fatal(err)
			}
		},
	}

	return mapCmd
}
