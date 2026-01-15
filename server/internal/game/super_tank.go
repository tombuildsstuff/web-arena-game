package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// SuperTank represents a heavy tank unit with 2x firepower and 3x health
type SuperTank struct {
	BaseUnit
}

// NewSuperTank creates a new super tank
func NewSuperTank(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *SuperTank {
	return &SuperTank{
		BaseUnit: BaseUnit{
			ID:              uuid.New().String(),
			Type:            "super_tank",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.SuperTankHealth,
			MaxHealth:       types.SuperTankHealth,
			TargetPosition:  targetPos,
			Speed:           types.SuperTankSpeed,
			Damage:          types.SuperTankDamage,
			AttackRange:     types.SuperTankAttackRange,
			AttackSpeed:     types.SuperTankAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.SuperTankCollisionRadius,
		},
	}
}
