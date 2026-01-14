package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// ObstacleType defines the type of obstacle
type ObstacleType string

const (
	ObstacleWall   ObstacleType = "wall"
	ObstaclePillar ObstacleType = "pillar"
	ObstacleCover  ObstacleType = "cover"
	ObstacleRamp   ObstacleType = "ramp"
)

// Obstacle represents a static obstacle in the arena
type Obstacle struct {
	ID       string       `json:"id"`
	Type     ObstacleType `json:"type"`
	Position types.Vector3 `json:"position"`
	Size     types.Vector3 `json:"size"`     // Width (X), Height (Y), Depth (Z)
	Rotation float64      `json:"rotation"` // Y-axis rotation in radians

	// For ramps only
	ElevationStart float64 `json:"elevationStart,omitempty"`
	ElevationEnd   float64 `json:"elevationEnd,omitempty"`

	// Precomputed bounding box for collision (not serialized)
	MinBounds types.Vector3 `json:"-"`
	MaxBounds types.Vector3 `json:"-"`
}

// NewObstacle creates a new obstacle with precomputed bounds
func NewObstacle(id string, obstacleType ObstacleType, pos, size types.Vector3, rotation float64) *Obstacle {
	o := &Obstacle{
		ID:       id,
		Type:     obstacleType,
		Position: pos,
		Size:     size,
		Rotation: rotation,
	}
	o.computeBounds()
	return o
}

// NewRamp creates a new ramp obstacle
func NewRamp(id string, pos, size types.Vector3, rotation, elevStart, elevEnd float64) *Obstacle {
	o := &Obstacle{
		ID:             id,
		Type:           ObstacleRamp,
		Position:       pos,
		Size:           size,
		Rotation:       rotation,
		ElevationStart: elevStart,
		ElevationEnd:   elevEnd,
	}
	o.computeBounds()
	return o
}

// computeBounds calculates the axis-aligned bounding box
// Note: For rotated obstacles, this creates a larger AABB that encompasses the rotated box
func (o *Obstacle) computeBounds() {
	halfX := o.Size.X / 2
	halfZ := o.Size.Z / 2

	// For simplicity, we use AABB that ignores rotation
	// This is conservative (larger than needed) but correct for collision detection
	o.MinBounds = types.Vector3{
		X: o.Position.X - halfX,
		Y: o.Position.Y,
		Z: o.Position.Z - halfZ,
	}
	o.MaxBounds = types.Vector3{
		X: o.Position.X + halfX,
		Y: o.Position.Y + o.Size.Y,
		Z: o.Position.Z + halfZ,
	}
}

// ContainsPoint checks if a point is inside the obstacle's bounding box
func (o *Obstacle) ContainsPoint(point types.Vector3) bool {
	return point.X >= o.MinBounds.X && point.X <= o.MaxBounds.X &&
		point.Y >= o.MinBounds.Y && point.Y <= o.MaxBounds.Y &&
		point.Z >= o.MinBounds.Z && point.Z <= o.MaxBounds.Z
}

// ContainsPointXZ checks if a point is inside the obstacle's XZ bounds (ignoring height)
func (o *Obstacle) ContainsPointXZ(x, z float64) bool {
	return x >= o.MinBounds.X && x <= o.MaxBounds.X &&
		z >= o.MinBounds.Z && z <= o.MaxBounds.Z
}

// IntersectsCircleXZ checks if a circle (unit with radius) intersects this obstacle in XZ plane
func (o *Obstacle) IntersectsCircleXZ(x, z, radius float64) bool {
	// Find the closest point on the AABB to the circle center
	closestX := clamp(x, o.MinBounds.X, o.MaxBounds.X)
	closestZ := clamp(z, o.MinBounds.Z, o.MaxBounds.Z)

	// Calculate distance from circle center to closest point
	distX := x - closestX
	distZ := z - closestZ
	distSquared := distX*distX + distZ*distZ

	return distSquared <= radius*radius
}

// BlocksLineOfSight returns whether this obstacle blocks line of sight
func (o *Obstacle) BlocksLineOfSight() bool {
	return o.Type != ObstacleRamp
}

// BlocksMovement returns whether this obstacle blocks ground unit movement
func (o *Obstacle) BlocksMovement() bool {
	return o.Type != ObstacleRamp
}

// GetElevationAt returns the ground elevation at a point on a ramp
// Returns 0 for non-ramp obstacles
func (o *Obstacle) GetElevationAt(x, z float64) float64 {
	if o.Type != ObstacleRamp {
		return 0
	}

	if !o.ContainsPointXZ(x, z) {
		return 0
	}

	// Calculate progress along the ramp (0 to 1)
	// Assuming ramp goes from MinZ to MaxZ
	progress := (z - o.MinBounds.Z) / (o.MaxBounds.Z - o.MinBounds.Z)
	progress = clamp(progress, 0, 1)

	return o.ElevationStart + (o.ElevationEnd-o.ElevationStart)*progress
}

// ToType converts Obstacle to types.Obstacle for JSON serialization
func (o *Obstacle) ToType() types.Obstacle {
	return types.Obstacle{
		ID:             o.ID,
		Type:           string(o.Type),
		Position:       o.Position,
		Size:           o.Size,
		Rotation:       o.Rotation,
		ElevationStart: o.ElevationStart,
		ElevationEnd:   o.ElevationEnd,
	}
}

// Helper function
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
