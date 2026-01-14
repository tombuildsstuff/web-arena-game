package game

import (
	"fmt"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// GetSymmetricObstacles returns a complex arena layout with many rooms and corridors
// Arena is 200x200 with bases at (-90, 0, 0) and (90, 0, 0)
// Layout is symmetric about the X=0 axis for fair gameplay
//
// All doorways/gaps are at least 1.5x tank diameter (6 units) to ensure passage
func GetSymmetricObstacles() []*Obstacle {
	obstacles := make([]*Obstacle, 0)
	idCounter := 0

	nextID := func() string {
		idCounter++
		return fmt.Sprintf("obs_%d", idCounter)
	}

	// Constants
	wallHeight := 6.0
	wallThickness := 2.0

	// Minimum gap calculation:
	// - Tank collision radius: 2.0 (diameter 4.0)
	// - Pathfinding buffer: 2.0 (cells within 2.0 of obstacles are unwalkable)
	// - Required gap: (2.0 + 2.0) * 2 = 8.0 minimum
	// - Using 10.0 for comfortable 1.5x passage
	minGap := 10.0

	// ============================================================
	// OUTER WALLS - Define the playable area boundary
	// ============================================================

	// North outer wall (full width)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 0, Y: 0, Z: -85},
			types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}, 0),
	)

	// South outer wall (full width)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 0, Y: 0, Z: 85},
			types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}, 0),
	)

	// ============================================================
	// HORIZONTAL LANE DIVIDERS (Z = ±30)
	// Creates 3 lanes: top, middle, bottom
	// Each wall section leaves gaps of minGap (10 units) for doorways
	// ============================================================

	// Upper lane divider (Z = -30)
	// Layout: [wall -85 to -55] [gap 10] [wall -45 to -5] [gap 10] [wall 5 to 45] [gap 10] [wall 55 to 85]
	//
	// Left outer section: center at -70, width 30 -> edges at -85 to -55
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -70, Y: 0, Z: -30},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Left inner section: center at -25, width 40 -> edges at -45 to -5
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -25, Y: 0, Z: -30},
			types.Vector3{X: 40, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Right inner section: center at 25, width 40 -> edges at 5 to 45
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 25, Y: 0, Z: -30},
			types.Vector3{X: 40, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Right outer section: center at 70, width 30 -> edges at 55 to 85
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 70, Y: 0, Z: -30},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Lower lane divider (Z = 30) - mirrors upper
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -70, Y: 0, Z: 30},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -25, Y: 0, Z: 30},
			types.Vector3{X: 40, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 25, Y: 0, Z: 30},
			types.Vector3{X: 40, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 70, Y: 0, Z: 30},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
	)

	// ============================================================
	// VERTICAL ROOM DIVIDERS (X = ±50)
	// Creates side rooms separated from center area
	// Doorways at Z = -40, 0, +40 (each 10 units wide)
	// ============================================================

	// Left vertical divider at X = -50
	// Top section: Z from -85 to -45 (height 40), leaves 10-unit gap at Z=-40
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -50, Y: 0, Z: -65},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 40}, 0),
	)
	// Upper-middle section: Z from -35 to -5 (height 30), 10-unit gaps at Z=-40 and Z=0
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -50, Y: 0, Z: -20},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 30}, 0),
	)
	// Lower-middle section: Z from 5 to 35 (height 30), 10-unit gaps at Z=0 and Z=40
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -50, Y: 0, Z: 20},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 30}, 0),
	)
	// Bottom section: Z from 45 to 85 (height 40)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -50, Y: 0, Z: 65},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 40}, 0),
	)

	// Right vertical divider at X = 50 - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 50, Y: 0, Z: -65},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 40}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 50, Y: 0, Z: -20},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 30}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 50, Y: 0, Z: 20},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 30}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 50, Y: 0, Z: 65},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 40}, 0),
	)

	// ============================================================
	// BASE APPROACH WALLS (X = ±75)
	// Smaller walls near bases - just outer sections to guide traffic
	// Wide center passage for easy tank access
	// ============================================================

	// Left base walls at X = -75
	// Top section only: Z from -85 to -50 (height 35)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -75, Y: 0, Z: -67.5},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 35}, 0),
	)
	// Bottom section only: Z from 50 to 85 (height 35)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -75, Y: 0, Z: 67.5},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 35}, 0),
	)

	// Right base walls at X = 75 - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 75, Y: 0, Z: -67.5},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 35}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 75, Y: 0, Z: 67.5},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 35}, 0),
	)

	// ============================================================
	// CENTRAL CORRIDOR WALLS (X = ±15)
	// Short walls in outer lanes only, leaving center completely open
	// ============================================================

	// Left central walls - only in top and bottom lanes
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -15, Y: 0, Z: -57},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -15, Y: 0, Z: 57},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// Right central walls - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 15, Y: 0, Z: -57},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 15, Y: 0, Z: 57},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// ============================================================
	// COVER PILLARS (provide cover without blocking movement)
	// All pillars are small enough (3x3) to leave minGap around them
	// ============================================================

	pillarSize := 3.0
	pillarHeight := 6.0

	// Central area pillars
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: -15},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: 15},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Mid-field pillars
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -35, Y: 0, Z: 0},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 35, Y: 0, Z: 0},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Lane pillars (top and bottom lanes)
	obstacles = append(obstacles,
		// Top lane
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -35, Y: 0, Z: -55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: -55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 35, Y: 0, Z: -55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		// Bottom lane
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -35, Y: 0, Z: 55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: 55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 35, Y: 0, Z: 55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Side room pillars
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -65, Y: 0, Z: -55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -65, Y: 0, Z: 0},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -65, Y: 0, Z: 55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 65, Y: 0, Z: -55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 65, Y: 0, Z: 0},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 65, Y: 0, Z: 55},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// ============================================================
	// COVER BLOCKS (low walls for tactical cover)
	// ============================================================

	coverHeight := 3.0
	coverSize := 4.0

	// Near-base cover
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -82, Y: 0, Z: -40},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -82, Y: 0, Z: 40},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 82, Y: 0, Z: -40},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 82, Y: 0, Z: 40},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
	)

	// Suppress unused variable warning
	_ = minGap

	return obstacles
}

// GetObstacleTypes returns all obstacles converted to types for JSON serialization
func GetObstacleTypes(obstacles []*Obstacle) []types.Obstacle {
	result := make([]types.Obstacle, len(obstacles))
	for i, obs := range obstacles {
		result[i] = obs.ToType()
	}
	return result
}
