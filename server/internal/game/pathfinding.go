package game

import (
	"container/heap"
	"math"

	"github.com/tombuildsstuff/web-arena-game/server/internal/types"
)

const (
	// PathGridCellSize is the size of each cell in the pathfinding grid
	PathGridCellSize = 4.0
	// UnitRadius is the radius to account for when checking walkability
	UnitRadius = 2.0
)

// PathNode represents a node in the A* search
type PathNode struct {
	X, Z   int     // Grid coordinates
	G      float64 // Cost from start
	H      float64 // Heuristic cost to goal
	F      float64 // G + H
	Parent *PathNode
	Index  int // Index in heap
}

// PathHeap implements heap.Interface for A* open list
type PathHeap []*PathNode

func (h PathHeap) Len() int           { return len(h) }
func (h PathHeap) Less(i, j int) bool { return h[i].F < h[j].F }
func (h PathHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *PathHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*PathNode)
	node.Index = n
	*h = append(*h, node)
}

func (h *PathHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.Index = -1
	*h = old[0 : n-1]
	return node
}

// PathGrid represents the navigation grid
type PathGrid struct {
	Width    int
	Height   int
	CellSize float64
	Walkable [][]bool
	OriginX  float64 // World X of grid origin
	OriginZ  float64 // World Z of grid origin
}

// PathfindingSystem handles pathfinding for units
type PathfindingSystem struct {
	Grid        *PathGrid
	SpatialGrid *SpatialGrid
}

// NewPathfindingSystem creates a new pathfinding system
func NewPathfindingSystem(obstacles []*Obstacle) *PathfindingSystem {
	spatialGrid := NewSpatialGrid(obstacles)

	// Create path grid covering the arena (200x200)
	arenaSize := 200.0
	gridSize := int(arenaSize / PathGridCellSize)

	grid := &PathGrid{
		Width:    gridSize,
		Height:   gridSize,
		CellSize: PathGridCellSize,
		Walkable: make([][]bool, gridSize),
		OriginX:  -arenaSize / 2,
		OriginZ:  -arenaSize / 2,
	}

	// Initialize walkability
	for x := 0; x < gridSize; x++ {
		grid.Walkable[x] = make([]bool, gridSize)
		for z := 0; z < gridSize; z++ {
			worldX := grid.OriginX + float64(x)*PathGridCellSize + PathGridCellSize/2
			worldZ := grid.OriginZ + float64(z)*PathGridCellSize + PathGridCellSize/2

			// Check if this cell is blocked
			grid.Walkable[x][z] = !spatialGrid.IsPositionBlocked(worldX, worldZ, UnitRadius)
		}
	}

	return &PathfindingSystem{
		Grid:        grid,
		SpatialGrid: spatialGrid,
	}
}

// WorldToGrid converts world coordinates to grid coordinates
func (g *PathGrid) WorldToGrid(worldX, worldZ float64) (int, int) {
	x := int((worldX - g.OriginX) / g.CellSize)
	z := int((worldZ - g.OriginZ) / g.CellSize)

	// Clamp to grid bounds
	if x < 0 {
		x = 0
	}
	if x >= g.Width {
		x = g.Width - 1
	}
	if z < 0 {
		z = 0
	}
	if z >= g.Height {
		z = g.Height - 1
	}

	return x, z
}

// GridToWorld converts grid coordinates to world coordinates (center of cell)
func (g *PathGrid) GridToWorld(gridX, gridZ int) (float64, float64) {
	worldX := g.OriginX + float64(gridX)*g.CellSize + g.CellSize/2
	worldZ := g.OriginZ + float64(gridZ)*g.CellSize + g.CellSize/2
	return worldX, worldZ
}

// IsWalkable checks if a grid cell is walkable
func (g *PathGrid) IsWalkable(x, z int) bool {
	if x < 0 || x >= g.Width || z < 0 || z >= g.Height {
		return false
	}
	return g.Walkable[x][z]
}

