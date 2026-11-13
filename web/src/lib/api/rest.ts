import { solveTurnstile } from './turnstile';
import type { HistoryResponse, StateResponse, ToggleResponse } from '$lib/types/rest.types';

const API_BASE = import.meta.env.VITE_API_BASE ?? '';

class RestApi {
	private verifying?: Promise<void>;

	/** Passes Turnstile while the visitor looks around, so clicking the bulb does not wait on it. */
	verifyInBackground(sitekey: string) {
		this.verifying ??= solveTurnstile(sitekey, 'interaction-only')
			.then((token) =>
				fetch(`${API_BASE}/api/v1/verify`, {
					method: 'POST',
					credentials: 'include',
					headers: { 'X-Turnstile-Token': token }
				})
			)
			.then(
				() => undefined,
				() => undefined
			);
	}

	async toggle(turnstileToken?: string): Promise<ToggleResponse> {
		await this.verifying;
		const res = await fetch(`${API_BASE}/api/v1/toggle`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				...(turnstileToken && { 'X-Turnstile-Token': turnstileToken })
			}
		});
		const data = await res.json();
		if (res.status === 403 && data.turnstile_site_key && !turnstileToken) {
			return this.toggle(await solveTurnstile(data.turnstile_site_key));
		}
		if (!res.ok) {
			throw { message: data.error || 'Failed to toggle bulb', cooldown_ms: data.cooldown_ms || 0 };
		}
		return data;
	}

	async state(): Promise<StateResponse> {
		const res = await fetch(`${API_BASE}/api/v1/state`);
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
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ id, reason })
		});
		if (!res.ok) {
			const data = await res.json();
			throw new Error(data.error || 'Failed to submit reason');
		}
	}

	async history(limit = 20, before?: number): Promise<HistoryResponse> {
		const params = new URLSearchParams();
		if (limit) params.set('limit', limit.toString());
		if (before) params.set('before', before.toString());

		const query = params.toString() ? `?${params.toString()}` : '';
		const res = await fetch(`${API_BASE}/api/v1/history${query}`);
		if (!res.ok) {
			throw new Error('Failed to get bulb history');
		}
		return res.json();
	}
}

export const restApi = new RestApi();
