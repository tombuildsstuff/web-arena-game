package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// AIController manages AI decision-making for a computer player
type AIController struct {
	playerID       int
	difficulty     string // "easy", "medium", "hard"
	lastDecision   time.Time
	decisionDelay  time.Duration
	lastPurchase   time.Time
	purchaseDelay  time.Duration
	currentTarget  *types.Vector3
	targetUpdateAt time.Time
}

// NewAIController creates a new AI controller for a player
func NewAIController(playerID int, difficulty string) *AIController {
	ai := &AIController{
		playerID:      playerID,
		difficulty:    difficulty,
		lastDecision:  time.Now(),
		lastPurchase:  time.Now(),
		targetUpdateAt: time.Now(),
	}

	// Set delays based on difficulty
	switch difficulty {
	case "easy":
		ai.decisionDelay = 500 * time.Millisecond
		ai.purchaseDelay = 5 * time.Second
	case "hard":
		ai.decisionDelay = 100 * time.Millisecond
		ai.purchaseDelay = 2 * time.Second
	default: // medium
		ai.decisionDelay = 250 * time.Millisecond
		ai.purchaseDelay = 3 * time.Second
	}

	return ai
}

// Update runs AI decision-making for one tick
func (ai *AIController) Update(state *State, room *GameRoom) {
	now := time.Now()

	// Check if it's time to make decisions
	if now.Sub(ai.lastDecision) < ai.decisionDelay {
		return
	}
	ai.lastDecision = now

	// Get AI's player unit
	playerUnit := state.GetPlayerUnit(ai.playerID)
	if playerUnit == nil {
		return
	}

	// Make purchasing decisions
	ai.decidePurchases(state, room, now)

	// If dead, don't do movement/combat
	if !playerUnit.IsAlive() {
		return
	}

	// Decide on claiming turrets/zones
	ai.decideClaimActions(state, room, playerUnit)

	// Decide on shooting
	ai.decideShootAction(state, room, playerUnit)

	// Decide on movement
	ai.decideMovement(state, room, playerUnit)
}

// decidePurchases handles AI economy decisions
func (ai *AIController) decidePurchases(state *State, room *GameRoom, now time.Time) {
	if now.Sub(ai.lastPurchase) < ai.purchaseDelay {
		return
	}

	player := state.GetPlayer(ai.playerID)
	if player == nil {
		return
	}

	// Priority: Buy tanks if we can afford them
	if player.CanAfford(types.TankCost) {
		// Randomly choose between base and owned forward zones
		var zoneID string
		ownedZones := ai.getOwnedBuyZones(state)

		if len(ownedZones) > 0 && rand.Float32() < 0.5 {
			// Buy from a forward zone
			zone := ownedZones[rand.Intn(len(ownedZones))]
			zoneID = zone.ID
		}

		if zoneID != "" {
			// Buy from zone
			ai.buyFromZone(state, room, player, zoneID)
		} else {
			// Buy from base
			ai.purchaseFromBase(state, room, player, "tank")
		}
		ai.lastPurchase = now
	}
}

// purchaseFromBase purchases a unit from the player's base
func (ai *AIController) purchaseFromBase(state *State, room *GameRoom, player *Player, unitType string) {
	var cost int
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
		return
	}

	if !player.CanAfford(cost) {
		return
	}

	player.Spend(cost)

	spawnPos := player.BasePosition
	targetPos := state.Players[1-ai.playerID].BasePosition

	var unit Unit
	switch unitType {
	case "tank":
		spawnPos.Y = types.TankYPosition
		unit = NewTank(ai.playerID, spawnPos, targetPos)
	case "airplane":
		spawnPos.Y = types.AirplaneYPosition
		unit = NewAirplane(ai.playerID, spawnPos, targetPos)
	case "super_tank":
		spawnPos.Y = types.TankYPosition
		unit = NewSuperTank(ai.playerID, spawnPos, targetPos)
	case "super_helicopter":
		spawnPos.Y = types.AirplaneYPosition
		unit = NewSuperHelicopter(ai.playerID, spawnPos, targetPos)
	}

	if unit != nil {
		state.AddUnit(unit)
	}
}

