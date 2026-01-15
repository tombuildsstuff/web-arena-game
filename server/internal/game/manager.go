package game

import (
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/tombuildsstuff/web-arena-game/server/internal/game/maps"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// Manager manages all game rooms and matchmaking
type Manager struct {
	rooms      map[string]*GameRoom
	roomsMutex sync.RWMutex

	// Matchmaking queue: clientID -> PlayerQueueEntry
	queue      map[string]*PlayerQueueEntry
	queueMutex sync.Mutex

	// Client to room mapping (for players)
	clientToRoom map[string]string // clientID -> roomID

	// Spectator to room mapping
	spectatorToRoom map[string]string // clientID -> roomID

	// Leaderboard
	leaderboard *Leaderboard
}

// PlayerQueueEntry represents a player in the matchmaking queue
type PlayerQueueEntry struct {
	ClientID      string
	Connection    ClientConnection
	DisplayName   string
	IsGuest       bool
	MapPreference string // Preferred map ID (empty = no preference)
}

// NewManager creates a new game manager
func NewManager() *Manager {
	leaderboardFile := os.Getenv("LEADERBOARD_FILE")
	if leaderboardFile == "" {
		leaderboardFile = "leaderboard.json"
	}

	return &Manager{
		rooms:           make(map[string]*GameRoom),
		queue:           make(map[string]*PlayerQueueEntry),
		clientToRoom:    make(map[string]string),
		spectatorToRoom: make(map[string]string),
		leaderboard:     NewLeaderboard(leaderboardFile),
	}
}

// GetLeaderboard returns the leaderboard instance
func (m *Manager) GetLeaderboard() *Leaderboard {
	return m.leaderboard
}

// AddToQueue adds a player to the matchmaking queue
func (m *Manager) AddToQueue(clientID string, conn ClientConnection, displayName string, isGuest bool, mapPreference string) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()

	// Check if already in queue
	if _, exists := m.queue[clientID]; exists {
		return
	}

	// Add to queue
	m.queue[clientID] = &PlayerQueueEntry{
		ClientID:      clientID,
		Connection:    conn,
		DisplayName:   displayName,
		IsGuest:       isGuest,
		MapPreference: mapPreference,
	}

	// Try to match players
	m.tryMatchPlayers()
}

// RemoveFromQueue removes a player from the queue
func (m *Manager) RemoveFromQueue(clientID string) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()

	if _, exists := m.queue[clientID]; exists {
		delete(m.queue, clientID)
	}
}

// resolveMapVote determines which map to use based on two players' preferences
func resolveMapVote(pref1, pref2 string) *types.MapDefinition {
	// Both have same preference (including both empty)
	if pref1 == pref2 {
		if pref1 == "" {
			return maps.GetDefault()
		}
		if m, err := maps.Get(pref1); err == nil {
			return m
		}
		return maps.GetDefault()
	}

	// One has preference, other doesn't
	if pref1 == "" && pref2 != "" {
		if m, err := maps.Get(pref2); err == nil {
			return m
		}
		return maps.GetDefault()
	}
	if pref2 == "" && pref1 != "" {
		if m, err := maps.Get(pref1); err == nil {
			return m
		}
		return maps.GetDefault()
	}

	// Both have different preferences - random pick
	chosen := pref1
	if rand.Intn(2) == 1 {
		chosen = pref2
	}
	if m, err := maps.Get(chosen); err == nil {
		return m
	}
	return maps.GetDefault()
}

// tryMatchPlayers attempts to match players from the queue
func (m *Manager) tryMatchPlayers() {
	if len(m.queue) < 2 {
		return
	}

	// Get first two players from queue
	var player1, player2 *PlayerQueueEntry
	for _, entry := range m.queue {
		if player1 == nil {
			player1 = entry
		} else if player2 == nil {
			player2 = entry
			break
		}
	}

	if player1 == nil || player2 == nil {
		return
	}

	// Remove from queue
	delete(m.queue, player1.ClientID)
	delete(m.queue, player2.ClientID)

	// Resolve map vote
	mapDef := resolveMapVote(player1.MapPreference, player2.MapPreference)
	log.Printf("Map vote resolved: P1=%q, P2=%q -> %s", player1.MapPreference, player2.MapPreference, mapDef.Name)

	// Create game room with display names
	gameID := uuid.New().String()
	room := NewGameRoomWithMap(gameID, mapDef,
		player1.ClientID, player1.DisplayName, player1.IsGuest,
		player2.ClientID, player2.DisplayName, player2.IsGuest,
	)

	// Set client connections
	room.SetClientConnection(0, player1.Connection)
	room.SetClientConnection(1, player2.Connection)

	// Set callback for when game ends
	room.SetOnGameEnd(m.handleGameEnd)

	// Set callback for recording game results to leaderboard
	room.SetOnGameResult(m.leaderboard.RecordGameResult)

	// Store room
	m.roomsMutex.Lock()
	m.rooms[gameID] = room
	m.clientToRoom[player1.ClientID] = gameID
	m.clientToRoom[player2.ClientID] = gameID
	m.roomsMutex.Unlock()

	log.Printf("Created game room %s with players %s/%s (P1) and %s/%s (P2)", gameID, player1.ClientID, player1.DisplayName, player2.ClientID, player2.DisplayName)

	// Start the game
	room.Start()
}

