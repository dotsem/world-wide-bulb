import { describe, it, expect, beforeEach } from 'vitest';
import { toast, toastState } from './toast.svelte';

describe('toastState', () => {
	beforeEach(() => {
		toastState.clear();
	});

	it('adds a toast item with default configuration', () => {
		const id = toast.show('Bulb turned on', { type: 'success' });

		expect(toastState.toasts).toHaveLength(1);
		expect(toastState.toasts[0].id).toBe(id);
		expect(toastState.toasts[0].message).toBe('Bulb turned on');
		expect(toastState.toasts[0].type).toBe('success');
	});

	it('dismisses a toast by id', () => {
		const id = toast.info('System standby');
		expect(toastState.toasts).toHaveLength(1);

		toast.dismiss(id);
		expect(toastState.toasts).toHaveLength(0);
	});

	it('enforces maximum toast capacity', () => {
		for (let i = 0; i < 7; i++) {
			toast.show(`Message ${i}`);
		}
		expect(toastState.toasts.length).toBeLessThanOrEqual(5);
	});
});
