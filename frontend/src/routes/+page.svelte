<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { getDevices, getStatus } from "$lib/api";

  let devices: any[] = [];
  let selectedDevice: number | null = null;
  let status: any = null;
  let error: string | null = null;
  let interval: ReturnType<typeof setInterval>;

  // Load devices on startup
  onMount(async () => {
    try {
      devices = await getDevices();
      if (devices.length > 0) {
        selectedDevice = devices[0].index;
        await loadStatus();
      }
    } catch (e) {
      error = "Failed to load devices";
    }

    // Poll status every 3 seconds
    interval = setInterval(loadStatus, 3000);
  });

  async function loadStatus() {
    if (selectedDevice === null) return;
    try {
      status = await getStatus(selectedDevice);
    } catch (e) {
      error = "Failed to load status";
    }
  }

  // Cleanup
  onDestroy(() => clearInterval(interval));
</script>

<main class="p-6 space-y-6">
  <h1 class="text-3xl font-bold">Liquid UI Dashboard</h1>

  {#if error}
    <div class="bg-red-500 text-white p-3 rounded">{error}</div>
  {/if}

  <!-- Device List -->
  <section>
    <h2 class="text-xl font-semibold mb-2">Devices</h2>

    {#if devices.length === 0}
      <p>No devices detected</p>
    {:else}
      <ul class="space-y-2">
        {#each devices as d}
          <li>
            <button
              class="px-4 py-2 rounded border w-full text-left
                     {selectedDevice === d.index ? 'bg-blue-600 text-white' : 'bg-gray-100'}"
              on:click={() => { selectedDevice = d.index; loadStatus(); }}
            >
              {d.name} (#{d.index})
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </section>

  <!-- Status Panel -->
  {#if status}
    <section class="mt-6 p-4 border rounded bg-gray-50">
      <h2 class="text-xl font-semibold mb-2">Status</h2>
      <p><strong>Temperature:</strong> {status.temperature} °C</p>
      <p><strong>Fan RPM:</strong> {status.fanRpm}</p>
      <p><strong>Pump RPM:</strong> {status.pumpRpm}</p>
    </section>
  {/if}
</main>

<style>
  main {
    max-width: 800px;
    margin: auto;
  }
</style>