// buyFromZone purchases from a forward buy zone
func (ai *AIController) buyFromZone(state *State, room *GameRoom, player *Player, zoneID string) {
	var zone *BuyZone
	for _, z := range state.BuyZones {
		if z.ID == zoneID && z.OwnerID == ai.playerID {
			zone = z
			break
		}
	}

	if zone == nil || !player.CanAfford(zone.Cost) {
		return
	}

	// Check super unit limit (only 1 super tank and 1 super helicopter per player)
	if zone.UnitType == "super_tank" || zone.UnitType == "super_helicopter" {
		if state.HasSuperUnit(ai.playerID, zone.UnitType) {
			return // AI already has this super unit type
		}
	}

	player.Spend(zone.Cost)

	spawnPos := zone.Position
	targetPos := state.Players[1-ai.playerID].BasePosition

	switch zone.UnitType {
	case "tank", "super_tank":
		spawnPos.Y = types.TankYPosition
	case "airplane", "super_helicopter":
		spawnPos.Y = types.AirplaneYPosition
	}

	state.SpawnQueue.Add(zone.UnitType, ai.playerID, spawnPos, targetPos, zoneID)
}

// getOwnedBuyZones returns buy zones owned by the AI (only unit purchase zones, not base zones)
func (ai *AIController) getOwnedBuyZones(state *State) []*BuyZone {
	var zones []*BuyZone
	for _, zone := range state.BuyZones {
		// Only include zones with a unit type (exclude base zones)
		if zone.OwnerID == ai.playerID && zone.UnitType != "" {
			zones = append(zones, zone)
		}
	}
	return zones
}

// decideClaimActions decides whether to claim nearby turrets or zones
func (ai *AIController) decideClaimActions(state *State, room *GameRoom, playerUnit *PlayerUnit) {
	pos := playerUnit.GetPosition()

	// Check for claimable turrets nearby
	for _, turret := range state.Turrets {
		if turret.CanBeClaimed(ai.playerID) && turret.IsPlayerInRange(pos) {
			turret.Claim(ai.playerID)
			// Reward AI for claiming turret
			player := state.GetPlayer(ai.playerID)
			if player != nil {
				player.Money += types.TurretClaimReward
			}
			return // One action per decision cycle
		}
	}

	// Check for claimable buy zones (forward bases) nearby
	player := state.GetPlayer(ai.playerID)
	if player == nil {
		return
	}

	for _, zone := range state.BuyZones {
		if zone.CanBeClaimed(ai.playerID) && zone.IsPlayerInRange(pos) {
			// Check if AI can afford the claim cost
			if player.Money >= zone.ClaimCost {
				player.Money -= zone.ClaimCost
				zone.Claim(ai.playerID)

				// If this is a forward base, also claim all child zones
				if zone.UnitType == "" && zone.IsClaimable {
					for _, childZone := range state.BuyZones {
						if childZone.ForwardBaseID == zone.ID {
							childZone.Claim(ai.playerID)
						}
					}
				}
				return
			}
		}
	}
}

// decideShootAction decides whether and what to shoot
func (ai *AIController) decideShootAction(state *State, room *GameRoom, playerUnit *PlayerUnit) {
	pos := playerUnit.GetPosition()
	attackRange := playerUnit.GetAttackRange()

	// Check attack cooldown
	now := time.Now().UnixMilli()
	timeSinceLastAttack := now - playerUnit.GetLastAttackTime()
	attackCooldown := int64(1000.0 / playerUnit.GetAttackSpeed())
	if timeSinceLastAttack < attackCooldown {
		return
	}

	// Find closest enemy unit in range
	var closestEnemy Unit
	closestDist := math.MaxFloat64

	for _, unit := range state.Units {
		if unit.GetOwnerID() == ai.playerID || !unit.IsAlive() {
			continue
		}

		dist := calculateDistance(pos, unit.GetPosition())
		if dist <= attackRange && dist < closestDist {
			// Check line of sight
			if room.losSystem.HasLineOfSight(pos, unit.GetPosition(), false) {
				closestEnemy = unit
				closestDist = dist
			}
		}
	}

	// Also check enemy turrets
	var closestTurret *Turret
	for _, turret := range state.Turrets {
		if turret.OwnerID == ai.playerID || !turret.IsAlive() {
			continue
		}

		dist := calculateDistance(pos, turret.Position)
		if dist <= attackRange && dist < closestDist {
			if room.losSystem.HasLineOfSight(pos, turret.Position, false) {
				closestTurret = turret
				closestDist = dist
				closestEnemy = nil // Prefer turret
			}
		}
	}

	// Shoot at target
	if closestTurret != nil {
		projectile := NewProjectileFromPlayerToTurret(playerUnit, closestTurret, now)
		state.AddProjectile(projectile)
		playerUnit.SetLastAttackTime(now)
	} else if closestEnemy != nil {
		// Add some inaccuracy based on difficulty
		targetPos := closestEnemy.GetPosition()
		if ai.difficulty == "easy" {
			targetPos.X += (rand.Float64() - 0.5) * 4
			targetPos.Z += (rand.Float64() - 0.5) * 4
		} else if ai.difficulty == "medium" {
			targetPos.X += (rand.Float64() - 0.5) * 2
			targetPos.Z += (rand.Float64() - 0.5) * 2
		}

		projectile := NewProjectileFromPlayer(playerUnit, targetPos, closestEnemy, now)
		state.AddProjectile(projectile)
		playerUnit.SetLastAttackTime(now)
	}
}

