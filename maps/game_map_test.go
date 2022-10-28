package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestMetadataValidate(t *testing.T) {
	for label, test := range map[string]struct {
		metadata   Metadata
		boardState *rules.BoardState
		expected   error
	}{
		"unlimited": {
			Metadata{
				BoardSizes: AnySize(),
			},
			rules.NewBoardState(99, 99),
			nil,
		},
		"in sizes": {
			Metadata{
				BoardSizes: OddSizes(7, 25),
			},
			rules.NewBoardState(7, 7),
			nil,
		},
		"too small": {
			Metadata{
				BoardSizes: OddSizes(7, 25),
			},
			rules.NewBoardState(6, 6),
			rules.RulesetError("This map can only be played on these board sizes: 7x7, 9x9, 11x11, 13x13, 15x15, 17x17, 19x19, 21x21, 23x23, 25x25"),
		},
		"too large": {
			Metadata{
				BoardSizes: OddSizes(7, 25),
			},
			rules.NewBoardState(26, 26),
			rules.RulesetError("This map can only be played on these board sizes: 7x7, 9x9, 11x11, 13x13, 15x15, 17x17, 19x19, 21x21, 23x23, 25x25"),
		},
		"valid players": {
			Metadata{
				BoardSizes: AnySize(),
				MinPlayers: 4,
				MaxPlayers: 4,
			},
			&rules.BoardState{
				Snakes: []rules.Snake{
					{ID: "1"},
					{ID: "2"},
					{ID: "3"},
					{ID: "4"},
				},
			},
			nil,
		},
		"too few players": {
			Metadata{
				BoardSizes: AnySize(),
				MinPlayers: 3,
				MaxPlayers: 4,
			},
			&rules.BoardState{
				Snakes: []rules.Snake{
					{ID: "1"},
					{ID: "2"},
				},
			},
			rules.RulesetError("This map can only be played with 3-4 players"),
		},
		"too many players": {
			Metadata{
				BoardSizes: AnySize(),
				MinPlayers: 3,
				MaxPlayers: 4,
			},
			&rules.BoardState{
				Snakes: []rules.Snake{
					{ID: "1"},
					{ID: "2"},
					{ID: "3"},
					{ID: "4"},
					{ID: "5"},
				},
			},
			rules.RulesetError("This map can only be played with 3-4 players"),
		},
	} {
		t.Run(label, func(t *testing.T) {
			actual := test.metadata.Validate(test.boardState)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestMapSizes(t *testing.T) {
	s := FixedSizes(Dimensions{11, 12})
	require.Equal(t, s[0].Width, 11)
	require.Equal(t, s[0].Height, 12)

	s = FixedSizes(Dimensions{11, 11}, Dimensions{19, 25})
	require.Len(t, s, 2)
	require.Equal(t, s[1].Width, 19)
	require.Equal(t, s[1].Height, 25)

	s = AnySize()
	require.Len(t, s, 1, "unlimited maps should have just one dimensions")
	require.True(t, s.IsUnlimited())
}

func TestBoardStateEditorInterface(t *testing.T) {
	var _ Editor = (*BoardStateEditor)(nil)
}

func TestBoardStateEditor(t *testing.T) {
	boardState := rules.NewBoardState(11, 11)
	boardState.Snakes = append(boardState.Snakes, rules.Snake{
		ID:     "existing_snake",
		Health: 100,
	})

	editor := BoardStateEditor{boardState: boardState}

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

	expected := rules.NewBoardState(11, 11).
		WithFood([]rules.Point{
			{X: 1, Y: 3},
			{X: 3, Y: 7},
		}).
		WithHazards([]rules.Point{
			{X: 1, Y: 3},
			{X: 3, Y: 7},
		}).
		WithSnakes([]rules.Snake{
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
		})
	require.Equal(t, expected, boardState)

	require.Equal(t, []rules.Point{
		{X: 1, Y: 3},
		{X: 3, Y: 7},
	}, editor.Food())

	require.Equal(t, []rules.Point{
		{X: 1, Y: 3},
		{X: 3, Y: 7},
	}, editor.Hazards())

	require.Equal(t, map[string][]rules.Point{
		"existing_snake": {
			{X: 5, Y: 2}, {X: 5, Y: 1}, {X: 5, Y: 0},
		},
		"new_snake": {

			{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1},
		},
	}, editor.SnakeBodies())

	editor.ClearFood()
	require.Equal(t, []rules.Point{}, boardState.Food)

	editor.ClearHazards()
	require.Equal(t, []rules.Point{}, boardState.Hazards)
}

func TestBoardStateEditorPlaceSnakesRandomlyAtPositions(t *testing.T) {
	for label, test := range map[string]struct {
		rand           rules.Rand
		initialSnakes  []rules.Snake
		heads          []rules.Point
		bodyLength     int
		expectedError  error
		expectedSnakes []rules.Snake
	}{
		"empty": {
			rules.MinRand,
			[]rules.Snake{},
			[]rules.Point{},
			0,
			nil,
			[]rules.Snake{},
		},
		"too many snakes": {
			rules.MinRand,
			[]rules.Snake{
				{ID: "1"}, {ID: "2"}, {ID: "3"},
			},
			[]rules.Point{{X: 3, Y: 3}, {X: 6, Y: 2}},
			3,
			rules.ErrorTooManySnakes,
			nil,
		},
		"success unshuffled": {
			rules.MinRand,
			[]rules.Snake{
				{ID: "1"}, {ID: "2"},
			},
			[]rules.Point{{X: 3, Y: 3}, {X: 6, Y: 2}},
			3,
			nil,
			[]rules.Snake{
				{
					ID:     "1",
					Body:   []rules.Point{{X: 3, Y: 3}, {X: 3, Y: 3}, {X: 3, Y: 3}},
					Health: rules.SnakeMaxHealth,
				}, {
					ID:     "2",
					Body:   []rules.Point{{X: 6, Y: 2}, {X: 6, Y: 2}, {X: 6, Y: 2}},
					Health: rules.SnakeMaxHealth,
				},
			},
		},
		"success shuffled": {
			rules.MaxRand,
			[]rules.Snake{
				{ID: "1"}, {ID: "2"},
			},
			[]rules.Point{{X: 3, Y: 3}, {X: 6, Y: 2}},
			3,
			nil,
			[]rules.Snake{
				{
					ID:     "1",
					Body:   []rules.Point{{X: 6, Y: 2}, {X: 6, Y: 2}, {X: 6, Y: 2}},
					Health: rules.SnakeMaxHealth,
				}, {
					ID:     "2",
					Body:   []rules.Point{{X: 3, Y: 3}, {X: 3, Y: 3}, {X: 3, Y: 3}},
					Health: rules.SnakeMaxHealth,
				},
			},
		},
	} {
		t.Run(label, func(t *testing.T) {
			boardState := rules.NewBoardState(rules.BoardSizeSmall, rules.BoardSizeSmall)
			boardState.Snakes = test.initialSnakes
			editor := NewBoardStateEditor(boardState)

			err := editor.PlaceSnakesRandomlyAtPositions(test.rand, test.initialSnakes, test.heads, test.bodyLength)
			if test.expectedError != nil {
				require.Equal(t, test.expectedError, err)
			} else {
				require.Equal(t, test.expectedSnakes, boardState.Snakes)
			}
		})
	}
}

func TestBoardStateEditorIsOccupied(t *testing.T) {
	for label, test := range map[string]struct {
		boardState            *rules.BoardState
		point                 rules.Point
		snakes, hazards, food bool
		expected              bool
	}{
		"empty board": {
			rules.NewBoardState(rules.BoardSizeSmall, rules.BoardSizeSmall),
			rules.Point{X: 3, Y: 3},
			true, true, true,
			false,
		},
		"unoccupied": {
			&rules.BoardState{
				Food:    []rules.Point{{X: 1, Y: 1}},
				Hazards: []rules.Point{{X: 2, Y: 2}},
				Snakes: []rules.Snake{
					{
						ID:   "1",
						Body: []rules.Point{{X: 3, Y: 3}},
					},
				},
			},
			rules.Point{X: 2, Y: 3},
			true, true, true,
			false,
		},
		"food": {
			&rules.BoardState{
				Food: []rules.Point{{X: 1, Y: 1}},
			},
			rules.Point{X: 1, Y: 1},
			false, false, true,
			true,
		},
		"ignored food": {
			&rules.BoardState{
				Food: []rules.Point{{X: 1, Y: 1}},
			},
			rules.Point{X: 1, Y: 1},
			false, false, false,
			false,
		},
		"hazard": {
			&rules.BoardState{
				Hazards: []rules.Point{{X: 1, Y: 1}},
			},
			rules.Point{X: 1, Y: 1},
			false, true, false,
			true,
		},
		"ignored hazard": {
			&rules.BoardState{
				Food: []rules.Point{{X: 1, Y: 1}},
			},
			rules.Point{X: 1, Y: 1},
			false, false, false,
			false,
		},
		"snake": {
			&rules.BoardState{
				Snakes: []rules.Snake{
					{
						ID:   "1",
						Body: []rules.Point{{X: 1, Y: 1}},
					},
				},
			},
			rules.Point{X: 1, Y: 1},
			true, false, false,
			true,
		},
		"ignored snake": {
			&rules.BoardState{
				Snakes: []rules.Snake{
					{
						ID:   "1",
						Body: []rules.Point{{X: 1, Y: 1}},
					},
				},
			},
			rules.Point{X: 1, Y: 1},
			false, false, false,
			false,
		},
	} {
		t.Run(label, func(t *testing.T) {
			editor := NewBoardStateEditor(test.boardState)

			actual := editor.IsOccupied(test.point, test.snakes, test.hazards, test.food)

			require.Equal(t, test.expected, actual)
		})
	}
}

func TestBoardStateEditorOccupiedPoints(t *testing.T) {
	testBoardState := &rules.BoardState{
		Food:    []rules.Point{{X: 1, Y: 1}},
		Hazards: []rules.Point{{X: 2, Y: 2}},
		Snakes: []rules.Snake{
			{
				ID:   "1",
				Body: []rules.Point{{X: 3, Y: 3}},
			},
		},
	}

	for label, test := range map[string]struct {
		boardState            *rules.BoardState
		snakes, hazards, food bool
		expected              map[rules.Point]bool
	}{
		"empty board": {
			rules.NewBoardState(rules.BoardSizeSmall, rules.BoardSizeSmall),
			true, true, true,
			map[rules.Point]bool{},
		},
		"all types": {
			testBoardState,
			true, true, true,
			map[rules.Point]bool{
				{X: 1, Y: 1}: true,
				{X: 2, Y: 2}: true,
				{X: 3, Y: 3}: true,
			},
		},
		"ignore snakes": {
			testBoardState,
			false, true, true,
			map[rules.Point]bool{
				{X: 1, Y: 1}: true,
				{X: 2, Y: 2}: true,
			},
		},
		"ignore hazards": {
			testBoardState,
			true, false, true,
			map[rules.Point]bool{
				{X: 1, Y: 1}: true,
				{X: 3, Y: 3}: true,
			},
		},
		"ignore food": {
			testBoardState,
			true, true, false,
			map[rules.Point]bool{
				{X: 2, Y: 2}: true,
				{X: 3, Y: 3}: true,
			},
		},
	} {
		t.Run(label, func(t *testing.T) {
			editor := NewBoardStateEditor(test.boardState)

			actual := editor.OccupiedPoints(test.snakes, test.hazards, test.food)

			require.Equal(t, test.expected, actual)
		})
	}
}

func TestBoardStateEditorFilterUnoccupiedPoints(t *testing.T) {
	testBoardState := &rules.BoardState{
		Food:    []rules.Point{{X: 1, Y: 1}},
		Hazards: []rules.Point{{X: 2, Y: 2}},
		Snakes: []rules.Snake{
			{
				ID:   "1",
				Body: []rules.Point{{X: 3, Y: 3}},
			},
		},
	}

	for label, test := range map[string]struct {
		boardState            *rules.BoardState
		targets               []rules.Point
		snakes, hazards, food bool
		expected              []rules.Point
	}{
		"empty": {
			rules.NewBoardState(rules.BoardSizeSmall, rules.BoardSizeSmall),
			[]rules.Point{},
			true, true, true,
			[]rules.Point{},
		},
		"all types": {
			testBoardState,
			[]rules.Point{{X: 3, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 2, Y: 1}},
			true, true, true,
			[]rules.Point{{X: 2, Y: 1}},
		},
		"ignore snakes": {
			testBoardState,
			[]rules.Point{{X: 3, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 2, Y: 1}},
			false, true, true,
			[]rules.Point{{X: 3, Y: 3}, {X: 2, Y: 1}},
		},
		"ignore hazards": {
			testBoardState,
			[]rules.Point{{X: 3, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 2, Y: 1}},
			true, false, true,
			[]rules.Point{{X: 2, Y: 2}, {X: 2, Y: 1}},
		},
		"ignore food": {
			testBoardState,
			[]rules.Point{{X: 3, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 2, Y: 1}},
			true, true, false,
			[]rules.Point{{X: 1, Y: 1}, {X: 2, Y: 1}},
		},
	} {
		t.Run(label, func(t *testing.T) {
			editor := NewBoardStateEditor(test.boardState)

			actual := editor.FilterUnoccupiedPoints(test.targets, test.snakes, test.hazards, test.food)

			require.Equal(t, test.expected, actual)
		})
	}
}

func TestBoardStateEditorShufflePoints(t *testing.T) {
	editor := NewBoardStateEditor(rules.NewBoardState(rules.BoardSizeSmall, rules.BoardSizeSmall))
	points := []rules.Point{{X: 4, Y: 0}, {X: 3, Y: 1}, {X: 2, Y: 2}, {X: 1, Y: 3}, {X: 0, Y: 4}}

	editor.ShufflePoints(rules.MaxRand, points)
	expected := []rules.Point{{X: 3, Y: 1}, {X: 2, Y: 2}, {X: 1, Y: 3}, {X: 0, Y: 4}, {X: 4, Y: 0}}

	require.Equal(t, expected, points)
}
