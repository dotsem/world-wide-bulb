<script lang="ts">
	import { Menu, X } from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { resolve } from '$app/paths';
	import type { ResolvedPathname } from '$app/types';
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

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && mobileOpen) {
			closeMobile();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#snippet navLinks()}
	<DocsBarSection
		title="Overview"
		href={resolve('/docs#overview') as ResolvedPathname}
		onclick={closeMobile}
	>
		<DocsBarEndpoint
			href={resolve('/docs#overview') as ResolvedPathname}
			title="About &amp; GitHub"
			active={activeId === 'overview'}
			onclick={closeMobile}
		/>
		<DocsBarEndpoint
			href={resolve('/docs#use-cases') as ResolvedPathname}
			title="Use Cases"
			active={activeId === 'use-cases'}
			onclick={closeMobile}
		/>
		<DocsBarEndpoint
			href={resolve('/docs#self-hosting') as ResolvedPathname}
			title="Self-Hosting"
			active={activeId === 'self-hosting'}
			onclick={closeMobile}
		/>
	</DocsBarSection>

	<DocsBarSection
		title="REST API"
		href={resolve('/docs#rest-api') as ResolvedPathname}
		onclick={closeMobile}
	>
		<DocsBarEndpoint
			href={resolve('/docs#get-state') as ResolvedPathname}
			title="/state"
			active={activeId === 'get-state'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="GET" />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#get-history') as ResolvedPathname}
			title="/history"
			active={activeId === 'get-history'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="GET" />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#post-toggle') as ResolvedPathname}
			title="/toggle"
			active={activeId === 'post-toggle'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="POST" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>

		<DocsBarEndpoint
			href={resolve('/docs#post-reason') as ResolvedPathname}
			title="/reason"
			active={activeId === 'post-reason'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="POST" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>

	<DocsBarSection
		title="Server-Sent Events"
		href={resolve('/docs#sse') as ResolvedPathname}
		onclick={closeMobile}
	>
		<DocsBarEndpoint
			href={resolve('/docs#sse-events') as ResolvedPathname}
			title="/events"
			active={activeId === 'sse-events'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="SSE" />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>

	<DocsBarSection
		title="WebSockets"
		href={resolve('/docs#ws') as ResolvedPathname}
		onclick={closeMobile}
	>
		<DocsBarEndpoint
			href={resolve('/docs#ws-stream') as ResolvedPathname}
			title="/ws"
			active={activeId === 'ws-stream'}
			onclick={closeMobile}
		>
			{#snippet badge()}
				<RestTypeLabel method="WS" isInternal={true} />
			{/snippet}
		</DocsBarEndpoint>
	</DocsBarSection>
{/snippet}

<!-- Mobile Floating Trigger Button -->
<div class="fixed bottom-6 right-4 sm:right-6 z-40 lg:hidden pointer-events-auto">
	<button
		type="button"
		onclick={() => (mobileOpen = true)}
		class="h-12 w-12 rounded-full flex items-center justify-center bg-app-surface-solid border border-app-border text-app-text shadow-xl shadow-black/40 hover:text-app-accent hover:border-app-accent/40 transition-all cursor-pointer"
		aria-label="Open documentation navigation menu"
		title="Documentation Navigation"
	>
		<Menu size={20} />
	</button>
</div>

<!-- Mobile Slide-Over Drawer with Backdrop -->
{#if mobileOpen}
	<div
		role="presentation"
		transition:fade={{ duration: 200 }}
		class="fixed inset-0 bg-black/60 backdrop-blur-xs z-50 lg:hidden"
		onclick={closeMobile}
	></div>

	<div
		role="dialog"
		aria-modal="true"
		aria-label="Documentation navigation drawer"
		transition:fly={{ x: 320, duration: 250 }}
		class="fixed inset-y-0 right-0 w-72 sm:w-80 max-w-[85vw] bg-app-surface-solid border-l border-app-border z-50 p-6 flex flex-col shadow-2xl lg:hidden pointer-events-auto"
	>
		<div class="flex items-center justify-between pb-4 mb-4 border-b border-app-border">
			<span class="text-sm font-bold text-app-text">Documentation</span>
			<button
				type="button"
				onclick={closeMobile}
				class="p-1.5 rounded-lg text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer"
				aria-label="Close navigation drawer"
			>
				<X size={18} />
			</button>
		</div>

		<div class="overflow-y-auto custom-scrollbar flex-1 space-y-6 pr-1">
			{@render navLinks()}
		</div>
	</div>
{/if}

<!-- Desktop Sticky Sidebar -->
<aside
	class="hidden lg:block space-y-6 max-h-[calc(100vh-8rem)] overflow-y-auto custom-scrollbar p-1"
>
	{@render navLinks()}
</aside>
