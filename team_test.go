package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestBodyPassthrough(t *testing.T) {
// 	r := TeamRuleset{
// 		BodyPassthrough: true,
// 		teams: map[string][]string{
// 			"A": {"1", "2"},
// 			"B": {"3"},
// 		},
// 	}
// 	initialState := &BoardState{
// 		Height: 10,
// 		Width:  10,
// 		Snakes: []Snake{
// 			{ID: "1", Health: 100, Body: []Point{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}},
// 			{ID: "2", Health: 100, Body: []Point{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}}},
// 			{ID: "3", Health: 100, Body: []Point{{X: 3, Y: 2}}},
// 		},
// 	}
// 	moves := []SnakeMove{
// 		{ID: "1", Move: "left"},
// 		{ID: "2", Move: "up"},
// 		{ID: "3", Move: "left"},
// 	}
// 	newState, err := r.CreateNextBoardState(initialState, moves)
// 	require.NoError(t, err)
// 	require.Equal(t, EliminatedByCollision, newState.Snakes[2].EliminatedCause)
// 	require.Equal(t, "2", newState.Snakes[2].EliminatedBy)

// 	require.Empty(t, newState.Snakes[1].EliminatedCause)
// 	require.Empty(t, newState.Snakes[1].EliminatedBy)

// 	r.BodyPassthrough = false

// 	newState, err = r.CreateNextBoardState(initialState, moves)
// 	require.NoError(t, err)

// 	require.Equal(t, EliminatedByCollision, newState.Snakes[1].EliminatedCause)
// 	require.Equal(t, "1", newState.Snakes[1].EliminatedBy)
// }

func TestCreateNextBoardStateSanity(t *testing.T) {
	boardState := &BoardState{}
	r := TeamRuleset{}
	_, err := r.CreateNextBoardState(boardState, []SnakeMove{})
	require.NoError(t, err)
}

func TestSharedAttributesSanity(t *testing.T) {
	boardState := &BoardState{}
	r := TeamRuleset{}
	err := r.shareTeamAttributes(boardState)
	require.NoError(t, err)
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
