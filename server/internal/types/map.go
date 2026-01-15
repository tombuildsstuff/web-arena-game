package types

// MapDefinition contains all the configurable data for a game map
type MapDefinition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// Arena dimensions
	ArenaSize     float64 `json:"arenaSize"`
	ArenaBoundary float64 `json:"arenaBoundary"` // Units should stay within this boundary

	// Player configuration
	Players []MapPlayerConfig `json:"players"`

	// Buy zones (base zones and forward bases)
	BuyZones []MapBuyZone `json:"buyZones"`

	// Turrets
	Turrets []MapTurret `json:"turrets"`

	// Barracks (infantry spawn points)
	Barracks []MapBarracks `json:"barracks"`

	// Obstacles (walls, pillars, platforms, ramps)
	Obstacles []MapObstacle `json:"obstacles"`

	// Health pack spawn configuration
	HealthPackSpawnBounds MapBounds `json:"healthPackSpawnBounds"`
}

// MapPlayerConfig defines player-specific map settings
type MapPlayerConfig struct {
	BasePosition Vector3 `json:"basePosition"`
	SpawnOffset  Vector3 `json:"spawnOffset"` // Offset from base for unit spawning
	Color        string  `json:"color"`
}

// MapBuyZone defines a buy zone in the map
type MapBuyZone struct {
	ID            string  `json:"id"`
	DefaultOwner  int     `json:"defaultOwner"` // -1 = neutral, 0 = player 1, 1 = player 2
	UnitType      string  `json:"unitType"`     // "tank", "airplane", "super_tank", "super_helicopter", or "" for base zones
	Position      Vector3 `json:"position"`
	Radius        float64 `json:"radius"`
	Cost          int     `json:"cost"`      // Cost to buy units (0 if not a purchase zone)
	ClaimCost     int     `json:"claimCost"` // Cost to claim (0 if not claimable)
	IsClaimable   bool    `json:"isClaimable"`
	ForwardBaseID string  `json:"forwardBaseId,omitempty"` // Parent forward base ID
}

// MapTurret defines a turret in the map
type MapTurret struct {
	ID           string  `json:"id"`
	Position     Vector3 `json:"position"`
	DefaultOwner int     `json:"defaultOwner"` // -1 = neutral, 0 = player 1, 1 = player 2
}

// MapBarracks defines a barracks (infantry spawn point) in the map
type MapBarracks struct {
	ID       string  `json:"id"`
	Position Vector3 `json:"position"`
}

// MapObstacle defines an obstacle in the map
type MapObstacle struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"` // "wall", "pillar", "cover", "ramp"
	Position       Vector3 `json:"position"`
	Size           Vector3 `json:"size"`
	Rotation       float64 `json:"rotation"`
	ElevationStart float64 `json:"elevationStart,omitempty"` // For ramps
	ElevationEnd   float64 `json:"elevationEnd,omitempty"`   // For ramps
}

// MapBounds defines a rectangular area for spawning
type MapBounds struct {
	MinX float64 `json:"minX"`
	MaxX float64 `json:"maxX"`
	MinZ float64 `json:"minZ"`
	MaxZ float64 `json:"maxZ"`
}
