<script lang="ts">
	import { Lightbulb, Radio, Users } from '@lucide/svelte';
</script>

<section id="use-cases" class="scroll-mt-28 space-y-5 pt-2">
	<div class="space-y-1">
		<h2 class="text-xl font-bold tracking-tight text-app-text">Use Cases</h2>
		<p class="text-xs sm:text-sm text-app-muted">
			World Wide Bulb can be used for multiple things.
		</p>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
		<!-- Physical IoT Switch -->
		<div class="flex flex-col p-5 rounded-2xl bg-app-surface border border-app-border space-y-3">
			<div class="text-app-accent">
				<Lightbulb size={20} />
			</div>
			<div class="space-y-1.5">
				<h3 class="text-sm font-bold text-app-text">Physical IoT Relay</h3>
				<p class="text-xs text-app-muted leading-relaxed">
					Connect an ESP32, Raspberry Pi, or Shelly smart relay to the <code
						class="px-1 py-0.5 rounded bg-app-surface-solid text-app-accent font-mono"
						>/api/v1/events</code
					> SSE stream to physically toggle a real room light when the state changes.
				</p>
			</div>
		</div>

		<!-- Community & Stream Gimmick -->
		<div class="flex flex-col p-5 rounded-2xl bg-app-surface border border-app-border space-y-3">
			<div class="text-app-accent">
				<Users size={20} />
			</div>
			<div class="space-y-1.5">
				<h3 class="text-sm font-bold text-app-text">Stream Overlay &amp; Gimmick</h3>
				<p class="text-xs text-app-muted leading-relaxed">
					Embed as an OBS browser source or live widget for streams and community dashboards.
					Viewers can toggle the lamp in real time and leave custom reasons.
				</p>
			</div>
		</div>

		<!-- Private Remote Switch -->
		<div class="flex flex-col p-5 rounded-2xl bg-app-surface border border-app-border space-y-3">
			<div class="text-app-accent">
				<Radio size={20} />
			</div>
			<div class="space-y-1.5">
				<h3 class="text-sm font-bold text-app-text">Private Remote Switch</h3>
				<p class="text-xs text-app-muted leading-relaxed">
					Run on a home server, Raspberry Pi, or Tailscale network as a persistent, zero-latency
					private switch with full audit history across reboots.
				</p>
			</div>
		</div>
	</div>
</section>
