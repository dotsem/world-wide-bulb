<script lang="ts">
    import { bulbState, Bulb } from "$lib";
</script>

<Bulb isOn={bulbState.isOn} durationMs={200} onclick={bulbState.toggle} />
