type Appearance = 'always' | 'interaction-only';

interface Turnstile {
	render(
		container: HTMLElement,
		options: {
			sitekey: string;
			appearance: Appearance;
			callback: (token: string) => void;
			'error-callback': () => void;
			'timeout-callback': () => void;
		}
	): string;
	remove(widgetId: string): void;
}

declare global {
	interface Window {
		turnstile?: Turnstile;
	}
}

let script: Promise<void> | undefined;

function loadScript(): Promise<void> {
	script ??= new Promise((resolve, reject) => {
		const el = document.createElement('script');
		el.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit';
		el.onload = () => resolve();
		el.onerror = () => {
			script = undefined;
			reject(new Error('Failed to load human verification'));
		};
		document.head.append(el);
	});
	return script;
}

/** Runs a Cloudflare Turnstile challenge and resolves with the token to send to the backend. */
export async function solveTurnstile(
	sitekey: string,
	appearance: Appearance = 'always'
): Promise<string> {
	await loadScript();
	const turnstile = window.turnstile;
	if (!turnstile) throw new Error('Failed to load human verification');

	const box = document.createElement('div');
	box.style.cssText = 'position:fixed;bottom:6rem;left:50%;transform:translateX(-50%);z-index:9999';
	document.body.append(box);
	let widgetId: string | undefined;
	try {
		return await new Promise<string>((resolve, reject) => {
			widgetId = turnstile.render(box, {
				sitekey,
				appearance,
				callback: resolve,
				'error-callback': () => reject(new Error('Human verification failed')),
				'timeout-callback': () => reject(new Error('Human verification timed out'))
			});
		});
	} finally {
		if (widgetId) turnstile.remove(widgetId);
		box.remove();
	}
}
