<script lang="ts">
	type Method = 'GET' | 'POST' | 'SSE' | 'WS';

	let {
		method,
		isInternal = false
	}: {
		method: Method;
		isInternal?: boolean;
	} = $props();

	const methodStyles: Record<Method, string> = {
		GET: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
		POST: 'bg-sky-500/15 text-sky-400 border-sky-500/30',
		SSE: 'bg-amber-500/15 text-amber-400 border-amber-500/30',
		WS: 'bg-purple-500/15 text-purple-400 border-purple-500/30'
	};
</script>

<div class="inline-flex items-center gap-1.5 shrink-0">
	<span
		class="px-2 py-0.5 rounded-md text-[11px] font-bold font-mono uppercase tracking-wider border {methodStyles[
			method
		] || methodStyles.GET}"
	>
		{method}
	</span>
	{#if isInternal}
		<span
			class="px-1.5 py-0.5 rounded-md text-[10px] font-medium font-mono uppercase tracking-wider bg-app-surface text-app-muted border border-app-border"
			title="Requires allowed origin or self-hosted deployment"
		>
			Internal
		</span>
	{/if}
</div>
