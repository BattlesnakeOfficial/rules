package maps

import (
	"math"

	"github.com/BattlesnakeOfficial/rules"
)

type SinkholesMap struct{}

func init() {
	globalRegistry.RegisterMap("sinkholes", SinkholesMap{})
}

func (m SinkholesMap) ID() string {
	return "sinkholes"
}

func (m SinkholesMap) Meta() Metadata {
	return Metadata{
		Name:        "Sinkholes",
		Description: "Spawns a rounded sinkhole on the board that grows every N turns, layering additional hazard squares over previously spawned ones.",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  FixedSizes(Dimensions{7, 7}, Dimensions{11, 11}, Dimensions{19, 19}),
		Tags:        []string{TAG_HAZARD_PLACEMENT},
	}
}

func (m SinkholesMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(initialBoardState, settings, editor)
}

func (m SinkholesMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m SinkholesMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.PostUpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn
	startTurn := 1
	spawnEveryNTurns := 10
	shrinkEveryNTurns := settings.Int(rules.ParamShrinkEveryNTurns, 0)
	if shrinkEveryNTurns > 0 {
		spawnEveryNTurns = shrinkEveryNTurns
	}
	maxRings := 5
	if lastBoardState.Width == 7 {
		maxRings = 3
	} else if lastBoardState.Width == 19 {
		maxRings = 7
	}

	spawnLocation := rules.Point{X: lastBoardState.Width / 2, Y: lastBoardState.Height / 2}

	if currentTurn == startTurn {
		editor.AddHazard(spawnLocation)
		return nil
	}

	// Are we at max size, if so stop try to generate hazards
	if currentTurn > spawnEveryNTurns*maxRings {
		return nil
	}

	// Is this a turn to grow the sinkhole?
	if (currentTurn-startTurn)%spawnEveryNTurns != 0 {
		return nil
	}

	offset := int(math.Floor(float64(currentTurn-startTurn) / float64(spawnEveryNTurns)))

	if offset > 0 && offset <= maxRings {
		for x := spawnLocation.X - offset; x <= spawnLocation.X+offset; x++ {
			for y := spawnLocation.Y - offset; y <= spawnLocation.Y+offset; y++ {
				// don't draw in the corners of the square so we get a rounded effect
				if !(x == spawnLocation.X-offset && y == spawnLocation.Y-offset) &&
					!(x == spawnLocation.X+offset && y == spawnLocation.Y-offset) &&
					!(x == spawnLocation.X-offset && y == spawnLocation.Y+offset) &&
					!(x == spawnLocation.X+offset && y == spawnLocation.Y+offset) {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
		}
	}

	return nil
}