// decideMovement decides where the AI player should move
func (ai *AIController) decideMovement(state *State, room *GameRoom, playerUnit *PlayerUnit) {
	pos := playerUnit.GetPosition()
	now := time.Now()

	// Update target periodically or if no target
	if ai.currentTarget == nil || now.After(ai.targetUpdateAt) {
		ai.currentTarget = ai.selectMovementTarget(state, pos)
		ai.targetUpdateAt = now.Add(2 * time.Second)
	}

	if ai.currentTarget == nil {
		playerUnit.SetMoveDirection(types.Vector3{X: 0, Y: 0, Z: 0})
		return
	}

	// Check if reached target
	dist := calculateDistance(pos, *ai.currentTarget)
	if dist < 3.0 {
		ai.currentTarget = nil
		playerUnit.SetMoveDirection(types.Vector3{X: 0, Y: 0, Z: 0})
		return
	}

	// Move toward target
	dir := types.Vector3{
		X: ai.currentTarget.X - pos.X,
		Y: 0,
		Z: ai.currentTarget.Z - pos.Z,
	}

	// Normalize
	length := math.Sqrt(dir.X*dir.X + dir.Z*dir.Z)
	if length > 0 {
		dir.X /= length
		dir.Z /= length
	}

	playerUnit.SetMoveDirection(dir)
}

// selectMovementTarget chooses where the AI should move
func (ai *AIController) selectMovementTarget(state *State, currentPos types.Vector3) *types.Vector3 {
	// Priority 1: Claim nearby neutral turrets
	for _, turret := range state.Turrets {
		if turret.CanBeClaimed(ai.playerID) {
			dist := calculateDistance(currentPos, turret.Position)
			if dist < 50 { // Within reasonable distance
				return &turret.Position
			}
		}
	}

	// Priority 2: Claim nearby neutral buy zones
	for _, zone := range state.BuyZones {
		if zone.CanBeClaimed(ai.playerID) {
			dist := calculateDistance(currentPos, zone.Position)
			if dist < 50 {
				return &zone.Position
			}
		}
	}

	// Priority 3: Move toward enemy units to engage
	var closestEnemy Unit
	closestDist := math.MaxFloat64
	for _, unit := range state.Units {
		if unit.GetOwnerID() == ai.playerID || !unit.IsAlive() {
			continue
		}
		dist := calculateDistance(currentPos, unit.GetPosition())
		if dist < closestDist {
			closestEnemy = unit
			closestDist = dist
		}
	}

	if closestEnemy != nil && closestDist < 60 {
		enemyPos := closestEnemy.GetPosition()
		return &enemyPos
	}

	// Priority 4: Patrol around mid-field area
	// Random position in the middle area of the map
	target := types.Vector3{
		X: (rand.Float64() - 0.5) * 60, // -30 to 30
		Y: 0,
		Z: (rand.Float64() - 0.5) * 80, // -40 to 40
	}
	return &target
}

// AIClientConnection is a dummy connection for AI players
type AIClientConnection struct{}

// SendMessage discards messages for AI (AI doesn't need to receive state updates)
func (c *AIClientConnection) SendMessage(msgType string, payload interface{}) {
	// No-op - AI doesn't process incoming messages
}
