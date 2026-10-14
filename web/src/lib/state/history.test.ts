import { describe, it, expect, vi, beforeEach } from 'vitest';
import { HistoryState } from './history.svelte';
import { restApi } from '$lib/api/rest';
import { wsClient } from '$lib/api/ws.svelte';

vi.mock('$lib/api/rest', () => ({
	restApi: {
		history: vi.fn()
	}
}));

type MockHistoryWsClient = {
	emitStateChange: (msg: {
		id: number;
		state: boolean;
		reason?: string;
		created_at?: string;
	}) => void;
	emitReasonUpdate: (msg: { type: string; toggle_id: number; reason: string }) => void;
};

vi.mock('$lib/api/ws.svelte', () => {
	let stateChangeHandler:
		| ((msg: { id: number; state: boolean; reason?: string; created_at?: string }) => void)
		| undefined;
	let reasonUpdateHandler:
		((msg: { type: string; toggle_id: number; reason: string }) => void) | undefined;

	return {
		wsClient: {
			onStateChange: vi.fn((handler) => {
				stateChangeHandler = handler;
				return () => {};
			}),
			onReasonUpdate: vi.fn((handler) => {
				reasonUpdateHandler = handler;
				return () => {};
			}),
			emitStateChange: (msg: {
				id: number;
				state: boolean;
				reason?: string;
				created_at?: string;
			}) => stateChangeHandler?.(msg),
			emitReasonUpdate: (msg: { type: string; toggle_id: number; reason: string }) =>
				reasonUpdateHandler?.(msg)
		}
	};
});

describe('HistoryState', () => {
	let state: HistoryState;

	beforeEach(() => {
		vi.clearAllMocks();
		state = new HistoryState();
	});

	it('initializes with default values', () => {
		expect(state.items).toEqual([]);
		expect(state.totalToggles).toBe(0);
		expect(state.loading).toBe(false);
		expect(state.initialLoading).toBe(true);
		expect(state.hasMore).toBe(true);
		expect(state.error).toBeNull();
	});

	it('refreshes items from restApi.history', async () => {
		vi.mocked(restApi.history).mockResolvedValueOnce({
			toggles: [
				{ id: 2, state: true, reason: 'Reading', created_at: '2026-10-12T12:00:00Z' },
				{ id: 1, state: false, reason: '', created_at: '2026-10-12T11:00:00Z' }
			],
			next_cursor: 1,
			has_more: true
		});

		await state.refresh();

		expect(state.items).toHaveLength(2);
		expect(state.totalToggles).toBe(2);
		expect(state.items[0].relativeTime).toBeDefined();
		expect(state.items[0].formattedDate).toBeDefined();
		expect(state.nextCursor).toBe(1);
		expect(state.hasMore).toBe(true);
		expect(state.initialLoading).toBe(false);
		expect(state.loading).toBe(false);
	});

	it('handles refresh error', async () => {
		vi.mocked(restApi.history).mockRejectedValueOnce(new Error('Network failure'));

		await state.refresh();

		expect(state.error).toBe('Network failure');
		expect(state.initialLoading).toBe(false);
		expect(state.loading).toBe(false);
	});

	it('loads more items and deduplicates them', async () => {
		state.items = [
			{
				id: 3,
				state: true,
				reason: 'Daylight',
				created_at: '2026-10-12T12:00:00Z',
				formattedDate: 'Oct 12',
				relativeTime: '1h ago'
			}
		];
		state.nextCursor = 2;
		state.hasMore = true;

		vi.mocked(restApi.history).mockResolvedValueOnce({
			toggles: [
				{ id: 3, state: true, reason: 'Daylight', created_at: '2026-10-12T12:00:00Z' },
				{ id: 2, state: false, reason: '', created_at: '2026-10-12T11:00:00Z' }
			],
			next_cursor: undefined,
			has_more: false
		});

		await state.loadMore();

		expect(state.items).toHaveLength(2);
		expect(state.items[1].id).toBe(2);
		expect(state.hasMore).toBe(false);
		expect(state.nextCursor).toBeUndefined();
	});

	it('guards loadMore against concurrent or completed calls', async () => {
		state.loading = true;
		await state.loadMore();
		expect(restApi.history).not.toHaveBeenCalled();

		state.loading = false;
		state.hasMore = false;
		await state.loadMore();
		expect(restApi.history).not.toHaveBeenCalled();
	});

	it('prepends new item on continuous websocket state change', () => {
		state.items = [
			{
				id: 5,
				state: true,
				reason: '',
				created_at: '2026-10-12T10:00:00Z',
				formattedDate: 'Oct 12',
				relativeTime: '1m ago'
			}
		];

		(wsClient as unknown as MockHistoryWsClient).emitStateChange({
			id: 6,
			state: false,
			reason: 'Leaving room',
			created_at: '2026-10-12T10:01:00Z'
		});

		expect(state.items).toHaveLength(2);
		expect(state.items[0].id).toBe(6);
		expect(state.items[0].reason).toBe('Leaving room');
	});

	it('updates existing item state on websocket state change', () => {
		state.items = [
			{
				id: 5,
				state: false,
				reason: '',
				created_at: '2026-10-12T10:00:00Z',
				formattedDate: 'Oct 12',
				relativeTime: '1m ago'
			}
		];

		(wsClient as unknown as MockHistoryWsClient).emitStateChange({
			id: 5,
			state: true,
			reason: 'Correction'
		});

		expect(state.items).toHaveLength(1);
		expect(state.items[0].state).toBe(true);
		expect(state.items[0].reason).toBe('Correction');
	});

	it('triggers refresh when state change gap is detected', () => {
		state.items = [
			{
				id: 5,
				state: true,
				reason: '',
				created_at: '2026-10-12T10:00:00Z',
				formattedDate: 'Oct 12',
				relativeTime: '1m ago'
			}
		];
		const refreshSpy = vi.spyOn(state, 'refresh').mockImplementation(async () => {});

		(wsClient as unknown as MockHistoryWsClient).emitStateChange({
			id: 8,
			state: false
		});

		expect(refreshSpy).toHaveBeenCalled();
	});

	it('updates reason on websocket reason update', () => {
		state.items = [
			{
				id: 10,
				state: true,
				reason: '',
				created_at: '2026-10-12T10:00:00Z',
				formattedDate: 'Oct 12',
				relativeTime: '1m ago'
			}
		];

		(wsClient as unknown as MockHistoryWsClient).emitReasonUpdate({
			type: 'reason_updated',
			toggle_id: 10,
			reason: 'Because darkness'
		});

		expect(state.items[0].reason).toBe('Because darkness');
	});
});
