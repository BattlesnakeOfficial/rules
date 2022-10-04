package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/board"
	"github.com/BattlesnakeOfficial/rules/client"
	"github.com/BattlesnakeOfficial/rules/test"
	"github.com/stretchr/testify/require"
)

func buildDefaultGameState() *GameState {
	gameState := &GameState{
		Width:               11,
		Height:              11,
		Names:               nil,
		Timeout:             500,
		Sequential:          false,
		GameType:            "standard",
		MapName:             "standard",
		ViewMap:             false,
		UseColor:            false,
		Seed:                1,
		TurnDelay:           0,
		TurnDuration:        0,
		Output:              "",
		FoodSpawnChance:     15,
		MinimumFood:         1,
		HazardDamagePerTurn: 14,
		ShrinkEveryNTurns:   25,
	}

	return gameState
}

func TestGetIndividualBoardStateForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}

	gameState := buildDefaultGameState()
	gameState.initialize()
	gameState.gameID = "GAME_ID"
	gameState.snakeStates = map[string]SnakeState{
		s1State.ID: s1State,
		s2State.ID: s2State,
	}

	snakeRequest := gameState.getRequestBodyForSnake(state, s1State)
	requestBody := serialiseSnakeRequest(snakeRequest)

	test.RequireJSONMatchesFixture(t, "testdata/snake_request_body.json", string(requestBody))
}

func TestSettingsRequestSerialization(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}

	for _, gt := range []string{
		rules.GameTypeStandard, rules.GameTypeRoyale, rules.GameTypeSolo,
		rules.GameTypeWrapped, rules.GameTypeConstrictor,
	} {
		t.Run(gt, func(t *testing.T) {
			gameState := buildDefaultGameState()

			gameState.FoodSpawnChance = 11
			gameState.MinimumFood = 7
			gameState.HazardDamagePerTurn = 19
			gameState.ShrinkEveryNTurns = 17
			gameState.GameType = gt

			gameState.initialize()
			gameState.gameID = "GAME_ID"
			gameState.snakeStates = map[string]SnakeState{s1State.ID: s1State, s2State.ID: s2State}

			snakeRequest := gameState.getRequestBodyForSnake(state, s1State)
			requestBody := serialiseSnakeRequest(snakeRequest)
			t.Log(string(requestBody))

			test.RequireJSONMatchesFixture(t, fmt.Sprintf("testdata/snake_request_body_%s.json", gt), string(requestBody))
		})
	}
}

