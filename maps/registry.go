package maps

import (
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

// MapRegistry is a mapping of map names to game maps.
type MapRegistry map[string]GameMap

var globalRegistry = MapRegistry{}

// RegisterMap adds a stage to the registry.
// If a map has already been registered this will panic.
func (registry MapRegistry) RegisterMap(id string, m GameMap) {
	if _, ok := registry[id]; ok {
		panic(fmt.Sprintf("map '%s' has already been registered", id))
	}

	registry[id] = m
}

// GetMap returns the map associated with the given ID.
func (registry MapRegistry) GetMap(id string) (GameMap, error) {
	if m, ok := registry[id]; ok {
		return m, nil
	}
	return nil, rules.ErrorMapNotFound
}

// GetMap returns the map associated with the given ID from the global registry.
func GetMap(id string) (GameMap, error) {
	return globalRegistry.GetMap(id)
}

// RegisterMap adds a map to the global registry.
func RegisterMap(id string, m GameMap) {
	globalRegistry.RegisterMap(id, m)
}
