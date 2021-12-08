package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Used to store state for each SnakeState while running a local game
type SnakeState struct {
	URL       string
	Name      string
	ID        string
	LastMove  string
	Squad     string
	Character rune
	Color     string
	Head      string
	Tail      string
}

var GameId string
var Turn int32
var HttpClient http.Client
var Width int32
var Height int32
var Names []string
var URLs []string
var Squads []string
var Timeout int32
var Sequential bool
var GameType string
var ViewMap bool
var UseColor bool
var Seed int64
var TurnDelay int32
var DebugRequests bool
var Output string

var FoodSpawnChance int32
var MinimumFood int32
var HazardDamagePerTurn int32
var ShrinkEveryNTurns int32

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play a game of Battlesnake locally.",
	Long:  "Play a game of Battlesnake locally.",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().Int32VarP(&Width, "width", "W", 11, "Width of Board")
	playCmd.Flags().Int32VarP(&Height, "height", "H", 11, "Height of Board")
	playCmd.Flags().StringArrayVarP(&Names, "name", "n", nil, "Name of Snake")
	playCmd.Flags().StringArrayVarP(&URLs, "url", "u", nil, "URL of Snake")
	playCmd.Flags().StringArrayVarP(&Names, "squad", "S", nil, "Squad of Snake")
	playCmd.Flags().Int32VarP(&Timeout, "timeout", "t", 500, "Request Timeout")
	playCmd.Flags().BoolVarP(&Sequential, "sequential", "s", false, "Use Sequential Processing")
	playCmd.Flags().StringVarP(&GameType, "gametype", "g", "standard", "Type of Game Rules")
	playCmd.Flags().BoolVarP(&ViewMap, "viewmap", "v", false, "View the Map Each Turn")
	playCmd.Flags().BoolVarP(&UseColor, "color", "c", false, "Use color to draw the map")
	playCmd.Flags().Int64VarP(&Seed, "seed", "r", time.Now().UTC().UnixNano(), "Random Seed")
	playCmd.Flags().Int32VarP(&TurnDelay, "delay", "d", 0, "Turn Delay in Milliseconds")
	playCmd.Flags().BoolVar(&DebugRequests, "debug-requests", false, "Log body of all requests sent")
	playCmd.Flags().StringVarP(&Output, "output", "o", "", "File path to output game state to. Existing files will be overwritten")

	playCmd.Flags().Int32Var(&FoodSpawnChance, "foodSpawnChance", 15, "Percentage chance of spawning a new food every round")
	playCmd.Flags().Int32Var(&MinimumFood, "minimumFood", 1, "Minimum food to keep on the board every turn")
	playCmd.Flags().Int32Var(&HazardDamagePerTurn, "hazardDamagePerTurn", 14, "Health damage a snake will take when ending its turn in a hazard")
	playCmd.Flags().Int32Var(&ShrinkEveryNTurns, "shrinkEveryNTurns", 25, "In Royale mode, the number of turns between generating new hazards")

	playCmd.Flags().SortFlags = false
}

