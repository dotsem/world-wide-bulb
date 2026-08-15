<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import NavToolBar from '$lib/components/NavToolBar.svelte';
	import { bulbState } from '$lib';

	let { children } = $props();

	let theme = $derived(bulbState.isOn ? 'on' : 'off');

	$effect(() => {
		if (typeof document !== 'undefined') {
			document.documentElement.setAttribute('data-theme', theme);
		}
	});
</script>

<svelte:head>
	<title>World Wide Bulb</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	data-theme={theme}
	class="min-h-screen bg-app-bg text-app-text flex flex-col font-sans selection:bg-amber-500/30 selection:text-amber-200 transition-colors duration-300 ease-out"
>
	<NavToolBar />

	<main class="flex-1 flex flex-col">
		{@render children()}
	</main>
</div>
