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
import { TouchControls } from '../input/TouchControls.js';
import { BuyZonePopup } from '../ui/BuyZonePopup.js';
import { Leaderboard } from '../ui/Leaderboard.js';
import { AuthService } from '../auth/AuthService.js';
import { SoundManager } from '../audio/SoundManager.js';

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
    this.touchControls = null;
    this.buyZonePopup = null;
    this.leaderboard = null;
    this.authService = null;
    this.soundManager = null;
    this.isSpectating = false;
    this.lobbyStatusInterval = null;
  }

  async init() {
    // Get container
    const container = document.getElementById('game-container');

    // Initialize auth service
    this.authService = new AuthService();
    this.setupAuthUI();
    this.authService.handleAuthCallback();

    // Wait for auth check to complete (ensures guest cookie is set before WebSocket connects)
    await this.authService.checkAuth();

    // Initialize renderer
    this.renderer = new Renderer(container);

    // Initialize camera
    this.camera = new Camera(this.renderer.getRenderer());

    // Initialize scene
    this.scene = new Scene();

    // Initialize sound manager
    this.soundManager = new SoundManager();
    this.soundManager.init();

    // Initialize game loop
    this.gameLoop = new GameLoop(
      this.renderer,
      this.scene,
      this.camera,
      this.gameState,
      this.soundManager
    );

    // Initialize UI
    this.hud = new HUD(this.gameState, this.gameLoop, this.soundManager);
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
      (zoneId) => this.sendBulkBuyFromZone(zoneId),
      (turretId) => this.sendClaimTurret(turretId),
      (zoneId) => this.sendClaimBuyZone(zoneId),
      (barracksId) => this.sendClaimBarracks(barracksId),
      this.gameLoop
    );

    // Update movement direction when camera rotates (for camera-relative controls)
    this.camera.setRotationChangeCallback(() => {
      this.playerInput.updateMovementDirection();
    });

    // Initialize touch controls for mobile devices (user agent detection only)
    if (this.isMobileDevice()) {
      this.touchControls = new TouchControls({
        camera: this.camera,
        gameState: this.gameState,
        canvas: this.renderer.getRenderer().domElement,
        onMove: (direction) => this.sendPlayerMove(direction),
        onShoot: (targetX, targetZ) => this.sendPlayerShoot(targetX, targetZ),
        onInteract: () => this.playerInput.handleInteraction()
      });
    }

    // Initialize WebSocket
    this.setupWebSocket();

    // Start game loop
    this.gameLoop.start();

    // Setup queue screen
    this.setupQueueScreen();

    // Setup info modals (About/Credits)
    this.setupInfoModals();
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

      // Store map definition if provided
      if (payload.map) {
        this.gameState.setMapDefinition(payload.map);
        console.log('Map loaded:', payload.map.name);
      }

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
      if (this.touchControls) {
        this.touchControls.enable();
      }

      // Play match start sound
      if (this.soundManager) {
        this.soundManager.resume();
        this.soundManager.playGlobal('match_start');
      }
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
      if (this.touchControls) {
        this.touchControls.disable();
      }

      // Play match end sound
      if (this.soundManager) {
        this.soundManager.playGlobal('match_end');
      }

      if (this.isSpectating) {
        // Spectator - just show game over and return to lobby
        this.stopSpectating();
      } else {
        // Player - show full game over screen
        this.gameOverScreen.show(payload.winner, payload.reason, payload.matchDuration, payload.stats);
      }
    });

    // Add handler for lobby status
    this.messageHandler.on('lobby_status', (payload) => {
      this.updateLobbyStatus(payload);
    });

    // Add handler for spectate start
    this.messageHandler.on('spectate_start', (payload) => {
      console.log('Spectating game:', payload.gameId);
      this.isSpectating = true;

      // Enable spectator mode on game loop
      this.gameLoop.setSpectatorMode(true);

      if (payload.state) {
        this.gameState.update(payload.state);
        this.scene.createBases(this.gameState.players);
      }

      // Hide queue screen, show spectator HUD
      document.getElementById('queue-screen').classList.add('hidden');
      document.getElementById('hud').classList.add('hidden');
      document.getElementById('spectator-hud').classList.remove('hidden');

      // Update spectator players display
      const players = this.gameState.players;
      if (players && players.length >= 2) {
        document.getElementById('spectator-players').textContent =
          `${players[0].displayName} vs ${players[1].displayName}`;
      }
    });

    // Add handler for spectate stopped
    this.messageHandler.on('spectate_stopped', () => {
      this.returnToLobby();
    });

    this.ws = new WebSocketClient((message) => {
      this.messageHandler.handle(message);
    });

    this.ws.setConnectionChangeHandler((connected) => {
      this.updateConnectionStatus(connected);

      if (connected) {
        // Request lobby status periodically when connected
        this.requestLobbyStatus();
        this.lobbyStatusInterval = setInterval(() => {
          if (this.ws.isConnected() && !this.isSpectating && this.gameState.gameStatus !== 'playing') {
            this.requestLobbyStatus();
          }
        }, 3000); // Update every 3 seconds
      } else {
        // Clear interval when disconnected
        if (this.lobbyStatusInterval) {
          clearInterval(this.lobbyStatusInterval);
          this.lobbyStatusInterval = null;
        }
      }
    });

    this.ws.connect();
  }

  requestLobbyStatus() {
    if (this.ws.isConnected()) {
      this.ws.send('get_lobby_status', {});
    }
  }

  updateLobbyStatus(payload) {
    const { queueSize, activeGames } = payload;

    // Update players waiting
    const playersWaitingEl = document.getElementById('players-waiting');
    if (queueSize > 0) {
      playersWaitingEl.textContent = queueSize === 1
        ? '1 player waiting'
        : `${queueSize} players waiting`;
      playersWaitingEl.classList.remove('hidden');
    } else {
      playersWaitingEl.classList.add('hidden');
    }

    // Update active games
    const activeGamesListEl = document.getElementById('active-games-list');
    const noActiveGamesEl = document.getElementById('no-active-games');

    if (activeGames && activeGames.length > 0) {
      noActiveGamesEl.classList.add('hidden');
      activeGamesListEl.innerHTML = activeGames.map(game => `
        <div class="active-game-item">
          <span class="game-players">
            ${this.escapeHtml(game.player1Name)}
            <span class="vs">vs</span>
            ${this.escapeHtml(game.player2Name)}
          </span>
          <div>
            <span class="game-spectators">${game.spectatorCount} watching</span>
            <button class="spectate-button" data-game-id="${game.gameId}">Watch</button>
          </div>
        </div>
      `).join('');

      // Add click handlers for spectate buttons
      activeGamesListEl.querySelectorAll('.spectate-button').forEach(button => {
        button.addEventListener('click', () => {
          const gameId = button.dataset.gameId;
          this.spectateGame(gameId);
        });
      });
    } else {
      activeGamesListEl.innerHTML = '';
      noActiveGamesEl.classList.remove('hidden');
    }
  }

  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  spectateGame(gameId) {
    if (this.ws.isConnected()) {
      this.ws.send('spectate_game', { gameId });
    }
  }

  stopSpectating() {
    if (this.ws.isConnected()) {
      this.ws.send('stop_spectating', {});
    }
    this.returnToLobby();
  }

  returnToLobby() {
    this.isSpectating = false;
    this.gameLoop.setSpectatorMode(false);
    this.gameState.clear();
    this.gameLoop.reset();

    // Hide spectator HUD, show queue screen
    document.getElementById('spectator-hud').classList.add('hidden');
    document.getElementById('hud').classList.remove('hidden');
    document.getElementById('queue-screen').classList.remove('hidden');
    document.getElementById('join-queue-button').disabled = false;
    document.getElementById('play-vs-ai-button').disabled = false;
    document.getElementById('queue-status').classList.add('hidden');

    // Refresh lobby status and leaderboard
    this.requestLobbyStatus();
    this.leaderboard.fetch();
  }

  setupQueueScreen() {
    const joinButton = document.getElementById('join-queue-button');
    const queueStatus = document.getElementById('queue-status');
    const mpMapSelect = document.getElementById('mp-map-select');

    joinButton.addEventListener('click', () => {
      if (this.ws.isConnected()) {
        const mapId = mpMapSelect ? mpMapSelect.value : 'classic';
        this.ws.send('join_queue', { mapId });
        joinButton.disabled = true;
        document.getElementById('play-vs-ai-button').disabled = true;
        queueStatus.classList.remove('hidden');
      } else {
        // Try to reconnect if not connected
        this.ws.forceReconnect();
      }
    });

    // Play vs AI button
    const playVsAIButton = document.getElementById('play-vs-ai-button');
    const aiDifficulty = document.getElementById('ai-difficulty');
    const aiMapSelect = document.getElementById('ai-map-select');
    if (playVsAIButton) {
      playVsAIButton.addEventListener('click', () => {
        if (this.ws.isConnected()) {
          const difficulty = aiDifficulty.value;
          const mapId = aiMapSelect ? aiMapSelect.value : 'classic';
          this.ws.send('start_vs_ai', { difficulty, mapId });
          joinButton.disabled = true;
          playVsAIButton.disabled = true;
        } else {
          // Try to reconnect if not connected
          this.ws.forceReconnect();
        }
      });
    }

    // Stop spectating button
    const stopSpectatingButton = document.getElementById('stop-spectating-button');
    if (stopSpectatingButton) {
      stopSpectatingButton.addEventListener('click', () => {
        this.stopSpectating();
      });
    }
  }

  updateConnectionStatus(connected) {
    const statusElement = document.getElementById('connection-status');
    const statusText = document.getElementById('status-text');
    const joinButton = document.getElementById('join-queue-button');
    const playVsAIButton = document.getElementById('play-vs-ai-button');

    if (connected) {
      statusElement.className = 'connected';
      statusText.textContent = 'Connected';
      joinButton.disabled = false;
      if (playVsAIButton) playVsAIButton.disabled = false;
    } else {
      statusElement.className = 'disconnected';
      statusText.textContent = 'Disconnected';
      joinButton.disabled = true;
      if (playVsAIButton) playVsAIButton.disabled = true;
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

  sendBulkBuyFromZone(zoneId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('bulk_buy_from_zone', { zoneId });
    }
  }

  sendClaimTurret(turretId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('claim_turret', { turretId });
    }
  }

  sendClaimBuyZone(zoneId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('claim_buy_zone', { zoneId });
    }
  }

  sendClaimBarracks(barracksId) {
    if (this.ws.isConnected() && this.gameState.gameStatus === 'playing') {
      this.ws.send('claim_barracks', { barracksId });
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
    if (this.touchControls) {
      this.touchControls.disable();
    }

    // Show queue screen
    document.getElementById('queue-screen').classList.remove('hidden');
    document.getElementById('join-queue-button').disabled = false;
    document.getElementById('play-vs-ai-button').disabled = false;
    document.getElementById('queue-status').classList.add('hidden');

    // Refresh leaderboard
    this.leaderboard.fetch();
  }

  isMobileDevice() {
    // User agent detection only - don't use touch capability checks
    // as those return true on many desktop devices
    const ua = navigator.userAgent.toLowerCase();
    return /android|iphone|ipad|ipod|mobile|tablet/.test(ua);
  }

  setupAuthUI() {
    const githubLoginButton = document.getElementById('github-login-button');
    const blueskyLoginButton = document.getElementById('bluesky-login-button');
    const logoutButton = document.getElementById('logout-button');
    const guestSection = document.getElementById('auth-guest');
    const loggedInSection = document.getElementById('auth-logged-in');
    const guestName = document.getElementById('guest-name');
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
        // User is logged in with GitHub/BlueSky
        guestSection.classList.add('hidden');
        loggedInSection.classList.remove('hidden');
        userName.textContent = user.displayName;
        if (user.avatarUrl) {
          userAvatar.src = user.avatarUrl;
        } else {
          // Default avatar for users without one
          userAvatar.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(user.displayName)}&background=random`;
        }
      } else if (user && user.isGuest) {
        // User is a guest
        guestSection.classList.remove('hidden');
        loggedInSection.classList.add('hidden');
        guestName.textContent = user.displayName;
      } else {
        // No user (loading state)
        guestSection.classList.add('hidden');
        loggedInSection.classList.add('hidden');
      }
    });
  }

  setupInfoModals() {
    // About modal
    const aboutButton = document.getElementById('about-button');
    const aboutModal = document.getElementById('about-modal');
    const aboutClose = document.getElementById('about-close');
    const aboutBackdrop = aboutModal?.querySelector('.modal-backdrop');

    if (aboutButton && aboutModal) {
      aboutButton.addEventListener('click', () => {
        aboutModal.classList.remove('hidden');
      });

      const closeAbout = () => aboutModal.classList.add('hidden');
      aboutClose?.addEventListener('click', closeAbout);
      aboutBackdrop?.addEventListener('click', closeAbout);
    }

    // Credits modal
    const creditsButton = document.getElementById('credits-button');
    const creditsModal = document.getElementById('credits-modal');
    const creditsClose = document.getElementById('credits-close');
    const creditsBackdrop = creditsModal?.querySelector('.modal-backdrop');

    if (creditsButton && creditsModal) {
      creditsButton.addEventListener('click', () => {
        creditsModal.classList.remove('hidden');
      });

      const closeCredits = () => creditsModal.classList.add('hidden');
      creditsClose?.addEventListener('click', closeCredits);
      creditsBackdrop?.addEventListener('click', closeCredits);
    }

    // How It Was Built modal
    const howBuiltButton = document.getElementById('how-built-button');
    const howBuiltModal = document.getElementById('how-built-modal');
    const howBuiltClose = document.getElementById('how-built-close');
    const howBuiltBackdrop = howBuiltModal?.querySelector('.modal-backdrop');

    if (howBuiltButton && howBuiltModal) {
      howBuiltButton.addEventListener('click', () => {
        howBuiltModal.classList.remove('hidden');
      });

      const closeHowBuilt = () => howBuiltModal.classList.add('hidden');
      howBuiltClose?.addEventListener('click', closeHowBuilt);
      howBuiltBackdrop?.addEventListener('click', closeHowBuilt);
    }
  }
}
