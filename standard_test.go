package rules

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandardRulesetInterface(t *testing.T) {
	var _ Ruleset = (*StandardRuleset)(nil)
}

func TestSanity(t *testing.T) {
	r := StandardRuleset{}

	state, err := r.CreateInitialBoardState(0, 0, []string{})
	require.NoError(t, err)
	require.NotNil(t, state)
	require.Equal(t, int32(0), state.Width)
	require.Equal(t, int32(0), state.Height)
	require.Len(t, state.Food, 0)
	require.Len(t, state.Snakes, 0)

	next, err := r.CreateNextBoardState(
		&BoardState{},
		[]SnakeMove{},
	)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.Equal(t, int32(0), state.Width)
	require.Equal(t, int32(0), state.Height)
	require.Len(t, state.Snakes, 0)
}

func TestStandardName(t *testing.T) {
	r := StandardRuleset{}
	require.Equal(t, "standard", r.Name())
}

func TestCreateInitialBoardState(t *testing.T) {
	tests := []struct {
		Height          int32
		Width           int32
		IDs             []string
		ExpectedNumFood int
		Err             error
	}{
		{1, 1, []string{"one"}, 0, nil},
		{1, 2, []string{"one"}, 0, nil},
		{1, 4, []string{"one"}, 1, nil},
		{2, 2, []string{"one"}, 1, nil},
		{9, 8, []string{"one"}, 1, nil},
		{2, 2, []string{"one", "two"}, 0, nil},
		{1, 1, []string{"one", "two"}, 2, ErrorNoRoomForSnake},
		{1, 2, []string{"one", "two"}, 2, ErrorNoRoomForSnake},
		{BoardSizeSmall, BoardSizeSmall, []string{"one", "two"}, 3, nil},
	}

	r := StandardRuleset{}
	for testNum, test := range tests {
		state, err := r.CreateInitialBoardState(test.Width, test.Height, test.IDs)
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
	}
}

