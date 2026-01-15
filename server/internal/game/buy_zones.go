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

// GetBuyZonesFromMap creates buy zones from a map definition
func GetBuyZonesFromMap(mapDef *types.MapDefinition) []*BuyZone {
	zones := make([]*BuyZone, 0, len(mapDef.BuyZones))

	for _, z := range mapDef.BuyZones {
		zones = append(zones, &BuyZone{
			ID:            z.ID,
			OwnerID:       z.DefaultOwner,
			UnitType:      z.UnitType,
			Position:      z.Position,
			Radius:        z.Radius,
			Cost:          z.Cost,
			ClaimCost:     z.ClaimCost,
			IsClaimable:   z.IsClaimable,
			ForwardBaseID: z.ForwardBaseID,
		})
	}

	return zones
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
