<script lang="ts">
	import { page } from '$app/state';
	import { RotateCcwClock, Sun } from '@lucide/svelte';
	import ToolBar from './ui/ToolBar.svelte';
	import ToolBarItem from './ui/ToolBarItem.svelte';

	const links = [
		{ href: '/', label: 'Bulb', icon: Sun },
		{ href: '/history', label: 'History', icon: RotateCcwClock }
	];

	let activeIndex = $derived(page.url.pathname === '/history' ? 1 : 0);
</script>

<nav
	class="fixed bottom-6 left-1/2 -translate-x-1/2 flex justify-center items-center z-50 pointer-events-auto"
>
	<div class="w-64 sm:w-72">
		<ToolBar {activeIndex} count={links.length}>
			{#each links as link, i (link.href)}
				<ToolBarItem href={link.href} active={activeIndex === i}>
					<link.icon size={20} class="shrink-0" />
					<span>{link.label}</span>
				</ToolBarItem>
			{/each}
		</ToolBar>
	</div>
</nav>
