package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tombuildsstuff/web-arena-game/server/internal/game"
	"github.com/tombuildsstuff/web-arena-game/server/internal/websocket"
)

//go:embed static/*
var staticContent embed.FS

func main() {
	// Create game manager
	gameManager := game.NewManager()

	// Create WebSocket hub
	hub := websocket.NewHub(gameManager)
	go hub.Run()

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))

	// API Routes
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/api/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		leaderboard := gameManager.GetLeaderboard()
		entries := leaderboard.GetTopPlayers(50) // Top 50 players

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(hub, w, r)
	})

	// Try to get the embedded filesystem
	staticFS, err := fs.Sub(staticContent, "static")
	if err != nil {
		panic(err)
	}

	// Serve everything else as a static file
	fileServer := http.FileServer(http.FS(staticFS))
	r.Get("/*", fileServer.ServeHTTP)

	// Start server
	port := ":3000"
	log.Printf("Server starting on %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", port)
	log.Printf("Web UI: http://localhost%s/", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
