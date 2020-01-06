package rules

import (
	"errors"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanity(t *testing.T) {
	r := StandardRuleset{}

	state, err := r.CreateInitialBoardState(0, 0, []string{})
	require.NoError(t, err)
	require.NotNil(t, state)
	require.Equal(t, int32(0), state.Width)
	require.Equal(t, int32(0), state.Height)
	require.Len(t, state.Food, 0)
	require.Len(t, state.Snakes, 0)

	next, err := r.ResolveMoves(
		&BoardState{},
		[]SnakeMove{},
	)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.Equal(t, int32(0), state.Width)
	require.Equal(t, int32(0), state.Height)
	require.Len(t, state.Snakes, 0)
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
		{1, 2, []string{"one"}, 1, nil},
		{9, 8, []string{"one"}, 1, nil},
		{2, 2, []string{"one", "two"}, 2, nil},
		{2, 2, []string{"one", "two"}, 2, nil},
		{1, 1, []string{"one", "two"}, 2, errors.New("not enough space to place snake")},
	}

	r := StandardRuleset{}
	for _, test := range tests {
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
		require.Len(t, state.Food, test.ExpectedNumFood)
	}
}

func TestPlaceSnakes(t *testing.T) {
	// Because placement is random, we only test to ensure
	// that snake bodies are populated correctly
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
			errors.New("not enough space to place snake"),
		},
		{
			&BoardState{
				Width:  10,
				Height: 5,
				Snakes: make([]Snake, 49),
			},
			nil,
		},
		{
			&BoardState{
				Width:  5,
				Height: 10,
				Snakes: make([]Snake, 50),
			},
			nil,
		},
		{
			&BoardState{
				Width:  25,
				Height: 2,
				Snakes: make([]Snake, 51),
			},
			errors.New("not enough space to place snake"),
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
			errors.New("too many snakes for fixed start positions"),
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
			errors.New("too many snakes for fixed start positions"),
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
			errors.New("too many snakes for fixed start positions"),
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		require.Equal(t, test.BoardState.Width*test.BoardState.Height, int32(len(r.getUnoccupiedPoints(test.BoardState))))
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

func TestResolveMoves(t *testing.T) {
	// TODO
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
			MoveUp, []Point{{10, 109}, {10, 110}},
			MoveDown, []Point{{23, 221}, {23, 220}, {22, 220}, {21, 220}},
			MoveUp, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{11, 109}, {10, 109}},
			MoveLeft, []Point{{22, 221}, {23, 221}, {23, 220}, {22, 220}},
			MoveUp, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{12, 109}, {11, 109}},
			MoveLeft, []Point{{21, 221}, {22, 221}, {23, 221}, {23, 220}},
			MoveUp, []Point{{0, 0}},
		},
		{
			MoveRight, []Point{{13, 109}, {12, 109}},
			MoveLeft, []Point{{20, 221}, {21, 221}, {22, 221}, {23, 221}},
			MoveUp, []Point{{0, 0}},
		},
		{
			MoveUp, []Point{{13, 108}, {13, 109}},
			MoveDown, []Point{{20, 222}, {20, 221}, {21, 221}, {22, 221}},
			MoveUp, []Point{{0, 0}},
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
	require.Equal(t, err, errors.New("snake not found"))
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
	require.Equal(t, err, errors.New("not enough snake moves"))
}

