// Client-side game state management

export class GameState {
  constructor() {
    this.timestamp = 0;
    this.players = [null, null];
    this.units = [];
    this.obstacles = [];
    this.projectiles = [];
    this.buyZones = [];
    this.turrets = [];
    this.pendingSpawns = [];
    this.gameStatus = 'waiting'; // 'waiting', 'playing', 'finished'
    this.winner = null;
    this.playerId = null;
    this.gameId = null;
  }

  update(newState) {
    this.timestamp = newState.timestamp || Date.now();
    this.players = newState.players || this.players;
    this.units = newState.units || [];
    this.obstacles = newState.obstacles || this.obstacles;
    this.projectiles = newState.projectiles || [];
    this.buyZones = newState.buyZones || this.buyZones;
    this.turrets = newState.turrets || this.turrets;
    this.pendingSpawns = newState.pendingSpawns || [];
    this.gameStatus = newState.gameStatus || this.gameStatus;
    this.winner = newState.winner !== undefined ? newState.winner : this.winner;
  }

  setPlayerInfo(playerId, gameId) {
    this.playerId = playerId;
    this.gameId = gameId;
  }

  getMyPlayer() {
    return this.playerId !== null ? this.players[this.playerId] : null;
  }

  getOpponentPlayer() {
    if (this.playerId === null) return null;
    return this.players[1 - this.playerId];
  }

  getMyUnits() {
    return this.units.filter(unit => unit.ownerId === this.playerId);
  }

  getOpponentUnits() {
    if (this.playerId === null) return [];
    return this.units.filter(unit => unit.ownerId !== this.playerId);
  }

  getMyPlayerUnit() {
    if (this.playerId === null) return null;
    return this.units.find(unit => unit.type === 'player' && unit.ownerId === this.playerId) || null;
  }

  clear() {
    this.timestamp = 0;
    this.players = [null, null];
    this.units = [];
    this.obstacles = [];
    this.projectiles = [];
    this.buyZones = [];
    this.turrets = [];
    this.pendingSpawns = [];
    this.gameStatus = 'waiting';
    this.winner = null;
    this.playerId = null;
    this.gameId = null;
  }

  // Get pending spawns for a specific player
  getMyPendingSpawns() {
    if (this.playerId === null) return [];
    return this.pendingSpawns.filter(spawn => spawn.ownerId === this.playerId);
  }

  // Get count of pending spawns by type for the current player
  getMyPendingSpawnCounts() {
    const counts = { tank: 0, airplane: 0 };
    for (const spawn of this.getMyPendingSpawns()) {
      if (counts[spawn.unitType] !== undefined) {
        counts[spawn.unitType]++;
      }
    }
    return counts;
  }
}
