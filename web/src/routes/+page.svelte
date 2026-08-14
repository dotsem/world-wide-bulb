<script lang="ts">
    import { bulbState, Bulb, CooldownNotice, ReasonPrompt } from "$lib";
</script>

<CooldownNotice />
<Bulb isOn={bulbState.isOn} disabled={bulbState.isCooldownActive} durationMs={200} onclick={bulbState.toggle} />
<ReasonPrompt />
