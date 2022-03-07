package rules

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSquadRulesetInterface(t *testing.T) {
	var _ Ruleset = (*SquadRuleset)(nil)
}

func TestSquadName(t *testing.T) {
	r := SquadRuleset{}
	require.Equal(t, "squad", r.Name())
}

func TestSquadCreateNextBoardStateSanity(t *testing.T) {
	boardState := &BoardState{}
	r := SquadRuleset{}
	_, err := r.CreateNextBoardState(boardState, []SnakeMove{})
	require.NoError(t, err)
}

func TestSquadResurrectSquadBodyCollisionsSanity(t *testing.T) {
	boardState := &BoardState{}
	r := SquadRuleset{}
	err := r.resurrectSquadBodyCollisions(boardState)
	require.NoError(t, err)
}

func TestSquadSharedAttributesSanity(t *testing.T) {
	boardState := &BoardState{}
	r := SquadRuleset{}
	err := r.shareSquadAttributes(boardState)
	require.NoError(t, err)
}

func TestSquadAllowBodyCollisions(t *testing.T) {
	testSnakes := []struct {
		SnakeID         string
		SquadID         string
		EliminatedCause string
		EliminatedBy    string
		ExpectedCause   string
		ExpectedBy      string
	}{
		// Red Squad
		{"R1", "red", NotEliminated, "", NotEliminated, ""},
		{"R2", "red", EliminatedByCollision, "R1", NotEliminated, ""},
		// Blue Squad
		{"B1", "blue", EliminatedByCollision, "R1", EliminatedByCollision, "R1"},
		{"B2", "blue", EliminatedBySelfCollision, "B1", EliminatedBySelfCollision, "B1"},
		{"B4", "blue", EliminatedByOutOfBounds, "", EliminatedByOutOfBounds, ""},
		{"B3", "blue", NotEliminated, "", NotEliminated, ""},
		// More Red Squad
		{"R3", "red", NotEliminated, "", NotEliminated, ""},
		{"R4", "red", EliminatedByCollision, "R4", EliminatedByCollision, "R4"}, // this is an error case but worth testing
		{"R5", "red", EliminatedByCollision, "R4", NotEliminated, ""},
		// Green Squad
		{"G1", "green", EliminatedByOutOfHealth, "x", EliminatedByOutOfHealth, "x"},
		// Yellow Squad
		{"Y1", "yellow", EliminatedByCollision, "B4", EliminatedByCollision, "B4"},
	}

	boardState := &BoardState{}
	squadMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:              testSnake.SnakeID,
			EliminatedCause: testSnake.EliminatedCause,
			EliminatedBy:    testSnake.EliminatedBy,
		})
		squadMap[testSnake.SnakeID] = testSnake.SquadID
	}
	require.Equal(t, len(squadMap), len(boardState.Snakes), "squad map is wrong size, error in test setup")

	r := SquadRuleset{SquadMap: squadMap, AllowBodyCollisions: true}
	err := r.resurrectSquadBodyCollisions(boardState)

	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(testSnakes))
	for i := 0; i < len(boardState.Snakes); i++ {
		require.Equal(
			t,
			testSnakes[i].ExpectedCause,
			boardState.Snakes[i].EliminatedCause,
			"snake %s failed shared eliminated cause",
			testSnakes[i].SnakeID,
		)
		require.Equal(
			t,
			testSnakes[i].ExpectedBy,
			boardState.Snakes[i].EliminatedBy,
			"snake %s failed shared eliminated by",
			testSnakes[i].SnakeID,
		)
	}
}

func TestSquadAllowBodyCollisionsEliminatedByNotSet(t *testing.T) {
	boardState := &BoardState{
		Snakes: []Snake{
			Snake{ID: "1", EliminatedCause: EliminatedByCollision},
			Snake{ID: "2"},
		},
	}
	r := SquadRuleset{
		AllowBodyCollisions: true,
		SquadMap: map[string]string{
			"1": "red",
			"2": "red",
		},
	}
	err := r.resurrectSquadBodyCollisions(boardState)
	require.Error(t, err)
}