func TestMoveSnakesTooManyMoves(t *testing.T) {
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
			Move: MoveUp,
		},
		{
			ID:   "two",
			Move: MoveUp,
		},
	}

	r := StandardRuleset{}
	err := r.moveSnakes(b, moves)
	require.Equal(t, err, errors.New("too many snake moves"))
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
			Expected: []Point{{0, -1}},
		},
		{
			Body:     []Point{{5, 5}, {5, 5}},
			Move:     "",
			Expected: []Point{{5, 4}, {5, 5}},
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

func TestSnakeHasStarved(t *testing.T) {
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
		require.Equal(t, test.Expected, r.snakeHasStarved(s), "Health: %+v", test.Health)
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

func TestEliminateSnakes(t *testing.T) {
	tests := []struct {
		Snakes                   []Snake
		ExpectedEliminatedCauses []string
		Err                      error
	}{
		{
			[]Snake{},
			[]string{},
			nil,
		},
		{
			[]Snake{
				Snake{},
			},
			[]string{NotEliminated},
			errors.New("snake is length zero"),
		},
		{
			[]Snake{
				Snake{Body: []Point{{1, 1}}},
			},
			[]string{EliminatedByStarvation},
			nil,
		},
		{
			[]Snake{
				Snake{Health: 1, Body: []Point{{1, 1}}},
			},
			[]string{NotEliminated},
			nil,
		},
		{
			[]Snake{
				Snake{Health: 1, Body: []Point{{-1, 1}}},
			},
			[]string{EliminatedByOutOfBounds},
			nil,
		},
		{
			[]Snake{
				Snake{Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 0}}},
			},
			[]string{EliminatedBySelfCollision},
			nil,
		},
		{
			[]Snake{
				Snake{Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 0}}},
				Snake{Health: 1, Body: []Point{{-1, 1}}},
			},
			[]string{
				EliminatedBySelfCollision,
				EliminatedByOutOfBounds},
			nil,
		},
		{
			[]Snake{
				Snake{ID: "1", Health: 1, Body: []Point{{0, 2}, {0, 3}, {0, 4}}},
				Snake{ID: "2", Health: 1, Body: []Point{{0, 0}, {0, 1}, {0, 2}}},
			},
			[]string{
				EliminatedByCollision,
				NotEliminated},
			nil,
		},
		{
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
			nil,
		},
		{
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
			nil,
		},
		{
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
			nil,
		},
		{
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
			nil,
		},
		{
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
			nil,
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Width:  10,
			Height: 10,
			Snakes: test.Snakes,
		}
		err := r.eliminateSnakes(b)
		require.Equal(t, test.Err, err)
		for i := 0; i < len(b.Snakes); i++ {
			require.Equal(t, test.ExpectedEliminatedCauses[i], b.Snakes[i].EliminatedCause)
		}
	}
}

func TestFeedSnakes(t *testing.T) {
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
		err := r.feedSnakes(b)
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
		unoccupiedPoints := r.getUnoccupiedPoints(test.Board)
		require.Equal(t, len(test.Expected), len(unoccupiedPoints))
		for i, e := range test.Expected {
			require.Equal(t, e, unoccupiedPoints[i])
		}
	}
}

func TestMaybeSpawnFood(t *testing.T) {
	tests := []struct {
		Seed         int64
		ExpectedFood []Point
	}{
		// Use pre-tested seeds and results
		{123, []Point{}},
		{456, []Point{}},
		{789, []Point{}},
		{1024, []Point{{2, 1}}},
		{511, []Point{{2, 0}}},
		{165, []Point{{3, 1}}},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Height: 4,
			Width:  5,
			Snakes: []Snake{
				{Body: []Point{{1, 0}, {1, 1}}},
				{Body: []Point{{0, 1}, {0, 2}, {0, 3}}},
			},
		}

		rand.Seed(test.Seed)
		err := r.maybeSpawnFood(b, 1)
		require.NoError(t, err)
		require.Equal(t, len(test.ExpectedFood), len(b.Food), "Seed %d", test.Seed)
		for i, e := range test.ExpectedFood {
			require.Equal(t, e, b.Food[i], "Seed %d", test.Seed)
		}
	}
}

// func TestCheckForSnakesEating(t *testing.T) {
// 	snake := &pb.Snake{
// 		Body: []*pb.Point{
// 			{X: 2, Y: 1},
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 2, Y: 2},
// 		},
// 	}
// 	checkForSnakesEating(&pb.GameFrame{
// 		Food: []*pb.Point{
// 			{X: 2, Y: 1},
// 		},
// 		Snakes: []*pb.Snake{snake},
// 	})
// 	require.Len(t, snake.Body, 4)
// 	require.Equal(t, snake.Body[2], snake.Body[3])
// }

