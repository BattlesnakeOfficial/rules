package maps

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

// MapRegistry is a mapping of map names to map generators
type MapRegistry map[string]Generator

var globalRegistry = MapRegistry{}

// RegisterMap adds a stage to the registry.
// If a map has already been registered this will panic.
func (registry MapRegistry) RegisterMap(id string, m Generator) {
	if _, ok := registry[id]; ok {
		panic(fmt.Sprintf("map '%s' has already been registered", id))
	}

	registry[id] = m
}

func (registry MapRegistry) GetMap(id string) (Generator, error) {
	if m, ok := registry[id]; ok {
		return m, nil
	}
	return nil, rules.ErrorMapNotFound
}

func GetMap(id string) (Generator, error) {
	return globalRegistry.GetMap(id)
}
