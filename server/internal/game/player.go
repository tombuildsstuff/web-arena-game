package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Player represents a player in the game
type Player struct {
	ID           int
	Money        int
	BasePosition types.Vector3
	Color        string
	ClientID     string // WebSocket client ID
	Kills        int    // Number of enemy units destroyed
}

// NewPlayer creates a new player
func NewPlayer(id int, clientID string) *Player {
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
	}
}

// ToType converts Player to types.Player for JSON serialization
func (p *Player) ToType() types.Player {
	return types.Player{
		ID:           p.ID,
		Money:        p.Money,
		BasePosition: p.BasePosition,
		Color:        p.Color,
		Kills:        p.Kills,
	}
}

// AddKill increments the player's kill count
func (p *Player) AddKill() {
	p.Kills++
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
