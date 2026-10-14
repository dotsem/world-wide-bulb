import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BulbState } from './bulb.svelte';
import { restApi } from '$lib/api/rest';
import { wsClient } from '$lib/api/ws.svelte';
import { toast } from '$lib/state/toast.svelte';

vi.mock('$lib/api/rest', () => ({
	restApi: {
		state: vi.fn().mockResolvedValue({ state: false, viewers: 1, cooldown_ms: 0 }),
		toggle: vi.fn(),
		postReason: vi.fn()
	}
}));

vi.mock('$lib/api/ws.svelte', () => {
	const handlers: Record<string, ((...args: unknown[]) => void)[]> = {};
	return {
		wsClient: {
			on: vi.fn((event, handler) => {
				handlers[event] = handlers[event] || [];
				handlers[event].push(handler);
				return () => {};
			}),
			onStateChange: vi.fn((handler) => {
				handlers['state_changed'] = handlers['state_changed'] || [];
				handlers['state_changed'].push(handler);
				return () => {};
			}),
			onViewerCount: vi.fn((handler) => {
				handlers['viewer_count'] = handlers['viewer_count'] || [];
				handlers['viewer_count'].push(handler);
			}),
			emitMock: (event: string, data: unknown) => {
				handlers[event]?.forEach((h) => h(data));
			}
		}
	};
});

vi.mock('$lib/state/toast.svelte', () => ({
	toast: {
		error: vi.fn(),
		info: vi.fn(),
		success: vi.fn(),
		show: vi.fn()
	}
}));

describe('BulbState', () => {
	let state: BulbState;

	beforeEach(() => {
		vi.clearAllMocks();
		state = new BulbState();
	});

	it('initializes with default values', () => {
		expect(state.isOn).toBe(false);
		expect(state.viewers).toBe(1);
		expect(state.remainingMs).toBe(0);
		expect(state.isCooldownActive).toBe(false);
		expect(state.remainingSeconds).toBe(0);
		expect(state.showReasonPrompt).toBe(false);
	});

	it('loads remote state on init()', async () => {
		vi.mocked(restApi.state).mockResolvedValueOnce({
			state: true,
			viewers: 42,
			cooldown_ms: 3000
		});

		await state.init();

		expect(state.isOn).toBe(true);
		expect(state.viewers).toBe(42);
		expect(state.isCooldownActive).toBe(true);
		expect(state.remainingSeconds).toBe(3);
	});

	type MockWsClient = { emitMock: (event: string, data: unknown) => void };

	it('updates state upon receiving websocket events and applies 1s debounce', () => {
		const mockEmitter = (wsClient as unknown as MockWsClient).emitMock;

		mockEmitter('state_changed', { state: true });
		expect(state.isOn).toBe(true);
		expect(state.isCooldownActive).toBe(true);
		expect(state.remainingSeconds).toBe(1);

		mockEmitter('viewer_count', { count: 15 });
		expect(state.viewers).toBe(15);
	});

	it('does not reduce longer personal cooldown upon receiving state_changed', () => {
		const mockEmitter = (wsClient as unknown as MockWsClient).emitMock;
		state.setCooldown(8000);

		mockEmitter('state_changed', { state: true });
		expect(state.remainingSeconds).toBe(8);
	});

	it('handles cooldown countdown and cancellation', () => {
		vi.useFakeTimers();

		state.setCooldown(2000);
		expect(state.isCooldownActive).toBe(true);
		expect(state.remainingSeconds).toBe(2);

		vi.advanceTimersByTime(1100);
		expect(state.remainingSeconds).toBe(1);

		vi.advanceTimersByTime(1000);
		expect(state.isCooldownActive).toBe(false);
		expect(state.remainingMs).toBe(0);

		state.setCooldown(5000);
		state.setCooldown(0);
		expect(state.isCooldownActive).toBe(false);
		expect(state.cooldownUntil).toBeNull();

		vi.useRealTimers();
	});

	it('toggles successfully and prepares reason prompt', async () => {
		vi.mocked(restApi.toggle).mockResolvedValueOnce({
			id: 'toggle-123',
			state: true,
			created_at: '2026-10-12T12:00:00Z',
			cooldown_ms: 10000
		});

		await state.toggle();

		expect(state.isOn).toBe(true);
		expect(state.lastToggleId).toBe('toggle-123');
		expect(state.lastActionState).toBe(true);
		expect(state.showReasonPrompt).toBe(true);
		expect(state.isCooldownActive).toBe(true);
	});

	it('pulses cooldown indicator and prevents toggle when active', async () => {
		state.setCooldown(5000);
		const initialPulse = state.cooldownPulseKey;

		await state.toggle();

		expect(state.cooldownPulseKey).toBe(initialPulse + 1);
		expect(restApi.toggle).not.toHaveBeenCalled();
	});

	it('reverts state and displays toast when toggle API fails', async () => {
		state.isOn = false;
		vi.mocked(restApi.toggle).mockRejectedValueOnce({
			message: 'Rate limit hit',
			cooldown_ms: 4000
		});

		await state.toggle();

		expect(state.isOn).toBe(false);
		expect(state.isCooldownActive).toBe(true);
		expect(toast.error).toHaveBeenCalledWith('Rate limit hit');
	});

	it('submits reason and cleans up prompt state', async () => {
		state.lastToggleId = 'toggle-456';
		state.showReasonPrompt = true;
		state.lastActionState = true;
		vi.mocked(restApi.postReason).mockResolvedValueOnce(undefined);

		await state.submitReason('Working late');

		expect(restApi.postReason).toHaveBeenCalledWith('toggle-456', 'Working late');
		expect(state.showReasonPrompt).toBe(false);
		expect(state.lastToggleId).toBeNull();
		expect(state.lastActionState).toBeNull();
	});

	it('skips submitting reason when no active toggle id', async () => {
		state.lastToggleId = null;
		await state.submitReason('Any reason');
		expect(restApi.postReason).not.toHaveBeenCalled();
	});

	it('dismisses reason prompt', () => {
		state.lastToggleId = 'toggle-789';
		state.lastActionState = false;
		state.showReasonPrompt = true;

		state.dismissReason();

		expect(state.showReasonPrompt).toBe(false);
		expect(state.lastToggleId).toBeNull();
		expect(state.lastActionState).toBeNull();
	});
});
