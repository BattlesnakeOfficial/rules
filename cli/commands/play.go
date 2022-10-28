package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/board"
	"github.com/BattlesnakeOfficial/rules/client"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	log "github.com/spf13/jwalterweatherman"
)

// Used to store state for each SnakeState while running a local game
type SnakeState struct {
	URL        string
	Name       string
	ID         string
	LastMove   string
	Character  rune
	Color      string
	Head       string
	Tail       string
	Author     string
	Version    string
	Error      error
	StatusCode int
	Latency    time.Duration
}

type GameState struct {
	// Options
	Width               int
	Height              int
	Names               []string
	URLs                []string
	Timeout             int
	TurnDuration        int
	Sequential          bool
	GameType            string
	MapName             string
	ViewMap             bool
	UseColor            bool
	Seed                int64
	TurnDelay           int
	OutputPath          string
	ViewInBrowser       bool
	BoardURL            string
	FoodSpawnChance     int
	MinimumFood         int
	HazardDamagePerTurn int
	ShrinkEveryNTurns   int

	// Internal game state
	settings    map[string]string
	snakeStates map[string]SnakeState
	gameID      string
	httpClient  TimedHttpClient
	ruleset     rules.Ruleset
	gameMap     maps.GameMap
	outputFile  io.WriteCloser
	idGenerator func(int) string
}

func NewPlayCommand() *cobra.Command {
	gameState := &GameState{}

	var playCmd = &cobra.Command{
		Use:   "play",
		Short: "Play a game of Battlesnake locally.",
		Long:  "Play a game of Battlesnake locally.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gameState.Initialize(); err != nil {
				log.ERROR.Fatalf("Error initializing game: %v", err)
			}
			if err := gameState.Run(); err != nil {
				log.ERROR.Fatalf("Error running game: %v", err)
			}
		},
	}

	playCmd.Flags().IntVarP(&gameState.Width, "width", "W", 11, "Width of Board")
	playCmd.Flags().IntVarP(&gameState.Height, "height", "H", 11, "Height of Board")
	playCmd.Flags().StringArrayVarP(&gameState.Names, "name", "n", nil, "Name of Snake")
	playCmd.Flags().StringArrayVarP(&gameState.URLs, "url", "u", nil, "URL of Snake")
	playCmd.Flags().IntVarP(&gameState.Timeout, "timeout", "t", 500, "Request Timeout")
	playCmd.Flags().BoolVarP(&gameState.Sequential, "sequential", "s", false, "Use Sequential Processing")
	playCmd.Flags().StringVarP(&gameState.GameType, "gametype", "g", "standard", "Type of Game Rules")
	playCmd.Flags().StringVarP(&gameState.MapName, "map", "m", "standard", "Game map to use to populate the board")
	playCmd.Flags().BoolVarP(&gameState.ViewMap, "viewmap", "v", false, "View the Map Each Turn")
	playCmd.Flags().BoolVarP(&gameState.UseColor, "color", "c", false, "Use color to draw the map")
	playCmd.Flags().Int64VarP(&gameState.Seed, "seed", "r", time.Now().UTC().UnixNano(), "Random Seed")
	playCmd.Flags().IntVarP(&gameState.TurnDelay, "delay", "d", 0, "Turn Delay in Milliseconds")
	playCmd.Flags().IntVarP(&gameState.TurnDuration, "duration", "D", 0, "Minimum Turn Duration in Milliseconds")
	playCmd.Flags().StringVarP(&gameState.OutputPath, "output", "o", "", "File path to output game state to. Existing files will be overwritten")
	playCmd.Flags().BoolVar(&gameState.ViewInBrowser, "browser", false, "View the game in the browser using the Battlesnake game board")
	playCmd.Flags().StringVar(&gameState.BoardURL, "board-url", "https://board.battlesnake.com", "Base URL for the game board when using --browser")

	playCmd.Flags().IntVar(&gameState.FoodSpawnChance, "foodSpawnChance", 15, "Percentage chance of spawning a new food every round")
	playCmd.Flags().IntVar(&gameState.MinimumFood, "minimumFood", 1, "Minimum food to keep on the board every turn")
	playCmd.Flags().IntVar(&gameState.HazardDamagePerTurn, "hazardDamagePerTurn", 14, "Health damage a snake will take when ending its turn in a hazard")
	playCmd.Flags().IntVar(&gameState.ShrinkEveryNTurns, "shrinkEveryNTurns", 25, "In Royale mode, the number of turns between generating new hazards")

	playCmd.Flags().SortFlags = false

	return playCmd
}

