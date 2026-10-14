<script lang="ts">
	import { Menu, X } from '@lucide/svelte';
	import { resolve } from '$app/paths';
	import DocsBarSection from './DocsBarSection.svelte';
	import DocsBarEndpoint from './DocsBarEndpoint.svelte';
	import RestTypeLabel from './RestTypeLabel.svelte';

	let {
		activeId = ''
	}: {
		activeId?: string;
	} = $props();

	let mobileOpen = $state(false);

	function closeMobile() {
		mobileOpen = false;
	}
</script>

<!-- Mobile Navigation Toggle Bar -->
<div
	class="lg:hidden mb-6 flex items-center justify-between p-3 rounded-2xl bg-app-surface border border-app-border"
>
	<span class="text-xs font-semibold text-app-muted uppercase tracking-wider">API Navigation</span>
	<button
		type="button"
		onclick={() => (mobileOpen = !mobileOpen)}
		class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-medium text-app-text bg-app-surface-hover border border-app-border transition-colors cursor-pointer"
		aria-label="Toggle navigation drawer"
	>
		{#if mobileOpen}
			<X size={16} />
			<span>Close</span>
		{:else}
			<Menu size={16} />
			<span>Menu</span>
		{/if}
	</button>
</div>

<!-- Navigation Sidebar Container -->
<aside
	class="space-y-6 max-h-[calc(100vh-8rem)] overflow-y-auto custom-scrollbar p-1 {mobileOpen
		? 'block'
		: 'hidden lg:block'}"
>
	<DocsBarSection title="REST API" href={resolve('/docs#rest-api')}>
		<DocsBarEndpoint
			href={resolve('/docs#get-state')}
			title="/state"
			active={activeId === 'get-state'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="GET" />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#get-history')}
			title="/history"
			active={activeId === 'get-history'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="GET" />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#post-toggle')}
			title="/toggle"
			active={activeId === 'post-toggle'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="POST" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#post-reason')}
			title="/reason"
			active={activeId === 'post-reason'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="POST" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>

	<DocsBarSection title="Server Side Events" href={resolve('/docs#sse')}>
		<DocsBarEndpoint
			href={resolve('/docs#sse-events')}
			title="/events"
			active={activeId === 'sse-events'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="SSE" />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>

	<DocsBarSection title="WebSockets" href={resolve('/docs#ws')}>
		<DocsBarEndpoint
			href={resolve('/docs#ws-stream')}
			title="/ws"
			active={activeId === 'ws-stream'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="WS" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>
</aside>
