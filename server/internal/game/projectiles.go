package game

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

const (
	// ProjectileSpeed is the speed of projectiles in units per second
	ProjectileSpeed = 50.0
)

// Projectile represents a traveling projectile
type Projectile struct {
	ID        string
	ShooterID string
	TargetID  string
	Position  types.Vector3
	StartPos  types.Vector3
	EndPos    types.Vector3
	Speed     float64
	Damage    int
	CreatedAt int64
}

// NewProjectile creates a new projectile targeting a unit
func NewProjectile(shooter, target Unit, timestamp int64) *Projectile {
	shooterPos := shooter.GetPosition()
	targetPos := target.GetPosition()

	return &Projectile{
		ID:        uuid.New().String(),
		ShooterID: shooter.GetID(),
		TargetID:  target.GetID(),
		Position:  shooterPos,
		StartPos:  shooterPos,
		EndPos:    targetPos,
		Speed:     ProjectileSpeed,
		Damage:    shooter.GetDamage(),
		CreatedAt: timestamp,
	}
}

// NewProjectileFromPlayer creates a projectile from a player to a position
// targetUnit can be nil if shooting at empty space
func NewProjectileFromPlayer(shooter *PlayerUnit, targetPos types.Vector3, targetUnit Unit, timestamp int64) *Projectile {
	shooterPos := shooter.GetPosition()

	targetID := ""
	if targetUnit != nil {
		targetID = targetUnit.GetID()
	}

	return &Projectile{
		ID:        uuid.New().String(),
		ShooterID: shooter.GetID(),
		TargetID:  targetID,
		Position:  shooterPos,
		StartPos:  shooterPos,
		EndPos:    targetPos,
		Speed:     ProjectileSpeed,
		Damage:    shooter.GetDamage(),
		CreatedAt: timestamp,
	}
}

// NewProjectileFromPlayerToTurret creates a projectile from a player to a turret
func NewProjectileFromPlayerToTurret(shooter *PlayerUnit, turret *Turret, timestamp int64) *Projectile {
	shooterPos := shooter.GetPosition()
	targetPos := turret.Position

	return &Projectile{
		ID:        uuid.New().String(),
		ShooterID: shooter.GetID(),
		TargetID:  turret.ID,
		Position:  shooterPos,
		StartPos:  shooterPos,
		EndPos:    targetPos,
		Speed:     ProjectileSpeed,
		Damage:    shooter.GetDamage(),
		CreatedAt: timestamp,
	}
}

// ToType converts Projectile to types.Projectile for JSON serialization
func (p *Projectile) ToType() types.Projectile {
	return types.Projectile{
		ID:        p.ID,
		ShooterID: p.ShooterID,
		TargetID:  p.TargetID,
		Position:  p.Position,
		StartPos:  p.StartPos,
		EndPos:    p.EndPos,
		Speed:     p.Speed,
		Damage:    p.Damage,
		CreatedAt: p.CreatedAt,
	}
}

// generateProjectileID generates a unique ID for a projectile
func generateProjectileID() string {
	return uuid.New().String()
}

// ProjectileSystem handles projectile movement and hit detection
type ProjectileSystem struct{}

// NewProjectileSystem creates a new projectile system
func NewProjectileSystem() *ProjectileSystem {
	return &ProjectileSystem{}
}

// Update moves all projectiles and handles hit detection
func (ps *ProjectileSystem) Update(state *State, deltaTime float64) {
	toRemove := make([]string, 0)

	for _, proj := range state.Projectiles {
		// Calculate direction to target
		direction := normalize(subtract(proj.EndPos, proj.Position))

		// Calculate movement distance
		movement := proj.Speed * deltaTime

		// Update position
		proj.Position.X += direction.X * movement
		proj.Position.Y += direction.Y * movement
		proj.Position.Z += direction.Z * movement

		// Check if reached target (within 2 units)
		distToTarget := calculateDistance(proj.Position, proj.EndPos)
		if distToTarget < 2.0 {
			hit := false

			// Try to apply damage to unit target
			target := state.GetUnitByID(proj.TargetID)
			if target != nil && target.IsAlive() {
				wasAlive := target.IsAlive()
				target.TakeDamage(proj.Damage)
				hit = true

				// If target died from this hit, credit the kill
				if wasAlive && !target.IsAlive() {
					// Find the shooter (could be unit or turret)
					shooter := state.GetUnitByID(proj.ShooterID)
					if shooter != nil {
						shooterOwner := state.GetPlayer(shooter.GetOwnerID())
						if shooterOwner != nil {
							shooterOwner.AddKill()
						}
					} else {
						// Check if shooter is a turret
						turret := state.GetTurretByID(proj.ShooterID)
						if turret != nil {
							shooterOwner := state.GetPlayer(turret.OwnerID)
							if shooterOwner != nil {
								shooterOwner.AddKill()
							}
						}
					}
				}
			}

			// Try to apply damage to turret target
			if !hit {
				turret := state.GetTurretByID(proj.TargetID)
				if turret != nil && turret.IsAlive() {
					turret.TakeDamage(proj.Damage)
					hit = true
				}
			}

			toRemove = append(toRemove, proj.ID)
		}

		// Also remove if projectile has traveled too far (miss)
		distFromStart := calculateDistance(proj.Position, proj.StartPos)
		if distFromStart > 200 { // Max range
			toRemove = append(toRemove, proj.ID)
		}
	}

	// Remove hit/expired projectiles
	if len(toRemove) > 0 {
		state.RemoveProjectiles(toRemove)
	}
}
