package game

import (
	"math"

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
			s.updateTankMovement(unit, deltaTime)
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

	// Check collision with obstacles using pathfinding system
	if s.Pathfinding != nil && !s.isPositionWalkable(newPos) {
		// Try sliding along walls
		// Try X movement only
		testPosX := types.Vector3{
			X: newPos.X,
			Y: playerUnit.Position.Y,
			Z: playerUnit.Position.Z,
		}
		if s.isPositionWalkable(testPosX) {
			newPos = testPosX
		} else {
			// Try Z movement only
			testPosZ := types.Vector3{
				X: playerUnit.Position.X,
				Y: playerUnit.Position.Y,
				Z: newPos.Z,
			}
			if s.isPositionWalkable(testPosZ) {
				newPos = testPosZ
			} else {
				// Can't move at all
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

// updateTankMovement handles waypoint-based movement for tanks
func (s *MovementSystem) updateTankMovement(unit Unit, deltaTime float64) {
	waypoints := unit.GetWaypoints()

	// Calculate path if no waypoints
	if len(waypoints) == 0 {
		path := s.Pathfinding.FindPath(unit.GetPosition(), unit.GetTargetPosition())
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

	// Get current waypoint
	currentIdx := unit.GetCurrentWaypoint()
	if currentIdx >= len(waypoints) {
		// Reached end of path, recalculate
		unit.ClearWaypoints()
		return
	}

	currentWaypoint := waypoints[currentIdx]
	pos := unit.GetPosition()

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

	// Update position
	unit.SetPosition(newPos)
}

// updateDirectMovement handles direct line movement (for airplanes)
func (s *MovementSystem) updateDirectMovement(unit Unit, deltaTime float64) {
	pos := unit.GetPosition()
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

	// Check if we've reached the target (within a small threshold)
	distanceToTarget := calculateDistance(newPos, target)
	if distanceToTarget < 1.0 {
		// Snap to target if very close
		newPos = target
	}

	// Update position
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
