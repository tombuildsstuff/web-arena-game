package types

import "time"

const (
	// Arena settings
	ArenaSize     = 200
	ArenaHalfSize = ArenaSize / 2
	ArenaBoundary = ArenaHalfSize - 5 // Units should stay within this boundary
	BaseRadius    = 5

	// Game loop settings
	TickRate     = 20 // Updates per second
	TickDuration = time.Second / TickRate

	// Economy settings
	StartingMoney          = 1000
	PassiveIncomePerSecond = 10
	ForwardBaseClaimCost   = 500
	KillRewardUnit         = 10  // Money awarded for destroying a unit
	TurretClaimReward      = 20  // Money awarded for claiming a turret
	BulkBuyQuantity        = 10  // Number of units in a bulk purchase
	BulkBuyDiscount        = 0.1 // 10% discount for bulk purchases

	// Tank stats
	TankCost        = 50
	TankSpeed       = 5.0
	TankHealth      = 30 // 3 hits to destroy (30 / 10 = 3)
	TankDamage      = 10
	TankAttackRange = 10.0
	TankAttackSpeed = 1.0 // attacks per second

	// Airplane stats
	AirplaneCost        = 80
	AirplaneSpeed       = 15.0
	AirplaneHealth      = 30 // 3 hits to destroy (30 / 10 = 3)
	AirplaneDamage      = 10 // Same damage as tank for consistency
	AirplaneAttackRange = 20.0
	AirplaneAttackSpeed = 1.0 // attacks per second

	// Super Tank stats (2x firepower, 3x health)
	SuperTankCost        = 150
	SuperTankSpeed       = 4.0  // Slower than regular tank
	SuperTankHealth      = 90   // 3x regular tank health
	SuperTankDamage      = 20   // 2x regular tank damage
	SuperTankAttackRange = 12.0 // Slightly longer range
	SuperTankAttackSpeed = 1.0  // attacks per second

	// Super Helicopter stats (2x firepower, 3x health)
	SuperHelicopterCost        = 150
	SuperHelicopterSpeed       = 12.0 // Slower than regular airplane
	SuperHelicopterHealth      = 90   // 3x regular airplane health
	SuperHelicopterDamage      = 20   // 2x regular airplane damage
	SuperHelicopterAttackRange = 22.0 // Slightly longer range
	SuperHelicopterAttackSpeed = 1.0  // attacks per second

	// Sniper stats (infantry - low health, high range, high damage, slow fire)
	SniperCost            = 60
	SniperSpeed           = 4.0  // Slow movement
	SniperHealth          = 15   // Very fragile - low pain tolerance
	SniperDamage          = 25   // High damage per shot
	SniperAttackRange     = 35.0 // Very long range
	SniperAttackSpeed     = 0.5  // Slow fire rate (1 shot per 2 seconds)
	SniperCollisionRadius = 0.8  // Small - human sized

	// Rocket Launcher stats (infantry - low health, long range, very high damage, very slow fire)
	RocketLauncherCost            = 75
	RocketLauncherSpeed           = 3.5  // Very slow movement (heavy weapon)
	RocketLauncherHealth          = 20   // Fragile but slightly tougher than sniper
	RocketLauncherDamage          = 40   // Very high damage
	RocketLauncherAttackRange     = 30.0 // Long range
	RocketLauncherAttackSpeed     = 0.3  // Very slow fire rate (1 shot per 3+ seconds)
	RocketLauncherCollisionRadius = 0.8  // Small - human sized

	// Infantry positioning (same as tanks - ground level)
	InfantryYPosition = 1.0

	// Barracks stats
	BarracksHealth        = 40   // Takes more hits than turrets
	BarracksRespawnTime   = 30.0 // Seconds to respawn as neutral after destruction
	BarracksClaimRadius   = 8.0  // Proximity to claim (infantry only)
	BarracksHealTime      = 5.0  // Seconds infantry must stay inside to heal
	BarracksHealAmount    = 10   // Health restored per heal tick
	BarracksScatterDamage = 5    // Damage dealt to infantry when barracks is destroyed
	BarracksScatterRadius = 12.0 // How far infantry scatter when barracks is destroyed

	// Unit positioning
	TankYPosition     = 1.0
	AirplaneYPosition = 10.0

	// Unit collision radii (for unit-to-unit collision)
	TankCollisionRadius            = 2.0 // Tanks are ~3 units wide
	AirplaneCollisionRadius        = 1.5 // Airplanes are smaller
	SuperTankCollisionRadius       = 2.5 // Super tanks are bigger
	SuperHelicopterCollisionRadius = 2.0 // Super helicopters are bigger
	PlayerCollisionRadius          = 1.0 // Players are smallest

	// Player unit stats
	PlayerUnitSpeed       = 12.0
	PlayerUnitHealth      = 80 // 8 hits before respawn (80 / 10 damage = 8)
	PlayerUnitDamage      = 10
	PlayerUnitAttackRange = 25.0
	PlayerUnitAttackSpeed = 2.0 // attacks per second
	PlayerUnitYPosition   = 1.5
	PlayerRespawnTime     = 5.0 // seconds

	// Health Pack stats
	HealthPackHealAmount       = 20    // Health restored when collected
	HealthPackRadius           = 3.0   // Collection radius
	HealthPackSpawnMinSeconds  = 15    // Minimum seconds between spawns
	HealthPackSpawnMaxSeconds  = 45    // Maximum seconds between spawns
	HealthPackMaxCount         = 3     // Maximum health packs on map at once
	HealthPackLifetimeSeconds  = 60    // Seconds before uncollected pack despawns

	// Turret stats
	TurretHealth          = 30   // 3 hits to destroy (30 / 10 = 3)
	TurretDamage          = 5    // 2 hits to kill a player (10 / 5 = 2)
	TurretAttackRange     = 10.0
	TurretVsTurretRange   = 50.0 // Range for turret-to-turret combat (center turrets only)
	TurretAttackSpeed     = 0.5  // attacks per second
	TurretRespawnTime     = 10.0 // seconds
	TurretClaimRadius     = 15   // radius for claiming (3x3 squares)
	TurretTrackingTime    = 1500 // milliseconds to lock on before firing
)

var (
	Base1Position = Vector3{X: -90, Y: 0, Z: 0}
	Base2Position = Vector3{X: 90, Y: 0, Z: 0}

	PlayerColors = [2]string{"#3b82f6", "#ef4444"} // Blue and Red
)
