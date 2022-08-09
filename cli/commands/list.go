package commands

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/spf13/cobra"
)

func NewMapListCommand() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available game maps",
		Long:  "List available game maps",
		Run: func(cmd *cobra.Command, args []string) {
			for _, m := range maps.List() {
				fmt.Println(m)
			}
		},
	}
	return listCmd
}
