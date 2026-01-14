package game

import (
	"math"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

const (
	// LOSRayStep is the step size for ray marching
	LOSRayStep = 1.0
)

// LOSSystem handles line-of-sight calculations
type LOSSystem struct {
	SpatialGrid *SpatialGrid
}

// NewLOSSystem creates a new line-of-sight system
func NewLOSSystem(spatialGrid *SpatialGrid) *LOSSystem {
	return &LOSSystem{
		SpatialGrid: spatialGrid,
	}
}

// HasLineOfSight checks if there's a clear line of sight between two positions
// ignoreHeight is true for units that can shoot over obstacles (helicopters)
func (l *LOSSystem) HasLineOfSight(from, to types.Vector3, ignoreHeight bool) bool {
	// Helicopters can always see (they fly over obstacles)
	if ignoreHeight {
		return true
	}

	dx := to.X - from.X
	dz := to.Z - from.Z
	length := math.Sqrt(dx*dx + dz*dz)

	if length == 0 {
		return true
	}

	// Normalize direction
	dirX := dx / length
	dirZ := dz / length

	// Step along the ray
	steps := int(length/LOSRayStep) + 1

	for i := 1; i < steps; i++ { // Start from 1 to skip the shooter's position
		t := float64(i) * LOSRayStep
		checkX := from.X + dirX*t
		checkZ := from.Z + dirZ*t

		// Check for obstacles at this point
		obstacles := l.SpatialGrid.GetObstaclesAt(checkX, checkZ)
		for _, obs := range obstacles {
			if !obs.BlocksLineOfSight() {
				continue
			}

			// Check if this point is inside the obstacle
			if obs.ContainsPointXZ(checkX, checkZ) {
				// Check height - if both shooter and target are above the obstacle, LOS is clear
				shooterHeight := from.Y
				targetHeight := to.Y
				obstacleTop := obs.MaxBounds.Y

				// If both are below the obstacle top, LOS is blocked
				if shooterHeight < obstacleTop && targetHeight < obstacleTop {
					return false
				}
			}
		}
	}

	return true
}

// HasLineOfSightBetweenUnits checks LOS between two units
func (l *LOSSystem) HasLineOfSightBetweenUnits(attacker, target Unit) bool {
	from := attacker.GetPosition()
	to := target.GetPosition()

	// Helicopters (airplanes) ignore obstacles for LOS
	ignoreHeight := attacker.GetType() == "airplane"

	return l.HasLineOfSight(from, to, ignoreHeight)
}
