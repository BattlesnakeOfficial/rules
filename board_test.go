package rules

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoardStateClone(t *testing.T) {
	empty := &BoardState{}
	require.Equal(t, NewBoardState(0, 0), empty.Clone())

	full := NewBoardState(11, 11).
		WithTurn(99).
		WithFood([]Point{{X: 1, Y: 2, TTL: 10, Value: 100}}).
		WithHazards([]Point{{X: 3, Y: 4, TTL: 5, Value: 50}}).
		WithSnakes([]Snake{
			{
				ID:               "1",
				Body:             []Point{{X: 1, Y: 2}},
				Health:           99,
				EliminatedCause:  EliminatedByCollision,
				EliminatedOnTurn: 45,
				EliminatedBy:     "2",
			},
		}).
		WithGameState(map[string]string{"example": "game data"}).
		WithPointState(map[Point]int{{X: 1, Y: 1}: 42})

	require.Equal(t, full, full.Clone())
}

func TestDev1235(t *testing.T) {
	// Small boards should no longer error and only get 1 food when num snakes > 4
	state, err := CreateDefaultBoardState(MaxRand, BoardSizeSmall, BoardSizeSmall, []string{
		"1", "2", "3", "4", "5", "6", "7", "8",
	})
	require.NoError(t, err)
	require.Len(t, state.Food, 1)
	state, err = CreateDefaultBoardState(MaxRand, BoardSizeSmall, BoardSizeSmall, []string{
		"1", "2", "3", "4", "5",
	})
	require.NoError(t, err)
	require.Len(t, state.Food, 1)

	// Small boards with <= 4 snakes should still get more than just center food
	state, err = CreateDefaultBoardState(MaxRand, BoardSizeSmall, BoardSizeSmall, []string{
		"1", "2", "3", "4",
	})
	require.NoError(t, err)
	require.Len(t, state.Food, 5)

	// Medium boards should still get 9 food
	state, err = CreateDefaultBoardState(MaxRand, BoardSizeMedium, BoardSizeMedium, []string{
		"1", "2", "3", "4", "5", "6", "7", "8",
	})
	require.NoError(t, err)
	require.Len(t, state.Food, 9)
}

func sortPoints(p []Point) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].X != p[j].X {
			return p[i].X < p[j].X
		}
		return p[i].Y < p[j].Y
	})
}

