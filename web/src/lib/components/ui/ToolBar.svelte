<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		activeIndex = 0,
		count = 2,
		children
	}: {
		activeIndex?: number;
		count?: number;
		children?: Snippet;
	} = $props();
</script>

<nav
	class="relative flex items-center p-1.5 rounded-full bg-app-surface border border-app-border shadow-xl shadow-black/40 backdrop-blur-xl w-full transition-colors duration-300"
>
	<div
		class="absolute inset-1.5 grid pointer-events-none"
		style="grid-template-columns: repeat({count}, minmax(0, 1fr));"
	>
		<div
			class="h-full w-full rounded-full bg-amber-500/20 border border-amber-500/40 shadow-sm shadow-amber-500/10 transition-transform duration-300 ease-out"
			style="transform: translateX({activeIndex * 100}%);"
		></div>
	</div>

	<div
		class="relative z-10 grid w-full"
		style="grid-template-columns: repeat({count}, minmax(0, 1fr));"
	>
		{@render children?.()}
	</div>
</nav>
