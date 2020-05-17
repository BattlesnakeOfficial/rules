package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTeamRulesetInterface(t *testing.T) {
	var _ Ruleset = (*TeamRuleset)(nil)
}

func TestCreateNextBoardStateSanity(t *testing.T) {
	boardState := &BoardState{}
	r := TeamRuleset{}
	_, err := r.CreateNextBoardState(boardState, []SnakeMove{})
	require.NoError(t, err)
}

func TestResurrectTeamBodyCollisionsSanity(t *testing.T) {
	boardState := &BoardState{}
	r := TeamRuleset{}
	err := r.resurrectTeamBodyCollisions(boardState)
	require.NoError(t, err)
}

func TestSharedAttributesSanity(t *testing.T) {
	boardState := &BoardState{}
	r := TeamRuleset{}
	err := r.shareTeamAttributes(boardState)
	require.NoError(t, err)
}

func TestAllowBodyCollisions(t *testing.T) {
	testSnakes := []struct {
		SnakeID         string
		TeamID          string
		EliminatedCause string
		EliminatedBy    string
		ExpectedCause   string
		ExpectedBy      string
	}{
		// Team Red
		{"R1", "red", NotEliminated, "", NotEliminated, ""},
		{"R2", "red", EliminatedByCollision, "R1", NotEliminated, ""},
		// Team Blue
		{"B1", "blue", EliminatedByCollision, "R1", EliminatedByCollision, "R1"},
		{"B2", "blue", EliminatedBySelfCollision, "B1", EliminatedBySelfCollision, "B1"},
		{"B4", "blue", EliminatedByOutOfBounds, "", EliminatedByOutOfBounds, ""},
		{"B3", "blue", NotEliminated, "", NotEliminated, ""},
		// More Team Red
		{"R3", "red", NotEliminated, "", NotEliminated, ""},
		{"R4", "red", EliminatedByCollision, "R4", EliminatedByCollision, "R4"}, // this is an error case but worth testing
		{"R5", "red", EliminatedByCollision, "R4", NotEliminated, ""},
		// // Team Green
		{"G1", "green", EliminatedByStarvation, "x", EliminatedByStarvation, "x"},
		// // Team Yellow
		{"Y1", "yellow", EliminatedByCollision, "B4", EliminatedByCollision, "B4"},
	}

	boardState := &BoardState{}
	teamMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:              testSnake.SnakeID,
			EliminatedCause: testSnake.EliminatedCause,
			EliminatedBy:    testSnake.EliminatedBy,
		})
		teamMap[testSnake.SnakeID] = testSnake.TeamID
	}
	require.Equal(t, len(teamMap), len(boardState.Snakes), "team map is wrong size, error in test setup")

	r := TeamRuleset{TeamMap: teamMap, AllowBodyCollisions: true}
	err := r.resurrectTeamBodyCollisions(boardState)

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

func TestAllowBodyCollisionsEliminatedByNotSet(t *testing.T) {
	boardState := &BoardState{
		Snakes: []Snake{
			Snake{ID: "1", EliminatedCause: EliminatedByCollision},
			Snake{ID: "2"},
		},
	}
	r := TeamRuleset{
		AllowBodyCollisions: true,
		TeamMap: map[string]string{
			"1": "red",
			"2": "red",
		},
	}
	err := r.resurrectTeamBodyCollisions(boardState)
	require.Error(t, err)
}

func TestShareTeamHealth(t *testing.T) {
	testSnakes := []struct {
		SnakeID        string
		TeamID         string
		Health         int32
		ExpectedHealth int32
	}{
		// Team Red
		{"R1", "red", 11, 88},
		{"R2", "red", 22, 88},
		// Team Blue
		{"B1", "blue", 33, 333},
		{"B2", "blue", 333, 333},
		{"B3", "blue", 3, 333},
		// More Team Red
		{"R3", "red", 77, 88},
		{"R4", "red", 88, 88},
		// Team Green
		{"G1", "green", 100, 100},
		// Team Yellow
		{"Y1", "yellow", 1, 1},
	}

	boardState := &BoardState{}
	teamMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:     testSnake.SnakeID,
			Health: testSnake.Health,
		})
		teamMap[testSnake.SnakeID] = testSnake.TeamID
	}
	require.Equal(t, len(teamMap), len(boardState.Snakes), "team map is wrong size, error in test setup")

	r := TeamRuleset{SharedHealth: true, TeamMap: teamMap}
	err := r.shareTeamAttributes(boardState)

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

