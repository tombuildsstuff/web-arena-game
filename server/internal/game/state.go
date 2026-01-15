package game

import (
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// State represents the game state
type State struct {
	Timestamp      int64
	Players        [2]*Player
	Units          []Unit
	Obstacles      []*Obstacle
	Projectiles    []*Projectile
	BuyZones       []*BuyZone
	Turrets        []*Turret
	HealthPacks    []*HealthPack
	SpawnQueue     *SpawnQueue
	GameStatus     string // "waiting", "playing", "finished"
	Winner         *int
	MatchStartTime int64 // Unix timestamp when match started
	MapDefinition  *types.MapDefinition
}

// NewStateWithMap creates a new game state using a map definition
func NewStateWithMap(mapDef *types.MapDefinition, player1ClientID, player1DisplayName string, player1IsGuest bool, player2ClientID, player2DisplayName string, player2IsGuest bool) *State {
	player1 := NewPlayerWithMap(0, player1ClientID, player1DisplayName, player1IsGuest, mapDef)
	player2 := NewPlayerWithMap(1, player2ClientID, player2DisplayName, player2IsGuest, mapDef)

	// Create player units at their bases
	playerUnit1 := NewPlayerUnit(0, player1.BasePosition)
	playerUnit2 := NewPlayerUnit(1, player2.BasePosition)

	now := time.Now().UnixMilli()
	return &State{
		Timestamp: now,
		Players: [2]*Player{
			player1,
			player2,
		},
		Units:          []Unit{playerUnit1, playerUnit2},
		Obstacles:      GetObstaclesFromMap(mapDef),
		Projectiles:    make([]*Projectile, 0),
		BuyZones:       GetBuyZonesFromMap(mapDef),
		Turrets:        GetTurretsFromMap(mapDef),
		HealthPacks:    make([]*HealthPack, 0),
		SpawnQueue:     NewSpawnQueue(),
		GameStatus:     "playing",
		Winner:         nil,
		MatchStartTime: now,
		MapDefinition:  mapDef,
	}
}

// ToType converts State to types.GameState for JSON serialization
func (s *State) ToType() types.GameState {
	unitsData := make([]types.Unit, len(s.Units))
	for i, unit := range s.Units {
		unitsData[i] = unit.ToType()
	}

	obstaclesData := make([]types.Obstacle, len(s.Obstacles))
	for i, obs := range s.Obstacles {
		obstaclesData[i] = obs.ToType()
	}

	projectilesData := make([]types.Projectile, len(s.Projectiles))
	for i, proj := range s.Projectiles {
		projectilesData[i] = proj.ToType()
	}

	buyZonesData := make([]types.BuyZone, len(s.BuyZones))
	for i, zone := range s.BuyZones {
		buyZonesData[i] = zone.ToType()
	}

	turretsData := make([]types.Turret, len(s.Turrets))
	for i, turret := range s.Turrets {
		turretsData[i] = turret.ToType()
	}

	healthPacksData := make([]types.HealthPack, len(s.HealthPacks))
	for i, pack := range s.HealthPacks {
		healthPacksData[i] = pack.ToType()
	}

	var pendingSpawnsData []types.PendingSpawn
	if s.SpawnQueue != nil {
		pendingSpawnsData = s.SpawnQueue.ToTypes()
	} else {
		pendingSpawnsData = make([]types.PendingSpawn, 0)
	}

	return types.GameState{
		Timestamp:     s.Timestamp,
		Players:       [2]types.Player{s.Players[0].ToType(), s.Players[1].ToType()},
		Units:         unitsData,
		Obstacles:     obstaclesData,
		Projectiles:   projectilesData,
		BuyZones:      buyZonesData,
		Turrets:       turretsData,
		HealthPacks:   healthPacksData,
		PendingSpawns: pendingSpawnsData,
		GameStatus:    s.GameStatus,
		Winner:        s.Winner,
	}
}

// AddUnit adds a unit to the game state
func (s *State) AddUnit(unit Unit) {
	s.Units = append(s.Units, unit)
}

// RemoveUnit removes a unit by ID
func (s *State) RemoveUnit(id string) {
	for i, unit := range s.Units {
		if unit.GetID() == id {
			s.Units = append(s.Units[:i], s.Units[i+1:]...)
			return
		}
	}
}

// GetPlayer returns a player by ID
func (s *State) GetPlayer(id int) *Player {
	if id >= 0 && id < len(s.Players) {
		return s.Players[id]
	}
	return nil
}

// UpdateTimestamp updates the state timestamp
func (s *State) UpdateTimestamp() {
	s.Timestamp = time.Now().UnixMilli()
}

// GetMatchDuration returns the match duration in seconds
func (s *State) GetMatchDuration() int {
	return int((time.Now().UnixMilli() - s.MatchStartTime) / 1000)
}

// AddProjectile adds a projectile to the game state
func (s *State) AddProjectile(proj *Projectile) {
	s.Projectiles = append(s.Projectiles, proj)
}

// RemoveProjectile removes a projectile by ID
func (s *State) RemoveProjectile(id string) {
	for i, proj := range s.Projectiles {
		if proj.ID == id {
			s.Projectiles = append(s.Projectiles[:i], s.Projectiles[i+1:]...)
			return
		}
	}
}

// RemoveProjectiles removes multiple projectiles by ID
func (s *State) RemoveProjectiles(ids []string) {
	idSet := make(map[string]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	remaining := make([]*Projectile, 0, len(s.Projectiles))
	for _, proj := range s.Projectiles {
		if !idSet[proj.ID] {
			remaining = append(remaining, proj)
		}
	}
	s.Projectiles = remaining
}

// GetUnitByID returns a unit by ID
func (s *State) GetUnitByID(id string) Unit {
	for _, unit := range s.Units {
		if unit.GetID() == id {
			return unit
		}
	}
	return nil
}

// GetObstacles returns all obstacles
func (s *State) GetObstacles() []*Obstacle {
	return s.Obstacles
}

// GetElevationAt returns the ground elevation at a given position
// Checks all ramp obstacles and returns the highest elevation
func (s *State) GetElevationAt(x, z float64) float64 {
	maxElevation := 0.0
	for _, obs := range s.Obstacles {
		elevation := obs.GetElevationAt(x, z)
		if elevation > maxElevation {
			maxElevation = elevation
		}
	}
	return maxElevation
}

// GetPlayerUnit returns the player unit for the given player ID
func (s *State) GetPlayerUnit(playerID int) *PlayerUnit {
	for _, unit := range s.Units {
		if pu, ok := unit.(*PlayerUnit); ok && pu.OwnerID == playerID {
			return pu
		}
	}
	return nil
}

// GetTurretByID returns a turret by ID
func (s *State) GetTurretByID(id string) *Turret {
	for _, turret := range s.Turrets {
		if turret.ID == id {
			return turret
		}
	}
	return nil
}

// GetTurrets returns all turrets
func (s *State) GetTurrets() []*Turret {
	return s.Turrets
}

// CountPlayerUnitsOfType counts how many units of a specific type a player has (alive)
func (s *State) CountPlayerUnitsOfType(playerID int, unitType string) int {
	count := 0
	for _, unit := range s.Units {
		if unit.GetOwnerID() == playerID && unit.GetType() == unitType && unit.IsAlive() {
			count++
		}
	}
	return count
}

// CountPlayerPendingUnitsOfType counts how many units of a specific type a player has pending in spawn queue
func (s *State) CountPlayerPendingUnitsOfType(playerID int, unitType string) int {
	if s.SpawnQueue == nil {
		return 0
	}
	count := 0
	for _, spawn := range s.SpawnQueue.Queue {
		if spawn.OwnerID == playerID && spawn.UnitType == unitType {
			count++
		}
	}
	return count
}

// HasSuperUnit checks if a player already has a super unit of the given type (alive or pending)
func (s *State) HasSuperUnit(playerID int, unitType string) bool {
	// Check alive units
	if s.CountPlayerUnitsOfType(playerID, unitType) > 0 {
		return true
	}
	// Check pending spawns
	if s.CountPlayerPendingUnitsOfType(playerID, unitType) > 0 {
		return true
	}
	return false
}

// RemoveHealthPack removes a health pack by ID
func (s *State) RemoveHealthPack(id string) {
	for i, pack := range s.HealthPacks {
		if pack.ID == id {
			s.HealthPacks = append(s.HealthPacks[:i], s.HealthPacks[i+1:]...)
			return
		}
	}
}
