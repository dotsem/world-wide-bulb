<script lang="ts">
	import { onMount } from 'svelte';
	import { LoaderCircle, BadgeQuestionMark } from '@lucide/svelte';
	import { bulbState, historyState, HistoryTile, Bulb } from '$lib';

	let sentinelEl = $state<HTMLElement | null>(null);

	$effect(() => {
		if (!sentinelEl) return;

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && historyState.hasMore && !historyState.loading) {
					historyState.loadMore();
				}
			},
			{ rootMargin: '300px' }
		);

		observer.observe(sentinelEl);
		return () => observer.disconnect();
	});

	onMount(() => {
		historyState.refresh();
	});
</script>

<svelte:head>
	<title>Toggle History - World Wide Bulb</title>
	<meta
		name="description"
		content="View historical toggle events and community reasons for World Wide Bulb."
	/>
</svelte:head>

<div
	class="sticky top-0 z-40 px-4 py-4 sm:px-8 bg-app-surface-solid mb-8 flex items-center justify-between border-b border-app-border overflow-hidden"
>
	<div>
		<h1 class="text-3xl font-extrabold tracking-tight text-app-text flex items-center gap-3">
			Toggle History
		</h1>
		<p class="text-sm text-app-muted mt-1">
			Total: {historyState.totalToggles} toggles
		</p>
	</div>
</div>
<div class="fixed -top-24 -right-24 z-40 pointer-events-none shrink-0">
	<Bulb
		isOn={bulbState.isOn}
		disabled={true}
		durationMs={300}
		class="h-72 w-72 rotate-215 transition-all"
	/>
</div>

<div class="w-full px-4 sm:px-8 mb-20">
	{#if historyState.initialLoading}
		<div class="space-y-3.5">
			{#each Array(5) as _, i (i)}
				<div
					class="animate-pulse flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl border border-app-border bg-app-surface/40 p-4 sm:p-5"
				>
					<div class="h-11 w-11 shrink-0 rounded-xl bg-app-border/40"></div>
					<div class="flex flex-col gap-2 min-w-0 flex-1">
						<div class="flex items-center gap-2">
							<div class="h-3.5 w-28 rounded-md bg-app-border/50"></div>
							<div class="h-3 w-16 rounded-md bg-app-border/30"></div>
						</div>
						<div class="h-4 w-3/5 rounded-md bg-app-border/40"></div>
					</div>
					<div class="h-4 w-10 shrink-0 rounded-md bg-app-border/40"></div>
				</div>
			{/each}
		</div>
	{:else if historyState.error && historyState.items.length === 0}
		<div
			class="rounded-2xl border border-app-danger/30 bg-app-surface-solid p-8 text-center text-app-danger"
		>
			<p class="font-medium">{historyState.error}</p>
			<button
				onclick={() => historyState.refresh()}
				class="mt-4 rounded-xl bg-app-danger/20 cursor-pointer px-5 py-2.5 text-sm font-semibold text-app-danger hover:bg-app-danger/30 transition-all shadow-sm"
			>
				Try Again
			</button>
		</div>
	{:else if historyState.items.length === 0}
		<div class="rounded-2xl border border-app-border bg-app-surface-solid p-12 text-center">
			<div
				class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-app-accent/10 text-app-accent mb-4 border border-app-accent/20"
			>
				<BadgeQuestionMark class="h-7 w-7" />
			</div>
			<h3 class="text-lg font-semibold text-app-text">No one toggled the bulb, I guess</h3>
		</div>
	{:else}
		<div class="space-y-3.5">
			{#each historyState.items as item (item.id)}
				<HistoryTile {item} />
			{/each}
		</div>

		<div bind:this={sentinelEl} class="mt-10 flex justify-center py-4 min-h-14">
			{#if historyState.loading}
				<div
					class="flex items-center gap-2.5 text-sm text-app-muted bg-app-surface border border-app-border px-4 py-2 rounded-full shadow-sm"
				>
					<LoaderCircle class="h-5 w-5 animate-spin text-app-accent" />
					<span>Loading more history...</span>
				</div>
			{:else if !historyState.hasMore}
				<div class="text-xs text-app-muted/60 italic flex items-center gap-3">
					<span class="h-px w-12 bg-app-border"></span>
					<span>Yes, the lamp started toggled off :)</span>
					<span class="h-px w-12 bg-app-border"></span>
				</div>
			{/if}
		</div>
	{/if}
</div>