// Setup a GameState once all the fields have been parsed from the command-line.
func (gameState *GameState) Initialize() error {
	// Generate game ID
	gameState.gameID = uuid.New().String()

	// Set up HTTP client with request timeout
	if gameState.Timeout == 0 {
		gameState.Timeout = 500
	}
	gameState.httpClient = timedHTTPClient{
		&http.Client{
			Timeout: time.Duration(gameState.Timeout) * time.Millisecond,
		},
	}

	// Load game map
	gameMap, err := maps.GetMap(gameState.MapName)
	if err != nil {
		return fmt.Errorf("Failed to load game map %#v: %v", gameState.MapName, err)
	}
	gameState.gameMap = gameMap

	// Create settings object
	gameState.settings = map[string]string{
		rules.ParamFoodSpawnChance:     fmt.Sprint(gameState.FoodSpawnChance),
		rules.ParamMinimumFood:         fmt.Sprint(gameState.MinimumFood),
		rules.ParamHazardDamagePerTurn: fmt.Sprint(gameState.HazardDamagePerTurn),
		rules.ParamShrinkEveryNTurns:   fmt.Sprint(gameState.ShrinkEveryNTurns),
	}

	// Build ruleset from settings
	ruleset := rules.NewRulesetBuilder().
		WithSeed(gameState.Seed).
		WithParams(gameState.settings).
		WithSolo(len(gameState.URLs) < 2).
		NamedRuleset(gameState.GameType)
	gameState.ruleset = ruleset

	// Initialize snake states as empty until we can ping the snake URLs
	gameState.snakeStates = map[string]SnakeState{}

	if gameState.OutputPath != "" {
		f, err := os.OpenFile(gameState.OutputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("Failed to open output file: %w", err)
		}
		gameState.outputFile = f
	}

	return nil
}

