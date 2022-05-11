package rules

import "math/rand"

type Rand interface {
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
}

// A Rand implementation that just uses the global math/rand generator.
var GlobalRand globalRand

type globalRand struct{}

func (globalRand) Intn(n int) int {
	return rand.Intn(n)
}

func (globalRand) Shuffle(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}

// For testing purposes

// A Rand implementation that always returns the minimum value for any method.
var MinRand minRand

type minRand struct{}

func (minRand) Intn(n int) int {
	return 0
}

func (minRand) Shuffle(n int, swap func(i, j int)) {
	// no shuffling
}

// A Rand implementation that always returns the maximum value for any method.
var MaxRand maxRand

type maxRand struct{}

func (maxRand) Intn(n int) int {
	return n - 1
}

func (maxRand) Shuffle(n int, swap func(i, j int)) {
	// rotate by one element so every element is moved
	if n < 2 {
		return
	}
	for i := 0; i < n-2; i++ {
		swap(i, i+1)
	}
	swap(n-2, n-1)
}
