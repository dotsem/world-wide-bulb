<script lang="ts">
	import { Package } from '@lucide/svelte';

	let {
		payloads
	}: {
		payloads?: {
			stateChanged?: string;
			reasonUpdated?: string;
			ping?: string;
		};
	} = $props();
</script>

{#snippet eventPayload(
	type: string,
	highlightedHtml: string | undefined,
	sentWhen: string,
	description?: string
)}
	<div class="p-4 rounded-xl bg-app-surface border border-app-border space-y-2">
		<p class="font-mono font-bold text-app-accent">event: {type}</p>
		<p class="text-xs text-app-muted">{sentWhen}</p>
		<div
			class="rounded-lg bg-app-surface-solid border border-app-border/60 text-xs font-mono text-app-text overflow-hidden"
		>
			<!-- eslint-disable-next-line svelte/no-at-html-tags -->
			{@html highlightedHtml}
		</div>
		{#if description}
			<p class="text-xs text-app-muted">{description}</p>
		{/if}
	</div>
{/snippet}

<section class="rounded-2xl border border-app-border bg-app-surface/60 p-6 sm:p-7 space-y-4">
	<div class="flex items-center gap-2">
		<Package size={20} class="text-app-accent" />
		<h2 class="text-lg font-bold text-app-text">Event Types & Payloads</h2>
	</div>

	{@render eventPayload(
		'state_changed',
		payloads?.stateChanged,
		'Sent on connect and whenever lamp is toggled'
	)}

	{@render eventPayload(
		'reason_updated',
		payloads?.reasonUpdated,
		'Sent when someone attaches a reason'
	)}

	{@render eventPayload(
		'ping',
		payloads?.ping,
		'Sent every 15s keep-alive heartbeat',
		'Empty payload data to prevent reverse-proxy and router timeouts.'
	)}
</section>
