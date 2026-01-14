package game

import (
	"fmt"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// PlayerUnit represents a player-controlled unit in the game
type PlayerUnit struct {
	BaseUnit
	RespawnTime    int64 // Unix timestamp when unit can respawn (0 if alive)
	IsRespawning   bool
	MoveDirection  types.Vector3 // Current movement direction from input
	BasePosition   types.Vector3 // Where to respawn
}

var playerUnitCounter = 0

// NewPlayerUnit creates a new player unit
func NewPlayerUnit(ownerID int, basePosition types.Vector3) *PlayerUnit {
	playerUnitCounter++
	spawnPos := basePosition
	spawnPos.Y = types.PlayerUnitYPosition

	return &PlayerUnit{
		BaseUnit: BaseUnit{
			ID:              fmt.Sprintf("player_%d", playerUnitCounter),
			Type:            "player",
			OwnerID:         ownerID,
			Position:        spawnPos,
			Health:          types.PlayerUnitHealth,
			TargetPosition:  spawnPos, // Players don't have AI targets
			Speed:           types.PlayerUnitSpeed,
			Damage:          types.PlayerUnitDamage,
			AttackRange:     types.PlayerUnitAttackRange,
			AttackSpeed:     types.PlayerUnitAttackSpeed,
			LastAttackTime:  0,
			CollisionRadius: types.PlayerCollisionRadius,
		},
		RespawnTime:   0,
		IsRespawning:  false,
		MoveDirection: types.Vector3{X: 0, Y: 0, Z: 0},
		BasePosition:  basePosition,
	}
}

// SetMoveDirection sets the movement direction from player input
func (p *PlayerUnit) SetMoveDirection(dir types.Vector3) {
	p.MoveDirection = dir
}

// GetMoveDirection returns the current movement direction
func (p *PlayerUnit) GetMoveDirection() types.Vector3 {
	return p.MoveDirection
}

// IsAlive returns true if the player unit is alive (not respawning)
func (p *PlayerUnit) IsAlive() bool {
	return p.Health > 0 && !p.IsRespawning
}

// TakeDamage damages the player unit
func (p *PlayerUnit) TakeDamage(amount int) {
	if p.IsRespawning {
		return // Can't take damage while respawning
	}

	p.Health -= amount
	if p.Health <= 0 {
		p.Health = 0
		p.startRespawn()
	}
}

// startRespawn begins the respawn timer
func (p *PlayerUnit) startRespawn() {
	p.IsRespawning = true
	p.RespawnTime = time.Now().UnixMilli() + int64(types.PlayerRespawnTime*1000)
}

// CheckRespawn checks if the respawn timer has expired and respawns if so
func (p *PlayerUnit) CheckRespawn() bool {
	if !p.IsRespawning {
		return false
	}

	now := time.Now().UnixMilli()
	if now >= p.RespawnTime {
		p.respawn()
		return true
	}
	return false
}

// respawn resets the player unit to their base
func (p *PlayerUnit) respawn() {
	p.Health = types.PlayerUnitHealth
	p.IsRespawning = false
	p.RespawnTime = 0
	p.Position = p.BasePosition
	p.Position.Y = types.PlayerUnitYPosition
	p.MoveDirection = types.Vector3{X: 0, Y: 0, Z: 0}
}

// GetRespawnTimeRemaining returns seconds until respawn (-1 if not respawning)
func (p *PlayerUnit) GetRespawnTimeRemaining() float64 {
	if !p.IsRespawning {
		return -1
	}
	remaining := float64(p.RespawnTime-time.Now().UnixMilli()) / 1000.0
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ToType converts PlayerUnit to types.Unit for JSON serialization
func (p *PlayerUnit) ToType() types.Unit {
	return types.Unit{
		ID:             p.ID,
		Type:           p.Type,
		OwnerID:        p.OwnerID,
		Position:       p.Position,
		Health:         p.Health,
		TargetPosition: p.TargetPosition,
		IsRespawning:   p.IsRespawning,
		RespawnTime:    p.GetRespawnTimeRemaining(),
	}
}
