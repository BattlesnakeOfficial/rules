package rules

type Rand interface {
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
}