var run = func(cmd *cobra.Command, args []string) {
	rand.Seed(Seed)

	GameId = uuid.New().String()
	Turn = 0

	snakeStates := buildSnakesFromOptions()

	ruleset := getRuleset(Seed, snakeStates)
	state := initializeBoardFromArgs(ruleset, snakeStates)
	exportGame := Output != ""

	gameExporter := GameExporter{
		game:          createClientGame(ruleset),
		snakeRequests: make([]client.SnakeRequest, 0),
		winner:        SnakeState{},
		isDraw:        false,
	}

	for v := false; !v; v, _ = ruleset.IsGameOver(state) {
		Turn++
		state = createNextBoardState(ruleset, state, snakeStates, Turn)

		if ViewMap {
			printMap(state, snakeStates, Turn)
		} else {
			log.Printf("[%v]: State: %v\n", Turn, state)
		}

		if TurnDelay > 0 {
			time.Sleep(time.Duration(TurnDelay) * time.Millisecond)
		}

		if exportGame {
			// The output file was designed in a way so that (nearly) every entry is equivalent to a valid API request.
			// This is meant to help unlock further development of tools such as replaying a saved game by simply copying each line and sending it as a POST request.
			// There was a design choice to be made here: the difference between SnakeRequest and BoardState is the `you` key.
			// We could choose to either store the SnakeRequest of each snake OR to omit the `you` key OR fill the `you` key with one of the snakes
			// In all cases the API request is technically non-compliant with how the actual API request should be.
			// The third option (filling the `you` key with an arbitrary snake) is the closest to the actual API request that would need the least manipulation to
			// be adjusted to look like an API call for a specific snake in the game.
			snakeState := snakeStates[state.Snakes[0].ID]
			snakeRequest := getIndividualBoardStateForSnake(state, snakeState, snakeStates, ruleset)
			gameExporter.AddSnakeRequest(snakeRequest)
		}
	}

	isDraw := true
	if GameType == "solo" {
		log.Printf("[DONE]: Game completed after %v turns.", Turn)
		if exportGame {
			// These checks for exportGame are present to avoid vacuuming up RAM when an export is not requred.
			gameExporter.winner = snakeStates[state.Snakes[0].ID]
		}
	} else {
		var winner SnakeState
		for _, snake := range state.Snakes {
			snakeState := snakeStates[snake.ID]
			if snake.EliminatedCause == rules.NotEliminated {
				isDraw = false
				winner = snakeState
			}
			sendEndRequest(ruleset, state, snakeState, snakeStates)
		}

		if isDraw {
			log.Printf("[DONE]: Game completed after %v turns. It was a draw.", Turn)
		} else {
			log.Printf("[DONE]: Game completed after %v turns. %v is the winner.", Turn, winner.Name)
		}
		if exportGame {
			gameExporter.winner = winner
		}
	}

	if exportGame {
		err := gameExporter.FlushToFile(Output, "JSONL")
		if err != nil {
			log.Printf("[WARN]: Unable to export game. Reason: %v\n", err.Error())
			os.Exit(1)
		}
	}
}

func getRuleset(seed int64, snakeStates map[string]SnakeState) rules.Ruleset {
	var ruleset rules.Ruleset
	var royale rules.RoyaleRuleset

	standard := rules.StandardRuleset{
		FoodSpawnChance:     FoodSpawnChance,
		MinimumFood:         MinimumFood,
		HazardDamagePerTurn: 0,
	}

	switch GameType {
	case "royale":
		standard.HazardDamagePerTurn = HazardDamagePerTurn
		royale = rules.RoyaleRuleset{
			StandardRuleset:   standard,
			Seed:              seed,
			ShrinkEveryNTurns: ShrinkEveryNTurns,
		}
		ruleset = &royale
	case "squad":
		squadMap := map[string]string{}
		for _, snakeState := range snakeStates {
			squadMap[snakeState.ID] = snakeState.Squad
		}
		ruleset = &rules.SquadRuleset{
			StandardRuleset:     standard,
			SquadMap:            squadMap,
			AllowBodyCollisions: true,
			SharedElimination:   true,
			SharedHealth:        true,
			SharedLength:        true,
		}
	case "solo":
		ruleset = &rules.SoloRuleset{
			StandardRuleset: standard,
		}
	case "wrapped":
		ruleset = &rules.WrappedRuleset{
			StandardRuleset: standard,
		}
	case "constrictor":
		ruleset = &rules.ConstrictorRuleset{
			StandardRuleset: standard,
		}
	default:
		ruleset = &standard
	}
	return ruleset
}

func initializeBoardFromArgs(ruleset rules.Ruleset, snakeStates map[string]SnakeState) *rules.BoardState {
	if Timeout == 0 {
		Timeout = 500
	}
	HttpClient = http.Client{
		Timeout: time.Duration(Timeout) * time.Millisecond,
	}

	snakeIds := []string{}
	for _, snakeState := range snakeStates {
		snakeIds = append(snakeIds, snakeState.ID)
	}
	state, err := rules.CreateDefaultBoardState(Width, Height, snakeIds)
	if err != nil {
		log.Panic("[PANIC]: Error Initializing Board State")
	}
	state, err = ruleset.ModifyInitialBoardState(state)
	if err != nil {
		log.Panic("[PANIC]: Error Initializing Board State")
	}

	for _, snakeState := range snakeStates {
		snakeRequest := getIndividualBoardStateForSnake(state, snakeState, snakeStates, ruleset)
		requestBody := serialiseSnakeRequest(snakeRequest)
		u, _ := url.ParseRequestURI(snakeState.URL)
		u.Path = path.Join(u.Path, "start")
		if DebugRequests {
			log.Printf("POST %s: %v", u, string(requestBody))
		}
		_, err = HttpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("[WARN]: Request to %v failed", u.String())
		}
	}
	return state
}

