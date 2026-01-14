package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// BuyZone represents a location where players can purchase units
type BuyZone struct {
	ID       string
	OwnerID  int
	UnitType string // "tank" or "airplane"
	Position types.Vector3
	Radius   float64
	Cost     int
}

// ToType converts BuyZone to types.BuyZone for JSON serialization
func (b *BuyZone) ToType() types.BuyZone {
	return types.BuyZone{
		ID:       b.ID,
		OwnerID:  b.OwnerID,
		UnitType: b.UnitType,
		Position: b.Position,
		Radius:   b.Radius,
		Cost:     b.Cost,
	}
}

// IsPlayerInRange checks if a player position is within the buy zone
func (b *BuyZone) IsPlayerInRange(pos types.Vector3) bool {
	dx := pos.X - b.Position.X
	dz := pos.Z - b.Position.Z
	distSq := dx*dx + dz*dz
	return distSq <= b.Radius*b.Radius
}

// GetBuyZones returns the buy zones for both players
// Each player has a tank buy zone and a helicopter buy zone in their base
func GetBuyZones() []*BuyZone {
	zones := make([]*BuyZone, 0)

	// Player 1 buy zones (base at X=-90)
	// Tank buy zone (north side of base)
	zones = append(zones, &BuyZone{
		ID:       "buy_p1_tank",
		OwnerID:  0,
		UnitType: "tank",
		Position: types.Vector3{X: -90, Y: 0, Z: -8},
		Radius:   4.0,
		Cost:     types.TankCost,
	})
	// Helicopter buy zone (south side of base)
	zones = append(zones, &BuyZone{
		ID:       "buy_p1_airplane",
		OwnerID:  0,
		UnitType: "airplane",
		Position: types.Vector3{X: -90, Y: 0, Z: 8},
		Radius:   4.0,
		Cost:     types.AirplaneCost,
	})

	// Player 2 buy zones (base at X=90)
	// Tank buy zone (north side of base)
	zones = append(zones, &BuyZone{
		ID:       "buy_p2_tank",
		OwnerID:  1,
		UnitType: "tank",
		Position: types.Vector3{X: 90, Y: 0, Z: -8},
		Radius:   4.0,
		Cost:     types.TankCost,
	})
	// Helicopter buy zone (south side of base)
	zones = append(zones, &BuyZone{
		ID:       "buy_p2_airplane",
		OwnerID:  1,
		UnitType: "airplane",
		Position: types.Vector3{X: 90, Y: 0, Z: 8},
		Radius:   4.0,
		Cost:     types.AirplaneCost,
	})

	return zones
}
