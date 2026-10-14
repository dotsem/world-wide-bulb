import { describe, it, expect, vi, beforeEach } from 'vitest';
import { WsClient } from './ws.svelte';

class MockWebSocket {
	onopen: (() => void) | null = null;
	onclose: (() => void) | null = null;
	onmessage: ((event: { data: string }) => void) | null = null;
	close = vi.fn();
}

type TestWsClient = {
	handleMessage: (msg: Record<string, unknown>) => void;
};

describe('WsClient', () => {
	let client: WsClient;
	let testClient: TestWsClient;

	beforeEach(() => {
		vi.stubGlobal('WebSocket', MockWebSocket);
		client = new WsClient();
		testClient = client as unknown as TestWsClient;
	});

	it('registers event handler and unsubscribes via returned function', () => {
		const handler = vi.fn();
		const unsub = client.on('test_event', handler);

		testClient.handleMessage({ type: 'test_event', data: 'hello' });
		expect(handler).toHaveBeenCalledWith({ type: 'test_event', data: 'hello' });

		unsub();
		testClient.handleMessage({ type: 'test_event', data: 'world' });
		expect(handler).toHaveBeenCalledTimes(1);
	});

	it('routes state_changed events when state property is defined', () => {
		const stateHandler = vi.fn();
		client.onStateChange(stateHandler);

		testClient.handleMessage({ state: true, id: 10 });
		expect(stateHandler).toHaveBeenCalledWith({ state: true, id: 10 });
	});

	it('routes reason_updated events via onReasonUpdate', () => {
		const reasonHandler = vi.fn();
		client.onReasonUpdate(reasonHandler);

		testClient.handleMessage({ type: 'reason_updated', toggle_id: 5, reason: 'Test' });
		expect(reasonHandler).toHaveBeenCalledWith({
			type: 'reason_updated',
			toggle_id: 5,
			reason: 'Test'
		});
	});

	it('routes viewer_count events via onViewerCount', () => {
		const countHandler = vi.fn();
		client.onViewerCount(countHandler);

		testClient.handleMessage({ type: 'viewer_count', count: 99 });
		expect(countHandler).toHaveBeenCalledWith({ type: 'viewer_count', count: 99 });
	});
});
