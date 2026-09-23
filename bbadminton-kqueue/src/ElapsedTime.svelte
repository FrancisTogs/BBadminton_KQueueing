<script lang="ts">
  import { onMount } from 'svelte';

  let { startUnix } = $props<{ startUnix: number }>();
  let currentTime = $state(Math.floor(Date.now() / 1000));

  onMount(() => {
    const interval = setInterval(() => {
      currentTime = Math.floor(Date.now() / 1000);
    }, 1000);
    return () => clearInterval(interval);
  });

  let formattedTime = $derived(() => {
    if (!startUnix && startUnix !== 0) return "0m 0s";
    const diff = Math.max(0, currentTime - startUnix);
    const m = Math.floor(diff / 60);
    const s = diff % 60;
    return `${m}m ${s}s`;
  });
</script>

<!-- Renders just the text, isolating the DOM update -->
{formattedTime()}