import { describe, it, expect, vi, beforeEach } from 'vitest';
import { WsClient } from './ws.svelte';

class MockWebSocket {
	onopen: (() => void) | null = null;
	onclose: (() => void) | null = null;
	onmessage: ((event: any) => void) | null = null;
	close = vi.fn();
}

describe('WsClient', () => {
	let client: WsClient;

	beforeEach(() => {
		vi.stubGlobal('WebSocket', MockWebSocket);
		client = new WsClient();
	});

	it('registers event handler and unsubscribes via returned function', () => {
		const handler = vi.fn();
		const unsub = client.on('test_event', handler);

		(client as any).handleMessage({ type: 'test_event', data: 'hello' });
		expect(handler).toHaveBeenCalledWith({ type: 'test_event', data: 'hello' });

		unsub();
		(client as any).handleMessage({ type: 'test_event', data: 'world' });
		expect(handler).toHaveBeenCalledTimes(1);
	});

	it('routes state_changed events when state property is defined', () => {
		const stateHandler = vi.fn();
		client.onStateChange(stateHandler);

		(client as any).handleMessage({ state: true, id: 10 });
		expect(stateHandler).toHaveBeenCalledWith({ state: true, id: 10 });
	});

	it('routes reason_updated events via onReasonUpdate', () => {
		const reasonHandler = vi.fn();
		client.onReasonUpdate(reasonHandler);

		(client as any).handleMessage({ type: 'reason_updated', toggle_id: 5, reason: 'Test' });
		expect(reasonHandler).toHaveBeenCalledWith({
			type: 'reason_updated',
			toggle_id: 5,
			reason: 'Test'
		});
	});

	it('routes viewer_count events via onViewerCount', () => {
		const countHandler = vi.fn();
		client.onViewerCount(countHandler);

		(client as any).handleMessage({ type: 'viewer_count', count: 99 });
		expect(countHandler).toHaveBeenCalledWith({ type: 'viewer_count', count: 99 });
	});
});