// func TestCheckForSnakesNotEating(t *testing.T) {
// 	snake := &pb.Snake{
// 		Body: []*pb.Point{
// 			{X: 2, Y: 1},
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 2, Y: 2},
// 		},
// 	}
// 	checkForSnakesEating(&pb.GameFrame{
// 		Food:   []*pb.Point{},
// 		Snakes: []*pb.Snake{snake},
// 	})
// 	require.Len(t, snake.Body, 3)
// 	require.NotEqual(t, snake.Body[2], snake.Body[1])
// }
// func TestGameFrameSnakeEats(t *testing.T) {
// 	snake := &pb.Snake{
// 		Health: 67,
// 		Body: []*pb.Point{
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 1, Y: 3},
// 		},
// 	}

// 	lastFrame.Snakes = []*pb.Snake{snake}

// 	gt, err := GameTick(commonGame, lastFrame)
// 	require.NoError(t, err)
// 	require.Len(t, gt.Snakes, 1)
// 	snake = gt.Snakes[0]
// 	require.Equal(t, int32(100), snake.Health)
// 	require.Len(t, snake.Body, 4)
// }

// func TestUpdateFood(t *testing.T) {
// 	updated, err := updateFood(&pb.Game{Width: 20, Height: 20}, &pb.GameFrame{
// 		Food: []*pb.Point{
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 		},
// 		Snakes: []*pb.Snake{
// 			{
// 				Body: []*pb.Point{
// 					{X: 1, Y: 2},
// 					{X: 2, Y: 2},
// 					{X: 3, Y: 2},
// 				},
// 			},
// 		},
// 	}, []*pb.Point{
// 		{X: 1, Y: 2},
// 	})
// 	require.NoError(t, err)
// 	require.Len(t, updated, 2)
// 	require.True(t, updated[0].Equal(&pb.Point{X: 1, Y: 1}))
// 	require.False(t, updated[1].Equal(&pb.Point{X: 1, Y: 2}))
// 	require.False(t, updated[1].Equal(&pb.Point{X: 2, Y: 2}))
// 	require.False(t, updated[1].Equal(&pb.Point{X: 3, Y: 2}))
// 	require.False(t, updated[1].Equal(&pb.Point{X: 1, Y: 1}))
// }

// func TestUpdateFoodWithFullBoard(t *testing.T) {
// 	updated, err := updateFood(&pb.Game{Width: 2, Height: 2}, &pb.GameFrame{
// 		Food: []*pb.Point{
// 			{X: 0, Y: 0},
// 		},
// 		Snakes: []*pb.Snake{
// 			{
// 				Body: []*pb.Point{
// 					{X: 0, Y: 0},
// 					{X: 0, Y: 1},
// 					{X: 1, Y: 1},
// 					{X: 1, Y: 0},
// 				},
// 			},
// 		},
// 	}, []*pb.Point{
// 		{X: 0, Y: 0},
// 	})
// 	require.NoError(t, err)
// 	require.Len(t, updated, 0)
// }

// func TestGetUnoccupiedPointEven(t *testing.T) {

// 	unoccupiedPoint := getUnoccupiedPointEven(2, 2,
// 		[]*pb.Point{},
// 		[]*pb.Snake{})
// 	require.True(t, (unoccupiedPoint.X+unoccupiedPoint.Y)%2 == 0, "Point coordinates should sum to an even number %o ", unoccupiedPoint)
// }

// func TestGetUnoccupiedPointOdd(t *testing.T) {
// 	unoccupiedPoint := getUnoccupiedPointOdd(2, 2,
// 		[]*pb.Point{{X: 0, Y: 1}},
// 		[]*pb.Snake{})
// 	require.True(t, (unoccupiedPoint.X+unoccupiedPoint.Y)%2 == 1, "Point coordinates should sum to an odd number %o ", unoccupiedPoint)
// }

