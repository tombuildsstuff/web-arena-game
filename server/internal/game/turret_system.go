package game

import (
	"time"
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

		// Skip combat if destroyed or unclaimed
		if turret.IsDestroyed || turret.OwnerID == -1 {
			continue
		}

		// Find and attack enemies
		s.processTurretCombat(turret, state, now)
	}
}

// processTurretCombat handles a single turret's combat logic
func (s *TurretSystem) processTurretCombat(turret *Turret, state *State, now int64) {
	if !turret.CanAttack() {
		return
	}

	// Find the closest enemy unit in range
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

	// Attack the closest target
	if closestTarget != nil {
		projectile := s.createTurretProjectile(turret, closestTarget, now)
		state.AddProjectile(projectile)
		turret.LastAttackTime = now
	}
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
