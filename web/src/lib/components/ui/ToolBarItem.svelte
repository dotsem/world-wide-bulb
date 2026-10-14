<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { ResolvedPathname } from '$app/types';

	let {
		href,
		active = false,
		ariaLabel,
		title,
		children
	}: {
		href: ResolvedPathname;
		active?: boolean;
		ariaLabel?: string;
		title?: string;
		children?: Snippet;
	} = $props();
</script>

<a
	{href}
	aria-label={ariaLabel}
	{title}
	class={`flex items-center justify-center gap-2 py-2 px-2 sm:px-4 rounded-full text-xs sm:text-sm font-medium transition-all duration-200 ${
		active
			? 'text-app-accent shadow-inner'
			: 'text-app-muted hover:text-app-text hover:bg-app-surface-hover/60'
	}`}
>
	{@render children?.()}
</a>
