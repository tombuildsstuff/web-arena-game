package game

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// Manager manages all game rooms and matchmaking
type Manager struct {
	rooms      map[string]*GameRoom
	roomsMutex sync.RWMutex

	// Matchmaking queue: clientID -> PlayerQueueEntry
	queue      map[string]*PlayerQueueEntry
	queueMutex sync.Mutex

	// Client to room mapping
	clientToRoom map[string]string // clientID -> roomID
}

// PlayerQueueEntry represents a player in the matchmaking queue
type PlayerQueueEntry struct {
	ClientID   string
	Connection ClientConnection
}

// NewManager creates a new game manager
func NewManager() *Manager {
	return &Manager{
		rooms:        make(map[string]*GameRoom),
		queue:        make(map[string]*PlayerQueueEntry),
		clientToRoom: make(map[string]string),
	}
}

// AddToQueue adds a player to the matchmaking queue
func (m *Manager) AddToQueue(clientID string, conn ClientConnection) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()

	// Check if already in queue
	if _, exists := m.queue[clientID]; exists {
		log.Printf("Client %s already in queue", clientID)
		return
	}

	// Add to queue
	m.queue[clientID] = &PlayerQueueEntry{
		ClientID:   clientID,
		Connection: conn,
	}

	log.Printf("Client %s added to queue (queue size: %d)", clientID, len(m.queue))

	// Try to match players
	m.tryMatchPlayers()
}

// RemoveFromQueue removes a player from the queue
func (m *Manager) RemoveFromQueue(clientID string) {
	m.queueMutex.Lock()
	defer m.queueMutex.Unlock()

	if _, exists := m.queue[clientID]; exists {
		delete(m.queue, clientID)
		log.Printf("Client %s removed from queue", clientID)
	}
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

	// Create game room
	gameID := uuid.New().String()
	room := NewGameRoom(gameID, player1.ClientID, player2.ClientID)

	// Set client connections
	room.SetClientConnection(0, player1.Connection)
	room.SetClientConnection(1, player2.Connection)

	// Store room
	m.roomsMutex.Lock()
	m.rooms[gameID] = room
	m.clientToRoom[player1.ClientID] = gameID
	m.clientToRoom[player2.ClientID] = gameID
	m.roomsMutex.Unlock()

	log.Printf("Created game room %s with players %s (P1) and %s (P2)", gameID, player1.ClientID, player2.ClientID)

	// Start the game
	room.Start()
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

// RemoveClient removes a client from their game room
func (m *Manager) RemoveClient(clientID string) {
	m.queueMutex.Lock()
	delete(m.queue, clientID)
	m.queueMutex.Unlock()

	m.roomsMutex.Lock()
	defer m.roomsMutex.Unlock()

	if roomID, exists := m.clientToRoom[clientID]; exists {
		if room, ok := m.rooms[roomID]; ok {
			room.Stop()
			delete(m.rooms, roomID)
		}
		delete(m.clientToRoom, clientID)

		log.Printf("Client %s removed from game room %s", clientID, roomID)
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
