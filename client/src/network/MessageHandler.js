// Handles incoming WebSocket messages

export class MessageHandler {
  constructor(gameState) {
    this.gameState = gameState;
    this.handlers = new Map();
  }

  handle(message) {
    const { type, payload } = message;

    if (this.handlers.has(type)) {
      this.handlers.get(type)(payload);
    } else {
      console.warn('Unhandled message type:', type);
    }
  }

  on(type, handler) {
    this.handlers.set(type, handler);
  }

  setupDefaultHandlers() {
    this.on('game_start', (payload) => {
      console.log('Game started:', payload);
      this.gameState.setPlayerInfo(payload.playerId, payload.gameId);
      if (payload.state) {
        this.gameState.update(payload.state);
      }
    });

    this.on('game_update', (payload) => {
      if (payload.state) {
        this.gameState.update(payload.state);
      }
    });

    this.on('game_over', (payload) => {
      console.log('Game over:', payload);
      this.gameState.winner = payload.winner;
      this.gameState.gameStatus = 'finished';
    });

    this.on('error', (payload) => {
      console.error('Server error:', payload.message);
      // Default handler just logs - Game.js overrides this for popup display
    });

    this.on('info', (payload) => {
      console.log('Server info:', payload);
    });
  }
}
