// Game constants (should match server)
export const ARENA_SIZE = 200;
export const BASE_RADIUS = 5;
export const BASE_SIZE = 25; // Size of base area (matches server)

export const STARTING_MONEY = 100;
export const PASSIVE_INCOME_PER_SECOND = 10;

// Unit constants
export const UNITS = {
  tank: {
    cost: 50,
    speed: 5.0,
    health: 30,  // 3 hits to destroy
    damage: 10,
    attackRange: 10.0,
    yPosition: 1.0
  },
  airplane: {
    cost: 80,
    speed: 15.0,
    health: 30,  // 3 hits to destroy
    damage: 10,
    attackRange: 20.0,
    yPosition: 10.0
  }
};

// Base positions
export const BASE_POSITIONS = {
  0: { x: -90, y: 0, z: 0 },
  1: { x: 90, y: 0, z: 0 }
};

// Player colors
export const PLAYER_COLORS = ['#3b82f6', '#ef4444']; // Blue and Red

// WebSocket server URL
// In production (embedded), use same host. In development, use explicit localhost:3000
const isProduction = import.meta.env.PROD;
const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const wsHost = isProduction ? window.location.host : 'localhost:3000';

export const WS_URL = `${wsProtocol}//${wsHost}/ws`;
