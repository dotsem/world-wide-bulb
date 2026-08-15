<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends HTMLInputAttributes {
		label: string;
		value?: string;
		maxlength?: number;
		error?: string | null;
	}

	let {
		label,
		value = $bindable(''),
		maxlength,
		error = null,
		disabled = false,
		id,
		class: className = '',
		...restProps
	}: Props = $props();

	let inputId = $derived(id || `floating-input-${Math.random().toString(36).substring(2, 9)}`);
	let isFocused = $state(false);
	let isFloating = $derived(isFocused || (value !== undefined && value.length > 0));
</script>

<div class="relative w-full {className}">
	<div class="relative flex items-center">
		<input
			{...restProps}
			id={inputId}
			bind:value
			{disabled}
			{maxlength}
			placeholder=" "
			onfocus={() => (isFocused = true)}
			onblur={() => (isFocused = false)}
			class="w-full px-3.5 py-3 rounded-xl bg-app-surface-hover/80 border text-sm text-app-text transition-all duration-200 outline-none pr-14 {error
				? 'border-rose-500/80 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/20'
				: isFocused
					? 'border-amber-400 focus:border-amber-400 focus:ring-2 focus:ring-amber-400/20'
					: 'border-app-border hover:border-app-border/80'} {disabled
				? 'opacity-50 cursor-not-allowed'
				: ''}"
		/>

		<label
			for={inputId}
			class="absolute left-3 transition-all duration-200 pointer-events-none select-none z-10 {isFloating
				? 'top-0 -translate-y-1/2 text-xs px-1.5 py-0.5 leading-none rounded-sm bg-app-surface-solid font-medium ' +
					(isFocused ? 'text-amber-400 font-semibold' : 'text-app-muted')
				: 'top-1/2 -translate-y-1/2 text-sm text-app-muted/70'}"
		>
			{label}
		</label>

		{#if maxlength}
			<span
				class="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-mono transition-colors {isFocused
					? 'text-amber-400/80'
					: 'text-app-muted/60'}"
			>
				{value.length}/{maxlength}
			</span>
		{/if}
	</div>

	{#if error}
		<p class="mt-1 text-xs text-rose-400 font-medium animate-in fade-in">{error}</p>
	{/if}
</div>