func TestPlaceSnakes(t *testing.T) {
	// Because placement is random, we only test to ensure
	// that snake bodies are populated correctly
	// Note: because snakes are randomly spawned on even diagonal points, the board can accomodate number of snakes equal to: width*height/2
	tests := []struct {
		BoardState *BoardState
		Err        error
	}{
		{
			&BoardState{
				Width:  1,
				Height: 1,
				Snakes: make([]Snake, 1),
			},
			nil,
		},
		{
			&BoardState{
				Width:  1,
				Height: 1,
				Snakes: make([]Snake, 2),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  2,
				Height: 1,
				Snakes: make([]Snake, 2),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  1,
				Height: 2,
				Snakes: make([]Snake, 2),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  10,
				Height: 5,
				Snakes: make([]Snake, 24),
			},
			nil,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
				Snakes: make([]Snake, 25),
			},
			nil,
		},
		{
			&BoardState{
				Width:  10,
				Height: 5,
				Snakes: make([]Snake, 49),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
				Snakes: make([]Snake, 50),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  25,
				Height: 2,
				Snakes: make([]Snake, 51),
			},
			ErrorNoRoomForSnake,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
				Snakes: make([]Snake, 1),
			},
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
				Snakes: make([]Snake, 8),
			},
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeSmall,
				Height: BoardSizeSmall,
				Snakes: make([]Snake, 9),
			},
			ErrorTooManySnakes,
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: make([]Snake, 8),
			},
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: make([]Snake, 9),
			},
			ErrorTooManySnakes,
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: make([]Snake, 8),
			},
			nil,
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: make([]Snake, 9),
			},
			ErrorTooManySnakes,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		require.Equal(t, test.BoardState.Width*test.BoardState.Height, int32(len(r.getUnoccupiedPoints(test.BoardState, true))))
		err := r.placeSnakes(test.BoardState)
		require.Equal(t, test.Err, err, "Snakes: %d", len(test.BoardState.Snakes))
		if err == nil {
			for i := 0; i < len(test.BoardState.Snakes); i++ {
				require.Len(t, test.BoardState.Snakes[i].Body, 3)
				for _, point := range test.BoardState.Snakes[i].Body {
					require.GreaterOrEqual(t, point.X, int32(0))
					require.GreaterOrEqual(t, point.Y, int32(0))
					require.Less(t, point.X, test.BoardState.Width)
					require.Less(t, point.Y, test.BoardState.Height)
				}

				// All snakes are expected to be placed on an even square - this is true even of fixed positions for known board sizes
				var snakePlacedOnEvenSquare bool = ((test.BoardState.Snakes[i].Body[0].X + test.BoardState.Snakes[i].Body[0].Y) % 2) == 0
				require.Equal(t, true, snakePlacedOnEvenSquare)
			}
		}
	}
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
					{Body: []Point{{5, 1}}},
					{Body: []Point{{5, 3}}},
					{Body: []Point{{5, 5}}},
				},
			},
			4, // +1 because of fixed spawn locations
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: []Snake{
					{Body: []Point{{1, 1}}},
					{Body: []Point{{1, 5}}},
					{Body: []Point{{1, 9}}},
					{Body: []Point{{5, 1}}},
					{Body: []Point{{5, 9}}},
					{Body: []Point{{9, 1}}},
					{Body: []Point{{9, 5}}},
					{Body: []Point{{9, 9}}},
				},
			},
			9, // +1 because of fixed spawn locations
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: []Snake{
					{Body: []Point{{1, 1}}},
					{Body: []Point{{1, 9}}},
					{Body: []Point{{1, 17}}},
					{Body: []Point{{17, 1}}},
					{Body: []Point{{17, 9}}},
					{Body: []Point{{17, 17}}},
				},
			},
			7, // +1 because of fixed spawn locations
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		require.Len(t, test.BoardState.Food, 0)
		err := r.placeFood(test.BoardState)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedFood, len(test.BoardState.Food))
		for _, point := range test.BoardState.Food {
			require.GreaterOrEqual(t, point.X, int32(0))
			require.GreaterOrEqual(t, point.Y, int32(0))
			require.Less(t, point.X, test.BoardState.Width)
			require.Less(t, point.Y, test.BoardState.Height)
		}
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
					{Body: []Point{{1, 3}}},
				},
			},
		},
		{
			&BoardState{
				Width:  BoardSizeMedium,
				Height: BoardSizeMedium,
				Snakes: []Snake{
					{Body: []Point{{1, 1}}},
					{Body: []Point{{1, 5}}},
					{Body: []Point{{9, 5}}},
					{Body: []Point{{9, 9}}},
				},
			},
		},
		{
			&BoardState{
				Width:  BoardSizeLarge,
				Height: BoardSizeLarge,
				Snakes: []Snake{
					{Body: []Point{{1, 1}}},
					{Body: []Point{{1, 9}}},
					{Body: []Point{{1, 17}}},
					{Body: []Point{{9, 1}}},
					{Body: []Point{{9, 17}}},
					{Body: []Point{{17, 1}}},
					{Body: []Point{{17, 9}}},
					{Body: []Point{{17, 17}}},
				},
			},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		require.Len(t, test.BoardState.Food, 0)

		err := r.placeFoodFixed(test.BoardState)
		require.NoError(t, err)
		require.Equal(t, len(test.BoardState.Snakes)+1, len(test.BoardState.Food))

		// Make sure every snake has food within 2 moves of it
		for _, snake := range test.BoardState.Snakes {
			head := snake.Body[0]

			bottomLeft := Point{head.X - 1, head.Y - 1}
			topLeft := Point{head.X - 1, head.Y + 1}
			bottomRight := Point{head.X + 1, head.Y - 1}
			topRight := Point{head.X + 1, head.Y + 1}

			foundFoodInTwoMoves := false
			for _, food := range test.BoardState.Food {
				if food == bottomLeft || food == topLeft || food == bottomRight || food == topRight {
					foundFoodInTwoMoves = true
					break
				}
			}
			require.True(t, foundFoodInTwoMoves)
		}

		// Make sure one food exists in center of board
		foundFoodInCenter := false
		midPoint := Point{(test.BoardState.Width - 1) / 2, (test.BoardState.Height - 1) / 2}
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
	r := StandardRuleset{}

	boardState := &BoardState{
		Width:  3,
		Height: 3,
		Snakes: []Snake{
			{Body: []Point{{1, 1}}},
		},
		Food: []Point{},
	}
	err := r.placeFoodFixed(boardState)
	require.Error(t, err)

	boardState = &BoardState{
		Width:  7,
		Height: 7,
		Snakes: []Snake{
			{Body: []Point{{1, 1}}},
		},
		Food: []Point{},
	}
	err = r.placeFoodFixed(boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 1, len(boardState.Food))

	err = r.placeFoodFixed(boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 2, len(boardState.Food))

	err = r.placeFoodFixed(boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 3, len(boardState.Food))

	err = r.placeFoodFixed(boardState)
	require.NoError(t, err)
	boardState.Food = boardState.Food[:len(boardState.Food)-1] // Center food
	require.Equal(t, 4, len(boardState.Food))

	// And now there should be no more room.
	err = r.placeFoodFixed(boardState)
	require.Error(t, err)
}

func TestCreateNextBoardState(t *testing.T) {
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{1, 1}, {1, 2}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{3, 4}, {3, 3}},
						Health: 100,
					},
				},
				Food: []Point{{0, 0}, {1, 0}},
			},
			[]SnakeMove{},
			ErrorNoMoveFound,
			nil,
		},
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{1, 1}, {1, 2}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{},
						Health: 100,
					},
				},
				Food: []Point{{0, 0}, {1, 0}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			ErrorZeroLengthSnake,
			nil,
		},
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{1, 1}, {1, 2}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{3, 4}, {3, 3}},
						Health: 100,
					},
					{
						ID:              "three",
						Body:            []Point{},
						Health:          100,
						EliminatedCause: EliminatedByOutOfBounds,
					},
				},
				Food: []Point{{0, 0}, {1, 0}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveDown},
				{ID: "two", Move: MoveUp},
				{ID: "three", Move: MoveLeft}, // Should be ignored
			},
			nil,
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{1, 0}, {1, 1}, {1, 1}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{3, 5}, {3, 4}},
						Health: 99,
					},
					{
						ID:              "three",
						Body:            []Point{},
						Health:          100,
						EliminatedCause: EliminatedByOutOfBounds,
					},
				},
				Food: []Point{{0, 0}},
			},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.Equal(t, test.expectedError, err)
		require.Equal(t, test.expectedState, nextState)
	}
}

