import { restApi } from "$lib/api/rest";
import { wsClient } from "$lib/api/ws.svelte";


class BulbState {
    isOn = $state(false);
    history = $state<Array<{ id: number; state: boolean; reason?: string }>>([]);
    cooldownUntil = $state<number | null>(null);
    remainingMs = $state(0);
    private timer?: ReturnType<typeof setInterval>;

    constructor() {
        wsClient.on('state_changed', (msg) => {
            this.isOn = msg.state;
        });
    }

    remainingSeconds = $derived(Math.ceil(this.remainingMs / 1000));
    isCooldownActive = $derived(this.remainingMs > 0);

    setCooldown(ms: number) {
        if (this.timer) {
            clearInterval(this.timer);
            this.timer = undefined;
        }

        if (ms <= 0) {
            this.cooldownUntil = null;
            this.remainingMs = 0;
            return;
        }

        const target = Date.now() + ms;
        this.cooldownUntil = target;
        this.remainingMs = ms;

        this.timer = setInterval(() => {
            const left = target - Date.now();
            if (left <= 0) {
                this.remainingMs = 0;
                this.cooldownUntil = null;
                if (this.timer) {
                    clearInterval(this.timer);
                    this.timer = undefined;
                }
            } else {
                this.remainingMs = left;
            }
        }, 100);
    }

    async init() {
        const res = await restApi.state();
        this.isOn = res.state;
        this.setCooldown(res.cooldown_ms);
    }

    toggle = async (reason?: unknown) => {
        if (this.isCooldownActive) return;
        const cleanReason = typeof reason === 'string' ? reason : undefined;
        try {
            const res = await restApi.toggle(cleanReason);
            this.isOn = res.state;
            this.setCooldown(res.cooldown_ms);
        } catch (err: any) {
            if (err?.cooldown_ms) {
                this.setCooldown(err.cooldown_ms);
            }
        }
    }
}

export const bulbState = new BulbState();
bulbState.init();