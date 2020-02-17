package rules

type TeamRuleset struct {
	StandardRuleset

	BodyPassthrough bool
	SharedStats     bool
	SharedDeath     bool

	teams map[string][]string
}

type teamInfo struct {
	maxSnakeLength int
	hasDeadSnake   bool
	firstDeadSnake string
	snakes         []*Snake
}

const EliminatedByTeamMemberDied = "team-member-died"

func (r *TeamRuleset) ResolveMoves(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	// We specifically want to copy prevState, so as not to alter it directly.
	nextState := &BoardState{
		Height: prevState.Height,
		Width:  prevState.Width,
		Food:   append([]Point{}, prevState.Food...),
		Snakes: make([]Snake, len(prevState.Snakes)),
	}
	for i := 0; i < len(prevState.Snakes); i++ {
		nextState.Snakes[i].ID = prevState.Snakes[i].ID
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
	}

	// TODO: Gut check the BoardState?

	// TODO: LOG?
	err := r.moveSnakes(nextState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.reduceSnakeHealth(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.eliminateSnakes(nextState)
	if err != nil {
		return nil, err
	}

	// TODO
	// bvanvugt: we specifically want this to happen before elimination
	// so that head-to-head collisions on food still remove the food.
	// It does create an artifact though, where head-to-head collisions
	// of equal length actually show length + 1

	// TODO: LOG?
	err = r.feedSnakes(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.maybeSpawnFood(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *TeamRuleset) feedSnakes(b *BoardState) error {
	err := r.StandardRuleset.feedSnakes(b)
	if err != nil {
		return err
	}

	if !r.SharedStats {
		return nil
	}

	teams := r.groupSnakesByTeam(b)

	for _, team := range teams {
		for _, snake := range team.snakes {
			if len(snake.Body) < team.maxSnakeLength {
				r.feedSnake(snake)
			}
		}
	}

	return nil
}

func (r *TeamRuleset) groupSnakesByTeam(b *BoardState) map[string]teamInfo {
	teams := map[string]teamInfo{}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		teamID := r.findTeam(snake.ID)
		if teamID == "" {
			continue
		}
		var team teamInfo
		var ok bool
		if team, ok = teams[teamID]; !ok {
			teams[teamID] = teamInfo{
				maxSnakeLength: 0,
				snakes:         []*Snake{},
				hasDeadSnake:   false,
			}
		}
		team.snakes = append(team.snakes, snake)
		if len(snake.Body) > team.maxSnakeLength {
			team.maxSnakeLength = len(snake.Body)
		}
		if snake.EliminatedCause != NotEliminated && team.firstDeadSnake == "" {
			team.hasDeadSnake = true
			team.firstDeadSnake = snake.ID
		}
		teams[teamID] = team
	}
	return teams
}

func (r *TeamRuleset) eliminateSnakes(b *BoardState) error {
	err := r.StandardRuleset.eliminateSnakes(b)
	if err != nil {
		return err
	}

	r.handleBodyPassthrough(b)
	r.handleSharedDeath(b)

	return nil
}

func (r *TeamRuleset) handleSharedDeath(b *BoardState) {
	if !r.SharedDeath {
		return
	}

	teams := r.groupSnakesByTeam(b)
	for _, team := range teams {
		if !team.hasDeadSnake {
			continue
		}
		for _, snake := range team.snakes {
			if snake.EliminatedCause == NotEliminated {
				snake.EliminatedCause = EliminatedByTeamMemberDied
				snake.EliminatedBy = team.firstDeadSnake
			}
		}
	}
}

func (r *TeamRuleset) handleBodyPassthrough(b *BoardState) {
	if !r.BodyPassthrough {
		return
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != EliminatedByCollision {
			continue
		}
		team1 := r.findTeam(snake.ID)
		team2 := r.findTeam(snake.EliminatedBy)
		if team1 == "" || team2 == "" {
			continue
		}

		if team1 == team2 {
			snake.EliminatedCause = NotEliminated
			snake.EliminatedBy = ""
		}
	}
}

func (r *TeamRuleset) findTeam(id string) string {
	for team, snakes := range r.teams {
		for _, snakeID := range snakes {
			if snakeID == id {
				return team
			}
		}
	}
	return ""
}
