import { restApi } from '$lib/api/rest';
import { wsClient } from '$lib/api/ws.svelte';
import type { HistoryItem } from '$lib/types/rest.types';

export interface FormattedHistoryItem extends HistoryItem {
	formattedDate: string;
	relativeTime: string;
}

const MAX_ITEMS = 100;

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

	get totalToggles(): number {
		return this.items[0]?.id ?? 0;
	}

	loadMore = async () => {
		if (this.loading || !this.hasMore) return;

		this.loading = true;
		this.error = null;

		try {
			const res = await restApi.history(20, this.nextCursor);
			const formatted = res.toggles.map(prepareItem);
			this.items = [...this.items, ...formatted].slice(0, MAX_ITEMS);
			this.nextCursor = res.next_cursor;
			this.hasMore = res.has_more && this.items.length < MAX_ITEMS;
		} catch (err: any) {
			this.error = err?.message || 'Failed to load history items';
		} finally {
			this.loading = false;
			this.initialLoading = false;
		}
	};

	initWsListeners = (): (() => void) => {
		const unsubState = wsClient.onStateChange((msg) => {
			if (msg.id === undefined || this.items.some((t) => t.id === msg.id)) return;
			const newItem: HistoryItem = {
				id: msg.id,
				state: msg.state,
				reason: msg.reason || '',
				created_at: msg.created_at || new Date().toISOString()
			};
			this.items = [prepareItem(newItem), ...this.items].slice(0, MAX_ITEMS);
		});

		const unsubReason = wsClient.onReasonUpdate((msg) => {
			const targetId = msg.toggle_id ?? (typeof msg.id === 'number' ? msg.id : undefined);
			if (targetId === undefined) return;
			this.items = this.items.map((t) => (t.id === targetId ? { ...t, reason: msg.reason } : t));
		});

		return () => {
			unsubState();
			unsubReason();
		};
	};
}

export const historyState = new HistoryState();
