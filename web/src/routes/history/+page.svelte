<script lang="ts">
	import { onMount } from 'svelte';
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
		historyState.loadMore();
		const cleanupWs = historyState.initWsListeners();
		return cleanupWs;
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
	class="sticky top-0 z-40 px-4 py-4 sm:px-8 bg-app-surface/50 backdrop-blur-xl mb-8 flex items-center justify-between border-b border-app-border overflow-hidden"
>
	<div>
		<h1 class="text-3xl font-extrabold tracking-tight text-app-text flex items-center gap-3">
			Toggle History
		</h1>
		<p class="text-sm text-app-muted mt-1">
			The bulb has been toggled {historyState.totalToggles} times
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
		<div class="space-y-4">
			{#each Array(5) as _}
				<div
					class="h-20 w-full animate-pulse rounded-2xl bg-app-surface border border-app-border/50"
				></div>
			{/each}
		</div>
	{:else if historyState.error && historyState.items.length === 0}
		<div
			class="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-8 text-center text-rose-300 backdrop-blur-md"
		>
			<p class="font-medium">{historyState.error}</p>
			<button
				onclick={() => historyState.loadMore()}
				class="mt-4 rounded-xl bg-rose-500/20 px-5 py-2.5 text-sm font-semibold text-rose-200 hover:bg-rose-500/30 transition-all shadow-sm"
			>
				Try Again
			</button>
		</div>
	{:else if historyState.items.length === 0}
		<div
			class="rounded-2xl border border-app-border bg-app-surface/60 p-12 text-center backdrop-blur-md"
		>
			<div
				class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-amber-500/10 text-amber-400 mb-4 border border-amber-500/20"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-7 w-7"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
					/>
				</svg>
			</div>
			<h3 class="text-lg font-semibold text-app-text">No toggle history found</h3>
			<p class="mt-1 text-sm text-app-muted">Be the first to flip the switch!</p>
			<a
				href="/"
				class="mt-5 inline-block rounded-xl bg-amber-500/20 border border-amber-500/40 px-5 py-2.5 text-sm font-semibold text-amber-300 hover:bg-amber-500/30 transition-all shadow-md shadow-amber-500/10"
			>
				Go to Bulb
			</a>
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
					<svg
						class="h-5 w-5 animate-spin text-amber-400"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					<span>Loading more history...</span>
				</div>
			{:else if !historyState.hasMore}
				<div class="text-xs text-app-muted/60 italic flex items-center gap-3">
					<span class="h-px w-12 bg-app-border"></span>
					<span>You scrolled far, maybe too far...</span>
					<span class="h-px w-12 bg-app-border"></span>
				</div>
			{/if}
		</div>
	{/if}
</div>
