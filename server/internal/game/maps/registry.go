package maps

import (
	"fmt"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Registry holds all available maps
var Registry = make(map[string]*types.MapDefinition)

// DefaultMapID is the ID of the default map
const DefaultMapID = "classic"

func init() {
	// Register all maps
	Register(ClassicMap())
	Register(TheDivideMap())
}

// Register adds a map to the registry
func Register(m *types.MapDefinition) {
	Registry[m.ID] = m
}

// Get retrieves a map by ID
func Get(id string) (*types.MapDefinition, error) {
	m, ok := Registry[id]
	if !ok {
		return nil, fmt.Errorf("map not found: %s", id)
	}
	return m, nil
}

// GetDefault retrieves the default map
func GetDefault() *types.MapDefinition {
	m, _ := Get(DefaultMapID)
	return m
}

// List returns all available map IDs
func List() []string {
	ids := make([]string, 0, len(Registry))
	for id := range Registry {
		ids = append(ids, id)
	}
	return ids
}
