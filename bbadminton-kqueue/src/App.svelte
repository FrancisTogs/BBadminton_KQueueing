<script lang="ts">
  import { onMount } from 'svelte';

  let players = $state<{id: number, name: string, level: string}[]>([]);
  let newPlayerName = $state('');
  let newPlayerLevel = $state('Low Beginner');

  onMount(async () => {
    await loadPlayers();
  });

  async function loadPlayers() {
    try {
      const response = await fetch('http://localhost:8080/api/players');
      const data = await response.json();
      players = data || []; // Fallback to empty array if DB is empty
    } catch (error) {
      console.error("Failed to fetch players:", error);
    }
  }

  async function addPlayer() {
    if (!newPlayerName.trim()) return;

    try {
      const response = await fetch('http://localhost:8080/api/players', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: newPlayerName,
          level: newPlayerLevel
        })
      });

      if (response.ok) {
        const addedPlayer = await response.json();
        players = [...players, addedPlayer];
        newPlayerName = '';
      }
    } catch (error) {
      console.error("Failed to add player:", error);
    }
  }

  async function removePlayer(id: number) {
    try {
      const response = await fetch(`http://localhost:8080/api/players?id=${id}`, {
        method: 'DELETE'
      });

      if (response.ok) {
        // Filter out the removed player to update the local state
        players = players.filter(player => player.id !== id);
      }
    } catch (error) {
      console.error("Failed to remove player:", error);
    }
  }
</script>

<main class="container">
  <h2>Live Player Registry</h2>

  <div class="add-player-box">
    <input type="text" bind:value={newPlayerName} placeholder="Player Name" />
    <select bind:value={newPlayerLevel}>
      <option>Low Beginner</option>
      <option>High Beginner</option>
      <option>Low Intermediate</option>
      <option>High Intermediate</option>
      <option>Low Advanced</option>
      <option>High Advanced</option>
    </select>
    <button onclick={addPlayer}>Register</button>
  </div>

  <ul class="player-list">
    <!-- Keying the loop by ID improves rendering performance during deletions -->
    {#each players as player (player.id)}
      <li>
        <div>
          <strong>{player.name}</strong> - {player.level} 
          <span class="id-badge">ID: {player.id}</span>
        </div>
        <button class="remove-btn" onclick={() => removePlayer(player.id)}>Remove</button>
      </li>
    {/each}
  </ul>
</main>

<style>
  .container { max-width: 500px; margin: 2rem auto; font-family: system-ui; }
  .add-player-box { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
  input, select { padding: 0.5rem; flex: 1; }
  button { background: #3498db; color: white; border: none; padding: 0.5rem 1rem; cursor: pointer; border-radius: 4px; }
  .player-list { list-style: none; padding: 0; }
  .player-list li { background: #f8f9fa; margin-bottom: 0.5rem; padding: 1rem; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
  .id-badge { color: #7f8c8d; font-size: 0.8rem; margin-left: 0.5rem; }
  .remove-btn { background: #e74c3c; padding: 0.25rem 0.5rem; font-size: 0.9rem; }
</style>