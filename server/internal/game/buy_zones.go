package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// BuyZone represents a location where players can purchase units
type BuyZone struct {
	ID            string
	OwnerID       int
	UnitType      string // "tank", "airplane", "super_tank", "super_helicopter", or "" for base zones
	Position      types.Vector3
	Radius        float64
	Cost          int    // Cost to buy units from this zone
	ClaimCost     int    // Cost to claim this zone (0 if not claimable)
	IsClaimable   bool   // Whether this zone can be claimed by players
	ForwardBaseID string // ID of parent forward base (empty if this IS a base or not part of one)
}

// ToType converts BuyZone to types.BuyZone for JSON serialization
func (b *BuyZone) ToType() types.BuyZone {
	return types.BuyZone{
		ID:            b.ID,
		OwnerID:       b.OwnerID,
		UnitType:      b.UnitType,
		Position:      b.Position,
		Radius:        b.Radius,
		Cost:          b.Cost,
		ClaimCost:     b.ClaimCost,
		IsClaimable:   b.IsClaimable,
		ForwardBaseID: b.ForwardBaseID,
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

	// ============================================================
	// CLAIMABLE FORWARD BASES
	// Each forward base is a single claimable unit with an enlarged ring.
	// Once claimed, the player can purchase from the 3 buy zones within it.
	// ============================================================

	// North forward base - the claimable base zone (large radius)
	zones = append(zones, &BuyZone{
		ID:          "forward_base_north",
		OwnerID:     -1, // Neutral - can be claimed
		UnitType:    "",  // Empty - this is a base, not a unit purchase zone
		Position:    types.Vector3{X: 0, Y: 3, Z: -70},
		Radius:      12.0, // Large radius covering the whole platform
		Cost:        0,
		ClaimCost:   types.ForwardBaseClaimCost,
		IsClaimable: true,
	})
	// North forward base child zones (not claimable, linked to parent base)
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_north_tank",
		OwnerID:       -1,
		UnitType:      "tank",
		Position:      types.Vector3{X: 0, Y: 3, Z: -70},
		Radius:        4.0,
		Cost:          types.TankCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_north",
	})
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_north_super_tank",
		OwnerID:       -1,
		UnitType:      "super_tank",
		Position:      types.Vector3{X: -6, Y: 3, Z: -70},
		Radius:        4.0,
		Cost:          types.SuperTankCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_north",
	})
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_north_super_helicopter",
		OwnerID:       -1,
		UnitType:      "super_helicopter",
		Position:      types.Vector3{X: 6, Y: 3, Z: -70},
		Radius:        4.0,
		Cost:          types.SuperHelicopterCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_north",
	})

	// South forward base - the claimable base zone (large radius)
	zones = append(zones, &BuyZone{
		ID:          "forward_base_south",
		OwnerID:     -1, // Neutral - can be claimed
		UnitType:    "",  // Empty - this is a base, not a unit purchase zone
		Position:    types.Vector3{X: 0, Y: 3, Z: 70},
		Radius:      12.0, // Large radius covering the whole platform
		Cost:        0,
		ClaimCost:   types.ForwardBaseClaimCost,
		IsClaimable: true,
	})
	// South forward base child zones (not claimable, linked to parent base)
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_south_tank",
		OwnerID:       -1,
		UnitType:      "tank",
		Position:      types.Vector3{X: 0, Y: 3, Z: 70},
		Radius:        4.0,
		Cost:          types.TankCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_south",
	})
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_south_super_tank",
		OwnerID:       -1,
		UnitType:      "super_tank",
		Position:      types.Vector3{X: -6, Y: 3, Z: 70},
		Radius:        4.0,
		Cost:          types.SuperTankCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_south",
	})
	zones = append(zones, &BuyZone{
		ID:            "buy_forward_south_super_helicopter",
		OwnerID:       -1,
		UnitType:      "super_helicopter",
		Position:      types.Vector3{X: 6, Y: 3, Z: 70},
		Radius:        4.0,
		Cost:          types.SuperHelicopterCost,
		ClaimCost:     0,
		IsClaimable:   false,
		ForwardBaseID: "forward_base_south",
	})

	return zones
}
