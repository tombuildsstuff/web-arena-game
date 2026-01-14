import { Scene } from './Scene.js';
import { Renderer } from './Renderer.js';
import { Camera } from './Camera.js';
import { GameLoop } from './GameLoop.js';
import { GameState } from '../state/GameState.js';
import { WebSocketClient } from '../network/WebSocketClient.js';
import { MessageHandler } from '../network/MessageHandler.js';
import { HUD } from '../ui/HUD.js';
import { GameOverScreen } from '../ui/GameOverScreen.js';
import { PlayerInput } from '../input/PlayerInput.js';
import { BuyZonePopup } from '../ui/BuyZonePopup.js';

export class Game {
  constructor() {
    this.gameState = new GameState();
    this.scene = null;
    this.renderer = null;
    this.camera = null;
    this.gameLoop = null;
    this.ws = null;
    this.messageHandler = null;
    this.hud = null;
    this.gameOverScreen = null;
    this.playerInput = null;
    this.buyZonePopup = null;
  }

  init() {
    // Get container
    const container = document.getElementById('game-container');

    // Initialize renderer
    this.renderer = new Renderer(container);

    // Initialize camera
    this.camera = new Camera(this.renderer.getRenderer());

    // Initialize scene
    this.scene = new Scene();

    // Initialize game loop
    this.gameLoop = new GameLoop(
      this.renderer,
      this.scene,
      this.camera,
      this.gameState
    );

    // Initialize UI
    this.hud = new HUD(this.gameState, this.gameLoop);
    this.gameOverScreen = new GameOverScreen(this.gameState, () => {
      this.resetGame();
    });
    this.buyZonePopup = new BuyZonePopup();
    this.buyZonePopup.setCamera(this.camera, this.renderer.getRenderer());

    // Give gameLoop access to the popup for position updates
    this.gameLoop.setBuyZonePopup(this.buyZonePopup);

    // Initialize player input
    this.playerInput = new PlayerInput(
      this.camera,
      this.renderer.getRenderer(),
      (direction) => this.sendPlayerMove(direction),
      (targetX, targetZ) => this.sendPlayerShoot(targetX, targetZ),
      (zoneId) => this.sendBuyFromZone(zoneId),
      this.gameLoop
    );

    // Initialize WebSocket
    this.setupWebSocket();

    // Start game loop
    this.gameLoop.start();

    // Setup queue screen
    this.setupQueueScreen();
  }

  setupWebSocket() {
    this.messageHandler = new MessageHandler(this.gameState);
    this.messageHandler.setupDefaultHandlers();

    // Override error handler to use popup for buy zone errors
    this.messageHandler.on('error', (payload) => {
      console.error('Server error:', payload.message);
      this.showBuyZoneError(payload.message);
    });

    // Add custom handler for game start
    this.messageHandler.on('game_start', (payload) => {
      console.log('Game started:', payload);
      this.gameState.setPlayerInfo(payload.playerId, payload.gameId);

      if (payload.state) {
        this.gameState.update(payload.state);
        this.scene.createBases(this.gameState.players);
      }

      // Hide queue screen
      document.getElementById('queue-screen').classList.add('hidden');

      // Update UI
      this.hud.update();

      // Enable player input
      this.playerInput.enable();
    });

    // Add custom handler for game updates
    this.messageHandler.on('game_update', (payload) => {
      if (payload.state) {
        this.gameState.update(payload.state);
        this.hud.update();
      }
    });

    // Add custom handler for game over
    this.messageHandler.on('game_over', (payload) => {
      console.log('Game over:', payload);
      this.gameState.winner = payload.winner;
      this.gameState.gameStatus = 'finished';
      this.playerInput.disable();
      this.gameOverScreen.show(payload.winner, payload.reason, payload.matchDuration, payload.stats);
    });

    this.ws = new WebSocketClient((message) => {
      this.messageHandler.handle(message);
    });

    this.ws.setConnectionChangeHandler((connected) => {
      this.updateConnectionStatus(connected);
    });

    this.ws.connect();
  }

  setupQueueScreen() {
    const joinButton = document.getElementById('join-queue-button');
    const queueStatus = document.getElementById('queue-status');

    joinButton.addEventListener('click', () => {
      if (this.ws.isConnected()) {
        this.ws.send('join_queue', {});
        joinButton.disabled = true;
        queueStatus.classList.remove('hidden');
      }
    });
  }

  updateConnectionStatus(connected) {
    const statusElement = document.getElementById('connection-status');
    const statusText = document.getElementById('status-text');
    const joinButton = document.getElementById('join-queue-button');

    if (connected) {
      statusElement.className = 'connected';
      statusText.textContent = 'Connected';
      joinButton.disabled = false;
    } else {
      statusElement.className = 'disconnected';
      statusText.textContent = 'Disconnected';
      joinButton.disabled = true;
    }
  }

  sendPlayerMove(direction) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('player_move', { direction });
    }
  }

  sendPlayerShoot(targetX, targetZ) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('player_shoot', { targetX, targetZ });
    }
  }

  sendBuyFromZone(zoneId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('buy_from_zone', { zoneId });
    }
  }

  showBuyZoneError(message) {
    // Get the nearby buy zone position for the popup
    const nearbyZone = this.gameLoop.getNearbyBuyZone();
    if (nearbyZone && this.buyZonePopup) {
      this.buyZonePopup.show(message, nearbyZone.position);
    }
  }

  resetGame() {
    this.gameState.clear();
    this.gameOverScreen.hide();
    this.playerInput.disable();
    document.getElementById('queue-screen').classList.remove('hidden');
    document.getElementById('join-queue-button').disabled = false;
    document.getElementById('queue-status').classList.add('hidden');
  }
}
