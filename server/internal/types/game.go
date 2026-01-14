package types

// Vector3 represents a 3D position
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// GameState represents the complete state of a game
type GameState struct {
	Timestamp     int64          `json:"timestamp"`
	Players       [2]Player      `json:"players"`
	Units         []Unit         `json:"units"`
	Obstacles     []Obstacle     `json:"obstacles"`
	Projectiles   []Projectile   `json:"projectiles"`
	BuyZones      []BuyZone      `json:"buyZones"`
	Turrets       []Turret       `json:"turrets"`
	PendingSpawns []PendingSpawn `json:"pendingSpawns"`
	GameStatus    string         `json:"gameStatus"` // "waiting", "playing", "finished"
	Winner        *int           `json:"winner"`
}

// Obstacle represents a static obstacle in the arena
type Obstacle struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"` // "wall", "pillar", "cover", "ramp"
	Position       Vector3 `json:"position"`
	Size           Vector3 `json:"size"`               // Width (X), Height (Y), Depth (Z)
	Rotation       float64 `json:"rotation"`           // Y-axis rotation in radians
	ElevationStart float64 `json:"elevationStart,omitempty"`
	ElevationEnd   float64 `json:"elevationEnd,omitempty"`
}

// Projectile represents a traveling projectile
type Projectile struct {
	ID        string  `json:"id"`
	ShooterID string  `json:"shooterId"`
	TargetID  string  `json:"targetId"`
	Position  Vector3 `json:"position"`
	StartPos  Vector3 `json:"startPos"`
	EndPos    Vector3 `json:"endPos"`
	Speed     float64 `json:"speed"`
	Damage    int     `json:"damage"`
	CreatedAt int64   `json:"createdAt"`
}

// Player represents a player in the game
type Player struct {
	ID           int     `json:"id"`
	Money        int     `json:"money"`
	BasePosition Vector3 `json:"basePosition"`
	Color        string  `json:"color"`
	DisplayName  string  `json:"displayName"`
	IsGuest      bool    `json:"isGuest"`
	Kills        int     `json:"kills"`
}

// BuyZone represents a location where players can purchase units
type BuyZone struct {
	ID          string  `json:"id"`
	OwnerID     int     `json:"ownerId"`    // Which player owns this buy zone (-1 = neutral)
	UnitType    string  `json:"unitType"`   // "tank" or "airplane"
	Position    Vector3 `json:"position"`
	Radius      float64 `json:"radius"`     // Interaction radius
	Cost        int     `json:"cost"`
	IsClaimable bool    `json:"isClaimable"` // Whether this zone can be claimed
}

// Unit represents a unit (tank, airplane, or player) in the game
type Unit struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"` // "tank", "airplane", or "player"
	OwnerID        int     `json:"ownerId"`
	Position       Vector3 `json:"position"`
	Health         int     `json:"health"`
	TargetPosition Vector3 `json:"targetPosition"`
	IsRespawning   bool    `json:"isRespawning,omitempty"`
	RespawnTime    float64 `json:"respawnTime,omitempty"` // Seconds remaining until respawn
}

// Turret represents a claimable turret that auto-attacks enemies
type Turret struct {
	ID               string  `json:"id"`
	Position         Vector3 `json:"position"`
	OwnerID          int     `json:"ownerId"`         // -1 = unclaimed, 0 = player 1, 1 = player 2
	DefaultOwnerID   int     `json:"defaultOwnerId"`  // Owner when respawning (-1 for middle turrets)
	Health           int     `json:"health"`
	MaxHealth        int     `json:"maxHealth"`
	IsDestroyed      bool    `json:"isDestroyed"`
	RespawnTime      float64 `json:"respawnTime"`     // Seconds remaining until respawn
	ClaimRadius      float64 `json:"claimRadius"`     // Radius for claiming
	IsTracking       bool    `json:"isTracking"`      // Whether turret is tracking a target
	TrackingProgress float64 `json:"trackingProgress"` // 0-1 progress to lock-on
}

// PendingSpawn represents a unit waiting to spawn
type PendingSpawn struct {
	UnitType  string  `json:"unitType"`  // "tank" or "airplane"
	OwnerID   int     `json:"ownerId"`   // Player who purchased the unit
	SpawnPos  Vector3 `json:"spawnPos"`  // Where the unit will spawn
	QueuedAt  int64   `json:"queuedAt"`  // When the spawn was queued (Unix millis)
	WaitTime  float64 `json:"waitTime"`  // Seconds the unit has been waiting
}
