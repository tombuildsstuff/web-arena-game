export class HUD {
  constructor(gameState, gameLoop) {
    this.gameState = gameState;
    this.gameLoop = gameLoop;
    this.playerIdElement = document.getElementById('player-id');
    this.playerMoneyElement = document.getElementById('player-money');
    this.playerHealthElement = document.getElementById('player-health');
    this.respawnTimerElement = document.getElementById('respawn-timer');
    this.buyZonePromptElement = document.getElementById('buy-zone-prompt');
    this.buyZoneTextElement = document.getElementById('buy-zone-text');
  }

  update() {
    const myPlayer = this.gameState.getMyPlayer();
    const myPlayerUnit = this.gameState.getMyPlayerUnit();

    if (myPlayer) {
      this.playerIdElement.textContent = `Player ${this.gameState.playerId + 1}`;
      this.playerMoneyElement.textContent = `Money: $${myPlayer.money}`;
    } else {
      this.playerIdElement.textContent = 'Player: -';
      this.playerMoneyElement.textContent = 'Money: $0';
    }

    // Update player unit health and respawn status
    if (myPlayerUnit) {
      this.playerHealthElement.textContent = `Health: ${myPlayerUnit.health}/5`;

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
    const myPlayer = this.gameState.getMyPlayer();

    if (nearbyZone && myPlayer) {
      const unitName = nearbyZone.unitType === 'tank' ? 'Tank' : 'Airplane';
      const cost = nearbyZone.cost;
      const canAfford = myPlayer.money >= cost;

      this.buyZoneTextElement.innerHTML = canAfford
        ? `Press <span class="key">E</span> to buy ${unitName} ($${cost})`
        : `${unitName} ($${cost}) - Not enough money`;

      this.buyZonePromptElement.classList.remove('hidden');
      this.buyZonePromptElement.style.borderColor = canAfford
        ? 'rgba(251, 191, 36, 0.6)'
        : 'rgba(239, 68, 68, 0.6)';
    } else {
      this.buyZonePromptElement.classList.add('hidden');
    }
  }
}
