package client

import "github.com/BattlesnakeOfficial/rules"

// The top-level message sent in /start, /move, and /end requests
type SnakeRequest struct {
	Game  Game  `json:"game"`
	Turn  int32 `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

// Game represents the current game state
type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
	Source  string  `json:"source"`
}

// Board provides information about the game board
type Board struct {
	Height      int32        `json:"height"`
	Width       int32        `json:"width"`
	Snakes      []Snake      `json:"snakes"`
	Food        []Coord      `json:"food"`
	Hazards     []Coord      `json:"hazards"`
	HazardSides *HazardSides `json:"hazardSides,omitempty"`
}

// What sides the hazards have spawned in from
type HazardSides struct {
	MinX int32 `json:"minX"`
	MaxX int32 `json:"maxX"`
	MinY int32 `json:"minY"`
	MaxY int32 `json:"maxY"`
}

// Snake represents information about a snake in the game
type Snake struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Latency        string         `json:"latency"`
	Health         int32          `json:"health"`
	Body           []Coord        `json:"body"`
	Head           Coord          `json:"head"`
	Length         int32          `json:"length"`
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

type RulesetSettings struct {
	FoodSpawnChance     int32          `json:"foodSpawnChance"`
	MinimumFood         int32          `json:"minimumFood"`
	HazardDamagePerTurn int32          `json:"hazardDamagePerTurn"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"`
}

type RoyaleSettings struct {
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

// Coord represents a point on the board
type Coord struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
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