func TestCreateDefaultBoardState(t *testing.T) {
	tests := []struct {
		Height          int
		Width           int
		IDs             []string
		ExpectedNumFood int
		Err             error
	}{
		{1, 1, []string{"one"}, 0, ErrorNoRoomForSnake},
		{1, 2, []string{"one"}, 0, ErrorNoRoomForSnake},
		{1, 4, []string{"one"}, 1, nil},
		{2, 2, []string{"one"}, 1, nil},
		{9, 8, []string{"one"}, 1, nil},
		{2, 2, []string{"one", "two"}, 0, ErrorNoRoomForSnake},
		{1, 1, []string{"one", "two"}, 2, ErrorNoRoomForSnake},
		{1, 2, []string{"one", "two"}, 2, ErrorNoRoomForSnake},
		{BoardSizeSmall, BoardSizeSmall, []string{"one", "two"}, 3, nil},
		{
			BoardSizeSmall,
			BoardSizeSmall,
			[]string{"1", "2", "3", "4"},
			5, // <= 4 snakes on a small board we get more than just center food
			nil,
		},
		{
			BoardSizeSmall,
			BoardSizeSmall,
			[]string{"1", "2", "3", "4", "5"},
			1, // for this size and this many snakes, food is only placed in the center
			nil,
		},
		{
			BoardSizeSmall,
			BoardSizeSmall,
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
			1, // for this size and this many snakes, food is only placed in the center
			nil,
		},
		{
			BoardSizeMedium,
			BoardSizeMedium,
			[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"},
			17, // > small boards and we get non-center food
			nil,
		},
	}

	for testNum, test := range tests {
		t.Logf("test case %d", testNum)
		state, err := CreateDefaultBoardState(MaxRand, test.Width, test.Height, test.IDs)
		require.Equal(t, test.Err, err)
		if err != nil {
			require.Nil(t, state)
			continue
		}
		require.NotNil(t, state)
		require.Equal(t, test.Width, state.Width)
		require.Equal(t, test.Height, state.Height)
		require.Equal(t, len(test.IDs), len(state.Snakes))
		for i, id := range test.IDs {
			require.Equal(t, id, state.Snakes[i].ID)
		}
		require.Len(t, state.Food, test.ExpectedNumFood, testNum)
		require.Len(t, state.Hazards, 0, testNum)
	}
}

func TestPlaceSnakesDefault(t *testing.T) {
	// Because placement is random, we only test to ensure
	// that snake bodies are populated correctly
	// Note: because snakes are randomly spawned on even diagonal points, the board can accomodate number of snakes equal to: width*height/2
	// Update: because we exclude the center point now, we can accommodate 1 less snake now (width*height/2 - 1)
	tests := []struct {
		BoardState *BoardState
		SnakeIDs   []string
		Err        error
	}{
		{
			&BoardState{
				Width:  1,
				Height: 1,
			},
			make([]string, 1),
			ErrorNoRoomForSnake, // we avoid placing snakes in the center, so a board size of 1 will error
		},
		{
			&BoardState{
				Width:  1,
				Height: 1,
			},
			make([]string, 2),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  2,
				Height: 1,
			},
			make([]string, 2),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  1,
				Height: 2,
			},
			make([]string, 2),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  10,
				Height: 5,
			},
			make([]string, 24),
			nil,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
			},
			make([]string, 24),
			nil,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
			},
			make([]string, 25),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  10,
				Height: 5,
			},
			make([]string, 49),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
			},
			make([]string, 50),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  25,
				Height: 2,
			},
			make([]string, 51),
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
			},
			make([]string, 1),
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
			},
			make([]string, 8),
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
			},
			make([]string, 8),
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
			},
			make([]string, 17),
			ErrorTooManySnakes,
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
			},
			make([]string, 8),
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
			},
			make([]string, 17),
			ErrorTooManySnakes,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.BoardState.Width, test.BoardState.Height, len(test.SnakeIDs)), func(t *testing.T) {
			require.Equal(t, test.BoardState.Width*test.BoardState.Height, len(GetUnoccupiedPoints(test.BoardState, true, false)))
			err := PlaceSnakesAutomatically(MaxRand, test.BoardState, test.SnakeIDs)
			require.Equal(t, test.Err, err, "Snakes: %d", len(test.BoardState.Snakes))
			if err == nil {
				for i := 0; i < len(test.BoardState.Snakes); i++ {
					require.Len(t, test.BoardState.Snakes[i].Body, 3)
					for _, point := range test.BoardState.Snakes[i].Body {
						require.GreaterOrEqual(t, point.X, 0)
						require.GreaterOrEqual(t, point.Y, 0)
						require.Less(t, point.X, test.BoardState.Width)
						require.Less(t, point.Y, test.BoardState.Height)
					}

					for j := 0; j < len(test.BoardState.Snakes); j++ {
						if j == i {
							continue
						}
						require.NotEqual(t, test.BoardState.Snakes[j].Body[0], test.BoardState.Snakes[i].Body[0], "Snakes placed at same square")
					}

					// All snakes are expected to be placed on an even square - this is true even of fixed positions for known board sizes
					var snakePlacedOnEvenSquare bool = ((test.BoardState.Snakes[i].Body[0].X + test.BoardState.Snakes[i].Body[0].Y) % 2) == 0
					require.Equal(t, true, snakePlacedOnEvenSquare)
				}
			}
		})
	}
}

