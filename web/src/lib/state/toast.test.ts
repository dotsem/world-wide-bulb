import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { toast, toastState } from './toast.svelte';

describe('toastState', () => {
	beforeEach(() => {
		toastState.clear();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it('adds a toast item with default configuration', () => {
		const id = toast.show('Bulb turned on', { type: 'success' });

		expect(toastState.toasts).toHaveLength(1);
		expect(toastState.toasts[0].id).toBe(id);
		expect(toastState.toasts[0].message).toBe('Bulb turned on');
		expect(toastState.toasts[0].type).toBe('success');
	});

	it('creates specialized toasts via helper methods', () => {
		const errId = toast.error('Error occurred');
		const warnId = toast.warning('Warning emitted');
		const succId = toast.success('Operation succeeded');
		const infoId = toast.info('Notice info');

		expect(toastState.toasts).toHaveLength(4);
		expect(toastState.toasts.find((t) => t.id === errId)?.type).toBe('error');
		expect(toastState.toasts.find((t) => t.id === warnId)?.type).toBe('warning');
		expect(toastState.toasts.find((t) => t.id === succId)?.type).toBe('success');
		expect(toastState.toasts.find((t) => t.id === infoId)?.type).toBe('info');
	});

	it('dismisses a toast by id', () => {
		const id = toast.info('System standby');
		expect(toastState.toasts).toHaveLength(1);

		toast.dismiss(id);
		expect(toastState.toasts).toHaveLength(0);
	});

	it('auto-dismisses after duration has passed', () => {
		vi.useFakeTimers();

		toast.show('Temporary toast', { duration: 3000 });
		expect(toastState.toasts).toHaveLength(1);

		vi.advanceTimersByTime(2999);
		expect(toastState.toasts).toHaveLength(1);

		vi.advanceTimersByTime(1);
		expect(toastState.toasts).toHaveLength(0);
	});

	it('enforces maximum toast capacity by evicting oldest', () => {
		for (let i = 0; i < 7; i++) {
			toast.show(`Message ${i}`);
		}
		expect(toastState.toasts.length).toBe(5);
		expect(toastState.toasts[0].message).toBe('Message 2');
		expect(toastState.toasts[4].message).toBe('Message 6');
	});

	it('clears all active toasts and timers', () => {
		toast.info('Item 1');
		toast.error('Item 2');
		expect(toastState.toasts).toHaveLength(2);

		toastState.clear();
		expect(toastState.toasts).toHaveLength(0);
	});
});
