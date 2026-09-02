<script lang="ts">
	import { ArrowLeft, Webhook } from '@lucide/svelte';
	import CodeExamples from './CodeExamples.svelte';
	import EventPayloads from './EventPayloads.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
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

<div class="max-w-7xl mx-auto w-full px-4 sm:px-8 pb-32">
	<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 items-start">
		<div class="lg:col-span-5 space-y-6">
			<section class="rounded-2xl border border-app-border bg-app-surface/60 p-5 sm:p-6 space-y-4">
				<div class="flex items-center gap-2.5 text-app-accent">
					<Webhook size={20} />
					<h2 class="text-base font-bold text-app-text">Server-Sent Events (SSE) Endpoint</h2>
				</div>
				<p class="text-xs sm:text-sm text-app-muted leading-relaxed">
					World Wide Bulb exposes a public, unauthenticated HTTP event stream. Anyone can subscribe
					to receive instantaneous notifications whenever the lamp flips state or a toggle reason is
					submitted. Go creative!
				</p>
				<div class="flex flex-wrap items-center gap-2.5 pt-1">
					<span
						class="px-2 py-0.5 rounded-md bg-app-accent/15 text-app-accent text-[11px] font-semibold uppercase tracking-wider"
						>GET</span
					>
					<code
						class="font-mono text-xs px-2.5 py-1 rounded-lg bg-app-surface border border-app-border text-app-text"
					>
						/api/v1/events
					</code>
					<span class="text-[11px] text-app-muted/80">Content-Type: text/event-stream</span>
				</div>
			</section>

			<EventPayloads payloads={data?.payloads} />
		</div>

		<div class="lg:col-span-7 lg:sticky lg:top-24">
			<CodeExamples highlighted={data?.highlighted} raw={data?.raw} />
		</div>
	</div>
</div>
