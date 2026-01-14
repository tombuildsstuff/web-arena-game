# Arena Game

A real-time 3D multiplayer strategy game where two players compete by purchasing and deploying tanks and airplanes to reach the opponent's base.

## Features

- **Real-time Multiplayer**: Play against another player with live updates via WebSocket
- **3D Graphics**: Built with Three.js for immersive 3D visualization
- **Strategic Gameplay**:
  - Purchase tanks ($50) and airplanes ($80)
  - Earn passive income: $10 per second
  - Units automatically move toward enemy base
  - Units engage in combat when in range
  - **Win condition**: Get a tank to the enemy base (airplanes cannot win!)
- **Modern Stack**: Vanilla JavaScript client + Go server

## Tech Stack

### Client
- **Vanilla JavaScript** with ES6+ modules
- **Three.js** for 3D rendering
- **Vite** for build tooling and dev server
- **Bun** as the JavaScript runtime

### Server
- **Go 1.21+** for high-performance game logic
- **gorilla/websocket** for WebSocket connections
- **chi** router for HTTP handling
- Authoritative server architecture (20 TPS)

## Project Structure

```
web-arena-game/
├── client/              # Frontend (Vanilla JS + Three.js)
│   ├── src/
│   │   ├── game/       # 3D game rendering
│   │   ├── network/    # WebSocket client
│   │   ├── ui/         # DOM-based UI
│   │   ├── state/      # Client state management
│   │   └── utils/      # Utilities and constants
│   ├── index.html
│   └── package.json
│
└── server/             # Backend (Go)
    ├── main.go         # Server entry point
    ├── internal/
    │   ├── game/       # Game logic and systems
    │   ├── websocket/  # WebSocket handling
    │   └── types/      # Shared types and constants
    └── go.mod
```

## Getting Started

### Prerequisites

- **Go 1.21+**: [Download Go](https://go.dev/dl/)
- **Bun**: [Install Bun](https://bun.sh/)
- **Make**: Usually pre-installed on macOS/Linux. For Windows, use WSL or install GNU Make.

### Quick Start (Recommended)

The client is **always embedded** in the server binary for simplicity:

```bash
# Install all dependencies
make tools

# Build server (client auto-embedded)
make build

# Run the server
make run
```

Then open two browser windows at `http://localhost:3000`

**For development with hot reload:**

```bash
# Run both servers with hot reload enabled
make run-dev
```

Then open two browser windows at `http://localhost:5173`

**Architecture Benefits:**
- ✅ Single binary deployment - client always embedded
- ✅ Simplified build process - one command does everything
- ✅ No coordination needed between frontend and backend
- ✅ WebSocket and HTTP routes on the same server
- ✅ Portable standalone binary

### Manual Setup

### Makefile Targets

**Root-level** (`make <target>` from project root):
- `make tools` - Install all tools and dependencies
- `make build` - Build server (client auto-embedded)
- `make run` - Build and run server
- `make run-dev` - Run dev servers with hot reload
- `make clean` - Clean all build artifacts
- `make test` - Run server tests
- `make help` - Show all available commands

**Server-specific** (`cd server && make <target>`):
- `make tools` - Install Go development tools
- `make build` - Build server with embedded client
- `make run` - Build and run the server
- `make test` - Run Go tests
- `make fmt` - Format Go code
- `make lint` - Run golangci-lint
- `make clean` - Remove build artifacts

**Client-specific** (`cd client && make <target>`):
- `make tools` - Install Bun dependencies
- `make build` - Build production bundle
- `make run` - Run Vite dev server
- `make preview` - Preview production build
- `make clean` - Remove build artifacts

**Note:** The server build **always** builds and embeds the client automatically.

## How to Play

1. Open two browser windows at `http://localhost:5173`
2. Click **"Join Queue"** in both windows
3. Players will be matched automatically
4. You'll start with $100 and earn $10/second
5. Purchase units:
   - **Tank** ($50): Slower but can win the game
   - **Airplane** ($80): Faster but cannot win
6. Units automatically move toward the enemy base
7. Units will attack enemies in range
8. **Win**: Get your tank to the enemy base!

## Game Mechanics

### Economy
- Starting money: $100
- Passive income: $10 per second
- Tank cost: $50
- Airplane cost: $80

### Unit Stats

| Unit | Cost | HP | Speed | Damage | Attack Range | Hits to Destroy |
|------|------|-----|-------|--------|--------------|-----------------|
| Tank | $50 | 30 | 5 | 10 | 10 | 3 |
| Airplane | $80 | 30 | 15 | 10 | 20 | 3 |

### Combat
- Units **autonomously fire** at enemies when within attack range
- Each unit deals 10 damage per hit
- Units are **destroyed after 3 hits** (30 HP / 10 damage = 3 hits)
- Units do NOT collide with each other - they pass through
- Destroyed units are removed from the battlefield
- Attack speed: 1 attack per second for both unit types

### Win Condition
- **Only tanks** can win by reaching the enemy base
- Airplanes serve as support units for offense/defense
- The game ends immediately when a tank reaches a base

## Development

### Building for Production

**Client:**
```bash
cd client
bun run build
```
Output will be in `client/dist/`

**Server:**
```bash
cd server
go build -o arena-server main.go
```

### Architecture Highlights

**Server (Authoritative)**:
- Game state is the single source of truth
- Runs at 20 ticks per second (TPS)
- Each game room runs in its own goroutine
- Systems: Movement, Combat, Economy, Win Condition

**Client (Rendering Only)**:
- Receives state updates from server
- Renders at 60 FPS for smooth visuals
- Sends purchase commands to server
- All game logic validated server-side

## Deployment

The server is built as a single standalone binary with the client embedded:

```bash
# Build the binary
make build

# Deploy the binary to your server
scp server/bin/arena-server user@your-server:/opt/arena-game/

# Run on the server
./arena-server
```

**Deployment Platforms:**
- **Fly.io**: Deploy the binary directly
- **Railway**: Use Dockerfile with the binary
- **Render**: Deploy as a native service
- **VPS/Cloud**: Run the binary with systemd or supervisor

**Example systemd service:**
```ini
[Unit]
Description=Arena Game Server
After=network.target

[Service]
Type=simple
User=arena
WorkingDirectory=/opt/arena-game
ExecStart=/opt/arena-game/arena-server
Restart=always

[Install]
WantedBy=multi-user.target
```

**Deployment Benefits:**
- ✅ Single file deployment (~12MB binary)
- ✅ No separate frontend build/deployment needed
- ✅ Automatic WebSocket configuration
- ✅ No CORS complexity
- ✅ Works offline/air-gapped environments

## Future Enhancements

- Multiple unit types (infantry, artillery, helicopters)
- Base upgrades and defenses
- Power-ups and special abilities
- Ranked matchmaking with ELO
- Replay system
- Mobile support
- Better 3D models (GLTF)
- Fog of war
- Map obstacles and terrain

## License

This project is open source and available for educational purposes.
