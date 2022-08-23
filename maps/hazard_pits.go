package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

func init() {
	globalRegistry.RegisterMap("hz_hazard_pits", HazardPitsMap{})
}

type HazardPitsMap struct{}

func (m HazardPitsMap) ID() string {
	return "hz_hazard_pits"
}

func (m HazardPitsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_hazard_pits",
		Description: "A map that that fills in grid-like pattern of squares with pits filled with hazard sauce. Every N turns the pits will fill with another layer of sauce. After 4 cycles the pits will drain and the process starts again",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  4,
		BoardSizes:  FixedSizes(Dimensions{11, 11}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m HazardPitsMap) AddHazardPits(board *rules.BoardState, settings rules.Settings, editor Editor) error {
	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			if x%2 == 1 && y%2 == 1 {
				point := rules.Point{X: x, Y: y}
				isStartPosition := false
				for _, startPos := range hazardPitStartPositions {
					if startPos == point {
						isStartPosition = true
					}
				}
				if !isStartPosition {
					editor.AddHazard(point)
				}
			}
		}
	}
	return nil
}

func (m HazardPitsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	rand.Shuffle(len(hazardPitStartPositions), func(i int, j int) {
		hazardPitStartPositions[i], hazardPitStartPositions[j] = hazardPitStartPositions[j], hazardPitStartPositions[i]
	})

	// Place snakes
	if len(initialBoardState.Snakes) > len(hazardPitStartPositions) {
		return rules.ErrorTooManySnakes
	}
	for index, snake := range initialBoardState.Snakes {
		head := hazardPitStartPositions[index]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, snake.Health)
	}

	err := rules.PlaceFoodFixed(rand, initialBoardState)
	if err != nil {
		return err
	}

	err = m.AddHazardPits(initialBoardState, settings, editor)
	if err != nil {
		return err
	}

	return nil
}

func (m HazardPitsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	if lastBoardState.Turn == 0 {
		return nil
	}
	// Turn 0 = 1 layer of hazards
	// Turn 1*<ShrinkEveryNTurns> = 2 layers of hazards
	// Turn 2*<ShrinkEveryNTurns> = 3 layers of hazards
	// Turn 3*<ShrinkEveryNTurns> = 4 layers of hazards
	// Turn 4*<ShrinkEveryNTurns> = restart the cycle, 1 layers of hazards

	if lastBoardState.Turn%settings.RoyaleSettings.ShrinkEveryNTurns == 0 {
		// Time to change the hazards
		phase := lastBoardState.Turn % 4

		// clear all existing hazards
		editor.ClearHazards()

		// Add a layers of hazard pits depending on the phase
		for n := 0; n < (phase + 1); n++ {
			m.AddHazardPits(lastBoardState, settings, editor)
		}

		// Remove food from the pit when it is full of sauce?
	}
	return nil
}

var hazardPitStartPositions = []rules.Point{
	{X: 1, Y: 1},
	{X: 9, Y: 1},
	{X: 1, Y: 9},
	{X: 9, Y: 9},
}
