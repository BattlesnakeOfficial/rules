package rules

func (b *BoardState) Copy() *BoardState {
	nextState := &BoardState{
		Height:  prevState.Height,
		Width:   prevState.Width,
		Food:    append([]Point{}, prevState.Food...),
		Snakes:  make([]Snake, len(prevState.Snakes)),
		Hazards: append([]Point{}, prevState.Hazards...),
	}
	for i := 0; i < len(prevState.Snakes); i++ {
		nextState.Snakes[i].ID = prevState.Snakes[i].ID
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
		nextState.Snakes[i].EliminatedCause = prevState.Snakes[i].EliminatedCause
		nextState.Snakes[i].EliminatedBy = prevState.Snakes[i].EliminatedBy
	}
	return nextState
}
