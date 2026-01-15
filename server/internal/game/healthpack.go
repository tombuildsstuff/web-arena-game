package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// HealthPack represents a collectible health pickup
type HealthPack struct {
	ID        string
	Position  types.Vector3
	HealAmount int
	Radius    float64 // Collection radius
	SpawnedAt int64   // When this pack was spawned
}

// NewHealthPack creates a new health pack at the given position
func NewHealthPack(id string, position types.Vector3) *HealthPack {
	return &HealthPack{
		ID:         id,
		Position:   position,
		HealAmount: types.HealthPackHealAmount,
		Radius:     types.HealthPackRadius,
		SpawnedAt:  time.Now().UnixMilli(),
	}
}

// ToType converts HealthPack to types.HealthPack for JSON serialization
func (h *HealthPack) ToType() types.HealthPack {
	return types.HealthPack{
		ID:       h.ID,
		Position: h.Position,
	}
}

// IsPlayerInRange checks if a player position is within collection range
func (h *HealthPack) IsPlayerInRange(pos types.Vector3) bool {
	dx := pos.X - h.Position.X
	dz := pos.Z - h.Position.Z
	distSq := dx*dx + dz*dz
	return distSq <= h.Radius*h.Radius
}

// HealthPackSystem manages health pack spawning and collection
type HealthPackSystem struct {
	lastSpawnTime  int64
	nextSpawnDelay int64 // milliseconds until next spawn
	idCounter      int
}

// NewHealthPackSystem creates a new health pack system
func NewHealthPackSystem() *HealthPackSystem {
	return &HealthPackSystem{
		lastSpawnTime:  time.Now().UnixMilli(),
		nextSpawnDelay: randomSpawnDelay(),
		idCounter:      0,
	}
}

// randomSpawnDelay returns a random delay between spawns (15-45 seconds)
func randomSpawnDelay() int64 {
	minDelay := int64(types.HealthPackSpawnMinSeconds * 1000)
	maxDelay := int64(types.HealthPackSpawnMaxSeconds * 1000)
	return minDelay + rand.Int63n(maxDelay-minDelay)
}

// Update checks for spawning new health packs and handles collection
func (s *HealthPackSystem) Update(state *State) {
	now := time.Now().UnixMilli()

	// Check if it's time to spawn a new health pack
	if now-s.lastSpawnTime >= s.nextSpawnDelay {
		// Only spawn if we don't have too many health packs
		if len(state.HealthPacks) < types.HealthPackMaxCount {
			s.spawnHealthPack(state)
		}
		s.lastSpawnTime = now
		s.nextSpawnDelay = randomSpawnDelay()
	}

	// Check for collection by player units
	s.checkCollection(state)

	// Remove old health packs that haven't been collected
	s.removeExpiredPacks(state, now)
}

// spawnHealthPack spawns a health pack at a random valid location
func (s *HealthPackSystem) spawnHealthPack(state *State) {
	// Try to find a valid spawn position
	for attempts := 0; attempts < 20; attempts++ {
		// Random position within the arena (avoiding edges and bases)
		x := (rand.Float64() - 0.5) * 140 // -70 to 70
		z := (rand.Float64() - 0.5) * 160 // -80 to 80

		pos := types.Vector3{X: x, Y: 1.0, Z: z}

		// Check if position is walkable (not inside an obstacle)
		// Simple check: avoid being too close to bases
		if x < -70 || x > 70 {
			continue // Too close to bases
		}

		// Get elevation at this position
		pos.Y = state.GetElevationAt(x, z) + 1.0

		// Create the health pack
		s.idCounter++
		pack := NewHealthPack(fmt.Sprintf("healthpack_%d", s.idCounter), pos)
		state.HealthPacks = append(state.HealthPacks, pack)
		return
	}
}

// checkCollection checks if any player unit is collecting a health pack
func (s *HealthPackSystem) checkCollection(state *State) {
	packsToRemove := make([]string, 0)

	for _, pack := range state.HealthPacks {
		for _, unit := range state.Units {
			// Only player units can collect health packs
			playerUnit, ok := unit.(*PlayerUnit)
			if !ok {
				continue
			}

			// Skip dead/respawning players
			if !playerUnit.IsAlive() {
				continue
			}

			// Check if player is in range
			if pack.IsPlayerInRange(playerUnit.GetPosition()) {
				// Heal the player (up to max health)
				currentHealth := playerUnit.GetHealth()
				maxHealth := types.PlayerUnitHealth
				newHealth := currentHealth + pack.HealAmount
				if newHealth > maxHealth {
					newHealth = maxHealth
				}
				playerUnit.SetHealth(newHealth)

				// Mark pack for removal
				packsToRemove = append(packsToRemove, pack.ID)
				break // Pack can only be collected once
			}
		}
	}

	// Remove collected packs
	for _, id := range packsToRemove {
		state.RemoveHealthPack(id)
	}
}

// removeExpiredPacks removes health packs that have been around too long
func (s *HealthPackSystem) removeExpiredPacks(state *State, now int64) {
	maxAge := int64(types.HealthPackLifetimeSeconds * 1000)
	packsToRemove := make([]string, 0)

	for _, pack := range state.HealthPacks {
		if now-pack.SpawnedAt > maxAge {
			packsToRemove = append(packsToRemove, pack.ID)
		}
	}

	for _, id := range packsToRemove {
		state.RemoveHealthPack(id)
	}
}