func TestEatingOnLastMove(t *testing.T) {
	// We want to specifically ensure that snakes eating food on their last turn
	// survive. It used to be that this wasn't the case, and snakes were eliminated
	// if they moved onto food with their final move. This behaviour wasn't "wrong" or incorrect,
	// it just was less fun to watch. So let's ensure we're always giving snakes every possible
	// changes to reach food before eliminating them.
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 2}, {0, 1}, {0, 0}},
						Health: 1,
					},
					{
						ID:     "two",
						Body:   []Point{{3, 2}, {3, 3}, {3, 4}},
						Health: 1,
					},
				},
				Food: []Point{{0, 3}, {9, 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			nil,
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 3}, {0, 2}, {0, 1}, {0, 1}},
						Health: 100,
					},
					{
						ID:              "two",
						Body:            []Point{{3, 1}, {3, 2}, {3, 3}},
						Health:          0,
						EliminatedCause: EliminatedByOutOfHealth,
					},
				},
				Food: []Point{{9, 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := StandardRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		require.Equal(t, nextState, test.expectedState)
	}
}

func TestHeadToHeadOnFood(t *testing.T) {
	// We want to specifically ensure that snakes that collide head-to-head
	// on top of food successfully remove the food - that's the core behaviour this test
	// is enforcing. There's a known side effect of this though, in that both snakes will
	// have eaten prior to being evaluated on the head-to-head (+1 length, full health).
	// We're okay with that since it does not impact the result of the head-to-head,
	// however that behaviour could change in the future and this test could be updated.
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 2}, {0, 1}, {0, 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []Point{{0, 4}, {0, 5}, {0, 6}},
						Health: 10,
					},
				},
				Food: []Point{{0, 3}, {9, 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			nil,
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:              "one",
						Body:            []Point{{0, 3}, {0, 2}, {0, 1}, {0, 1}},
						Health:          100,
						EliminatedCause: EliminatedByHeadToHeadCollision,
						EliminatedBy:    "two",
					},
					{
						ID:              "two",
						Body:            []Point{{0, 3}, {0, 4}, {0, 5}, {0, 5}},
						Health:          100,
						EliminatedCause: EliminatedByHeadToHeadCollision,
						EliminatedBy:    "one",
					},
				},
				Food: []Point{{9, 9}},
			},
		},
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 2}, {0, 1}, {0, 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []Point{{0, 4}, {0, 5}, {0, 6}, {0, 7}},
						Health: 10,
					},
				},
				Food: []Point{{0, 3}, {9, 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			nil,
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:              "one",
						Body:            []Point{{0, 3}, {0, 2}, {0, 1}, {0, 1}},
						Health:          100,
						EliminatedCause: EliminatedByHeadToHeadCollision,
						EliminatedBy:    "two",
					},
					{
						ID:     "two",
						Body:   []Point{{0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 6}},
						Health: 100,
					},
				},
				Food: []Point{{9, 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := StandardRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.Equal(t, test.expectedError, err)
		require.Equal(t, test.expectedState, nextState)
	}
}

func TestRegressionIssue19(t *testing.T) {
	// Eliminated snakes passed to CreateNextBoardState should not impact next game state
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 2}, {0, 1}, {0, 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{0, 5}, {0, 6}, {0, 7}},
						Health: 100,
					},
					{
						ID:              "eliminated",
						Body:            []Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}},
						Health:          0,
						EliminatedCause: EliminatedByOutOfHealth,
					},
				},
				Food: []Point{{9, 9}},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			nil,
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 3}, {0, 2}, {0, 1}},
						Health: 99,
					},
					{
						ID:     "two",
						Body:   []Point{{0, 4}, {0, 5}, {0, 6}},
						Health: 99,
					},
					{
						ID:              "eliminated",
						Body:            []Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}},
						Health:          0,
						EliminatedCause: EliminatedByOutOfHealth,
					},
				},
				Food: []Point{{9, 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := StandardRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		require.Equal(t, nextState, test.expectedState)
	}

}

func TestMoveSnakes(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{10, 110}, {11, 110}},
				Health: 111111,
			},
			{
				ID:     "two",
				Body:   []Point{{23, 220}, {22, 220}, {21, 220}, {20, 220}},
				Health: 222222,
			},
			{
				ID:              "three",
				Body:            []Point{{0, 0}},
				Health:          1,
				EliminatedCause: EliminatedByOutOfBounds,
			},
		},
	}

	tests := []struct {
		MoveOne       string
		ExpectedOne   []Point
		MoveTwo       string
		ExpectedTwo   []Point
		MoveThree     string
		ExpectedThree []Point
	}{
		{
			MoveDown, []Point{{10, 109}, {10, 110}},
			MoveUp, []Point{{23, 221}, {23, 220}, {22, 220}, {21, 220}},
			MoveDown, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{11, 109}, {10, 109}},
			MoveLeft, []Point{{22, 221}, {23, 221}, {23, 220}, {22, 220}},
			MoveDown, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{12, 109}, {11, 109}},
			MoveLeft, []Point{{21, 221}, {22, 221}, {23, 221}, {23, 220}},
			MoveDown, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{13, 109}, {12, 109}},
			MoveLeft, []Point{{20, 221}, {21, 221}, {22, 221}, {23, 221}},
			MoveDown, []Point{{0, 0}},
		},
		{
			MoveDown, []Point{{13, 108}, {13, 109}},
			MoveUp, []Point{{20, 222}, {20, 221}, {21, 221}, {22, 221}},
			MoveDown, []Point{{0, 0}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		moves := []SnakeMove{
			{ID: "one", Move: test.MoveOne},
			{ID: "two", Move: test.MoveTwo},
			{ID: "three", Move: test.MoveThree},
		}
		err := r.moveSnakes(b, moves)

		require.NoError(t, err)
		require.Len(t, b.Snakes, 3)

		require.Equal(t, int32(111111), b.Snakes[0].Health)
		require.Equal(t, int32(222222), b.Snakes[1].Health)
		require.Equal(t, int32(1), b.Snakes[2].Health)

		require.Len(t, b.Snakes[0].Body, 2)
		require.Len(t, b.Snakes[1].Body, 4)
		require.Len(t, b.Snakes[2].Body, 1)

		require.Equal(t, len(b.Snakes[0].Body), len(test.ExpectedOne))
		for i, e := range test.ExpectedOne {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
		require.Equal(t, len(b.Snakes[1].Body), len(test.ExpectedTwo))
		for i, e := range test.ExpectedTwo {
			require.Equal(t, e, b.Snakes[1].Body[i])
		}
		require.Equal(t, len(b.Snakes[2].Body), len(test.ExpectedThree))
		for i, e := range test.ExpectedThree {
			require.Equal(t, e, b.Snakes[2].Body[i])
		}
	}
}

func TestMoveSnakesWrongID(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{1, 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "not found",
			Move: MoveUp,
		},
	}

	r := StandardRuleset{}
	err := r.moveSnakes(b, moves)
	require.Equal(t, ErrorNoMoveFound, err)
}

