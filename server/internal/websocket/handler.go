package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	// Generate a unique ID for this client
	clientID := uuid.New().String()

	client := NewClient(hub, conn, clientID)
	hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
