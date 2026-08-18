<script lang="ts">
	import { bulbState, Bulb, CooldownNotice, ReasonPrompt } from '$lib';
	import Header from '$lib/components/Header.svelte';
</script>

<Header />
<div
	class="flex-1 flex flex-col items-center justify-center relative overflow-hidden px-4 py-8 min-h-[calc(100vh-5rem)]"
>
	<CooldownNotice />
	<Bulb
		isOn={bulbState.isOn}
		disabled={bulbState.isCooldownActive}
		class="h-[60vh] w-[60vh] md:h-[70vh] md:w-[70vh]"
		onclick={bulbState.toggle}
	/>
	<ReasonPrompt />
</div>
