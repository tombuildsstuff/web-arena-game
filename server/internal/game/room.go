package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// GameEndCallback is called when a game ends
type GameEndCallback func(roomID string)

// GameResultCallback is called with game results for leaderboard
type GameResultCallback func(player1Name, player2Name string, winner int, matchDuration int, p1Stats, p2Stats types.PlayerStats)

// GameRoom represents a single game instance
type GameRoom struct {
	ID                string
	State             *State
	IsRunning         bool
	mu                sync.RWMutex
	stopChan          chan bool
	clientConnections map[int]ClientConnection    // Map player ID to client connection
	spectators        map[string]ClientConnection // Map client ID to spectator connection
	lastIncomeTime    time.Time
	onGameEnd         GameEndCallback    // Callback when game ends
	onGameResult      GameResultCallback // Callback for game results (leaderboard)

	// AI controller (nil for human vs human games)
	aiController *AIController

	// Game systems
	pathfindingSystem  *PathfindingSystem
	spatialGrid        *SpatialGrid
	losSystem          *LOSSystem
	movementSystem     *MovementSystem
	combatSystem       *CombatSystem
	turretSystem       *TurretSystem
	healthPackSystem   *HealthPackSystem
	winConditionSystem *WinConditionSystem
}

// ClientConnection interface for sending messages to clients
type ClientConnection interface {
	SendMessage(msgType string, payload interface{})
}

// NewGameRoomWithMap creates a new game room using a specific map
func NewGameRoomWithMap(id string, mapDef *types.MapDefinition, player1ClientID, player1DisplayName string, player1IsGuest bool, player2ClientID, player2DisplayName string, player2IsGuest bool) *GameRoom {
	state := NewStateWithMap(mapDef, player1ClientID, player1DisplayName, player1IsGuest, player2ClientID, player2DisplayName, player2IsGuest)

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
		spectators:         make(map[string]ClientConnection),
		lastIncomeTime:     time.Now(),
		pathfindingSystem:  pathfindingSystem,
		spatialGrid:        spatialGrid,
		losSystem:          losSystem,
		movementSystem:     NewMovementSystem(pathfindingSystem),
		combatSystem:       NewCombatSystem(losSystem),
		turretSystem:       NewTurretSystem(losSystem),
		healthPackSystem:   NewHealthPackSystem(),
		winConditionSystem: NewWinConditionSystem(),
	}
}

// SetClientConnection sets the client connection for a player
func (r *GameRoom) SetClientConnection(playerID int, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clientConnections[playerID] = conn
}

// SetOnGameEnd sets the callback for when the game ends
func (r *GameRoom) SetOnGameEnd(callback GameEndCallback) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onGameEnd = callback
}

// SetOnGameResult sets the callback for game results (leaderboard)
func (r *GameRoom) SetOnGameResult(callback GameResultCallback) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onGameResult = callback
}

// SetAIController sets the AI controller for computer-controlled player
func (r *GameRoom) SetAIController(ai *AIController) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.aiController = ai
}

// AddSpectator adds a spectator to the game room
func (r *GameRoom) AddSpectator(clientID string, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spectators[clientID] = conn

	// Send current game state to the new spectator
	stateData := r.State.ToType()
	conn.SendMessage("spectate_start", types.SpectateStartPayload{
		GameID: r.ID,
		State:  stateData,
	})
}

// RemoveSpectator removes a spectator from the game room
func (r *GameRoom) RemoveSpectator(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.spectators[clientID]; exists {
		delete(r.spectators, clientID)
	}
}

// GetSpectatorCount returns the number of spectators
func (r *GameRoom) GetSpectatorCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.spectators)
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

// Stop stops the game room (called externally)
func (r *GameRoom) Stop() {
	r.mu.Lock()
	callback := r.stopInternal()
	r.mu.Unlock()

	// Call the callback outside of the lock to avoid deadlock
	if callback != nil {
		callback(r.ID)
	}
}

