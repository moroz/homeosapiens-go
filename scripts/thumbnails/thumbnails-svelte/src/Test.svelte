<script lang="ts">
  import Thumbnail from "./Thumbnail.svelte";

  const data = [
    {
      date: "2026-02-08",
      host: "Dr Asher Shaikh",
      locale: "en" as any,
      pp: "/dr-asher.jpeg",
      title: "What prevents me from moving on?",
    },
    {
      date: "2025-10-24",
      host: "Dr Sanjay Modi",
      title: "Udoskonalić kunszt homeopatyczny: Dzień 1",
      locale: "pl",
      pp: "/modi.jpeg",
    },
  ];
</script>

<div class="grid w-full h-screen place-items-center bg-blue-100">
  <div class="grid gap-6">
    {#each data as item}
      <Thumbnail {...item} />
    {/each}
  </div>
</div>
