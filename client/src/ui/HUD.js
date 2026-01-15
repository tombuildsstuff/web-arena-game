export class HUD {
  constructor(gameState, gameLoop, soundManager = null) {
    this.gameState = gameState;
    this.gameLoop = gameLoop;
    this.soundManager = soundManager;
    this.playerIdElement = document.getElementById('player-id');
    this.playerMoneyElement = document.getElementById('player-money');
    this.healthBarElement = document.getElementById('health-bar');
    this.healthTextElement = document.getElementById('health-text');
    this.respawnTimerElement = document.getElementById('respawn-timer');
    this.buyZonePromptElement = document.getElementById('buy-zone-prompt');
    this.buyZoneTextElement = document.getElementById('buy-zone-text');
    this.muteButton = document.getElementById('mute-button');
    this.controlsOverlay = document.getElementById('controls-overlay');

    // Unit count elements
    this.myTanksElement = document.getElementById('my-tanks');
    this.myHelicoptersElement = document.getElementById('my-helicopters');
    this.enemyTanksElement = document.getElementById('enemy-tanks');
    this.enemyHelicoptersElement = document.getElementById('enemy-helicopters');

    this.setupMuteButton();
    this.setupHelpToggle();
  }

  setupMuteButton() {
    if (!this.muteButton) return;

    this.muteButton.addEventListener('click', () => {
      if (this.soundManager) {
        const isMuted = this.soundManager.toggleMute();
        this.updateMuteButtonIcon(isMuted);
      }
    });

    // Initialize icon state
    this.updateMuteButtonIcon(this.soundManager?.getMuted() || false);
  }

  updateMuteButtonIcon(isMuted) {
    if (!this.muteButton) return;
    this.muteButton.textContent = isMuted ? '🔇' : '🔊';
    this.muteButton.title = isMuted ? 'Unmute sounds' : 'Mute sounds';
  }

  setupHelpToggle() {
    document.addEventListener('keydown', (e) => {
      if (e.key.toLowerCase() === 'h' && !e.repeat) {
        // Don't toggle if typing in an input
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
        this.toggleHelp();
      }
    });
  }

  toggleHelp() {
    if (this.controlsOverlay) {
      this.controlsOverlay.classList.toggle('hidden');
    }
  }

  hideHelp() {
    if (this.controlsOverlay) {
      this.controlsOverlay.classList.add('hidden');
    }
  }

  setSoundManager(soundManager) {
    this.soundManager = soundManager;
    this.updateMuteButtonIcon(soundManager?.getMuted() || false);
  }

  update() {
    const myPlayer = this.gameState.getMyPlayer();
    const myPlayerUnit = this.gameState.getMyPlayerUnit();

    if (myPlayer) {
      // Use displayName if available, fallback to "Player X"
      const displayName = myPlayer.displayName || `Player ${this.gameState.playerId + 1}`;
      this.playerIdElement.textContent = displayName;
      this.playerMoneyElement.textContent = `$${myPlayer.money}`;
    } else {
      this.playerIdElement.textContent = '-';
      this.playerMoneyElement.textContent = '$0';
    }

    // Update player unit health and respawn status
    if (myPlayerUnit) {
      const maxHealth = 80;
      const currentHealth = myPlayerUnit.health || 0;
      const healthPercent = Math.max(0, Math.min(100, (currentHealth / maxHealth) * 100));

      // Update health bar width
      if (this.healthBarElement) {
        this.healthBarElement.style.width = `${healthPercent}%`;

        // Update health bar color based on health level
        this.healthBarElement.classList.remove('low', 'critical', 'dead');
        if (myPlayerUnit.isRespawning) {
          this.healthBarElement.classList.add('dead');
        } else if (healthPercent <= 25) {
          this.healthBarElement.classList.add('critical');
        } else if (healthPercent <= 50) {
          this.healthBarElement.classList.add('low');
        }
      }

      // Update health text
      if (this.healthTextElement) {
        this.healthTextElement.textContent = `${currentHealth}/${maxHealth}`;
      }

      if (myPlayerUnit.isRespawning && myPlayerUnit.respawnTime > 0) {
        this.respawnTimerElement.textContent = `Respawning: ${Math.ceil(myPlayerUnit.respawnTime)}s`;
        this.respawnTimerElement.classList.remove('hidden');
      } else {
        this.respawnTimerElement.classList.add('hidden');
      }
    } else {
      if (this.healthBarElement) {
        this.healthBarElement.style.width = '0%';
      }
      if (this.healthTextElement) {
        this.healthTextElement.textContent = '-/-';
      }
      this.respawnTimerElement.classList.add('hidden');
    }

    // Update buy zone prompt
    this.updateBuyZonePrompt();

    // Update unit counts
    this.updateUnitCounts();
  }

  updateUnitCounts() {
    const playerId = this.gameState.playerId;
    if (playerId === null) return;

    const units = this.gameState.units || [];

    // Count units by type and owner
    let myTanks = 0;
    let myHelicopters = 0;
    let enemyTanks = 0;
    let enemyHelicopters = 0;

    for (const unit of units) {
      // Skip player units and dead units
      if (unit.type === 'player' || !unit.health || unit.health <= 0) continue;

      const isMyUnit = unit.ownerId === playerId;

      if (unit.type === 'tank' || unit.type === 'super_tank') {
        if (isMyUnit) {
          myTanks++;
        } else {
          enemyTanks++;
        }
      } else if (unit.type === 'airplane' || unit.type === 'super_helicopter') {
        if (isMyUnit) {
          myHelicopters++;
        } else {
          enemyHelicopters++;
        }
      }
    }

    // Update DOM
    if (this.myTanksElement) this.myTanksElement.textContent = myTanks;
    if (this.myHelicoptersElement) this.myHelicoptersElement.textContent = myHelicopters;
    if (this.enemyTanksElement) this.enemyTanksElement.textContent = enemyTanks;
    if (this.enemyHelicoptersElement) this.enemyHelicoptersElement.textContent = enemyHelicopters;
  }

  updateBuyZonePrompt() {
    if (!this.gameLoop || !this.buyZonePromptElement) return;

    const nearbyZone = this.gameLoop.getNearbyBuyZone();
    const nearbyClaimableZone = this.gameLoop.getNearbyClaimableBuyZone();
    const nearbyTurret = this.gameLoop.getNearbyTurret();
    const myPlayer = this.gameState.getMyPlayer();
    const pendingCounts = this.gameState.getMyPendingSpawnCounts();

    // Priority: owned buy zones first, then claimable zones, then turrets
    if (nearbyZone && myPlayer) {
      const unitNames = {
        'tank': 'Tank',
        'airplane': 'Helicopter',
        'super_tank': 'Super Tank',
        'super_helicopter': 'Super Helicopter'
      };
      const unitName = unitNames[nearbyZone.unitType] || nearbyZone.unitType;
      const cost = nearbyZone.cost;
      const canAfford = myPlayer.money >= cost;
      const pendingCount = pendingCounts[nearbyZone.unitType] || 0;

      let promptText = canAfford
        ? `Press <span class="key">E</span> to buy ${unitName} ($${cost})`
        : `${unitName} ($${cost}) - Not enough money`;

      // Show pending spawn count if any
      if (pendingCount > 0) {
        promptText += `<br><span class="pending-spawn">${pendingCount} ${unitName}${pendingCount > 1 ? 's' : ''} waiting to spawn...</span>`;
      }

      this.buyZoneTextElement.innerHTML = promptText;

      this.buyZonePromptElement.classList.remove('hidden');
      this.buyZonePromptElement.style.borderColor = canAfford
        ? 'rgba(251, 191, 36, 0.6)'
        : 'rgba(239, 68, 68, 0.6)';
    } else if (nearbyClaimableZone && myPlayer) {
      // Show claim prompt for neutral claimable buy zones
      const claimCost = nearbyClaimableZone.claimCost;
      const canAfford = myPlayer.money >= claimCost;

      this.buyZoneTextElement.innerHTML = canAfford
        ? `Press <span class="key">E</span> to claim Forward Base ($${claimCost})`
        : `Forward Base ($${claimCost}) - Not enough money`;

      this.buyZonePromptElement.classList.remove('hidden');
      this.buyZonePromptElement.style.borderColor = canAfford
        ? 'rgba(136, 136, 136, 0.6)'  // Grey for claimable bases (matches zone color)
        : 'rgba(239, 68, 68, 0.6)';   // Red if can't afford
    } else if (nearbyTurret) {
      // Show turret claiming prompt (only shown for neutral turrets)
      this.buyZoneTextElement.innerHTML = `Press <span class="key">E</span> to claim Turret`;

      this.buyZonePromptElement.classList.remove('hidden');
      this.buyZonePromptElement.style.borderColor = 'rgba(139, 92, 246, 0.6)'; // Purple for turrets
    } else {
      this.buyZonePromptElement.classList.add('hidden');
    }
  }
}
