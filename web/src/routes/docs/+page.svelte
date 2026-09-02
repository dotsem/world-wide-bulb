<script lang="ts">
	import { ArrowLeft, Webhook } from '@lucide/svelte';
	import { bulbState, Bulb } from '$lib';
	import CodeExamples from './CodeExamples.svelte';
	import EventPayloads from './EventPayloads.svelte';

	function getBaseUrl(): string {
		if (typeof window === 'undefined') return 'http://localhost:8080';
		return window.location.origin;
	}

	const sseUrl = $derived(`${getBaseUrl()}/api/v1/events`);
</script>

<svelte:head>
	<title>API & Stream Docs - World Wide Bulb</title>
	<meta
		name="description"
		content="Developer documentation and code examples (JavaScript, Python, Go, cURL) for the World Wide Bulb Server-Sent Events stream."
	/>
</svelte:head>

<div
	class="sticky top-0 z-40 px-4 py-4 sm:px-8 bg-app-surface-solid mb-8 flex items-center justify-between border-b border-app-border overflow-hidden"
>
	<div class="flex items-center gap-3.5">
		<a
			href="/"
			class="p-2 rounded-xl text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors border border-app-border/40 shrink-0"
			aria-label="Back to bulb"
			title="Back to bulb"
		>
			<ArrowLeft size={20} />
		</a>
		<div>
			<h1
				class="text-2xl sm:text-3xl font-extrabold tracking-tight text-app-text flex items-center gap-3"
			>
				Documentation
			</h1>
			<p class="text-sm text-app-muted mt-0.5">Fine things take time 🍷, just like these docs...</p>
		</div>
	</div>
</div>

<div class="fixed -top-24 -right-24 z-40 pointer-events-none shrink-0">
	<Bulb
		isOn={bulbState.isOn}
		disabled={true}
		durationMs={300}
		class="h-72 w-72 rotate-215 transition-all"
	/>
</div>

<div class="max-w-4xl mx-auto w-full px-4 sm:px-8 pb-32 space-y-8">
	<!-- Overview Card -->
	<section
		class="rounded-2xl border border-app-border bg-app-surface/60 p-6 sm:p-7 space-y-4 backdrop-blur-sm"
	>
		<div class="flex items-center gap-2.5 text-app-accent">
			<Webhook size={22} />
			<h2 class="text-lg font-bold text-app-text">Server-Sent Events (SSE) Endpoint</h2>
		</div>
		<p class="text-sm text-app-muted leading-relaxed">
			World Wide Bulb exposes a public, unauthenticated HTTP event stream. Anyone can subscribe to
			receive instantaneous notifications whenever the lamp flips state or a toggle reason is
			submitted.
		</p>
		<div class="flex flex-wrap items-center gap-3 pt-2">
			<span
				class="px-2.5 py-1 rounded-md bg-app-accent/15 text-app-accent text-xs font-semibold uppercase tracking-wider"
				>GET</span
			>
			<code
				class="font-mono text-sm px-3 py-1.5 rounded-lg bg-app-surface border border-app-border text-app-text"
			>
				/api/v1/events
			</code>
			<span class="text-xs text-app-muted/80">Content-Type: text/event-stream</span>
		</div>
	</section>

	<CodeExamples {sseUrl} />

	<EventPayloads />
</div>
