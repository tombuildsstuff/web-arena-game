package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Airplane represents an airplane unit
type Airplane struct {
	BaseUnit
}

// NewAirplane creates a new airplane
func NewAirplane(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *Airplane {
	// Airplanes spawn at a higher Y position
	spawnPos.Y = types.AirplaneYPosition

	return &Airplane{
		BaseUnit: BaseUnit{
			ID:             uuid.New().String(),
			Type:           "airplane",
			OwnerID:        ownerID,
			Position:       spawnPos,
			Health:         types.AirplaneHealth,
			TargetPosition: targetPos,
			Speed:          types.AirplaneSpeed,
			Damage:         types.AirplaneDamage,
			AttackRange:    types.AirplaneAttackRange,
			AttackSpeed:    types.AirplaneAttackSpeed,
			LastAttackTime: 0,
		},
	}
}
