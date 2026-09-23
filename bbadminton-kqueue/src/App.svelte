<script lang="ts">
  import { onMount } from 'svelte';

  // State is initially empty, we will fill it from the Go API
  let players = $state<{id: number, name: string, level: string}[]>([]);
  
  let newPlayerName = $state('');
  let newPlayerLevel = $state('Low Beginner');

  // Fetch players when the component loads
  onMount(async () => {
    await loadPlayers();
  });

  async function loadPlayers() {
    try {
      const response = await fetch('http://localhost:8080/api/players');
      players = await response.json();
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
        // Update local Svelte state with the database-confirmed player
        players = [...players, addedPlayer];
        newPlayerName = '';
      }
    } catch (error) {
      console.error("Failed to add player:", error);
    }
  }
</script>

<main class="container">
  <h2>Live Player Registry</h2>

  <div class="add-player-box">
    <input type="text" bind:value={newPlayerName} placeholder="Player Name" />
    <select bind:value={newPlayerLevel}>
      <option>Low Beginner</option>
      <option>Low Intermediate</option>
      <option>Low Advanced</option>
    </select>
    <button onclick={addPlayer}>Register</button>
  </div>

  <ul class="player-list">
    {#each players as player}
      <li>
        <strong>{player.name}</strong> - {player.level} 
        <span class="id-badge">ID: {player.id}</span>
      </li>
    {/each}
  </ul>
</main>

<style>
  .container { max-width: 500px; margin: 2rem auto; font-family: system-ui; }
  .add-player-box { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
  input, select { padding: 0.5rem; flex: 1; }
  button { background: #3498db; color: white; border: none; padding: 0.5rem 1rem; cursor: pointer; }
  .player-list { list-style: none; padding: 0; }
  .player-list li { background: #f8f9fa; margin-bottom: 0.5rem; padding: 1rem; border-radius: 4px; display: flex; justify-content: space-between; }
  .id-badge { color: #7f8c8d; font-size: 0.8rem; }
</style>