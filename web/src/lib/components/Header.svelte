<script lang="ts">
	import { CircleQuestionMark, Webhook } from '@lucide/svelte';
	import ViewerBadge from './ViewerBadge.svelte';
	import { HookDialog, InfoDialog } from '$lib';

	type HeaderPosition = 'top' | 'bottom';

	let { position = 'top' } = $props<{
		position?: HeaderPosition;
	}>();

	let isInfoOpen = $state(false);
	let hookInfoOpen = $state(false);

	// TODO: maybe create new component?
	const buttonStyle = `pointer-events-auto p-2 rounded-full text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-all border border-transparent hover:border-app-border focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer`;
</script>

<header
	class="fixed left-4 right-4 z-40 flex items-center justify-between pointer-events-none {position ===
	'bottom'
		? 'bottom-4'
		: 'top-4'}"
>
	<ViewerBadge />
	<div class="flex items-center gap-1">
		<button
			onclick={() => (hookInfoOpen = true)}
			aria-label="Hook into the lamp"
			title="Hook into the lamp"
			class={buttonStyle}
		>
			<Webhook size={22} />
		</button>
		<button
			onclick={() => (isInfoOpen = true)}
			aria-label="About World Wide Bulb"
			title="About World Wide Bulb"
			class={buttonStyle}
		>
			<CircleQuestionMark size={22} />
		</button>
	</div>
</header>

<InfoDialog bind:isOpen={isInfoOpen} />
<HookDialog bind:isOpen={hookInfoOpen} />
