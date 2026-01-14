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
import { Leaderboard } from '../ui/Leaderboard.js';
import { AuthService } from '../auth/AuthService.js';

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
    this.leaderboard = null;
    this.authService = null;
  }

  init() {
    // Get container
    const container = document.getElementById('game-container');

    // Initialize auth service
    this.authService = new AuthService();
    this.setupAuthUI();
    this.authService.handleAuthCallback();
    this.authService.checkAuth();

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
    this.leaderboard = new Leaderboard();

    // Fetch leaderboard on startup
    this.leaderboard.fetch();

    // Give gameLoop access to the popup for position updates
    this.gameLoop.setBuyZonePopup(this.buyZonePopup);

    // Initialize player input
    this.playerInput = new PlayerInput(
      this.camera,
      this.renderer.getRenderer(),
      (direction) => this.sendPlayerMove(direction),
      (targetX, targetZ) => this.sendPlayerShoot(targetX, targetZ),
      (zoneId) => this.sendBuyFromZone(zoneId),
      (turretId) => this.sendClaimTurret(turretId),
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

  sendClaimTurret(turretId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('claim_turret', { turretId });
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
    // Clear game state
    this.gameState.clear();

    // Reset the game loop (clears all meshes and resets flags)
    this.gameLoop.reset();

    // Hide game over screen
    this.gameOverScreen.hide();

    // Disable player input until new game starts
    this.playerInput.disable();

    // Show queue screen
    document.getElementById('queue-screen').classList.remove('hidden');
    document.getElementById('join-queue-button').disabled = false;
    document.getElementById('queue-status').classList.add('hidden');

    // Refresh leaderboard
    this.leaderboard.fetch();
  }

  setupAuthUI() {
    const githubLoginButton = document.getElementById('github-login-button');
    const blueskyLoginButton = document.getElementById('bluesky-login-button');
    const logoutButton = document.getElementById('logout-button');
    const loggedOutSection = document.getElementById('auth-logged-out');
    const loggedInSection = document.getElementById('auth-logged-in');
    const userAvatar = document.getElementById('user-avatar');
    const userName = document.getElementById('user-name');

    // BlueSky modal elements
    const blueskyModal = document.getElementById('bluesky-modal');
    const blueskyForm = document.getElementById('bluesky-login-form');
    const blueskyHandle = document.getElementById('bluesky-handle');
    const blueskyAppPassword = document.getElementById('bluesky-app-password');
    const blueskyError = document.getElementById('bluesky-error');
    const blueskyCancel = document.getElementById('bluesky-cancel');
    const blueskySubmit = document.getElementById('bluesky-submit');
    const blueskyBackdrop = blueskyModal.querySelector('.modal-backdrop');

    // GitHub login button click
    githubLoginButton.addEventListener('click', () => {
      this.authService.loginWithGitHub();
    });

    // BlueSky login button click - show modal
    blueskyLoginButton.addEventListener('click', () => {
      blueskyModal.classList.remove('hidden');
      blueskyHandle.focus();
    });

    // Close modal on cancel or backdrop click
    const closeBlueskyModal = () => {
      blueskyModal.classList.add('hidden');
      blueskyForm.reset();
      blueskyError.classList.add('hidden');
    };

    blueskyCancel.addEventListener('click', closeBlueskyModal);
    blueskyBackdrop.addEventListener('click', closeBlueskyModal);

    // Handle BlueSky form submit
    blueskyForm.addEventListener('submit', async (e) => {
      e.preventDefault();

      const handle = blueskyHandle.value.trim();
      const appPassword = blueskyAppPassword.value.trim();

      if (!handle || !appPassword) {
        blueskyError.textContent = 'Please enter both handle and app password';
        blueskyError.classList.remove('hidden');
        return;
      }

      // Disable submit button while logging in
      blueskySubmit.disabled = true;
      blueskySubmit.textContent = 'Logging in...';
      blueskyError.classList.add('hidden');

      const result = await this.authService.loginWithBlueSky(handle, appPassword);

      if (result.success) {
        closeBlueskyModal();
        // Reconnect WebSocket to pick up new auth identity
        if (this.ws) {
          this.ws.reconnect();
        }
      } else {
        blueskyError.textContent = result.error;
        blueskyError.classList.remove('hidden');
      }

      blueskySubmit.disabled = false;
      blueskySubmit.textContent = 'Login';
    });

    // Logout button click
    logoutButton.addEventListener('click', async () => {
      await this.authService.logout();
    });

    // Listen for auth state changes
    this.authService.onAuthChange((user) => {
      if (user && !user.isGuest) {
        // User is logged in
        loggedOutSection.classList.add('hidden');
        loggedInSection.classList.remove('hidden');
        userName.textContent = user.displayName;
        if (user.avatarUrl) {
          userAvatar.src = user.avatarUrl;
        } else {
          // Default avatar for users without one
          userAvatar.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(user.displayName)}&background=random`;
        }
      } else {
        // User is logged out or guest
        loggedOutSection.classList.remove('hidden');
        loggedInSection.classList.add('hidden');
      }
    });
  }
}
