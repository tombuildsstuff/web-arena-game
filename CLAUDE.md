# Claude Code Context

This file provides context for Claude Code when working on this project.

## Project Overview

Arena Game is a real-time 3D multiplayer strategy game. Players control a soldier character while purchasing and deploying autonomous units (tanks, helicopters, infantry) to defeat opponents. The game uses a client-server architecture with authoritative server logic.

## Quick Commands

```bash
# Development (hot reload)
make run-dev
# Then open http://localhost:5173

# Production build and run
make build && make run
# Then open http://localhost:3000

# Run server tests
make test

# Format and lint Go code
cd server && make fmt && make lint
```

## Architecture

### Server (Go) - Authoritative
- **Location**: `server/`
- **Entry point**: `server/main.go`
- **Game tick rate**: 20 TPS (50ms per tick)
- **Key packages**:
  - `internal/game/` - Core game logic (rooms, AI, combat, spawning)
  - `internal/websocket/` - WebSocket hub and client management
  - `internal/types/` - Shared types, events, constants
  - `internal/auth/` - OAuth providers (GitHub, BlueSky)

### Client (JavaScript + Three.js) - Rendering Only
- **Location**: `client/`
- **Entry point**: `client/src/main.js`
- **Render rate**: 60 FPS
- **Key directories**:
  - `src/game/` - 3D rendering (units, arena, effects, projectiles)
  - `src/input/` - Player input handling (keyboard, mouse, touch)
  - `src/network/` - WebSocket client
  - `src/ui/` - HUD and menu systems
  - `src/state/` - Client-side game state
  - `src/audio/` - Sound effects manager

## Key Files

### Server
| File | Purpose |
|------|---------|
| `internal/game/room.go` | Game room logic, handles player actions, game loop |
| `internal/game/ai.go` | AI opponent behavior |
| `internal/game/spawn_queue.go` | Delayed unit spawning system |
| `internal/game/turret.go` | Turret placement and behavior |
| `internal/game/barracks.go` | Barracks (infantry spawn points) logic |
| `internal/game/sniper.go` | Sniper infantry unit behavior |
| `internal/game/rocket_launcher.go` | Rocket launcher infantry unit behavior |
| `internal/game/health_pack.go` | Health pack spawning |
| `internal/game/leaderboard.go` | Leaderboard persistence and total matches tracking |
| `internal/websocket/hub.go` | WebSocket message routing |
| `internal/types/constants.go` | Game balance constants (costs, damage, etc.) |
| `internal/types/events.go` | WebSocket message payload types |

### Client
| File | Purpose |
|------|---------|
| `src/game/Game.js` | Main game class, coordinates all systems |
| `src/game/GameLoop.js` | Render loop, entity updates, proximity checks |
| `src/game/Arena.js` | 3D arena rendering (ground, boundaries) |
| `src/game/Base.js` | Player base rendering |
| `src/game/BuyZone.js` | Buy zone rendering and highlighting |
| `src/game/Turret.js` | Turret 3D model and animations |
| `src/game/Barracks.js` | Barracks 3D model and claim state |
| `src/game/Sniper.js` | Sniper infantry 3D model and animations |
| `src/game/RocketLauncher.js` | Rocket launcher infantry 3D model |
| `src/input/PlayerInput.js` | Keyboard/mouse input handling |
| `src/input/MobileInput.js` | Touch controls for mobile |
| `src/ui/HUD.js` | Health bar, money, unit counts, prompts |
| `src/ui/Leaderboard.js` | Leaderboard display and total matches counter |
| `src/network/WebSocketClient.js` | Server communication |
| `src/state/GameState.js` | Client-side state management |

## Game Constants

Key values are defined in `server/internal/types/constants.go`:
- Unit costs, HP, damage, speed, attack rates
- Economy values (income rate, kill rewards)
- Spawn delays and intervals
- Game duration (5 minutes)

Client-side rendering constants are in `client/src/utils/constants.js`.

## WebSocket Messages

### Client to Server
| Type | Purpose |
|------|---------|
| `join_queue` | Join matchmaking queue |
| `start_vs_ai` | Start game against AI |
| `player_move` | Send movement direction |
| `player_shoot` | Fire at target position |
| `buy_from_zone` | Purchase single unit |
| `bulk_buy_from_zone` | Purchase 10 units at discount |
| `claim_turret` | Claim neutral turret |
| `claim_buy_zone` | Claim forward base |
| `leave_game` | Exit current game |

### Server to Client
| Type | Purpose |
|------|---------|
| `game_start` | Game begins, includes initial state |
| `game_state` | Full state update (20/sec) |
| `game_over` | Game ended, includes final stats |
| `error` | Error message to display |
| `lobby_status` | Queue size and active games |

## Coding Patterns

### Server-side
- Game rooms run in goroutines with their own tick loop
- Use channels for thread-safe communication
- All game state mutations happen in the room's main loop
- Constants for game balance are centralized in `types/constants.go`

### Client-side
- Entity pattern: each game object has a corresponding class (Tank, Helicopter, Sniper, RocketLauncher, Turret, Barracks, etc.)
- State is received from server; client never modifies authoritative state
- Three.js scene management through GameLoop
- DOM-based UI separate from 3D canvas

## Common Tasks

### Adding a new unit type
1. Add constants in `server/internal/types/constants.go`
2. Add spawn handling in `server/internal/game/room.go`
3. Create rendering class in `client/src/game/`
4. Update `GameLoop.js` to instantiate the new unit type
5. Add to HUD if needed

### Adding a new player action
1. Define payload type in `server/internal/types/events.go`
2. Add handler in `server/internal/websocket/hub.go`
3. Implement logic in `server/internal/game/room.go`
4. Add client-side trigger in `PlayerInput.js`
5. Add WebSocket send method in `Game.js`

### Adjusting game balance
- Modify values in `server/internal/types/constants.go`
- No client changes needed (server is authoritative)

## Testing

```bash
cd server
make test        # Run all tests
go test -v ./... # Verbose output
```

## Build Output

- Server binary: `server/bin/arena-server`
- Client bundle: `client/dist/` (embedded in server binary)
- Single deployable binary with embedded static files

## Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `PORT` | HTTP server port | 3000 |
| `GITHUB_CLIENT_ID` | GitHub OAuth app ID | - |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth secret | - |

## Notes

- The client is always embedded in the server binary for production
- Hot reload in dev mode: client on :5173, server on :3000 with proxy
- Mobile touch controls auto-enable on touch devices
- Leaderboard data persists in `server/leaderboard.json`
