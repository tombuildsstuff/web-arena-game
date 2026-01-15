export class Leaderboard {
  constructor() {
    this.loading = document.getElementById('leaderboard-loading');
    this.top3Container = document.getElementById('leaderboard-top3');
    this.viewMoreButton = document.getElementById('leaderboard-view-more');
    this.empty = document.getElementById('leaderboard-empty');

    // Modal elements
    this.modal = document.getElementById('leaderboard-modal');
    this.modalBackdrop = this.modal?.querySelector('.modal-backdrop');
    this.table = document.getElementById('leaderboard-table');
    this.tbody = document.getElementById('leaderboard-body');
    this.closeButton = document.getElementById('leaderboard-close');

    // Store entries for modal
    this.entries = [];

    this.setupEventListeners();
  }

  setupEventListeners() {
    // View more button
    if (this.viewMoreButton) {
      this.viewMoreButton.addEventListener('click', () => this.showModal());
    }

    // Close button
    if (this.closeButton) {
      this.closeButton.addEventListener('click', () => this.hideModal());
    }

    // Backdrop click
    if (this.modalBackdrop) {
      this.modalBackdrop.addEventListener('click', () => this.hideModal());
    }
  }

  async fetch() {
    try {
      this.loading.classList.remove('hidden');
      this.top3Container.classList.add('hidden');
      this.viewMoreButton.classList.add('hidden');
      this.empty.classList.add('hidden');

      const response = await fetch('/api/leaderboard');
      if (!response.ok) {
        throw new Error('Failed to fetch leaderboard');
      }

      this.entries = await response.json();
      this.render();
    } catch (error) {
      console.error('Error fetching leaderboard:', error);
      this.loading.textContent = 'Failed to load leaderboard';
    }
  }

  render() {
    this.loading.classList.add('hidden');

    if (!this.entries || this.entries.length === 0) {
      this.empty.classList.remove('hidden');
      return;
    }

    // Render compact top 3 view
    this.renderTop3();

    // Show "View More" if there are more than 3 entries
    if (this.entries.length > 3) {
      this.viewMoreButton.classList.remove('hidden');
    }

    // Render full table in modal
    this.renderFullTable();
  }

  renderTop3() {
    const top3 = this.entries.slice(0, 3);
    this.top3Container.innerHTML = '';

    const medals = ['gold', 'silver', 'bronze'];
    const medalEmojis = ['🥇', '🥈', '🥉'];

    top3.forEach((entry, index) => {
      const item = document.createElement('div');
      item.className = `top3-item ${medals[index]}`;

      item.innerHTML = `
        <span class="top3-rank">${medalEmojis[index]}</span>
        <span class="top3-name">${this.escapeHtml(entry.playerName)}</span>
        <span class="top3-points">${entry.totalPoints.toLocaleString()} pts</span>
      `;

      this.top3Container.appendChild(item);
    });

    this.top3Container.classList.remove('hidden');
  }

  renderFullTable() {
    // Clear existing rows
    this.tbody.innerHTML = '';

    // Add rows for each entry
    this.entries.forEach((entry, index) => {
      const row = document.createElement('tr');

      // Rank
      const rankCell = document.createElement('td');
      rankCell.textContent = index + 1;
      rankCell.className = 'rank';
      if (index === 0) rankCell.classList.add('gold');
      else if (index === 1) rankCell.classList.add('silver');
      else if (index === 2) rankCell.classList.add('bronze');
      row.appendChild(rankCell);

      // Player name
      const nameCell = document.createElement('td');
      nameCell.textContent = entry.playerName;
      nameCell.className = 'player-name';
      row.appendChild(nameCell);

      // Points
      const pointsCell = document.createElement('td');
      pointsCell.textContent = entry.totalPoints.toLocaleString();
      pointsCell.className = 'points';
      row.appendChild(pointsCell);

      // Wins
      const winsCell = document.createElement('td');
      winsCell.textContent = entry.gamesWon;
      row.appendChild(winsCell);

      // Games
      const gamesCell = document.createElement('td');
      gamesCell.textContent = entry.gamesPlayed;
      row.appendChild(gamesCell);

      // Play time
      const timeCell = document.createElement('td');
      timeCell.textContent = this.formatPlayTime(entry.totalPlayTime);
      row.appendChild(timeCell);

      this.tbody.appendChild(row);
    });
  }

  showModal() {
    if (this.modal) {
      this.modal.classList.remove('hidden');
    }
  }

  hideModal() {
    if (this.modal) {
      this.modal.classList.add('hidden');
    }
  }

  formatPlayTime(seconds) {
    if (seconds < 60) {
      return `${seconds}s`;
    } else if (seconds < 3600) {
      const mins = Math.floor(seconds / 60);
      const secs = seconds % 60;
      return `${mins}m ${secs}s`;
    } else {
      const hours = Math.floor(seconds / 3600);
      const mins = Math.floor((seconds % 3600) / 60);
      return `${hours}h ${mins}m`;
    }
  }

  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
}