func createNextBoardState(ruleset rules.Ruleset, state *rules.BoardState, snakeStates map[string]SnakeState, turn int32) *rules.BoardState {
	var moves []rules.SnakeMove
	if Sequential {
		for _, snakeState := range snakeStates {
			for _, snake := range state.Snakes {
				if snakeState.ID == snake.ID && snake.EliminatedCause == rules.NotEliminated {
					moves = append(moves, getMoveForSnake(ruleset, state, snakeState, snakeStates))
				}
			}
		}
	} else {
		var wg sync.WaitGroup
		c := make(chan rules.SnakeMove, len(snakeStates))

		for _, snakeState := range snakeStates {
			for _, snake := range state.Snakes {
				if snakeState.ID == snake.ID && snake.EliminatedCause == rules.NotEliminated {
					wg.Add(1)
					go func(snakeState SnakeState) {
						defer wg.Done()
						c <- getMoveForSnake(ruleset, state, snakeState, snakeStates)
					}(snakeState)
				}
			}
		}

		wg.Wait()
		close(c)

		for move := range c {
			moves = append(moves, move)
		}
	}
	for _, move := range moves {
		snakeState := snakeStates[move.ID]
		snakeState.LastMove = move.Move
		snakeStates[move.ID] = snakeState
	}
	state, err := ruleset.CreateNextBoardState(state, moves)
	if err != nil {
		log.Panicf("[PANIC]: Error Producing Next Board State: %v", err)
	}

	state.Turn = turn

	return state
}

func getMoveForSnake(ruleset rules.Ruleset, state *rules.BoardState, snakeState SnakeState, snakeStates map[string]SnakeState) rules.SnakeMove {
	snakeRequest := getIndividualBoardStateForSnake(state, snakeState, snakeStates, ruleset)
	requestBody := serialiseSnakeRequest(snakeRequest)
	u, _ := url.ParseRequestURI(snakeState.URL)
	u.Path = path.Join(u.Path, "move")
	if DebugRequests {
		log.Printf("POST %s: %v", u, string(requestBody))
	}
	res, err := HttpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	move := snakeState.LastMove
	if err != nil {
		log.Printf("[WARN]: Request to %v failed\n", u.String())
		log.Printf("Body --> %v\n", string(requestBody))
	} else if res.Body != nil {
		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		} else {
			playerResponse := client.MoveResponse{}
			jsonErr := json.Unmarshal(body, &playerResponse)
			if jsonErr != nil {
				log.Fatal(jsonErr)
			} else {
				move = playerResponse.Move
			}
		}
	}
	return rules.SnakeMove{ID: snakeState.ID, Move: move}
}

func sendEndRequest(ruleset rules.Ruleset, state *rules.BoardState, snakeState SnakeState, snakeStates map[string]SnakeState) {
	snakeRequest := getIndividualBoardStateForSnake(state, snakeState, snakeStates, ruleset)
	requestBody := serialiseSnakeRequest(snakeRequest)
	u, _ := url.ParseRequestURI(snakeState.URL)
	u.Path = path.Join(u.Path, "end")
	if DebugRequests {
		log.Printf("POST %s: %v", u, string(requestBody))
	}
	_, err := HttpClient.Post(u.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("[WARN]: Request to %v failed", u.String())
	}
}

func getIndividualBoardStateForSnake(state *rules.BoardState, snakeState SnakeState, snakeStates map[string]SnakeState, ruleset rules.Ruleset) client.SnakeRequest {
	var youSnake rules.Snake
	for _, snk := range state.Snakes {
		if snakeState.ID == snk.ID {
			youSnake = snk
			break
		}
	}
	request := client.SnakeRequest{
		Game:  createClientGame(ruleset),
		Turn:  Turn,
		Board: convertStateToBoard(state, snakeStates),
		You:   convertRulesSnake(youSnake, snakeStates[youSnake.ID]),
	}
	return request
}