// stopInternal stops the game room - must be called with r.mu held
// Returns the callback to call (if any) after releasing the lock
func (r *GameRoom) stopInternal() GameEndCallback {
	if !r.IsRunning {
		return nil
	}

	r.IsRunning = false
	close(r.stopChan)
	log.Printf("Game room %s stopped", r.ID)
	return r.onGameEnd
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
	// Use a variable to store any callback that needs to be called after releasing the lock
	var endCallback GameEndCallback

	r.mu.Lock()
	defer func() {
		r.mu.Unlock()
		// Call the game end callback outside the lock to avoid deadlock with manager
		if endCallback != nil {
			endCallback(r.ID)
		}
	}()

	if r.State.GameStatus != "playing" {
		return
	}

	deltaTime := float64(types.TickDuration) / float64(time.Second)

	// Update passive income
	r.updateIncome()

	// Update AI (if present)
	if r.aiController != nil {
		r.aiController.Update(r.State, r)
	}

	// Update movement
	r.movementSystem.Update(r.State, deltaTime)

	// Process spawn queue - spawn units when their spawn area is clear
	if r.State.SpawnQueue != nil {
		spawnedUnits := r.State.SpawnQueue.ProcessQueue(r.State)
		for _, unit := range spawnedUnits {
			r.State.AddUnit(unit)
		}
	}

	// Update combat
	r.combatSystem.Update(r.State, deltaTime)

	// Update turrets (combat and respawns)
	r.turretSystem.Update(r.State, deltaTime)

	// Update health packs (spawning and collection)
	r.healthPackSystem.Update(r.State)

	// Update barracks (respawn timer)
	r.updateBarracks(deltaTime)

	// Check player respawns
	r.checkPlayerRespawns()

	// Check win condition
	if hasWinner, winnerID, reason := r.winConditionSystem.Check(r.State); hasWinner {
		r.State.GameStatus = "finished"
		r.State.Winner = &winnerID
		r.broadcastGameOver(winnerID, reason)
		// Use stopInternal since we already hold the lock, and capture the callback
		endCallback = r.stopInternal()
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

// updateBarracks updates all barracks (respawn timers, occupants, scatter)
func (r *GameRoom) updateBarracks(deltaTime float64) {
	for _, barracks := range r.State.Barracks {
		// Handle respawn timer
		barracks.Update(deltaTime)

		// Process any pending scatter from destruction
		scatteredIDs := barracks.GetAndClearPendingScatter()
		for _, unitID := range scatteredIDs {
			r.scatterInfantryFromBarracks(unitID, barracks.Position)
		}

		// Skip occupant tracking if barracks is destroyed
		if barracks.IsDestroyed {
			continue
		}

		// Track infantry inside this barracks
		for _, unit := range r.State.Units {
			if !unit.IsInfantry() {
				continue
			}
			if unit.GetHealth() <= 0 {
				continue
			}

			// Check if infantry is inside barracks
			if barracks.IsUnitInRange(unit.GetPosition()) {
				// Update occupant time and heal if ready
				healAmount := barracks.UpdateOccupant(unit.GetID(), deltaTime)
				if healAmount > 0 {
					unit.Heal(healAmount)
				}
			} else {
				// Remove from occupants if they left
				barracks.RemoveOccupant(unit.GetID())
			}
		}
	}
}

// scatterInfantryFromBarracks moves an infantry unit outside the barracks and applies damage
func (r *GameRoom) scatterInfantryFromBarracks(unitID string, barracksPos types.Vector3) {
	unit := r.State.GetUnitByID(unitID)
	if unit == nil || !unit.IsInfantry() {
		return
	}

	// Apply scatter damage
	unit.TakeDamage(types.BarracksScatterDamage)

	// Scatter in a random direction
	// Use unit ID hash for deterministic scatter direction
	hash := 0
	for _, c := range unitID {
		hash = hash*31 + int(c)
	}
	angle := float64(hash%360) * math.Pi / 180.0

	// Calculate new position outside barracks
	scatterDist := types.BarracksScatterRadius
	newX := barracksPos.X + math.Cos(angle)*scatterDist
	newZ := barracksPos.Z + math.Sin(angle)*scatterDist

	// Update unit position
	unit.SetPosition(types.Vector3{
		X: newX,
		Y: types.InfantryYPosition,
		Z: newZ,
	})
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
	case "super_tank":
		cost = types.SuperTankCost
	case "super_helicopter":
		cost = types.SuperHelicopterCost
	default:
		log.Printf("Unknown unit type: %s", unitType)
		return
	}

	// Check if player can afford
	if !player.CanAfford(cost) {
		if conn, ok := r.clientConnections[playerID]; ok {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "Not enough money",
			})
		}
		return
	}

	// Check super unit limit (only 1 super tank and 1 super helicopter per player)
	if unitType == "super_tank" || unitType == "super_helicopter" {
		if r.State.HasSuperUnit(playerID, unitType) {
			if conn, ok := r.clientConnections[playerID]; ok {
				unitName := "Super Tank"
				if unitType == "super_helicopter" {
					unitName = "Super Helicopter"
				}
				conn.SendMessage("error", types.ErrorPayload{
					Message: "You can only have one " + unitName + " at a time",
				})
			}
			return
		}
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
	case "super_tank":
		unit = NewSuperTank(playerID, spawnPos, targetPos)
	case "super_helicopter":
		unit = NewSuperHelicopter(playerID, spawnPos, targetPos)
	}

	// Add to state
	r.State.AddUnit(unit)
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

	// Check super unit limit (only 1 super tank and 1 super helicopter per player)
	if zone.UnitType == "super_tank" || zone.UnitType == "super_helicopter" {
		if r.State.HasSuperUnit(playerID, zone.UnitType) {
			unitName := "Super Tank"
			if zone.UnitType == "super_helicopter" {
				unitName = "Super Helicopter"
			}
			conn.SendMessage("error", types.ErrorPayload{
				Message: "You can only have one " + unitName + " at a time",
			})
			return
		}
	}

	// Deduct cost
	player.Spend(zone.Cost)

	// Queue the spawn instead of creating immediately
	spawnPos := zone.Position
	targetPos := r.State.Players[1-playerID].BasePosition // Target enemy base

	switch zone.UnitType {
	case "tank", "super_tank":
		spawnPos.Y = types.TankYPosition
	case "airplane", "super_helicopter":
		spawnPos.Y = types.AirplaneYPosition
	case "sniper", "rocket_launcher":
		spawnPos.Y = types.InfantryYPosition
		// Infantry spawn at closest owned barracks, or base if none owned
		if closestBarracks := r.getClosestOwnedBarracks(playerID, zone.Position); closestBarracks != nil {
			spawnPos = closestBarracks.Position
			spawnPos.Y = types.InfantryYPosition
		}
	}

	// Add to spawn queue
	r.State.SpawnQueue.Add(zone.UnitType, playerID, spawnPos, targetPos, zoneID)
}

