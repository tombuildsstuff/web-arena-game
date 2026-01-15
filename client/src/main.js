import { Game } from './game/Game.js';

// Initialize the game when DOM is loaded
document.addEventListener('DOMContentLoaded', async () => {
  console.log('Initializing Arena Game...');

  const game = new Game();
  await game.init();

  console.log('Game initialized!');
});
