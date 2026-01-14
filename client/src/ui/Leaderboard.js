export class Leaderboard {
  constructor() {
    this.loading = document.getElementById('leaderboard-loading');
    this.table = document.getElementById('leaderboard-table');
    this.tbody = document.getElementById('leaderboard-body');
    this.empty = document.getElementById('leaderboard-empty');
  }

  async fetch() {
    try {
      this.loading.classList.remove('hidden');
      this.table.classList.add('hidden');
      this.empty.classList.add('hidden');

      const response = await fetch('/api/leaderboard');
      if (!response.ok) {
        throw new Error('Failed to fetch leaderboard');
      }

      const entries = await response.json();
      this.render(entries);
    } catch (error) {
      console.error('Error fetching leaderboard:', error);
      this.loading.textContent = 'Failed to load leaderboard';
    }
  }

  render(entries) {
    this.loading.classList.add('hidden');

    if (!entries || entries.length === 0) {
      this.empty.classList.remove('hidden');
      return;
    }

    // Clear existing rows
    this.tbody.innerHTML = '';

    // Add rows for each entry
    entries.forEach((entry, index) => {
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

    this.table.classList.remove('hidden');
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
}
