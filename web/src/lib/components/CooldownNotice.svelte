<script lang="ts">
    import { bulbState } from "$lib/state/bulb.svelte";

    let remainingSeconds = $derived(bulbState.remainingSeconds);
    let isActive = $derived(bulbState.isCooldownActive);
</script>

{#if isActive}
    <div
        class="fixed top-6 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 px-5 py-2.5 rounded-full bg-amber-950/80 border border-amber-500/30 text-amber-200 shadow-lg shadow-amber-950/40 backdrop-blur-md transition-all duration-300 animate-in fade-in slide-in-from-top-4"
        role="status"
        aria-live="polite"
    >
        <span class="relative flex h-3 w-3">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 bg-amber-500"></span>
        </span>

        <span class="text-sm font-medium tracking-wide">
            Cooldown active: <span class="font-semibold tabular-nums text-amber-300">{remainingSeconds}s</span>
        </span>
    </div>
{/if}