func TestMoveSnakesNotEnoughMoves(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{1, 1}},
			},
			{
				ID:   "two",
				Body: []Point{{2, 2}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "two",
			Move: MoveUp,
		},
	}

	r := StandardRuleset{}
	err := r.moveSnakes(b, moves)
	require.Equal(t, ErrorNoMoveFound, err)
}

func TestMoveSnakesExtraMovesIgnored(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{1, 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "one",
			Move: MoveDown,
		},
		{
			ID:   "two",
			Move: MoveLeft,
		},
	}

	r := StandardRuleset{}
	err := r.moveSnakes(b, moves)
	require.NoError(t, err)
	require.Equal(t, []Point{{1, 0}}, b.Snakes[0].Body)
}

func TestIsKnownBoardSize(t *testing.T) {
	tests := []struct {
		Width    int32
		Height   int32
		Expected bool
	}{
		{1, 1, false},
		{0, 0, false},
		{0, 45, false},
		{45, 1, false},
		{7, 7, true},
		{11, 11, true},
		{19, 19, true},
		{7, 11, false},
		{11, 19, false},
		{19, 7, false},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		result := r.isKnownBoardSize(&BoardState{Width: test.Width, Height: test.Height})
		require.Equal(t, test.Expected, result)
	}
}

func TestMoveSnakesDefault(t *testing.T) {
	tests := []struct {
		Body     []Point
		Move     string
		Expected []Point
	}{
		{
			Body:     []Point{{0, 0}},
			Move:     "invalid",
			Expected: []Point{{0, 1}},
		},
		{
			Body:     []Point{{5, 5}, {5, 5}},
			Move:     "",
			Expected: []Point{{5, 6}, {5, 5}},
		},
		{
			Body:     []Point{{5, 5}, {5, 4}},
			Expected: []Point{{5, 6}, {5, 5}},
		},
		{
			Body:     []Point{{5, 4}, {5, 5}},
			Expected: []Point{{5, 3}, {5, 4}},
		},
		{
			Body:     []Point{{5, 4}, {5, 5}},
			Expected: []Point{{5, 3}, {5, 4}},
		},
		{
			Body:     []Point{{4, 5}, {5, 5}},
			Expected: []Point{{3, 5}, {4, 5}},
		},
		{
			Body:     []Point{{5, 5}, {4, 5}},
			Expected: []Point{{6, 5}, {5, 5}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Snakes: []Snake{
				{ID: "one", Body: test.Body},
			},
		}
		moves := []SnakeMove{{ID: "one", Move: test.Move}}

		err := r.moveSnakes(b, moves)
		require.NoError(t, err)
		require.Len(t, b.Snakes, 1)
		require.Equal(t, len(test.Body), len(b.Snakes[0].Body))
		require.Equal(t, len(test.Expected), len(b.Snakes[0].Body))
		for i, e := range test.Expected {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
	}
}

func TestReduceSnakeHealth(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				Body:   []Point{{0, 0}, {0, 1}},
				Health: 99,
			},
			{
				Body:   []Point{{5, 8}, {6, 8}, {7, 8}},
				Health: 2,
			},
			{
				Body:            []Point{{0, 0}, {0, 1}},
				Health:          50,
				EliminatedCause: EliminatedByCollision,
			},
		},
	}

	r := StandardRuleset{}
	err := r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(98))
	require.Equal(t, b.Snakes[1].Health, int32(1))
	require.Equal(t, b.Snakes[2].Health, int32(50))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(97))
	require.Equal(t, b.Snakes[1].Health, int32(0))
	require.Equal(t, b.Snakes[2].Health, int32(50))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(96))
	require.Equal(t, b.Snakes[1].Health, int32(-1))
	require.Equal(t, b.Snakes[2].Health, int32(50))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(95))
	require.Equal(t, b.Snakes[1].Health, int32(-2))
	require.Equal(t, b.Snakes[2].Health, int32(50))
}

func TestSnakeIsOutOfHealth(t *testing.T) {
	tests := []struct {
		Health   int32
		Expected bool
	}{
		{Health: math.MinInt32, Expected: true},
		{Health: -10, Expected: true},
		{Health: -2, Expected: true},
		{Health: -1, Expected: true},
		{Health: 0, Expected: true},
		{Health: 1, Expected: false},
		{Health: 2, Expected: false},
		{Health: 10, Expected: false},
		{Health: math.MaxInt32, Expected: false},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		s := &Snake{Health: test.Health}
		require.Equal(t, test.Expected, r.snakeIsOutOfHealth(s), "Health: %+v", test.Health)
	}
}

