package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestLimitInfoMapInterface(t *testing.T) {
	var _ maps.GameMap = maps.LimitInfoMap{}
}

func TestLimitInfoMapSetupBoard(t *testing.T) {
	m := maps.LimitInfoMap{}
	settings := rules.Settings{}

	tests := []struct {
		name              string
		initialBoardState *rules.BoardState
		rand              rules.Rand

		expected *rules.BoardState
		err      error
	}{
		{
			"empty 7x7",
			rules.NewBoardState(7, 7),
			rules.MinRand,
			rules.NewBoardState(7, 7).WithFood([]rules.Point{{X: 3, Y: 3}}),
			nil,
		},
		{
			"not enough room for snakes 7x7",
			rules.NewBoardState(7, 7).WithSnakes(generateSnakes(17)),
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"not enough room for snakes 5x5",
			rules.NewBoardState(5, 5).WithSnakes(generateSnakes(14)),
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"full 11x11 min",
			rules.NewBoardState(11, 11).WithSnakes(generateSnakes(8)),
			rules.MinRand,
			rules.NewBoardState(11, 11).
				WithFood([]rules.Point{
					{X: 0, Y: 2},
					{X: 0, Y: 8},
					{X: 8, Y: 0},
					{X: 8, Y: 10},
					{X: 0, Y: 4},
					{X: 4, Y: 0},
					{X: 4, Y: 10},
					{X: 10, Y: 4},
					{X: 5, Y: 5},
				}).
				WithSnakes([]rules.Snake{
					{ID: "1", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
					{ID: "2", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
				}),
			nil,
		},
		{
			"full 11x11 max",
			rules.NewBoardState(11, 11).WithSnakes(generateSnakes(8)),
			rules.MaxRand,
			rules.NewBoardState(11, 11).
				WithFood([]rules.Point{
					{X: 6, Y: 0},
					{X: 6, Y: 10},
					{X: 10, Y: 6},
					{X: 0, Y: 6},
					{X: 2, Y: 10},
					{X: 10, Y: 2},
					{X: 10, Y: 8},
					{X: 2, Y: 0},
					{X: 5, Y: 5},
				}).
				WithSnakes([]rules.Snake{
					{ID: "1", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "2", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
				}),
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextBoardState := rules.NewBoardState(test.initialBoardState.Width, test.initialBoardState.Height)
			editor := maps.NewBoardStateEditor(nextBoardState)
			settings := settings.WithRand(test.rand)

			err := m.SetupBoard(test.initialBoardState, settings, editor)

			if test.err != nil {
				require.Equal(t, test.err, err)
			} else {
				require.Equalf(t, test.expected, nextBoardState, "%#v", nextBoardState.Food)
			}
		})
	}
}

func TestLimitInfoMapUpdateBoard(t *testing.T) {
	m := maps.LimitInfoMap{}

	tests := []struct {
		name              string
		initialBoardState *rules.BoardState
		settings          rules.Settings
		rand              rules.Rand

		expected *rules.BoardState
	}{
		{
			"empty no food",
			rules.NewBoardState(2, 2),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "0", rules.ParamMinimumFood, "0"),
			rules.MinRand,
			rules.NewBoardState(2, 2),
		},
		{
			"empty MinimumFood",
			rules.NewBoardState(2, 2),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "0", rules.ParamMinimumFood, "2"),
			rules.MinRand,
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}}).WithGameState(map[string]string{"food_spawn_0_0": "1", "food_spawn_0_1": "1"}),
		},
		{
			"not empty MinimumFood",
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 1}}).WithGameState(map[string]string{"food_spawn_0_1": "1"}),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "0", rules.ParamMinimumFood, "2"),
			rules.MinRand,
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 1}, {X: 0, Y: 0}}).WithGameState(map[string]string{"food_spawn_0_0": "1", "food_spawn_0_1": "1"}),
		},
		{
			"empty FoodSpawnChance inactive",
			rules.NewBoardState(2, 2),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "50", rules.ParamMinimumFood, "0"),
			rules.MinRand,
			rules.NewBoardState(2, 2),
		},
		{
			"empty FoodSpawnChance active",
			rules.NewBoardState(2, 2),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "50", rules.ParamMinimumFood, "0"),
			rules.MaxRand,
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 1}}).WithGameState(map[string]string{"food_spawn_0_1": "1"}),
		},
		{
			"not empty FoodSpawnChance active",
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}}).WithGameState(map[string]string{"food_spawn_0_0": "1"}),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "50", rules.ParamMinimumFood, "0"),
			rules.MaxRand,
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}}).WithGameState(map[string]string{"food_spawn_0_0": "1", "food_spawn_1_0": "1"}),
		},
		{
			"not empty FoodSpawnChance no room",
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}}),
			rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "50", rules.ParamMinimumFood, "0"),
			rules.MaxRand,
			rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextBoardState := test.initialBoardState.Clone()
			settings := test.settings.WithRand(test.rand)
			editor := maps.NewBoardStateEditor(nextBoardState)

			err := m.PostUpdateBoard(test.initialBoardState.Clone(), settings, editor)

			require.NoError(t, err)
			require.Equal(t, test.expected, nextBoardState)
		})
	}
}

func TestLimitInfoMapCleanGameState(t *testing.T) {
	m := maps.LimitInfoMap{}
	state := rules.NewBoardState(3, 3)
	state.GameState = map[string]string{
		"food_spawn_1_1": "9",
		"food_spawn_2_2": "10",
		"food_spawn_3_3": "11",
		"other":          "abc",
	}
	state.Turn = 10
	editor := maps.NewBoardStateEditor(state)
	r := rules.MinRand
	settings := rules.Settings{}.WithRand(r)

	// call exported func to trigger cleaning
	err := m.PostUpdateBoard(state, settings, editor)
	require.NoError(t, err)

	cleaned_gamestate := editor.GameState()

	// first entry should be gone
	if _, exists := cleaned_gamestate["food_spawn_1_1"]; exists {
		t.Errorf("expected old food_spawn_1_1 to be deleted")
	}

	// entry of current turn stays
	if _, exists := cleaned_gamestate["food_spawn_2_2"]; !exists {
		t.Errorf("expected food_spawn_2_2 to still exist")
	}

	//future entries stay
	if _, exists := cleaned_gamestate["food_spawn_3_3"]; !exists {
		t.Errorf("expected food_spawn_3_3 to still exist")
	}

	// other entry stays
	if _, exists := cleaned_gamestate["other"]; !exists {
		t.Errorf("expected other to still exist")
	}
}
