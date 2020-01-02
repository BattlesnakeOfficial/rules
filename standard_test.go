package rulesets

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanity(t *testing.T) {
	r := StandardRuleset{}
	next, err := r.ResolveMoves(
		&BoardState{},
		[]*SnakeMove{},
	)

	require.NoError(t, err)
	require.NotNil(t, next)
}

func TestMoveSnakes(t *testing.T) {
	b := &BoardState{
		Snakes: []*Snake{
			{
				ID:     "one",
				Body:   []*Point{{10, 110}, {11, 110}},
				Health: 111111,
			},
			{
				ID:     "two",
				Body:   []*Point{{23, 220}, {22, 220}, {21, 220}, {20, 220}},
				Health: 222222,
			},
		},
	}

	tests := []struct {
		MoveOne     string
		ExpectedOne []*Point
		MoveTwo     string
		ExpectedTwo []*Point
	}{
		{
			MOVE_UP,
			[]*Point{{10, 109}, {10, 110}},
			MOVE_DOWN,
			[]*Point{{23, 221}, {23, 220}, {22, 220}, {21, 220}},
		},
		{
			MOVE_RIGHT,
			[]*Point{{11, 109}, {10, 109}},
			MOVE_LEFT,
			[]*Point{{22, 221}, {23, 221}, {23, 220}, {22, 220}},
		},
		{
			MOVE_RIGHT,
			[]*Point{{12, 109}, {11, 109}},
			MOVE_LEFT,
			[]*Point{{21, 221}, {22, 221}, {23, 221}, {23, 220}},
		},
		{
			MOVE_RIGHT,
			[]*Point{{13, 109}, {12, 109}},
			MOVE_LEFT,
			[]*Point{{20, 221}, {21, 221}, {22, 221}, {23, 221}},
		},
		{
			MOVE_UP,
			[]*Point{{13, 108}, {13, 109}},
			MOVE_DOWN,
			[]*Point{{20, 222}, {20, 221}, {21, 221}, {22, 221}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		moves := []*SnakeMove{
			{Snake: b.Snakes[0], Move: test.MoveOne},
			{Snake: b.Snakes[1], Move: test.MoveTwo},
		}
		err := r.moveSnakes(b, moves)

		require.NoError(t, err)
		require.Len(t, b.Snakes, 2)
		require.Equal(t, int32(111111), b.Snakes[0].Health)
		require.Equal(t, int32(222222), b.Snakes[1].Health)
		require.Len(t, b.Snakes[0].Body, 2)
		require.Len(t, b.Snakes[1].Body, 4)

		require.Equal(t, len(b.Snakes[0].Body), len(test.ExpectedOne))
		for i, e := range test.ExpectedOne {
			require.Equal(t, *e, *b.Snakes[0].Body[i])
		}
		require.Equal(t, len(b.Snakes[1].Body), len(test.ExpectedTwo))
		for i, e := range test.ExpectedTwo {
			require.Equal(t, *e, *b.Snakes[1].Body[i])
		}
	}
}

func TestMoveSnakesDefault(t *testing.T) {
	tests := []struct {
		Body     []*Point
		Move     string
		Expected []*Point
	}{
		{
			Body:     []*Point{{0, 0}},
			Move:     "asdf",
			Expected: []*Point{{0, -1}},
		},
		{
			Body:     []*Point{{5, 5}, {5, 5}},
			Move:     "",
			Expected: []*Point{{5, 4}, {5, 5}},
		},
		{
			Body:     []*Point{{5, 5}, {5, 4}},
			Expected: []*Point{{5, 6}, {5, 5}},
		},
		{
			Body:     []*Point{{5, 4}, {5, 5}},
			Expected: []*Point{{5, 3}, {5, 4}},
		},
		{
			Body:     []*Point{{5, 4}, {5, 5}},
			Expected: []*Point{{5, 3}, {5, 4}},
		},
		{
			Body:     []*Point{{4, 5}, {5, 5}},
			Expected: []*Point{{3, 5}, {4, 5}},
		},
		{
			Body:     []*Point{{5, 5}, {4, 5}},
			Expected: []*Point{{6, 5}, {5, 5}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Snakes: []*Snake{
				{Body: test.Body},
			},
		}
		moves := []*SnakeMove{{Snake: b.Snakes[0], Move: test.Move}}

		err := r.moveSnakes(b, moves)
		require.NoError(t, err)
		require.Len(t, b.Snakes, 1)
		require.Equal(t, len(test.Body), len(b.Snakes[0].Body))
		require.Equal(t, len(test.Expected), len(b.Snakes[0].Body))
		for i, e := range test.Expected {
			require.Equal(t, *e, *b.Snakes[0].Body[i])
		}
	}
}

func TestReduceSnakeHealth(t *testing.T) {
	var err error
	r := StandardRuleset{}
	b := &BoardState{
		Snakes: []*Snake{
			&Snake{
				Body:   []*Point{{0, 0}, {0, 1}},
				Health: 99,
			},
			&Snake{
				Body:   []*Point{{5, 8}, {6, 8}, {7, 8}},
				Health: 2,
			},
		},
	}

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(98))
	require.Equal(t, b.Snakes[1].Health, int32(1))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(97))
	require.Equal(t, b.Snakes[1].Health, int32(0))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(96))
	require.Equal(t, b.Snakes[1].Health, int32(-1))

	err = r.reduceSnakeHealth(b)
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, int32(95))
	require.Equal(t, b.Snakes[1].Health, int32(-2))
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
	var boardWidth int32 = 10
	var boardHeight int32 = 100

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

	var s *Snake
	r := StandardRuleset{}
	for _, test := range tests {
		// Test with point as head
		s = &Snake{Body: []*Point{&test.Point}}
		require.Equal(t, test.Expected, r.snakeIsOutOfBounds(s, boardWidth, boardHeight), "Head%+v", test.Point)
		// Test with point as body
		s = &Snake{Body: []*Point{&Point{0, 0}, &Point{0, 0}, &test.Point}}
		require.Equal(t, test.Expected, r.snakeIsOutOfBounds(s, boardWidth, boardHeight), "Body%+v", test.Point)
	}
}

