package client

import "github.com/BattlesnakeOfficial/rules"

// The top-level message sent in /start, /move, and /end requests
type SnakeRequest struct {
	Game  Game  `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

// Game represents the current game state
type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Map     string  `json:"map"`
	Timeout int     `json:"timeout"`
	Source  string  `json:"source"`
}

// Board provides information about the game board
type Board struct {
	Height  int     `json:"height"`
	Width   int     `json:"width"`
	Snakes  []Snake `json:"snakes"`
	Food    []Coord `json:"food"`
	Hazards []Coord `json:"hazards"`
}

// Snake represents information about a snake in the game
type Snake struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Latency        string         `json:"latency"`
	Health         int            `json:"health"`
	Body           []Coord        `json:"body"`
	Head           Coord          `json:"head"`
	Length         int            `json:"length"`
	Shout          string         `json:"shout"`
	Squad          string         `json:"squad"`
	Customizations Customizations `json:"customizations"`
}

type Customizations struct {
	Color string `json:"color"`
	Head  string `json:"head"`
	Tail  string `json:"tail"`
}

type Ruleset struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	Settings RulesetSettings `json:"settings"`
}

// RulesetSettings contains a static collection of a few settings that are exposed through the API.
type RulesetSettings struct {
	FoodSpawnChance     int            `json:"foodSpawnChance"`
	MinimumFood         int            `json:"minimumFood"`
	HazardDamagePerTurn int            `json:"hazardDamagePerTurn"`
	HazardMap           string         `json:"hazardMap"`       // Deprecated, replaced by Game.Map
	HazardMapAuthor     string         `json:"hazardMapAuthor"` // Deprecated, no planned replacement
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"` // Deprecated, provided with default fields for API compatibility
}

// RoyaleSettings contains settings that are specific to the "royale" game mode
type RoyaleSettings struct {
	ShrinkEveryNTurns int `json:"shrinkEveryNTurns"`
}

// SquadSettings contains settings that are specific to the "squad" game mode
type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

// Converts a rules.Settings (which can contain arbitrary settings) into the static RulesetSettings used in the client API.
func ConvertRulesetSettings(settings rules.Settings) RulesetSettings {
	return RulesetSettings{
		FoodSpawnChance:     settings.Int(rules.ParamFoodSpawnChance, 0),
		MinimumFood:         settings.Int(rules.ParamMinimumFood, 0),
		HazardDamagePerTurn: settings.Int(rules.ParamHazardDamagePerTurn, 0),
		RoyaleSettings: RoyaleSettings{
			ShrinkEveryNTurns: settings.Int(rules.ParamShrinkEveryNTurns, 0),
		},
		SquadSettings: SquadSettings{},
	}
}

// Coord represents a point on the board
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// The expected format of the response body from a /move request
type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

// The expected format of the response body from a GET request to a Battlesnake's index URL
type SnakeMetadataResponse struct {
	APIVersion string `json:"apiversion,omitempty"`
	Author     string `json:"author,omitempty"`
	Color      string `json:"color,omitempty"`
	Head       string `json:"head,omitempty"`
	Tail       string `json:"tail,omitempty"`
	Version    string `json:"version,omitempty"`
}

func CoordFromPoint(pt rules.Point) Coord {
	return Coord{X: pt.X, Y: pt.Y}
}

func CoordFromPointArray(ptArray []rules.Point) []Coord {
	a := make([]Coord, 0)
	for _, pt := range ptArray {
		a = append(a, CoordFromPoint(pt))
	}
	return a
}
