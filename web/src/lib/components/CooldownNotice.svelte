<script lang="ts">
	import { bulbState } from '$lib/state/bulb.svelte';
	import { Lock } from '@lucide/svelte';
	import { fly } from 'svelte/transition';

	let remainingSeconds = $derived(bulbState.remainingSeconds);
	let isActive = $derived(bulbState.isCooldownActive);

	let isPulsing = $state(false);
	let pulseTimeout: ReturnType<typeof setTimeout> | undefined;

	$effect(() => {
		if (bulbState.cooldownPulseKey > 0) {
			isPulsing = true;
			if (pulseTimeout) clearTimeout(pulseTimeout);
			pulseTimeout = setTimeout(() => {
				isPulsing = false;
			}, 300);
		}
	});
</script>

{#if isActive}
	<div
		transition:fly={{ y: -40, duration: 300 }}
		class="fixed top-0 left-1/2 -translate-x-1/2 z-50 flex items-center gap-2.5 px-8 py-2.5 rounded-b-full bg-app-surface-solid border border-t-0 border-app-border text-app-text shadow-xl shadow-black/40 transition-all duration-300 whitespace-nowrap {isPulsing
			? 'scale-110 border-app-accent text-app-text shadow-app-glow ring-2 ring-app-accent/40'
			: ''}"
		role="status"
		aria-live="polite"
	>
		<Lock class="h-4 w-4 shrink-0 text-app-accent" />
		<span class="text-sm font-medium tracking-wide flex items-center gap-1.5 whitespace-nowrap">
			<span class="whitespace-nowrap">Cooldown active:</span>
			<span class="font-semibold tabular-nums text-app-accent inline-block min-w-[3.5ch] text-left"
				>{remainingSeconds}s</span
			>
		</span>
	</div>
{/if}
