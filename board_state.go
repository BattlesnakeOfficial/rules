package rules

func (b *BoardState) Copy() *BoardState {
	nextState := &BoardState{
		Height: b.Height,
		Width:  b.Width,
		Food:   append([]Point{}, b.Food...),
		Snakes: make([]Snake, len(b.Snakes)),
	}
	for i := 0; i < len(b.Snakes); i++ {
		nextState.Snakes[i].ID = b.Snakes[i].ID
		nextState.Snakes[i].Health = b.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, b.Snakes[i].Body...)
		nextState.Snakes[i].EliminatedCause = b.Snakes[i].EliminatedCause
		nextState.Snakes[i].EliminatedBy = b.Snakes[i].EliminatedBy
	}
	return nextState
}
