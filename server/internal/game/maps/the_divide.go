package maps

import (
	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

// TheDivideMap returns a map with a massive central ridge
// Players must fight for control of the high ground or flank around
func TheDivideMap() *types.MapDefinition {
	return &types.MapDefinition{
		ID:          "the_divide",
		Name:        "The Divide",
		Description: "Two elevated forward bases face each other across a central battleground. Fight for control of the high ground.",

		// Arena dimensions
		ArenaSize:     200,
		ArenaBoundary: 95,

		// Player configuration - bases on opposite ends
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
		BuyZones: divideBuyZones(),

		// Turrets
		Turrets: divideTurrets(),

		// Barracks (infantry spawn points)
		Barracks: divideBarracks(),

		// Obstacles
		Obstacles: divideObstacles(),

		// Health pack spawn bounds
		HealthPackSpawnBounds: types.MapBounds{
			MinX: -70,
			MaxX: 70,
			MinZ: -80,
			MaxZ: 80,
		},
	}
}

func divideBuyZones() []types.MapBuyZone {
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

		// North forward base - raised platform facing south (toward center)
		{
			ID:           "forward_base_north",
			DefaultOwner: -1,
			UnitType:     "",
			Position:     types.Vector3{X: 0, Y: 4, Z: -55},
			Radius:       12.0,
			Cost:         0,
			ClaimCost:    types.ForwardBaseClaimCost,
			IsClaimable:  true,
		},
		{
			ID:            "buy_forward_north_tank",
			DefaultOwner:  -1,
			UnitType:      "tank",
			Position:      types.Vector3{X: 0, Y: 4, Z: -55},
			Radius:        4.0,
			Cost:          types.TankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},
		{
			ID:            "buy_forward_north_super_tank",
			DefaultOwner:  -1,
			UnitType:      "super_tank",
			Position:      types.Vector3{X: -6, Y: 4, Z: -55},
			Radius:        4.0,
			Cost:          types.SuperTankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},
		{
			ID:            "buy_forward_north_super_helicopter",
			DefaultOwner:  -1,
			UnitType:      "super_helicopter",
			Position:      types.Vector3{X: 6, Y: 4, Z: -55},
			Radius:        4.0,
			Cost:          types.SuperHelicopterCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_north",
		},

		// South forward base - raised platform facing north (toward center)
		{
			ID:           "forward_base_south",
			DefaultOwner: -1,
			UnitType:     "",
			Position:     types.Vector3{X: 0, Y: 4, Z: 55},
			Radius:       12.0,
			Cost:         0,
			ClaimCost:    types.ForwardBaseClaimCost,
			IsClaimable:  true,
		},
		{
			ID:            "buy_forward_south_tank",
			DefaultOwner:  -1,
			UnitType:      "tank",
			Position:      types.Vector3{X: 0, Y: 4, Z: 55},
			Radius:        4.0,
			Cost:          types.TankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
		{
			ID:            "buy_forward_south_super_tank",
			DefaultOwner:  -1,
			UnitType:      "super_tank",
			Position:      types.Vector3{X: -6, Y: 4, Z: 55},
			Radius:        4.0,
			Cost:          types.SuperTankCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
		{
			ID:            "buy_forward_south_super_helicopter",
			DefaultOwner:  -1,
			UnitType:      "super_helicopter",
			Position:      types.Vector3{X: 6, Y: 4, Z: 55},
			Radius:        4.0,
			Cost:          types.SuperHelicopterCost,
			IsClaimable:   false,
			ForwardBaseID: "forward_base_south",
		},
	}
}

func divideTurrets() []types.MapTurret {
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

		// === ALL NEUTRAL TURRETS ON GROUND LEVEL ===

		// Central corridor - main battle line (ground level)
		{ID: "turret_9", Position: types.Vector3{X: -40, Y: 3, Z: 0}, DefaultOwner: -1},
		{ID: "turret_10", Position: types.Vector3{X: -20, Y: 3, Z: 0}, DefaultOwner: -1},
		{ID: "turret_11", Position: types.Vector3{X: 0, Y: 3, Z: 0}, DefaultOwner: -1},
		{ID: "turret_12", Position: types.Vector3{X: 20, Y: 3, Z: 0}, DefaultOwner: -1},
		{ID: "turret_13", Position: types.Vector3{X: 40, Y: 3, Z: 0}, DefaultOwner: -1},

		// North platform approach turrets - ground level below the platform
		{ID: "turret_14", Position: types.Vector3{X: -30, Y: 3, Z: -40}, DefaultOwner: -1},
		{ID: "turret_15", Position: types.Vector3{X: 0, Y: 3, Z: -40}, DefaultOwner: -1},
		{ID: "turret_16", Position: types.Vector3{X: 30, Y: 3, Z: -40}, DefaultOwner: -1},

		// South platform approach turrets - ground level below the platform
		{ID: "turret_17", Position: types.Vector3{X: -30, Y: 3, Z: 40}, DefaultOwner: -1},
		{ID: "turret_18", Position: types.Vector3{X: 0, Y: 3, Z: 40}, DefaultOwner: -1},
		{ID: "turret_19", Position: types.Vector3{X: 30, Y: 3, Z: 40}, DefaultOwner: -1},

		// Outer flanking turrets - west side
		{ID: "turret_20", Position: types.Vector3{X: -60, Y: 3, Z: -35}, DefaultOwner: -1},
		{ID: "turret_21", Position: types.Vector3{X: -60, Y: 3, Z: 35}, DefaultOwner: -1},

		// Outer flanking turrets - east side
		{ID: "turret_22", Position: types.Vector3{X: 60, Y: 3, Z: -35}, DefaultOwner: -1},
		{ID: "turret_23", Position: types.Vector3{X: 60, Y: 3, Z: 35}, DefaultOwner: -1},

		// Corner turrets - cover the far edges
		{ID: "turret_24", Position: types.Vector3{X: -55, Y: 3, Z: -65}, DefaultOwner: -1},
		{ID: "turret_25", Position: types.Vector3{X: 55, Y: 3, Z: -65}, DefaultOwner: -1},
		{ID: "turret_26", Position: types.Vector3{X: -55, Y: 3, Z: 65}, DefaultOwner: -1},
		{ID: "turret_27", Position: types.Vector3{X: 55, Y: 3, Z: 65}, DefaultOwner: -1},

		// Mid-lane turrets - between center and bases (moved off the walls)
		{ID: "turret_28", Position: types.Vector3{X: -50, Y: 3, Z: -15}, DefaultOwner: -1},
		{ID: "turret_29", Position: types.Vector3{X: -50, Y: 3, Z: 15}, DefaultOwner: -1},
		{ID: "turret_30", Position: types.Vector3{X: 50, Y: 3, Z: -15}, DefaultOwner: -1},
		{ID: "turret_31", Position: types.Vector3{X: 50, Y: 3, Z: 15}, DefaultOwner: -1},
	}
}

func divideObstacles() []types.MapObstacle {
	obstacles := make([]types.MapObstacle, 0)

	wallHeight := 6.0
	wallThickness := 2.0
	pillarSize := 3.0
	pillarHeight := 6.0
	coverHeight := 3.0
	coverSize := 4.0
	platformHeight := 4.0
	platformSize := 20.0

	// ============================================================
	// OUTER WALLS
	// ============================================================
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_1", Type: "wall", Position: types.Vector3{X: 0, Y: 0, Z: -85}, Size: types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_2", Type: "wall", Position: types.Vector3{X: 0, Y: 0, Z: 85}, Size: types.Vector3{X: 180, Y: wallHeight, Z: wallThickness}},
	)

	// ============================================================
	// NORTH FORWARD BASE PLATFORM - Raised platform facing south
	// ============================================================

	// North platform (flat elevated area)
	obstacles = append(obstacles,
		types.MapObstacle{
			ID:             "obs_3",
			Type:           "ramp",
			Position:       types.Vector3{X: 0, Y: 0, Z: -55},
			Size:           types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize},
			ElevationStart: platformHeight,
			ElevationEnd:   platformHeight,
		},
	)

	// North platform ramp (approaching from center/south)
	obstacles = append(obstacles,
		types.MapObstacle{
			ID:             "obs_4",
			Type:           "ramp",
			Position:       types.Vector3{X: 0, Y: 0, Z: -39},
			Size:           types.Vector3{X: 14, Y: platformHeight, Z: 12},
			ElevationStart: platformHeight, // Top at minZ (toward platform)
			ElevationEnd:   0,              // Bottom at maxZ (toward center)
		},
	)

	// North platform barriers
	barrierHeight := 1.5
	barrierThickness := 1.0
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_5", Type: "cover", Position: types.Vector3{X: -10, Y: platformHeight, Z: -55}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_6", Type: "cover", Position: types.Vector3{X: 10, Y: platformHeight, Z: -55}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_7", Type: "cover", Position: types.Vector3{X: 0, Y: platformHeight, Z: -65}, Size: types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}},
	)

	// ============================================================
	// SOUTH FORWARD BASE PLATFORM - Raised platform facing north
	// ============================================================

	// South platform (flat elevated area)
	obstacles = append(obstacles,
		types.MapObstacle{
			ID:             "obs_8",
			Type:           "ramp",
			Position:       types.Vector3{X: 0, Y: 0, Z: 55},
			Size:           types.Vector3{X: platformSize, Y: platformHeight, Z: platformSize},
			ElevationStart: platformHeight,
			ElevationEnd:   platformHeight,
		},
	)

	// South platform ramp (approaching from center/north)
	obstacles = append(obstacles,
		types.MapObstacle{
			ID:             "obs_9",
			Type:           "ramp",
			Position:       types.Vector3{X: 0, Y: 0, Z: 39},
			Size:           types.Vector3{X: 14, Y: platformHeight, Z: 12},
			ElevationStart: 0,              // Bottom at minZ (toward center)
			ElevationEnd:   platformHeight, // Top at maxZ (toward platform)
		},
	)

	// South platform barriers
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_10", Type: "cover", Position: types.Vector3{X: -10, Y: platformHeight, Z: 55}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_11", Type: "cover", Position: types.Vector3{X: 10, Y: platformHeight, Z: 55}, Size: types.Vector3{X: barrierThickness, Y: barrierHeight, Z: platformSize}},
		types.MapObstacle{ID: "obs_12", Type: "cover", Position: types.Vector3{X: 0, Y: platformHeight, Z: 65}, Size: types.Vector3{X: platformSize, Y: barrierHeight, Z: barrierThickness}},
	)

	// ============================================================
	// CENTRAL BATTLE ZONE - The "divide" between north and south
	// Walls creating lanes and chokepoints
	// ============================================================

	// Central horizontal walls creating the divide
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_13", Type: "wall", Position: types.Vector3{X: -50, Y: 0, Z: -20}, Size: types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_14", Type: "wall", Position: types.Vector3{X: 50, Y: 0, Z: -20}, Size: types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_15", Type: "wall", Position: types.Vector3{X: -50, Y: 0, Z: 20}, Size: types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}},
		types.MapObstacle{ID: "obs_16", Type: "wall", Position: types.Vector3{X: 50, Y: 0, Z: 20}, Size: types.Vector3{X: 30, Y: wallHeight, Z: wallThickness}},
	)

	// Central vertical walls creating corridors
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_17", Type: "wall", Position: types.Vector3{X: -25, Y: 0, Z: 0}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}},
		types.MapObstacle{ID: "obs_18", Type: "wall", Position: types.Vector3{X: 25, Y: 0, Z: 0}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 25}},
	)

	// ============================================================
	// BASE APPROACH WALLS
	// ============================================================

	// Near P1 base
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_19", Type: "wall", Position: types.Vector3{X: -70, Y: 0, Z: -50}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_20", Type: "wall", Position: types.Vector3{X: -70, Y: 0, Z: 50}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
	)

	// Near P2 base
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_21", Type: "wall", Position: types.Vector3{X: 70, Y: 0, Z: -50}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
		types.MapObstacle{ID: "obs_22", Type: "wall", Position: types.Vector3{X: 70, Y: 0, Z: 50}, Size: types.Vector3{X: wallThickness, Y: wallHeight, Z: 20}},
	)

	// ============================================================
	// COVER BLOCKS AND PILLARS
	// ============================================================

	// Near-base cover
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_23", Type: "cover", Position: types.Vector3{X: -82, Y: 0, Z: -40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_24", Type: "cover", Position: types.Vector3{X: -82, Y: 0, Z: 40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_25", Type: "cover", Position: types.Vector3{X: 82, Y: 0, Z: -40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
		types.MapObstacle{ID: "obs_26", Type: "cover", Position: types.Vector3{X: 82, Y: 0, Z: 40}, Size: types.Vector3{X: coverSize, Y: coverHeight, Z: coverSize}},
	)

	// Central area pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_27", Type: "pillar", Position: types.Vector3{X: -10, Y: 0, Z: -8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_28", Type: "pillar", Position: types.Vector3{X: 10, Y: 0, Z: -8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_29", Type: "pillar", Position: types.Vector3{X: -10, Y: 0, Z: 8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_30", Type: "pillar", Position: types.Vector3{X: 10, Y: 0, Z: 8}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Mid-field pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_31", Type: "pillar", Position: types.Vector3{X: -45, Y: 0, Z: -45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_32", Type: "pillar", Position: types.Vector3{X: 45, Y: 0, Z: -45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_33", Type: "pillar", Position: types.Vector3{X: -45, Y: 0, Z: 45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_34", Type: "pillar", Position: types.Vector3{X: 45, Y: 0, Z: 45}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Corner pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_35", Type: "pillar", Position: types.Vector3{X: -70, Y: 0, Z: -70}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_36", Type: "pillar", Position: types.Vector3{X: 70, Y: 0, Z: -70}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_37", Type: "pillar", Position: types.Vector3{X: -70, Y: 0, Z: 70}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_38", Type: "pillar", Position: types.Vector3{X: 70, Y: 0, Z: 70}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	// Flanking route pillars
	obstacles = append(obstacles,
		types.MapObstacle{ID: "obs_39", Type: "pillar", Position: types.Vector3{X: -55, Y: 0, Z: 0}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
		types.MapObstacle{ID: "obs_40", Type: "pillar", Position: types.Vector3{X: 55, Y: 0, Z: 0}, Size: types.Vector3{X: pillarSize, Y: pillarHeight, Z: pillarSize}},
	)

	return obstacles
}

func divideBarracks() []types.MapBarracks {
	return []types.MapBarracks{
		// Central barracks - key strategic point between the two platforms
		{ID: "barracks_center", Position: types.Vector3{X: 0, Y: 0, Z: 0}},

		// Flanking barracks - west side
		{ID: "barracks_west_north", Position: types.Vector3{X: -55, Y: 0, Z: -35}},
		{ID: "barracks_west_south", Position: types.Vector3{X: -55, Y: 0, Z: 35}},

		// Flanking barracks - east side
		{ID: "barracks_east_north", Position: types.Vector3{X: 55, Y: 0, Z: -35}},
		{ID: "barracks_east_south", Position: types.Vector3{X: 55, Y: 0, Z: 35}},

		// Near-platform barracks (below the raised forward bases)
		{ID: "barracks_north_approach", Position: types.Vector3{X: 0, Y: 0, Z: -30}},
		{ID: "barracks_south_approach", Position: types.Vector3{X: 0, Y: 0, Z: 30}},
	}
}