func TestSquadShareSquadHealth(t *testing.T) {
	testSnakes := []struct {
		SnakeID        string
		SquadID        string
		Health         int32
		ExpectedHealth int32
	}{
		// Red Squad
		{"R1", "red", 11, 88},
		{"R2", "red", 22, 88},
		// Blue Squad
		{"B1", "blue", 33, 333},
		{"B2", "blue", 333, 333},
		{"B3", "blue", 3, 333},
		// More Red Squad
		{"R3", "red", 77, 88},
		{"R4", "red", 88, 88},
		// Green Squad
		{"G1", "green", 100, 100},
		// Yellow Squad
		{"Y1", "yellow", 1, 1},
	}

	boardState := &BoardState{}
	squadMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:     testSnake.SnakeID,
			Health: testSnake.Health,
		})
		squadMap[testSnake.SnakeID] = testSnake.SquadID
	}
	require.Equal(t, len(squadMap), len(boardState.Snakes), "squad map is wrong size, error in test setup")

	r := SquadRuleset{SharedHealth: true, SquadMap: squadMap}
	err := r.shareSquadAttributes(boardState)

	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(testSnakes))
	for i := 0; i < len(boardState.Snakes); i++ {
		require.Equal(
			t,
			testSnakes[i].ExpectedHealth,
			boardState.Snakes[i].Health,
			"snake %s failed shared health",
			testSnakes[i].SnakeID,
		)
	}
}

func TestSquadSharedLength(t *testing.T) {
	testSnakes := []struct {
		SnakeID      string
		SquadID      string
		Body         []Point
		ExpectedBody []Point
	}{
		// Red Squad
		{"R1", "red", []Point{{1, 1}}, []Point{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}},
		{"R2", "red", []Point{{2, 2}, {2, 2}}, []Point{{2, 2}, {2, 2}, {2, 2}, {2, 2}, {2, 2}}},
		// Blue Squad
		{"B1", "blue", []Point{{1, 1}, {1, 2}}, []Point{{1, 1}, {1, 2}}},
		{"B2", "blue", []Point{{2, 1}}, []Point{{2, 1}, {2, 1}}},
		{"B3", "blue", []Point{{3, 3}}, []Point{{3, 3}, {3, 3}}},
		// More Red Squad
		{"R3", "red", []Point{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}}, []Point{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}}},
		{"R4", "red", []Point{{4, 4}}, []Point{{4, 4}, {4, 4}, {4, 4}, {4, 4}, {4, 4}}},
		// Green Squad
		{"G1", "green", []Point{{1, 1}}, []Point{{1, 1}}},
		// Yellow Squad
		{"Y1", "yellow", []Point{{1, 3}, {1, 4}, {1, 5}, {1, 6}}, []Point{{1, 3}, {1, 4}, {1, 5}, {1, 6}}},
	}

	boardState := &BoardState{}
	squadMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:   testSnake.SnakeID,
			Body: testSnake.Body,
		})
		squadMap[testSnake.SnakeID] = testSnake.SquadID
	}
	require.Equal(t, len(squadMap), len(boardState.Snakes), "squad map is wrong size, error in test setup")

	r := SquadRuleset{SharedLength: true, SquadMap: squadMap}
	err := r.shareSquadAttributes(boardState)

	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(testSnakes))
	for i := 0; i < len(boardState.Snakes); i++ {
		require.Equal(
			t,
			testSnakes[i].ExpectedBody,
			boardState.Snakes[i].Body,
			"snake %s failed shared length",
			testSnakes[i].SnakeID,
		)
	}
}

