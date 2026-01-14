import { Game } from './game/Game.js';

// Initialize the game when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  console.log('Initializing Arena Game...');

  const game = new Game();
  game.init();

  console.log('Game initialized!');
});
