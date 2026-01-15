package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Barracks represents a claimable spawn point for infantry units
type Barracks struct {
	ID          string
	Position    types.Vector3
	OwnerID     int     // -1 = neutral, 0 = player 1, 1 = player 2
	Health      int
	MaxHealth   int
	IsDestroyed bool
	RespawnTime float64 // Seconds remaining until respawn as neutral
	ClaimRadius float64

	// Occupant tracking for healing
	Occupants map[string]float64 // UnitID -> time spent inside (seconds)

	// Pending scatter - units that need to be scattered after destruction
	// Populated when barracks is destroyed, cleared by room.go after processing
	PendingScatter []string
}

// NewBarracks creates a new barracks
func NewBarracks(id string, position types.Vector3) *Barracks {
	return &Barracks{
		ID:          id,
		Position:    position,
		OwnerID:     -1, // Start neutral
		Health:      types.BarracksHealth,
		MaxHealth:   types.BarracksHealth,
		IsDestroyed: false,
		RespawnTime: 0,
		ClaimRadius: types.BarracksClaimRadius,
		Occupants:   make(map[string]float64),
	}
}

// ToType converts Barracks to types.Barracks for JSON serialization
func (b *Barracks) ToType() types.Barracks {
	return types.Barracks{
		ID:            b.ID,
		Position:      b.Position,
		OwnerID:       b.OwnerID,
		Health:        b.Health,
		MaxHealth:     b.MaxHealth,
		IsDestroyed:   b.IsDestroyed,
		RespawnTime:   b.RespawnTime,
		ClaimRadius:   b.ClaimRadius,
		OccupantCount: len(b.Occupants),
	}
}

// IsUnitInRange checks if a unit position is within claiming range
func (b *Barracks) IsUnitInRange(pos types.Vector3) bool {
	dx := pos.X - b.Position.X
	dz := pos.Z - b.Position.Z
	distSq := dx*dx + dz*dz
	return distSq <= b.ClaimRadius*b.ClaimRadius
}

// CanBeClaimed returns true if the barracks can be claimed by the given player
// Only infantry units can claim barracks
func (b *Barracks) CanBeClaimed(playerID int) bool {
	// Can't claim if destroyed
	if b.IsDestroyed {
		return false
	}
	// Can't claim your own barracks
	if b.OwnerID == playerID {
		return false
	}
	// Can claim neutral OR enemy barracks (infantry can recapture)
	return true
}

// Claim sets the barracks owner
func (b *Barracks) Claim(playerID int) {
	b.OwnerID = playerID
}

// TakeDamage applies damage to the barracks
func (b *Barracks) TakeDamage(amount int) {
	if b.IsDestroyed {
		return
	}
	b.Health -= amount
	if b.Health <= 0 {
		b.Health = 0
		b.IsDestroyed = true
		b.RespawnTime = types.BarracksRespawnTime
		b.OwnerID = -1 // Reset to neutral when destroyed

		// Store occupants for scatter - room.go will process this
		b.PendingScatter = b.GetOccupantIDs()
		b.ClearOccupants()
	}
}

// Update handles respawn timing
func (b *Barracks) Update(deltaTime float64) {
	if !b.IsDestroyed {
		return
	}

	b.RespawnTime -= deltaTime
	if b.RespawnTime <= 0 {
		b.Respawn()
	}
}

// Respawn resets the barracks to neutral state
func (b *Barracks) Respawn() {
	b.IsDestroyed = false
	b.Health = b.MaxHealth
	b.RespawnTime = 0
	b.OwnerID = -1 // Always respawn as neutral
}

// IsAlive returns true if the barracks is not destroyed
func (b *Barracks) IsAlive() bool {
	return !b.IsDestroyed
}

// UpdateOccupant tracks an infantry unit inside the barracks
// Returns the amount of health to heal (0 if not ready yet)
func (b *Barracks) UpdateOccupant(unitID string, deltaTime float64) int {
	if b.IsDestroyed {
		return 0
	}

	// Get or initialize time tracking
	timeInside, exists := b.Occupants[unitID]
	if !exists {
		b.Occupants[unitID] = 0
		return 0
	}

	// Update time
	timeInside += deltaTime
	b.Occupants[unitID] = timeInside

	// Check if enough time has passed for healing
	// Heal every BarracksHealTime seconds
	healTicks := int(timeInside / types.BarracksHealTime)
	previousTicks := int((timeInside - deltaTime) / types.BarracksHealTime)

	if healTicks > previousTicks {
		return types.BarracksHealAmount
	}

	return 0
}

// RemoveOccupant removes an infantry unit from tracking
func (b *Barracks) RemoveOccupant(unitID string) {
	delete(b.Occupants, unitID)
}

// GetOccupantIDs returns all unit IDs currently inside the barracks
func (b *Barracks) GetOccupantIDs() []string {
	ids := make([]string, 0, len(b.Occupants))
	for id := range b.Occupants {
		ids = append(ids, id)
	}
	return ids
}

// ClearOccupants removes all occupants (called when barracks is destroyed)
func (b *Barracks) ClearOccupants() {
	b.Occupants = make(map[string]float64)
}

// GetAndClearPendingScatter returns units that need to scatter and clears the list
func (b *Barracks) GetAndClearPendingScatter() []string {
	if len(b.PendingScatter) == 0 {
		return nil
	}
	scattered := b.PendingScatter
	b.PendingScatter = nil
	return scattered
}

// GetBarracksFromMap creates barracks from a map definition
func GetBarracksFromMap(mapDef *types.MapDefinition) []*Barracks {
	if mapDef.Barracks == nil {
		return make([]*Barracks, 0)
	}

	barracks := make([]*Barracks, 0, len(mapDef.Barracks))
	for _, b := range mapDef.Barracks {
		barracks = append(barracks, NewBarracks(b.ID, b.Position))
	}
	return barracks
}
