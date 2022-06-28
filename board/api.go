package board

import (
	"github.com/BattlesnakeOfficial/rules"
)

// Types used to implement the JSON API expected by the board client.

// JSON structure returned by the game status endpoint.
type Game struct {
	ID           string            `json:"ID"`
	Status       string            `json:"Status"`
	Width        int               `json:"Width"`
	Height       int               `json:"Height"`
	Ruleset      map[string]string `json:"Ruleset"`
	SnakeTimeout int               `json:"SnakeTimeout"`
	Source       string            `json:"Source"`
	RulesetName  string            `json:"RulesetName"`
	RulesStages  []string          `json:"RulesStages"`
	Map          string            `json:"Map"`
}

// The websocket stream has support for returning different types of events, along with a "type" attribute.
type GameEventType string

const (
	EVENT_TYPE_FRAME    GameEventType = "frame"
	EVENT_TYPE_GAME_END GameEventType = "game_end"
)

// Top-level JSON structure sent in each websocket frame.
type GameEvent struct {
	EventType GameEventType `json:"Type"`
	Data      interface{}   `json:"Data"`
}

// Represents a single turn in the game.
type GameFrame struct {
	Turn    int           `json:"Turn"`
	Snakes  []Snake       `json:"Snakes"`
	Food    []rules.Point `json:"Food"`
	Hazards []rules.Point `json:"Hazards"`
}

type GameEnd struct {
	Game Game `json:"game"`
}

type Snake struct {
	ID            string        `json:"ID"`
	Name          string        `json:"Name"`
	Body          []rules.Point `json:"Body"`
	Health        int           `json:"Health"`
	Death         *Death        `json:"Death"`
	Color         string        `json:"Color"`
	HeadType      string        `json:"HeadType"`
	TailType      string        `json:"TailType"`
	Latency       string        `json:"Latency"`
	Shout         string        `json:"Shout"`
	Squad         string        `json:"Squad"`
	Author        string        `json:"Author"`
	StatusCode    int           `json:"StatusCode"`
	Error         string        `json:"Error"`
	IsBot         bool          `json:"IsBot"`
	IsEnvironment bool          `json:"IsEnvironment"`
}

type Death struct {
	Cause        string `json:"Cause"`
	Turn         int    `json:"Turn"`
	EliminatedBy string `json:"EliminatedBy"`
}
