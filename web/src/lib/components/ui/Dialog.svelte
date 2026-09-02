<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { X } from '@lucide/svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		isOpen?: boolean;
		title?: string;
		maxWidth?: string;
		header?: Snippet;
		children?: Snippet;
		footer?: Snippet;
		onClose?: () => void;
	}

	let {
		isOpen = $bindable(false),
		title,
		maxWidth = 'max-w-md',
		header,
		children,
		footer,
		onClose
	}: Props = $props();

	function close() {
		isOpen = false;
		onClose?.();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && isOpen) {
			close();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
	<div
		transition:fade={{ duration: 200 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs"
		role="presentation"
		onclick={close}
		onkeydown={(e) => e.key === 'Escape' && close()}
	>
		<div
			transition:fly={{ y: 20, duration: 250 }}
			class="relative w-full {maxWidth} p-6 rounded-2xl bg-app-surface-solid border border-app-border text-app-text shadow-2xl shadow-black/50 space-y-4"
			role="dialog"
			aria-modal="true"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between border-b border-app-border/50 pb-3">
				{#if header}
					{@render header()}
				{:else if title}
					<h2 class="text-lg font-bold tracking-tight">{title}</h2>
				{/if}
				<button
					type="button"
					onclick={close}
					aria-label="Close dialog"
					class="p-1.5 rounded-lg text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors focus:outline-none focus:ring-2 focus:ring-app-accent/50 cursor-pointer ml-auto"
				>
					<X size={18} />
				</button>
			</div>

			{#if children}
				<div class="space-y-3 text-sm text-app-muted leading-relaxed">
					{@render children()}
				</div>
			{/if}

			{#if footer}
				<div
					class="pt-2 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-t border-app-border/40"
				>
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
