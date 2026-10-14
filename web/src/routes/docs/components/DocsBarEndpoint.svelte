<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { ResolvedPathname } from '$app/types';

	let {
		href,
		title,
		active = false,
		badge,
		onclick
	}: {
		href: ResolvedPathname;
		title: string;
		active?: boolean;
		badge?: Snippet;
		onclick?: () => void;
	} = $props();
</script>

<a
	{href}
	{onclick}
	class="group flex items-center justify-between gap-2 px-3 py-2 rounded-xl text-xs font-mono transition-all {active
		? 'bg-app-accent/15 text-app-accent font-semibold border border-app-accent/30 shadow-xs'
		: 'text-app-muted hover:text-app-text hover:bg-app-surface border border-transparent'}"
>
	<span class="truncate">{title}</span>
	{#if badge}
		{@render badge()}
	{/if}
</a>