func TestPlaceSnakesFixed(t *testing.T) {
	snakeIDs := make([]string, 8)

	for _, test := range []struct {
		label              string
		rand               Rand
		expectedSnakeHeads []Point
	}{
		{
			label: "corners before cardinal directions",
			rand:  MinRand,
			expectedSnakeHeads: []Point{
				{X: 1, Y: 1},
				{X: 1, Y: 9},
				{X: 9, Y: 1},
				{X: 9, Y: 9},

				{X: 1, Y: 5},
				{X: 5, Y: 1},
				{X: 5, Y: 9},
				{X: 9, Y: 5},
			},
		},
		{
			label: "cardinal directions before corners",
			rand:  MaxRand,
			expectedSnakeHeads: []Point{
				{X: 5, Y: 1},
				{X: 5, Y: 9},
				{X: 9, Y: 5},
				{X: 1, Y: 5},

				{X: 1, Y: 9},
				{X: 9, Y: 1},
				{X: 9, Y: 9},
				{X: 1, Y: 1},
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			boardState := &BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
			}

			err := PlaceSnakesAutomatically(test.rand, boardState, snakeIDs)
			require.NoError(t, err)

			var snakeHeads []Point
			for _, snake := range boardState.Snakes {
				require.Len(t, snake.Body, 3)
				snakeHeads = append(snakeHeads, snake.Body[0])
			}
			require.Equalf(t, test.expectedSnakeHeads, snakeHeads, "%#v", snakeHeads)
		})
	}
}

func TestPlaceSnake(t *testing.T) {
	// TODO: Should PlaceSnake check for boundaries?
	boardState := NewBoardState(BoardSizeSmall, BoardSizeSmall)
	require.Empty(t, boardState.Snakes)

	_ = PlaceSnake(boardState, "a", []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}})

	require.Len(t, boardState.Snakes, 1)
	require.Equal(t, Snake{
		ID:              "a",
		Body:            []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}},
		Health:          SnakeMaxHealth,
		EliminatedCause: NotEliminated,
		EliminatedBy:    "",
	}, boardState.Snakes[0])

	_ = PlaceSnake(boardState, "b", []Point{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 3, Y: 2}})

	require.Len(t, boardState.Snakes, 2)
	require.Equal(t, Snake{
		ID:              "b",
		Body:            []Point{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 3, Y: 2}},
		Health:          SnakeMaxHealth,
		EliminatedCause: NotEliminated,
		EliminatedBy:    "",
	}, boardState.Snakes[1])
}

func TestPlaceFood(t *testing.T) {
	tests := []struct {
		BoardState   *BoardState
		ExpectedFood int
	}{
		{
			&BoardState{
				Width:  1,
				Height: 1,
				Snakes: make([]Snake, 1),
			},
			1,
		},
		{
			&BoardState{
				Width:  1,
				Height: 2,
				Snakes: make([]Snake, 2),
			},
			2,
		},
		{
			&BoardState{
				Width:  101,
				Height: 202,
				Snakes: make([]Snake, 17),
			},
			17,
		},
		{
			&BoardState{
				Width:  10,
				Height: 20,
				Snakes: make([]Snake, 305),
			},
			200,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
				Snakes: []Snake{
					{Body: []Point{{X: 5, Y: 1}}},
					{Body: []Point{{X: 5, Y: 3}}},
					{Body: []Point{{X: 5, Y: 5}}},
				},
			},
			4, // +1 because of fixed spawn locations
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: []Snake{
					{Body: []Point{{X: 1, Y: 1}}},
					{Body: []Point{{X: 1, Y: 5}}},
					{Body: []Point{{X: 1, Y: 9}}},
					{Body: []Point{{X: 5, Y: 1}}},
					{Body: []Point{{X: 5, Y: 9}}},
					{Body: []Point{{X: 9, Y: 1}}},
					{Body: []Point{{X: 9, Y: 5}}},
					{Body: []Point{{X: 9, Y: 9}}},
				},
			},
			9, // +1 because of fixed spawn locations
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: []Snake{
					{Body: []Point{{X: 1, Y: 1}}},
					{Body: []Point{{X: 1, Y: 9}}},
					{Body: []Point{{X: 1, Y: 17}}},
					{Body: []Point{{X: 17, Y: 1}}},
					{Body: []Point{{X: 17, Y: 9}}},
					{Body: []Point{{X: 17, Y: 17}}},
				},
			},
			7, // +1 because of fixed spawn locations
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {

			require.Len(t, test.BoardState.Food, 0)
			err := PlaceFoodAutomatically(MaxRand, test.BoardState)
			require.NoError(t, err)
			require.Equal(t, test.ExpectedFood, len(test.BoardState.Food))
			for _, point := range test.BoardState.Food {
				require.GreaterOrEqual(t, point.X, 0)
				require.GreaterOrEqual(t, point.Y, 0)
				require.Less(t, point.X, test.BoardState.Width)
				require.Less(t, point.Y, test.BoardState.Height)
			}
		})
	}
}

