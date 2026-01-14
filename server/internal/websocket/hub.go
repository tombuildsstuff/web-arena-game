package websocket

import (
	"encoding/json"
	"log"

	"github.com/tombuildsstuff/web-arena-game/server/internal/game"
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// ClientMessage pairs a client with their message
type ClientMessage struct {
	Client  *Client
	Message types.Message
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Inbound messages from clients
	HandleMessage chan *ClientMessage

	// Map of client ID to client
	clientsByID map[string]*Client

	// Game manager
	gameManager *game.Manager
}

// NewHub creates a new Hub
func NewHub(gameManager *game.Manager) *Hub {
	return &Hub{
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		HandleMessage: make(chan *ClientMessage),
		clients:       make(map[*Client]bool),
		clientsByID:   make(map[string]*Client),
		gameManager:   gameManager,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			h.clientsByID[client.ID] = client
			log.Printf("Client connected: %s (total: %d)", client.ID, len(h.clients))

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByID, client.ID)
				close(client.send)

				// Remove from game manager (handles both queue and active games)
				h.gameManager.RemoveClient(client.ID)

				log.Printf("Client disconnected: %s (total: %d)", client.ID, len(h.clients))
			}

		case clientMsg := <-h.HandleMessage:
			h.handleClientMessage(clientMsg)
		}
	}
}

// handleClientMessage processes messages from clients
func (h *Hub) handleClientMessage(clientMsg *ClientMessage) {
	msg := clientMsg.Message
	client := clientMsg.Client

	log.Printf("Received message from %s: %s", client.ID, msg.Type)

	switch msg.Type {
	case "join_queue":
		h.handleJoinQueue(client)

	case "purchase_unit":
		h.handlePurchaseUnit(client, msg.Payload)

	case "player_move":
		h.handlePlayerMove(client, msg.Payload)

	case "player_shoot":
		h.handlePlayerShoot(client, msg.Payload)

	case "buy_from_zone":
		h.handleBuyFromZone(client, msg.Payload)

	case "claim_turret":
		h.handleClaimTurret(client, msg.Payload)

	case "leave_game":
		h.handleLeaveGame(client)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
		client.SendMessage("error", types.ErrorPayload{
			Message: "Unknown message type",
		})
	}
}

// handleJoinQueue adds a client to the matchmaking queue
func (h *Hub) handleJoinQueue(client *Client) {
	h.gameManager.AddToQueue(client.ID, client)
}

// handlePurchaseUnit processes a unit purchase request
func (h *Hub) handlePurchaseUnit(client *Client, payload interface{}) {
	// Parse payload
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshaling payload: %v", err)
		return
	}

	var purchase types.PurchaseUnitPayload
	if err := json.Unmarshal(data, &purchase); err != nil {
		log.Printf("error unmarshaling purchase: %v", err)
		client.SendMessage("error", types.ErrorPayload{
			Message: "Invalid purchase payload",
		})
		return
	}

	// Get the game room for this client
	room := h.gameManager.GetRoomByClient(client.ID)
	if room == nil {
		client.SendMessage("error", types.ErrorPayload{
			Message: "You are not in a game",
		})
		return
	}

	// Get player ID
	playerID := h.gameManager.GetPlayerIDInRoom(client.ID)
	if playerID < 0 {
		client.SendMessage("error", types.ErrorPayload{
			Message: "Invalid player",
		})
		return
	}

	// Forward to game room
	room.HandlePurchase(playerID, purchase.UnitType)
}

// handlePlayerMove processes player movement input
func (h *Hub) handlePlayerMove(client *Client, payload interface{}) {
	// Parse payload
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var move types.PlayerMovePayload
	if err := json.Unmarshal(data, &move); err != nil {
		return
	}

	// Get the game room for this client
	room := h.gameManager.GetRoomByClient(client.ID)
	if room == nil {
		return
	}

	// Get player ID
	playerID := h.gameManager.GetPlayerIDInRoom(client.ID)
	if playerID < 0 {
		return
	}

	// Forward to game room
	room.HandlePlayerMove(playerID, move.Direction)
}

// handlePlayerShoot processes player shoot command
func (h *Hub) handlePlayerShoot(client *Client, payload interface{}) {
	// Parse payload
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var shoot types.PlayerShootPayload
	if err := json.Unmarshal(data, &shoot); err != nil {
		return
	}

	// Get the game room for this client
	room := h.gameManager.GetRoomByClient(client.ID)
	if room == nil {
		return
	}

	// Get player ID
	playerID := h.gameManager.GetPlayerIDInRoom(client.ID)
	if playerID < 0 {
		return
	}

	// Forward to game room
	room.HandlePlayerShoot(playerID, shoot.TargetX, shoot.TargetZ)
}

// handleBuyFromZone processes a buy from zone request
func (h *Hub) handleBuyFromZone(client *Client, payload interface{}) {
	// Parse payload
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var buy types.BuyFromZonePayload
	if err := json.Unmarshal(data, &buy); err != nil {
		return
	}

	// Get the game room for this client
	room := h.gameManager.GetRoomByClient(client.ID)
	if room == nil {
		return
	}

	// Get player ID
	playerID := h.gameManager.GetPlayerIDInRoom(client.ID)
	if playerID < 0 {
		return
	}

	// Forward to game room
	room.HandleBuyFromZone(playerID, buy.ZoneID, client)
}

// handleClaimTurret processes a turret claiming request
func (h *Hub) handleClaimTurret(client *Client, payload interface{}) {
	// Parse payload
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var claim types.ClaimTurretPayload
	if err := json.Unmarshal(data, &claim); err != nil {
		return
	}

	// Get the game room for this client
	room := h.gameManager.GetRoomByClient(client.ID)
	if room == nil {
		return
	}

	// Get player ID
	playerID := h.gameManager.GetPlayerIDInRoom(client.ID)
	if playerID < 0 {
		return
	}

	// Forward to game room
	room.HandleClaimTurret(playerID, claim.TurretID, client)
}

// handleLeaveGame removes a client from their current game
func (h *Hub) handleLeaveGame(client *Client) {
	log.Printf("Client %s leaving game", client.ID)
	h.gameManager.RemoveClient(client.ID)
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(msgType string, payload interface{}) {
	msg := types.Message{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling broadcast message: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
