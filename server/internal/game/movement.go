package game

import (
	"math"
	"math/rand"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// MovementSystem handles unit movement
type MovementSystem struct {
	Pathfinding *PathfindingSystem
}

// NewMovementSystem creates a new movement system
func NewMovementSystem(pathfinding *PathfindingSystem) *MovementSystem {
	return &MovementSystem{
		Pathfinding: pathfinding,
	}
}

// Update updates all unit positions
func (s *MovementSystem) Update(state *State, deltaTime float64) {
	for i := range state.Units {
		unit := state.Units[i]

		// Handle different unit types
		switch unit.GetType() {
		case "tank":
			s.updateTankMovement(unit, state, deltaTime)
		case "player":
			s.updatePlayerMovement(unit, state, deltaTime)
		default: // airplane
			s.updateDirectMovement(unit, deltaTime)
		}
	}
}

// updatePlayerMovement handles player-controlled movement
func (s *MovementSystem) updatePlayerMovement(unit Unit, state *State, deltaTime float64) {
	playerUnit, ok := unit.(*PlayerUnit)
	if !ok {
		return
	}

	// Check for respawn
	if playerUnit.CheckRespawn() {
		return // Just respawned, skip movement this frame
	}

	// Skip if dead/respawning
	if !playerUnit.IsAlive() {
		return
	}

	// Get movement direction from input
	dir := playerUnit.GetMoveDirection()

	// If no movement input, don't move
	if dir.X == 0 && dir.Z == 0 {
		return
	}

	// Normalize direction
	dir = normalize(dir)

	// Calculate new position
	speed := playerUnit.GetSpeed()
	newPos := types.Vector3{
		X: playerUnit.Position.X + dir.X*speed*deltaTime,
		Y: playerUnit.Position.Y,
		Z: playerUnit.Position.Z + dir.Z*speed*deltaTime,
	}

	// Check collision with obstacles using pathfinding system (with player radius)
	playerRadius := playerUnit.GetCollisionRadius()
	obstacleBlocked := false
	if s.Pathfinding != nil && !s.isPositionWalkableWithRadius(newPos, playerRadius) {
		obstacleBlocked = true
		// Try sliding along walls
		// Try X movement only
		testPosX := types.Vector3{
			X: newPos.X,
			Y: playerUnit.Position.Y,
			Z: playerUnit.Position.Z,
		}
		if s.isPositionWalkableWithRadius(testPosX, playerRadius) {
			newPos = testPosX
			obstacleBlocked = false
		} else {
			// Try Z movement only
			testPosZ := types.Vector3{
				X: playerUnit.Position.X,
				Y: playerUnit.Position.Y,
				Z: newPos.Z,
			}
			if s.isPositionWalkableWithRadius(testPosZ, playerRadius) {
				newPos = testPosZ
				obstacleBlocked = false
			}
		}
	}

	// If completely blocked by obstacles, can't move at all
	if obstacleBlocked {
		return
	}

	// Check collision with other units
	if s.wouldCollideWithUnit(playerUnit, newPos, state.Units) {
		// Try sliding along X axis only
		testPosX := types.Vector3{
			X: newPos.X,
			Y: playerUnit.Position.Y,
			Z: playerUnit.Position.Z,
		}
		if !s.wouldCollideWithUnit(playerUnit, testPosX, state.Units) &&
			(s.Pathfinding == nil || s.isPositionWalkableWithRadius(testPosX, playerRadius)) {
			newPos = testPosX
		} else {
			// Try sliding along Z axis only
			testPosZ := types.Vector3{
				X: playerUnit.Position.X,
				Y: playerUnit.Position.Y,
				Z: newPos.Z,
			}
			if !s.wouldCollideWithUnit(playerUnit, testPosZ, state.Units) &&
				(s.Pathfinding == nil || s.isPositionWalkableWithRadius(testPosZ, playerRadius)) {
				newPos = testPosZ
			} else {
				// Can't move at all due to unit collision
				return
			}
		}
	}

	// Clamp to arena bounds
	arenaHalfSize := 95.0 // Slightly less than 100 for player radius
	newPos.X = clamp(newPos.X, -arenaHalfSize, arenaHalfSize)
	newPos.Z = clamp(newPos.Z, -arenaHalfSize, arenaHalfSize)

	playerUnit.SetPosition(newPos)
}

// isPositionWalkable checks if a position is walkable
func (s *MovementSystem) isPositionWalkable(pos types.Vector3) bool {
	if s.Pathfinding == nil || s.Pathfinding.Grid == nil {
		return true
	}
	// Use the pathfinding grid to check walkability
	gridX, gridZ := s.Pathfinding.Grid.WorldToGrid(pos.X, pos.Z)
	return s.Pathfinding.Grid.IsWalkable(gridX, gridZ)
}

// isPositionWalkableWithRadius checks if a position is walkable accounting for unit radius
// Uses direct obstacle collision detection for accuracy
func (s *MovementSystem) isPositionWalkableWithRadius(pos types.Vector3, radius float64) bool {
	if s.Pathfinding == nil {
		return true
	}

	// Use SpatialGrid for direct obstacle collision if available (more accurate)
	if s.Pathfinding.SpatialGrid != nil {
		return !s.Pathfinding.SpatialGrid.IsPositionBlocked(pos.X, pos.Z, radius)
	}

	// Fallback to grid check if no spatial grid
	if s.Pathfinding.Grid == nil {
		return true
	}

	// Check center and 4 cardinal points at the radius distance
	offsets := []types.Vector3{
		{X: 0, Y: 0, Z: 0},           // Center
		{X: radius, Y: 0, Z: 0},      // Right
		{X: -radius, Y: 0, Z: 0},     // Left
		{X: 0, Y: 0, Z: radius},      // Front
		{X: 0, Y: 0, Z: -radius},     // Back
	}

	for _, offset := range offsets {
		checkPos := types.Vector3{
			X: pos.X + offset.X,
			Y: pos.Y,
			Z: pos.Z + offset.Z,
		}
		gridX, gridZ := s.Pathfinding.Grid.WorldToGrid(checkPos.X, checkPos.Z)
		if !s.Pathfinding.Grid.IsWalkable(gridX, gridZ) {
			return false
		}
	}

	return true
}

// wouldCollideWithUnit checks if moving a unit to a new position would collide with any other unit
func (s *MovementSystem) wouldCollideWithUnit(movingUnit Unit, newPos types.Vector3, allUnits []Unit) bool {
	movingRadius := movingUnit.GetCollisionRadius()
	movingY := newPos.Y

	for _, other := range allUnits {
		// Skip self
		if other.GetID() == movingUnit.GetID() {
			continue
		}

		// Skip dead/respawning units
		if !other.IsAlive() {
			continue
		}

		// Skip units at different Y levels (airplanes vs ground units)
		otherPos := other.GetPosition()
		yDiff := math.Abs(movingY - otherPos.Y)
		if yDiff > 5.0 { // Different altitude, no collision
			continue
		}

		// Check circle-circle collision (2D, ignoring Y)
		otherRadius := other.GetCollisionRadius()
		combinedRadius := movingRadius + otherRadius

		dx := newPos.X - otherPos.X
		dz := newPos.Z - otherPos.Z
		distSq := dx*dx + dz*dz

		if distSq < combinedRadius*combinedRadius {
			return true // Would collide
		}
	}

	return false
}

// findAvoidanceDirection finds a direction to move around a blocking unit
func (s *MovementSystem) findAvoidanceDirection(movingUnit Unit, desiredDir types.Vector3, allUnits []Unit) types.Vector3 {
	pos := movingUnit.GetPosition()
	movingRadius := movingUnit.GetCollisionRadius()

	// Find the closest blocking unit
	var closestBlocker Unit
	closestDist := math.MaxFloat64

	for _, other := range allUnits {
		if other.GetID() == movingUnit.GetID() {
			continue
		}
		if !other.IsAlive() {
			continue
		}

		otherPos := other.GetPosition()
		yDiff := math.Abs(pos.Y - otherPos.Y)
		if yDiff > 5.0 {
			continue
		}

		dx := pos.X - otherPos.X
		dz := pos.Z - otherPos.Z
		dist := math.Sqrt(dx*dx + dz*dz)

		if dist < closestDist {
			closestDist = dist
			closestBlocker = other
		}
	}

	if closestBlocker == nil {
		return desiredDir
	}

	// Calculate direction to avoid the blocker
	blockerPos := closestBlocker.GetPosition()
	toBlocker := types.Vector3{
		X: blockerPos.X - pos.X,
		Y: 0,
		Z: blockerPos.Z - pos.Z,
	}

	// Perpendicular directions (left and right of blocker direction)
	perpLeft := types.Vector3{X: -toBlocker.Z, Y: 0, Z: toBlocker.X}
	perpRight := types.Vector3{X: toBlocker.Z, Y: 0, Z: -toBlocker.X}

	perpLeft = normalize(perpLeft)
	perpRight = normalize(perpRight)

	// Choose the perpendicular direction closer to our desired direction
	dotLeft := desiredDir.X*perpLeft.X + desiredDir.Z*perpLeft.Z
	dotRight := desiredDir.X*perpRight.X + desiredDir.Z*perpRight.Z

	var avoidDir types.Vector3
	if dotLeft > dotRight {
		avoidDir = perpLeft
	} else {
		avoidDir = perpRight
	}

	// Test if we can move in the avoidance direction
	testDist := movingRadius + 1.0
	testPos := types.Vector3{
		X: pos.X + avoidDir.X*testDist,
		Y: pos.Y,
		Z: pos.Z + avoidDir.Z*testDist,
	}

	if !s.wouldCollideWithUnit(movingUnit, testPos, allUnits) {
		return avoidDir
	}

	// If both directions blocked, return zero (stay put)
	return types.Vector3{X: 0, Y: 0, Z: 0}
}

// updateTankMovement handles waypoint-based movement for tanks
// Tanks will either follow a friendly tank ahead of them, or take a dynamic path to the enemy base
func (s *MovementSystem) updateTankMovement(unit Unit, state *State, deltaTime float64) {
	waypoints := unit.GetWaypoints()
	boundary := float64(types.ArenaBoundary)
	pos := unit.GetPosition()

	// Check if we should follow another friendly tank
	leaderTank := s.findLeaderTank(unit, state)

	// Calculate path if no waypoints
	if len(waypoints) == 0 {
		var targetPos types.Vector3

		if leaderTank != nil {
			// Follow the leader tank - path to a position near the leader
			leaderPos := leaderTank.GetPosition()
			// Stay a bit behind the leader (offset towards our position)
			dirToUs := normalize(subtract(pos, leaderPos))
			targetPos = types.Vector3{
				X: leaderPos.X + dirToUs.X*8.0, // Stay 8 units behind
				Y: pos.Y,
				Z: leaderPos.Z + dirToUs.Z*8.0,
			}
		} else {
			// No leader - take a dynamic path with random intermediate waypoint
			targetPos = s.getDynamicTargetPosition(unit, state)
		}

		path := s.Pathfinding.FindPath(pos, targetPos)
		if len(path) > 0 {
			// Skip the first waypoint (current position) and start from second
			if len(path) > 1 {
				path = path[1:]
			}
			unit.SetWaypoints(path)
			waypoints = path
		} else {
			// Fallback to direct movement if no path found
			s.updateDirectMovement(unit, deltaTime)
			return
		}
	}

	// If following a leader, periodically recalculate path to keep up
	if leaderTank != nil {
		currentIdx := unit.GetCurrentWaypoint()
		// Recalculate every few waypoints or when we're close to end of path
		if currentIdx > 0 && (currentIdx%3 == 0 || currentIdx >= len(waypoints)-1) {
			// Check if leader has moved significantly
			leaderPos := leaderTank.GetPosition()
			if len(waypoints) > 0 {
				lastWaypoint := waypoints[len(waypoints)-1]
				distToLeader := calculateDistance2D(lastWaypoint, leaderPos)
				if distToLeader > 15.0 {
					// Leader moved, recalculate path
					unit.ClearWaypoints()
					return
				}
			}
		}
	}

	// Get current waypoint
	currentIdx := unit.GetCurrentWaypoint()
	if currentIdx >= len(waypoints) {
		// Reached end of path, recalculate
		unit.ClearWaypoints()
		return
	}

	currentWaypoint := waypoints[currentIdx]

	// Calculate direction to current waypoint
	direction := normalize(subtract(currentWaypoint, pos))

	// Calculate movement distance
	distance := unit.GetSpeed() * deltaTime

	// Calculate new position
	newPos := types.Vector3{
		X: pos.X + direction.X*distance,
		Y: pos.Y, // Keep Y constant for tanks
		Z: pos.Z + direction.Z*distance,
	}

	// Clamp to arena boundary
	newPos.X = clamp(newPos.X, -boundary, boundary)
	newPos.Z = clamp(newPos.Z, -boundary, boundary)

	// Check collision with obstacles FIRST (using unit's collision radius)
	unitRadius := unit.GetCollisionRadius()
	if s.Pathfinding != nil && !s.isPositionWalkableWithRadius(newPos, unitRadius) {
		// Try sliding along X axis
		testPosX := types.Vector3{X: newPos.X, Y: pos.Y, Z: pos.Z}
		if s.isPositionWalkableWithRadius(testPosX, unitRadius) {
			newPos = testPosX
		} else {
			// Try sliding along Z axis
			testPosZ := types.Vector3{X: pos.X, Y: pos.Y, Z: newPos.Z}
			if s.isPositionWalkableWithRadius(testPosZ, unitRadius) {
				newPos = testPosZ
			} else {
				// Completely blocked by obstacle, recalculate path
				unit.ClearWaypoints()
				return
			}
		}
	}

	// Check collision with other units
	if s.wouldCollideWithUnit(unit, newPos, state.Units) {
		// Try to find an avoidance direction
		avoidDir := s.findAvoidanceDirection(unit, direction, state.Units)
		if avoidDir.X != 0 || avoidDir.Z != 0 {
			// Move in avoidance direction instead
			newPos = types.Vector3{
				X: pos.X + avoidDir.X*distance,
				Y: pos.Y,
				Z: pos.Z + avoidDir.Z*distance,
			}
			// Clamp avoidance position too
			newPos.X = clamp(newPos.X, -boundary, boundary)
			newPos.Z = clamp(newPos.Z, -boundary, boundary)
			// Check if avoidance position is walkable and doesn't collide
			if s.Pathfinding != nil && !s.isPositionWalkableWithRadius(newPos, unitRadius) {
				return // Can't move at all
			}
			if s.wouldCollideWithUnit(unit, newPos, state.Units) {
				return // Can't move at all
			}
		} else {
			// No avoidance possible, wait
			return
		}
	}

	// Check if we've reached the current waypoint
	distanceToWaypoint := calculateDistance2D(newPos, currentWaypoint)
	if distanceToWaypoint < 2.0 {
		// Move to next waypoint
		unit.SetCurrentWaypoint(currentIdx + 1)

		// If this was the last waypoint, clear the path
		if currentIdx+1 >= len(waypoints) {
			unit.ClearWaypoints()
		}
	}

	// Stuck detection: check if we've barely moved
	lastPos := unit.GetLastPosition()
	movedDistance := calculateDistance2D(newPos, lastPos)

	if movedDistance < 0.5 {
		// Barely moved, increment stuck counter
		stuckTicks := unit.GetStuckTicks() + 1
		unit.SetStuckTicks(stuckTicks)

		// If stuck for too long, try a different path
		if stuckTicks > 30 { // ~0.5 seconds at 60 ticks/sec
			pathAttempts := unit.GetPathAttempts() + 1
			unit.SetPathAttempts(pathAttempts)
			unit.SetStuckTicks(0)
			unit.ClearWaypoints()

			// After several failed attempts, try moving to a random nearby position
			if pathAttempts > 3 {
				// Reset path attempts and try a completely random direction
				unit.SetPathAttempts(0)
				// Apply a random nudge to get unstuck
				nudgeX := (rand.Float64() - 0.5) * 10.0
				nudgeZ := (rand.Float64() - 0.5) * 10.0
				nudgedPos := types.Vector3{
					X: clamp(newPos.X+nudgeX, -boundary, boundary),
					Y: newPos.Y,
					Z: clamp(newPos.Z+nudgeZ, -boundary, boundary),
				}
				if s.isPositionWalkableWithRadius(nudgedPos, unitRadius) {
					newPos = nudgedPos
				}
			}
		}
	} else {
		// Moving fine, reset stuck counters
		unit.SetStuckTicks(0)
		if movedDistance > 2.0 {
			// Making good progress, reset path attempts too
			unit.SetPathAttempts(0)
		}
	}

	// Update position and last position
	unit.SetPosition(newPos)
	unit.SetLastPosition(newPos)
}

// findLeaderTank finds a friendly tank ahead of this tank that it should follow
// Returns nil if no suitable leader is found
func (s *MovementSystem) findLeaderTank(unit Unit, state *State) Unit {
	pos := unit.GetPosition()
	ownerID := unit.GetOwnerID()
	targetPos := unit.GetTargetPosition() // Enemy base

	// Calculate our distance to the enemy base
	ourDistToTarget := calculateDistance2D(pos, targetPos)

	var bestLeader Unit
	bestLeaderDist := math.MaxFloat64

	for _, other := range state.Units {
		// Only consider friendly tanks
		if other.GetType() != "tank" || other.GetOwnerID() != ownerID {
			continue
		}
		// Skip self
		if other.GetID() == unit.GetID() {
			continue
		}
		// Skip dead units
		if !other.IsAlive() {
			continue
		}

		otherPos := other.GetPosition()

		// Check if this tank is closer to the enemy base (ahead of us)
		otherDistToTarget := calculateDistance2D(otherPos, targetPos)
		if otherDistToTarget >= ourDistToTarget {
			continue // Not ahead of us
		}

		// Check if this tank is within following range (not too far away)
		distToOther := calculateDistance2D(pos, otherPos)
		if distToOther > 40.0 {
			continue // Too far to follow
		}

		// Check if this tank is roughly in the direction of our target
		dirToTarget := normalize(subtract(targetPos, pos))
		dirToOther := normalize(subtract(otherPos, pos))
		dot := dirToTarget.X*dirToOther.X + dirToTarget.Z*dirToOther.Z
		if dot < 0.3 {
			continue // Not in the right direction
		}

		// Prefer the closest suitable leader
		if distToOther < bestLeaderDist {
			bestLeaderDist = distToOther
			bestLeader = other
		}
	}

	return bestLeader
}

// getDynamicTargetPosition returns a varied path target for tanks not following a leader
// This creates more interesting movement patterns instead of all tanks taking the same path
func (s *MovementSystem) getDynamicTargetPosition(unit Unit, state *State) types.Vector3 {
	pos := unit.GetPosition()
	finalTarget := unit.GetTargetPosition() // Enemy base

	// Distance to the final target
	distToTarget := calculateDistance2D(pos, finalTarget)

	// If we're close to the target, go directly there
	if distToTarget < 30.0 {
		return finalTarget
	}

	// Pick a random intermediate point to route through
	// This creates varied paths for different tanks

	// Calculate the general direction to the target
	dirToTarget := normalize(subtract(finalTarget, pos))

	// Pick a point roughly halfway to the target, with random lateral offset
	midDist := distToTarget * 0.5
	lateralOffset := (rand.Float64() - 0.5) * 40.0 // Random offset -20 to +20

	// Perpendicular direction
	perpDir := types.Vector3{X: -dirToTarget.Z, Y: 0, Z: dirToTarget.X}

	intermediatePoint := types.Vector3{
		X: pos.X + dirToTarget.X*midDist + perpDir.X*lateralOffset,
		Y: pos.Y,
		Z: pos.Z + dirToTarget.Z*midDist + perpDir.Z*lateralOffset,
	}

	// Clamp to arena bounds
	boundary := float64(types.ArenaBoundary) - 5.0
	intermediatePoint.X = clamp(intermediatePoint.X, -boundary, boundary)
	intermediatePoint.Z = clamp(intermediatePoint.Z, -boundary, boundary)

	// Check if intermediate point is walkable, if not go directly to target
	// Use tank collision radius (2.0) for the check
	if s.Pathfinding != nil && !s.isPositionWalkableWithRadius(intermediatePoint, 2.0) {
		return finalTarget
	}

	return intermediatePoint
}

// Arena boundary for perimeter patrol
var patrolCorners = []types.Vector3{
	{X: -types.ArenaBoundary, Y: types.AirplaneYPosition, Z: -types.ArenaBoundary}, // NW
	{X: types.ArenaBoundary, Y: types.AirplaneYPosition, Z: -types.ArenaBoundary},  // NE
	{X: types.ArenaBoundary, Y: types.AirplaneYPosition, Z: types.ArenaBoundary},   // SE
	{X: -types.ArenaBoundary, Y: types.AirplaneYPosition, Z: types.ArenaBoundary},  // SW
}

// updateDirectMovement handles direct line movement (for airplanes)
func (s *MovementSystem) updateDirectMovement(unit Unit, deltaTime float64) {
	pos := unit.GetPosition()
	boundary := float64(types.ArenaBoundary)

	// Check if unit is patrolling
	if unit.IsPatrolling() {
		s.updatePatrolMovement(unit, deltaTime, boundary)
		return
	}

	target := unit.GetTargetPosition()

	// Calculate direction to target
	direction := normalize(subtract(target, pos))

	// Calculate movement distance
	distance := unit.GetSpeed() * deltaTime

	// Calculate new position
	newPos := types.Vector3{
		X: pos.X + direction.X*distance,
		Y: pos.Y + direction.Y*distance,
		Z: pos.Z + direction.Z*distance,
	}

	// Check if we hit or would exceed the arena boundary
	hitBoundary := false
	if newPos.X <= -boundary || newPos.X >= boundary ||
		newPos.Z <= -boundary || newPos.Z >= boundary {
		hitBoundary = true
		// Clamp position to boundary
		newPos.X = clamp(newPos.X, -boundary, boundary)
		newPos.Z = clamp(newPos.Z, -boundary, boundary)
	}

	// Check if we've reached the target (within a small threshold)
	distanceToTarget := calculateDistance(newPos, target)
	if distanceToTarget < 1.0 {
		// Snap to target if very close
		newPos = target
	}

	// Update position
	unit.SetPosition(newPos)

	// If we hit the boundary, switch to perimeter patrol mode
	if hitBoundary {
		// Find the nearest corner to start patrolling towards
		nearestCorner := s.findNearestPatrolCorner(newPos)
		unit.SetPatrolling(true)
		unit.SetPatrolCorner(nearestCorner)
	}
}

// updatePatrolMovement handles perimeter sweep movement
func (s *MovementSystem) updatePatrolMovement(unit Unit, deltaTime float64, boundary float64) {
	pos := unit.GetPosition()
	cornerIdx := unit.GetPatrolCorner()
	target := patrolCorners[cornerIdx]

	// Calculate direction to current patrol corner
	direction := normalize(subtract(target, pos))

	// Calculate movement distance
	distance := unit.GetSpeed() * deltaTime

	// Calculate new position
	newPos := types.Vector3{
		X: pos.X + direction.X*distance,
		Y: pos.Y, // Keep Y constant for patrol
		Z: pos.Z + direction.Z*distance,
	}

	// Clamp to boundary (stay on perimeter)
	newPos.X = clamp(newPos.X, -boundary, boundary)
	newPos.Z = clamp(newPos.Z, -boundary, boundary)

	// Check if we've reached the current corner
	distanceToCorner := calculateDistance2D(newPos, target)
	if distanceToCorner < 3.0 {
		// Move to next corner (clockwise)
		nextCorner := (cornerIdx + 1) % 4
		unit.SetPatrolCorner(nextCorner)
	}

	unit.SetPosition(newPos)
}

// Helper functions for vector math

func subtract(a, b types.Vector3) types.Vector3 {
	return types.Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func normalize(v types.Vector3) types.Vector3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if length == 0 {
		return types.Vector3{X: 0, Y: 0, Z: 0}
	}
	return types.Vector3{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func calculateDistance(a, b types.Vector3) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func calculateDistance2D(a, b types.Vector3) float64 {
	dx := b.X - a.X
	dz := b.Z - a.Z
	return math.Sqrt(dx*dx + dz*dz)
}

// findNearestPatrolCorner finds the index of the nearest patrol corner to start sweeping
func (s *MovementSystem) findNearestPatrolCorner(pos types.Vector3) int {
	nearestIdx := 0
	nearestDist := math.MaxFloat64

	for i, corner := range patrolCorners {
		dist := calculateDistance2D(pos, corner)
		if dist < nearestDist {
			nearestDist = dist
			nearestIdx = i
		}
	}

	// Return the NEXT corner (so we start moving, not staying in place)
	return (nearestIdx + 1) % 4
}