func TestPlaceFoodFixed(t *testing.T) {
	tests := []struct {
		BoardState *BoardState
	}{
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
				Snakes: []Snake{
					{Body: []Point{{X: 1, Y: 3}}},
				},
			},
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: []Snake{
					{Body: []Point{{X: 1, Y: 1}}},
					{Body: []Point{{X: 1, Y: 5}}},
					{Body: []Point{{X: 9, Y: 5}}},
					{Body: []Point{{X: 9, Y: 9}}},
				},
			},
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: []Snake{
					{Body: []Point{{X: 1, Y: 1}}},
					{Body: []Point{{X: 1, Y: 9}}},
					{Body: []Point{{X: 1, Y: 17}}},
					{Body: []Point{{X: 9, Y: 1}}},
					{Body: []Point{{X: 9, Y: 17}}},
					{Body: []Point{{X: 17, Y: 1}}},
					{Body: []Point{{X: 17, Y: 9}}},
					{Body: []Point{{X: 17, Y: 17}}},
				},
			},
		},
	}

	for _, test := range tests {
		require.Len(t, test.BoardState.Food, 0)

		err := PlaceFoodFixed(MaxRand, test.BoardState)
		require.NoError(t, err)
		require.Equal(t, len(test.BoardState.Snakes)+1, len(test.BoardState.Food))

		midPoint := Point{X: (test.BoardState.Width - 1) / 2, Y: (test.BoardState.Height - 1) / 2}

		// Make sure every snake has food within 2 moves of it
		for _, snake := range test.BoardState.Snakes {
			head := snake.Body[0]

			bottomLeft := Point{X: head.X - 1, Y: head.Y - 1}
			topLeft := Point{X: head.X - 1, Y: head.Y + 1}
			bottomRight := Point{X: head.X + 1, Y: head.Y - 1}
			topRight := Point{X: head.X + 1, Y: head.Y + 1}

			foundFoodInTwoMoves := false
			for _, food := range test.BoardState.Food {
				if food == bottomLeft || food == topLeft || food == bottomRight || food == topRight {
					foundFoodInTwoMoves = true
					// Ensure it's not closer to the center than snake
					require.True(t, getDistanceBetweenPoints(head, midPoint) <= getDistanceBetweenPoints(food, midPoint))
					break
				}
			}
			require.True(t, foundFoodInTwoMoves)
		}

		// Make sure one food exists in center of board
		foundFoodInCenter := false
		for _, food := range test.BoardState.Food {
			if food == midPoint {
				foundFoodInCenter = true
				break
			}
		}
		require.True(t, foundFoodInCenter)
	}
}

func TestPlaceFoodFixedNoRoom(t *testing.T) {
	boardState := &BoardState{
		Width:  3,
		Height: 3,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 1}}},
		},
		Food: []Point{},
	}
	err := PlaceFoodFixed(MaxRand, boardState)
	require.Error(t, err)
}

