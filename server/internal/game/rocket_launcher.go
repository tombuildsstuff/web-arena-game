package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// RocketLauncher represents an infantry rocket launcher unit
// Long range, very high damage, low health, very slow fire rate
type RocketLauncher struct {
	BaseUnit
}

// NewRocketLauncher creates a new rocket launcher
func NewRocketLauncher(ownerID int, spawnPos types.Vector3, targetPos types.Vector3) *RocketLauncher {
	return &RocketLauncher{
		BaseUnit: BaseUnit{
			ID:              uuid.New().String(),
			Type:            "rocket_launcher",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.RocketLauncherHealth,
			MaxHealth:       types.RocketLauncherHealth,
			TargetPosition:  targetPos,
			Speed:           types.RocketLauncherSpeed,
			Damage:          types.RocketLauncherDamage,
			AttackRange:     types.RocketLauncherAttackRange,
			AttackSpeed:     types.RocketLauncherAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.RocketLauncherCollisionRadius,
		},
	}
}

// IsInfantry returns true - rocket launchers can claim barracks
func (r *RocketLauncher) IsInfantry() bool {
	return true
}
