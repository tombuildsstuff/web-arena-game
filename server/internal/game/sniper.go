package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Sniper represents an infantry sniper unit
// High range, high damage, low health, slow fire rate
type Sniper struct {
	BaseUnit
}

// NewSniper creates a new sniper
func NewSniper(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *Sniper {
	return &Sniper{
		BaseUnit: BaseUnit{
			ID:              uuid.New().String(),
			Type:            "sniper",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.SniperHealth,
			MaxHealth:       types.SniperHealth,
			TargetPosition:  targetPos,
			Speed:           types.SniperSpeed,
			Damage:          types.SniperDamage,
			AttackRange:     types.SniperAttackRange,
			AttackSpeed:     types.SniperAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.SniperCollisionRadius,
		},
	}
}

// IsInfantry returns true - snipers can claim barracks
func (s *Sniper) IsInfantry() bool {
	return true
}
