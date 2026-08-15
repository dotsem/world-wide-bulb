<script lang="ts">
	import { onDestroy } from 'svelte';
	import { bulbState } from '$lib/state/bulb.svelte';

	let reason = $state('');
	let isSubmitting = $state(false);
	let isSuccess = $state(false);
	let errorMsg = $state<string | null>(null);
	let successTimeout: ReturnType<typeof setTimeout> | undefined;
	let autoDismissTimeout: ReturnType<typeof setTimeout> | undefined;

	let isVisible = $derived(bulbState.showReasonPrompt);
	let actionText = $derived(bulbState.lastActionState ? 'turn on' : 'turn off');

	$effect(() => {
		if (isVisible) {
			startAutoDismissTimer();
		} else {
			clearAutoDismissTimer();
		}
	});

	function startAutoDismissTimer() {
		clearAutoDismissTimer();
		autoDismissTimeout = setTimeout(() => {
			if (!reason.trim()) {
				handleSkip();
			}
		}, 60000);
	}

	function clearAutoDismissTimer() {
		if (autoDismissTimeout) {
			clearTimeout(autoDismissTimeout);
			autoDismissTimeout = undefined;
		}
	}

	onDestroy(() => {
		if (successTimeout) clearTimeout(successTimeout);
		clearAutoDismissTimer();
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const trimmed = reason.trim();
		if (!trimmed || isSubmitting) return;

		isSubmitting = true;
		errorMsg = null;

		try {
			clearAutoDismissTimer();
			await bulbState.submitReason(trimmed);
			isSuccess = true;
			reason = '';
			if (successTimeout) clearTimeout(successTimeout);
			successTimeout = setTimeout(() => {
				isSuccess = false;
			}, 1500);
		} catch (err: any) {
			errorMsg = err?.message || 'Failed to submit reason';
		} finally {
			isSubmitting = false;
		}
	}

	function handleSkip() {
		clearAutoDismissTimer();
		reason = '';
		errorMsg = null;
		bulbState.dismissReason();
	}
</script>

{#if isVisible}
	<div
		class="fixed bottom-8 left-1/2 -translate-x-1/2 z-40 w-[90vw] max-w-md p-5 rounded-2xl bg-slate-900/90 border border-amber-500/20 text-slate-100 shadow-2xl shadow-amber-950/20 backdrop-blur-xl transition-all duration-300 animate-in fade-in slide-in-from-bottom-4"
		role="dialog"
		aria-labelledby="reason-prompt-title"
	>
		<div class="flex items-start justify-between gap-3 mb-3">
			<div class="flex items-center gap-2">
				<span class="text-amber-400 text-lg">💡</span>
				<h3 id="reason-prompt-title" class="text-sm font-semibold text-amber-200 tracking-wide">
					Why did you decide to {actionText} the lamp?
				</h3>
			</div>
			<button
				type="button"
				onclick={handleSkip}
				class="text-slate-400 hover:text-slate-200 p-1 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-amber-500/50"
				aria-label="Close reason prompt"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</div>

		{#if isSuccess}
			<div
				class="flex items-center gap-2 py-2 text-emerald-400 text-sm font-medium animate-in fade-in"
			>
				<span>✨</span> Reason recorded! Thank you.
			</div>
		{:else}
			<form onsubmit={handleSubmit} class="space-y-3">
				<div class="relative">
					<input
						type="text"
						bind:value={reason}
						maxlength={100}
						placeholder="e.g. It was too dark, felt like turning it on..."
						class="w-full px-3.5 py-2.5 rounded-xl bg-slate-800/80 border border-slate-700 text-sm text-slate-100 placeholder-slate-400 focus:outline-none focus:border-amber-500/60 focus:ring-2 focus:ring-amber-500/20 transition-all pr-14"
						disabled={isSubmitting}
					/>
					<span class="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-mono text-slate-500">
						{reason.length}/100
					</span>
				</div>

				{#if errorMsg}
					<p class="text-xs text-rose-400 font-medium">{errorMsg}</p>
				{/if}

				<div class="flex items-center justify-end gap-2 pt-1">
					<button
						type="button"
						onclick={handleSkip}
						class="px-3.5 py-1.5 rounded-lg text-xs font-medium text-slate-400 hover:text-slate-200 transition-colors focus:outline-none"
					>
						Skip
					</button>
					<button
						type="submit"
						disabled={!reason.trim() || isSubmitting}
						class="px-4 py-1.5 rounded-lg text-xs font-semibold text-amber-950 bg-amber-400 hover:bg-amber-300 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-amber-400/50"
					>
						{isSubmitting ? 'Saving...' : 'Submit'}
					</button>
				</div>
			</form>
		{/if}
	</div>
{/if}
