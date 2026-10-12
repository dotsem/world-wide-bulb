import { test, expect } from '@playwright/test';

test.describe('Bulb UI Component', () => {
	test.beforeEach(async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 0 })
			});
		});

		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});

		await page.goto('/');
	});

	test('renders initial state with accessible turn-on label', async ({ page }) => {
		const bulbButton = page.locator('button[aria-label="Turn on lightbulb"]');
		await expect(bulbButton).toBeVisible();
		await expect(bulbButton).toHaveAttribute('type', 'button');
		await expect(bulbButton).toHaveAttribute('aria-disabled', 'false');
	});

	test('toggles to on state upon click and updates aria-label', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ id: 'toggle-1', state: true, cooldown_ms: 0 })
			});
		});

		const bulbButton = page.locator('button[aria-label="Turn on lightbulb"]');
		await bulbButton.click();

		const offButton = page.locator('button[aria-label="Turn off lightbulb"]');
		await expect(offButton).toBeVisible();
	});

	test('reflects disabled state and prevents click during cooldown', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ id: 'toggle-2', state: true, cooldown_ms: 10000 })
			});
		});

		let toggleCallCount = 0;
		page.on('request', (req) => {
			if (req.url().includes('/api/v1/toggle')) toggleCallCount += 1;
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click();

		await expect(bulbButton).toHaveAttribute('aria-disabled', 'true');
		await expect(bulbButton).toHaveClass(/cursor-not-allowed/);

		// Click again while disabled
		await bulbButton.click({ force: true });
		expect(toggleCallCount).toBe(1);
	});
});