func TestSnakeHasBodyCollidedSelf(t *testing.T) {
	tests := []struct {
		Body     []*Point
		Expected bool
	}{
		{[]*Point{{1, 1}}, false},
		// Self stacks should self collide
		// (we rely on snakes moving before we check self-collision on turn one)
		{[]*Point{{2, 2}, {2, 2}}, true},
		{[]*Point{{3, 3}, {3, 3}, {3, 3}}, true},
		{[]*Point{{5, 5}, {5, 5}, {5, 5}, {5, 5}, {5, 5}}, true},
		// Non-collision cases
		{[]*Point{{0, 0}, {1, 0}, {1, 0}}, false},
		{[]*Point{{0, 0}, {1, 0}, {2, 0}}, false},
		{[]*Point{{0, 0}, {1, 0}, {2, 0}, {2, 0}, {2, 0}}, false},
		{[]*Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}}, false},
		{[]*Point{{0, 0}, {0, 1}, {0, 2}}, false},
		{[]*Point{{0, 0}, {0, 1}, {0, 2}, {0, 2}, {0, 2}}, false},
		{[]*Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}}, false},
		// Collision cases
		{[]*Point{{0, 0}, {1, 0}, {0, 0}}, true},
		{[]*Point{{0, 0}, {0, 0}, {1, 0}}, true},
		{[]*Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}, true},
		{[]*Point{{4, 4}, {3, 4}, {3, 3}, {4, 4}, {4, 4}}, true},
		{[]*Point{{3, 3}, {3, 4}, {3, 3}, {4, 4}, {4, 5}}, true},
	}

	var s *Snake
	r := StandardRuleset{}
	for _, test := range tests {
		s = &Snake{Body: test.Body}
		require.Equal(t, test.Expected, r.snakeHasBodyCollided(s, s), "Body%q", s.Body)
	}
}

func TestSnakeHasBodyCollidedOther(t *testing.T) {
	tests := []struct {
		SnakeBody []*Point
		OtherBody []*Point
		Expected  bool
	}{
		{
			// Just heads
			[]*Point{{0, 0}},
			[]*Point{{1, 1}},
			false,
		},
		{
			// Head-to-heads are not considered in body collisions
			[]*Point{{0, 0}},
			[]*Point{{0, 0}},
			false,
		},
		{
			// Stacked bodies
			[]*Point{{0, 0}},
			[]*Point{{0, 0}, {0, 0}},
			true,
		},
		{
			// Separate stacked bodies
			[]*Point{{0, 0}, {0, 0}, {0, 0}},
			[]*Point{{1, 1}, {1, 1}, {1, 1}},
			false,
		},
		{
			// Stacked bodies, separated heads
			[]*Point{{0, 0}, {1, 0}, {1, 0}},
			[]*Point{{2, 0}, {1, 0}, {1, 0}},
			false,
		},
		{
			// Mid-snake collision
			[]*Point{{1, 1}},
			[]*Point{{0, 1}, {1, 1}, {2, 1}},
			true,
		},
	}

	var s *Snake
	var o *Snake
	r := StandardRuleset{}
	for _, test := range tests {
		s = &Snake{Body: test.SnakeBody}
		o = &Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, r.snakeHasBodyCollided(s, o), "Snake%q Other%q", s.Body, o.Body)
	}
}

