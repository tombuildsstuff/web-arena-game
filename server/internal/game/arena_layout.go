package game

import (
	"fmt"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// GetObstaclesFromMap creates obstacles from a map definition
func GetObstaclesFromMap(mapDef *types.MapDefinition) []*Obstacle {
	obstacles := make([]*Obstacle, 0, len(mapDef.Obstacles))

	for _, obs := range mapDef.Obstacles {
		var obstacle *Obstacle
		switch obs.Type {
		case "ramp":
			obstacle = NewRamp(obs.ID, obs.Position, obs.Size, obs.Rotation, obs.ElevationStart, obs.ElevationEnd)
		case "wall":
			obstacle = NewObstacle(obs.ID, ObstacleWall, obs.Position, obs.Size, obs.Rotation)
		case "pillar":
			obstacle = NewObstacle(obs.ID, ObstaclePillar, obs.Position, obs.Size, obs.Rotation)
		case "cover":
			obstacle = NewObstacle(obs.ID, ObstacleCover, obs.Position, obs.Size, obs.Rotation)
		default:
			obstacle = NewObstacle(obs.ID, ObstacleWall, obs.Position, obs.Size, obs.Rotation)
		}
		obstacles = append(obstacles, obstacle)
	}

	return obstacles
}

// GetSymmetricObstacles returns an arena layout with standalone walls providing cover
// Arena is 200x200 with bases at (-90, 0, 0) and (90, 0, 0)
// Layout is symmetric about the X=0 axis for fair gameplay
//
// Deprecated: Use GetObstaclesFromMap with a MapDefinition instead
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
	// HORIZONTAL LANE DIVIDERS (Z = ±35)
	// Short standalone segments with large gaps between them
	// No walls connect to vertical walls
	// ============================================================

	// Upper lane divider segments (Z = -35)
	// Segment 1: near left base
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -65, Y: 0, Z: -35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Segment 2: left-center
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -25, Y: 0, Z: -35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Segment 3: right-center
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 25, Y: 0, Z: -35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	// Segment 4: near right base
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 65, Y: 0, Z: -35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)

	// Lower lane divider segments (Z = 35) - mirrors upper
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -65, Y: 0, Z: 35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -25, Y: 0, Z: 35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 25, Y: 0, Z: 35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 65, Y: 0, Z: 35},
			types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}, 0),
	)

	// ============================================================
	// VERTICAL COVER WALLS (X = ±45)
	// Short standalone segments - NOT connected to horizontal walls
	// Positioned to avoid creating rooms
	// ============================================================

	// Left vertical segments at X = -45
	// Top lane segment
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -45, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)
	// Center segment
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -45, Y: 0, Z: 0},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	// Bottom lane segment
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -45, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)

	// Right vertical segments at X = 45 - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 45, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 45, Y: 0, Z: 0},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 45, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)

	// ============================================================
	// BASE APPROACH WALLS (X = ±70)
	// Short segments to provide cover near bases
	// ============================================================

	// Left base walls at X = -70
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -70, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -70, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// Right base walls at X = 70 - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 70, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 70, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}, 0),
	)

	// ============================================================
	// CENTRAL CORRIDOR WALLS (X = ±15)
	// Short walls in outer lanes for cover
	// ============================================================

	// Left central walls
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -15, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: -15, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)

	// Right central walls - mirrors left
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 15, Y: 0, Z: -60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleWall,
			types.Vector3{X: 15, Y: 0, Z: 60},
			types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}, 0),
	)

	// ============================================================
	// COVER PILLARS (provide cover without blocking turret sight lines)
	// Spread in clusters of 3 at reasonable distance from turrets
	// ============================================================

	pillarSize := 3.0
	pillarHeight := 6.0

	// Central area pillars (near center room turrets at ±25, 0)
	// Spread across the middle lane
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -10, Y: 0, Z: -8},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: 0},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 10, Y: 0, Z: 8},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Cover near central corridor turrets (at 0, ±50)
	// Top turret cover - 3 pillars spread in arc, ~20 units from turret
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -15, Y: 0, Z: -35},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: -30},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 15, Y: 0, Z: -35},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)
	// Bottom turret cover
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -15, Y: 0, Z: 35},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 0, Y: 0, Z: 30},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 15, Y: 0, Z: 35},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Cover near corner turrets (at ±55, ±65)
	// Top-left corner - 3 pillars spread, ~20 units from turret
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -70, Y: 0, Z: -50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -55, Y: 0, Z: -45},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -40, Y: 0, Z: -50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)
	// Top-right corner
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 40, Y: 0, Z: -50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 55, Y: 0, Z: -45},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 70, Y: 0, Z: -50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)
	// Bottom-left corner
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -70, Y: 0, Z: 50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -55, Y: 0, Z: 45},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -40, Y: 0, Z: 50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)
	// Bottom-right corner
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 40, Y: 0, Z: 50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 55, Y: 0, Z: 45},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 70, Y: 0, Z: 50},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)

	// Mid-lane cover (near turrets at ±50, ±40)
	// Left side - between base and center
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -35, Y: 0, Z: -25},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: -35, Y: 0, Z: 25},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
	)
	// Right side
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 35, Y: 0, Z: -25},
			types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}, 0),
		NewObstacle(nextID(), ObstaclePillar,
			types.Vector3{X: 35, Y: 0, Z: 25},
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

	// Mid-field cover blocks
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -55, Y: 0, Z: 0},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 55, Y: 0, Z: 0},
			types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}, 0),
	)

	// ============================================================
	// FORWARD BASE PLATFORMS (raised areas at Z = ±70)
	// These are the claimable forward bases
	// Uses flat ramps so units can walk on them
	// Enlarged to fit multiple buy zones (tank, super tank, super helicopter)
	// ============================================================

	platformHeight := 3.0
	platformSize := 20.0 // Enlarged from 12 to fit super unit buy zones

	// North forward base platform (Z = -70) - flat ramp so units can walk on it
	obstacles = append(obstacles,
		NewRamp(nextID(),
			types.Vector3{X: 0, Y: 0, Z: -70},
			types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize},
			0,
			platformHeight, // Flat - same elevation at both ends
			platformHeight,
		),
	)

	// South forward base platform (Z = 70) - flat ramp so units can walk on it
	obstacles = append(obstacles,
		NewRamp(nextID(),
			types.Vector3{X: 0, Y: 0, Z: 70},
			types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize},
			0,
			platformHeight, // Flat - same elevation at both ends
			platformHeight,
		),
	)

	// ============================================================
	// RAMPS TO FORWARD BASES
	// Ramps leading up to the raised platforms
	// ============================================================

	rampWidth := 14.0  // Wider ramp for larger platform
	rampLength := 12.0

	// North platform ramp (approaching from center, Z increasing towards -70)
	// Ramp goes from Z=-48 (ground level) to Z=-60 (platform level)
	obstacles = append(obstacles,
		NewRamp(nextID(),
			types.Vector3{X: 0, Y: 0, Z: -54}, // Center of ramp
			types.Vector3{X: rampWidth, Y: platformHeight, Z: rampLength},
			0,                 // No rotation
			platformHeight,    // ElevationStart at MinZ (-60) = 3 (top)
			0,                 // ElevationEnd at MaxZ (-48) = 0 (bottom)
		),
	)

	// South platform ramp (approaching from center, Z decreasing towards 70)
	// Ramp goes from Z=48 (ground level) to Z=60 (platform level)
	obstacles = append(obstacles,
		NewRamp(nextID(),
			types.Vector3{X: 0, Y: 0, Z: 54}, // Center of ramp
			types.Vector3{X: rampWidth, Y: platformHeight, Z: rampLength},
			0,                 // No rotation
			0,                 // ElevationStart at MinZ (48) = 0 (bottom)
			platformHeight,    // ElevationEnd at MaxZ (60) = 3 (top)
		),
	)

	// ============================================================
	// PLATFORM EDGE BARRIERS
	// Prevent units from walking off platform edges - must use ramps
	// Platform is 20x20, ramp is 14 wide, so leave 14-unit gap for ramp
	// ============================================================

	barrierHeight := 1.5 // Low barrier, just enough to block movement
	barrierThickness := 1.0

	// North platform barriers (platform at Z=-70, spans X:-10 to 10, Z:-80 to -60)
	// Left edge (X = -10)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -10, Y: platformHeight, Z: -70},
			types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}, 0),
	)
	// Right edge (X = 10)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 10, Y: platformHeight, Z: -70},
			types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}, 0),
	)
	// Back edge (Z = -80) - full width
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 0, Y: platformHeight, Z: -80},
			types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}, 0),
	)

	// South platform barriers (platform at Z=70, spans X:-10 to 10, Z:60 to 80)
	// Left edge (X = -10)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: -10, Y: platformHeight, Z: 70},
			types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}, 0),
	)
	// Right edge (X = 10)
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 10, Y: platformHeight, Z: 70},
			types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}, 0),
	)
	// Back edge (Z = 80) - full width
	obstacles = append(obstacles,
		NewObstacle(nextID(), ObstacleCover,
			types.Vector3{X: 0, Y: platformHeight, Z: 80},
			types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}, 0),
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