// HandleBulkBuyFromZone handles a bulk purchase of 10 units at 10% discount
func (r *GameRoom) HandleBulkBuyFromZone(playerID int, zoneID string, conn ClientConnection) {
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

	// Bulk buy only available for regular tanks and helicopters, not super units
	if zone.UnitType != "tank" && zone.UnitType != "airplane" {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Bulk purchase only available for tanks and helicopters",
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

	// Calculate bulk price with 10% discount
	quantity := types.BulkBuyQuantity
	totalCost := int(float64(zone.Cost*quantity) * (1.0 - types.BulkBuyDiscount))

	// Check if player can afford
	if !player.CanAfford(totalCost) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: fmt.Sprintf("Not enough money! Need $%d for %d units", totalCost, quantity),
		})
		return
	}

	// Deduct cost
	player.Spend(totalCost)

	// Queue all spawns
	spawnPos := zone.Position
	targetPos := r.State.Players[1-playerID].BasePosition // Target enemy base

	switch zone.UnitType {
	case "tank":
		spawnPos.Y = types.TankYPosition
	case "airplane":
		spawnPos.Y = types.AirplaneYPosition
	}

	// Add all units to spawn queue
	for i := 0; i < quantity; i++ {
		r.State.SpawnQueue.Add(zone.UnitType, playerID, spawnPos, targetPos, zoneID)
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
		} else if turret.OwnerID == playerID {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "You already own this turret",
			})
		} else {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "Destroy enemy turret first",
			})
		}
		return
	}

	// Claim the turret
	turret.Claim(playerID)

	// Reward player for claiming turret
	player := r.State.GetPlayer(playerID)
	if player != nil {
		player.Money += types.TurretClaimReward
	}
}