func TestConvertRulesSnakes(t *testing.T) {
	tests := []struct {
		name     string
		snakes   []rules.Snake
		state    map[string]SnakeState
		expected []client.Snake
	}{
		{
			name:     "empty",
			snakes:   []rules.Snake{},
			state:    map[string]SnakeState{},
			expected: []client.Snake{},
		},
		{
			name: "all properties",
			snakes: []rules.Snake{
				{ID: "one", Body: []rules.Point{{X: 3, Y: 3}, {X: 2, Y: 3}}, Health: 100},
			},
			state: map[string]SnakeState{
				"one": {
					ID:        "one",
					Name:      "ONE",
					URL:       "http://example1.com",
					Head:      "a",
					Tail:      "b",
					Color:     "#012345",
					LastMove:  "up",
					Character: '+',
				},
			},
			expected: []client.Snake{
				{
					ID:      "one",
					Name:    "ONE",
					Latency: "0",
					Health:  100,
					Body:    []client.Coord{{X: 3, Y: 3}, {X: 2, Y: 3}},
					Head:    client.Coord{X: 3, Y: 3},
					Length:  2,
					Shout:   "",
					Customizations: client.Customizations{
						Color: "#012345",
						Head:  "a",
						Tail:  "b",
					},
				},
			},
		},
		{
			name: "some eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
				{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
				"two": {ID: "two"},
			},
			expected: []client.Snake{
				{
					ID:      "two",
					Latency: "0",
					Body:    []client.Coord{{X: 4, Y: 3}},
					Head:    client.Coord{X: 4, Y: 3},
					Length:  1,
				},
			},
		},
		{
			name: "all eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
			},
			expected: []client.Snake{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := convertRulesSnakes(test.snakes, test.state)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestBuildFrameEvent(t *testing.T) {
	tests := []struct {
		name        string
		boardState  *rules.BoardState
		snakeStates map[string]SnakeState
		expected    board.GameEvent
	}{
		{
			name:        "empty",
			boardState:  rules.NewBoardState(11, 11),
			snakeStates: map[string]SnakeState{},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,
				Data: board.GameFrame{
					Turn:    0,
					Snakes:  []board.Snake{},
					Food:    []rules.Point{},
					Hazards: []rules.Point{},
				},
			},
		},
		{
			name: "snake fields",
			boardState: &rules.BoardState{
				Turn:   99,
				Height: 19,
				Width:  25,
				Food:   []rules.Point{{X: 9, Y: 4}},
				Snakes: []rules.Snake{
					{
						ID: "1",
						Body: []rules.Point{
							{X: 9, Y: 4},
							{X: 8, Y: 4},
							{X: 7, Y: 4},
						},
						Health:           97,
						EliminatedCause:  rules.EliminatedBySelfCollision,
						EliminatedOnTurn: 45,
						EliminatedBy:     "1",
					},
				},
				Hazards: []rules.Point{{X: 8, Y: 6}},
			},
			snakeStates: map[string]SnakeState{
				"1": {
					URL:        "http://example.com",
					Name:       "One",
					ID:         "1",
					LastMove:   "left",
					Color:      "#ff00ff",
					Head:       "silly",
					Tail:       "default",
					Author:     "AUTHOR",
					Version:    "1.5",
					Error:      nil,
					StatusCode: 200,
					Latency:    54 * time.Millisecond,
				},
			},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,

				Data: board.GameFrame{
					Turn: 99,
					Snakes: []board.Snake{
						{
							ID:     "1",
							Name:   "One",
							Body:   []rules.Point{{X: 9, Y: 4}, {X: 8, Y: 4}, {X: 7, Y: 4}},
							Health: 97,
							Death: &board.Death{
								Cause:        rules.EliminatedBySelfCollision,
								Turn:         45,
								EliminatedBy: "1",
							},
							Color:         "#ff00ff",
							HeadType:      "silly",
							TailType:      "default",
							Latency:       "54",
							Author:        "AUTHOR",
							StatusCode:    200,
							Error:         "",
							IsBot:         false,
							IsEnvironment: false,
						},
					},
					Food:    []rules.Point{{X: 9, Y: 4}},
					Hazards: []rules.Point{{X: 8, Y: 6}},
				},
			},
		},
		{
			name: "snake errors",
			boardState: &rules.BoardState{
				Height: 19,
				Width:  25,
				Snakes: []rules.Snake{
					{
						ID: "bad_status",
					},
					{
						ID: "connection_error",
					},
				},
			},
			snakeStates: map[string]SnakeState{
				"bad_status": {
					StatusCode: 504,
					Latency:    54 * time.Millisecond,
				},
				"connection_error": {
					Error:   fmt.Errorf("error connecting to host"),
					Latency: 0,
				},
			},
			expected: board.GameEvent{
				EventType: board.EVENT_TYPE_FRAME,

				Data: board.GameFrame{
					Snakes: []board.Snake{
						{
							ID:         "bad_status",
							Latency:    "54",
							StatusCode: 504,
							Error:      "7:Bad HTTP status code 504",
						},
						{
							ID:         "connection_error",
							Latency:    "1",
							StatusCode: 0,
							Error:      "0:Error communicating with server",
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gameState := GameState{
				snakeStates: test.snakeStates,
			}
			actual := gameState.buildFrameEvent(test.boardState)
			require.Equalf(t, test.expected, actual, "%#v", actual)
		})
	}
}

func TestGetMoveForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	boardState := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}

	tests := []struct {
		name            string
		boardState      *rules.BoardState
		snakeState      SnakeState
		responseErr     error
		responseCode    int
		responseBody    string
		responseLatency time.Duration

		expectedMove  rules.SnakeMove
		expectedState SnakeState
	}{
		{
			name:       "invalid URL",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "",
				LastMove: rules.MoveLeft,
			},
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveLeft,
			},
		},
		{
			name:       "error response",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseErr: errors.New("connection error"),
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveLeft,
			},
		},
		{
			name:       "bad response body",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    200,
			responseBody:    `right`,
			responseLatency: 54 * time.Millisecond,
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveLeft,
			},
		},
		{
			name:       "bad move value",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    200,
			responseBody:    `{"move": "north"}`,
			responseLatency: 54 * time.Millisecond,
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveLeft,
			},
		},
		{
			name:       "bad status code",
			boardState: boardState,
			snakeState: SnakeState{
				ID:       "one",
				URL:      "http://example.com",
				LastMove: rules.MoveLeft,
			},
			responseCode:    500,
			responseLatency: 54 * time.Millisecond,
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveLeft,
			},
		},
		{
			name:       "successful move",
			boardState: boardState,
			snakeState: SnakeState{
				ID:  "one",
				URL: "http://example.com",
			},
			responseCode:    200,
			responseBody:    `{"move": "right"}`,
			responseLatency: 54 * time.Millisecond,
			expectedMove: rules.SnakeMove{
				ID:   "one",
				Move: rules.MoveRight,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gameState := buildDefaultGameState()
			gameState.initialize()
			gameState.snakeStates = map[string]SnakeState{test.snakeState.ID: test.snakeState}
			gameState.httpClient = stubHTTPClient{test.responseErr, test.responseCode, test.responseBody, test.responseLatency}

			move, statusCode, latency := gameState.getMoveForSnake(test.boardState, test.snakeState)
			require.Equal(t, test.expectedMove, move)
			require.Equal(t, test.responseCode, statusCode)
			require.Equal(t, test.responseLatency, latency)
		})
	}
}

