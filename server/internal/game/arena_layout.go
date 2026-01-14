package game

import (
	"fmt"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// GetSymmetricObstacles returns a layout with 10 rooms around a central courtyard
// Arena is 200x200 with bases at (-90, 0, 0) and (90, 0, 0)
// Layout is symmetric about the X=0 axis for fair gameplay
//
// Layout overview (top-down view):
//
//	+-------+-------+-------+-------+-------+
//	|       |       |       |       |       |
//	| Room1 | Room2 |       | Room4 | Room5 |
//	|  (P1) |       | Court |       |  (P2) |
//	+---  --+---  --+   Y   +--  ---+--  ---+
//	|       |       |   A   |       |       |
//	| Room6 | Room7 |   R   | Room8 | Room9 |
//	|       |       |   D   |       |       |
//	+-------+---  --+--   --+--  ---+-------+
//	        |       |       |       |
//	        |Room10 |       |Room10 |
//	        |       |       |       |
//	        +-------+-------+-------+
//
func GetSymmetricObstacles() []*Obstacle {
	obstacles := make([]*Obstacle, 0)
	idCounter := 0

	nextID := func() string {
		idCounter++
		return fmt.Sprintf("obs_%d", idCounter)
	}

	wallHeight := 6.0
	wallThickness := 2.0

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
	// CENTRAL COURTYARD WALLS
	// ============================================================

	// The central courtyard is roughly 40x60 in the center
	courtyardHalfWidth := 20.0
	courtyardHalfDepth := 40.0

	// North courtyard wall (with opening in center)
	obstacles = append(obstacles,
		// Left section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -courtyardHalfWidth/2 - 5, Y: 0, Z: -courtyardHalfDepth},
			types.Vector3{X: courtyardHalfWidth - 10, Y: wallHeight, Z: wallThickness}, 0),
		// Right section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: courtyardHalfWidth/2 + 5, Y: 0, Z: -courtyardHalfDepth},
			types.Vector3{X: courtyardHalfWidth - 10, Y: wallHeight, Z: wallThickness}, 0),
	)

	// South courtyard wall (with opening in center)
	obstacles = append(obstacles,
		// Left section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -courtyardHalfWidth/2 - 5, Y: 0, Z: courtyardHalfDepth},
			types.Vector3{X: courtyardHalfWidth - 10, Y: wallHeight, Z: wallThickness}, 0),
		// Right section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: courtyardHalfWidth/2 + 5, Y: 0, Z: courtyardHalfDepth},
			types.Vector3{X: courtyardHalfWidth - 10, Y: wallHeight, Z: wallThickness}, 0),
	)

	// West courtyard wall (with openings for room connections)
	obstacles = append(obstacles,
		// Top section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -courtyardHalfWidth, Y: 0, Z: -courtyardHalfDepth/2 - 10},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: courtyardHalfDepth/2 - 5}, 0),
		// Bottom section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -courtyardHalfWidth, Y: 0, Z: courtyardHalfDepth/2 + 10},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: courtyardHalfDepth/2 - 5}, 0),
	)

	// East courtyard wall (with openings for room connections)
	obstacles = append(obstacles,
		// Top section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: courtyardHalfWidth, Y: 0, Z: -courtyardHalfDepth/2 - 10},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: courtyardHalfDepth/2 - 5}, 0),
		// Bottom section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: courtyardHalfWidth, Y: 0, Z: courtyardHalfDepth/2 + 10},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: courtyardHalfDepth/2 - 5}, 0),
	)

	// ============================================================
	// LEFT SIDE ROOMS (Player 1 side) - 5 rooms
	// ============================================================

	// Room divider positions on left side
	leftRoomX := -55.0 // Center X of left rooms area

	// Horizontal wall dividing top and bottom rows (with doorway)
	obstacles = append(obstacles,
		// Left section of horizontal divider
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -70, Y: 0, Z: 0},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
		// Right section (near courtyard)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -30, Y: 0, Z: 0},
			types.Vector3{X: 15, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Horizontal wall in front of Player 1 base (provides cover for spawned units)
	// Runs east-west with openings at top and bottom for passage
	obstacles = append(obstacles,
		// Center section (main cover wall)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: leftRoomX, Y: 0, Z: 0},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
		// Top section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: leftRoomX, Y: 0, Z: -50},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
		// Bottom section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: leftRoomX, Y: 0, Z: 50},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Wall between left rooms and courtyard approach
	obstacles = append(obstacles,
		// Top diagonal approach wall
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -35, Y: 0, Z: -55},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}, 0),
		// Bottom diagonal approach wall
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -35, Y: 0, Z: 55},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}, 0),
	)

	// ============================================================
	// RIGHT SIDE ROOMS (Player 2 side) - 5 rooms (symmetric)
	// ============================================================

	rightRoomX := 55.0

	// Horizontal wall dividing top and bottom rows (with doorway)
	obstacles = append(obstacles,
		// Right section of horizontal divider
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 70, Y: 0, Z: 0},
			types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}, 0),
		// Left section (near courtyard)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 30, Y: 0, Z: 0},
			types.Vector3{X: 15, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Horizontal wall in front of Player 2 base (provides cover for spawned units)
	// Runs east-west with openings at top and bottom for passage
	obstacles = append(obstacles,
		// Center section (main cover wall)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: rightRoomX, Y: 0, Z: 0},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
		// Top section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: rightRoomX, Y: 0, Z: -50},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
		// Bottom section
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: rightRoomX, Y: 0, Z: 50},
			types.Vector3{X: 35, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Wall between right rooms and courtyard approach
	obstacles = append(obstacles,
		// Top approach wall
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 35, Y: 0, Z: -55},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}, 0),
		// Bottom approach wall
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 35, Y: 0, Z: 55},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}, 0),
	)

	// ============================================================
	// BOTTOM ROOMS (Room 10 areas - symmetric on both sides)
	// ============================================================

	// These are the rooms at the south end, connecting both sides

	// Left bottom room walls
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -55, Y: 0, Z: 70},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// Right bottom room walls
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 55, Y: 0, Z: 70},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// ============================================================
	// TOP ROOMS EXTENSION
	// ============================================================

	// Left top room walls
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -55, Y: 0, Z: -70},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// Right top room walls
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 55, Y: 0, Z: -70},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// ============================================================
	// COVER OBJECTS IN ROOMS AND COURTYARD
	// ============================================================

	coverHeight := 3.0
	pillarHeight := 6.0

	// Central courtyard pillars (4 corners)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -12, Y: 0, Z: -25},
			types.Vector3{X: 4, Y: pillarHeight, Z: 4}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 12, Y: 0, Z: -25},
			types.Vector3{X: 4, Y: pillarHeight, Z: 4}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -12, Y: 0, Z: 25},
			types.Vector3{X: 4, Y: pillarHeight, Z: 4}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 12, Y: 0, Z: 25},
			types.Vector3{X: 4, Y: pillarHeight, Z: 4}, 0),
	)

	// Cover in left side rooms
	obstacles = append(obstacles,
		// Room 1 (top-left, near P1 base)
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -75, Y: 0, Z: -45},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 2
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -45, Y: 0, Z: -50},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 6 (bottom-left)
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -75, Y: 0, Z: 45},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 7
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -45, Y: 0, Z: 50},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
	)

	// Cover in right side rooms (symmetric)
	obstacles = append(obstacles,
		// Room 5 (top-right, near P2 base)
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 75, Y: 0, Z: -45},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 4
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 45, Y: 0, Z: -50},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 9 (bottom-right)
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 75, Y: 0, Z: 45},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		// Room 8
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 45, Y: 0, Z: 50},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
	)

	// Cover in bottom rooms
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -30, Y: 0, Z: 70},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 30, Y: 0, Z: 70},
			types.Vector3{X: 5, Y: coverHeight, Z: 5}, 0),
	)

	// ============================================================
	// BASE PERIMETER WALLS
	// Each base has walls on 3 sides (back, north, south) with open front
	// ============================================================

	baseWallHeight := 4.0
	baseWallThickness := 2.0
	baseSize := 25.0 // Size of the base enclosure

	// Player 1 base perimeter (at X=-90, open front facing the arena)
	obstacles = append(obstacles,
		// Back wall (west side, vertical running north-south)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -90 - baseSize/2, Y: 0, Z: 0},
			types.Vector3{X: baseWallThickness, Y: baseWallHeight, Z: baseSize}, 0),
		// North wall (horizontal running east-west)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -90, Y: 0, Z: -baseSize / 2},
			types.Vector3{X: baseSize, Y: baseWallHeight, Z: baseWallThickness}, 0),
		// South wall (horizontal running east-west)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -90, Y: 0, Z: baseSize / 2},
			types.Vector3{X: baseSize, Y: baseWallHeight, Z: baseWallThickness}, 0),
	)

	// Player 2 base perimeter (at X=90, open front facing the arena)
	obstacles = append(obstacles,
		// Back wall (east side, vertical running north-south)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 90 + baseSize/2, Y: 0, Z: 0},
			types.Vector3{X: baseWallThickness, Y: baseWallHeight, Z: baseSize}, 0),
		// North wall (horizontal running east-west)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 90, Y: 0, Z: -baseSize / 2},
			types.Vector3{X: baseSize, Y: baseWallHeight, Z: baseWallThickness}, 0),
		// South wall (horizontal running east-west)
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 90, Y: 0, Z: baseSize / 2},
			types.Vector3{X: baseSize, Y: baseWallHeight, Z: baseWallThickness}, 0),
	)

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
