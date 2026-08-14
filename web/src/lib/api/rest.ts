import type { StateResponse, ToggleResponse } from "$lib/types/rest.types";

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:5000';


class RestApi {

    async toggle(): Promise<ToggleResponse> {
        const res = await fetch(`${API_BASE}/api/v1/toggle`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        const data = await res.json();
        if (!res.ok) {
            throw { message: data.error || 'Failed to toggle bulb', cooldown_ms: data.cooldown_ms || 0 };
        }
        return data;
    }

    async state(): Promise<StateResponse> {
        const res = await fetch(`${API_BASE}/api/v1/state`, {
            method: 'GET',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!res.ok) {
            throw new Error('Failed to get bulb state');
        }
        return res.json();
    }

    async postReason(id: string, reason: string): Promise<void> {
        const res = await fetch(`${API_BASE}/api/v1/reason`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id, reason }),
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error || 'Failed to submit reason');
        }
    }

    async history() {

    }
}

export const restApi = new RestApi();