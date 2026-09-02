<script lang="ts">
	import { Sparkles } from '@lucide/svelte';
</script>

<section class="rounded-2xl border border-app-border bg-app-surface/60 p-6 sm:p-7 space-y-4">
	<div class="flex items-center gap-2">
		<Sparkles size={20} class="text-app-accent" />
		<h2 class="text-lg font-bold text-app-text">Event Types & Payloads</h2>
	</div>

	<div class="space-y-4 text-xs sm:text-sm">
		<div class="p-4 rounded-xl bg-app-surface border border-app-border space-y-2">
			<div class="flex items-center gap-2">
				<span class="font-mono font-bold text-app-accent">event: state_changed</span>
				<span class="text-xs text-app-muted">(Sent on connect and whenever lamp is toggled)</span>
			</div>
			<pre
				class="p-3 rounded-lg bg-app-surface-solid border border-app-border/60 text-xs font-mono text-app-text overflow-x-auto"><code
					>&#123;
  "type": "state_changed",
  "id": 42,
  "state": true,
  "reason": "Midnight coding",
  "created_at": "2026-09-02T17:35:00Z"
&#125;</code
				></pre>
		</div>

		<div class="p-4 rounded-xl bg-app-surface border border-app-border space-y-2">
			<div class="flex items-center gap-2">
				<span class="font-mono font-bold text-app-accent">event: reason_updated</span>
				<span class="text-xs text-app-muted">(Sent when someone attaches a reason)</span>
			</div>
			<pre
				class="p-3 rounded-lg bg-app-surface-solid border border-app-border/60 text-xs font-mono text-app-text overflow-x-auto"><code
					>&#123;
  "type": "reason_updated",
  "id": 42,
  "reason": "Party at the office!"
&#125;</code
				></pre>
		</div>

		<div class="p-4 rounded-xl bg-app-surface border border-app-border space-y-2">
			<div class="flex items-center gap-2">
				<span class="font-mono font-bold text-app-accent">event: ping</span>
				<span class="text-xs text-app-muted">(Sent every 15s keep-alive heartbeat)</span>
			</div>
			<p class="text-xs text-app-muted">
				Empty payload data to prevent reverse-proxy and router timeouts.
			</p>
		</div>
	</div>
</section>