func TestSnakeIsOutOfBounds(t *testing.T) {
	boardWidth := int32(10)
	boardHeight := int32(100)

	tests := []struct {
		Point    Point
		Expected bool
	}{
		{Point{X: math.MinInt32, Y: math.MinInt32}, true},
		{Point{X: math.MinInt32, Y: 0}, true},
		{Point{X: 0, Y: math.MinInt32}, true},
		{Point{X: -1, Y: -1}, true},
		{Point{X: -1, Y: 0}, true},
		{Point{X: 0, Y: -1}, true},
		{Point{X: 0, Y: 0}, false},
		{Point{X: 1, Y: 0}, false},
		{Point{X: 0, Y: 1}, false},
		{Point{X: 1, Y: 1}, false},
		{Point{X: 9, Y: 9}, false},
		{Point{X: 9, Y: 10}, false},
		{Point{X: 9, Y: 11}, false},
		{Point{X: 10, Y: 9}, true},
		{Point{X: 10, Y: 10}, true},
		{Point{X: 10, Y: 11}, true},
		{Point{X: 11, Y: 9}, true},
		{Point{X: 11, Y: 10}, true},
		{Point{X: 11, Y: 11}, true},
		{Point{X: math.MaxInt32, Y: 11}, true},
		{Point{X: 9, Y: 99}, false},
		{Point{X: 9, Y: 100}, true},
		{Point{X: 9, Y: 101}, true},
		{Point{X: 9, Y: math.MaxInt32}, true},
		{Point{X: math.MaxInt32, Y: math.MaxInt32}, true},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		// Test with point as head
		s := Snake{Body: []Point{test.Point}}
		require.Equal(t, test.Expected, r.snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Head%+v", test.Point)
		// Test with point as body
		s = Snake{Body: []Point{Point{0, 0}, Point{0, 0}, test.Point}}
		require.Equal(t, test.Expected, r.snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Body%+v", test.Point)
	}
}

func TestSnakeHasBodyCollidedSelf(t *testing.T) {
	tests := []struct {
		Body     []Point
		Expected bool
	}{
		{[]Point{{1, 1}}, false},
		// Self stacks should self collide
		// (we rely on snakes moving before we check self-collision on turn one)
		{[]Point{{2, 2}, {2, 2}}, true},
		{[]Point{{3, 3}, {3, 3}, {3, 3}}, true},
		{[]Point{{5, 5}, {5, 5}, {5, 5}, {5, 5}, {5, 5}}, true},
		// Non-collision cases
		{[]Point{{0, 0}, {1, 0}, {1, 0}}, false},
		{[]Point{{0, 0}, {1, 0}, {2, 0}}, false},
		{[]Point{{0, 0}, {1, 0}, {2, 0}, {2, 0}, {2, 0}}, false},
		{[]Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}}, false},
		{[]Point{{0, 0}, {0, 1}, {0, 2}}, false},
		{[]Point{{0, 0}, {0, 1}, {0, 2}, {0, 2}, {0, 2}}, false},
		{[]Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}}, false},
		// Collision cases
		{[]Point{{0, 0}, {1, 0}, {0, 0}}, true},
		{[]Point{{0, 0}, {0, 0}, {1, 0}}, true},
		{[]Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}, true},
		{[]Point{{4, 4}, {3, 4}, {3, 3}, {4, 4}, {4, 4}}, true},
		{[]Point{{3, 3}, {3, 4}, {3, 3}, {4, 4}, {4, 5}}, true},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		s := Snake{Body: test.Body}
		require.Equal(t, test.Expected, r.snakeHasBodyCollided(&s, &s), "Body%q", s.Body)
	}
}

func TestSnakeHasBodyCollidedOther(t *testing.T) {
	tests := []struct {
		SnakeBody []Point
		OtherBody []Point
		Expected  bool
	}{
		{
			// Just heads
			[]Point{{0, 0}},
			[]Point{{1, 1}},
			false,
		},
		{
			// Head-to-heads are not considered in body collisions
			[]Point{{0, 0}},
			[]Point{{0, 0}},
			false,
		},
		{
			// Stacked bodies
			[]Point{{0, 0}},
			[]Point{{0, 0}, {0, 0}},
			true,
		},
		{
			// Separate stacked bodies
			[]Point{{0, 0}, {0, 0}, {0, 0}},
			[]Point{{1, 1}, {1, 1}, {1, 1}},
			false,
		},
		{
			// Stacked bodies, separated heads
			[]Point{{0, 0}, {1, 0}, {1, 0}},
			[]Point{{2, 0}, {1, 0}, {1, 0}},
			false,
		},
		{
			// Mid-snake collision
			[]Point{{1, 1}},
			[]Point{{0, 1}, {1, 1}, {2, 1}},
			true,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		s := &Snake{Body: test.SnakeBody}
		o := &Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, r.snakeHasBodyCollided(s, o), "Snake%q Other%q", s.Body, o.Body)
	}
}