func TestSquadSharedElimination(t *testing.T) {
	testSnakes := []struct {
		SnakeID         string
		SquadID         string
		EliminatedCause string
		EliminatedBy    string
		ExpectedCause   string
		ExpectedBy      string
	}{
		// Red Squad
		{"R1", "red", NotEliminated, "", EliminatedBySquad, ""},
		{"R2", "red", EliminatedByHeadToHeadCollision, "y", EliminatedByHeadToHeadCollision, "y"},
		// Blue Squad
		{"B1", "blue", EliminatedByOutOfBounds, "z", EliminatedByOutOfBounds, "z"},
		{"B2", "blue", NotEliminated, "", EliminatedBySquad, ""},
		{"B3", "blue", NotEliminated, "", EliminatedBySquad, ""},
		// More Red Squad
		{"R3", "red", NotEliminated, "", EliminatedBySquad, ""},
		{"R4", "red", EliminatedByCollision, "B1", EliminatedByCollision, "B1"},
		// Green Squad
		{"G1", "green", EliminatedByOutOfHealth, "x", EliminatedByOutOfHealth, "x"},
		// Yellow Squad
		{"Y1", "yellow", NotEliminated, "", NotEliminated, ""},
	}

	boardState := &BoardState{}
	squadMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:              testSnake.SnakeID,
			EliminatedCause: testSnake.EliminatedCause,
			EliminatedBy:    testSnake.EliminatedBy,
		})
		squadMap[testSnake.SnakeID] = testSnake.SquadID
	}
	require.Equal(t, len(squadMap), len(boardState.Snakes), "squad map is wrong size, error in test setup")

	r := SquadRuleset{SharedElimination: true, SquadMap: squadMap}
	err := r.shareSquadAttributes(boardState)

	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(testSnakes))
	for i := 0; i < len(boardState.Snakes); i++ {
		require.Equal(
			t,
			testSnakes[i].ExpectedCause,
			boardState.Snakes[i].EliminatedCause,
			"snake %s failed shared eliminated cause",
			testSnakes[i].SnakeID,
		)
		require.Equal(
			t,
			testSnakes[i].ExpectedBy,
			boardState.Snakes[i].EliminatedBy,
			"snake %s failed shared eliminated by",
			testSnakes[i].SnakeID,
		)
	}
}

func TestSquadSharedAttributesErrorLengthZero(t *testing.T) {
	boardState := &BoardState{
		Snakes: []Snake{
			Snake{ID: "1"},
			Snake{ID: "2"},
		},
	}
	r := SquadRuleset{
		SharedLength: true,
		SquadMap: map[string]string{
			"1": "red",
			"2": "red",
		},
	}
	err := r.shareSquadAttributes(boardState)
	require.Error(t, err)
}

func TestSquadIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []Snake
		SquadMap map[string]string
		Expected bool
	}{
		{[]Snake{}, map[string]string{}, true},
		{[]Snake{{ID: "R1"}}, map[string]string{"R1": "red"}, true},
		{
			[]Snake{{ID: "R1"}, {ID: "R2"}, {ID: "R3"}},
			map[string]string{"R1": "red", "R2": "red", "R3": "red"},
			true,
		},
		{
			[]Snake{{ID: "R1"}, {ID: "B1"}},
			map[string]string{"R1": "red", "B1": "blue"},
			false,
		},
		{
			[]Snake{{ID: "R1"}, {ID: "B1"}, {ID: "B2"}, {ID: "G1"}},
			map[string]string{"R1": "red", "B1": "blue", "B2": "blue", "G1": "green"},
			false,
		},
		{
			[]Snake{
				{ID: "R1", EliminatedCause: EliminatedByOutOfBounds},
				{ID: "B1", EliminatedCause: EliminatedBySelfCollision, EliminatedBy: "B1"},
				{ID: "B2", EliminatedCause: EliminatedByCollision, EliminatedBy: "B2"},
				{ID: "G1"},
			},
			map[string]string{"R1": "red", "B1": "blue", "B2": "blue", "G1": "green"},
			true,
		},
	}

	for _, test := range tests {
		b := &BoardState{
			Height: 11,
			Width:  11,
			Snakes: test.Snakes,
			Food:   []Point{},
		}
		r := SquadRuleset{SquadMap: test.SquadMap}

		actual, err := r.IsGameOver(b)
		require.NoError(t, err)
		require.Equal(t, test.Expected, actual)
	}
}