func TestCreateNextBoardState(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	boardState := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1},
	}
	snakeState := SnakeState{
		ID:  s1.ID,
		URL: "http://example.com",
	}

	for _, sequential := range []bool{false, true} {
		t.Run(fmt.Sprintf("sequential_%v", sequential), func(t *testing.T) {
			gameState := buildDefaultGameState()
			gameState.Sequential = sequential
			gameState.initialize()
			gameState.snakeStates = map[string]SnakeState{s1.ID: snakeState}
			gameState.httpClient = stubHTTPClient{nil, 200, `{"move": "right"}`, 54 * time.Millisecond}

			nextBoardState := gameState.createNextBoardState(boardState)
			snakeState = gameState.snakeStates[s1.ID]

			require.NotNil(t, nextBoardState)
			require.Equal(t, nextBoardState.Turn, 1)
			require.Equal(t, nextBoardState.Snakes[0].Body[0], rules.Point{X: 4, Y: 3})
			require.Equal(t, snakeState.LastMove, rules.MoveRight)
			require.Equal(t, snakeState.StatusCode, 200)
			require.Equal(t, snakeState.Latency, 54*time.Millisecond)
		})
	}
}

type stubHTTPClient struct {
	err        error
	statusCode int
	body       string
	latency    time.Duration
}

func (client stubHTTPClient) request() (*http.Response, time.Duration, error) {
	if client.err != nil {
		return nil, client.latency, client.err
	}
	body := ioutil.NopCloser(bytes.NewBufferString(client.body))

	response := &http.Response{
		Header:     make(http.Header),
		Body:       body,
		StatusCode: client.statusCode,
	}

	return response, client.latency, nil
}

func (client stubHTTPClient) Get(url string) (*http.Response, time.Duration, error) {
	return client.request()
}

func (client stubHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, time.Duration, error) {
	return client.request()
}
