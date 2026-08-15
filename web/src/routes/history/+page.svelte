<script lang="ts">
    import { onMount } from "svelte";
    import { restApi } from "$lib/api/rest";
    import { wsClient } from "$lib/api/ws.svelte";
    import type { HistoryItem } from "$lib/types/rest.types";

    let toggles = $state<HistoryItem[]>([]);
    let nextCursor = $state<number | undefined>(undefined);
    let hasMore = $state(true);
    let loading = $state(false);
    let initialLoading = $state(true);
    let error = $state<string | null>(null);

    let sentinelEl = $state<HTMLElement | null>(null);

    function formatDate(dateStr: string): string {
        try {
            const date = new Date(dateStr);
            return date.toLocaleString(undefined, {
                dateStyle: "medium",
                timeStyle: "short",
            });
        } catch {
            return dateStr;
        }
    }

    function getRelativeTime(dateStr: string): string {
        try {
            const now = Date.now();
            const elapsedSec = Math.floor(
                (now - new Date(dateStr).getTime()) / 1000,
            );
            if (elapsedSec < 5) return "just now";
            if (elapsedSec < 60) return `${elapsedSec}s ago`;
            const elapsedMin = Math.floor(elapsedSec / 60);
            if (elapsedMin < 60) return `${elapsedMin}m ago`;
            const elapsedHours = Math.floor(elapsedMin / 60);
            if (elapsedHours < 24) return `${elapsedHours}h ago`;
            const elapsedDays = Math.floor(elapsedHours / 24);
            return `${elapsedDays}d ago`;
        } catch {
            return "";
        }
    }

    async function loadMore() {
        if (loading || !hasMore) return;

        loading = true;
        error = null;

        try {
            const res = await restApi.history(20, nextCursor);
            toggles = [...toggles, ...res.toggles];
            nextCursor = res.next_cursor;
            hasMore = res.has_more;
        } catch (err: any) {
            error = err?.message || "Failed to load history items";
        } finally {
            loading = false;
            initialLoading = false;
        }
    }

    $effect(() => {
        if (!sentinelEl) return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !loading) {
                    loadMore();
                }
            },
            { rootMargin: "300px" },
        );

        observer.observe(sentinelEl);
        return () => observer.disconnect();
    });

    onMount(() => {
        loadMore();

        const unsubState = wsClient.onStateChange((msg) => {
            if (msg.id === undefined || toggles.some((t) => t.id === msg.id)) return;
            const newItem: HistoryItem = {
                id: msg.id,
                state: msg.state,
                reason: msg.reason || "",
                created_at: msg.created_at || new Date().toISOString(),
            };
            toggles = [newItem, ...toggles];
        });

        const unsubReason = wsClient.onReasonUpdate((msg) => {
            const targetId = msg.toggle_id ?? (typeof msg.id === "number" ? msg.id : undefined);
            if (targetId === undefined) return;
            toggles = toggles.map((t) =>
                t.id === targetId ? { ...t, reason: msg.reason } : t
            );
        });

        return () => {
            unsubState();
            unsubReason();
        };
    });
</script>

<svelte:head>
    <title>Toggle History - World Wide Bulb</title>
    <meta
        name="description"
        content="View historical toggle events and community reasons for World Wide Bulb."
    />
</svelte:head>

<div class="mx-auto max-w-4xl px-4 py-8 sm:px-6">
    <div
        class="mb-8 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between border-b border-slate-800/80 pb-6"
    >
        <h1
            class="text-3xl font-extrabold tracking-tight text-white flex items-center gap-3"
        >
            Toggle History
        </h1>
    </div>

    {#if initialLoading}
        <div class="space-y-4">
            {#each Array(5) as _}
                <div
                    class="h-20 w-full animate-pulse rounded-xl bg-slate-900/60 border border-slate-800/50"
                ></div>
            {/each}
        </div>
    {:else if error && toggles.length === 0}
        <div
            class="rounded-xl border border-red-500/20 bg-red-500/10 p-6 text-center text-red-300"
        >
            <p class="font-medium">{error}</p>
            <button
                onclick={() => loadMore()}
                class="mt-4 rounded-lg bg-red-500/20 px-4 py-2 text-sm font-semibold text-red-200 hover:bg-red-500/30 transition-colors"
            >
                Try Again
            </button>
        </div>
    {:else if toggles.length === 0}
        <div
            class="rounded-xl border border-slate-800 bg-slate-900/40 p-12 text-center"
        >
            <div
                class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-slate-800/80 text-slate-400 mb-3"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                    />
                </svg>
            </div>
            <h3 class="text-base font-semibold text-slate-200">
                No toggle history found
            </h3>
            <p class="mt-1 text-sm text-slate-400">
                Be the first to flip the switch!
            </p>
            <a
                href="/"
                class="mt-4 inline-block rounded-lg bg-amber-500/20 border border-amber-500/40 px-4 py-2 text-sm font-semibold text-amber-300 hover:bg-amber-500/30 transition-colors"
            >
                Go to Bulb
            </a>
        </div>
    {:else}
        <div class="space-y-3">
            {#each toggles as item (item.id)}
                <div
                    class="group relative flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-xl border border-slate-800/80 bg-slate-900/60 p-4 backdrop-blur-sm hover:border-slate-700/80 transition-all duration-200 hover:shadow-lg hover:shadow-amber-500/5"
                >
                    <div
                        class={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl font-bold text-xs border transition-all ${
                            item.state
                                ? "bg-amber-500/10 text-amber-300 border-amber-500/40 shadow-sm shadow-amber-500/20"
                                : "bg-slate-800/80 text-slate-400 border-slate-700/50"
                        }`}
                    >
                        {item.state ? "ON" : "OFF"}
                    </div>

                    <div class="flex flex-col gap-0.5 min-w-0 flex-1">
                        <div class="flex items-center gap-2 flex-wrap">
                            <span class="text-xs font-mono text-slate-500"
                                >#{item.id}</span
                            >
                            <span class="text-xs text-slate-400 font-medium"
                                >{formatDate(item.created_at)}</span
                            >
                            <span class="text-[11px] text-slate-500"
                                >({getRelativeTime(item.created_at)})</span
                            >
                        </div>

                        {#if item.reason}
                            <p
                                class="text-sm font-medium text-slate-200 flex items-start gap-1.5 mt-0.5 min-w-0"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 shrink-0 text-amber-400/80 mt-0.5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"
                                    />
                                </svg>
                                <span class="wrap-break-word min-w-0"
                                    >"{item.reason}"</span
                                >
                            </p>
                        {:else}
                            <span class="text-xs italic text-slate-600"
                                >No reason specified</span
                            >
                        {/if}
                    </div>
                </div>
            {/each}
        </div>

        <div
            bind:this={sentinelEl}
            class="mt-8 flex justify-center py-4 min-h-12.5"
        >
            {#if loading}
                <div class="flex items-center gap-2 text-sm text-slate-400">
                    <svg
                        class="h-5 w-5 animate-spin text-amber-400"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                        ></circle>
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                    <span>Loading more history...</span>
                </div>
            {:else if !hasMore}
                <div
                    class="text-xs text-slate-500 italic flex items-center gap-2"
                >
                    <span class="h-px w-8 bg-slate-800"></span>
                    <span>End of toggle history</span>
                    <span class="h-px w-8 bg-slate-800"></span>
                </div>
            {/if}
        </div>
    {/if}
</div>