func TestRegressionIssue16(t *testing.T) {
	// This is a specific test case to detect this issue:
	// https://github.com/BattlesnakeOfficial/rules/issues/16
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "teamBoi", Health: 10, Body: []Point{{1, 4}, {1, 3}, {0, 3}, {0, 2}, {1, 2}, {2, 2}}},
			{ID: "Node-Red-Bellied-Black-Snake", Health: 10, Body: []Point{{1, 8}, {2, 8}, {2, 9}, {3, 9}, {4, 9}, {4, 10}}},
			{ID: "Crash Override", Health: 10, Body: []Point{{2, 7}, {2, 6}, {3, 6}, {4, 6}, {4, 5}, {5, 5}, {6, 5}}},
			{ID: "Zero Cool", Health: 10, Body: []Point{{6, 5}, {5, 5}, {5, 4}, {5, 3}, {4, 3}, {3, 3}, {3, 4}}},
		},
	}
	squadMap := map[string]string{
		"teamBoi":                      "BirdSnakers",
		"Node-Red-Bellied-Black-Snake": "BirdSnakers",
		"Crash Override":               "Hackers",
		"Zero Cool":                    "Hackers",
	}
	snakeMoves := []SnakeMove{
		{ID: "teamBoi", Move: "up"},
		{ID: "Node-Red-Bellied-Black-Snake", Move: "left"},
		{ID: "Crash Override", Move: "left"},
		{ID: "Zero Cool", Move: "left"},
	}

	require.Equal(t, len(squadMap), len(boardState.Snakes), "squad map is wrong size, error in test setup")

	r := SquadRuleset{
		AllowBodyCollisions: true,
		SquadMap:            squadMap,
	}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "teamBoi", Body: []Point{{1, 5}, {1, 4}, {1, 3}, {0, 3}, {0, 2}, {1, 2}}},
		{ID: "Node-Red-Bellied-Black-Snake", Body: []Point{{0, 8}, {1, 8}, {2, 8}, {2, 9}, {3, 9}, {4, 9}}},
		{ID: "Crash Override", Body: []Point{{1, 7}, {2, 7}, {2, 6}, {3, 6}, {4, 6}, {4, 5}, {5, 5}}},
		{ID: "Zero Cool", Body: []Point{{5, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {4, 3}, {3, 3}}, EliminatedCause: EliminatedBySelfCollision, EliminatedBy: "Zero Cool"},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

var squadCaseMoveSquadCollisions = gameTestCase{
	"Squad Case Move Squad Collisions",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "snake1squad1",
				Body:   []Point{{1, 1}, {2, 1}},
				Health: 100,
			},
			{
				ID:     "snake2squad1",
				Body:   []Point{{1, 2}, {2, 2}},
				Health: 100,
			},
			{
				ID:     "snake3squad2",
				Body:   []Point{{4, 4}, {4, 5}},
				Health: 100,
			},
			{
				ID:     "snake4squad2",
				Body:   []Point{{5, 4}, {5, 5}},
				Health: 100,
			},
		},
		Food:    []Point{},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "snake1squad1", Move: MoveUp},
		{ID: "snake2squad1", Move: MoveDown},
		{ID: "snake3squad2", Move: MoveRight},
		{ID: "snake4squad2", Move: MoveLeft},
	},
	nil,
	&BoardState{Width: 10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "snake1squad1",
				Body:   []Point{{1, 2}, {1, 1}},
				Health: 99,
			},
			{
				ID:     "snake2squad1",
				Body:   []Point{{1, 1}, {1, 2}},
				Health: 99,
			},
			{
				ID:     "snake3squad2",
				Body:   []Point{{5, 4}, {4, 4}},
				Health: 99,
			},
			{
				ID:     "snake4squad2",
				Body:   []Point{{4, 4}, {5, 4}},
				Health: 99,
			},
		},
		Food:    []Point{},
		Hazards: []Point{}},
}

func TestSquadCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		// inherits these test cases from standard
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		standardCaseMoveEatAndGrow,
		squadCaseMoveSquadCollisions,
	}
	rand.Seed(0)
	r := SquadRuleset{
		AllowBodyCollisions: true,
		SquadMap: map[string]string{
			"snake1squad1": "squad1",
			"snake2squad1": "squad1",
			"snake3squad2": "squad2",
			"snake4squad2": "squad2",
		},
	}
	for i, gc := range cases {
		t.Logf("Running test case %d", i)
		gc.requireCasesEqual(t, &r)
	}
}
