// This file uses material from the Wikipedia article <a href="https://en.wikipedia.org/wiki/List_of_snakes_by_common_name">"List of snakes by common name"</a>, which is released under the <a href="https://creativecommons.org/licenses/by-sa/3.0/">Creative Commons Attribution-Share-Alike License 3.0</a>.
package commands

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var snakeNames = []string{
	"Adder",
	"Aesculapian snake",
	"Anaconda",
	"Arafura file snake",
	"Asp",
	"African beaked snake",
	"Ball Python",
	"Bird snake",
	"Black-headed snake",
	"Mexican black kingsnake",
	"Black rat snake",
	"Black snake",
	"Blind snake",
	"Boa",
	"Boiga",
	"Boomslang",
	"Brown snake",
	"Bull snake",
	"Bushmaster",
	"Dwarf beaked snake",
	"Rufous beaked snake",
	"Canebrake",
	"Cantil",
	"Cascabel",
	"Cat-eyed snake",
	"Cat snake",
	"Chicken snake",
	"Coachwhip snake",
	"Cobra",
	"Collett's snake",
	"Congo snake",
	"Copperhead",
	"Coral snake",
	"Corn snake",
	"Cottonmouth",
	"Crowned snake",
	"Cuban wood snake",
	"Eastern hognose snake",
	"Egg-eater",
	"Eyelash viper",
	"Eastern coral snake",
	"Fer-de-lance",
	"Fierce snake",
	"Fishing snake",
	"Flying snake",
	"Fox snake",
	"Forest flame snake",
	"Garter snake",
	"Glossy snake",
	"Gopher snake",
	"Grass snake",
	"Green snake",
	"Ground snake",
	"Habu",
	"Harlequin snake",
	"Herald snake",
	"Hognose snake",
	"Hoop snake",
	"Hundred pacer",
	"Ikaheka snake",
	"Indigo snake",
	"Jamaican Tree Snake",
	"Jamaican Tree Snake",
	"Jararacussu",
	"Keelback",
	"King brown",
	"King cobra",
	"King snake",
	"Krait",
	"Large shield snake",
	"Lancehead",
	"Lora",
	"Lyre snake",
	"Machete savane",
	"Mamba",
	"Mamushi",
	"Mangrove snake",
	"Milk snake",
	"Moccasin snake",
	"Montpellier snake",
	"Mud snake",
	"Mussurana",
	"Night snake",
	"Nose-horned viper",
	"Parrot snake",
	"Patchnose snake",
	"Perrotet's shieldtail snake",
	"Pine snake",
	"Pipe snake",
	"Python",
	"Queen snake",
	"Racer",
	"Raddysnake",
	"Rat snake",
	"Rattlesnake",
	"Ribbon snake",
	"Rinkhals",
	"River jack",
	"Sea snake",
	"Shield-tailed snake",
	"Sidewinder",
	"Small-eyed snake",
	"Stiletto snake",
	"Striped snake",
	"Sunbeam snake",
	"Taipan",
	"Tentacled snake",
	"Tic polonga",
	"Tiger snake",
	"Tigre snake",
	"Tree snake",
	"Trinket snake",
	"Twig snake",
	"Twin Headed King Snake",
	"Titanoboa",
	"Urutu",
	"Vine snake",
	"Viper",
	"Wart snake",
	"Water moccasin",
	"Water snake",
	"Whip snake",
	"Wolf snake",
	"Worm snake",
	"Wutu",
	"Yarara",
	"Zebra snake",
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
