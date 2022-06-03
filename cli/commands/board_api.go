package commands

// TODO: These types should be moved to a separate package.

type GameEventType string

const (
	EVENT_TYPE_FRAME    GameEventType = "frame"
	EVENT_TYPE_GAME_END GameEventType = "game_end"
)

type GameEvent struct {
	EventType GameEventType `json:"Type"`
	Data      interface{}   `json:"Data"`
}

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

type GameFrame struct {
	Turn    int     `json:"Turn"`
	Snakes  []Snake `json:"Snakes"`
	Food    []Point `json:"Food"`
	Hazards []Point `json:"Hazards"`
}

type Snake struct {
	ID            string  `json:"ID"`
	Name          string  `json:"Name"`
	Body          []Point `json:"Body"`
	Health        int     `json:"Health"`
	Death         *Death  `json:"Death"`
	Color         string  `json:"Color"`
	HeadType      string  `json:"HeadType"`
	TailType      string  `json:"TailType"`
	Latency       string  `json:"Latency"`
	Shout         string  `json:"Shout"`
	Squad         string  `json:"Squad"`
	Author        string  `json:"Author"`
	StatusCode    int     `json:"StatusCode"`
	Error         string  `json:"Error"`
	IsBot         bool    `json:"IsBot"`
	IsEnvironment bool    `json:"IsEnvironment"`
}

type Point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type Death struct {
	Cause        string `json:"Cause"`
	Turn         int    `json:"Turn"`
	EliminatedBy string `json:"EliminatedBy"`
}
