package game

import (
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// TurretSystem handles turret combat and updates
type TurretSystem struct {
	LOSSystem *LOSSystem
}

// NewTurretSystem creates a new turret system
func NewTurretSystem(losSystem *LOSSystem) *TurretSystem {
	return &TurretSystem{
		LOSSystem: losSystem,
	}
}

// Update processes turret updates including respawns and combat
func (s *TurretSystem) Update(state *State, deltaTime float64) {
	now := time.Now().UnixMilli()

	for _, turret := range state.Turrets {
		// Update respawn timers
		turret.Update(deltaTime)

		// Check for auto-claiming by tanks/helicopters passing by unclaimed turrets
		if !turret.IsDestroyed && turret.OwnerID == -1 {
			s.checkAutoClaimByUnits(turret, state)
		}

		// Skip combat if destroyed or unclaimed
		if turret.IsDestroyed || turret.OwnerID == -1 {
			continue
		}

		// Find and attack enemies
		s.processTurretCombat(turret, state, now)
	}
}

// checkAutoClaimByUnits checks if any tank or helicopter is near an unclaimed turret and claims it
func (s *TurretSystem) checkAutoClaimByUnits(turret *Turret, state *State) {
	for _, unit := range state.Units {
		// Only tanks and airplanes can auto-claim
		unitType := unit.GetType()
		if unitType != "tank" && unitType != "airplane" {
			continue
		}

		// Skip dead units
		if !unit.IsAlive() {
			continue
		}

		// Check if unit is in range
		if turret.IsPlayerInRange(unit.GetPosition()) {
			// Claim the turret for this unit's team
			turret.Claim(unit.GetOwnerID())
			return // Only one unit can claim per tick
		}
	}
}

// processTurretCombat handles a single turret's combat logic
func (s *TurretSystem) processTurretCombat(turret *Turret, state *State, now int64) {
	// Find the closest enemy unit in range with LOS
	var closestTarget Unit
	closestDistance := turret.AttackRange + 1 // Start beyond range

	for _, unit := range state.Units {
		if !unit.IsAlive() {
			continue
		}

		// Only attack enemy units (units owned by the other player)
		if unit.GetOwnerID() == turret.OwnerID {
			continue
		}

		// Calculate distance
		distance := calculateDistance(turret.Position, unit.GetPosition())

		if distance <= turret.AttackRange && distance < closestDistance {
			// Check line of sight
			turretPos := turret.Position
			targetPos := unit.GetPosition()

			// Turrets are at Y=3, need to check LOS
			if s.LOSSystem.HasLineOfSight(turretPos, targetPos, false) {
				closestTarget = unit
				closestDistance = distance
			}
		}
	}

	// Handle target tracking
	if closestTarget == nil {
		// No target - clear tracking
		turret.CurrentTargetID = ""
		turret.TargetAcquiredAt = 0
		return
	}

	targetID := closestTarget.GetID()

	// Check if this is a new target
	if turret.CurrentTargetID != targetID {
		// New target - start tracking
		turret.CurrentTargetID = targetID
		turret.TargetAcquiredAt = now
		return // Don't fire yet, need to track first
	}

	// Check if we've tracked long enough
	trackingDuration := now - turret.TargetAcquiredAt
	if trackingDuration < types.TurretTrackingTime {
		return // Still tracking, don't fire yet
	}

	// Check attack cooldown
	if !turret.CanAttack() {
		return
	}

	// Fire at the target
	projectile := s.createTurretProjectile(turret, closestTarget, now)
	state.AddProjectile(projectile)
	turret.LastAttackTime = now
}

// createTurretProjectile creates a projectile from a turret
func (s *TurretSystem) createTurretProjectile(turret *Turret, target Unit, now int64) *Projectile {
	startPos := turret.Position
	endPos := target.GetPosition()

	return &Projectile{
		ID:        generateProjectileID(),
		ShooterID: turret.ID, // Use turret ID as shooter
		TargetID:  target.GetID(),
		Position:  startPos,
		StartPos:  startPos,
		EndPos:    endPos,
		Speed:     ProjectileSpeed,
		Damage:    turret.Damage,
		CreatedAt: now,
	}
}

// ProcessTurretDamage handles projectiles hitting turrets
func (s *TurretSystem) ProcessTurretDamage(state *State, projectile *Projectile) bool {
	// Check if the projectile target is a turret
	for _, turret := range state.Turrets {
		if turret.ID == projectile.TargetID {
			if turret.IsAlive() {
				turret.TakeDamage(projectile.Damage)
				return true
			}
		}
	}
	return false
}
