<script lang="ts">
	import { Check, Container, Copy, Server, Terminal } from '@lucide/svelte';

	let activeTab = $state<'binary' | 'docker'>('binary');
	let copied = $state(false);

	const binarySnippet = `./wwb -port 8080 -db ./wwb.db`;
	const dockerSnippet = `docker run -d \\
  --name world-wide-bulb \\
  -p 8080:8080 \\
  -v ./data:/data \\
  -e DB_PATH=/data/wwb.db \\
  ghcr.io/dotsem/world-wide-bulb:latest`;

	async function copySnippet() {
		const text = activeTab === 'binary' ? binarySnippet : dockerSnippet;
		if (typeof navigator !== 'undefined' && navigator.clipboard) {
			await navigator.clipboard.writeText(text);
			copied = true;
			setTimeout(() => (copied = false), 2000);
		}
	}
</script>

<section id="self-hosting" class="scroll-mt-28 space-y-6 pt-4">
	<div class="space-y-1.5">
		<h2 class="text-2xl font-bold tracking-tight text-app-text flex items-center gap-2.5">
			<Server size={22} class="text-app-accent" />
			<span>Self-Hosting</span>
		</h2>
		<p class="text-sm text-app-muted">
			World Wide Bulb is designed for friction-free self-hosting with zero external dependencies.
		</p>
	</div>

	<!-- Why Single Binary Architecture Box -->
	<div class="p-5 rounded-2xl bg-app-surface border border-app-border space-y-3">
		<h3 class="text-sm font-bold text-app-text flex items-center gap-2">
			<Terminal size={16} class="text-app-accent" />
			<span>Single Self-Contained Binary</span>
		</h3>
		<p class="text-xs text-app-muted leading-relaxed">
			The Go binary embeds the entire compiled SvelteKit web client using Go's <code
				class="px-1.5 py-0.5 rounded bg-app-bg text-app-accent font-mono">embed.FS</code
			> and runs an in-process pure Go SQLite engine. No Node.js runtime, no external database container,
			and no web server reverse proxy are required to get started.
		</p>
	</div>

	<!-- Quickstart Snippet Card -->
	<div class="rounded-2xl border border-app-border bg-app-surface overflow-hidden">
		<!-- Tab Switcher Header -->
		<div
			class="flex items-center justify-between border-b border-app-border bg-app-surface-solid px-4 py-2.5"
		>
			<div class="flex items-center gap-2">
				<button
					type="button"
					onclick={() => (activeTab = 'binary')}
					class={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors cursor-pointer ${
						activeTab === 'binary'
							? 'bg-app-accent text-app-accent-contrast shadow-sm'
							: 'text-app-muted hover:text-app-text hover:bg-app-surface-hover'
					}`}
				>
					<Terminal size={14} />
					<span>Standalone Binary</span>
				</button>
				<button
					type="button"
					onclick={() => (activeTab = 'docker')}
					class={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors cursor-pointer ${
						activeTab === 'docker'
							? 'bg-app-accent text-app-accent-contrast shadow-sm'
							: 'text-app-muted hover:text-app-text hover:bg-app-surface-hover'
					}`}
				>
					<Container size={14} />
					<span>Docker Container</span>
				</button>
			</div>

			<button
				type="button"
				onclick={copySnippet}
				aria-label="Copy command to clipboard"
				title="Copy command"
				class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer"
			>
				{#if copied}
					<Check size={14} class="text-emerald-400" />
					<span class="text-emerald-400">Copied!</span>
				{:else}
					<Copy size={14} />
					<span>Copy</span>
				{/if}
			</button>
		</div>

		<!-- Code Display -->
		<div class="p-4 bg-app-bg overflow-x-auto text-xs font-mono text-app-text leading-relaxed">
			{#if activeTab === 'binary'}
				<div class="text-app-muted mb-2"># Run directly (or build with `just build`)</div>
				<pre class="select-all">{binarySnippet}</pre>
			{:else}
				<div class="text-app-muted mb-2">
					# Persist database across container rebuilds via volume
				</div>
				<pre class="select-all">{dockerSnippet}</pre>
			{/if}
		</div>

		<!-- Flags Reference Footer -->
		<div
			class="p-4 border-t border-app-border bg-app-surface-solid/50 text-xs text-app-muted space-y-1.5"
		>
			<div class="font-semibold text-app-text">Key Configuration Flags &amp; Environment:</div>
			<ul class="space-y-1 list-disc list-inside">
				<li>
					<code class="font-mono text-app-accent">-port</code> (or
					<code class="font-mono text-app-accent">PORT</code>): Web server port (default:
					<code class="font-mono">8080</code>)
				</li>
				<li>
					<code class="font-mono text-app-accent">-db</code> (or
					<code class="font-mono text-app-accent">DB_PATH</code>): SQLite file path (default:
					<code class="font-mono">./wwb.db</code>)
				</li>
				<li>
					<code class="font-mono text-app-accent">-cooldown</code> (or
					<code class="font-mono text-app-accent">COOLDOWN_SEC</code>): Toggle rate limit in seconds
					(default: <code class="font-mono">10</code>)
				</li>
			</ul>
		</div>
	</div>
</section>
