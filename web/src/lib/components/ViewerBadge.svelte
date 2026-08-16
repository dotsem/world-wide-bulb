<script lang="ts">
	import { Eye, EyeClosed } from '@lucide/svelte';
	import { bulbState } from '$lib';

	const labels = [
		{ singular: 'viewer', plural: 'viewers' },
		{ singular: 'onlooker', plural: 'onlookers' },
		{ singular: 'spectator', plural: 'spectators' },
		{ singular: 'stalker', plural: 'stalkers' },
		{ singular: 'moth', plural: 'moths' },
		{ singular: 'curious soul', plural: 'curious souls' },
		{ singular: 'alien', plural: 'aliens' },
		{ singular: 'person', plural: 'people' }
		// TODO: populate more
	];

	let labelIndex = $state(0);
	let isBlinking = $state(false);

	let currentLabel = $derived(
		bulbState.viewers === 1 ? labels[labelIndex].singular : labels[labelIndex].plural
	);

	$effect(() => {
		const interval = setInterval(() => {
			isBlinking = true;
			setTimeout(() => {
				labelIndex = (labelIndex + 1) % labels.length;
				isBlinking = false;
			}, 180);
		}, 4000);

		return () => clearInterval(interval);
	});
</script>

<div
	title="{bulbState.viewers} {bulbState.viewers === 1
		? 'person'
		: 'people'} are currently watching this lamp somewhere in this world, god knows what they use it for..."
	class="pointer-events-auto flex items-center gap-2 px-3 py-1.5 rounded-full bg-app-surface-solid/80 backdrop-blur-md border border-app-border text-xs font-medium text-app-text shadow-sm transition-all"
>
	{#if isBlinking}
		<EyeClosed size={14} class="text-app-muted shrink-0" />
	{:else}
		<Eye size={14} class="text-app-text shrink-0" />
	{/if}
	<span class="tabular-nums transition-all duration-200">
		{bulbState.viewers}
		{currentLabel}
	</span>
</div>
