package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tombuildsstuff/web-arena-game/server/internal/auth"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, restrict this to your client domain
		return true
	},
}

// HandleWebSocket upgrades HTTP connections to WebSocket
func HandleWebSocket(hub *Hub, authHandler *auth.Handler, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	// Generate a unique ID for this client
	clientID := uuid.New().String()

	// Get user info from auth token (if present)
	userInfo := authHandler.GetUserFromRequest(r)
	if userInfo == nil {
		// No auth token - check for guest token
		userInfo = authHandler.GetGuestFromRequest(r)
	}
	if userInfo == nil {
		// No guest token either - create ephemeral guest (shouldn't happen if /api/me was called)
		userInfo = auth.GenerateGuestUser()
		log.Printf("Warning: WebSocket connected without guest session, created ephemeral guest")
	}

	client := NewClient(hub, conn, clientID, userInfo.DisplayName, userInfo.IsGuest)
	hub.Register <- client

	log.Printf("Client connected: %s (%s, guest=%v)", clientID, userInfo.DisplayName, userInfo.IsGuest)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
