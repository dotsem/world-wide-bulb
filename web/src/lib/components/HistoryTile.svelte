<script lang="ts">
	import type { FormattedHistoryItem } from '$lib/state/history.svelte';
	import { Lightbulb, LightbulbOff, MessageSquare, Trophy } from '@lucide/svelte';

	let { item }: { item: FormattedHistoryItem } = $props();

	let isMilestone = $derived(item.id % 100 === 0);
</script>

<div
	class="group relative flex items-center justify-between gap-3 sm:gap-4 rounded-2xl border border-app-border p-3.5 sm:p-5 [content-visibility:auto] [contain-intrinsic-size:0_80px] transition-all duration-200 {item.state
		? ' bg-app-accent/15'
		: ' bg-app-surface/40'}"
>
	<div
		class="flex h-10 w-10 sm:h-11 sm:w-11 shrink-0 items-center justify-center rounded-xl font-extrabold text-xs border transition-all z-10 {item.state
			? 'bg-app-accent/20 text-app-accent border-app-accent/50 shadow-md shadow-app-accent/20 ring-2 ring-app-accent/20'
			: 'bg-app-surface text-app-muted border-app-border'}"
	>
		{#if item.state}
			<Lightbulb class="h-5 w-5" />
		{:else}
			<LightbulbOff class="h-5 w-5" />
		{/if}
	</div>

	<div class="flex flex-col gap-1 min-w-0 flex-1 justify-center">
		{#if isMilestone}
			<div>
				<span
					class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-app-accent/20 text-app-accent border border-app-accent/30"
				>
					<Trophy class="h-3 w-3" />
					<span>Milestone #{item.id}</span>
				</span>
			</div>
		{/if}

		{#if item.reason}
			<div class="flex items-center gap-1.5 text-xs sm:text-sm">
				<MessageSquare class="h-4 w-4 shrink-0 text-app-accent" />
				<p class="min-w-0 wrap-break-word italic text-app-text/95">
					"{item.reason}"
				</p>
			</div>
		{:else}
			<span class="text-xs italic text-app-muted/60">No reason specified</span>
		{/if}
	</div>

	<div class="flex flex-col items-end text-right shrink-0">
		<span class="text-xs font-mono font-semibold text-app-accent/80">#{item.id}</span>
		<span class="text-[11px] text-app-muted font-medium">{item.formattedDate}</span>
		<span class="text-[10px] text-app-muted/70">({item.relativeTime})</span>
	</div>
</div>
