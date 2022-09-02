// This file uses material from the Wikipedia article <a href="https://en.wikipedia.org/wiki/List_of_snakes_by_common_name">"List of snakes by common name"</a>, which is released under the <a href="https://creativecommons.org/licenses/by-sa/3.0/">Creative Commons Attribution-Share-Alike License 3.0</a>.
package commands

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var snakeNames = []string{
	"Adder",
	"Aesculapian Snake",
	"Anaconda",
	"Arafura File Snake",
	"Asp",
	"African Beaked Snake",
	"Ball Python",
	"Bird Snake",
	"Black-headed Snake",
	"Mexican Black Kingsnake",
	"Black Rat Snake",
	"Black Snake",
	"Blind Snake",
	"Boa",
	"Boiga",
	"Boomslang",
	"Brown Snake",
	"Bull Snake",
	"Bushmaster",
	"Dwarf Beaked Snake",
	"Rufous Beaked Snake",
	"Canebrake",
	"Cantil",
	"Cascabel",
	"Cat-eyed Snake",
	"Cat Snake",
	"Chicken Snake",
	"Coachwhip Snake",
	"Cobra",
	"Collett's Snake",
	"Congo Snake",
	"Copperhead",
	"Coral Snake",
	"Corn Snake",
	"Cottonmouth",
	"Crowned Snake",
	"Cuban Wood Snake",
	"Egg-eater",
	"Eyelash Viper",
	"Fer-de-lance",
	"Fierce Snake",
	"Fishing Snake",
	"Flying Snake",
	"Fox Snake",
	"Forest Flame Snake",
	"Garter Snake",
	"Glossy Snake",
	"Gopher Snake",
	"Grass Snake",
	"Green Snake",
	"Ground Snake",
	"Habu",
	"Harlequin Snake",
	"Herald Snake",
	"Hognose Snake",
	"Hoop Snake",
	"Hundred Pacer",
	"Ikaheka Snake",
	"Indigo Snake",
	"Jamaican Tree Snake",
	"Jararacussu",
	"Keelback",
	"King Brown",
	"King Cobra",
	"King Snake",
	"Krait",
	"Lancehead",
	"Lora",
	"Lyre Snake",
	"Machete Savane",
	"Mamba",
	"Mamushi",
	"Mangrove Snake",
	"Milk Snake",
	"Moccasin Snake",
	"Montpellier Snake",
	"Mud Snake",
	"Mussurana",
	"Night Snake",
	"Nose-horned Viper",
	"Parrot Snake",
	"Patchnose Snake",
	"Pine Snake",
	"Pipe Snake",
	"Python",
	"Queen Snake",
	"Racer",
	"Raddysnake",
	"Rat Snake",
	"Rattlesnake",
	"Ribbon Snake",
	"Rinkhals",
	"River Jack",
	"Sea Snake",
	"Shield-tailed Snake",
	"Sidewinder",
	"Small-eyed Snake",
	"Stiletto Snake",
	"Striped Snake",
	"Sunbeam Snake",
	"Taipan",
	"Tentacled Snake",
	"Tic Polonga",
	"Tiger Snake",
	"Tigre Snake",
	"Tree Snake",
	"Trinket Snake",
	"Twig Snake",
	"Twin Headed King Snake",
	"Titanoboa",
	"Urutu",
	"Vine Snake",
	"Viper",
	"Wart Snake",
	"Water Moccasin",
	"Water Snake",
	"Whip Snake",
	"Wolf Snake",
	"Worm Snake",
	"Wutu",
	"Yarara",
	"Zebra Snake",
}

func init() {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randGen.Shuffle(len(snakeNames), func(i, j int) {
		snakeNames[i], snakeNames[j] = snakeNames[j], snakeNames[i]
	})
}

// Generate a random unique snake name, or return a UUID if there are no more names available.
func GenerateSnakeName() string {
	if len(snakeNames) == 0 {
		return uuid.New().String()
	}

	name := snakeNames[0]
	snakeNames = snakeNames[1:]

	return name
}
