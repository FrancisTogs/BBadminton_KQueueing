<script lang="ts">
  import { onMount } from 'svelte';

  let players = $state<{id: number, name: string, level: string}[]>([]);
  let newPlayerName = $state('');
  let newPlayerLevel = $state('Low Beginner');

  // Edit State
  let editingId = $state<number | null>(null);
  let editName = $state('');
  let editLevel = $state('');

  onMount(async () => {
    await loadPlayers();
  });

  async function loadPlayers() {
    try {
      const response = await fetch('http://localhost:8080/api/players');
      const data = await response.json();
      players = data || []; 
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
        players = players.filter(player => player.id !== id);
      }
    } catch (error) {
      console.error("Failed to remove player:", error);
    }
  }

  // Edit Mode Triggers
  function startEditing(player: {id: number, name: string, level: string}) {
    editingId = player.id;
    editName = player.name;
    editLevel = player.level;
  }

  function cancelEditing() {
    editingId = null;
  }

  async function saveEdit() {
    if (!editName.trim() || editingId === null) return;

    try {
      const response = await fetch('http://localhost:8080/api/players', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: editingId,
          name: editName,
          level: editLevel
        })
      });

      if (response.ok) {
        const updatedPlayer = await response.json();
        // Swap out the old player data for the new data in the local state
        players = players.map(p => p.id === editingId ? updatedPlayer : p);
        editingId = null; // Exit edit mode
      }
    } catch (error) {
      console.error("Failed to update player:", error);
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
    <!-- Notice the added ", i" here to get the array index -->
    {#each players as player, i (player.id)}
      <li>
        <span class="queue-number">#{i + 1}</span>

        {#if editingId === player.id}
          <!-- Edit Form Mode -->
          <div class="edit-box">
            <input type="text" bind:value={editName} />
            <select bind:value={editLevel}>
              <option>Low Beginner</option>
              <option>High Beginner</option>
              <option>Low Intermediate</option>
              <option>High Intermediate</option>
              <option>Low Advanced</option>
              <option>High Advanced</option>
            </select>
            <button class="save-btn" onclick={saveEdit}>Save</button>
            <button class="cancel-btn" onclick={cancelEditing}>Cancel</button>
          </div>
        {:else}
          <!-- Standard Display Mode -->
          <div>
            <span 
              class="clickable-name" 
              onclick={() => startEditing(player)} 
              onkeydown={(e) => e.key === 'Enter' && startEditing(player)}
              tabindex="0"
              role="button"
            >
              <strong>{player.name}</strong> - {player.level}
            </span> 
          </div>
          <button class="remove-btn" onclick={() => removePlayer(player.id)}>Remove</button>
        {/if}
      </li>
    {/each}
  </ul>
</main>

<style>
  .container { max-width: 500px; margin: 2rem auto; font-family: system-ui; }
  .add-player-box { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
  input, select { padding: 0.5rem; flex: 1; border: 1px solid #ccc; border-radius: 4px; }
  button { background: #3498db; color: white; border: none; padding: 0.5rem 1rem; cursor: pointer; border-radius: 4px; }
  
  .player-list { list-style: none; padding: 0; }
  .player-list li { background: #f8f9fa; margin-bottom: 0.5rem; padding: 1rem; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
  
  .id-badge { color: #7f8c8d; font-size: 0.8rem; margin-left: 0.5rem; pointer-events: none; }
  .remove-btn { background: #e74c3c; padding: 0.25rem 0.5rem; font-size: 0.9rem; }
  .queue-number { font-weight: bold; color: #2c3e50; margin-right: 1rem; font-size: 1.2rem; }

  /* Edit Mode Styles */
  .clickable-name { cursor: pointer; padding: 0.25rem; border-radius: 4px; transition: background 0.2s; }
  .clickable-name:hover { background: #e0e0e0; }
  .edit-box { display: flex; gap: 0.5rem; width: 100%; align-items: center; }
  .save-btn { background: #2ecc71; padding: 0.25rem 0.5rem; }
  .cancel-btn { background: #95a5a6; padding: 0.25rem 0.5rem; }
</style>