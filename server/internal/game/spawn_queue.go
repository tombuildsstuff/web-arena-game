package game

import (
	"math"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// PendingSpawn represents a unit waiting to spawn
type PendingSpawn struct {
	UnitType   string        // "tank" or "airplane"
	OwnerID    int           // Player who purchased the unit
	SpawnPos   types.Vector3 // Where to spawn
	TargetPos  types.Vector3 // Target position for AI units
	QueuedAt   int64         // When the spawn was queued (Unix millis)
	ZoneID     string        // Which buy zone this came from
}

// SpawnQueue manages pending unit spawns
type SpawnQueue struct {
	Queue []*PendingSpawn
}

// NewSpawnQueue creates a new spawn queue
func NewSpawnQueue() *SpawnQueue {
	return &SpawnQueue{
		Queue: make([]*PendingSpawn, 0),
	}
}

// Add adds a pending spawn to the queue
func (q *SpawnQueue) Add(unitType string, ownerID int, spawnPos, targetPos types.Vector3, zoneID string) {
	q.Queue = append(q.Queue, &PendingSpawn{
		UnitType:  unitType,
		OwnerID:   ownerID,
		SpawnPos:  spawnPos,
		TargetPos: targetPos,
		QueuedAt:  time.Now().UnixMilli(),
		ZoneID:    zoneID,
	})
}

// GetPendingForPlayer returns the number of pending spawns for a player
func (q *SpawnQueue) GetPendingForPlayer(playerID int) int {
	count := 0
	for _, spawn := range q.Queue {
		if spawn.OwnerID == playerID {
			count++
		}
	}
	return count
}

// GetPendingByType returns pending spawns grouped by type for a player
func (q *SpawnQueue) GetPendingByType(playerID int) map[string]int {
	result := make(map[string]int)
	for _, spawn := range q.Queue {
		if spawn.OwnerID == playerID {
			result[spawn.UnitType]++
		}
	}
	return result
}

// ProcessQueue attempts to spawn units from the queue
// Returns a list of units that were successfully spawned
func (q *SpawnQueue) ProcessQueue(state *State) []Unit {
	spawnedUnits := make([]Unit, 0)
	remainingQueue := make([]*PendingSpawn, 0)

	for _, pending := range q.Queue {
		// Check if spawn position is clear
		if q.isSpawnPositionClear(pending, state.Units) {
			// Create the unit
			var unit Unit
			switch pending.UnitType {
			case "tank":
				unit = NewTank(pending.OwnerID, pending.SpawnPos, pending.TargetPos)
			case "airplane":
				unit = NewAirplane(pending.OwnerID, pending.SpawnPos, pending.TargetPos)
			}

			if unit != nil {
				spawnedUnits = append(spawnedUnits, unit)
			}
		} else {
			// Keep in queue for next tick
			remainingQueue = append(remainingQueue, pending)
		}
	}

	q.Queue = remainingQueue
	return spawnedUnits
}

// isSpawnPositionClear checks if a spawn position is free of other units
func (q *SpawnQueue) isSpawnPositionClear(pending *PendingSpawn, units []Unit) bool {
	// Get the collision radius for the unit type
	var spawnRadius float64
	var spawnY float64
	switch pending.UnitType {
	case "tank":
		spawnRadius = types.TankCollisionRadius
		spawnY = types.TankYPosition
	case "airplane":
		spawnRadius = types.AirplaneCollisionRadius
		spawnY = types.AirplaneYPosition
	default:
		spawnRadius = 2.0
		spawnY = 1.0
	}

	// Check against all existing units
	for _, unit := range units {
		if !unit.IsAlive() {
			continue
		}

		unitPos := unit.GetPosition()

		// Skip units at different Y levels (airplanes vs ground units)
		yDiff := math.Abs(spawnY - unitPos.Y)
		if yDiff > 5.0 {
			continue
		}

		// Check circle-circle collision
		otherRadius := unit.GetCollisionRadius()
		combinedRadius := spawnRadius + otherRadius

		dx := pending.SpawnPos.X - unitPos.X
		dz := pending.SpawnPos.Z - unitPos.Z
		distSq := dx*dx + dz*dz

		// Add a small buffer (1.0 units) to ensure units have room to move
		bufferRadius := combinedRadius + 1.0
		if distSq < bufferRadius*bufferRadius {
			return false // Position is blocked
		}
	}

	// Also check against other pending spawns at the same location
	for _, other := range q.Queue {
		if other == pending {
			continue
		}

		// Skip different Y levels
		var otherY float64
		var otherRadius float64
		switch other.UnitType {
		case "tank":
			otherY = types.TankYPosition
			otherRadius = types.TankCollisionRadius
		case "airplane":
			otherY = types.AirplaneYPosition
			otherRadius = types.AirplaneCollisionRadius
		default:
			otherY = 1.0
			otherRadius = 2.0
		}

		yDiff := math.Abs(spawnY - otherY)
		if yDiff > 5.0 {
			continue
		}

		// Check if spawning at same position
		dx := pending.SpawnPos.X - other.SpawnPos.X
		dz := pending.SpawnPos.Z - other.SpawnPos.Z
		distSq := dx*dx + dz*dz

		combinedRadius := spawnRadius + otherRadius + 1.0
		if distSq < combinedRadius*combinedRadius {
			// Same position - only allow the one queued first
			if other.QueuedAt < pending.QueuedAt {
				return false // Other spawn has priority
			}
		}
	}

	return true
}

// ToTypes converts pending spawns to types for JSON serialization
func (q *SpawnQueue) ToTypes() []types.PendingSpawn {
	result := make([]types.PendingSpawn, len(q.Queue))
	now := time.Now().UnixMilli()
	for i, spawn := range q.Queue {
		result[i] = types.PendingSpawn{
			UnitType:  spawn.UnitType,
			OwnerID:   spawn.OwnerID,
			SpawnPos:  spawn.SpawnPos,
			QueuedAt:  spawn.QueuedAt,
			WaitTime:  float64(now-spawn.QueuedAt) / 1000.0, // Seconds waiting
		}
	}
	return result
}