// HandleClaimBuyZone handles a buy zone claiming request
func (r *GameRoom) HandleClaimBuyZone(playerID int, zoneID string, conn ClientConnection) {
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

	// Get the player unit
	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil || !playerUnit.IsAlive() {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "You must be alive to claim a base",
		})
		return
	}

	// Check if player is near the buy zone
	if !zone.IsPlayerInRange(playerUnit.GetPosition()) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Get closer to the base",
		})
		return
	}

	// Check if zone can be claimed
	if !zone.CanBeClaimed(playerID) {
		if !zone.IsClaimable {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "This base cannot be claimed",
			})
		} else if zone.OwnerID == playerID {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "You already own this base",
			})
		} else {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "This base is already owned",
			})
		}
		return
	}

	// Check if player has enough money to claim
	player := r.State.Players[playerID]
	if player.Money < zone.ClaimCost {
		conn.SendMessage("error", types.ErrorPayload{
			Message: fmt.Sprintf("Not enough money! Need $%d", zone.ClaimCost),
		})
		return
	}

	// Deduct the claim cost
	player.Money -= zone.ClaimCost

	// Claim the zone
	zone.Claim(playerID)

	// If this is a forward base, also claim all child zones
	if zone.UnitType == "" && zone.IsClaimable {
		for _, childZone := range r.State.BuyZones {
			if childZone.ForwardBaseID == zone.ID {
				childZone.Claim(playerID)
			}
		}
	}
}

// HandleClaimBarracks handles a barracks claiming request
func (r *GameRoom) HandleClaimBarracks(playerID int, barracksID string, conn ClientConnection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State.GameStatus != "playing" {
		return
	}

	// Find the barracks
	var barracks *Barracks
	for _, b := range r.State.Barracks {
		if b.ID == barracksID {
			barracks = b
			break
		}
	}

	if barracks == nil {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Invalid barracks",
		})
		return
	}

	// Get the player unit
	playerUnit := r.State.GetPlayerUnit(playerID)
	if playerUnit == nil || !playerUnit.IsAlive() {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "You must be alive to claim a barracks",
		})
		return
	}

	// Check if player is near the barracks
	if !barracks.IsUnitInRange(playerUnit.GetPosition()) {
		conn.SendMessage("error", types.ErrorPayload{
			Message: "Get closer to the barracks",
		})
		return
	}

	// Check if barracks can be claimed
	if !barracks.CanBeClaimed(playerID) {
		if barracks.IsDestroyed {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "Barracks is destroyed",
			})
		} else if barracks.OwnerID == playerID {
			conn.SendMessage("error", types.ErrorPayload{
				Message: "You already own this barracks",
			})
		}
		return
	}

	// Claim the barracks (free for infantry)
	barracks.Claim(playerID)
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

	// Check if there's a barracks at the target position (enemy barracks only)
	var targetBarracks *Barracks
	if targetUnit == nil && targetTurret == nil {
		for _, barracks := range r.State.Barracks {
			if barracks.OwnerID == playerID {
				continue // Don't shoot own barracks
			}
			if !barracks.IsAlive() {
				continue
			}

			distToTarget := calculateDistance(barracks.Position, targetPos)
			if distToTarget < 5.0 { // Close enough to target position
				targetBarracks = barracks
				break
			}
		}
	}

	// Create projectile
	var projectile *Projectile
	if targetTurret != nil {
		projectile = NewProjectileFromPlayerToTurret(playerUnit, targetTurret, now)
	} else if targetBarracks != nil {
		projectile = NewProjectileToBarracks(playerUnit, targetBarracks, now)
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
			Map:      r.State.MapDefinition,
		}
		conn.SendMessage("game_start", payload)
	}
}

