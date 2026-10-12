import { test, expect } from '@playwright/test';

test.describe('ReasonPrompt and FloatingInput UI Components', () => {
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

	test('surfaces prompt upon toggle and submits reason with feedback', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ id: 'toggle-99', state: true, cooldown_ms: 10000 })
			});
		});

		let reasonSubmitted = false;
		await page.route('**/api/v1/reason', async (route) => {
			reasonSubmitted = true;
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ success: true })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click();

		const dialog = page.locator('div[role="dialog"]');
		await expect(dialog).toBeVisible();
		await expect(dialog).toContainText('Why did you decide to turn on the lamp?');

		const submitBtn = dialog.locator('button[type="submit"]');
		await expect(submitBtn).toBeDisabled();

		const input = dialog.locator('input');
		await input.fill('Late night study session');
		await expect(dialog).toContainText('24/60');
		await expect(submitBtn).toBeEnabled();

		const reasonPromise = page.waitForRequest((req) => req.url().includes('/api/v1/reason'));
		await submitBtn.click();
		await reasonPromise;
		expect(reasonSubmitted).toBe(true);

		await expect(dialog).toHaveCount(0);
	});

	test('dismisses prompt when skip is clicked', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ id: 'toggle-100', state: true, cooldown_ms: 10000 })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click();

		const dialog = page.locator('div[role="dialog"]');
		await expect(dialog).toBeVisible();

		const skipBtn = dialog.locator('button', { hasText: 'Skip' });
		await skipBtn.click();

		await expect(dialog).toHaveCount(0);
	});
});