func TestPlaceFoodFixedNoRoom_Corners(t *testing.T) {
	boardState := &BoardState{
		Width:  7,
		Height: 7,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 1}}},
			{Body: []Point{{X: 1, Y: 5}}},
			{Body: []Point{{X: 5, Y: 1}}},
			{Body: []Point{{X: 5, Y: 5}}},
		},
		Food: []Point{},
	}

	// There are only two possible food spawn locations for each snake,
	// so repeat calls to place food should fail after 2 successes
	err := PlaceFoodFixed(MaxRand, boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 4, len(boardState.Food))

	err = PlaceFoodFixed(MaxRand, boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 8, len(boardState.Food))

	// And now there should be no more room.
	err = PlaceFoodFixed(MaxRand, boardState)
	require.Error(t, err)

	expectedFood := []Point{
		{X: 0, Y: 2}, {X: 2, Y: 0}, // Snake @ {X: 1, Y: 1}
		{X: 0, Y: 4}, {X: 2, Y: 6}, // Snake @ {X: 1, Y: 5}
		{X: 4, Y: 0}, {X: 6, Y: 2}, // Snake @ {X: 5, Y: 1}
		{X: 4, Y: 6}, {X: 6, Y: 4}, // Snake @ {X: 5, Y: 5}
	}
	sortPoints(expectedFood)
	sortPoints(boardState.Food)
	require.Equal(t, expectedFood, boardState.Food)
}

func TestPlaceFoodFixedNoRoom_Cardinal(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 5}}},
			{Body: []Point{{X: 5, Y: 1}}},
			{Body: []Point{{X: 5, Y: 9}}},
			{Body: []Point{{X: 9, Y: 5}}},
		},
		Food: []Point{},
	}

	// There are only two possible spawn locations for each snake,
	// so repeat calls to place food should fail after 2 successes
	err := PlaceFoodFixed(MaxRand, boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 4, len(boardState.Food))

	err = PlaceFoodFixed(MaxRand, boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 8, len(boardState.Food))

	// And now there should be no more room.
	err = PlaceFoodFixed(MaxRand, boardState)
	require.Error(t, err)

	expectedFood := []Point{
		{X: 0, Y: 4}, {X: 0, Y: 6}, // Snake @ {X: 1, Y: 5}
		{X: 4, Y: 0}, {X: 6, Y: 0}, // Snake @ {X: 5, Y: 1}
		{X: 4, Y: 10}, {X: 6, Y: 10}, // Snake @ {X: 5, Y: 9}
		{X: 10, Y: 4}, {X: 10, Y: 6}, // Snake @ {X: 9, Y: 5}
	}
	sortPoints(expectedFood)
	sortPoints(boardState.Food)
	require.Equal(t, expectedFood, boardState.Food)
}

func TestGetDistanceBetweenPoints(t *testing.T) {
	tests := []struct {
		A        Point
		B        Point
		Expected int
	}{
		{Point{X: 0, Y: 0}, Point{X: 0, Y: 0}, 0},
		{Point{X: 0, Y: 0}, Point{X: 1, Y: 0}, 1},
		{Point{X: 0, Y: 0}, Point{X: 0, Y: 1}, 1},
		{Point{X: 0, Y: 0}, Point{X: 1, Y: 1}, 2},
		{Point{X: 0, Y: 0}, Point{X: 4, Y: 4}, 8},
		{Point{X: 0, Y: 0}, Point{X: 4, Y: 6}, 10},
		{Point{X: 8, Y: 0}, Point{X: 8, Y: 0}, 0},
		{Point{X: 8, Y: 0}, Point{X: 8, Y: 8}, 8},
		{Point{X: 8, Y: 0}, Point{X: 0, Y: 8}, 16},
	}

	for _, test := range tests {
		require.Equal(t, getDistanceBetweenPoints(test.A, test.B), test.Expected)
		require.Equal(t, getDistanceBetweenPoints(test.B, test.A), test.Expected)
	}
}

func TestIsSquareBoard(t *testing.T) {
	tests := []struct {
		Width    int
		Height   int
		Expected bool
	}{
		{1, 1, true},
		{0, 0, true},
		{0, 45, false},
		{45, 1, false},
		{7, 7, true},
		{11, 11, true},
		{19, 19, true},
		{7, 11, false},
		{11, 19, false},
		{19, 7, false},
	}

	for _, test := range tests {
		result := isSquareBoard(&BoardState{Width: test.Width, Height: test.Height})
		require.Equal(t, test.Expected, result)
	}
}

