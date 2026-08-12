import { restApi } from "$lib/api/rest";
import { wsClient } from "$lib/api/ws.svelte";


class BulbState {
    isOn = $state(false);
    history = $state<Array<{ id: number; state: boolean; reason?: string }>>([]);
    cooldownUntil = $state<number | null>(null);

    async init() {
        const res = await restApi.state();
        this.isOn = res.state;
    }

    constructor() {
        wsClient.on('state_changed', (msg) => {
            this.isOn = msg.state;
        });
    }

    toggle = async (reason?: unknown) => {
        const cleanReason = typeof reason === 'string' ? reason : undefined;
        const res = await restApi.toggle(cleanReason);
        this.isOn = res.state;
    }
}

export const bulbState = new BulbState()
bulbState.init();