func TestSnakeHasLostHeadToHead(t *testing.T) {
	tests := []struct {
		SnakeBody        []Point
		OtherBody        []Point
		Expected         bool
		ExpectedOpposite bool
	}{
		{
			// Just heads
			[]Point{{0, 0}},
			[]Point{{1, 1}},
			false, false,
		},
		{
			// Just heads colliding
			[]Point{{0, 0}},
			[]Point{{0, 0}},
			true, true,
		},
		{
			// One snake larger
			[]Point{{0, 0}, {1, 0}, {2, 0}},
			[]Point{{0, 0}},
			false, true,
		},
		{
			// Other snake equal
			[]Point{{0, 0}, {1, 0}, {2, 0}},
			[]Point{{0, 0}, {0, 1}, {0, 2}},
			true, true,
		},
		{
			// Other snake longer
			[]Point{{0, 0}, {1, 0}, {2, 0}},
			[]Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			true, false,
		},
		{
			// Body collision
			[]Point{{0, 1}, {1, 1}, {2, 1}},
			[]Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			false, false,
		},
		{
			// Separate stacked bodies, head collision
			[]Point{{3, 10}, {2, 10}, {2, 10}},
			[]Point{{3, 10}, {4, 10}, {4, 10}},
			true, true,
		},
		{
			// Separate stacked bodies, head collision
			[]Point{{10, 3}, {10, 2}, {10, 1}, {10, 0}},
			[]Point{{10, 3}, {10, 4}, {10, 5}},
			false, true,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		s := Snake{Body: test.SnakeBody}
		o := Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, r.snakeHasLostHeadToHead(&s, &o), "Snake%q Other%q", s.Body, o.Body)
		require.Equal(t, test.ExpectedOpposite, r.snakeHasLostHeadToHead(&o, &s), "Snake%q Other%q", s.Body, o.Body)
	}

}

func TestMaybeEliminateSnakes(t *testing.T) {
	tests := []struct {
		Name                     string
		Snakes                   []Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
		Err                      error
	}{
		{
			"Empty",
			[]Snake{},
			[]string{},
			[]string{},
			nil,
		},
		{
			"Zero Snake",
			[]Snake{
				Snake{},
			},
			[]string{NotEliminated},
			[]string{""},
			ErrorZeroLengthSnake,
		},
		{
			"Single Starvation",
			[]Snake{
				Snake{ID: "1", Body: []Point{{1, 1}}},
			},
			[]string{EliminatedByOutOfHealth},
			[]string{""},
			nil,
		},
		{
			"Not Eliminated",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{1, 1}}},
			},
			[]string{NotEliminated},
			[]string{""},
			nil,
		},
		{
			"Out of Bounds",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{-1, 1}}},
			},
			[]string{EliminatedByOutOfBounds},
			[]string{""},
			nil,
		},
		{
			"Self Collision",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 0}}},
			},
			[]string{EliminatedBySelfCollision},
			[]string{"1"},
			nil,
		},
		{
			"Multiple Separate Deaths",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 0}}},
				Snake{ID: "2", Health: 1, Body: []Point{{-1, 1}}},
			},
			[]string{
				EliminatedBySelfCollision,
				EliminatedByOutOfBounds},
			[]string{"1", ""},
			nil,
		},
		{
			"Other Collision",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{0, 2}, {0, 3}, {0, 4}}},
				Snake{ID: "2", Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			[]string{
				EliminatedByCollision,
				NotEliminated},
			[]string{"2", ""},
			nil,
		},
		{
			"All Eliminated Head 2 Head",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{1, 1}}},
				Snake{ID: "2", Health: 1, Body: []Point{{1, 1}}},
				Snake{ID: "3", Health: 1, Body: []Point{{1, 1}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1"},
			nil,
		},
		{
			"One Snake wins Head 2 Head",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{1, 1}, {0, 1}}},
				Snake{ID: "2", Health: 1, Body: []Point{{1, 1}, {1, 2}, {1, 3}}},
				Snake{ID: "3", Health: 1, Body: []Point{{1, 1}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				NotEliminated,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "", "2"},
			nil,
		},
		{
			"All Snakes Body Eliminated",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{4, 4}, {3, 3}}},
				Snake{ID: "2", Health: 1, Body: []Point{{3, 3}, {2, 2}}},
				Snake{ID: "3", Health: 1, Body: []Point{{2, 2}, {1, 1}}},
				Snake{ID: "4", Health: 1, Body: []Point{{1, 1}, {4, 4}}},
				Snake{ID: "5", Health: 1, Body: []Point{{4, 4}}}, // Body collision takes priority
			},
			[]string{
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
			},
			[]string{"4", "1", "2", "3", "4"},
			nil,
		},
		{
			"All Snakes Eliminated Head 2 Head",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{4, 4}, {4, 5}}},
				Snake{ID: "2", Health: 1, Body: []Point{{4, 4}, {4, 3}}},
				Snake{ID: "3", Health: 1, Body: []Point{{4, 4}, {5, 4}}},
				Snake{ID: "4", Health: 1, Body: []Point{{4, 4}, {3, 4}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1", "1"},
			nil,
		},
		{
			"4 Snakes Head 2 Head",
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{4, 4}, {4, 5}}},
				Snake{ID: "2", Health: 1, Body: []Point{{4, 4}, {4, 3}}},
				Snake{ID: "3", Health: 1, Body: []Point{{4, 4}, {5, 4}, {6, 4}}},
				Snake{ID: "4", Health: 1, Body: []Point{{4, 4}, {3, 4}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				NotEliminated,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"3", "3", "", "3"},
			nil,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b := &BoardState{
				Width:  10,
				Height: 10,
				Snakes: test.Snakes,
			}
			err := r.maybeEliminateSnakes(b)
			require.Equal(t, test.Err, err)
			for i, snake := range b.Snakes {
				require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
				require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy)
			}
		})
	}
}

