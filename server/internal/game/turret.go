package game

import (
	"fmt"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Turret represents a claimable defensive turret
type Turret struct {
	ID             string
	Position       types.Vector3
	OwnerID        int     // -1 = unclaimed, 0 = player 1, 1 = player 2
	DefaultOwnerID int     // Owner when respawning (-1 for middle turrets)
	Health         int
	MaxHealth      int
	IsDestroyed    bool
	RespawnTime    float64 // Seconds remaining until respawn
	ClaimRadius    float64
	LastAttackTime int64
	AttackRange    float64
	AttackSpeed    float64
	Damage         int
	// Tracking delay fields
	CurrentTargetID   string // ID of the unit being tracked
	TargetAcquiredAt  int64  // When the turret started tracking this target
}

// NewTurret creates a new turret
func NewTurret(id string, position types.Vector3, defaultOwnerID int) *Turret {
	return &Turret{
		ID:             id,
		Position:       position,
		OwnerID:        defaultOwnerID,
		DefaultOwnerID: defaultOwnerID,
		Health:         types.TurretHealth,
		MaxHealth:      types.TurretHealth,
		IsDestroyed:    false,
		RespawnTime:    0,
		ClaimRadius:    types.TurretClaimRadius,
		LastAttackTime: 0,
		AttackRange:    types.TurretAttackRange,
		AttackSpeed:    types.TurretAttackSpeed,
		Damage:         types.TurretDamage,
	}
}

// ToType converts Turret to types.Turret for JSON serialization
func (t *Turret) ToType() types.Turret {
	// Calculate tracking progress
	isTracking := t.CurrentTargetID != ""
	trackingProgress := 0.0
	if isTracking && t.TargetAcquiredAt > 0 {
		now := time.Now().UnixMilli()
		elapsed := float64(now - t.TargetAcquiredAt)
		trackingProgress = elapsed / float64(types.TurretTrackingTime)
		if trackingProgress > 1.0 {
			trackingProgress = 1.0
		}
	}

	return types.Turret{
		ID:               t.ID,
		Position:         t.Position,
		OwnerID:          t.OwnerID,
		DefaultOwnerID:   t.DefaultOwnerID,
		Health:           t.Health,
		MaxHealth:        t.MaxHealth,
		IsDestroyed:      t.IsDestroyed,
		RespawnTime:      t.RespawnTime,
		ClaimRadius:      t.ClaimRadius,
		IsTracking:       isTracking,
		TrackingProgress: trackingProgress,
	}
}

// IsPlayerInRange checks if a player position is within claiming range
func (t *Turret) IsPlayerInRange(pos types.Vector3) bool {
	dx := pos.X - t.Position.X
	dz := pos.Z - t.Position.Z
	distSq := dx*dx + dz*dz
	return distSq <= t.ClaimRadius*t.ClaimRadius
}

// CanBeClaimed returns true if the turret can be claimed by the given player
func (t *Turret) CanBeClaimed(playerID int) bool {
	// Can't claim if destroyed
	if t.IsDestroyed {
		return false
	}
	// Can't claim your own turret
	if t.OwnerID == playerID {
		return false
	}
	// Can only claim neutral turrets - must destroy enemy turrets first
	if t.OwnerID != -1 {
		return false
	}
	return true
}

// Claim sets the turret owner
func (t *Turret) Claim(playerID int) {
	t.OwnerID = playerID
}

// TakeDamage applies damage to the turret
func (t *Turret) TakeDamage(amount int) {
	if t.IsDestroyed {
		return
	}
	t.Health -= amount
	if t.Health <= 0 {
		t.Health = 0
		t.IsDestroyed = true
		t.RespawnTime = types.TurretRespawnTime
	}
}

// Update handles respawn timing
func (t *Turret) Update(deltaTime float64) {
	if !t.IsDestroyed {
		return
	}

	t.RespawnTime -= deltaTime
	if t.RespawnTime <= 0 {
		t.Respawn()
	}
}

// Respawn resets the turret to its default state
func (t *Turret) Respawn() {
	t.IsDestroyed = false
	t.Health = t.MaxHealth
	t.RespawnTime = 0
	t.OwnerID = t.DefaultOwnerID
	t.LastAttackTime = 0
	t.CurrentTargetID = ""
	t.TargetAcquiredAt = 0
}

// IsAlive returns true if the turret is not destroyed
func (t *Turret) IsAlive() bool {
	return !t.IsDestroyed
}

// CanAttack checks if the turret can attack (has owner and not on cooldown)
func (t *Turret) CanAttack() bool {
	// Must have an owner to attack
	if t.OwnerID == -1 {
		return false
	}
	if t.IsDestroyed {
		return false
	}

	now := time.Now().UnixMilli()
	timeSinceLastAttack := now - t.LastAttackTime
	attackCooldown := int64(1000.0 / t.AttackSpeed)

	return timeSinceLastAttack >= attackCooldown
}

// GetTurretsFromMap creates turrets from a map definition
func GetTurretsFromMap(mapDef *types.MapDefinition) []*Turret {
	turrets := make([]*Turret, 0, len(mapDef.Turrets))

	for _, t := range mapDef.Turrets {
		turrets = append(turrets, NewTurret(t.ID, t.Position, t.DefaultOwner))
	}

	return turrets
}

// GetTurrets returns turrets for the arena
// Base turrets are on corners of each base, middle turrets are neutral
//
// Deprecated: Use GetTurretsFromMap with a MapDefinition instead
func GetTurrets() []*Turret {
	turrets := make([]*Turret, 0)
	idCounter := 0

	nextID := func() string {
		idCounter++
		return fmt.Sprintf("turret_%d", idCounter)
	}

	// ============================================================
	// PLAYER 1 BASE TURRETS (4 corners of base at X=-90)
	// Base dimensions: width=25 (X: -102.5 to -77.5, clamped to -100)
	//                  depth=28.75 (Z: -14.375 to 14.375)
	// ============================================================
	// Front-North corner (toward center, north side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -77, Y: 3, Z: -15},
		0, // Owned by player 1 by default
	))
	// Front-South corner (toward center, south side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -77, Y: 3, Z: 15},
		0, // Owned by player 1 by default
	))
	// Back-North corner (toward wall, north side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -98, Y: 3, Z: -15},
		0, // Owned by player 1 by default
	))
	// Back-South corner (toward wall, south side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -98, Y: 3, Z: 15},
		0, // Owned by player 1 by default
	))

	// ============================================================
	// PLAYER 2 BASE TURRETS (4 corners of base at X=90)
	// Base dimensions: width=25 (X: 77.5 to 102.5, clamped to 100)
	//                  depth=28.75 (Z: -14.375 to 14.375)
	// ============================================================
	// Front-North corner (toward center, north side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 77, Y: 3, Z: -15},
		1, // Owned by player 2 by default
	))
	// Front-South corner (toward center, south side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 77, Y: 3, Z: 15},
		1, // Owned by player 2 by default
	))
	// Back-North corner (toward wall, north side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 98, Y: 3, Z: -15},
		1, // Owned by player 2 by default
	))
	// Back-South corner (toward wall, south side)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 98, Y: 3, Z: 15},
		1, // Owned by player 2 by default
	))

	// ============================================================
	// NEUTRAL TURRETS (claimable by either player)
	// ============================================================

	// Central corridor turrets (high value - control the middle)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 0, Y: 3, Z: -50},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 0, Y: 3, Z: 50},
		-1, // Unclaimed
	))

	// Mid-lane turrets (between bases and center)
	// Left side mid turrets
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -50, Y: 3, Z: -40},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -50, Y: 3, Z: 40},
		-1, // Unclaimed
	))
	// Right side mid turrets
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 50, Y: 3, Z: -40},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 50, Y: 3, Z: 40},
		-1, // Unclaimed
	))

	// Center room turrets (control the middle lane)
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -25, Y: 3, Z: 0},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 25, Y: 3, Z: 0},
		-1, // Unclaimed
	))

	// Corner room turrets (outer flanking positions)
	// Top corners
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -55, Y: 3, Z: -65},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 55, Y: 3, Z: -65},
		-1, // Unclaimed
	))
	// Bottom corners
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: -55, Y: 3, Z: 65},
		-1, // Unclaimed
	))
	turrets = append(turrets, NewTurret(
		nextID(),
		types.Vector3{X: 55, Y: 3, Z: 65},
		-1, // Unclaimed
	))

	return turrets
}
