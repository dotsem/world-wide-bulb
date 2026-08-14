
export interface ToggleResponse {
    state: boolean;
    created_at: string;
    cooldown_ms: number;
}

export interface StateResponse {
    state: boolean;
    cooldown_ms: number;
}