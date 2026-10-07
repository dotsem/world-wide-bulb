export type ToastType = 'error' | 'warning' | 'info' | 'success';

export interface ToastItem {
	id: string;
	message: string;
	type: ToastType;
	duration: number;
	createdAt: number;
}

const MAX_TOASTS = 5;
const DEFAULT_DURATION = 5000;

class ToastState {
	toasts = $state<ToastItem[]>([]);
	private timers = new Map<string, ReturnType<typeof setTimeout>>();

	show(message: string, options?: { type?: ToastType; duration?: number }) {
		const id =
			typeof crypto !== 'undefined' && crypto.randomUUID
				? crypto.randomUUID()
				: Math.random().toString(36).substring(2, 9);
		const type = options?.type ?? 'error';
		const duration = options?.duration ?? DEFAULT_DURATION;

		const item: ToastItem = {
			id,
			message,
			type,
			duration,
			createdAt: Date.now()
		};

		if (this.toasts.length >= MAX_TOASTS) {
			const oldest = this.toasts[0];
			if (oldest) {
				this.dismiss(oldest.id);
			}
		}

		this.toasts.push(item);

		if (duration > 0) {
			const timer = setTimeout(() => {
				this.dismiss(id);
			}, duration);
			this.timers.set(id, timer);
		}

		return id;
	}

	error(message: string, duration?: number) {
		return this.show(message, { type: 'error', duration });
	}

	warning(message: string, duration?: number) {
		return this.show(message, { type: 'warning', duration });
	}

	info(message: string, duration?: number) {
		return this.show(message, { type: 'info', duration });
	}

	success(message: string, duration?: number) {
		return this.show(message, { type: 'success', duration });
	}

	dismiss(id: string) {
		const timer = this.timers.get(id);
		if (timer) {
			clearTimeout(timer);
			this.timers.delete(id);
		}
		this.toasts = this.toasts.filter((t) => t.id !== id);
	}

	clear() {
		for (const timer of this.timers.values()) {
			clearTimeout(timer);
		}
		this.timers.clear();
		this.toasts = [];
	}
}

export const toastState = new ToastState();
export const toast = toastState;