// func TestGetUnoccupiedPointWithFullBoard(t *testing.T) {
// 	unoccupiedPoint := getUnoccupiedPoint(2, 2,
// 		[]*pb.Point{{X: 0, Y: 0}},
// 		[]*pb.Snake{
// 			{
// 				Body: []*pb.Point{
// 					{X: 0, Y: 1},
// 					{X: 1, Y: 1},
// 					{X: 1, Y: 0},
// 				},
// 			},
// 		})
// 	require.True(t, unoccupiedPoint.Equal(nil))
// }

// func TestGetUnoccupiedPointsWithEmptySpots(t *testing.T) {
// 	unoccupiedPoints := getUnoccupiedPoints(2, 2,
// 		[]*pb.Point{{X: 0, Y: 0}},
// 		[]*pb.Snake{
// 			{
// 				Body: []*pb.Point{
// 					{X: 0, Y: 1},
// 				},
// 			},
// 		})

// 	require.Len(t, unoccupiedPoints, 2)
// 	require.True(t, unoccupiedPoints[0].Equal(&pb.Point{X: 1, Y: 0}))
// 	require.True(t, unoccupiedPoints[1].Equal(&pb.Point{X: 1, Y: 1}))
// }

// func TestGetUniqOccupiedPoints(t *testing.T) {
// 	unoccupiedPoints := getUniqOccupiedPoints(
// 		[]*pb.Point{
// 			{X: 0, Y: 0},
// 		},
// 		[]*pb.Snake{
// 			{
// 				Body: []*pb.Point{
// 					{X: 0, Y: 1},
// 					{X: 1, Y: 1},
// 					{X: 1, Y: 1},
// 					{X: 1, Y: 0},
// 				},
// 			},
// 		})

// 	require.Len(t, unoccupiedPoints, 4)
// }

// func TestGameTickUpdatesTurnCounter(t *testing.T) {
// 	gt, err := GameTick(commonGame, &pb.GameFrame{Turn: 5})
// 	require.NoError(t, err)
// 	require.Equal(t, int32(6), gt.Turn)
// }

// func TestGameTickUpdatesSnake(t *testing.T) {
// 	snake := &pb.Snake{
// 		Health: 67,
// 		Body: []*pb.Point{
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 1, Y: 3},
// 		},
// 	}
// 	game := &pb.Game{
// 		Width:  20,
// 		Height: 20,
// 	}
// 	gt, err := GameTick(game, &pb.GameFrame{
// 		Turn: 5,
// 		Snakes: []*pb.Snake{
// 			snake,
// 		},
// 	})
// 	require.NoError(t, err)
// 	require.Len(t, gt.Snakes, 1)
// 	snake = gt.Snakes[0]
// 	require.Equal(t, int32(66), snake.Health)
// 	require.Len(t, snake.Body, 3)
// 	require.Equal(t, &pb.Point{X: 1, Y: 0}, snake.Body[0])
// 	require.Equal(t, &pb.Point{X: 1, Y: 1}, snake.Body[1])
// 	require.Equal(t, &pb.Point{X: 1, Y: 2}, snake.Body[2])
// }

// var commonGame = &pb.Game{
// 	Width:  20,
// 	Height: 20,
// }
// var lastFrame = &pb.GameFrame{
// 	Turn:   5,
// 	Snakes: []*pb.Snake{},
// 	Food: []*pb.Point{
// 		{X: 1, Y: 0},
// 	},
// }

// func TestGameTickDeadSnakeDoNotUpdate(t *testing.T) {
// 	snake := &pb.Snake{
// 		Health: 87,
// 		Body: []*pb.Point{
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 1, Y: 3},
// 		},
// 		Death: &pb.Death{
// 			Turn:  4,
// 			Cause: DeathCauseSnakeCollision,
// 		},
// 	}

// 	lastFrame.Snakes = []*pb.Snake{snake}

