package types

import "time"

const (
	// Arena settings
	ArenaSize       = 200
	ArenaHalfSize   = ArenaSize / 2
	ArenaBoundary   = ArenaHalfSize - 5 // Units should stay within this boundary
	BaseRadius      = 5

	// Game loop settings
	TickRate     = 20 // Updates per second
	TickDuration = time.Second / TickRate

	// Economy settings
	StartingMoney          = 500
	PassiveIncomePerSecond = 10

	// Tank stats
	TankCost         = 50
	TankSpeed        = 5.0
	TankHealth       = 30  // 3 hits to destroy (30 / 10 = 3)
	TankDamage       = 10
	TankAttackRange  = 10.0
	TankAttackSpeed  = 1.0 // attacks per second

	// Airplane stats
	AirplaneCost        = 80
	AirplaneSpeed       = 15.0
	AirplaneHealth      = 30  // 3 hits to destroy (30 / 10 = 3)
	AirplaneDamage      = 10  // Same damage as tank for consistency
	AirplaneAttackRange = 20.0
	AirplaneAttackSpeed = 1.0 // attacks per second

	// Unit positioning
	TankYPosition     = 1.0
	AirplaneYPosition = 10.0

	// Unit collision radii (for unit-to-unit collision)
	TankCollisionRadius     = 2.0  // Tanks are ~3 units wide
	AirplaneCollisionRadius = 1.5  // Airplanes are smaller
	PlayerCollisionRadius   = 1.0  // Players are smallest

	// Player unit stats
	PlayerUnitSpeed       = 12.0
	PlayerUnitHealth      = 5   // 5 hits before respawn
	PlayerUnitDamage      = 10
	PlayerUnitAttackRange = 25.0
	PlayerUnitAttackSpeed = 2.0  // attacks per second
	PlayerUnitYPosition   = 1.5
	PlayerRespawnTime     = 10.0 // seconds

	// Turret stats
	TurretHealth       = 30   // 3 hits to destroy (30 / 10 = 3)
	TurretDamage       = 10
	TurretAttackRange  = 20.0 // ~4 squares, requires line of sight
	TurretAttackSpeed  = 1.5  // attacks per second
	TurretRespawnTime  = 10.0 // seconds
	TurretClaimRadius  = 7.5  // radius for claiming (3x3 squares)
	TurretTrackingTime = 1000 // milliseconds to lock on before firing
)

var (
	Base1Position = Vector3{X: -90, Y: 0, Z: 0}
	Base2Position = Vector3{X: 90, Y: 0, Z: 0}

	PlayerColors = [2]string{"#3b82f6", "#ef4444"} // Blue and Red
)
