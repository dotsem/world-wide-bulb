export interface ToggleResponse {
	id: string;
	state: boolean;
	created_at: string;
	cooldown_ms: number;
}

export interface StateResponse {
	state: boolean;
	cooldown_ms: number;
	viewers?: number;
}

export interface HistoryItem {
	id: number;
	state: boolean;
	reason?: string;
	created_at: string;
}

export interface HistoryResponse {
	toggles: HistoryItem[];
	next_cursor?: number;
	has_more: boolean;
}
