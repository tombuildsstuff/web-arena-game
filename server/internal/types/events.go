package types

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// JoinQueuePayload is empty for now
type JoinQueuePayload struct{}

// PurchaseUnitPayload represents a request to purchase a unit
type PurchaseUnitPayload struct {
	UnitType string `json:"unitType"` // "tank" or "airplane"
}

// GameStartPayload is sent when a game begins
type GameStartPayload struct {
	GameID   string    `json:"gameId"`
	PlayerID int       `json:"playerId"`
	State    GameState `json:"state"`
}

// GameUpdatePayload is sent periodically with the current game state
type GameUpdatePayload struct {
	Timestamp int64     `json:"timestamp"`
	State     GameState `json:"state"`
}

// GameOverPayload is sent when the game ends
type GameOverPayload struct {
	Winner        int        `json:"winner"`
	Reason        string     `json:"reason"`
	MatchDuration int        `json:"matchDuration"` // Duration in seconds
	Stats         MatchStats `json:"stats"`
}

// PlayerStats contains detailed kill statistics for a player
type PlayerStats struct {
	TankKills     int `json:"tankKills"`
	AirplaneKills int `json:"airplaneKills"`
	TurretKills   int `json:"turretKills"`
	PlayerKills   int `json:"playerKills"`
	TotalPoints   int `json:"totalPoints"`
}

// MatchStats contains end-of-match statistics
type MatchStats struct {
	Player1Kills int         `json:"player1Kills"`
	Player2Kills int         `json:"player2Kills"`
	Player1Stats PlayerStats `json:"player1Stats"`
	Player2Stats PlayerStats `json:"player2Stats"`
}

// ErrorPayload is sent when an error occurs
type ErrorPayload struct {
	Message string `json:"message"`
}

// PlayerMovePayload represents player movement input
type PlayerMovePayload struct {
	Direction Vector3 `json:"direction"` // Movement direction (normalized by client)
}

// PlayerShootPayload represents player shoot command
type PlayerShootPayload struct {
	TargetX float64 `json:"targetX"` // World X coordinate to shoot at
	TargetZ float64 `json:"targetZ"` // World Z coordinate to shoot at
}

// BuyFromZonePayload represents a request to buy from a buy zone
type BuyFromZonePayload struct {
	ZoneID string `json:"zoneId"` // ID of the buy zone to purchase from
}

// ClaimTurretPayload represents a request to claim a turret
type ClaimTurretPayload struct {
	TurretID string `json:"turretId"` // ID of the turret to claim
}

// ClaimBuyZonePayload represents a request to claim a buy zone
type ClaimBuyZonePayload struct {
	ZoneID string `json:"zoneId"` // ID of the buy zone to claim
}

// SpectateStartPayload is sent when a spectator joins a game
type SpectateStartPayload struct {
	GameID string    `json:"gameId"`
	State  GameState `json:"state"`
}

// SpectateGamePayload represents a request to spectate a game
type SpectateGamePayload struct {
	GameID string `json:"gameId"` // ID of the game to spectate
}

// ActiveGame represents a game available for spectating
type ActiveGame struct {
	GameID         string `json:"gameId"`
	Player1Name    string `json:"player1Name"`
	Player2Name    string `json:"player2Name"`
	SpectatorCount int    `json:"spectatorCount"`
}

// LobbyStatusPayload is sent to clients with queue and game information
type LobbyStatusPayload struct {
	QueueSize   int          `json:"queueSize"`
	ActiveGames []ActiveGame `json:"activeGames"`
}

// StartVsAIPayload represents a request to start a game vs AI
type StartVsAIPayload struct {
	Difficulty string `json:"difficulty"` // "easy", "medium", or "hard"
}
