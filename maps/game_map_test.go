package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestBoardStateEditorInterface(t *testing.T) {
	var _ Editor = (*BoardStateEditor)(nil)
}

func TestBoardStateEditor(t *testing.T) {
	boardState := rules.NewBoardState(11, 11)
	boardState.Snakes = append(boardState.Snakes, rules.Snake{
		ID:     "existing_snake",
		Health: 100,
	})

	editor := BoardStateEditor{BoardState: boardState}

	editor.AddFood(rules.Point{X: 1, Y: 3})
	editor.AddFood(rules.Point{X: 3, Y: 6})
	editor.AddFood(rules.Point{X: 3, Y: 7})
	editor.RemoveFood(rules.Point{X: 3, Y: 6})
	editor.AddHazard(rules.Point{X: 1, Y: 3})
	editor.AddHazard(rules.Point{X: 3, Y: 6})
	editor.AddHazard(rules.Point{X: 3, Y: 7})
	editor.RemoveHazard(rules.Point{X: 3, Y: 6})
	editor.PlaceSnake("existing_snake", []rules.Point{{X: 5, Y: 2}, {X: 5, Y: 1}, {X: 5, Y: 0}}, 99)
	editor.PlaceSnake("new_snake", []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}, 98)

	require.Equal(t, &rules.BoardState{
		Width:  11,
		Height: 11,
		Food: []rules.Point{
			{X: 1, Y: 3},
			{X: 3, Y: 7},
		},
		Hazards: []rules.Point{
			{X: 1, Y: 3},
			{X: 3, Y: 7},
		},
		Snakes: []rules.Snake{
			{
				ID:     "existing_snake",
				Health: 99,
				Body:   []rules.Point{{X: 5, Y: 2}, {X: 5, Y: 1}, {X: 5, Y: 0}},
			},
			{
				ID:     "new_snake",
				Health: 98,
				Body:   []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			},
		},
	}, boardState)

	editor.ClearFood()
	require.Equal(t, []rules.Point{}, boardState.Food)

	editor.ClearHazards()
	require.Equal(t, []rules.Point{}, boardState.Hazards)
}

func TestMapSizes(t *testing.T) {
	s := FixedSizes(Dimensions{11, 12})
	require.Equal(t, s[0].Width, uint(11))
	require.Equal(t, s[0].Height, uint(12))

	s = FixedSizes(Dimensions{11, 11}, Dimensions{19, 25})
	require.Len(t, s, 2)
	require.Equal(t, s[1].Width, uint(19))
	require.Equal(t, s[1].Height, uint(25))

	s = AnySize()
	require.Len(t, s, 1, "unlimited maps should have just one dimensions")
	require.True(t, s.IsUnlimited())
}
