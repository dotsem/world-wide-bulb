<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		title,
		href,
		children
	}: {
		title: string;
		href?: string;
		children: Snippet;
	} = $props();
</script>

<div class="space-y-1.5">
	{#if href}
		<a
			{href}
			class="block text-[11px] font-bold uppercase tracking-wider text-app-muted hover:text-app-text px-3 py-1 transition-colors"
		>
			{title}
		</a>
	{:else}
		<div class="text-[11px] font-bold uppercase tracking-wider text-app-muted px-3 py-1">
			{title}
		</div>
	{/if}
	<div class="space-y-1 pl-1">
		{@render children()}
	</div>
</div>
