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
	SetTargetPosition(pos types.Vector3)
	GetSpeed() float64
	GetDamage() int
	GetAttackRange() float64
	GetAttackSpeed() float64
	GetLastAttackTime() int64
	SetLastAttackTime(time int64)
	GetCollisionRadius() float64
	ToType() types.Unit
	IsAlive() bool
	TakeDamage(amount int)
	Heal(amount int)
	GetMaxHealth() int
	// Pathfinding support
	GetWaypoints() []types.Vector3
	SetWaypoints(waypoints []types.Vector3)
	GetCurrentWaypoint() int
	SetCurrentWaypoint(index int)
	ClearWaypoints()
	// Patrol mode support (for perimeter sweep)
	IsPatrolling() bool
	SetPatrolling(patrol bool)
	GetPatrolCorner() int
	SetPatrolCorner(corner int)
	// Stuck detection support
	GetLastPosition() types.Vector3
	SetLastPosition(pos types.Vector3)
	GetStuckTicks() int
	SetStuckTicks(ticks int)
	GetPathAttempts() int
	SetPathAttempts(attempts int)
	// Avoidance direction persistence (to prevent flickering)
	GetAvoidanceDirection() int // -1 = left, 0 = none, 1 = right
	SetAvoidanceDirection(dir int)
	GetAvoidanceTicks() int
	SetAvoidanceTicks(ticks int)
	// Infantry check (for barracks claiming)
	IsInfantry() bool
}

// BaseUnit provides common functionality for all units
type BaseUnit struct {
	ID              string
	Type            string
	OwnerID         int
	Position        types.Vector3
	Health          int
	MaxHealth       int
	TargetPosition  types.Vector3
	Speed           float64
	Damage          int
	AttackRange     float64
	AttackSpeed     float64
	LastAttackTime  int64   // Unix timestamp in milliseconds
	CollisionRadius float64 // Radius for unit-to-unit collision
	// Pathfinding fields
	Waypoints       []types.Vector3
	CurrentWaypoint int
	// Stuck detection fields
	LastPosition types.Vector3
	StuckTicks   int
	PathAttempts int // Number of path recalculations when stuck
	// Patrol mode fields (for perimeter sweep)
	Patrolling   bool
	PatrolCorner int // Index of current patrol corner target (0-3)
	// Avoidance direction persistence (to prevent flickering when tanks collide)
	AvoidanceDirection int // -1 = left, 0 = none, 1 = right
	AvoidanceTicks     int // Ticks remaining to keep using the same avoidance direction
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

func (u *BaseUnit) SetTargetPosition(pos types.Vector3) {
	u.TargetPosition = pos
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

func (u *BaseUnit) GetCollisionRadius() float64 {
	return u.CollisionRadius
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

func (u *BaseUnit) Heal(amount int) {
	u.Health += amount
	if u.Health > u.MaxHealth {
		u.Health = u.MaxHealth
	}
}

func (u *BaseUnit) GetMaxHealth() int {
	return u.MaxHealth
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

// Patrol mode methods

func (u *BaseUnit) IsPatrolling() bool {
	return u.Patrolling
}

func (u *BaseUnit) SetPatrolling(patrol bool) {
	u.Patrolling = patrol
}

func (u *BaseUnit) GetPatrolCorner() int {
	return u.PatrolCorner
}

func (u *BaseUnit) SetPatrolCorner(corner int) {
	u.PatrolCorner = corner
}

// Stuck detection methods

func (u *BaseUnit) GetLastPosition() types.Vector3 {
	return u.LastPosition
}

func (u *BaseUnit) SetLastPosition(pos types.Vector3) {
	u.LastPosition = pos
}

func (u *BaseUnit) GetStuckTicks() int {
	return u.StuckTicks
}

func (u *BaseUnit) SetStuckTicks(ticks int) {
	u.StuckTicks = ticks
}

func (u *BaseUnit) GetPathAttempts() int {
	return u.PathAttempts
}

func (u *BaseUnit) SetPathAttempts(attempts int) {
	u.PathAttempts = attempts
}

// Avoidance direction methods

func (u *BaseUnit) GetAvoidanceDirection() int {
	return u.AvoidanceDirection
}

func (u *BaseUnit) SetAvoidanceDirection(dir int) {
	u.AvoidanceDirection = dir
}

func (u *BaseUnit) GetAvoidanceTicks() int {
	return u.AvoidanceTicks
}

func (u *BaseUnit) SetAvoidanceTicks(ticks int) {
	u.AvoidanceTicks = ticks
}

// IsInfantry returns false by default - only infantry units override this
func (u *BaseUnit) IsInfantry() bool {
	return false
}