// broadcastState sends the current game state to all players and spectators
func (r *GameRoom) broadcastState() {
	stateData := r.State.ToType()
	payload := types.GameUpdatePayload{
		Timestamp: r.State.Timestamp,
		State:     stateData,
	}

	// Send to players
	for _, conn := range r.clientConnections {
		conn.SendMessage("game_update", payload)
	}

	// Send to spectators
	for _, conn := range r.spectators {
		conn.SendMessage("game_update", payload)
	}
}

// broadcastGameOver sends the game_over message to all players
// Note: Caller must hold the lock (called from update() which holds r.mu.Lock())
func (r *GameRoom) broadcastGameOver(winner int, reason string) {
	p1Stats := r.State.Players[0].GetStats()
	p2Stats := r.State.Players[1].GetStats()
	matchDuration := r.State.GetMatchDuration()

	// Add win bonus points (200 for winning the round)
	const winBonus = 200
	if winner == 0 {
		p1Stats.TotalPoints += winBonus
	} else if winner == 1 {
		p2Stats.TotalPoints += winBonus
	}

	payload := types.GameOverPayload{
		Winner:        winner,
		Reason:        reason,
		MatchDuration: matchDuration,
		Stats: types.MatchStats{
			Player1Kills: r.State.Players[0].Kills,
			Player2Kills: r.State.Players[1].Kills,
			Player1Stats: p1Stats,
			Player2Stats: p2Stats,
		},
	}

	// Send to players
	for _, conn := range r.clientConnections {
		conn.SendMessage("game_over", payload)
	}

	// Send to spectators
	for _, conn := range r.spectators {
		conn.SendMessage("game_over", payload)
	}

	log.Printf("Game %s ended: Player %d won (%s) - Duration: %ds, P1: %d pts, P2: %d pts",
		r.ID, winner, reason, matchDuration,
		p1Stats.TotalPoints, p2Stats.TotalPoints)

	// Record to leaderboard using display names (GitHub username or Guest_XXXX)
	if r.onGameResult != nil {
		// Call callback outside lock - use goroutine to avoid blocking
		callback := r.onGameResult
		p1Name := r.State.Players[0].DisplayName
		p2Name := r.State.Players[1].DisplayName
		go callback(p1Name, p2Name, winner, matchDuration, p1Stats, p2Stats)
	}
}

// GetState returns a copy of the current game state (thread-safe)
func (r *GameRoom) GetState() types.GameState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.State.ToType()
}

// GetClientIDs returns the client IDs of both players
func (r *GameRoom) GetClientIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clientIDs := make([]string, 0, 2)
	for _, player := range r.State.Players {
		if player.ClientID != "" {
			clientIDs = append(clientIDs, player.ClientID)
		}
	}
	return clientIDs
}

// GetGameInfo returns summary information about the game for lobby display
func (r *GameRoom) GetGameInfo() types.ActiveGame {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return types.ActiveGame{
		GameID:         r.ID,
		Player1Name:    r.State.Players[0].DisplayName,
		Player2Name:    r.State.Players[1].DisplayName,
		SpectatorCount: len(r.spectators),
	}
}

// getClosestOwnedBarracks returns the closest barracks owned by the player, or nil if none owned
func (r *GameRoom) getClosestOwnedBarracks(playerID int, fromPos types.Vector3) *Barracks {
	var closest *Barracks
	closestDistSq := float64(0)

	for _, barracks := range r.State.Barracks {
		if barracks.OwnerID != playerID || barracks.IsDestroyed {
			continue
		}

		dx := barracks.Position.X - fromPos.X
		dz := barracks.Position.Z - fromPos.Z
		distSq := dx*dx + dz*dz

		if closest == nil || distSq < closestDistSq {
			closest = barracks
			closestDistSq = distSq
		}
	}

	return closest
}

// Marshal game room to JSON (for debugging)
func (r *GameRoom) MarshalJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return json.Marshal(struct {
		ID        string          `json:"id"`
		State     types.GameState `json:"state"`
		IsRunning bool            `json:"isRunning"`
	}{
		ID:        r.ID,
		State:     r.State.ToType(),
		IsRunning: r.IsRunning,
	})
}
