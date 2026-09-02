<script lang="ts">
	import jsCode from '../../../../examples/sse/javascript/client.js?raw';
	import pyCode from '../../../../examples/sse/python/client.py?raw';
	import goCode from '../../../../examples/sse/go/main.go?raw';
	import { Check, Code, Copy } from '@lucide/svelte';

	type Language = 'javascript' | 'python' | 'go' | 'curl';

	let { sseUrl }: { sseUrl: string } = $props();

	let activeLang = $state<Language>('javascript');
	let copiedLang = $state<string | null>(null);

	const curlCode = $derived(`curl -N ${sseUrl}`);

	const snippets = $derived<Record<Language, string>>({
		javascript: jsCode,
		python: pyCode,
		go: goCode,
		curl: curlCode
	});

	function copySnippet(lang: Language) {
		if (typeof navigator !== 'undefined') {
			navigator.clipboard.writeText(snippets[lang]);
			copiedLang = lang;
			setTimeout(() => {
				if (copiedLang === lang) copiedLang = null;
			}, 2000);
		}
	}
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
				{activeLang === 'curl'
					? 'terminal'
					: `examples/sse/${activeLang === 'javascript' ? 'javascript/client.js' : activeLang === 'python' ? 'python/client.py' : 'go/main.go'}`}
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
		<pre
			class="p-5 text-xs sm:text-sm font-mono text-app-text overflow-x-auto leading-relaxed"><code
				>{snippets[activeLang]}</code
			></pre>
	</div>
</section>
