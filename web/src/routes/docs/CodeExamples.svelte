<script lang="ts">
	import { Check, Code, Copy } from '@lucide/svelte';

	type Language = 'javascript' | 'python' | 'go' | 'curl';

	let {
		highlighted,
		raw
	}: {
		highlighted: Record<Language, string>;
		raw: Record<Language, string>;
	} = $props();

	let activeLang = $state<Language>('javascript');
	let copiedLang = $state<string | null>(null);

	function copySnippet(lang: Language) {
		if (typeof navigator !== 'undefined') {
			navigator.clipboard.writeText(raw[lang]);
			copiedLang = lang;
			setTimeout(() => {
				if (copiedLang === lang) copiedLang = null;
			}, 2000);
		}
	}

	const langToPath: Record<Exclude<Language, 'curl'>, string> = {
		javascript: 'examples/sse/javascript/client.js',
		python: 'examples/sse/python/client.py',
		go: 'examples/sse/go/main.go'
	};

	const repoPath = 'https://github.com/dotsem/world-wide-bulb/tree/main/';
</script>

<section class="space-y-4">
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<Code size={20} class="text-app-accent" />
			<h2 class="text-lg font-bold text-app-text">Code Examples</h2>
		</div>

		<!-- Language Switcher Tabs -->
		<div class="inline-flex p-1 rounded-xl bg-app-surface border border-app-border">
			{#each ['javascript', 'python', 'go', 'curl'] as const as lang (lang)}
				<button
					type="button"
					onclick={() => (activeLang = lang)}
					class="px-3 py-1.5 rounded-lg text-xs font-semibold capitalize transition-all cursor-pointer {activeLang ===
					lang
						? 'bg-app-accent text-app-accent-contrast shadow-sm'
						: 'text-app-muted hover:text-app-text'}"
				>
					{lang === 'javascript' ? 'JavaScript' : lang === 'curl' ? 'cURL' : lang}
				</button>
			{/each}
		</div>
	</div>

	<!-- Code Snippet Display Card -->
	<div
		class="relative rounded-2xl border border-app-border bg-app-surface-solid overflow-hidden shadow-xl"
	>
		<div
			class="flex items-center justify-between px-4 py-2.5 bg-app-surface border-b border-app-border text-xs text-app-muted"
		>
			<span class="font-mono">
				{#if activeLang === 'curl'}
					terminal
				{:else}
					<a
						target="_blank"
						rel="noreferrer"
						href={repoPath + langToPath[activeLang]}
						class="text-app-accent underline">{langToPath[activeLang]}</a
					>
				{/if}
			</span>
			<button
				type="button"
				onclick={() => copySnippet(activeLang)}
				class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer"
			>
				{#if copiedLang === activeLang}
					<Check size={14} class="text-emerald-400" />
					<span class="text-emerald-400">Copied</span>
				{:else}
					<Copy size={14} />
					<span>Copy Code</span>
				{/if}
			</button>
		</div>

		<div class="max-h-145 overflow-y-auto text-xs sm:text-sm font-mono leading-relaxed">
			{#if highlighted && highlighted[activeLang]}
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html highlighted[activeLang]}
			{:else}
				<pre class="p-4 sm:p-5 text-app-text overflow-x-auto"><code>{raw[activeLang]}</code></pre>
			{/if}
		</div>
	</div>
</section>
