package game

import (
	"math"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

const (
	// SpatialGridCellSize is the size of each cell in the spatial grid
	SpatialGridCellSize = 10.0
)

// GridKey represents a cell coordinate in the spatial grid
type GridKey struct {
	X, Z int
}

// SpatialGrid provides efficient spatial queries for obstacles
type SpatialGrid struct {
	Cells     map[GridKey][]*Obstacle
	CellSize  float64
	Obstacles []*Obstacle
}

// NewSpatialGrid creates a new spatial grid from obstacles
func NewSpatialGrid(obstacles []*Obstacle) *SpatialGrid {
	grid := &SpatialGrid{
		Cells:     make(map[GridKey][]*Obstacle),
		CellSize:  SpatialGridCellSize,
		Obstacles: obstacles,
	}

	// Add each obstacle to the grid
	for _, obs := range obstacles {
		grid.addObstacle(obs)
	}

	return grid
}

// GetCellKey returns the grid cell key for a world position
func (g *SpatialGrid) GetCellKey(x, z float64) GridKey {
	return GridKey{
		X: int(math.Floor(x / g.CellSize)),
		Z: int(math.Floor(z / g.CellSize)),
	}
}

// addObstacle adds an obstacle to all cells it overlaps
func (g *SpatialGrid) addObstacle(obs *Obstacle) {
	// Get cells that this obstacle overlaps
	minKey := g.GetCellKey(obs.MinBounds.X, obs.MinBounds.Z)
	maxKey := g.GetCellKey(obs.MaxBounds.X, obs.MaxBounds.Z)

	for x := minKey.X; x <= maxKey.X; x++ {
		for z := minKey.Z; z <= maxKey.Z; z++ {
			key := GridKey{X: x, Z: z}
			g.Cells[key] = append(g.Cells[key], obs)
		}
	}
}

// GetObstaclesAt returns obstacles in the cell containing the given position
func (g *SpatialGrid) GetObstaclesAt(x, z float64) []*Obstacle {
	key := g.GetCellKey(x, z)
	return g.Cells[key]
}

// GetObstaclesInRadius returns all unique obstacles within radius of a position
func (g *SpatialGrid) GetObstaclesInRadius(pos types.Vector3, radius float64) []*Obstacle {
	// Determine cell range to check
	minKey := g.GetCellKey(pos.X-radius, pos.Z-radius)
	maxKey := g.GetCellKey(pos.X+radius, pos.Z+radius)

	// Use map to deduplicate
	seen := make(map[string]bool)
	result := make([]*Obstacle, 0)

	for x := minKey.X; x <= maxKey.X; x++ {
		for z := minKey.Z; z <= maxKey.Z; z++ {
			key := GridKey{X: x, Z: z}
			for _, obs := range g.Cells[key] {
				if !seen[obs.ID] {
					seen[obs.ID] = true
					result = append(result, obs)
				}
			}
		}
	}

	return result
}

// GetObstaclesAlongLine returns obstacles along a line from start to end
func (g *SpatialGrid) GetObstaclesAlongLine(start, end types.Vector3) []*Obstacle {
	// Use Bresenham-like approach to find cells along the line
	seen := make(map[string]bool)
	result := make([]*Obstacle, 0)

	dx := end.X - start.X
	dz := end.Z - start.Z
	length := math.Sqrt(dx*dx + dz*dz)

	if length == 0 {
		return g.GetObstaclesAt(start.X, start.Z)
	}

	// Step along the line
	stepSize := g.CellSize / 2 // Step half a cell at a time
	steps := int(length/stepSize) + 1

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := start.X + dx*t
		z := start.Z + dz*t

		// Get obstacles in this cell
		for _, obs := range g.GetObstaclesAt(x, z) {
			if !seen[obs.ID] {
				seen[obs.ID] = true
				result = append(result, obs)
			}
		}
	}

	return result
}

// CheckCollision checks if a circle at position with radius collides with any obstacle
func (g *SpatialGrid) CheckCollision(pos types.Vector3, radius float64) *Obstacle {
	obstacles := g.GetObstaclesInRadius(pos, radius+g.CellSize)

	for _, obs := range obstacles {
		if !obs.BlocksMovement() {
			continue
		}
		if obs.IntersectsCircleXZ(pos.X, pos.Z, radius) {
			return obs
		}
	}

	return nil
}

// IsPositionBlocked checks if a position is blocked by any obstacle
func (g *SpatialGrid) IsPositionBlocked(x, z, radius float64) bool {
	obstacles := g.GetObstaclesInRadius(types.Vector3{X: x, Y: 0, Z: z}, radius+g.CellSize)

	for _, obs := range obstacles {
		if !obs.BlocksMovement() {
			continue
		}
		if obs.IntersectsCircleXZ(x, z, radius) {
			return true
		}
	}

	return false
}
