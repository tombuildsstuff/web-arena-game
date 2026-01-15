package maps

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// ClassicMap returns the default arena map
// This is the original map layout with symmetric design
func ClassicMap() *types.MapDefinition {
	return &types.MapDefinition{
		ID:          "classic",
		Name:        "Classic Arena",
		Description: "The original symmetric arena with forward bases and neutral turrets",

		// Arena dimensions
		ArenaSize:     200,
		ArenaBoundary: 95, // Units should stay within this boundary

		// Player configuration
		Players: []types.MapPlayerConfig{
			{
				BasePosition: types.Vector3{X: -90, Y: 0, Z: 0},
				SpawnOffset:  types.Vector3{X: 5, Y: 0, Z: 0},
				Color:        "#3b82f6", // Blue
			},
			{
				BasePosition: types.Vector3{X: 90, Y: 0, Z: 0},
				SpawnOffset:  types.Vector3{X: -5, Y: 0, Z: 0},
				Color:        "#ef4444", // Red
			},
		},

		// Buy zones
		BuyZones: classicBuyZones(),

		// Turrets
		Turrets: classicTurrets(),

		// Barracks (infantry spawn points)
		Barracks: classicBarracks(),

		// Obstacles
		Obstacles: classicObstacles(),

		// Health pack spawn bounds (avoiding bases)
		HealthPackSpawnBounds: types.MapBounds{
			MinX: -70,
			MaxX: 70,
			MinZ: -80,
			MaxZ: 80,
		},
	}
}

