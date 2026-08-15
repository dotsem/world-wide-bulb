<script lang="ts">
    import "./layout.css";
    import favicon from "$lib/assets/favicon.svg";
    import NavToolBar from "$lib/components/NavToolBar.svelte";

    let { children } = $props();
</script>

<svelte:head>
    <title>World Wide Bulb</title>
    <link rel="icon" href={favicon} />
</svelte:head>

<div
    class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-amber-500/30 selection:text-amber-200"
>
    <NavToolBar />

    <main class="flex-1">
        {@render children()}
    </main>
</div>
