package game

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// GameRoom represents a single game instance
type GameRoom struct {
	ID                string
	State             *State
	IsRunning         bool
	mu                sync.RWMutex
	stopChan          chan bool
	clientConnections map[int]ClientConnection // Map player ID to client connection
	lastIncomeTime    time.Time

	// Game systems
	pathfindingSystem  *PathfindingSystem
	spatialGrid        *SpatialGrid
	losSystem          *LOSSystem
	movementSystem     *MovementSystem
	combatSystem       *CombatSystem
	turretSystem       *TurretSystem
	winConditionSystem *WinConditionSystem
}

// ClientConnection interface for sending messages to clients
type ClientConnection interface {
	SendMessage(msgType string, payload interface{})
}

// NewGameRoom creates a new game room
func NewGameRoom(id string, player1ClientID, player2ClientID string) *GameRoom {
	state := NewState(player1ClientID, player2ClientID)

	// Initialize spatial systems
	spatialGrid := NewSpatialGrid(state.Obstacles)
	pathfindingSystem := NewPathfindingSystem(state.Obstacles)
	losSystem := NewLOSSystem(spatialGrid)

	return &GameRoom{
		ID:                 id,
		State:              state,
		IsRunning:          false,
		stopChan:           make(chan bool),
		clientConnections:  make(map[int]ClientConnection),
		lastIncomeTime:     time.Now(),
		pathfindingSystem:  pathfindingSystem,
		spatialGrid:        spatialGrid,
		losSystem:          losSystem,
		movementSystem:     NewMovementSystem(pathfindingSystem),
		combatSystem:       NewCombatSystem(losSystem),
		turretSystem:       NewTurretSystem(losSystem),
		winConditionSystem: NewWinConditionSystem(),
	}
}

// SetClientConnection sets the client connection for a player
func (r *GameRoom) SetClientConnection(playerID int, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clientConnections[playerID] = conn
}

// Start starts the game room update loop
func (r *GameRoom) Start() {
	r.mu.Lock()
	r.IsRunning = true
	r.mu.Unlock()

	log.Printf("Game room %s started", r.ID)

	// Send game start message to both players
	r.broadcastGameStart()

	// Start the game loop in a goroutine
	go r.gameLoop()
}

// Stop stops the game room
func (r *GameRoom) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsRunning {
		r.IsRunning = false
		close(r.stopChan)
		log.Printf("Game room %s stopped", r.ID)
	}
}

// gameLoop is the main game loop running at 20 TPS
func (r *GameRoom) gameLoop() {
	ticker := time.NewTicker(types.TickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.update()
		case <-r.stopChan:
			return
		}
	}
}

// update updates the game state for one tick
func (r *GameRoom) update() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	deltaTime := float64(types.TickDuration) / float64(time.Second)

	// Update passive income
	r.updateIncome()

	// Update movement
	r.movementSystem.Update(r.State, deltaTime)

	// Update combat
	r.combatSystem.Update(r.State, deltaTime)

	// Update turrets (combat and respawns)
	r.turretSystem.Update(r.State, deltaTime)

	// Check player respawns
	r.checkPlayerRespawns()

	// Check win condition
	if hasWinner, winnerID, reason := r.winConditionSystem.Check(r.State); hasWinner {
		r.State.GameStatus = "finished"
		r.State.Winner = &winnerID
		r.broadcastGameOver(winnerID, reason)
		r.Stop()
		return
	}

	// Update timestamp
	r.State.UpdateTimestamp()

	// Broadcast state to all clients
	r.broadcastState()
}

// updateIncome handles passive income for both players
func (r *GameRoom) updateIncome() {
	now := time.Now()
	elapsed := now.Sub(r.lastIncomeTime).Seconds()

	// Give income approximately every second
	if elapsed >= 1.0 {
		incomeAmount := int(elapsed * types.PassiveIncomePerSecond)

		for _, player := range r.State.Players {
			player.AddMoney(incomeAmount)
		}

		r.lastIncomeTime = now
	}
}

// checkPlayerRespawns checks if any player units need to respawn
func (r *GameRoom) checkPlayerRespawns() {
	for _, unit := range r.State.Units {
		if playerUnit, ok := unit.(*PlayerUnit); ok {
			playerUnit.CheckRespawn()
		}
	}
}