// Setup and run a full game.
func (gameState *GameState) Run() error {
	var gameOver bool
	var err error

	// Setup local state for snakes
	gameState.snakeStates, err = gameState.buildSnakesFromOptions()
	if err != nil {
		return fmt.Errorf("Error getting snake metadata: %w", err)
	}

	rand.Seed(gameState.Seed)

	gameOver, boardState, err := gameState.initializeBoardFromArgs()
	if err != nil {
		return fmt.Errorf("Error initializing board: %w", err)
	}

	gameExporter := GameExporter{
		game:          gameState.createClientGame(),
		snakeRequests: make([]client.SnakeRequest, 0),
		winner:        SnakeState{},
		isDraw:        false,
	}
	exportGame := gameState.outputFile != nil
	if exportGame {
		defer gameState.outputFile.Close()
	}

	boardGame := board.Game{
		ID:     gameState.gameID,
		Status: "running",
		Width:  gameState.Width,
		Height: gameState.Height,
		Ruleset: map[string]string{
			rules.ParamGameType: gameState.GameType,
		},
		RulesetName: gameState.GameType,
		RulesStages: []string{},
		Map:         gameState.MapName,
	}
	boardServer := board.NewBoardServer(boardGame)

	if gameState.ViewInBrowser {
		serverURL, err := boardServer.Listen()
		if err != nil {
			return fmt.Errorf("Error starting HTTP server: %w", err)
		}
		defer boardServer.Shutdown()
		log.INFO.Printf("Board server listening on %s", serverURL)

		boardURL := fmt.Sprintf(gameState.BoardURL+"?engine=%s&game=%s&autoplay=true", serverURL, gameState.gameID)

		log.INFO.Printf("Opening board URL: %s", boardURL)
		if err := browser.OpenURL(boardURL); err != nil {
			log.ERROR.Printf("Failed to open browser: %v", err)
		}

		// send turn zero to websocket server
		boardServer.SendEvent(gameState.buildFrameEvent(boardState))
	}

	log.INFO.Printf("Ruleset: %v, Seed: %v", gameState.GameType, gameState.Seed)

	if gameState.ViewMap {
		gameState.printMap(boardState)
	} else {
		gameState.printState(boardState)
	}

	// Export game first, if enabled, so that we capture the request for turn zero.
	if exportGame {
		// The output file was designed in a way so that (nearly) every entry is equivalent to a valid API request.
		// This is meant to help unlock further development of tools such as replaying a saved game by simply copying each line and sending it as a POST request.
		// There was a design choice to be made here: the difference between SnakeRequest and BoardState is the `you` key.
		// We could choose to either store the SnakeRequest of each snake OR to omit the `you` key OR fill the `you` key with one of the snakes
		// In all cases the API request is technically non-compliant with how the actual API request should be.
		// The third option (filling the `you` key with an arbitrary snake) is the closest to the actual API request that would need the least manipulation to
		// be adjusted to look like an API call for a specific snake in the game.
		for _, snakeState := range gameState.snakeStates {
			snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
			gameExporter.AddSnakeRequest(snakeRequest)
			break
		}
	}

	var endTime time.Time
	for !gameOver {
		if gameState.TurnDuration > 0 {
			endTime = time.Now().Add(time.Duration(gameState.TurnDuration) * time.Millisecond)
		}

		gameOver, boardState, err = gameState.createNextBoardState(boardState)
		if err != nil {
			return fmt.Errorf("Error processing game: %w", err)
		}

		if gameOver {
			// Stop processing here - because game over is detected at the start of the pipeline, nothing will have changed.
			break
		}

		if gameState.ViewMap {
			gameState.printMap(boardState)
		} else {
			gameState.printState(boardState)
		}

		if gameState.TurnDelay > 0 {
			time.Sleep(time.Duration(gameState.TurnDelay) * time.Millisecond)
		}

		if gameState.TurnDuration > 0 {
			time.Sleep(time.Until(endTime))
		}

		if gameState.ViewInBrowser {
			boardServer.SendEvent(gameState.buildFrameEvent(boardState))
		}

		if exportGame {
			for _, snakeState := range gameState.snakeStates {
				snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
				gameExporter.AddSnakeRequest(snakeRequest)
				break
			}
		}
	}

	gameExporter.isDraw = false

	if len(gameState.snakeStates) > 1 {
		// A draw is possible if there is more than one snake in the game.
		gameExporter.isDraw = true
	}

	for _, snake := range boardState.Snakes {
		snakeState := gameState.snakeStates[snake.ID]
		if snake.EliminatedCause == rules.NotEliminated {
			gameExporter.isDraw = false
			gameExporter.winner = snakeState
		}

		gameState.sendEndRequest(boardState, snakeState)
	}

	if gameExporter.isDraw {
		log.INFO.Printf("Game completed after %v turns. It was a draw.", boardState.Turn)
	} else if gameExporter.winner.Name != "" {
		log.INFO.Printf("Game completed after %v turns. %v was the winner.", boardState.Turn, gameExporter.winner.Name)
	} else {
		log.INFO.Printf("Game completed after %v turns.", boardState.Turn)
	}

	if gameState.ViewInBrowser {
		boardServer.SendEvent(board.GameEvent{
			EventType: board.EVENT_TYPE_GAME_END,
			Data:      boardGame,
		})
	}

	if exportGame {
		lines, err := gameExporter.FlushToFile(gameState.outputFile)
		if err != nil {
			return fmt.Errorf("Unable to export game: %w", err)
		}
		log.INFO.Printf("Wrote %d lines to output file: %s", lines, gameState.OutputPath)
	}

	return nil
}

func (gameState *GameState) initializeBoardFromArgs() (bool, *rules.BoardState, error) {
	snakeIds := []string{}
	for _, snakeState := range gameState.snakeStates {
		snakeIds = append(snakeIds, snakeState.ID)
	}
	boardState, err := maps.SetupBoard(gameState.gameMap.ID(), gameState.ruleset.Settings(), gameState.Width, gameState.Height, snakeIds)
	if err != nil {
		return false, nil, fmt.Errorf("Error initializing BoardState with map: %w", err)
	}
	gameOver, boardState, err := gameState.ruleset.Execute(boardState, nil)
	if err != nil {
		return false, nil, fmt.Errorf("Error initializing BoardState with ruleset: %w", err)
	}

	for _, snakeState := range gameState.snakeStates {
		snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
		requestBody := serialiseSnakeRequest(snakeRequest)
		u, _ := url.ParseRequestURI(snakeState.URL)
		u.Path = path.Join(u.Path, "start")
		log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
		_, _, err = gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.WARN.Printf("Request to %v failed", u.String())
		}
	}
	return gameOver, boardState, nil
}

