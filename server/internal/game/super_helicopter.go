package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// SuperHelicopter represents a heavy helicopter unit with 2x firepower and 3x health
type SuperHelicopter struct {
	BaseUnit
}

// NewSuperHelicopter creates a new super helicopter
func NewSuperHelicopter(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *SuperHelicopter {
	// Super helicopters spawn at a higher Y position like regular airplanes
	spawnPos.Y = types.AirplaneYPosition

	return &SuperHelicopter{
		BaseUnit: BaseUnit{
			ID:              uuid.New().String(),
			Type:            "super_helicopter",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.SuperHelicopterHealth,
			MaxHealth:       types.SuperHelicopterHealth,
			TargetPosition:  targetPos,
			Speed:           types.SuperHelicopterSpeed,
			Damage:          types.SuperHelicopterDamage,
			AttackRange:     types.SuperHelicopterAttackRange,
			AttackSpeed:     types.SuperHelicopterAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.SuperHelicopterCollisionRadius,
		},
	}
}