func TestGetUnoccupiedPoints(t *testing.T) {
	tests := []struct {
		Board    *BoardState
		Expected []Point
	}{
		{
			&BoardState{
				Height: 1,
				Width:  1,
			},
			[]Point{{X: 0, Y: 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  2,
			},
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  1,
				Food:   []Point{{X: 0, Y: 0}, {X: 101, Y: 202}, {X: -4, Y: -5}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
			},
			[]Point{{X: 0, Y: 1}, {X: 1, Y: 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 4,
				Width:  1,
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}}},
				},
			},
			[]Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}},
				},
			},
			[]Point{{X: 0, Y: 1}, {X: 2, Y: 0}, {X: 2, Y: 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 0}},
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}},
					{Body: []Point{{X: 0, Y: 1}}},
				},
			},
			[]Point{{X: 2, Y: 1}},
		},
		{
			&BoardState{
				Height:  1,
				Width:   1,
				Hazards: []Point{{X: 0, Y: 0}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height:  2,
				Width:   2,
				Hazards: []Point{{X: 1, Y: 1}},
			},
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []Point{{X: 1, Y: 1}, {X: 2, Y: 0}},
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}},
					{Body: []Point{{X: 0, Y: 1}}},
				},
				Hazards: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
			},
			[]Point{{X: 2, Y: 1}},
		},
	}

	for _, test := range tests {
		unoccupiedPoints := GetUnoccupiedPoints(test.Board, true, true)
		require.Equal(t, len(test.Expected), len(unoccupiedPoints))
		for i, e := range test.Expected {
			require.Equal(t, e, unoccupiedPoints[i])
		}
	}
}

func TestGetEvenUnoccupiedPoints(t *testing.T) {
	tests := []struct {
		Board    *BoardState
		Expected []Point
	}{
		{
			&BoardState{
				Height: 1,
				Width:  1,
			},
			[]Point{{X: 0, Y: 0}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
			},
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 1}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  1,
				Food:   []Point{{X: 0, Y: 0}, {X: 101, Y: 202}, {X: -4, Y: -5}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
			},
			[]Point{{X: 1, Y: 1}},
		},
		{
			&BoardState{
				Height: 4,
				Width:  4,
				Food:   []Point{{X: 0, Y: 0}, {X: 0, Y: 2}, {X: 1, Y: 1}, {X: 1, Y: 3}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 3, Y: 1}, {X: 3, Y: 3}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 4,
				Width:  1,
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}}},
				},
			},
			[]Point{{X: 0, Y: 2}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}},
				},
			},
			[]Point{{X: 2, Y: 0}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 1}},
				Snakes: []Snake{
					{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}},
					{Body: []Point{{X: 0, Y: 1}}},
				},
			},
			[]Point{{X: 2, Y: 0}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			evenUnoccupiedPoints := GetEvenUnoccupiedPoints(test.Board)
			require.Equal(t, len(test.Expected), len(evenUnoccupiedPoints))
			for i, e := range test.Expected {
				require.Equal(t, e, evenUnoccupiedPoints[i])
			}
		})
	}
}

func TestPlaceFoodRandomly(t *testing.T) {
	b := &BoardState{
		Height: 1,
		Width:  3,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 0}}},
		},
	}
	// Food should never spawn, no room
	err := PlaceFoodRandomly(MaxRand, b, 99)
	require.NoError(t, err)
	require.Equal(t, len(b.Food), 0)
}

func TestEliminateSnake(t *testing.T) {
	s := &Snake{}
	EliminateSnake(s, "test-cause", "", 2)
	require.Equal(t, "test-cause", s.EliminatedCause)
	require.Equal(t, "", s.EliminatedBy)
	require.Equal(t, 2, s.EliminatedOnTurn)
}
