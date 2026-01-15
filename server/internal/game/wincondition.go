package game

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// WinConditionSystem checks for win conditions
type WinConditionSystem struct{}

// NewWinConditionSystem creates a new win condition system
func NewWinConditionSystem() *WinConditionSystem {
	return &WinConditionSystem{}
}

// Check checks if any player has won
// Returns (hasWinner, winnerID, reason)
func (s *WinConditionSystem) Check(state *State) (bool, int, string) {
	// Check if any tank has reached the enemy base
	// Only tanks can capture bases (not players or helicopters)
	for _, unit := range state.Units {
		unitType := unit.GetType()
		if unitType != "tank" && unitType != "super_tank" {
			continue
		}

		ownerID := unit.GetOwnerID()
		enemyID := 1 - ownerID
		enemyBase := state.Players[enemyID].BasePosition

		// Check distance to enemy base
		distance := calculateDistance(unit.GetPosition(), enemyBase)

		if distance < types.BaseRadius+5 {
			return true, ownerID, "Tank captured enemy base"
		}
	}

	return false, -1, ""
}
