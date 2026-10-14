<script lang="ts">
	import type { Component, Snippet } from 'svelte';

	let {
		id,
		title,
		description,
		icon: Icon,
		children
	}: {
		id: string;
		title: string;
		description: string;
		icon?: Component<{ size?: number; class?: string }>;
		children: Snippet;
	} = $props();
</script>

<section {id} class="space-y-6 scroll-mt-28">
	<div class="border-b border-app-border pb-4">
		<div class="flex items-center gap-2 text-app-accent mb-1">
			{#if Icon}
				<Icon size={20} />
			{/if}
			<h2 class="text-xl font-extrabold text-app-text">{title}</h2>
		</div>
		<p class="text-xs sm:text-sm text-app-muted leading-relaxed">
			{description}
		</p>
	</div>

	<div class="space-y-8">
		{@render children()}
	</div>
</section>
