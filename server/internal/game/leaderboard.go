package game

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// LeaderboardEntry represents a player's all-time statistics
type LeaderboardEntry struct {
	PlayerName    string `json:"playerName"`
	TankKills     int    `json:"tankKills"`
	AirplaneKills int    `json:"airplaneKills"`
	TurretKills   int    `json:"turretKills"`
	PlayerKills   int    `json:"playerKills"`
	TotalPoints   int    `json:"totalPoints"`
	GamesPlayed   int    `json:"gamesPlayed"`
	GamesWon      int    `json:"gamesWon"`
	TotalPlayTime int    `json:"totalPlayTime"` // in seconds
	LastPlayed    int64  `json:"lastPlayed"`    // Unix timestamp
}

// Leaderboard manages player statistics
type Leaderboard struct {
	entries      map[string]*LeaderboardEntry
	totalMatches int
	mu           sync.RWMutex
	filePath     string
}

// NewLeaderboard creates a new leaderboard, loading from file if exists
func NewLeaderboard(filePath string) *Leaderboard {
	lb := &Leaderboard{
		entries:  make(map[string]*LeaderboardEntry),
		filePath: filePath,
	}
	lb.load()
	return lb
}

// RecordGameResult records the result of a game for both players
func (lb *Leaderboard) RecordGameResult(player1Name, player2Name string, winner int, matchDuration int, p1Stats, p2Stats types.PlayerStats) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now().Unix()

	// Update player 1 stats
	entry1 := lb.getOrCreateEntry(player1Name)
	entry1.TankKills += p1Stats.TankKills
	entry1.AirplaneKills += p1Stats.AirplaneKills
	entry1.TurretKills += p1Stats.TurretKills
	entry1.PlayerKills += p1Stats.PlayerKills
	entry1.TotalPoints += p1Stats.TotalPoints
	entry1.GamesPlayed++
	entry1.TotalPlayTime += matchDuration
	entry1.LastPlayed = now
	if winner == 0 {
		entry1.GamesWon++
	}

	// Update player 2 stats
	entry2 := lb.getOrCreateEntry(player2Name)
	entry2.TankKills += p2Stats.TankKills
	entry2.AirplaneKills += p2Stats.AirplaneKills
	entry2.TurretKills += p2Stats.TurretKills
	entry2.PlayerKills += p2Stats.PlayerKills
	entry2.TotalPoints += p2Stats.TotalPoints
	entry2.GamesPlayed++
	entry2.TotalPlayTime += matchDuration
	entry2.LastPlayed = now
	if winner == 1 {
		entry2.GamesWon++
	}

	// Save to file
	lb.saveUnlocked()
}

// GetTopPlayers returns the top N players by points
func (lb *Leaderboard) GetTopPlayers(limit int) []LeaderboardEntry {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Convert map to slice
	entries := make([]LeaderboardEntry, 0, len(lb.entries))
	for _, entry := range lb.entries {
		entries = append(entries, *entry)
	}

	// Sort by total points descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TotalPoints > entries[j].TotalPoints
	})

	// Limit results
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	return entries
}

// GetPlayerStats returns stats for a specific player
func (lb *Leaderboard) GetPlayerStats(playerName string) *LeaderboardEntry {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if entry, exists := lb.entries[playerName]; exists {
		// Return a copy
		entryCopy := *entry
		return &entryCopy
	}
	return nil
}

// GetTotalMatches returns the total number of matches played
func (lb *Leaderboard) GetTotalMatches() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return lb.totalMatches
}

// IncrementTotalMatches increments the total matches counter (called when a game starts)
func (lb *Leaderboard) IncrementTotalMatches() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.totalMatches++
	lb.saveUnlocked()
}

// getOrCreateEntry gets or creates an entry for a player (must hold write lock)
func (lb *Leaderboard) getOrCreateEntry(playerName string) *LeaderboardEntry {
	if entry, exists := lb.entries[playerName]; exists {
		return entry
	}

	entry := &LeaderboardEntry{
		PlayerName: playerName,
	}
	lb.entries[playerName] = entry
	return entry
}

// leaderboardData is the structure for JSON persistence
type leaderboardData struct {
	TotalMatches int                `json:"totalMatches"`
	Entries      []LeaderboardEntry `json:"entries"`
}

// load reads the leaderboard from the file
func (lb *Leaderboard) load() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	data, err := os.ReadFile(lb.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Leaderboard file not found, starting fresh")
			return
		}
		log.Printf("Error reading leaderboard file: %v", err)
		return
	}

	var lbData leaderboardData
	if err := json.Unmarshal(data, &lbData); err != nil {
		log.Printf("Error parsing leaderboard file: %v", err)
		return
	}

	lb.totalMatches = lbData.TotalMatches
	for _, entry := range lbData.Entries {
		entryCopy := entry
		lb.entries[entry.PlayerName] = &entryCopy
	}
	log.Printf("Loaded %d leaderboard entries, %d total matches", len(lb.entries), lb.totalMatches)
}

// saveUnlocked saves the leaderboard to file (must hold lock)
func (lb *Leaderboard) saveUnlocked() {
	entries := make([]LeaderboardEntry, 0, len(lb.entries))
	for _, entry := range lb.entries {
		entries = append(entries, *entry)
	}

	lbData := leaderboardData{
		TotalMatches: lb.totalMatches,
		Entries:      entries,
	}

	data, err := json.MarshalIndent(lbData, "", "  ")
	if err != nil {
		log.Printf("Error marshaling leaderboard: %v", err)
		return
	}

	if err := os.WriteFile(lb.filePath, data, 0644); err != nil {
		log.Printf("Error writing leaderboard file: %v", err)
		return
	}
}

// Save explicitly saves the leaderboard (for graceful shutdown)
func (lb *Leaderboard) Save() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.saveUnlocked()
}
