export class HUD {
  constructor(gameState, gameLoop, soundManager = null) {
    this.gameState = gameState;
    this.gameLoop = gameLoop;
    this.soundManager = soundManager;
    this.playerIdElement = document.getElementById('player-id');
    this.playerMoneyElement = document.getElementById('player-money');
    this.playerHealthElement = document.getElementById('player-health');
    this.respawnTimerElement = document.getElementById('respawn-timer');
    this.buyZonePromptElement = document.getElementById('buy-zone-prompt');
    this.buyZoneTextElement = document.getElementById('buy-zone-text');
    this.muteButton = document.getElementById('mute-button');

    this.setupMuteButton();
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
      this.playerMoneyElement.textContent = `Money: $${myPlayer.money}`;
    } else {
      this.playerIdElement.textContent = 'Player: -';
      this.playerMoneyElement.textContent = 'Money: $0';
    }

    // Update player unit health and respawn status
    if (myPlayerUnit) {
      this.playerHealthElement.textContent = `Health: ${myPlayerUnit.health}/80`;

      if (myPlayerUnit.isRespawning && myPlayerUnit.respawnTime > 0) {
        this.respawnTimerElement.textContent = `Respawning: ${Math.ceil(myPlayerUnit.respawnTime)}s`;
        this.respawnTimerElement.classList.remove('hidden');
        this.playerHealthElement.classList.add('dead');
      } else {
        this.respawnTimerElement.classList.add('hidden');
        this.playerHealthElement.classList.remove('dead');
      }
    } else {
      this.playerHealthElement.textContent = 'Health: -';
      this.respawnTimerElement.classList.add('hidden');
    }

    // Update buy zone prompt
    this.updateBuyZonePrompt();
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
      const unitName = nearbyZone.unitType === 'tank' ? 'Tank' : 'Airplane';
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