func TestSnakeHasLostHeadToHead(t *testing.T) {
	tests := []struct {
		SnakeBody        []*Point
		OtherBody        []*Point
		Expected         bool
		ExpectedOpposite bool
	}{
		{
			// Just heads
			[]*Point{{0, 0}},
			[]*Point{{1, 1}},
			false, false,
		},
		{
			// Just heads colliding
			[]*Point{{0, 0}},
			[]*Point{{0, 0}},
			true, true,
		},
		{
			// One snake larger
			[]*Point{{0, 0}, {1, 0}, {2, 0}},
			[]*Point{{0, 0}},
			false, true,
		},
		{
			// Other snake equal
			[]*Point{{0, 0}, {1, 0}, {2, 0}},
			[]*Point{{0, 0}, {0, 1}, {0, 2}},
			true, true,
		},
		{
			// Other snake longer
			[]*Point{{0, 0}, {1, 0}, {2, 0}},
			[]*Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			true, false,
		},
		{
			// Body collision
			[]*Point{{0, 1}, {1, 1}, {2, 1}},
			[]*Point{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			false, false,
		},
		{
			// Separate stacked bodies, head collision
			[]*Point{{3, 10}, {2, 10}, {2, 10}},
			[]*Point{{3, 10}, {4, 10}, {4, 10}},
			true, true,
		},
		{
			// Separate stacked bodies, head collision
			[]*Point{{10, 3}, {10, 2}, {10, 1}, {10, 0}},
			[]*Point{{10, 3}, {10, 4}, {10, 5}},
			false, true,
		},
	}

	var s *Snake
	var o *Snake
	r := StandardRuleset{}
	for _, test := range tests {
		s = &Snake{Body: test.SnakeBody}
		o = &Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, r.snakeHasLostHeadToHead(s, o), "Snake%q Other%q", s.Body, o.Body)
		require.Equal(t, test.ExpectedOpposite, r.snakeHasLostHeadToHead(o, s), "Snake%q Other%q", s.Body, o.Body)
	}

}

func TestFeedSnakes(t *testing.T) {
	r := StandardRuleset{}
	b := &BoardState{
		Snakes: []*Snake{
			{Body: []*Point{
				{2, 1}, {1, 1}, {1, 2}, {2, 2},
			}},
		},
		Food: []*Point{
			{2, 1},
		},
	}

	err := r.feedSnakes(b)
	require.NoError(t, err)
	require.Equal(t, 0, len(b.Food))

}

func TestGetUnoccupiedPoints(t *testing.T) {
	tests := []struct {
		Board    *BoardState
		Expected []*Point
	}{
		{
			&BoardState{
				Height: 1,
				Width:  1,
			},
			[]*Point{{0, 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  2,
			},
			[]*Point{{0, 0}, {1, 0}},
		},
		{
			&BoardState{
				Height: 1,
				Width:  1,
				Food:   []*Point{{0, 0}, {101, 202}, {-4, -5}},
			},
			[]*Point{},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []*Point{{0, 0}, {1, 0}},
			},
			[]*Point{{0, 1}, {1, 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  2,
				Food:   []*Point{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			},
			[]*Point{},
		},
		{
			&BoardState{
				Height: 4,
				Width:  1,
				Snakes: []*Snake{
					{Body: []*Point{{0, 0}}},
				},
			},
			[]*Point{{0, 1}, {0, 2}, {0, 3}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Snakes: []*Snake{
					{Body: []*Point{{0, 0}, {1, 0}, {1, 1}}},
				},
			},
			[]*Point{{0, 1}, {2, 0}, {2, 1}},
		},
		{
			&BoardState{
				Height: 2,
				Width:  3,
				Food:   []*Point{{0, 0}, {1, 0}, {1, 1}, {2, 0}},
				Snakes: []*Snake{
					{Body: []*Point{{0, 0}, {1, 0}, {1, 1}}},
					{Body: []*Point{{0, 1}}},
				},
			},
			[]*Point{{2, 1}},
		},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		unoccupiedPoints := r.getUnoccupiedPoints(test.Board)
		require.Equal(t, len(test.Expected), len(unoccupiedPoints))
		for i, e := range test.Expected {
			require.Equal(t, *e, *unoccupiedPoints[i])
		}
	}
}

func TestMaybeSpawnFood(t *testing.T) {
	tests := []struct {
		Seed         int64
		ExpectedFood []*Point
	}{
		// Use pre-tested seeds and results
		{123, []*Point{}},
		{456, []*Point{}},
		{789, []*Point{}},
		{1024, []*Point{{2, 1}}},
		{511, []*Point{{2, 0}}},
		{165, []*Point{{3, 1}}},
	}

	r := StandardRuleset{}
	for _, test := range tests {
		b := &BoardState{
			Height: 4,
			Width:  5,
			Snakes: []*Snake{
				{Body: []*Point{{1, 0}, {1, 1}}},
				{Body: []*Point{{0, 1}, {0, 2}, {0, 3}}},
			},
		}

		rand.Seed(test.Seed)
		err := r.maybeSpawnFood(b, 1)
		require.NoError(t, err)
		require.Equal(t, len(test.ExpectedFood), len(b.Food), "Seed %d", test.Seed)
		for i, e := range test.ExpectedFood {
			require.Equal(t, *e, *b.Food[i], "Seed %d", test.Seed)
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
