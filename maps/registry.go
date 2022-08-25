package maps

import (
	"fmt"
	"sort"

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

// List returns all registered map IDs in alphabetical order
func (registry MapRegistry) List() []string {
	var keys []string
	for k := range registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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

// List returns a list of maps registered to the global registry.
func List() []string {
	return globalRegistry.List()
}

// RegisterMap adds a map to the global registry.
func RegisterMap(id string, m GameMap) {
	globalRegistry.RegisterMap(id, m)
}

func TestMap(id string, m GameMap, callback func()) {
	globalRegistry[id] = m
	callback()
	delete(globalRegistry, id)
}