func TestSharedLength(t *testing.T) {
	testSnakes := []struct {
		SnakeID      string
		TeamID       string
		Body         []Point
		ExpectedBody []Point
	}{
		// Team Red
		{"R1", "red", []Point{{1, 1}}, []Point{{1, 1}, {1, 1}, {1, 1}, {1, 1}, {1, 1}}},
		{"R2", "red", []Point{{2, 2}, {2, 2}}, []Point{{2, 2}, {2, 2}, {2, 2}, {2, 2}, {2, 2}}},
		// Team Blue
		{"B1", "blue", []Point{{1, 1}, {1, 2}}, []Point{{1, 1}, {1, 2}}},
		{"B2", "blue", []Point{{2, 1}}, []Point{{2, 1}, {2, 1}}},
		{"B3", "blue", []Point{{3, 3}}, []Point{{3, 3}, {3, 3}}},
		// More Team Red
		{"R3", "red", []Point{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}}, []Point{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}}},
		{"R4", "red", []Point{{4, 4}}, []Point{{4, 4}, {4, 4}, {4, 4}, {4, 4}, {4, 4}}},
		// Team Green
		{"G1", "green", []Point{{1, 1}}, []Point{{1, 1}}},
		// Team Yellow
		{"Y1", "yellow", []Point{{1, 3}, {1, 4}, {1, 5}, {1, 6}}, []Point{{1, 3}, {1, 4}, {1, 5}, {1, 6}}},
	}

	boardState := &BoardState{}
	teamMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:   testSnake.SnakeID,
			Body: testSnake.Body,
		})
		teamMap[testSnake.SnakeID] = testSnake.TeamID
	}
	require.Equal(t, len(teamMap), len(boardState.Snakes), "team map is wrong size, error in test setup")

	r := TeamRuleset{SharedLength: true, TeamMap: teamMap}
	err := r.shareTeamAttributes(boardState)

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

func TestSharedElimination(t *testing.T) {
	testSnakes := []struct {
		SnakeID         string
		TeamID          string
		EliminatedCause string
		EliminatedBy    string
		ExpectedCause   string
		ExpectedBy      string
	}{
		// Team Red
		{"R1", "red", NotEliminated, "", EliminatedByTeam, ""},
		{"R2", "red", EliminatedByHeadToHeadCollision, "y", EliminatedByHeadToHeadCollision, "y"},
		// Team Blue
		{"B1", "blue", EliminatedByOutOfBounds, "z", EliminatedByOutOfBounds, "z"},
		{"B2", "blue", NotEliminated, "", EliminatedByTeam, ""},
		{"B3", "blue", NotEliminated, "", EliminatedByTeam, ""},
		// More Team Red
		{"R3", "red", NotEliminated, "", EliminatedByTeam, ""},
		{"R4", "red", EliminatedByCollision, "B1", EliminatedByCollision, "B1"},
		// Team Green
		{"G1", "green", EliminatedByStarvation, "x", EliminatedByStarvation, "x"},
		// Team Yellow
		{"Y1", "yellow", NotEliminated, "", NotEliminated, ""},
	}

	boardState := &BoardState{}
	teamMap := make(map[string]string)
	for _, testSnake := range testSnakes {
		boardState.Snakes = append(boardState.Snakes, Snake{
			ID:              testSnake.SnakeID,
			EliminatedCause: testSnake.EliminatedCause,
			EliminatedBy:    testSnake.EliminatedBy,
		})
		teamMap[testSnake.SnakeID] = testSnake.TeamID
	}
	require.Equal(t, len(teamMap), len(boardState.Snakes), "team map is wrong size, error in test setup")

	r := TeamRuleset{SharedElimination: true, TeamMap: teamMap}
	err := r.shareTeamAttributes(boardState)

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

func TestSharedAttributesErrorLengthZero(t *testing.T) {
	boardState := &BoardState{
		Snakes: []Snake{
			Snake{ID: "1"},
			Snake{ID: "2"},
		},
	}
	r := TeamRuleset{
		SharedLength: true,
		TeamMap: map[string]string{
			"1": "red",
			"2": "red",
		},
	}
	err := r.shareTeamAttributes(boardState)
	require.Error(t, err)
}

func TestTeamIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []Snake
		TeamMap  map[string]string
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
		r := TeamRuleset{TeamMap: test.TeamMap}

		actual, err := r.IsGameOver(b)
		require.NoError(t, err)
		require.Equal(t, test.Expected, actual)
	}
}
