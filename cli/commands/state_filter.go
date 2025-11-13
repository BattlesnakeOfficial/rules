package commands

import (
	"fmt"
	"strconv"

	"github.com/BattlesnakeOfficial/rules"
)

func manhattan_d(a, b rules.Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func FilterBoardStateForSnake(boardState *rules.BoardState, self SnakeState, viewRadius int) *rules.BoardState {
	filtered := &rules.BoardState{
		Turn:      boardState.Turn,
		Height:    boardState.Height,
		Width:     boardState.Width,
		GameState: boardState.GameState,
		Food:      []rules.Point{},
		Hazards:   []rules.Point{},
		Snakes:    []rules.Snake{},
	}

	// find the head of self snake
	var head rules.Point
	for _, s := range boardState.Snakes {
		if s.ID == self.ID && len(s.Body) > 0 {
			head = s.Body[0]
			break
		}
	}

	// FILTER FOOD on view radius or spawn turn
	for _, f := range boardState.Food {
		visible := manhattan_d(f, head) <= viewRadius
		key := fmt.Sprintf("food_spawn_%d_%d", f.X, f.Y)
		spawnTurnStr, spawned := boardState.GameState[key]

		spawnVisible := false
		if spawned {
			spawnTurn, _ := strconv.Atoi(spawnTurnStr)
			spawnVisible = spawnTurn == boardState.Turn
		}

		if visible || spawnVisible {
			filtered.Food = append(filtered.Food, f)
		}
	}

	// FILTER HAZARDS with view radius
	for _, h := range boardState.Hazards {
		if manhattan_d(h, head) <= viewRadius {
			filtered.Hazards = append(filtered.Hazards, h)
		}
	}

	// FILTER SNAKE bodies
	for _, s := range boardState.Snakes {
		if s.ID == self.ID {
			// self snake sees whole body
			filtered.Snakes = append(filtered.Snakes, rules.Snake{
				ID:     s.ID,
				Body:   append([]rules.Point(nil), s.Body...),
				Health: s.Health,
			})
			continue
		}

		filteredBody := []rules.Point{}
		for i, seg := range s.Body {
			if manhattan_d(seg, head) <= viewRadius {
				filteredBody = append(filteredBody, seg)
			} else if i == 0 || (i > 0 && manhattan_d(s.Body[i-1], head) <= viewRadius) {
				// mark end with -1
				filteredBody = append(filteredBody, rules.Point{X: -1, Y: -1})
			}
		}

		if len(filteredBody) == 0 {
			// mark end with -1
			filteredBody = append(filteredBody, rules.Point{X: -1, Y: -1})
		}

		filtered.Snakes = append(filtered.Snakes, rules.Snake{
			ID:     s.ID,
			Body:   filteredBody,
			Health: 0,
		})
	}

	return filtered
}
