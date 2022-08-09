package commands

import (
	"fmt"
	"log"

	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/spf13/cobra"
)

func NewMapInfoCommand() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:   "info [flags] map_name [...map_name]",
		Short: "Display metadata for given map(s)",
		Long:  "Display metadata for given map(s)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
			}
			for i, m := range args {
				gameMap, err := maps.GetMap(m)
				if err != nil {
					log.Fatalf("Failed to load game map %#v: %v", m, err)
				}
				meta := gameMap.Meta()
				fmt.Println("Name:", meta.Name)
				fmt.Println("Author:", meta.Author)
				fmt.Println("Description:", meta.Description)
				fmt.Println("Version:", meta.Version)
				fmt.Println("Min Players:", meta.MinPlayers)
				fmt.Println("Max Players:", meta.MaxPlayers)
				fmt.Print("Board Sizes (WxH):")
				for j, s := range meta.BoardSizes {
					fmt.Printf(" %dx%d", s.Width, s.Height)
					if j == (len(meta.BoardSizes) - 1) {
						fmt.Print("\n")
					}
				}
				// separate map information when querying multiple maps
				if i < (len(args) - 1) {
					fmt.Print("\n")
				}
			}
		},
	}
	return infoCmd
}
