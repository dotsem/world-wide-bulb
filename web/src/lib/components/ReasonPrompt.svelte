<script lang="ts">
	import { onDestroy } from 'svelte';
	import { bulbState } from '$lib/state/bulb.svelte';
	import { fly } from 'svelte/transition';
	import FloatingInput from './ui/FloatingInput.svelte';

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
		transition:fly={{ y: 100, duration: 300 }}
		class="fixed bottom-24 left-1/2 -translate-x-1/2 z-40 w-[90vw] max-w-md p-5 rounded-2xl bg-app-surface-solid border border-app-border text-app-text shadow-2xl shadow-black/50 transition-all duration-300"
		role="dialog"
		aria-labelledby="reason-prompt-title"
	>
		{#if isSuccess}
			<div
				class="flex items-center gap-2 py-2 text-app-success text-sm font-medium animate-in fade-in"
			>
				Thank you for your feedback!
			</div>
		{:else}
			<form onsubmit={handleSubmit} class="space-y-3 pt-2">
				<FloatingInput
					label="Why did you decide to {actionText} the lamp?"
					bind:value={reason}
					maxlength={100}
					disabled={isSubmitting}
					error={errorMsg}
				/>

				<div class="flex items-center justify-end gap-2 pt-1">
					<button
						type="button"
						onclick={handleSkip}
						class="px-3.5 py-1.5 rounded-lg cursor-pointer text-xs font-medium text-app-muted hover:text-app-text transition-colors focus:outline-none"
					>
						Skip
					</button>
					<button
						type="submit"
						disabled={!reason.trim() || isSubmitting}
						class="px-4 py-1.5 rounded-lg text-xs font-semibold cursor-pointer text-app-accent-contrast bg-app-accent hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm focus:outline-none focus:ring-2 focus:ring-app-accent/50"
					>
						{isSubmitting ? 'Saving...' : 'Submit'}
					</button>
				</div>
			</form>
		{/if}
	</div>
{/if}