// HandlePurchase handles a unit purchase request
func (r *GameRoom) HandlePurchase(playerID int, unitType string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	player := r.State.GetPlayer(playerID)
	if player == nil {
		log.Printf("Player %d not found", playerID)
		return
	}

	// Determine cost and spawn position
	var cost int
	var spawnPos, targetPos types.Vector3

	spawnPos = player.BasePosition
	spawnPos.Y = types.TankYPosition // Will be overridden for airplanes

	// Target is the enemy base
	enemyPlayer := r.State.GetPlayer(1 - playerID)
	targetPos = enemyPlayer.BasePosition

	switch unitType {
	case "tank":
		cost = types.TankCost
	case "airplane":
		cost = types.AirplaneCost
	default:
		log.Printf("Unknown unit type: %s", unitType)
		return
	}

	// Check if player can afford
	if !player.CanAfford(cost) {
		log.Printf("Player %d cannot afford %s (cost: %d, money: %d)", playerID, unitType, cost, player.Money)
		if conn, ok := r.clientConnections[playerID]; ok {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "Not enough money",
			})
		}
		return
	}

	// Deduct cost
	player.Spend(cost)

	// Create unit
	var unit Unit
	switch unitType {
	case "tank":
		unit = NewTank(playerID, spawnPos, targetPos)
	case "airplane":
		unit = NewAirplane(playerID, spawnPos, targetPos)
	}

	// Add to state
	r.State.AddUnit(unit)

	log.Printf("Player %d purchased %s (remaining money: %d)", playerID, unitType, player.Money)
}

// HandlePlayerMove handles player movement input
func (r *GameRoom) HandlePlayerMove(playerID int, direction types.Vector3) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil {
		return
	}

	playerUnit.SetMoveDirection(direction)
}

// HandleBuyFromZone handles a buy from zone request
func (r *GameRoom) HandleBuyFromZone(playerID int, zoneID string, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	// Find the buy zone
	var zone *BuyZone
	for _, z := range r.State.BuyZones {
		if z.ID == zoneID {
			zone = z
			break
		}
	}

	if zone == nil {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Invalid buy zone",
		})
		return
	}

	// Check if zone belongs to this player
	if zone.OwnerID != playerID {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "This is not your buy zone",
		})
		return
	}

	// Get the player
	player := r.State.GetPlayer(playerID)
	if player == nil {
		return
	}

	// Get the player unit
	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil || !playerUnit.IsAlive() {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "You must be alive to purchase",
		})
		return
	}

	// Check if player is near the buy zone
	if !zone.IsPlayerInRange(playerUnit.GetPosition()) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Get closer to the buy zone",
		})
		return
	}

	// Check if player can afford
	if !player.CanAfford(zone.Cost) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Not enough money",
		})
		return
	}

	// Deduct cost
	player.Spend(zone.Cost)

	// Create unit at the buy zone position
	var unit Unit
	spawnPos := zone.Position
	targetPos := r.State.Players[1-playerID].BasePosition // Target enemy base

	switch zone.UnitType {
	case "tank":
		spawnPos.Y = types.TankYPosition
		unit = NewTank(playerID, spawnPos, targetPos)
	case "airplane":
		spawnPos.Y = types.AirplaneYPosition
		unit = NewAirplane(playerID, spawnPos, targetPos)
	}

	if unit != nil {
		r.State.AddUnit(unit)
		log.Printf("Player %d purchased %s from zone %s (remaining: $%d)", playerID, zone.UnitType, zoneID, player.Money)
	}
}

// HandleClaimTurret handles a turret claiming request
func (r *GameRoom) HandleClaimTurret(playerID int, turretID string, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	// Find the turret
	turret := r.State.GetTurretByID(turretID)
	if turret == nil {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Invalid turret",
		})
		return
	}

	// Get the player unit
	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil || !playerUnit.IsAlive() {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "You must be alive to claim a turret",
		})
		return
	}

	// Check if player is near the turret
	if !turret.IsPlayerInRange(playerUnit.GetPosition()) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Get closer to the turret",
		})
		return
	}

	// Check if turret can be claimed
	if !turret.CanBeClaimed(playerID) {
		if turret.IsDestroyed {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "Turret is destroyed",
			})
		} else {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "You already own this turret",
			})
		}
		return
	}

	// Claim the turret
	turret.Claim(playerID)
	log.Printf("Player %d claimed turret %s", playerID, turretID)
}

