import { restApi } from '$lib/api/rest';
import { wsClient, type StateChangedEvent, type ReasonUpdatedEvent } from '$lib/api/ws.svelte';
import type { HistoryItem } from '$lib/types/rest.types';

export interface FormattedHistoryItem extends HistoryItem {
	formattedDate: string;
	relativeTime: string;
}

function formatDate(dateStr: string): string {
	try {
		return new Date(dateStr).toLocaleString(undefined, {
			dateStyle: 'medium',
			timeStyle: 'short'
		});
	} catch {
		return dateStr;
	}
}

function getRelativeTime(dateStr: string): string {
	try {
		const elapsedSec = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
		if (elapsedSec < 5) return 'just now';
		if (elapsedSec < 60) return `${elapsedSec}s ago`;
		const elapsedMin = Math.floor(elapsedSec / 60);
		if (elapsedMin < 60) return `${elapsedMin}m ago`;
		const elapsedHours = Math.floor(elapsedMin / 60);
		if (elapsedHours < 24) return `${elapsedHours}h ago`;
		return `${Math.floor(elapsedHours / 24)}d ago`;
	} catch {
		return '';
	}
}

function prepareItem(item: HistoryItem): FormattedHistoryItem {
	return {
		...item,
		formattedDate: formatDate(item.created_at),
		relativeTime: getRelativeTime(item.created_at)
	};
}

class HistoryState {
	items = $state<FormattedHistoryItem[]>([]);
	nextCursor = $state<number | undefined>(undefined);
	hasMore = $state(true);
	loading = $state(false);
	initialLoading = $state(true);
	error = $state<string | null>(null);

	constructor() {
		wsClient.onStateChange((msg) => this.handleStateChange(msg));
		wsClient.onReasonUpdate((msg) => this.handleReasonUpdate(msg));
	}

	get totalToggles(): number {
		return this.items[0]?.id ?? 0;
	}

	private handleStateChange(msg: StateChangedEvent) {
		if (msg.id === undefined) return;

		const existingIdx = this.items.findIndex((t) => t.id === msg.id);
		if (existingIdx !== -1) {
			this.items[existingIdx] = {
				...this.items[existingIdx],
				state: msg.state,
				...(msg.reason ? { reason: msg.reason } : {})
			};
			return;
		}

		const newItem: HistoryItem = {
			id: msg.id,
			state: msg.state,
			reason: msg.reason || '',
			created_at: msg.created_at || new Date().toISOString()
		};

		if (this.items.length > 0 && msg.id > this.items[0].id + 1) {
			this.refresh();
			return;
		}

		this.items = [prepareItem(newItem), ...this.items];
	}

	private handleReasonUpdate(msg: ReasonUpdatedEvent) {
		const targetId = msg.toggle_id ?? (typeof msg.id === 'number' ? msg.id : undefined);
		if (targetId === undefined) return;
		this.items = this.items.map((t) => (t.id === targetId ? { ...t, reason: msg.reason } : t));
	}

	refresh = async () => {
		this.loading = true;
		this.error = null;

		try {
			const res = await restApi.history(20);
			this.items = res.toggles.map(prepareItem);
			this.nextCursor = res.next_cursor;
			this.hasMore = res.has_more;
		} catch (err: any) {
			this.error = err?.message || 'Failed to load history items';
		} finally {
			this.loading = false;
			this.initialLoading = false;
		}
	};

	loadMore = async () => {
		if (this.loading || !this.hasMore) return;

		this.loading = true;
		this.error = null;

		try {
			const res = await restApi.history(20, this.nextCursor);
			const formatted = res.toggles.map(prepareItem);
			const existingIds = new Set(this.items.map((i) => i.id));
			const newItems = formatted.filter((i) => !existingIds.has(i.id));

			this.items = [...this.items, ...newItems];
			this.nextCursor = res.next_cursor;
			this.hasMore = res.has_more;
		} catch (err: any) {
			this.error = err?.message || 'Failed to load history items';
		} finally {
			this.loading = false;
			this.initialLoading = false;
		}
	};
}

export const historyState = new HistoryState();
