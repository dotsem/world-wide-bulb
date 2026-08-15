<script lang="ts">
	import { page } from '$app/state';
	import ToolBar from './ui/ToolBar.svelte';
	import ToolBarItem from './ui/ToolBarItem.svelte';

	const links = [
		{ href: '/', label: 'Bulb', icon: 'bulb' },
		{ href: '/history', label: 'History', icon: 'history' }
	];

	let activeIndex = $derived(page.url.pathname === '/history' ? 1 : 0);
</script>

<nav class="fixed bottom-4 left-4 right-4 flex justify-center items-center z-50">
	<div class="w-64 sm:w-72">
		<ToolBar {activeIndex} count={links.length}>
			{#each links as link, i (link.href)}
				<ToolBarItem href={link.href} active={activeIndex === i}>
					{#if link.icon === 'bulb'}
						<!-- //TODO: add icons from lucide package-->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-4 w-4 shrink-0"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M12 2v2m0 16v2m8-10h2M2 12H4m14.364-6.364l-1.414 1.414M6.05 17.95l-1.414 1.414M17.95 17.95l1.414 1.414M6.05 6.05L4.636 4.636M12 7a5 5 0 00-3 9h6a5 5 0 00-3-9z"
							/>
						</svg>
					{:else if link.icon === 'history'}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-4 w-4 shrink-0"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					{/if}
					<span>{link.label}</span>
				</ToolBarItem>
			{/each}
		</ToolBar>
	</div>
</nav>
