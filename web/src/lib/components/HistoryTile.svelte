<script lang="ts">
	import type { FormattedHistoryItem } from '$lib/state/history.svelte';
	import { Lightbulb, LightbulbOff, MessageSquare, Trophy } from '@lucide/svelte';

	let { item }: { item: FormattedHistoryItem } = $props();

	let isMilestone = $derived(item.id % 100 === 0);
</script>

<div
	class="group relative flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl border border-app-border p-4 sm:p-5 [content-visibility:auto] [contain-intrinsic-size:0_80px] transition-all duration-200 {item.state
		? ' bg-amber-400/20'
		: ' bg-slate-600/30'}"
>
	<div
		class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl font-extrabold text-xs border transition-all z-10 {item.state
			? 'bg-amber-500/20 text-amber-300 border-amber-400/50 shadow-md shadow-amber-500/20 ring-2 ring-amber-400/20'
			: 'bg-slate-800/60 text-slate-400 border-slate-700'}"
	>
		{#if item.state}
			<Lightbulb />
		{:else}
			<LightbulbOff />
		{/if}
	</div>

	<div class="flex flex-col gap-1 min-w-0 flex-1">
		<div class="flex items-center gap-2 flex-wrap">
			{#if isMilestone}
				<span
					class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-amber-400/20 text-amber-300 border border-amber-400/30"
				>
					<Trophy class="h-3 w-3" />
					<span>Milestone #{item.id}</span>
				</span>
			{/if}
			<span class="text-xs text-app-muted font-medium">{item.formattedDate}</span>
			<span class="text-[11px] text-app-muted/70">({item.relativeTime})</span>
		</div>

		{#if item.reason}
			<div class="mt-1 flex items-start gap-2 text-sm">
				<MessageSquare class="h-4 w-4 shrink-0 text-amber-400 mt-0.5" />
				<p class="min-w-0 wrap-break-word italic text-app-text/95">
					"{item.reason}"
				</p>
			</div>
		{:else}
			<span class="text-xs italic text-app-muted/60 mt-0.5">No reason specified</span>
		{/if}
	</div>

	<span class="text-xs font-mono font-semibold text-amber-400/80">#{item.id}</span>
</div>
