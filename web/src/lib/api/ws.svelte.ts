function getWsUrl(): string {
	if (import.meta.env.VITE_WS_BASE) {
		return `${import.meta.env.VITE_WS_BASE}/ws`;
	}
	if (typeof window === 'undefined') {
		return '';
	}
	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	return `${protocol}//${window.location.host}/ws`;
}

export interface StateChangedEvent {
	id?: number;
	state: boolean;
	reason?: string;
	created_at?: string;
}

export interface ReasonUpdatedEvent {
	type: 'reason_updated';
	id?: string;
	toggle_id?: number;
	reason: string;
}

export interface ViewerCountEvent {
	type: 'viewer_count';
	count: number;
}

class WsClient {
	private socket?: WebSocket;
	private handlers: Record<string, Array<(msg: any) => void>> = {};
	isConnected = $state(false);

	constructor() {
		this.connect();
	}

	private connect() {
		const url = getWsUrl();
		if (!url) {
			return;
		}
		this.socket = new WebSocket(url);
		this.socket.onopen = () => {
			this.isConnected = true;
		};
		this.socket.onmessage = (event) => {
			const msg = JSON.parse(event.data);
			this.handleMessage(msg);
		};
		this.socket.onclose = () => {
			this.isConnected = false;
			setTimeout(() => this.connect(), 1000);
		};
	}

	private handleMessage(msg: any) {
		if (msg.state !== undefined && this.handlers['state_changed']) {
			this.handlers['state_changed'].forEach((handler) => handler(msg));
		}
		if (msg.type && this.handlers[msg.type]) {
			this.handlers[msg.type].forEach((handler) => handler(msg));
		}
	}

	on(eventType: string, handler: (msg: any) => void): () => void {
		if (!this.handlers[eventType]) {
			this.handlers[eventType] = [];
		}
		this.handlers[eventType].push(handler);
		return () => this.off(eventType, handler);
	}

	off(eventType: string, handler: (msg: any) => void) {
		if (!this.handlers[eventType]) return;
		this.handlers[eventType] = this.handlers[eventType].filter((h) => h !== handler);
	}

	onStateChange(handler: (msg: StateChangedEvent) => void): () => void {
		return this.on('state_changed', handler);
	}

	onReasonUpdate(handler: (msg: ReasonUpdatedEvent) => void): () => void {
		return this.on('reason_updated', handler);
	}

	onViewerCount(handler: (msg: ViewerCountEvent) => void): () => void {
		return this.on('viewer_count', handler);
	}
}

export const wsClient = new WsClient();
