<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import NavToolBar from '$lib/components/NavToolBar.svelte';
	import InfoDialog from '$lib/components/InfoDialog.svelte';
	import { bulbState } from '$lib';
	import { CircleQuestionMark } from '@lucide/svelte';

	let { children } = $props();
	let isInfoOpen = $state(false);

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
	class="min-h-screen bg-app-bg text-app-text flex flex-col font-sans selection:bg-app-accent/30 selection:text-app-accent transition-colors duration-300 ease-out relative"
>
	<header class="fixed top-4 right-4 z-40">
		<button
			type="button"
			onclick={() => (isInfoOpen = true)}
			aria-label="About World Wide Bulb"
			title="About World Wide Bulb"
			class="p-2 rounded-full text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-all border border-transparent hover:border-app-border focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer"
		>
			<CircleQuestionMark size={22} />
		</button>
	</header>

	<NavToolBar />

	<main class="flex-1 flex flex-col">
		{@render children()}
	</main>

	<InfoDialog bind:isOpen={isInfoOpen} />
</div>
