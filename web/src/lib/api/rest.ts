import type { StateResponse, ToggleResponse } from "$lib/types/rest.types";

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:5000';


class RestApi {

    async toggle(reason?: string): Promise<ToggleResponse> {
        const res = await fetch(`${API_BASE}/api/v1/toggle`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ reason }),
        });
        if (!res.ok) {
            throw new Error('Failed to toggle bulb');
        }
        return res.json();
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

    async history() {

    }
}

export const restApi = new RestApi();