func classicBuyZones() []types.MapBuyZone {
	return []types.MapBuyZone{
		// Player 1 base buy zones (at X=-90)
		{
			ID:           "buy_p1_tank",
			DefaultOwner: 0,
			UnitType:     "tank",
			Position:     types.Vector3{X: -90, Y: 0, Z: -8},
			Radius:       4.0,
			Cost:         types.TankCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p1_airplane",
			DefaultOwner: 0,
			UnitType:     "airplane",
			Position:     types.Vector3{X: -90, Y: 0, Z: 8},
			Radius:       4.0,
			Cost:         types.AirplaneCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p1_sniper",
			DefaultOwner: 0,
			UnitType:     "sniper",
			Position:     types.Vector3{X: -95, Y: 0, Z: -4},
			Radius:       4.0,
			Cost:         types.SniperCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p1_rocket_launcher",
			DefaultOwner: 0,
			UnitType:     "rocket_launcher",
			Position:     types.Vector3{X: -95, Y: 0, Z: 4},
			Radius:       4.0,
			Cost:         types.RocketLauncherCost,
			IsClaimable:  false,
		},

		// Player 2 base buy zones (at X=90)
		{
			ID:           "buy_p2_tank",
			DefaultOwner: 1,
			UnitType:     "tank",
			Position:     types.Vector3{X: 90, Y: 0, Z: -8},
			Radius:       4.0,
			Cost:         types.TankCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p2_airplane",
			DefaultOwner: 1,
			UnitType:     "airplane",
			Position:     types.Vector3{X: 90, Y: 0, Z: 8},
			Radius:       4.0,
			Cost:         types.AirplaneCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p2_sniper",
			DefaultOwner: 1,
			UnitType:     "sniper",
			Position:     types.Vector3{X: 95, Y: 0, Z: -4},
			Radius:       4.0,
			Cost:         types.SniperCost,
			IsClaimable:  false,
		},
		{
			ID:           "buy_p2_rocket_launcher",
			DefaultOwner: 1,
			UnitType:     "rocket_launcher",
			Position:     types.Vector3{X: 95, Y: 0, Z: 4},
			Radius:       4.0,
			Cost:         types.RocketLauncherCost,
			IsClaimable:  false,
		},

		// North forward base (claimable)
		{
			ID:           "forward_base_north",
			DefaultOwner: -1,
			UnitType:     "",
			Position:     types.Vector3{X: 0, Y: 3, Z: -70},
			Radius:       12.0,
			Cost:         0,
			ClaimCost:    types.ForwardBaseClaimCost,
			IsClaimable:  true,
		},
		{
			ID:            "buy_forward_north_tank",
			DefaultOwner:  -1,
			UnitType:      "tank",
			Position:      types.Vector3{X: 0, Y: 3, Z: -70},
			Radius:        4.0,
			Cost:          types.TankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},
		{
			ID:            "buy_forward_north_super_tank",
			DefaultOwner:  -1,
			UnitType:      "super_tank",
			Position:      types.Vector3{X: -6, Y: 3, Z: -70},
			Radius:        4.0,
			Cost:          types.SuperTankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},
		{
			ID:            "buy_forward_north_super_helicopter",
			DefaultOwner:  -1,
			UnitType:      "super_helicopter",
			Position:      types.Vector3{X: 6, Y: 3, Z: -70},
			Radius:        4.0,
			Cost:          types.SuperHelicopterCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},

		// South forward base (claimable)
		{
			ID:           "forward_base_south",
			DefaultOwner: -1,
			UnitType:     "",
			Position:     types.Vector3{X: 0, Y: 3, Z: 70},
			Radius:       12.0,
			Cost:         0,
			ClaimCost:    types.ForwardBaseClaimCost,
			IsClaimable:  true,
		},
		{
			ID:            "buy_forward_south_tank",
			DefaultOwner:  -1,
			UnitType:      "tank",
			Position:      types.Vector3{X: 0, Y: 3, Z: 70},
			Radius:        4.0,
			Cost:          types.TankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
		{
			ID:            "buy_forward_south_super_tank",
			DefaultOwner:  -1,
			UnitType:      "super_tank",
			Position:      types.Vector3{X: -6, Y: 3, Z: 70},
			Radius:        4.0,
			Cost:          types.SuperTankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
		{
			ID:            "buy_forward_south_super_helicopter",
			DefaultOwner:  -1,
			UnitType:      "super_helicopter",
			Position:      types.Vector3{X: 6, Y: 3, Z: 70},
			Radius:        4.0,
			Cost:          types.SuperHelicopterCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
	}
}

func classicTurrets() []types.MapTurret {
	return []types.MapTurret{
		// Player 1 base turrets (4 corners)
		{ID: "turret_1", Position: types.Vector3{X: -77, Y: 3, Z: -15}, DefaultOwner: 0},
		{ID: "turret_2", Position: types.Vector3{X: -77, Y: 3, Z: 15}, DefaultOwner: 0},
		{ID: "turret_3", Position: types.Vector3{X: -98, Y: 3, Z: -15}, DefaultOwner: 0},
		{ID: "turret_4", Position: types.Vector3{X: -98, Y: 3, Z: 15}, DefaultOwner: 0},

		// Player 2 base turrets (4 corners)
		{ID: "turret_5", Position: types.Vector3{X: 77, Y: 3, Z: -15}, DefaultOwner: 1},
		{ID: "turret_6", Position: types.Vector3{X: 77, Y: 3, Z: 15}, DefaultOwner: 1},
		{ID: "turret_7", Position: types.Vector3{X: 98, Y: 3, Z: -15}, DefaultOwner: 1},
		{ID: "turret_8", Position: types.Vector3{X: 98, Y: 3, Z: 15}, DefaultOwner: 1},

		// Central corridor turrets (neutral)
		{ID: "turret_9", Position: types.Vector3{X: 0, Y: 3, Z: -50}, DefaultOwner: -1},
		{ID: "turret_10", Position: types.Vector3{X: 0, Y: 3, Z: 50}, DefaultOwner: -1},

		// Mid-lane turrets (neutral)
		{ID: "turret_11", Position: types.Vector3{X: -50, Y: 3, Z: -40}, DefaultOwner: -1},
		{ID: "turret_12", Position: types.Vector3{X: -50, Y: 3, Z: 40}, DefaultOwner: -1},
		{ID: "turret_13", Position: types.Vector3{X: 50, Y: 3, Z: -40}, DefaultOwner: -1},
		{ID: "turret_14", Position: types.Vector3{X: 50, Y: 3, Z: 40}, DefaultOwner: -1},

		// Center room turrets (neutral)
		{ID: "turret_15", Position: types.Vector3{X: -25, Y: 3, Z: 0}, DefaultOwner: -1},
		{ID: "turret_16", Position: types.Vector3{X: 25, Y: 3, Z: 0}, DefaultOwner: -1},

		// Corner turrets (neutral)
		{ID: "turret_17", Position: types.Vector3{X: -55, Y: 3, Z: -65}, DefaultOwner: -1},
		{ID: "turret_18", Position: types.Vector3{X: 55, Y: 3, Z: -65}, DefaultOwner: -1},
		{ID: "turret_19", Position: types.Vector3{X: -55, Y: 3, Z: 65}, DefaultOwner: -1},
		{ID: "turret_20", Position: types.Vector3{X: 55, Y: 3, Z: 65}, DefaultOwner: -1},
	}
}

func classicObstacles() []types.MapObstacle {
	obstacles := make([]types.MapObstacle, 0)

	wallHeight := 6.0
	wallThickness := 2.0
	pillarSize := 3.0
	pillarHeight := 6.0
	coverHeight := 3.0
	coverSize := 4.0
	platformHeight := 3.0
	platformSize := 20.0

	// Outer walls
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_1", Type: "wall", Position: types.Vector3{X: 0, Y: 0, Z: -85}, Size: types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_2", Type: "wall", Position: types.Vector3{X: 0, Y: 0, Z: 85}, Size: types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}},
	)

	// Upper lane divider segments (Z = -35)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_3", Type: "wall", Position: types.Vector3{X: -65, Y: 0, Z: -35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_4", Type: "wall", Position: types.Vector3{X: -25, Y: 0, Z: -35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_5", Type: "wall", Position: types.Vector3{X: 25, Y: 0, Z: -35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_6", Type: "wall", Position: types.Vector3{X: 65, Y: 0, Z: -35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
	)

	// Lower lane divider segments (Z = 35)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_7", Type: "wall", Position: types.Vector3{X: -65, Y: 0, Z: 35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_8", Type: "wall", Position: types.Vector3{X: -25, Y: 0, Z: 35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_9", Type: "wall", Position: types.Vector3{X: 25, Y: 0, Z: 35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_10", Type: "wall", Position: types.Vector3{X: 65, Y: 0, Z: 35}, Size: types.Vector3{X: 20, Y: wallHeight, Z: wallThickness}},
	)

	// Vertical cover walls (X = -45)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_11", Type: "wall", Position: types.Vector3{X: -45, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
		types.MapObstacle{ID: "obs_12", Type: "wall", Position: types.Vector3{X: -45, Y: 0, Z: 0}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_13", Type: "wall", Position: types.Vector3{X: -45, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
	)

	// Vertical cover walls (X = 45)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_14", Type: "wall", Position: types.Vector3{X: 45, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
		types.MapObstacle{ID: "obs_15", Type: "wall", Position: types.Vector3{X: 45, Y: 0, Z: 0}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_16", Type: "wall", Position: types.Vector3{X: 45, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
	)

	// Base approach walls (X = -70)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_17", Type: "wall", Position: types.Vector3{X: -70, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_18", Type: "wall", Position: types.Vector3{X: -70, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
	)

	// Base approach walls (X = 70)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_19", Type: "wall", Position: types.Vector3{X: 70, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_20", Type: "wall", Position: types.Vector3{X: 70, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
	)

	// Central corridor walls (X = -15)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_21", Type: "wall", Position: types.Vector3{X: -15, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
		types.MapObstacle{ID: "obs_22", Type: "wall", Position: types.Vector3{X: -15, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
	)

	// Central corridor walls (X = 15)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_23", Type: "wall", Position: types.Vector3{X: 15, Y: 0, Z: -60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
		types.MapObstacle{ID: "obs_24", Type: "wall", Position: types.Vector3{X: 15, Y: 0, Z: 60}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 15}},
	)

	// Central area pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_25", Type: "pillar", Position: types.Vector3{X: -10, Y: 0, Z: -8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_26", Type: "pillar", Position: types.Vector3{X: 0, Y: 0, Z: 0}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_27", Type: "pillar", Position: types.Vector3{X: 10, Y: 0, Z: 8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Top turret cover pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_28", Type: "pillar", Position: types.Vector3{X: -15, Y: 0, Z: -35}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_29", Type: "pillar", Position: types.Vector3{X: 0, Y: 0, Z: -30}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_30", Type: "pillar", Position: types.Vector3{X: 15, Y: 0, Z: -35}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Bottom turret cover pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_31", Type: "pillar", Position: types.Vector3{X: -15, Y: 0, Z: 35}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_32", Type: "pillar", Position: types.Vector3{X: 0, Y: 0, Z: 30}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_33", Type: "pillar", Position: types.Vector3{X: 15, Y: 0, Z: 35}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Top-left corner pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_34", Type: "pillar", Position: types.Vector3{X: -70, Y: 0, Z: -50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_35", Type: "pillar", Position: types.Vector3{X: -55, Y: 0, Z: -45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_36", Type: "pillar", Position: types.Vector3{X: -40, Y: 0, Z: -50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Top-right corner pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_37", Type: "pillar", Position: types.Vector3{X: 40, Y: 0, Z: -50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_38", Type: "pillar", Position: types.Vector3{X: 55, Y: 0, Z: -45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_39", Type: "pillar", Position: types.Vector3{X: 70, Y: 0, Z: -50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Bottom-left corner pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_40", Type: "pillar", Position: types.Vector3{X: -70, Y: 0, Z: 50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_41", Type: "pillar", Position: types.Vector3{X: -55, Y: 0, Z: 45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_42", Type: "pillar", Position: types.Vector3{X: -40, Y: 0, Z: 50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Bottom-right corner pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_43", Type: "pillar", Position: types.Vector3{X: 40, Y: 0, Z: 50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_44", Type: "pillar", Position: types.Vector3{X: 55, Y: 0, Z: 45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_45", Type: "pillar", Position: types.Vector3{X: 70, Y: 0, Z: 50}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Mid-lane cover pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_46", Type: "pillar", Position: types.Vector3{X: -35, Y: 0, Z: -25}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_47", Type: "pillar", Position: types.Vector3{X: -35, Y: 0, Z: 25}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_48", Type: "pillar", Position: types.Vector3{X: 35, Y: 0, Z: -25}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_49", Type: "pillar", Position: types.Vector3{X: 35, Y: 0, Z: 25}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Near-base cover blocks
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_50", Type: "cover", Position: types.Vector3{X: -82, Y: 0, Z: -40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_51", Type: "cover", Position: types.Vector3{X: -82, Y: 0, Z: 40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_52", Type: "cover", Position: types.Vector3{X: 82, Y: 0, Z: -40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_53", Type: "cover", Position: types.Vector3{X: 82, Y: 0, Z: 40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
	)

	// Mid-field cover blocks
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_54", Type: "cover", Position: types.Vector3{X: -55, Y: 0, Z: 0}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_55", Type: "cover", Position: types.Vector3{X: 55, Y: 0, Z: 0}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
	)

	// Forward base platforms (ramps)
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_56", Type: "ramp", Position: types.Vector3{X: 0, Y: 0, Z: -70}, Size: types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize}, ElevationStart: platformHeight, ElevationEnd: platformHeight},
		types.MapObstacle{ID: "obs_57", Type: "ramp", Position: types.Vector3{X: 0, Y: 0, Z: 70}, Size: types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize}, ElevationStart: platformHeight, ElevationEnd: platformHeight},
	)

	// Ramps to forward bases
	rampWidth := 14.0
	rampLength := 12.0
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_58", Type: "ramp", Position: types.Vector3{X: 0, Y: 0, Z: -54}, Size: types.Vector3{X: rampWidth, Y: platformHeight, Z: rampLength}, ElevationStart: platformHeight, ElevationEnd: 0},
		types.MapObstacle{ID: "obs_59", Type: "ramp", Position: types.Vector3{X: 0, Y: 0, Z: 54}, Size: types.Vector3{X: rampWidth, Y: platformHeight, Z: rampLength}, ElevationStart: 0, ElevationEnd: platformHeight},
	)

	// Platform edge barriers
	barrierHeight := 1.5
	barrierThickness := 1.0

	// North platform barriers
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_60", Type: "cover", Position: types.Vector3{X: -10, Y: platformHeight, Z: -70}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_61", Type: "cover", Position: types.Vector3{X: 10, Y: platformHeight, Z: -70}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_62", Type: "cover", Position: types.Vector3{X: 0, Y: platformHeight, Z: -80}, Size: types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}},
	)

	// South platform barriers
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_63", Type: "cover", Position: types.Vector3{X: -10, Y: platformHeight, Z: 70}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_64", Type: "cover", Position: types.Vector3{X: 10, Y: platformHeight, Z: 70}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_65", Type: "cover", Position: types.Vector3{X: 0, Y: platformHeight, Z: 80}, Size: types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}},
	)

	return obstacles
}

func classicBarracks() []types.MapBarracks {
	return []types.MapBarracks{
		// Mid-field barracks (neutral, key strategic points)
		{ID: "barracks_center_north", Position: types.Vector3{X: 0, Y: 0, Z: -40}},
		{ID: "barracks_center_south", Position: types.Vector3{X: 0, Y: 0, Z: 40}},

		// Side corridor barracks (flanking positions)
		{ID: "barracks_west_north", Position: types.Vector3{X: -50, Y: 0, Z: -55}},
		{ID: "barracks_west_south", Position: types.Vector3{X: -50, Y: 0, Z: 55}},
		{ID: "barracks_east_north", Position: types.Vector3{X: 50, Y: 0, Z: -55}},
		{ID: "barracks_east_south", Position: types.Vector3{X: 50, Y: 0, Z: 55}},
	}
}
