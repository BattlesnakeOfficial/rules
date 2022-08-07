package commands

import (
	"fmt"
	"log"

	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/spf13/cobra"
)

type mapCmdFlags struct {
	List bool
	Info string
}

func NewMapCommand() *cobra.Command {
	mapFlags := &mapCmdFlags{}

	var mapCmd = &cobra.Command{
		Use:   "map",
		Short: "Display map information",
		Long:  "Display map information",
		Run: func(cmd *cobra.Command, args []string) {
			mapFlags.Process()
		},
	}

	mapCmd.Flags().BoolVarP(&mapFlags.List, "list", "l", false, "List all available maps")
	mapCmd.Flags().StringVarP(&mapFlags.Info, "info", "i", "", "Display metadata for given map")

	mapCmd.Flags().SortFlags = false

	return mapCmd
}

func (m *mapCmdFlags) Process() {

	if m.List {
		for _, m := range maps.List() {
			fmt.Println(m)
		}
	}

	if m.Info != "" {
		gameMap, err := maps.GetMap(m.Info)
		if err != nil {
			log.Fatalf("Failed to load game map %#v: %v", m.Info, err)
		}
		meta := gameMap.Meta()
		fmt.Println("Name:", meta.Name)
		fmt.Println("Author:", meta.Author)
		fmt.Println("Description:", meta.Description)
		fmt.Println("Version:", meta.Version)
		fmt.Println("Min Players:", meta.MinPlayers)
		fmt.Println("Max Players:", meta.MaxPlayers)
		fmt.Print("Board Sizes (WxH):")
		for _, s := range meta.BoardSizes {
			fmt.Printf(" %dx%d", s.Width, s.Height)
		}
		fmt.Print("\n")
	}
}
