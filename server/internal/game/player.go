package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Player represents a player in the game
type Player struct {
	ID            int
	Money         int
	BasePosition  types.Vector3
	Color         string
	ClientID      string // WebSocket client ID
	DisplayName   string // GitHub username or "Guest_XXXX"
	IsGuest       bool   // Whether this is a guest player
	Kills         int    // Total number of enemy units destroyed
	TankKills     int    // Tanks destroyed
	AirplaneKills int    // Airplanes destroyed
	TurretKills   int    // Turrets destroyed
	PlayerKills   int    // Enemy player deaths caused
}

// NewPlayer creates a new player
func NewPlayer(id int, clientID string, displayName string, isGuest bool) *Player {
	basePos := types.Base1Position
	if id == 1 {
		basePos = types.Base2Position
	}

	return &Player{
		ID:           id,
		Money:        types.StartingMoney,
		BasePosition: basePos,
		Color:        types.PlayerColors[id],
		ClientID:     clientID,
		DisplayName:  displayName,
		IsGuest:      isGuest,
	}
}

// ToType converts Player to types.Player for JSON serialization
func (p *Player) ToType() types.Player {
	return types.Player{
		ID:           p.ID,
		Money:        p.Money,
		BasePosition: p.BasePosition,
		Color:        p.Color,
		DisplayName:  p.DisplayName,
		IsGuest:      p.IsGuest,
		Kills:        p.Kills,
	}
}

// AddKill increments the player's kill count (legacy, use AddKillByType)
func (p *Player) AddKill() {
	p.Kills++
}

// AddKillByType increments the kill count for a specific unit type
func (p *Player) AddKillByType(unitType string) {
	p.Kills++
	switch unitType {
	case "tank":
		p.TankKills++
	case "airplane":
		p.AirplaneKills++
	case "turret":
		p.TurretKills++
	case "player":
		p.PlayerKills++
	}
}

// GetStats returns the player's detailed statistics
func (p *Player) GetStats() types.PlayerStats {
	// Points: 10 per tank, 20 per airplane, 20 per turret, 50 per player kill
	// Win bonus (200 points) is added separately in room.go when game ends
	points := p.TankKills*10 + p.AirplaneKills*20 + p.TurretKills*20 + p.PlayerKills*50
	return types.PlayerStats{
		TankKills:     p.TankKills,
		AirplaneKills: p.AirplaneKills,
		TurretKills:   p.TurretKills,
		PlayerKills:   p.PlayerKills,
		TotalPoints:   points,
	}
}

// CanAfford checks if the player can afford a purchase
func (p *Player) CanAfford(cost int) bool {
	return p.Money >= cost
}

// Spend deducts money from the player
func (p *Player) Spend(amount int) {
	p.Money -= amount
}

// AddMoney adds money to the player
func (p *Player) AddMoney(amount int) {
	p.Money += amount
}
