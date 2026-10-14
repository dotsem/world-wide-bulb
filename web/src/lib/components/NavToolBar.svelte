<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { BookOpen, RotateCcwClock, Sun } from '@lucide/svelte';
	import ToolBar from './ui/ToolBar.svelte';
	import ToolBarItem from './ui/ToolBarItem.svelte';

	const links = [
		{ href: resolve('/'), label: 'Bulb', icon: Sun },
		{ href: resolve('/history'), label: 'History', icon: RotateCcwClock },
		{ href: resolve('/docs'), label: 'Docs', icon: BookOpen }
	];

	let activeIndex = $derived.by(() => {
		if (page.url.pathname.startsWith('/docs')) return 2;
		if (page.url.pathname === '/history') return 1;
		return 0;
	});
</script>

<nav
	class="fixed bottom-6 left-1/2 -translate-x-1/2 flex justify-center items-center z-50 pointer-events-auto"
>
	<div class="w-48 sm:w-84">
		<ToolBar {activeIndex} count={links.length}>
			{#each links as link, i (link.href)}
				<ToolBarItem
					href={link.href}
					active={activeIndex === i}
					ariaLabel={link.label}
					title={link.label}
				>
					<link.icon size={18} class="shrink-0" />
					<span class="hidden sm:inline">{link.label}</span>
				</ToolBarItem>
			{/each}
		</ToolBar>
	</div>
</nav>
