export class GameOverScreen {
  constructor(gameState, onPlayAgain) {
    this.gameState = gameState;
    this.onPlayAgain = onPlayAgain;
    this.screen = document.getElementById('game-over-screen');
    this.title = document.getElementById('game-over-title');
    this.message = document.getElementById('game-over-message');
    this.playAgainButton = document.getElementById('play-again-button');
    this.matchDuration = document.getElementById('match-duration');
    this.yourKills = document.getElementById('your-kills');
    this.enemyKills = document.getElementById('enemy-kills');

    this.setupEventListeners();
  }

  setupEventListeners() {
    this.playAgainButton.addEventListener('click', () => {
      this.onPlayAgain();
    });
  }

  show(winnerId, reason, matchDuration, stats) {
    // Determine if this player won
    const myPlayerId = this.gameState.playerId;
    const didWin = winnerId === myPlayerId;

    console.log('Game over:', { winnerId, myPlayerId, didWin, reason, matchDuration, stats });

    if (didWin) {
      this.title.textContent = 'Victory!';
      this.title.className = 'victory';
      this.message.textContent = reason || 'You won the game!';
    } else {
      this.title.textContent = 'Defeat';
      this.title.className = 'defeat';
      this.message.textContent = reason || 'You lost the game.';
    }

    // Display match duration
    if (matchDuration !== undefined) {
      const minutes = Math.floor(matchDuration / 60);
      const seconds = matchDuration % 60;
      this.matchDuration.textContent = `Match Duration: ${minutes}:${seconds.toString().padStart(2, '0')}`;
    }

    // Display kills
    if (stats) {
      const myKills = myPlayerId === 0 ? stats.player1Kills : stats.player2Kills;
      const theirKills = myPlayerId === 0 ? stats.player2Kills : stats.player1Kills;
      this.yourKills.textContent = myKills;
      this.enemyKills.textContent = theirKills;
    }

    this.screen.classList.remove('hidden');
  }

  hide() {
    this.screen.classList.add('hidden');
  }
}
