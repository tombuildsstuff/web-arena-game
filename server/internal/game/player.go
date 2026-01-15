package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Player represents a player in the game
type Player struct {
	ID                  int
	Money               int
	BasePosition        types.Vector3
	Color               string
	ClientID            string // WebSocket client ID
	DisplayName         string // GitHub username or "Guest_XXXX"
	IsGuest             bool   // Whether this is a guest player
	Kills               int    // Total number of enemy units destroyed
	TankKills           int    // Tanks destroyed
	AirplaneKills       int    // Airplanes destroyed
	TurretKills         int    // Turrets destroyed
	PlayerKills         int    // Enemy player deaths caused
	SniperKills         int    // Snipers destroyed
	RocketLauncherKills int    // Rocket launchers destroyed
	BarracksKills       int    // Barracks destroyed
}

// NewPlayerWithMap creates a new player using map configuration
func NewPlayerWithMap(id int, clientID string, displayName string, isGuest bool, mapDef *types.MapDefinition) *Player {
	playerConfig := mapDef.Players[id]
	return &Player{
		ID:           id,
		Money:        types.StartingMoney,
		BasePosition: playerConfig.BasePosition,
		Color:        playerConfig.Color,
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

// AddKillByType increments the kill count for a specific unit type and awards money
func (p *Player) AddKillByType(unitType string) {
	p.Kills++
	switch unitType {
	case "tank":
		p.TankKills++
		p.Money += types.KillRewardUnit // 10 money per unit kill
	case "super_tank":
		p.TankKills++
		p.Money += types.KillRewardUnit * 2 // 20 money for super units (double)
	case "airplane":
		p.AirplaneKills++
		p.Money += types.KillRewardUnit // 10 money per unit kill
	case "super_helicopter":
		p.AirplaneKills++
		p.Money += types.KillRewardUnit * 2 // 20 money for super units (double)
	case "sniper":
		p.SniperKills++
		p.Money += types.KillRewardUnit // 10 money per sniper kill
	case "rocket_launcher":
		p.RocketLauncherKills++
		p.Money += types.KillRewardUnit // 10 money per rocket launcher kill
	case "turret":
		p.TurretKills++
		p.Money += types.KillRewardUnit // 10 money per turret destroyed
	case "barracks":
		p.BarracksKills++
		p.Money += types.KillRewardUnit // 10 money per barracks destroyed
	case "player":
		p.PlayerKills++
		p.Money += types.KillRewardUnit // 10 money per player kill
	}
}

// GetStats returns the player's detailed statistics
func (p *Player) GetStats() types.PlayerStats {
	// Points: 10 per tank, 20 per airplane, 15 per infantry, 20 per turret, 25 per barracks, 50 per player kill
	// Win bonus (200 points) is added separately in room.go when game ends
	points := p.TankKills*10 +
		p.AirplaneKills*20 +
		p.SniperKills*15 +
		p.RocketLauncherKills*15 +
		p.TurretKills*20 +
		p.BarracksKills*25 +
		p.PlayerKills*50
	return types.PlayerStats{
		TankKills:           p.TankKills,
		AirplaneKills:       p.AirplaneKills,
		SniperKills:         p.SniperKills,
		RocketLauncherKills: p.RocketLauncherKills,
		TurretKills:         p.TurretKills,
		BarracksKills:       p.BarracksKills,
		PlayerKills:         p.PlayerKills,
		TotalPoints:         points,
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
