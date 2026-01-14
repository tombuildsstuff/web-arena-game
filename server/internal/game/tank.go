package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Tank represents a tank unit
type Tank struct {
	BaseUnit
}

// NewTank creates a new tank
func NewTank(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *Tank {
	return &Tank{
		BaseUnit: BaseUnit{
			ID:              uuid.New().String(),
			Type:            "tank",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.TankHealth,
			TargetPosition:  targetPos,
			Speed:           types.TankSpeed,
			Damage:          types.TankDamage,
			AttackRange:     types.TankAttackRange,
			AttackSpeed:     types.TankAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.TankCollisionRadius,
		},
	}
}
