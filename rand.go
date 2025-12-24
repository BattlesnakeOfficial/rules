package rules

import (
	"math/rand"
	"sync"
)

type Rand interface {
	Intn(n int) int
	// Range produces a random integer in the range of [min,max] (inclusive)
	// For example, Range(1,3) could produce the values 1, 2 or 3.
	// Panics if max < min (like how Intn(n) panics for n <=0)
	Range(min, max int) int
	Shuffle(n int, swap func(i, j int))
}

// A Rand implementation that just uses a shared math/rand.Rand guarded by a lock.
// This preserves the behaviour of the package-level math/rand global while allowing
// deterministic seeding without rand.Seed.
var GlobalRand Rand = globalRandInstance()

type lockedRand struct {
	mu sync.Mutex
	r  *rand.Rand
}

func globalRandInstance() *lockedRand {
	return &lockedRand{r: rand.New(rand.NewSource(1))}
}

func (g *lockedRand) Range(min, max int) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.r.Intn(max-min+1) + min
}

func (g *lockedRand) Intn(n int) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.r.Intn(n)
}

func (g *lockedRand) Shuffle(n int, swap func(i, j int)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.r.Shuffle(n, swap)
}

// SetGlobalSeed seeds the shared global RNG using rand.New(rand.NewSource(seed)),
// matching the previous rand.Seed semantics without triggering deprecation lints.
func SetGlobalSeed(seed int64) {
	gr := GlobalRand.(*lockedRand)
	gr.mu.Lock()
	gr.r = rand.New(rand.NewSource(seed))
	gr.mu.Unlock()
}

type seedRand struct {
	seed int64
	rand *rand.Rand
}

func NewSeedRand(seed int64) *seedRand {
	return &seedRand{
		seed: seed,
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (s seedRand) Intn(n int) int {
	return s.rand.Intn(n)
}

func (s seedRand) Range(min, max int) int {
	return s.rand.Intn(max-min+1) + min
}

func (s seedRand) Shuffle(n int, swap func(i, j int)) {
	s.rand.Shuffle(n, swap)
}

// For testing purposes

// A Rand implementation that always returns the minimum value for any method.
var MinRand minRand

type minRand struct{}

func (minRand) Intn(n int) int {
	return 0
}

func (minRand) Range(min, max int) int {
	return min
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

func (maxRand) Range(min, max int) int {
	return max
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
