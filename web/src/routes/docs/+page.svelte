<script lang="ts">
	import { ArrowLeft, BookOpen, Globe, Radio, Webhook } from '@lucide/svelte';
	import { resolve } from '$app/paths';
	import DocsBar from './components/DocsBar.svelte';
	import DocsHero from './components/DocsHero.svelte';
	import DocsUseCases from './components/DocsUseCases.svelte';
	import DocsSelfHosting from './components/DocsSelfHosting.svelte';
	import DocsSection from './components/DocsSection.svelte';
	import EndpointCard from './components/EndpointCard.svelte';
	import EventPayloads from './EventPayloads.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let activeSection = $state('overview');

	const sectionDelegates: Record<string, string> = {
		'rest-api': 'get-state',
		sse: 'sse-events',
		ws: 'ws-stream'
	};

	$effect(() => {
		if (typeof window === 'undefined') return;

		function updateFromHash() {
			const hash = window.location.hash.replace('#', '');
			if (hash) {
				activeSection = sectionDelegates[hash] ?? hash;
			}
		}

		window.addEventListener('hashchange', updateFromHash);
		updateFromHash();

		const observer = new IntersectionObserver(
			(entries) => {
				const visibleArticle = entries.find(
					(e) => e.isIntersecting && e.target.tagName === 'ARTICLE'
				);
				const visibleSection = entries.find(
					(e) => e.isIntersecting && e.target.tagName === 'SECTION'
				);
				const visible = visibleArticle || visibleSection;

				if (visible?.target?.id) {
					const targetId = visible.target.id;
					activeSection = sectionDelegates[targetId] ?? targetId;
				}
			},
			{ rootMargin: '-20% 0px -70% 0px' }
		);

		const elements = document.querySelectorAll('article[id], section[id]');
		elements.forEach((el) => observer.observe(el));

		return () => {
			window.removeEventListener('hashchange', updateFromHash);
			observer.disconnect();
		};
	});

	const ep = $derived(data.endpoints);
</script>

<svelte:head>
	<title>API & Stream Docs - World Wide Bulb</title>
	<meta
		name="description"
		content="Developer documentation and code examples (JavaScript, Python, Go, cURL) for World Wide Bulb REST and Realtime endpoints."
	/>
</svelte:head>

<!-- Sticky Top Bar -->
<div
	class="sticky top-0 z-40 px-4 py-4 sm:px-8 bg-app-surface-solid mb-8 flex items-center justify-between border-b border-app-border"
>
	<div class="flex items-center gap-3.5">
		<a
			href={resolve('/')}
			class="p-2 rounded-xl text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors border border-app-border/40 shrink-0"
			aria-label="Back to bulb"
			title="Back to bulb"
		>
			<ArrowLeft size={20} />
		</a>
		<div class="flex gap-4 items-center">
			<h1
				class="text-xl sm:text-2xl font-extrabold tracking-tight text-app-text flex items-center gap-2.5"
			>
				<BookOpen size={22} class="text-app-accent" />
				<span>Documentation</span>
			</h1>
			<p class="text-xs text-app-muted mt-0.5">Yes we have them!</p>
		</div>
	</div>
</div>

<!-- Main Documentation Layout -->
<div class="w-full pb-32">
	<div class="docs-layout-flow px-4 sm:px-8">
		<!-- Sidebar Navigation (Desktop sticky rail + Mobile Drawer) -->
		<div class="lg:sticky lg:top-24 lg:self-start lg:w-56 lg:shrink-0">
			<DocsBar activeId={activeSection} />
		</div>

		<div class="min-w-0 max-w-4xl w-full space-y-16 mx-auto lg:mx-0">
			<DocsHero />

			<DocsUseCases />

			<DocsSelfHosting />

			<DocsSection
				id="rest-api"
				title="REST API"
				description="Standard HTTP endpoints for querying lamp status, viewing toggle history, and performing mutations."
				icon={Globe}
			>
				<EndpointCard
					id="get-state"
					method="GET"
					path="/api/v1/state"
					title="Get Current State"
					description="Returns whether the bulb is currently ON or OFF, remaining cooldown for your device, and active connected viewers."
					responseHtml={ep.state.responseHtml}
					codeExamples={ep.state.examples}
				/>

				<EndpointCard
					id="get-history"
					method="GET"
					path="/api/v1/history"
					title="Get Toggle History"
					description="Fetch recent state changes in reverse chronological order with cursor-based pagination."
					params={[
						{
							name: 'limit',
							type: 'integer',
							required: false,
							description: 'Number of entries to return (default 20, max 100)'
						},
						{
							name: 'before',
							type: 'integer',
							required: false,
							description: 'Cursor toggle ID for pagination (fetches items before this ID)'
						}
					]}
					responseHtml={ep.history.responseHtml}
					codeExamples={ep.history.examples}
				/>

				<EndpointCard
					id="post-toggle"
					method="POST"
					path="/api/v1/toggle"
					title="Toggle Bulb State"
					description="Flips the lamp on or off. Rate-limited per device and debounced globally."
					isInternal={true}
					responseHtml={ep.toggle.responseHtml}
					codeExamples={ep.toggle.examples}
				/>

				<EndpointCard
					id="post-reason"
					method="POST"
					path="/api/v1/reason"
					title="Submit Toggle Reason"
					description="Attaches an optional text reason to your latest toggle within the allowed time window."
					isInternal={true}
					params={[
						{
							name: 'id',
							type: 'string (uuid)',
							required: true,
							description: 'The unique UUID returned by /api/v1/toggle'
						},
						{
							name: 'reason',
							type: 'string',
							required: true,
							description: 'Reason text (maximum 60 characters, no URLs or profanity)'
						}
					]}
					responseHtml={ep.reason.responseHtml}
					codeExamples={ep.reason.examples}
				/>
			</DocsSection>

			<DocsSection
				id="sse"
				title="Server-Sent Events"
				description="Public, unauthenticated HTTP stream pushing instant updates whenever the lamp toggles or a reason is attached."
				icon={Webhook}
			>
				<EndpointCard
					id="sse-events"
					method="SSE"
					path="/api/v1/events"
					title="Server-Sent Events Stream"
					description="Subscribe to real-time events over standard HTTP."
					responseHtml={ep.events.responseHtml}
					codeExamples={ep.events.examples}
				/>

				<EventPayloads payloads={ep.events.payloads} />
			</DocsSection>

			<DocsSection
				id="ws"
				title="WebSockets"
				description="Bidirectional WebSocket channel used by the client for state synchronization and heartbeats."
				icon={Radio}
			>
				<EndpointCard
					id="ws-stream"
					method="WS"
					path="/ws"
					title="WebSocket Client Channel"
					description="Internal connection for real-time bulb state events."
					isInternal={true}
					responseHtml={ep.ws.responseHtml}
					codeExamples={ep.ws.examples}
				/>
			</DocsSection>
		</div>
	</div>
</div>

<style>
	@media (min-width: 1024px) {
		.docs-layout-flow {
			display: flex;
			align-items: flex-start;
			gap: 2.5rem;
			margin-left: max(2rem, calc(50% - 44.5rem));
			padding-right: 2rem;
		}
	}
</style>
