<script lang="ts">
	import { Package } from '@lucide/svelte';

	const stateChangedEventPayload = JSON.stringify(
		{
			type: 'state_changed',
			id: 42,
			state: true,
			reason: 'Midnight coding',
			created_at: '2026-09-02T17:35:00Z'
		},
		null,
		2
	);

	const reasonUpdatedEventPayload = JSON.stringify(
		{
			type: 'reason_updated',
			id: 42,
			reason: 'Party at the office!'
		},
		null,
		2
	);
</script>

{#snippet eventPayload(type: string, data: string, sentWhen: string, description?: string)}
	<div class="p-4 rounded-xl bg-app-surface border border-app-border space-y-2">
		<p class="font-mono font-bold text-app-accent">event: {type}</p>
		<p class="text-xs text-app-muted">{sentWhen}</p>
		<pre
			class="p-3 rounded-lg bg-app-surface-solid border border-app-border/60 text-xs font-mono text-app-text overflow-x-auto"><code
				>{data}</code
			></pre>
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
		stateChangedEventPayload,
		'Sent on connect and whenever lamp is toggled'
	)}

	{@render eventPayload(
		'reason_updated',
		reasonUpdatedEventPayload,
		'Sent when someone attaches a reason'
	)}

	{@render eventPayload(
		'ping',
		'""',
		'Sent every 15s keep-alive heartbeat',
		'Empty payload data to prevent reverse-proxy and router timeouts.'
	)}
</section>
