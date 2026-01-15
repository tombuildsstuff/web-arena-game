package game

import (
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

// GetObstacleTypes returns all obstacles converted to types for JSON serialization
func GetObstacleTypes(obstacles []*Obstacle) []types.Obstacle {
	result := make([]types.Obstacle, len(obstacles))
	for i, obs := range obstacles {
		result[i] = obs.ToType()
	}
	return result
}
