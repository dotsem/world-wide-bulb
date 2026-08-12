
export interface ToggleResponse {
    state: boolean;
    reason?: string;
    created_at: string;
}

export interface StateResponse {
    state: boolean;
}