// 	gt, err := GameTick(commonGame, lastFrame)
// 	require.NoError(t, err)
// 	require.Len(t, gt.Snakes, 1)
// 	snake = gt.Snakes[0]
// 	require.Equal(t, int32(87), snake.Health)
// 	require.Len(t, snake.Body, 3)
// 	require.Equal(t, &pb.Point{X: 1, Y: 1}, snake.Body[0])
// 	require.Equal(t, &pb.Point{X: 1, Y: 2}, snake.Body[1])
// 	require.Equal(t, &pb.Point{X: 1, Y: 3}, snake.Body[2])
// }

// func TestGameTickUpdatesDeath(t *testing.T) {
// 	snake := &pb.Snake{
// 		Health: 0,
// 		Body: []*pb.Point{
// 			{X: 3, Y: 1},
// 			{X: 3, Y: 2},
// 			{X: 3, Y: 3},
// 		},
// 	}

// 	lastFrame.Snakes = []*pb.Snake{snake}

// 	gt, err := GameTick(commonGame, lastFrame)
// 	require.NoError(t, err)
// 	require.NotNil(t, gt.Snakes[0].Death)
// }

// func TestUpdateSnakes(t *testing.T) {
// 	snake := &pb.Snake{
// 		Body: []*pb.Point{
// 			{X: 1, Y: 1},
// 		},
// 	}
// 	moves := []*SnakeUpdate{
// 		{
// 			Snake: snake,
// 			Err:   errors.New("some error"),
// 		},
// 	}
// 	updateSnakes(&pb.Game{}, &pb.GameFrame{
// 		Snakes: []*pb.Snake{snake},
// 	}, moves)
// 	require.Equal(t, &pb.Point{X: 1, Y: 0}, snake.Head(), "snake did not move up")

// 	moves = []*SnakeUpdate{
// 		{
// 			Snake: snake,
// 			Move:  "left",
// 		},
// 	}
// 	updateSnakes(&pb.Game{}, &pb.GameFrame{
// 		Snakes: []*pb.Snake{snake},
// 	}, moves)
// 	require.Equal(t, &pb.Point{X: 0, Y: 0}, snake.Head(), "snake did not move left")
// }

// func TestCanFollowTail(t *testing.T) {
// 	url := setupSnakeServer(t, MoveResponse{
// 		Move: "down",
// 	}, StartResponse{})
// 	snake := &pb.Snake{
// 		Body: []*pb.Point{
// 			{X: 2, Y: 1},
// 			{X: 1, Y: 1},
// 			{X: 1, Y: 2},
// 			{X: 2, Y: 2},
// 		},
// 		URL:    url,
// 		Health: 100,
// 	}
// 	next, err := GameTick(&pb.Game{
// 		Width:  20,
// 		Height: 20,
// 	}, &pb.GameFrame{
// 		Snakes: []*pb.Snake{snake},
// 	})
// 	require.NoError(t, err)
// 	require.NotNil(t, next)
// 	require.Nil(t, next.Snakes[0].Death)
// }

// func TestNextFoodSpawn(t *testing.T) {
// 	rand.Seed(1) // random order is 65, 85, 29
// 	snakes := []*pb.Snake{
// 		{URL: setupSnakeServer(t, MoveResponse{}, StartResponse{})},
// 		{URL: setupSnakeServer(t, MoveResponse{}, StartResponse{})},
// 		{URL: setupSnakeServer(t, MoveResponse{}, StartResponse{})},
// 		{URL: setupSnakeServer(t, MoveResponse{}, StartResponse{})},
// 	}
// 	next, err := GameTick(&pb.Game{
// 		Width:                   20,
// 		Height:                  20,
// 		TurnsSinceLastFoodSpawn: 5,
// 		MaxTurnsToNextFoodSpawn: 5,
// 	}, &pb.GameFrame{
// 		Snakes: snakes,
// 	})
// 	require.NoError(t, err)
// 	require.Len(t, next.Food, 2)
// }