func (gameState *GameState) createNextBoardState(boardState *rules.BoardState) (bool, *rules.BoardState, error) {
	// apply PreUpdateBoard before making requests to snakes
	boardState, err := maps.PreUpdateBoard(gameState.gameMap, boardState, gameState.ruleset.Settings())
	if err != nil {
		return false, boardState, fmt.Errorf("Error pre-updating board with game map: %w", err)
	}

	// get moves from snakes
	stateUpdates := make(chan SnakeState, len(gameState.snakeStates))
	if gameState.Sequential {
		for _, snakeState := range gameState.snakeStates {
			for _, snake := range boardState.Snakes {
				if snakeState.ID == snake.ID && snake.EliminatedCause == rules.NotEliminated {
					nextSnakeState := gameState.getSnakeUpdate(boardState, snakeState)
					stateUpdates <- nextSnakeState
				}
			}
		}
		close(stateUpdates)
	} else {
		var wg sync.WaitGroup

		for _, snakeState := range gameState.snakeStates {
			for _, snake := range boardState.Snakes {
				if snakeState.ID == snake.ID && snake.EliminatedCause == rules.NotEliminated {
					wg.Add(1)
					go func(snakeState SnakeState) {
						defer wg.Done()
						nextSnakeState := gameState.getSnakeUpdate(boardState, snakeState)
						stateUpdates <- nextSnakeState
					}(snakeState)
				}
			}
		}

		wg.Wait()
		close(stateUpdates)
	}

	var moves []rules.SnakeMove
	for snakeState := range stateUpdates {
		gameState.snakeStates[snakeState.ID] = snakeState
		moves = append(moves, rules.SnakeMove{ID: snakeState.ID, Move: snakeState.LastMove})
	}

	gameOver, boardState, err := gameState.ruleset.Execute(boardState, moves)
	if err != nil {
		return false, boardState, fmt.Errorf("Error updating board state from ruleset: %w", err)
	}

	// apply PostUpdateBoard after ruleset operates on snake moves
	boardState, err = maps.PostUpdateBoard(gameState.gameMap, boardState, gameState.ruleset.Settings())
	if err != nil {
		return false, boardState, fmt.Errorf("Error post-updating board with game map: %w", err)
	}

	boardState.Turn += 1

	return gameOver, boardState, nil
}

func (gameState *GameState) getSnakeUpdate(boardState *rules.BoardState, snakeState SnakeState) SnakeState {
	snakeState.StatusCode = 0
	snakeState.Error = nil
	snakeState.Latency = 0

	snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
	requestBody := serialiseSnakeRequest(snakeRequest)

	u, err := url.ParseRequestURI(snakeState.URL)
	if err != nil {
		log.ERROR.Printf("Error parsing snake URL %#v: %v", snakeState.URL, err)
		snakeState.Error = err
		return snakeState
	}
	u.Path = path.Join(u.Path, "move")
	log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
	res, responseTime, err := gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))

	snakeState.Latency = responseTime

	if err != nil {
		log.WARN.Printf(
			"Request to %v failed\n"+
				"\tError: %s", u.String(), err)
		snakeState.Error = err
		return snakeState
	}

	snakeState.StatusCode = res.StatusCode

	if res.Body == nil {
		log.WARN.Printf(
			"Failed to parse response from %v\n"+
				"\tError: body is empty", u.String())
		return snakeState
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.WARN.Printf(
			"Failed to read response body from %v\n"+
				"\tError: %v", u.String(), readErr)
		snakeState.Error = readErr
		return snakeState
	}
	if res.StatusCode != http.StatusOK {
		log.WARN.Printf(
			"Got non-ok status code from %v\n"+
				"\tStatusCode: %d (expected %d)\n"+
				"\tBody: %q", u.String(), res.StatusCode, http.StatusOK, body)
		return snakeState
	}

	playerResponse := client.MoveResponse{}
	jsonErr := json.Unmarshal(body, &playerResponse)
	if jsonErr != nil {
		log.WARN.Printf(
			"Failed to decode JSON from %v\n"+
				"\tError: %v\n"+
				"\tBody: %q\n"+
				"\tSee https://docs.battlesnake.com/references/api#post-move", u.String(), jsonErr, body)
		snakeState.Error = jsonErr
		return snakeState
	}
	if playerResponse.Move != "up" && playerResponse.Move != "down" && playerResponse.Move != "left" && playerResponse.Move != "right" {
		log.WARN.Printf(
			"Failed to parse JSON data from %v\n"+
				"\tError: invalid move %q, valid moves are \"up\", \"down\", \"left\" or \"right\"\n"+
				"\tBody: %q\n"+
				"\tSee https://docs.battlesnake.com/references/api#post-move", u.String(), playerResponse.Move, body)
		return snakeState
	}

	snakeState.LastMove = playerResponse.Move

	return snakeState
}

