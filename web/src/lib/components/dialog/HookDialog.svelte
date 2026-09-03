<script lang="ts">
	import { Check, Copy, Terminal, Webhook } from '@lucide/svelte';
	import Dialog from '$lib/components/ui/Dialog.svelte';

	let { isOpen = $bindable(false) }: { isOpen?: boolean } = $props();

	let copied = $state(false);

	function getSseUrl(): string {
		if (typeof window === 'undefined') return '/api/v1/events';
		return `${window.location.origin}/api/v1/events`;
	}

	const curlExample = $derived(`curl -N ${getSseUrl()}`);

	function copyCommand() {
		if (typeof navigator !== 'undefined') {
			navigator.clipboard.writeText(curlExample);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		}
	}
</script>

<Dialog bind:isOpen maxWidth="max-w-lg">
	{#snippet header()}
		<div class="flex items-center gap-2">
			<h2 class="text-lg font-bold tracking-tight">Hook into the Lamp</h2>
			<Webhook size={20} class="text-app-accent animate-spin" />
		</div>
	{/snippet}

	<p>
		Want to synchronize a physical lamp, Home Assistant automation, or custom script with World Wide
		Bulb?
	</p>
	<p>
		Subscribe to the real-time <span class="font-semibold">Server-Sent Events (SSE)</span> stream to receive
		state changes instantly.
	</p>

	<div class="space-y-1.5 pt-1">
		<span class="text-xs font-semibold text-app-text flex items-center gap-1.5">
			<Terminal size={14} class="text-app-accent" /> Stream Endpoint:
		</span>
		<div
			class="flex items-center justify-between p-2.5 rounded-xl bg-app-surface border border-app-border text-xs font-mono text-app-text"
		>
			<span class="truncate">{curlExample}</span>
			<button
				type="button"
				onclick={copyCommand}
				aria-label="Copy curl command"
				class="p-1 rounded-md text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer shrink-0 ml-2"
			>
				{#if copied}
					<Check size={14} class="text-emerald-400" />
				{:else}
					<Copy size={14} />
				{/if}
			</button>
		</div>
	</div>

	{#snippet footer()}
		<a
			href="/docs"
			onclick={() => (isOpen = false)}
			class="text-xs font-semibold text-app-accent hover:underline inline-flex items-center gap-1"
		>
			View Code Examples &rarr;
		</a>
		<button
			type="button"
			onclick={() => (isOpen = false)}
			class="px-4 py-1.5 rounded-lg text-xs font-semibold cursor-pointer text-app-accent-contrast bg-app-accent hover:opacity-90 transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-app-accent/50 self-end sm:self-center"
		>
			Got it!
		</button>
	{/snippet}
</Dialog>