func TestMaybeEliminateSnakesPriority(t *testing.T) {
	tests := []struct {
		Snakes                   []Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
	}{
		{
			[]Snake{
				{ID: "1", Health: 0, Body: []Point{{-1, 0}, {0, 0}, {1, 0}}},
				{ID: "2", Health: 1, Body: []Point{{-1, 0}, {0, 0}, {1, 0}}},
				{ID: "3", Health: 1, Body: []Point{{1, 0}, {0, 0}, {1, 0}}},
				{ID: "4", Health: 1, Body: []Point{{1, 0}, {1, 1}, {1, 2}}},
				{ID: "5", Health: 1, Body: []Point{{2, 2}, {2, 1}, {2, 0}}},
				{ID: "6", Health: 1, Body: []Point{{2, 2}, {2, 3}, {2, 4}, {2, 5}}},
			},
			[]string{
				EliminatedByOutOfHealth,
				EliminatedByOutOfBounds,
				EliminatedBySelfCollision,
				EliminatedByCollision,
				EliminatedByHeadToHeadCollision,
				NotEliminated,
			},
			[]string{"", "", "3", "3", "6", ""},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{Width: 10, Height: 10, Snakes: test.Snakes}
		err := r.maybeEliminateSnakes(b)
		require.NoError(t, err)
		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause, snake.ID)
			require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy, snake.ID)
		}
	}
}