func (gameState *GameState) sendEndRequest(boardState *rules.BoardState, snakeState SnakeState) {
	snakeRequest := gameState.getRequestBodyForSnake(boardState, snakeState)
	requestBody := serialiseSnakeRequest(snakeRequest)
	u, _ := url.ParseRequestURI(snakeState.URL)
	u.Path = path.Join(u.Path, "end")
	log.DEBUG.Printf("POST %s: %v", u, string(requestBody))
	_, _, err := gameState.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.WARN.Printf("Request to %v failed", u.String())
	}
}

func (gameState *GameState) getRequestBodyForSnake(boardState *rules.BoardState, snakeState SnakeState) client.SnakeRequest {
	var youSnake rules.Snake
	for _, snk := range boardState.Snakes {
		if snakeState.ID == snk.ID {
			youSnake = snk
			break
		}
	}
	request := client.SnakeRequest{
		Game:  gameState.createClientGame(),
		Turn:  boardState.Turn,
		Board: convertStateToBoard(boardState, gameState.snakeStates),
		You:   convertRulesSnake(youSnake, snakeState),
	}
	return request
}

func (gameState *GameState) createClientGame() client.Game {
	return client.Game{
		ID:      gameState.gameID,
		Timeout: gameState.Timeout,
		Ruleset: client.Ruleset{
			Name:     gameState.ruleset.Name(),
			Version:  "cli", // TODO: Use GitHub Release Version
			Settings: client.ConvertRulesetSettings(gameState.ruleset.Settings()),
		},
		Map: gameState.gameMap.ID(),
	}
}

func (gameState *GameState) buildSnakesFromOptions() (map[string]SnakeState, error) {
	bodyChars := []rune{'■', '⌀', '●', '☻', '◘', '☺', '□', '⍟'}
	var numSnakes int
	snakes := map[string]SnakeState{}
	numNames := len(gameState.Names)
	numURLs := len(gameState.URLs)
	if numNames > numURLs {
		numSnakes = numNames
	} else {
		numSnakes = numURLs
	}
	for i := int(0); i < numSnakes; i++ {
		var snakeName string
		var snakeURL string

		var id string
		if gameState.idGenerator != nil {
			id = gameState.idGenerator(i)
		} else {
			id = uuid.New().String()
		}

		if i < numNames {
			snakeName = gameState.Names[i]
		} else {
			log.DEBUG.Printf("Name for URL %v is missing: a name will be generated automatically", gameState.URLs[i])
			snakeName = GenerateSnakeName()
		}

		if i < numURLs {
			u, err := url.ParseRequestURI(gameState.URLs[i])
			if err != nil {
				return nil, fmt.Errorf("URL %v is not valid: %w", gameState.URLs[i], err)
			}
			snakeURL = u.String()
		} else {
			return nil, fmt.Errorf("URL for name %v is missing", gameState.Names[i])
		}

		snakeState := SnakeState{
			Name: snakeName, URL: snakeURL, ID: id, LastMove: "up", Character: bodyChars[i%8],
		}
		var snakeErr error
		res, _, err := gameState.httpClient.Get(snakeURL)
		if err != nil {
			return nil, fmt.Errorf("Snake metadata request to %v failed: %w", snakeURL, err)
		}

		snakeState.StatusCode = res.StatusCode

		if res.Body == nil {
			return nil, fmt.Errorf("Empty response body from snake metadata URL: %v", snakeURL)
		}

		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("Error reading from snake metadata URL %v: %w", snakeURL, readErr)
		}

		pingResponse := client.SnakeMetadataResponse{}
		jsonErr := json.Unmarshal(body, &pingResponse)
		if jsonErr != nil {
			return nil, fmt.Errorf("Failed to parse response from %v: %w", snakeURL, jsonErr)
		}

		snakeState.Head = pingResponse.Head
		snakeState.Tail = pingResponse.Tail
		snakeState.Color = pingResponse.Color
		snakeState.Author = pingResponse.Author
		snakeState.Version = pingResponse.Version

		if snakeErr != nil {
			snakeState.Error = snakeErr
		}

		snakes[snakeState.ID] = snakeState

		log.INFO.Printf("Snake ID: %v URL: %v, Name: \"%v\"", snakeState.ID, snakeURL, snakeState.Name)
	}
	return snakes, nil
}