// CreateAIGame creates a game with a human player vs AI
func (m *Manager) CreateAIGame(clientID string, conn ClientConnection, displayName string, isGuest bool, difficulty string, mapPreference string) {
	// Remove from queue if present
	m.queueMutex.Lock()
	delete(m.queue, clientID)
	m.queueMutex.Unlock()

	// Get the requested map or default
	mapDef := maps.GetDefault()
	if mapPreference != "" {
		if m, err := maps.Get(mapPreference); err == nil {
			mapDef = m
		}
	}

	// Create game room with AI as player 2
	gameID := uuid.New().String()
	aiDisplayName := "AI (" + difficulty + ")"

	room := NewGameRoomWithMap(gameID, mapDef,
		clientID, displayName, isGuest,
		"ai-"+gameID, aiDisplayName, false,
	)

	// Set human player connection
	room.SetClientConnection(0, conn)

	// Set AI connection (dummy - AI doesn't need to receive messages)
	room.SetClientConnection(1, &AIClientConnection{})

	// Create and set AI controller (AI is player 1, index 1)
	aiController := NewAIController(1, difficulty)
	room.SetAIController(aiController)

	// Set callback for when game ends
	room.SetOnGameEnd(m.handleGameEnd)

	// Set callback for recording game results to leaderboard
	room.SetOnGameResult(m.leaderboard.RecordGameResult)

	// Store room
	m.roomsMutex.Lock()
	m.rooms[gameID] = room
	m.clientToRoom[clientID] = gameID
	m.roomsMutex.Unlock()

	log.Printf("Created AI game room %s: %s vs %s", gameID, displayName, aiDisplayName)

	// Start the game
	room.Start()
}

// handleGameEnd cleans up when a game finishes
func (m *Manager) handleGameEnd(roomID string) {
	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return
	}

	// Get client IDs before removing room
	clientIDs := room.GetClientIDs()

	// Remove client-to-room mappings
	for _, clientID := range clientIDs {
		delete(m.clientToRoom, clientID)
	}

	// Remove room
	delete(m.rooms, roomID)

	// Verbose: log.Printf("Cleaned up game room %s (removed %d client mappings)", roomID, len(clientIDs))
}

// GetRoomByClient returns the game room for a client
func (m *Manager) GetRoomByClient(clientID string) *GameRoom {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()

	if roomID, exists := m.clientToRoom[clientID]; exists {
		return m.rooms[roomID]
	}

	return nil
}

// GetPlayerIDInRoom returns the player ID for a client in their game room
func (m *Manager) GetPlayerIDInRoom(clientID string) int {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()

	roomID, exists := m.clientToRoom[clientID]
	if !exists {
		return -1
	}

	room := m.rooms[roomID]
	if room == nil {
		return -1
	}

	// Check which player this client is
	for i, player := range room.State.Players {
		if player.ClientID == clientID {
			return i
		}
	}

	return -1
}

// RemoveClient removes a client from their game room or spectating session
func (m *Manager) RemoveClient(clientID string) {
	m.queueMutex.Lock()
	delete(m.queue, clientID)
	m.queueMutex.Unlock()

	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	// Check if client is a player
	if roomID, exists := m.clientToRoom[clientID]; exists {
		if room, ok := m.rooms[roomID]; ok {
			room.Stop()
			delete(m.rooms, roomID)
		}
		delete(m.clientToRoom, clientID)
		return
	}

	// Check if client is a spectator
	if gameID, exists := m.spectatorToRoom[clientID]; exists {
		if room, ok := m.rooms[gameID]; ok {
			room.RemoveSpectator(clientID)
		}
		delete(m.spectatorToRoom, clientID)
	}
}

// GetRoom returns a game room by ID
func (m *Manager) GetRoom(roomID string) *GameRoom {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()
	return m.rooms[roomID]
}

// GetActiveRoomsCount returns the number of active game rooms
func (m *Manager) GetActiveRoomsCount() int {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()
	return len(m.rooms)
}

// GetQueueSize returns the current queue size
func (m *Manager) GetQueueSize() int {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()
	return len(m.queue)
}

// GetActiveGames returns a list of active games for the lobby
func (m *Manager) GetActiveGames() []ActiveGameInfo {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()

	games := make([]ActiveGameInfo, 0, len(m.rooms))
	for _, room := range m.rooms {
		if room.IsRunning {
			info := room.GetGameInfo()
			games = append(games, ActiveGameInfo{
				GameID:         info.GameID,
				Player1Name:    info.Player1Name,
				Player2Name:    info.Player2Name,
				SpectatorCount: info.SpectatorCount,
			})
		}
	}
	return games
}

// ActiveGameInfo contains information about an active game
type ActiveGameInfo struct {
	GameID         string
	Player1Name    string
	Player2Name    string
	SpectatorCount int
}

// AddSpectator adds a spectator to a game room
func (m *Manager) AddSpectator(clientID string, gameID string, conn ClientConnection) bool {
	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	room, exists := m.rooms[gameID]
	if !exists || !room.IsRunning {
		return false
	}

	// Track spectator
	m.spectatorToRoom[clientID] = gameID
	room.AddSpectator(clientID, conn)
	return true
}

// RemoveSpectator removes a spectator from their game
func (m *Manager) RemoveSpectator(clientID string) {
	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	if gameID, exists := m.spectatorToRoom[clientID]; exists {
		if room, ok := m.rooms[gameID]; ok {
			room.RemoveSpectator(clientID)
		}
		delete(m.spectatorToRoom, clientID)
	}
}

// IsSpectating checks if a client is currently spectating a game
func (m *Manager) IsSpectating(clientID string) bool {
	m.roomsMutex.RLock()
	defer m.roomsMutex.RUnlock()
	_, exists := m.spectatorToRoom[clientID]
	return exists
}
