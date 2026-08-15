<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { Bug, Code, Lightbulb, X } from '@lucide/svelte';
	import type { LucideIcon } from '@lucide/svelte';

	let { isOpen = $bindable(false) }: { isOpen?: boolean } = $props();

	$effect(() => {
		if (typeof localStorage !== 'undefined') {
			const hasSeenIntro = localStorage.getItem('wwb_seen_intro');
			if (!hasSeenIntro) {
				isOpen = true;
			}
		}
	});

	function closeDialog() {
		isOpen = false;
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem('wwb_seen_intro', 'true');
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && isOpen) {
			closeDialog();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#snippet externalLink({
	href,
	label,
	icon: Icon
}: {
	href: string;
	label: string;
	icon: LucideIcon;
})}
	<a
		{href}
		target="_blank"
		rel="noopener noreferrer"
		class="inline-flex items-center gap-1.5 text-xs font-medium text-app-muted hover:text-app-accent transition-color duration-400 underline underline-offset-2 focus:outline-none focus:ring-2 focus:ring-app-accent/50"
	>
		<Icon size={14} />
		<span>{label}</span>
	</a>
{/snippet}

{#if isOpen}
	<div
		transition:fade={{ duration: 200 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
		role="presentation"
		onclick={closeDialog}
		onkeydown={(e) => e.key === 'Escape' && closeDialog()}
	>
		<div
			transition:fly={{ y: 20, duration: 250 }}
			class="relative w-full max-w-md p-6 rounded-2xl bg-app-surface-solid border border-app-border text-app-text shadow-2xl shadow-black/50 space-y-4"
			role="dialog"
			aria-modal="true"
			aria-labelledby="info-dialog-title"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between border-b border-app-border/50 pb-3">
				<div class="flex items-center gap-2">
					<h2 id="info-dialog-title" class="text-lg font-bold tracking-tight">World Wide Bulb</h2>
					<Lightbulb size={20} class="text-app-accent animate-spin" />
				</div>
				<button
					type="button"
					onclick={closeDialog}
					aria-label="Close dialog"
					class="p-1.5 rounded-lg text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer"
				>
					<X size={18} />
				</button>
			</div>

			<div class="space-y-2 text-sm text-app-muted leading-relaxed">
				<p>A real-time collaborative light switch shared by everyone on the internet.</p>
				<p>Click on the lightbulb to flip its state instantly for all online visitors.</p>
				<p>Explore real-time toggle history and see reasons left by people worldwide.</p>
				<small>*You can toggle the lightbulb once every 10 seconds.</small>
			</div>

			<div
				class="pt-2 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-t border-app-border/40"
			>
				<div class="flex flex-col gap-1.5">
					{@render externalLink({
						href: 'https://github.com/dotsem/world-wide-bulb',
						label: 'Check out the code on GitHub',
						icon: Code
					})}
					{@render externalLink({
						href: 'https://github.com/dotsem/world-wide-bulb/issues',
						label: 'Report Bugs & Feature Requests',
						icon: Bug
					})}
				</div>

				<button
					type="button"
					onclick={closeDialog}
					class="px-4 py-1.5 rounded-lg text-xs font-semibold cursor-pointer text-app-accent-contrast bg-app-accent hover:opacity-90 transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-app-accent/50 self-end sm:self-center"
				>
					Cool!
				</button>
			</div>
		</div>
	</div>
{/if}