func (gameState *GameState) printState(boardState *rules.BoardState) {
	var aliveSnakeNames []string
	for _, snake := range boardState.Snakes {
		if snake.EliminatedCause == rules.NotEliminated {
			aliveSnakeNames = append(aliveSnakeNames, gameState.snakeStates[snake.ID].Name)
		}
	}
	log.INFO.Printf(
		"Turn: %d, Snakes Alive: [%v], Food: %d, Hazards: %d",
		boardState.Turn, strings.Join(aliveSnakeNames, ", "), len(boardState.Food), len(boardState.Hazards),
	)
}

func (gameState *GameState) printMap(boardState *rules.BoardState) {
	var o bytes.Buffer
	o.WriteString(fmt.Sprintf("Turn: %d\n", boardState.Turn))
	board := make([][]string, boardState.Width)
	for i := range board {
		board[i] = make([]string, boardState.Height)
	}
	for y := int(0); y < boardState.Height; y++ {
		for x := int(0); x < boardState.Width; x++ {
			if gameState.UseColor {
				board[x][y] = TERM_FG_LIGHTGRAY + "□"
			} else {
				board[x][y] = "◦"
			}
		}
	}
	for _, oob := range boardState.Hazards {
		if gameState.UseColor {
			board[oob.X][oob.Y] = TERM_BG_GRAY + " " + TERM_BG_WHITE
		} else {
			board[oob.X][oob.Y] = "░"
		}
	}
	if gameState.UseColor {
		o.WriteString(fmt.Sprintf("Hazards "+TERM_BG_GRAY+" "+TERM_RESET+": %v\n", boardState.Hazards))
	} else {
		o.WriteString(fmt.Sprintf("Hazards ░: %v\n", boardState.Hazards))
	}
	for _, f := range boardState.Food {
		if gameState.UseColor {
			board[f.X][f.Y] = TERM_FG_FOOD + "●"
		} else {
			board[f.X][f.Y] = "⚕"
		}
	}
	if gameState.UseColor {
		o.WriteString(fmt.Sprintf("Food "+TERM_FG_FOOD+TERM_BG_WHITE+"●"+TERM_RESET+": %v\n", boardState.Food))
	} else {
		o.WriteString(fmt.Sprintf("Food ⚕: %v\n", boardState.Food))
	}
	for _, s := range boardState.Snakes {
		state := gameState.snakeStates[s.ID]

		red, green, blue := parseSnakeColor(state.Color)
		for _, b := range s.Body {
			if b.X >= 0 && b.X < boardState.Width && b.Y >= 0 && b.Y < boardState.Height {
				if gameState.UseColor {
					board[b.X][b.Y] = fmt.Sprintf(TERM_FG_RGB+"■", red, green, blue)
				} else {
					board[b.X][b.Y] = string(state.Character)
				}
			}
		}
		if gameState.UseColor {
			o.WriteString(fmt.Sprintf("%v "+TERM_FG_RGB+TERM_BG_WHITE+"■■■"+TERM_RESET+": ", state.Name, red, green, blue))
		} else {
			o.WriteString(fmt.Sprintf("%v %c: ", state.Name, state.Character))
		}
		o.WriteString(fmt.Sprintf("Health: %d", s.Health))
		if s.EliminatedCause != rules.NotEliminated {
			o.WriteString(fmt.Sprintf(", Eliminated: %v, Turn: %d", s.EliminatedCause, s.EliminatedOnTurn))
		}
		o.WriteString("\n")
	}
	for y := boardState.Height - 1; y >= 0; y-- {
		if gameState.UseColor {
			o.WriteString(TERM_BG_WHITE)
		}
		for x := int(0); x < boardState.Width; x++ {
			o.WriteString(board[x][y])
		}
		if gameState.UseColor {
			o.WriteString(TERM_RESET)
		}
		o.WriteString("\n")
	}
	fmt.Println(o.String())
}

