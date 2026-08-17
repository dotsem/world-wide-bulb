<script lang="ts">
	import { CircleQuestionMark } from '@lucide/svelte';
	import ViewerBadge from './ViewerBadge.svelte';
	import { InfoDialog } from '$lib';

	type HeaderPosition = 'top' | 'bottom';

	let { position = 'top' } = $props<{
		position?: HeaderPosition;
	}>();

	let isInfoOpen = $state(false);
</script>

<header
	class="fixed left-4 right-4 z-40 flex items-center justify-between pointer-events-none {position ===
	'bottom'
		? 'bottom-4'
		: 'top-4'}"
>
	<ViewerBadge />
	<button
		type="button"
		onclick={() => (isInfoOpen = true)}
		aria-label="About World Wide Bulb"
		title="About World Wide Bulb"
		class="pointer-events-auto p-2 rounded-full text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-all border border-transparent hover:border-app-border focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer"
	>
		<CircleQuestionMark size={22} />
	</button>
</header>

<InfoDialog bind:isOpen={isInfoOpen} />
