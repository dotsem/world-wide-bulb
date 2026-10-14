<script lang="ts">
	import { Check, Copy, ExternalLink } from '@lucide/svelte';
	import RestTypeLabel from './RestTypeLabel.svelte';

	type Param = {
		name: string;
		type: string;
		required: boolean;
		description: string;
	};

	type Language = 'curl' | 'javascript' | 'python' | 'go';

	let {
		id,
		method,
		path,
		title,
		description,
		isInternal = false,
		params = [],
		responseHtml,
		codeExamples = {}
	}: {
		id: string;
		method: 'GET' | 'POST' | 'SSE' | 'WS';
		path: string;
		title: string;
		description: string;
		isInternal?: boolean;
		params?: Param[];
		responseHtml?: string;
		codeExamples?: Partial<
			Record<Language, { highlighted: string; raw: string; repoPath?: string }>
		>;
	} = $props();

	const availableLangs = $derived(
		(['curl', 'javascript', 'python', 'go'] as const).filter((l) => codeExamples && codeExamples[l])
	);

	let activeLang = $state<Language>('curl');

	$effect(() => {
		if (availableLangs.length > 0 && !availableLangs.includes(activeLang)) {
			activeLang = availableLangs[0];
		}
	});

	let copied = $state(false);

	function copyActiveCode() {
		const snippet = codeExamples?.[activeLang]?.raw;
		if (snippet && typeof navigator !== 'undefined') {
			navigator.clipboard.writeText(snippet);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		}
	}

	const repoBase = 'https://github.com/dotsem/world-wide-bulb/tree/main/';
</script>

<article
	{id}
	class="scroll-mt-28 rounded-2xl border border-app-border bg-app-surface/60 p-5 sm:p-7 space-y-6 shadow-sm"
>
	<!-- Header -->
	<div class="space-y-3">
		<div class="flex flex-wrap items-center gap-2.5">
			<RestTypeLabel {method} {isInternal} />
			<code
				class="font-mono text-sm px-2.5 py-0.5 rounded-lg bg-app-surface border border-app-border text-app-text font-bold"
			>
				{path}
			</code>
		</div>
		<h3 class="text-xl font-extrabold text-app-text tracking-tight">{title}</h3>
		<p class="text-xs sm:text-sm text-app-muted leading-relaxed">{description}</p>
		{#if isInternal}
			<div
				class="text-[11px] px-3 py-2 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-300 leading-normal"
			>
				<strong>Self-hosted / Web only:</strong> External browsers will be rejected in production by CORS
				unless their origin is whitelisted in server configuration.
			</div>
		{/if}
	</div>

	<!-- Query / Body Parameters -->
	{#if params.length > 0}
		<div class="space-y-2">
			<h4 class="text-xs font-bold uppercase tracking-wider text-app-muted">Parameters</h4>
			<div class="overflow-x-auto rounded-xl border border-app-border bg-app-surface-solid">
				<table class="w-full text-left text-xs font-mono">
					<thead>
						<tr class="border-b border-app-border bg-app-surface/50 text-app-muted">
							<th class="px-3.5 py-2">Name</th>
							<th class="px-3.5 py-2">Type</th>
							<th class="px-3.5 py-2">Required</th>
							<th class="px-3.5 py-2 font-sans">Description</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-app-border/40 text-app-text">
						{#each params as p (p.name)}
							<tr>
								<td class="px-3.5 py-2 text-app-accent font-semibold">{p.name}</td>
								<td class="px-3.5 py-2 text-app-muted">{p.type}</td>
								<td class="px-3.5 py-2">{p.required ? 'Yes' : 'No'}</td>
								<td class="px-3.5 py-2 font-sans text-app-muted">{p.description}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}

	<!-- Code Examples & Return Value Grid -->
	<div class="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
		<!-- Response Schema -->
		{#if responseHtml}
			<div class="lg:col-span-5 space-y-2">
				<h4 class="text-xs font-bold uppercase tracking-wider text-app-muted">Response Payload</h4>
				<div
					class="rounded-xl border border-app-border bg-app-surface-solid overflow-hidden text-xs font-mono p-3 leading-relaxed custom-scrollbar max-h-80"
				>
					<!-- eslint-disable-next-line svelte/no-at-html-tags -->
					{@html responseHtml}
				</div>
			</div>
		{/if}

		<!-- Code Snippets -->
		{#if availableLangs.length > 0}
			<div class="{responseHtml ? 'lg:col-span-7' : 'lg:col-span-12'} space-y-2">
				<div class="flex items-center justify-between gap-2">
					<h4 class="text-xs font-bold uppercase tracking-wider text-app-muted">Code Example</h4>
					<div class="inline-flex p-0.5 rounded-lg bg-app-surface border border-app-border">
						{#each availableLangs as lang (lang)}
							<button
								type="button"
								onclick={() => (activeLang = lang)}
								class="px-2.5 py-1 rounded-md text-[11px] font-semibold capitalize transition-all cursor-pointer {activeLang ===
								lang
									? 'bg-app-accent text-app-accent-contrast shadow-xs'
									: 'text-app-muted hover:text-app-text'}"
							>
								{lang === 'javascript' ? 'JS' : lang === 'curl' ? 'cURL' : lang}
							</button>
						{/each}
					</div>
				</div>

				<div
					class="rounded-xl border border-app-border bg-app-surface-solid overflow-hidden shadow-sm"
				>
					<div
						class="flex items-center justify-between px-3 py-1.5 bg-app-surface border-b border-app-border text-[11px] text-app-muted"
					>
						<span class="font-mono">
							{#if codeExamples[activeLang]?.repoPath}
								<a
									target="_blank"
									rel="noreferrer"
									href={repoBase + codeExamples[activeLang]?.repoPath}
									class="inline-flex items-center gap-1 text-app-accent hover:underline"
								>
									<span>{codeExamples[activeLang]?.repoPath}</span>
									<ExternalLink size={12} />
								</a>
							{:else}
								terminal
							{/if}
						</span>
						<button
							type="button"
							onclick={copyActiveCode}
							class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer"
						>
							{#if copied}
								<Check size={12} class="text-emerald-400" />
								<span class="text-emerald-400">Copied</span>
							{:else}
								<Copy size={12} />
								<span>Copy</span>
							{/if}
						</button>
					</div>

					<div
						class="p-3 text-xs font-mono leading-relaxed custom-scrollbar max-h-80 overflow-auto"
					>
						<!-- eslint-disable-next-line svelte/no-at-html-tags -->
						{@html codeExamples[activeLang]?.highlighted || ''}
					</div>
				</div>
			</div>
		{/if}
	</div>
</article>