// HandlePlayerShoot handles player shoot command
func (r *GameRoom) HandlePlayerShoot(playerID int, targetX, targetZ float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil || !playerUnit.IsAlive() {
		return
	}

	// Check attack cooldown
	now := time.Now().UnixMilli()
	timeSinceLastAttack := now - playerUnit.GetLastAttackTime()
	attackCooldown := int64(1000.0 / playerUnit.GetAttackSpeed())

	if timeSinceLastAttack < attackCooldown {
		return // Still on cooldown
	}

	// Create target position
	targetPos := types.Vector3{
		X: targetX,
		Y: playerUnit.Position.Y,
		Z: targetZ,
	}

	// Check if target is in range
	distance := calculateDistance(playerUnit.GetPosition(), targetPos)
	if distance > playerUnit.GetAttackRange() {
		return // Out of range
	}

	// Check line of sight
	if !r.losSystem.HasLineOfSight(playerUnit.GetPosition(), targetPos, false) {
		return // Blocked by obstacle
	}

	// Find if there's a unit at the target position (or close to it)
	var targetUnit Unit
	for _, unit := range r.State.Units {
		if unit.GetOwnerID() == playerID {
			continue // Don't shoot own units
		}
		if !unit.IsAlive() {
			continue
		}

		unitPos := unit.GetPosition()
		distToTarget := calculateDistance(unitPos, targetPos)
		if distToTarget < 3.0 { // Close enough to target position
			targetUnit = unit
			break
		}
	}

	// Check if there's a turret at the target position (enemy turrets only)
	var targetTurret *Turret
	if targetUnit == nil {
		for _, turret := range r.State.Turrets {
			if turret.OwnerID == playerID {
				continue // Don't shoot own turrets
			}
			if !turret.IsAlive() {
				continue
			}

			distToTarget := calculateDistance(turret.Position, targetPos)
			if distToTarget < 4.0 { // Close enough to target position (turrets are bigger)
				targetTurret = turret
				break
			}
		}
	}

	// Create projectile
	var projectile *Projectile
	if targetTurret != nil {
		projectile = NewProjectileFromPlayerToTurret(playerUnit, targetTurret, now)
	} else {
		projectile = NewProjectileFromPlayer(playerUnit, targetPos, targetUnit, now)
	}
	r.State.AddProjectile(projectile)
	playerUnit.SetLastAttackTime(now)
}

// broadcastGameStart sends the game_start message to all players
func (r *GameRoom) broadcastGameStart() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for playerID, conn := range r.clientConnections {
		payload := types.GameStartPayload{
			GameID:   r.ID,
			PlayerID: playerID,
			State:    r.State.ToType(),
		}
		conn.SendMessage("game_start", payload)
	}
}

// broadcastState sends the current game state to all players
func (r *GameRoom) broadcastState() {
	stateData := r.State.ToType()

	for _, conn := range r.clientConnections {
		payload := types.GameUpdatePayload{
			Timestamp: r.State.Timestamp,
			State:     stateData,
		}
		conn.SendMessage("game_update", payload)
	}
}

// broadcastGameOver sends the game_over message to all players
// Note: Caller must hold the lock (called from update() which holds r.mu.Lock())
func (r *GameRoom) broadcastGameOver(winner int, reason string) {
	payload := types.GameOverPayload{
		Winner:        winner,
		Reason:        reason,
		MatchDuration: r.State.GetMatchDuration(),
		Stats: types.MatchStats{
			Player1Kills: r.State.Players[0].Kills,
			Player2Kills: r.State.Players[1].Kills,
		},
	}

	for _, conn := range r.clientConnections {
		conn.SendMessage("game_over", payload)
	}

	log.Printf("Game %s ended: Player %d won (%s) - Duration: %ds, P1 Kills: %d, P2 Kills: %d",
		r.ID, winner, reason, payload.MatchDuration,
		payload.Stats.Player1Kills, payload.Stats.Player2Kills)
}

// GetState returns a copy of the current game state (thread-safe)
func (r *GameRoom) GetState() types.GameState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.State.ToType()
}

// Marshal game room to JSON (for debugging)
func (r *GameRoom) MarshalJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return json.Marshal(struct {
		ID        string           `json:"id"`
		State     types.GameState  `json:"state"`
		IsRunning bool             `json:"isRunning"`
	}{
		ID:        r.ID,
		State:     r.State.ToType(),
		IsRunning: r.IsRunning,
	})
}
