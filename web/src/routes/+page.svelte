<script lang="ts">
    import { bulbState, Bulb, CooldownNotice } from "$lib";
</script>

<CooldownNotice />
<Bulb isOn={bulbState.isOn} disabled={bulbState.isCooldownActive} durationMs={200} onclick={bulbState.toggle} />