func TestMaybeFeedSnakes(t *testing.T) {
	tests := []struct {
		Name           string
		Snakes         []Snake
		Food           []Point
		ExpectedSnakes []Snake
		ExpectedFood   []Point
	}{
		{
			Name: "snake not on food",
			Snakes: []Snake{
				{Health: 5, Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			Food: []Point{{3, 3}},
			ExpectedSnakes: []Snake{
				{Health: 5, Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			ExpectedFood: []Point{{3, 3}},
		},
		{
			Name: "snake on food",
			Snakes: []Snake{
				{Health: SnakeMaxHealth - 1, Body: []Point{{2, 1}, {1, 1}, {1, 2}, {2, 2}}},
			},
			Food: []Point{{2, 1}},
			ExpectedSnakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{2, 1}, {1, 1}, {1, 2}, {2, 2}, {2, 2}}},
			},
			ExpectedFood: []Point{},
		},
		{
			Name: "food under body",
			Snakes: []Snake{
				{Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			Food: []Point{{0, 1}},
			ExpectedSnakes: []Snake{
				{Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			ExpectedFood: []Point{{0, 1}},
		},
		{
			Name: "snake on food but already eliminated",
			Snakes: []Snake{
				{Body: []Point{{0, 0}, {0, 1}, {0, 2}}, EliminatedCause: "EliminatedByOutOfBounds"},
			},
			Food: []Point{{0, 0}},
			ExpectedSnakes: []Snake{
				{Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			ExpectedFood: []Point{{0, 0}},
		},
		{
			Name: "multiple snakes on same food",
			Snakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
				{Health: SnakeMaxHealth - 9, Body: []Point{{0, 0}, {1, 0}, {2, 0}}},
			},
			Food: []Point{{0, 0}, {4, 4}},
			ExpectedSnakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{0, 0}, {0, 1}, {0, 2}, {0, 2}}},
				{Health: SnakeMaxHealth, Body: []Point{{0, 0}, {1, 0}, {2, 0}, {2, 0}}},
			},
			ExpectedFood: []Point{{4, 4}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Snakes: test.Snakes,
			Food:   test.Food,
		}
		err := r.maybeFeedSnakes(b)
		require.NoError(t, err, test.Name)
		require.Equal(t, len(test.ExpectedSnakes), len(b.Snakes), test.Name)
		for i := 0; i < len(b.Snakes); i++ {
			require.Equal(t, test.ExpectedSnakes[i].Health, b.Snakes[i].Health, test.Name)
			require.Equal(t, test.ExpectedSnakes[i].Body, b.Snakes[i].Body, test.Name)
		}
		require.Equal(t, test.ExpectedFood, b.Food, test.Name)
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
			[]Point{{0, 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  2,
			},
			[]Point{{0, 0}, {1, 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  1,
				Food:   []Point{{0, 0}, {101, 202}, {-4, -5}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{0, 0}, {1, 0}},
			},
			[]Point{{0, 1}, {1, 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 4,
				Width:  1,
				Snakes: []Snake{
					{Body: []Point{{0, 0}}},
				},
			},
			[]Point{{0, 1}, {0, 2}, {0, 3}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Snakes: []Snake{
					{Body: []Point{{0, 0}, {1, 0}, {1, 1}}},
				},
			},
			[]Point{{0, 1}, {2, 0}, {2, 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []Point{{0, 0}, {1, 0}, {1, 1}, {2, 0}},
				Snakes: []Snake{
					{Body: []Point{{0, 0}, {1, 0}, {1, 1}}},
					{Body: []Point{{0, 1}}},
				},
			},
			[]Point{{2, 1}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		unoccupiedPoints := r.getUnoccupiedPoints(test.Board, true)
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
			[]Point{{0, 0}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
			},
			[]Point{{0, 0}, {1, 1}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  1,
				Food:   []Point{{0, 0}, {101, 202}, {-4, -5}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []Point{{0, 0}, {1, 0}},
			},
			[]Point{{1, 1}},
		},
		{
			&BoardState{
				Height: 4,
				Width:  4,
				Food:   []Point{{0, 0}, {0, 2}, {1, 1}, {1, 3}, {2, 0}, {2, 2}, {3, 1}, {3, 3}},
			},
			[]Point{},
		},
		{
			&BoardState{
				Height: 4,
				Width:  1,
				Snakes: []Snake{
					{Body: []Point{{0, 0}}},
				},
			},
			[]Point{{0, 2}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Snakes: []Snake{
					{Body: []Point{{0, 0}, {1, 0}, {1, 1}}},
				},
			},
			[]Point{{2, 0}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []Point{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
				Snakes: []Snake{
					{Body: []Point{{0, 0}, {1, 0}, {1, 1}}},
					{Body: []Point{{0, 1}}},
				},
			},
			[]Point{{2, 0}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		evenUnoccupiedPoints := r.getEvenUnoccupiedPoints(test.Board)
		require.Equal(t, len(test.Expected), len(evenUnoccupiedPoints))
		for i, e := range test.Expected {
			require.Equal(t, e, evenUnoccupiedPoints[i])
		}
	}
}

func TestMaybeSpawnFoodMinimum(t *testing.T) {
	tests := []struct {
		MinimumFood  int32
		Food         []Point
		ExpectedFood int
	}{
		// Use pre-tested seeds and results
		{0, []Point{}, 0},
		{1, []Point{}, 1},
		{9, []Point{}, 9},
		{7, []Point{{4, 5}, {4, 4}, {4, 1}}, 7},
	}

	for _, test := range tests {
		r := StandardRuleset{MinimumFood: test.MinimumFood}
		b := &BoardState{
			Height: 11,
			Width:  11,
			Snakes: []Snake{
				{Body: []Point{{1, 0}, {1, 1}}},
				{Body: []Point{{0, 1}, {0, 2}, {0, 3}}},
			},
			Food: test.Food,
		}

		err := r.maybeSpawnFood(b)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedFood, len(b.Food))
	}
}

func TestMaybeSpawnFoodZeroChance(t *testing.T) {
	r := StandardRuleset{FoodSpawnChance: 0}
	b := &BoardState{
		Height: 11,
		Width:  11,
		Snakes: []Snake{
			{Body: []Point{{1, 0}, {1, 1}}},
			{Body: []Point{{0, 1}, {0, 2}, {0, 3}}},
		},
		Food: []Point{},
	}
	for i := 0; i < 1000; i++ {
		err := r.maybeSpawnFood(b)
		require.NoError(t, err)
		require.Equal(t, len(b.Food), 0)
	}
}

func TestMaybeSpawnFoodHundredChance(t *testing.T) {
	r := StandardRuleset{FoodSpawnChance: 100}
	b := &BoardState{
		Height: 11,
		Width:  11,
		Snakes: []Snake{
			{Body: []Point{{1, 0}, {1, 1}}},
			{Body: []Point{{0, 1}, {0, 2}, {0, 3}}},
		},
		Food: []Point{},
	}
	for i := 1; i <= 22; i++ {
		err := r.maybeSpawnFood(b)
		require.NoError(t, err)
		require.Equal(t, i, len(b.Food))
	}
}

func TestMaybeSpawnFoodHalfChance(t *testing.T) {
	tests := []struct {
		Seed         int64
		Food         []Point
		ExpectedFood int32
	}{
		// Use pre-tested seeds and results
		{123, []Point{}, 1},
		{12345, []Point{}, 0},
		{456, []Point{{4, 4}}, 1},
		{789, []Point{{4, 4}}, 2},
		{511, []Point{{4, 4}}, 1},
		{165, []Point{{4, 4}}, 2},
	}

	r := StandardRuleset{FoodSpawnChance: 50}
	for _, test := range tests {
		b := &BoardState{
			Height: 4,
			Width:  5,
			Snakes: []Snake{
				{Body: []Point{{1, 0}, {1, 1}}},
				{Body: []Point{{0, 1}, {0, 2}, {0, 3}}},
			},
			Food: test.Food,
		}

		rand.Seed(test.Seed)
		err := r.maybeSpawnFood(b)
		require.NoError(t, err)
		require.Equal(t, test.ExpectedFood, int32(len(b.Food)), "Seed %d", test.Seed)
	}
}

func TestSpawnFood(t *testing.T) {
	b := &BoardState{
		Height: 1,
		Width:  3,
		Snakes: []Snake{
			{Body: []Point{{1, 0}}},
		},
	}
	// Food should never spawn, no room
	r := StandardRuleset{}
	err := r.spawnFood(b, 99)
	require.NoError(t, err)
	require.Equal(t, len(b.Food), 0)
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []Snake
		Expected bool
	}{
		{[]Snake{}, true},
		{[]Snake{{}}, true},
		{[]Snake{{}, {}}, false},
		{[]Snake{{}, {}, {}, {}, {}}, false},
		{
			[]Snake{
				{EliminatedCause: EliminatedByCollision},
				{EliminatedCause: NotEliminated},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: NotEliminated},
				{EliminatedCause: EliminatedByCollision},
				{EliminatedCause: NotEliminated},
				{EliminatedCause: NotEliminated},
			},
			false,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: NotEliminated},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: NotEliminated},
				{EliminatedCause: NotEliminated},
			},
			false,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Height: 11,
			Width:  11,
			Snakes: test.Snakes,
			Food:   []Point{},
		}

		actual, err := r.IsGameOver(b)
		require.NoError(t, err)
		require.Equal(t, test.Expected, actual)
	}
}
