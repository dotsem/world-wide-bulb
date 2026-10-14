<script lang="ts">
	import { resolve } from '$app/paths';
	import type { ResolvedPathname } from '$app/types';
</script>

<section
	id="overview"
	class="scroll-mt-28 min-h-[60vh] flex flex-col justify-center space-y-6 py-12 border-b border-app-border"
>
	<div class="space-y-4 max-w-3xl">
		<h1 class="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-app-text tracking-tight">
			World Wide Bulb
		</h1>

		<p class="text-base sm:text-lg text-app-muted leading-relaxed max-w-2xl">
			A real-time collaborative light switch shared across the internet. Built as a single,
			self-contained Go binary with an embedded frontend and SQLite database.
		</p>
	</div>

	<div class="flex flex-wrap items-center gap-3 pt-2">
		<a
			href="https://github.com/dotsem/world-wide-bulb"
			target="_blank"
			rel="noopener noreferrer external"
			class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs sm:text-sm font-semibold text-app-accent-contrast bg-app-accent hover:opacity-90 transition-opacity focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer"
		>
			<svg class="w-4 h-4 fill-current" viewBox="0 0 24 24" aria-hidden="true">
				<path
					fill-rule="evenodd"
					clip-rule="evenodd"
					d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
				/>
			</svg>
			<span>GitHub</span>
		</a>

		<a
			href={resolve('/docs#self-hosting') as ResolvedPathname}
			class="inline-flex items-center px-4 py-2.5 rounded-xl text-xs sm:text-sm font-semibold text-app-text bg-app-surface hover:bg-app-surface-hover border border-app-border transition-colors cursor-pointer"
		>
			<span>Self-Hosting</span>
		</a>

		<a
			href={resolve('/docs#rest-api') as ResolvedPathname}
			class="inline-flex items-center px-4 py-2.5 rounded-xl text-xs sm:text-sm font-semibold text-app-muted hover:text-app-text hover:bg-app-surface-hover/60 border border-transparent hover:border-app-border transition-colors cursor-pointer"
		>
			<span>API Reference</span>
		</a>
	</div>
</section>
