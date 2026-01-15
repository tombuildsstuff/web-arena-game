export class GameOverScreen {
  constructor(gameState, onPlayAgain) {
    this.gameState = gameState;
    this.onPlayAgain = onPlayAgain;
    this.screen = document.getElementById('game-over-screen');
    this.title = document.getElementById('game-over-title');
    this.message = document.getElementById('game-over-message');
    this.playAgainButton = document.getElementById('play-again-button');
    this.matchDuration = document.getElementById('match-duration');

    // Detailed stats elements
    this.yourPoints = document.getElementById('your-points');
    this.yourTankKills = document.getElementById('your-tank-kills');
    this.yourAirplaneKills = document.getElementById('your-airplane-kills');
    this.yourSniperKills = document.getElementById('your-sniper-kills');
    this.yourRocketKills = document.getElementById('your-rocket-kills');
    this.yourTurretKills = document.getElementById('your-turret-kills');
    this.yourBarracksKills = document.getElementById('your-barracks-kills');
    this.yourPlayerKills = document.getElementById('your-player-kills');

    this.enemyPoints = document.getElementById('enemy-points');
    this.enemyTankKills = document.getElementById('enemy-tank-kills');
    this.enemyAirplaneKills = document.getElementById('enemy-airplane-kills');
    this.enemySniperKills = document.getElementById('enemy-sniper-kills');
    this.enemyRocketKills = document.getElementById('enemy-rocket-kills');
    this.enemyTurretKills = document.getElementById('enemy-turret-kills');
    this.enemyBarracksKills = document.getElementById('enemy-barracks-kills');
    this.enemyPlayerKills = document.getElementById('enemy-player-kills');

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

    // Display detailed stats
    if (stats) {
      const myStats = myPlayerId === 0 ? stats.player1Stats : stats.player2Stats;
      const theirStats = myPlayerId === 0 ? stats.player2Stats : stats.player1Stats;

      // Your stats
      if (myStats) {
        this.yourPoints.textContent = myStats.totalPoints || 0;
        this.yourTankKills.textContent = myStats.tankKills || 0;
        this.yourAirplaneKills.textContent = myStats.airplaneKills || 0;
        this.yourSniperKills.textContent = myStats.sniperKills || 0;
        this.yourRocketKills.textContent = myStats.rocketLauncherKills || 0;
        this.yourTurretKills.textContent = myStats.turretKills || 0;
        this.yourBarracksKills.textContent = myStats.barracksKills || 0;
        this.yourPlayerKills.textContent = myStats.playerKills || 0;
      }

      // Enemy stats
      if (theirStats) {
        this.enemyPoints.textContent = theirStats.totalPoints || 0;
        this.enemyTankKills.textContent = theirStats.tankKills || 0;
        this.enemyAirplaneKills.textContent = theirStats.airplaneKills || 0;
        this.enemySniperKills.textContent = theirStats.sniperKills || 0;
        this.enemyRocketKills.textContent = theirStats.rocketLauncherKills || 0;
        this.enemyTurretKills.textContent = theirStats.turretKills || 0;
        this.enemyBarracksKills.textContent = theirStats.barracksKills || 0;
        this.enemyPlayerKills.textContent = theirStats.playerKills || 0;
      }
    }

    this.screen.classList.remove('hidden');
  }

  hide() {
    this.screen.classList.add('hidden');
  }
}
