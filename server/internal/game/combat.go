package game

import (
	"time"
)

// CombatSystem handles unit-to-unit combat
type CombatSystem struct {
	LOSSystem        *LOSSystem
	ProjectileSystem *ProjectileSystem
}

// NewCombatSystem creates a new combat system
func NewCombatSystem(losSystem *LOSSystem) *CombatSystem {
	return &CombatSystem{
		LOSSystem:        losSystem,
		ProjectileSystem: NewProjectileSystem(),
	}
}

// Update processes combat between units
func (s *CombatSystem) Update(state *State, deltaTime float64) {
	now := time.Now().UnixMilli()

	// Check each unit against all other units
	for i := range state.Units {
		attacker := state.Units[i]

		if !attacker.IsAlive() {
			continue
		}

		// Find enemies in range
		for j := range state.Units {
			if i == j {
				continue
			}

			target := state.Units[j]

			if !target.IsAlive() {
				continue
			}

			// Only attack enemy units
			if attacker.GetOwnerID() == target.GetOwnerID() {
				continue
			}

			// Check if target is in attack range
			distance := calculateDistance(attacker.GetPosition(), target.GetPosition())

			if distance <= attacker.GetAttackRange() {
				// Check line of sight (helicopters ignore LOS)
				if !s.LOSSystem.HasLineOfSightBetweenUnits(attacker, target) {
					continue // Can't see target, skip
				}

				// Check if enough time has passed since last attack
				timeSinceLastAttack := now - attacker.GetLastAttackTime()
				attackCooldown := int64(1000.0 / attacker.GetAttackSpeed()) // Convert attacks/sec to ms

				if timeSinceLastAttack >= attackCooldown {
					// Create projectile instead of instant damage
					projectile := NewProjectile(attacker, target, now)
					state.AddProjectile(projectile)
					attacker.SetLastAttackTime(now)

					break // Only attack one target per tick
				}
			}
		}
	}

	// Update projectiles
	s.ProjectileSystem.Update(state, deltaTime)

	// Remove dead units
	s.removeDeadUnits(state)
}

// removeDeadUnits removes units with health <= 0 (except players who respawn)
func (s *CombatSystem) removeDeadUnits(state *State) {
	aliveUnits := make([]Unit, 0, len(state.Units))

	for _, unit := range state.Units {
		// Keep player units (they respawn instead of being removed)
		if unit.GetType() == "player" {
			aliveUnits = append(aliveUnits, unit)
			continue
		}

		// Remove dead non-player units
		if unit.IsAlive() {
			aliveUnits = append(aliveUnits, unit)
		}
	}

	state.Units = aliveUnits
}