func serialiseSnakeRequest(snakeRequest client.SnakeRequest) []byte {
	requestJSON, err := json.Marshal(snakeRequest)
	if err != nil {
		log.Panic("[PANIC]: Error Marshalling JSON from State")
		panic(err)
	}
	return requestJSON
}

func createClientGame(ruleset rules.Ruleset) client.Game {
	return client.Game{ID: GameId, Timeout: Timeout, Ruleset: client.Ruleset{
		Name:    ruleset.Name(),
		Version: "cli", // TODO: Use GitHub Release Version
		Settings: client.RulesetSettings{
			HazardDamagePerTurn: HazardDamagePerTurn,
			FoodSpawnChance:     FoodSpawnChance,
			MinimumFood:         MinimumFood,
			RoyaleSettings: client.RoyaleSettings{
				ShrinkEveryNTurns: ShrinkEveryNTurns,
			},
			SquadSettings: client.SquadSettings{
				AllowBodyCollisions: true,
				SharedElimination:   true,
				SharedHealth:        true,
				SharedLength:        true,
			},
		},
	}}
}

func convertRulesSnake(snake rules.Snake, snakeState SnakeState) client.Snake {
	return client.Snake{
		ID:      snake.ID,
		Name:    snakeState.Name,
		Health:  snake.Health,
		Body:    client.CoordFromPointArray(snake.Body),
		Latency: "0",
		Head:    client.CoordFromPoint(snake.Body[0]),
		Length:  int32(len(snake.Body)),
		Shout:   "",
		Squad:   snakeState.Squad,
		Customizations: client.Customizations{
			Head:  snakeState.Head,
			Tail:  snakeState.Tail,
			Color: snakeState.Color,
		},
	}
}

func convertRulesSnakes(snakes []rules.Snake, snakeStates map[string]SnakeState) []client.Snake {
	var a []client.Snake
	for _, snake := range snakes {
		if snake.EliminatedCause == rules.NotEliminated {
			a = append(a, convertRulesSnake(snake, snakeStates[snake.ID]))
		}
	}
	return a
}

func convertStateToBoard(state *rules.BoardState, snakeStates map[string]SnakeState) client.Board {
	return client.Board{
		Height:  state.Height,
		Width:   state.Width,
		Food:    client.CoordFromPointArray(state.Food),
		Hazards: client.CoordFromPointArray(state.Hazards),
		Snakes:  convertRulesSnakes(state.Snakes, snakeStates),
	}
}

