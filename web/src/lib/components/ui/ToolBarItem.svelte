<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		href,
		active = false,
		children
	}: {
		href: string;
		active?: boolean;
		children?: Snippet;
	} = $props();
</script>

<a
	{href}
	class={`flex items-center justify-center gap-2 py-2 px-4 rounded-full text-xs sm:text-sm font-medium transition-all duration-200 ${
		active
			? 'text-amber-300 font-semibold shadow-inner'
			: 'text-app-muted hover:text-app-text hover:bg-app-surface-hover/60'
	}`}
>
	{@render children?.()}
</a>
