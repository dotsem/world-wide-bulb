<script lang="ts">
	import { fly } from 'svelte/transition';
	import { CircleAlert, TriangleAlert, Info, CheckCircle2, X } from '@lucide/svelte';
	import { toastState, type ToastType } from '$lib/state/toast.svelte';

	interface Props {
		position?: 'bottom-right' | 'bottom-center' | 'top-right' | 'top-center';
		class?: string;
	}

	let { position = 'bottom-right', class: className = '' }: Props = $props();

	const positionClasses = {
		'bottom-right':
			'bottom-4 inset-x-0 mx-auto items-center sm:inset-x-auto sm:mx-0 sm:right-4 sm:items-end',
		'bottom-center': 'bottom-4 inset-x-0 mx-auto items-center',
		'top-right':
			'top-4 inset-x-0 mx-auto items-center sm:inset-x-auto sm:mx-0 sm:right-4 sm:items-end',
		'top-center': 'top-4 inset-x-0 mx-auto items-center'
	};

	function getTypeStyles(type: ToastType) {
		switch (type) {
			case 'error':
				return {
					border: 'border-app-danger/40',
					iconClass: 'text-app-danger'
				};
			case 'warning':
				return {
					border: 'border-app-accent/40',
					iconClass: 'text-app-accent'
				};
			case 'success':
				return {
					border: 'border-app-success/40',
					iconClass: 'text-app-success'
				};
			case 'info':
			default:
				return {
					border: 'border-app-border',
					iconClass: 'text-app-muted'
				};
		}
	}
</script>

{#if toastState.toasts.length > 0}
	<aside
		class="fixed z-50 flex flex-col gap-2.5 pointer-events-none p-4 max-w-md w-full {positionClasses[
			position
		]} {className}"
		aria-label="Notifications"
	>
		{#each toastState.toasts as item (item.id)}
			{@const style = getTypeStyles(item.type)}
			<div
				transition:fly={{ y: 16, duration: 200 }}
				class="pointer-events-auto w-full flex items-start gap-3 p-3.5 rounded-xl bg-app-surface-solid border {style.border} text-app-text shadow-xl shadow-black/40 backdrop-blur-md"
				role={item.type === 'error' ? 'alert' : 'status'}
				aria-live={item.type === 'error' ? 'assertive' : 'polite'}
				data-testid="toast-card"
				data-toast-type={item.type}
			>
				<div class="mt-0.5 shrink-0 {style.iconClass}">
					{#if item.type === 'error'}
						<CircleAlert size={18} />
					{:else if item.type === 'warning'}
						<TriangleAlert size={18} />
					{:else if item.type === 'success'}
						<CheckCircle2 size={18} />
					{:else}
						<Info size={18} />
					{/if}
				</div>

				<div class="flex-1 text-sm font-medium leading-snug wrap-break-word">
					{item.message}
				</div>

				<button
					type="button"
					onclick={() => toastState.dismiss(item.id)}
					aria-label="Dismiss notification"
					class="shrink-0 p-1 -mr-1 -mt-1 rounded-lg text-app-muted hover:text-app-text hover:bg-app-surface-hover transition-colors cursor-pointer focus:outline-none focus:ring-2 focus:ring-app-accent/50"
					data-testid="toast-dismiss"
				>
					<X size={16} />
				</button>
			</div>
		{/each}
	</aside>
{/if}