func buildSnakesFromOptions() map[string]SnakeState {
	bodyChars := []rune{'■', '⌀', '●', '☻', '◘', '☺', '□', '⍟'}
	var numSnakes int
	snakes := map[string]SnakeState{}
	numNames := len(Names)
	numURLs := len(URLs)
	numSquads := len(Squads)
	if numNames > numURLs {
		numSnakes = numNames
	} else {
		numSnakes = numURLs
	}
	if numNames != numURLs {
		log.Println("[WARN]: Number of Names and URLs do not match: defaults will be applied to missing values")
	}
	for i := int(0); i < numSnakes; i++ {
		var snakeName string
		var snakeURL string
		var snakeSquad string

		id := uuid.New().String()

		if i < numNames {
			snakeName = Names[i]
		} else {
			log.Printf("[WARN]: Name for URL %v is missing: a default name will be applied\n", URLs[i])
			snakeName = id
		}

		if i < numURLs {
			u, err := url.ParseRequestURI(URLs[i])
			if err != nil {
				log.Printf("[WARN]: URL %v is not valid: a default will be applied\n", URLs[i])
				snakeURL = "https://example.com"
			} else {
				snakeURL = u.String()
			}
		} else {
			log.Printf("[WARN]: URL for Name %v is missing: a default URL will be applied\n", Names[i])
			snakeURL = "https://example.com"
		}

		if GameType == "squad" {
			if i < numSquads {
				snakeSquad = Squads[i]
			} else {
				log.Printf("[WARN]: Squad for URL %v is missing: a default squad will be applied\n", URLs[i])
				snakeSquad = strconv.Itoa(i / 2)
			}
		}
		snakeState := SnakeState{
			Name: snakeName, URL: snakeURL, ID: id, LastMove: "up", Character: bodyChars[i%8],
		}
		res, err := HttpClient.Get(snakeURL)
		if err != nil {
			log.Printf("[WARN]: Request to %v failed: %v", snakeURL, err)
		} else if res.Body != nil {
			defer res.Body.Close()
			body, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			pingResponse := client.SnakeMetadataResponse{}
			jsonErr := json.Unmarshal(body, &pingResponse)
			if jsonErr != nil {
				log.Printf("Error reading response from %v: %v", snakeURL, jsonErr)
			} else {
				snakeState.Head = pingResponse.Head
				snakeState.Tail = pingResponse.Tail
				snakeState.Color = pingResponse.Color
			}
		}
		if GameType == "squad" {
			snakeState.Squad = snakeSquad
		}
		snakes[snakeState.ID] = snakeState
	}
	return snakes
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

func printMap(state *rules.BoardState, snakeStates map[string]SnakeState, gameTurn int32) {
	var o bytes.Buffer
	o.WriteString(fmt.Sprintf("Ruleset: %s, Seed: %d, Turn: %v\n", GameType, Seed, gameTurn))
	board := make([][]string, state.Width)
	for i := range board {
		board[i] = make([]string, state.Height)
	}
	for y := int32(0); y < state.Height; y++ {
		for x := int32(0); x < state.Width; x++ {
			if UseColor {
				board[x][y] = TERM_FG_LIGHTGRAY + "□"
			} else {
				board[x][y] = "◦"
			}
		}
	}
	for _, oob := range state.Hazards {
		if UseColor {
			board[oob.X][oob.Y] = TERM_BG_GRAY + " " + TERM_BG_WHITE
		} else {
			board[oob.X][oob.Y] = "░"
		}
	}
	if UseColor {
		o.WriteString(fmt.Sprintf("Hazards "+TERM_BG_GRAY+" "+TERM_RESET+": %v\n", state.Hazards))
	} else {
		o.WriteString(fmt.Sprintf("Hazards ░: %v\n", state.Hazards))
	}
	for _, f := range state.Food {
		if UseColor {
			board[f.X][f.Y] = TERM_FG_FOOD + "●"
		} else {
			board[f.X][f.Y] = "⚕"
		}
	}
	if UseColor {
		o.WriteString(fmt.Sprintf("Food "+TERM_FG_FOOD+TERM_BG_WHITE+"●"+TERM_RESET+": %v\n", state.Food))
	} else {
		o.WriteString(fmt.Sprintf("Food ⚕: %v\n", state.Food))
	}
	for _, s := range state.Snakes {
		red, green, blue := parseSnakeColor(snakeStates[s.ID].Color)
		for _, b := range s.Body {
			if b.X >= 0 && b.X < state.Width && b.Y >= 0 && b.Y < state.Height {
				if UseColor {
					board[b.X][b.Y] = fmt.Sprintf(TERM_FG_RGB+"■", red, green, blue)
				} else {
					board[b.X][b.Y] = string(snakeStates[s.ID].Character)
				}
			}
		}
		if UseColor {
			o.WriteString(fmt.Sprintf("%v "+TERM_FG_RGB+TERM_BG_WHITE+"■■■"+TERM_RESET+": %v\n", snakeStates[s.ID].Name, red, green, blue, s))
		} else {
			o.WriteString(fmt.Sprintf("%v %c: %v\n", snakeStates[s.ID].Name, snakeStates[s.ID].Character, s))
		}
	}
	for y := state.Height - 1; y >= 0; y-- {
		if UseColor {
			o.WriteString(TERM_BG_WHITE)
		}
		for x := int32(0); x < state.Width; x++ {
			o.WriteString(board[x][y])
		}
		if UseColor {
			o.WriteString(TERM_RESET)
		}
		o.WriteString("\n")
	}
	log.Print(o.String())
}
