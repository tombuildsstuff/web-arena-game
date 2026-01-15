# Arena Game

A real-time 3D multiplayer strategy game where two players compete by purchasing and deploying tanks, helicopters, and turrets while controlling a player character that can shoot enemies directly.

**[Play Now at arena.ibuildstuff.eu](https://arena.ibuildstuff.eu)**

## Features

- **Real-time Multiplayer**: Play against another player or AI with live updates via WebSocket
- **3D Graphics**: Built with Three.js for immersive 3D visualization
- **Player Character**: Control a soldier that can move and shoot enemies
- **Strategic Gameplay**:
  - Purchase tanks, helicopters, and super units
  - Claim neutral turrets and forward bases
  - Earn passive income and kill rewards
  - Units automatically move toward enemy base
  - Units engage in combat when in range
  - Collect health packs to heal your player
- **Point-based Scoring**: Earn points for kills (tanks: 10, helicopters: 20, turrets: 20, players: 50)
- **AI Opponents**: Practice against Easy, Medium, or Hard AI
- **Spectator Mode**: Watch live games in progress
- **Leaderboard**: Track top players by points and wins
- **Authentication**: Login with GitHub, BlueSky, or play as guest
- **Sound Effects**: Immersive audio for shooting, explosions, and more
- **Mobile Support**: Touch controls with dual joysticks

## Controls

| Key | Action |
|-----|--------|
| WASD / Arrow Keys | Move player |
| Left Click / X | Shoot |
| C | Buy unit / Claim turret or base |
| V | Bulk buy (10 units at 10% discount) |
| Right-drag | Rotate camera |
| Scroll | Zoom in/out |
| H | Toggle help overlay |

## Tech Stack

### Client
- **Vanilla JavaScript** with ES6+ modules
- **Three.js 0.160** for 3D rendering
- **Vite 5.4** for build tooling and dev server
- **Bun 1.x** as the JavaScript runtime

### Server
- **Go 1.25** for high-performance game logic
- **gorilla/websocket 1.5** for WebSocket connections
- **chi 5.2** router for HTTP handling
- Authoritative server architecture (20 TPS)

## Project Structure

```
web-arena-game/
├── client/                  # Frontend (Vanilla JS + Three.js)
│   ├── src/
│   │   ├── game/           # 3D game rendering (units, arena, effects)
│   │   ├── input/          # Player input handling
│   │   ├── network/        # WebSocket client
│   │   ├── ui/             # DOM-based UI (HUD, menus)
│   │   ├── state/          # Client state management
│   │   ├── audio/          # Sound effects manager
│   │   └── utils/          # Utilities and constants
│   ├── public/sounds/      # Audio files
│   ├── index.html
│   └── package.json
│
└── server/                  # Backend (Go)
    ├── main.go              # Server entry point
    ├── internal/
    │   ├── game/           # Game logic (room, AI, spawning, turrets)
    │   ├── websocket/      # WebSocket hub and client handling
    │   ├── types/          # Shared types, events, constants
    │   └── auth/           # OAuth (GitHub, BlueSky)
    ├── leaderboard.json    # Persistent leaderboard data
    └── go.mod
```

## Getting Started

### Prerequisites

- **Go 1.25+**: [Download Go](https://go.dev/dl/)
- **Bun 1.x**: [Install Bun](https://bun.sh/)
- **Make**: Usually pre-installed on macOS/Linux

### Quick Start

```bash
# Install all dependencies
make tools

# Build server (client auto-embedded)
make build

# Run the server
make run
```

Then open browser at `http://localhost:3000`

**For development with hot reload:**

```bash
make run-dev
```

Then open browser at `http://localhost:5173`

### Makefile Targets

**Root-level** (`make <target>` from project root):
- `make tools` - Install all tools and dependencies
- `make build` - Build server (client auto-embedded)
- `make run` - Build and run server
- `make run-dev` - Run dev servers with hot reload
- `make clean` - Clean all build artifacts
- `make test` - Run server tests

## How to Play

1. Open your browser at `http://localhost:3000` (or 5173 for dev)
2. Optionally login with GitHub or BlueSky for leaderboard tracking
3. Click **"Join Queue"** to find an opponent, or **"Play vs AI"**
4. Control your player with WASD, shoot with click or X
5. Walk to buy zones and press C to purchase units
6. Claim neutral turrets and forward bases for tactical advantage
7. Use V for bulk purchases (10 units at 10% discount)
8. Collect health packs (green crosses) to heal
9. Earn points by destroying enemy units
10. **Win**: Have more points when the 5-minute timer ends!

## Game Mechanics

### Economy
| Item | Cost | Notes |
|------|------|-------|
| Starting money | $200 | |
| Passive income | $10/second | |
| Tank | $50 | From base or forward bases |
| Helicopter | $80 | From base or forward bases |
| Super Tank | $250 | One per player, from base only |
| Super Helicopter | $400 | One per player, from base only |
| Bulk buy (10 units) | 10% discount | $450 for tanks, $720 for helicopters |
| Forward Base claim | $500 | Top/bottom claimable zones |
| Turret claim | Free | Walk to neutral turret and press C |

### Kill Rewards
| Target | Reward |
|--------|--------|
| Tank | $15 |
| Helicopter | $25 |
| Super Tank | $75 |
| Super Helicopter | $100 |
| Player | $50 |
| Turret | $30 |

### Unit Stats

| Unit | HP | Speed | Damage | Attack Rate | Range |
|------|-----|-------|--------|-------------|-------|
| Player | 80 | 12 | 15 | 3/sec | 25 |
| Tank | 30 | 5 | 10 | 1/sec | 10 |
| Helicopter | 30 | 15 | 10 | 1/sec | 20 |
| Super Tank | 150 | 4 | 25 | 0.5/sec | 15 |
| Super Helicopter | 100 | 12 | 20 | 1.5/sec | 25 |
| Turret | 100 | - | 15 | 2/sec | 30 |
| Base Turret | 150 | - | 20 | 2/sec | 35 |

### Scoring System
| Kill | Points |
|------|--------|
| Tank | 10 |
| Helicopter | 20 |
| Super Tank | 30 |
| Super Helicopter | 40 |
| Turret | 20 |
| Player | 50 |

### Win Condition
- Games last 5 minutes
- Player with the most points wins
- Points are earned by destroying enemy units
- If tied, the game ends in a draw

### Map Features
- **Player Bases**: Protected spawn areas at each end
- **Base Turrets**: 4 turrets at each base corner, defend the base
- **Buy Zones**: Areas to purchase tanks and helicopters
- **Forward Bases**: Claimable zones (top/bottom) for forward spawning
- **Neutral Turrets**: Claimable defensive positions in the middle
- **Health Packs**: Spawn periodically, heal 30 HP

## AI Difficulty

| Difficulty | Buy Interval | Turret Claim | Zone Claim | Spawn Delay |
|------------|--------------|--------------|------------|-------------|
| Easy | 8s | 20% chance | 10% chance | +3s |
| Medium | 5s | 50% chance | 30% chance | Normal |
| Hard | 3s | 80% chance | 60% chance | -1s |

## Development

### Architecture

**Server (Authoritative)**:
- Game state is the single source of truth
- Runs at 20 ticks per second (TPS)
- Each game room runs in its own goroutine
- Systems: Movement, Combat, Economy, Spawning, AI

**Client (Rendering)**:
- Receives state updates from server
- Renders at 60 FPS for smooth visuals
- Handles input and sends commands to server
- All game logic validated server-side

### Building for Production

```bash
# Build single binary with embedded client
make build

# Binary located at server/bin/arena-server
```

## Deployment

The server is built as a single standalone binary with the client embedded:

```bash
# Build and deploy
make build
scp server/bin/arena-server user@your-server:/opt/arena-game/

# Run on the server
./arena-server
```

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

## Built With AI

This project was built as a collaboration between a human and an AI. Here's how the work was divided:

### Claude Code (AI)

- **All source code**: Go server, JavaScript client, HTML, CSS
- **Game design**: Mechanics, balancing, unit stats, economy system
- **3D rendering**: Three.js scene setup, models, effects, animations
- **Networking**: WebSocket protocol, state synchronization, client prediction
- **AI opponents**: Behavior trees, difficulty scaling
- **Authentication**: OAuth flows for GitHub and BlueSky
- **Documentation**: README, code comments, CLAUDE.md
- **Sound research**: Finding and suggesting sound effects from Freesound.org

### Human ([@tombuildsstuff](https://github.com/tombuildsstuff))

- **Prompting & direction**: Describing features, requesting changes, guiding development
- **GitHub OAuth app**: Creating and configuring the GitHub application for SSO
- **Sound selection**: Listening to suggested sounds and picking the final choices
- **Deployment**: Configuring the Railway hosting platform
- **DNS**: Setting up the arena.ibuildstuff.eu domain

## Credits

- **Code**: [Claude Code](https://www.anthropic.com/claude-code) by Anthropic
- **Direction**: [@tombuildsstuff](https://github.com/tombuildsstuff)
- **Sound Effects**: Various artists from [Freesound.org](https://freesound.org) (see in-game credits)

## Inspiration

This game is inspired by the "Precinct Assault" mode from [**Future Cop: LAPD**](https://en.wikipedia.org/wiki/Future_Cop:_LAPD) (1998), developed by EA Redwood Shores and published by Electronic Arts. The original game featured a similar two-base strategy mode where players captured turrets and spawned units while controlling a player character directly.

All rights to Future Cop: LAPD are reserved by Electronic Arts.

## License

This project is open source and available for educational purposes.
