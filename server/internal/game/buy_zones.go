package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// BuyZone represents a location where players can purchase units
type BuyZone struct {
	ID          string
	OwnerID     int
	UnitType    string // "tank" or "airplane"
	Position    types.Vector3
	Radius      float64
	Cost        int  // Cost to buy units from this zone
	ClaimCost   int  // Cost to claim this zone (0 if not claimable)
	IsClaimable bool // Whether this zone can be claimed by players
}

// ToType converts BuyZone to types.BuyZone for JSON serialization
func (b *BuyZone) ToType() types.BuyZone {
	return types.BuyZone{
		ID:          b.ID,
		OwnerID:     b.OwnerID,
		UnitType:    b.UnitType,
		Position:    b.Position,
		Radius:      b.Radius,
		Cost:        b.Cost,
		ClaimCost:   b.ClaimCost,
		IsClaimable: b.IsClaimable,
	}
}

// CanBeClaimed returns true if this zone can be claimed by the given player
func (b *BuyZone) CanBeClaimed(playerID int) bool {
	if !b.IsClaimable {
		return false
	}
	// Can only claim neutral zones
	if b.OwnerID != -1 {
		return false
	}
	return true
}

// Claim sets the zone owner
func (b *BuyZone) Claim(playerID int) {
	b.OwnerID = playerID
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
// Plus claimable forward bases at top/bottom middle
func GetBuyZones() []*BuyZone {
	zones := make([]*BuyZone, 0)

	// Player 1 buy zones (base at X=-90)
	// Tank buy zone (north side of base)
	zones = append(zones, &BuyZone{
		ID:          "buy_p1_tank",
		OwnerID:     0,
		UnitType:    "tank",
		Position:    types.Vector3{X: -90, Y: 0, Z: -8},
		Radius:      4.0,
		Cost:        types.TankCost,
		IsClaimable: false,
	})
	// Helicopter buy zone (south side of base)
	zones = append(zones, &BuyZone{
		ID:          "buy_p1_airplane",
		OwnerID:     0,
		UnitType:    "airplane",
		Position:    types.Vector3{X: -90, Y: 0, Z: 8},
		Radius:      4.0,
		Cost:        types.AirplaneCost,
		IsClaimable: false,
	})

	// Player 2 buy zones (base at X=90)
	// Tank buy zone (north side of base)
	zones = append(zones, &BuyZone{
		ID:          "buy_p2_tank",
		OwnerID:     1,
		UnitType:    "tank",
		Position:    types.Vector3{X: 90, Y: 0, Z: -8},
		Radius:      4.0,
		Cost:        types.TankCost,
		IsClaimable: false,
	})
	// Helicopter buy zone (south side of base)
	zones = append(zones, &BuyZone{
		ID:          "buy_p2_airplane",
		OwnerID:     1,
		UnitType:    "airplane",
		Position:    types.Vector3{X: 90, Y: 0, Z: 8},
		Radius:      4.0,
		Cost:        types.AirplaneCost,
		IsClaimable: false,
	})

	// Claimable forward bases (neutral, on raised platforms at top/bottom)
	// North forward base (on raised platform at Y=3)
	zones = append(zones, &BuyZone{
		ID:          "buy_forward_north",
		OwnerID:     -1, // Neutral - can be claimed
		UnitType:    "tank",
		Position:    types.Vector3{X: 0, Y: 3, Z: -70},
		Radius:      5.0,
		Cost:        types.TankCost,
		ClaimCost:   types.ForwardBaseClaimCost,
		IsClaimable: true,
	})
	// South forward base (on raised platform at Y=3)
	zones = append(zones, &BuyZone{
		ID:          "buy_forward_south",
		OwnerID:     -1, // Neutral - can be claimed
		UnitType:    "tank",
		Position:    types.Vector3{X: 0, Y: 3, Z: 70},
		Radius:      5.0,
		Cost:        types.TankCost,
		ClaimCost:   types.ForwardBaseClaimCost,
		IsClaimable: true,
	})

	return zones
}
