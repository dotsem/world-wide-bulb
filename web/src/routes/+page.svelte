<script lang="ts">
    import Bulb from "$lib/components/Bulb.svelte";

    let state: boolean = $state(false);

    function toggle() {
        state = !state;
    }
</script>

<Bulb isOn={state} durationMs={200} onclick={toggle} />
