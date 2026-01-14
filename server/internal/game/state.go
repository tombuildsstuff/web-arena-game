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
	GameStatus     string // "waiting", "playing", "finished"
	Winner         *int
	MatchStartTime int64 // Unix timestamp when match started
}

// NewState creates a new game state
func NewState(player1ClientID, player2ClientID string) *State {
	player1 := NewPlayer(0, player1ClientID)
	player2 := NewPlayer(1, player2ClientID)

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
		Obstacles:      GetSymmetricObstacles(),
		Projectiles:    make([]*Projectile, 0),
		BuyZones:       GetBuyZones(),
		GameStatus:     "playing",
		Winner:         nil,
		MatchStartTime: now,
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

	return types.GameState{
		Timestamp:   s.Timestamp,
		Players:     [2]types.Player{s.Players[0].ToType(), s.Players[1].ToType()},
		Units:       unitsData,
		Obstacles:   obstaclesData,
		Projectiles: projectilesData,
		BuyZones:    buyZonesData,
		GameStatus:  s.GameStatus,
		Winner:      s.Winner,
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

// GetPlayerUnit returns the player unit for the given player ID
func (s *State) GetPlayerUnit(playerID int) *PlayerUnit {
	for _, unit := range s.Units {
		if pu, ok := unit.(*PlayerUnit); ok && pu.OwnerID == playerID {
			return pu
		}
	}
	return nil
}
