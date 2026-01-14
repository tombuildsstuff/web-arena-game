package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Unit interface for all units
type Unit interface {
	GetID() string
	GetType() string
	GetOwnerID() int
	GetPosition() types.Vector3
	SetPosition(pos types.Vector3)
	GetHealth() int
	SetHealth(health int)
	GetTargetPosition() types.Vector3
	GetSpeed() float64
	GetDamage() int
	GetAttackRange() float64
	GetAttackSpeed() float64
	GetLastAttackTime() int64
	SetLastAttackTime(time int64)
	ToType() types.Unit
	IsAlive() bool
	TakeDamage(amount int)
	// Pathfinding support
	GetWaypoints() []types.Vector3
	SetWaypoints(waypoints []types.Vector3)
	GetCurrentWaypoint() int
	SetCurrentWaypoint(index int)
	ClearWaypoints()
}

// BaseUnit provides common functionality for all units
type BaseUnit struct {
	ID             string
	Type           string
	OwnerID        int
	Position       types.Vector3
	Health         int
	TargetPosition types.Vector3
	Speed          float64
	Damage         int
	AttackRange    float64
	AttackSpeed    float64
	LastAttackTime int64 // Unix timestamp in milliseconds
	// Pathfinding fields
	Waypoints       []types.Vector3
	CurrentWaypoint int
}

func (u *BaseUnit) GetID() string {
	return u.ID
}

func (u *BaseUnit) GetType() string {
	return u.Type
}

func (u *BaseUnit) GetOwnerID() int {
	return u.OwnerID
}

func (u *BaseUnit) GetPosition() types.Vector3 {
	return u.Position
}

func (u *BaseUnit) SetPosition(pos types.Vector3) {
	u.Position = pos
}

func (u *BaseUnit) GetHealth() int {
	return u.Health
}

func (u *BaseUnit) SetHealth(health int) {
	u.Health = health
}

func (u *BaseUnit) GetTargetPosition() types.Vector3 {
	return u.TargetPosition
}

func (u *BaseUnit) GetSpeed() float64 {
	return u.Speed
}

func (u *BaseUnit) GetDamage() int {
	return u.Damage
}

func (u *BaseUnit) GetAttackRange() float64 {
	return u.AttackRange
}

func (u *BaseUnit) GetAttackSpeed() float64 {
	return u.AttackSpeed
}

func (u *BaseUnit) GetLastAttackTime() int64 {
	return u.LastAttackTime
}

func (u *BaseUnit) SetLastAttackTime(time int64) {
	u.LastAttackTime = time
}

func (u *BaseUnit) IsAlive() bool {
	return u.Health > 0
}

func (u *BaseUnit) TakeDamage(amount int) {
	u.Health -= amount
	if u.Health < 0 {
		u.Health = 0
	}
}

func (u *BaseUnit) ToType() types.Unit {
	return types.Unit{
		ID:             u.ID,
		Type:           u.Type,
		OwnerID:        u.OwnerID,
		Position:       u.Position,
		Health:         u.Health,
		TargetPosition: u.TargetPosition,
	}
}

// Pathfinding methods

func (u *BaseUnit) GetWaypoints() []types.Vector3 {
	return u.Waypoints
}

func (u *BaseUnit) SetWaypoints(waypoints []types.Vector3) {
	u.Waypoints = waypoints
	u.CurrentWaypoint = 0
}

func (u *BaseUnit) GetCurrentWaypoint() int {
	return u.CurrentWaypoint
}

func (u *BaseUnit) SetCurrentWaypoint(index int) {
	u.CurrentWaypoint = index
}

func (u *BaseUnit) ClearWaypoints() {
	u.Waypoints = nil
	u.CurrentWaypoint = 0
}
