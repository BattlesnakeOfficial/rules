package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"

	"github.com/BattlesnakeOfficial/rules/maps"
)

type mapInfo struct {
	All bool
}

func NewMapInfoCommand() *cobra.Command {
	info := mapInfo{}
	var infoCmd = &cobra.Command{
		Use:   "info [flags] map_name [...map_name]",
		Short: "Display metadata for given map(s)",
		Long:  "Display metadata for given map(s)",
		Run: func(cmd *cobra.Command, args []string) {
			// handle --all flag first as there would be no args
			if info.All {
				mapList := maps.List()
				for i, m := range mapList {
					info.display(m)
					if i < (len(mapList) - 1) {
						fmt.Print("\n")
					}
				}
				return
			}

			// display help when no map(s) provided via args
			if len(args) < 1 {
				err := cmd.Help()
				if err != nil {
					log.ERROR.Fatal(err)
				}
				return
			}

			// display all maps via command args
			for i, m := range args {
				info.display(m)
				if i < (len(args) - 1) {
					fmt.Print("\n")
				}
			}

		},
	}

	infoCmd.Flags().BoolVarP(&info.All, "all", "a", false, "Display information for all maps")

	return infoCmd
}

func (m *mapInfo) display(id string) {
	gameMap, err := maps.GetMap(id)
	if err != nil {
		log.ERROR.Fatalf("Failed to load game map %v: %v", id, err)
	}
	meta := gameMap.Meta()
	fmt.Println("Name:", meta.Name)
	fmt.Println("Author:", meta.Author)
	fmt.Println("Description:", meta.Description)
	fmt.Println("Version:", meta.Version)
	fmt.Println("Min Players:", meta.MinPlayers)
	fmt.Println("Max Players:", meta.MaxPlayers)
	fmt.Print("Board Sizes (WxH):")
	for i, s := range meta.BoardSizes {
		fmt.Printf(" %dx%d", s.Width, s.Height)
		if i == (len(meta.BoardSizes) - 1) {
			fmt.Print("\n")
		}
	}
	fmt.Print("Tags:")
	if len(meta.Tags) < 1 {
		fmt.Print("\n")
	}
	for i, t := range meta.Tags {
		fmt.Printf(" %s", t)
		if i == (len(meta.Tags) - 1) {
			fmt.Print("\n")
		}
	}
}
