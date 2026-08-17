<script lang="ts">
	import './layout.css';
	import bulb_on from '$lib/assets/bulb_on.svg';
	import bulb_off from '$lib/assets/bulb_off.svg';
	import NavToolBar from '$lib/components/NavToolBar.svelte';
	import { bulbState } from '$lib';

	let { children } = $props();

	let theme = $derived(bulbState.isOn ? 'on' : 'off');

	$effect(() => {
		if (typeof document !== 'undefined') {
			document.documentElement.setAttribute('data-theme', theme);
		}
	});

	let favicon = $derived(bulbState.isOn ? bulb_on : bulb_off);
</script>

<svelte:head>
	<title>World Wide Bulb</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div
	data-theme={theme}
	class="min-h-screen bg-app-bg text-app-text flex flex-col font-sans selection:bg-app-accent/30 selection:text-app-accent relative"
>
	<NavToolBar />

	<main class="flex-1 flex flex-col">
		{@render children()}
	</main>
</div>