func (gameState *GameState) buildFrameEvent(boardState *rules.BoardState) board.GameEvent {
	snakes := []board.Snake{}

	for _, snake := range boardState.Snakes {
		snakeState := gameState.snakeStates[snake.ID]

		latencyMS := snakeState.Latency.Milliseconds()
		// round up latency of 0 to 1, to avoid legacy error display in board
		if latencyMS == 0 {
			latencyMS = 1
		}
		convertedSnake := board.Snake{
			ID:            snake.ID,
			Name:          snakeState.Name,
			Body:          snake.Body,
			Health:        snake.Health,
			Color:         snakeState.Color,
			HeadType:      snakeState.Head,
			TailType:      snakeState.Tail,
			Author:        snakeState.Author,
			StatusCode:    snakeState.StatusCode,
			IsBot:         false,
			IsEnvironment: false,
			Latency:       fmt.Sprint(latencyMS),
		}
		if snakeState.Error != nil {
			// Instead of trying to keep in sync with the production engine's
			// error detection and messages, just show a generic error and rely
			// on the CLI logs to show what really happened.
			convertedSnake.Error = "0:Error communicating with server"
		} else if snakeState.StatusCode != http.StatusOK {
			convertedSnake.Error = fmt.Sprintf("7:Bad HTTP status code %d", snakeState.StatusCode)
		}
		if snake.EliminatedCause != rules.NotEliminated {
			convertedSnake.Death = &board.Death{
				Cause:        snake.EliminatedCause,
				Turn:         snake.EliminatedOnTurn,
				EliminatedBy: snake.EliminatedBy,
			}
		}
		snakes = append(snakes, convertedSnake)
	}

	gameFrame := board.GameFrame{
		Turn:    boardState.Turn,
		Snakes:  snakes,
		Food:    boardState.Food,
		Hazards: boardState.Hazards,
	}

	return board.GameEvent{
		EventType: board.EVENT_TYPE_FRAME,
		Data:      gameFrame,
	}
}

func serialiseSnakeRequest(snakeRequest client.SnakeRequest) []byte {
	requestJSON, err := json.Marshal(snakeRequest)
	if err != nil {
		// This is likely to be a programming error like a unsupported type or cyclical reference
		log.ERROR.Panicf("Error marshalling JSON from State: %v", err)
	}
	return requestJSON
}

func convertRulesSnake(snake rules.Snake, snakeState SnakeState) client.Snake {
	latencyMS := snakeState.Latency.Milliseconds()
	return client.Snake{
		ID:      snake.ID,
		Name:    snakeState.Name,
		Health:  snake.Health,
		Body:    client.CoordFromPointArray(snake.Body),
		Latency: fmt.Sprint(latencyMS),
		Head:    client.CoordFromPoint(snake.Body[0]),
		Length:  int(len(snake.Body)),
		Shout:   "",
		Customizations: client.Customizations{
			Head:  snakeState.Head,
			Tail:  snakeState.Tail,
			Color: snakeState.Color,
		},
	}
}

func convertRulesSnakes(snakes []rules.Snake, snakeStates map[string]SnakeState) []client.Snake {
	a := make([]client.Snake, 0)
	for _, snake := range snakes {
		if snake.EliminatedCause == rules.NotEliminated {
			a = append(a, convertRulesSnake(snake, snakeStates[snake.ID]))
		}
	}
	return a
}

func convertStateToBoard(boardState *rules.BoardState, snakeStates map[string]SnakeState) client.Board {
	return client.Board{
		Height:  boardState.Height,
		Width:   boardState.Width,
		Food:    client.CoordFromPointArray(boardState.Food),
		Hazards: client.CoordFromPointArray(boardState.Hazards),
		Snakes:  convertRulesSnakes(boardState.Snakes, snakeStates),
	}
}

// Parses a color string like "#ef03d3" to rgb values from 0 to 255 or returns
// the default gray if any errors occure
func parseSnakeColor(color string) (int64, int64, int64) {
	if len(color) == 7 {
		red, err_r := strconv.ParseInt(color[1:3], 16, 64)
		green, err_g := strconv.ParseInt(color[3:5], 16, 64)
		blue, err_b := strconv.ParseInt(color[5:], 16, 64)
		if err_r == nil && err_g == nil && err_b == nil {
			return red, green, blue
		}
	}
	// Default gray color from Battlesnake board
	return 136, 136, 136
}