// FindPath finds a path from start to end using A*
func (ps *PathfindingSystem) FindPath(start, end types.Vector3) []types.Vector3 {
	startX, startZ := ps.Grid.WorldToGrid(start.X, start.Z)
	endX, endZ := ps.Grid.WorldToGrid(end.X, end.Z)

	// If start or end is blocked, find nearest walkable cell
	if !ps.Grid.IsWalkable(startX, startZ) {
		startX, startZ = ps.findNearestWalkable(startX, startZ)
	}
	if !ps.Grid.IsWalkable(endX, endZ) {
		endX, endZ = ps.findNearestWalkable(endX, endZ)
	}

	// A* search
	openList := &PathHeap{}
	heap.Init(openList)

	closedSet := make(map[GridKey]bool)
	nodeMap := make(map[GridKey]*PathNode)

	startNode := &PathNode{
		X: startX,
		Z: startZ,
		G: 0,
		H: ps.heuristic(startX, startZ, endX, endZ),
	}
	startNode.F = startNode.G + startNode.H
	heap.Push(openList, startNode)
	nodeMap[GridKey{X: startX, Z: startZ}] = startNode

	// Directions: 8-directional movement
	dirs := [][2]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Cardinal
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonal
	}

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*PathNode)
		currentKey := GridKey{X: current.X, Z: current.Z}

		// Check if we reached the goal
		if current.X == endX && current.Z == endZ {
			return ps.reconstructPath(current, start.Y, start, end)
		}

		closedSet[currentKey] = true

		// Check neighbors
		for _, dir := range dirs {
			nx, nz := current.X+dir[0], current.Z+dir[1]
			neighborKey := GridKey{X: nx, Z: nz}

			if !ps.Grid.IsWalkable(nx, nz) || closedSet[neighborKey] {
				continue
			}

			// Check diagonal movement (can't cut corners)
			if dir[0] != 0 && dir[1] != 0 {
				if !ps.Grid.IsWalkable(current.X+dir[0], current.Z) ||
					!ps.Grid.IsWalkable(current.X, current.Z+dir[1]) {
					continue
				}
			}

			// Calculate cost
			moveCost := 1.0
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = 1.414 // Diagonal cost
			}
			tentativeG := current.G + moveCost

			neighbor, exists := nodeMap[neighborKey]
			if !exists {
				neighbor = &PathNode{
					X: nx,
					Z: nz,
					G: math.MaxFloat64,
				}
				nodeMap[neighborKey] = neighbor
			}

			if tentativeG < neighbor.G {
				neighbor.Parent = current
				neighbor.G = tentativeG
				neighbor.H = ps.heuristic(nx, nz, endX, endZ)
				neighbor.F = neighbor.G + neighbor.H

				if neighbor.Index == -1 || neighbor.Index == 0 && !exists {
					heap.Push(openList, neighbor)
				} else {
					heap.Fix(openList, neighbor.Index)
				}
			}
		}
	}

	// No path found - return direct path
	return []types.Vector3{end}
}

// heuristic calculates the A* heuristic (diagonal distance)
func (ps *PathfindingSystem) heuristic(x1, z1, x2, z2 int) float64 {
	dx := math.Abs(float64(x2 - x1))
	dz := math.Abs(float64(z2 - z1))
	return dx + dz + (1.414-2)*math.Min(dx, dz)
}

// reconstructPath builds the path from the goal node back to start
func (ps *PathfindingSystem) reconstructPath(node *PathNode, y float64, actualStart, actualEnd types.Vector3) []types.Vector3 {
	path := make([]types.Vector3, 0)

	for node != nil {
		worldX, worldZ := ps.Grid.GridToWorld(node.X, node.Z)
		path = append([]types.Vector3{{X: worldX, Y: y, Z: worldZ}}, path...)
		node = node.Parent
	}

	// Replace first and last waypoints with actual positions
	if len(path) > 0 {
		path[0] = actualStart
		path[len(path)-1] = actualEnd
	}

	// Smooth the path
	return ps.smoothPath(path)
}

// smoothPath removes unnecessary waypoints
func (ps *PathfindingSystem) smoothPath(path []types.Vector3) []types.Vector3 {
	if len(path) <= 2 {
		return path
	}

	smoothed := []types.Vector3{path[0]}

	for i := 1; i < len(path)-1; i++ {
		// Check if we can skip this waypoint
		if !ps.hasDirectPath(smoothed[len(smoothed)-1], path[i+1]) {
			smoothed = append(smoothed, path[i])
		}
	}

	smoothed = append(smoothed, path[len(path)-1])
	return smoothed
}

// hasDirectPath checks if there's a clear line between two points
func (ps *PathfindingSystem) hasDirectPath(start, end types.Vector3) bool {
	dx := end.X - start.X
	dz := end.Z - start.Z
	length := math.Sqrt(dx*dx + dz*dz)

	if length == 0 {
		return true
	}

	// Step along the line checking for obstacles
	stepSize := PathGridCellSize / 2
	steps := int(length/stepSize) + 1

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := start.X + dx*t
		z := start.Z + dz*t

		if ps.SpatialGrid.IsPositionBlocked(x, z, UnitRadius) {
			return false
		}
	}

	return true
}

// findNearestWalkable finds the nearest walkable cell to the given position
func (ps *PathfindingSystem) findNearestWalkable(x, z int) (int, int) {
	// Spiral outward search
	for radius := 1; radius < 20; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dz := -radius; dz <= radius; dz++ {
				if math.Abs(float64(dx)) == float64(radius) || math.Abs(float64(dz)) == float64(radius) {
					nx, nz := x+dx, z+dz
					if ps.Grid.IsWalkable(nx, nz) {
						return nx, nz
					}
				}
			}
		}
	}
	return x, z
}

// CheckLineOfSight checks if there's line of sight between two positions
func (ps *PathfindingSystem) CheckLineOfSight(from, to types.Vector3) bool {
	return ps.hasDirectPath(from, to)
}
