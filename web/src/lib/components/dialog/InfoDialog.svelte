<script lang="ts">
	import { Bug, Code, Lightbulb } from '@lucide/svelte';
	import type { LucideIcon } from '@lucide/svelte';
	import Dialog from '$lib/components/ui/Dialog.svelte';

	let { isOpen = $bindable(false) }: { isOpen?: boolean } = $props();

	$effect(() => {
		if (typeof localStorage !== 'undefined') {
			const hasSeenIntro = localStorage.getItem('wwb_seen_intro');
			if (!hasSeenIntro) {
				isOpen = true;
			}
		}
	});

	function handleClose() {
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem('wwb_seen_intro', 'true');
		}
	}
</script>

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
		class="inline-flex items-center gap-1.5 text-xs font-medium text-app-muted hover:text-app-accent transition-colors duration-200 underline underline-offset-2 focus:outline-none focus:ring-2 focus:ring-app-accent/50"
	>
		<Icon size={14} />
		<span>{label}</span>
	</a>
{/snippet}

<Dialog bind:isOpen onClose={handleClose}>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<h2 class="text-lg font-bold tracking-tight">World Wide Bulb</h2>
			<Lightbulb size={20} class="text-app-accent animate-spin" />
		</div>
	{/snippet}

	<p>A real-time collaborative light switch shared by everyone on the internet.</p>
	<p>Click on the lightbulb to flip its state instantly for all online visitors.</p>
	<p>Explore real-time toggle history and see reasons left by people worldwide.</p>
	<small class="text-app-muted/75">*You can toggle the lightbulb once every 10 seconds.</small>

	{#snippet footer()}
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
			onclick={() => {
				isOpen = false;
				handleClose();
			}}
			class="px-4 py-1.5 rounded-lg text-xs font-semibold cursor-pointer text-app-accent-contrast bg-app-accent hover:opacity-90 transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-app-accent/50 self-end sm:self-center"
		>
			Cool!
		</button>
	{/snippet}
</Dialog>
