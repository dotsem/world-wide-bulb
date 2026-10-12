import { test, expect } from '@playwright/test';

test.describe('CooldownNotice UI Component', () => {
	test.beforeEach(async ({ page }) => {
		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});
	});

	test('is hidden when cooldown is inactive', async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 0 })
			});
		});

		await page.goto('/');
		const notice = page.locator('div[role="status"]');
		await expect(notice).toHaveCount(0);
	});

	test('displays countdown when cooldown is active', async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 8000 })
			});
		});

		await page.goto('/');
		const notice = page.locator('div[role="status"]');
		await expect(notice).toBeVisible();
		await expect(notice).toContainText('Cooldown active:');
		await expect(notice).toContainText('8s');
	});

	test('pulses indicator when clicking bulb during cooldown', async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 6000 })
			});
		});

		await page.goto('/');
		const notice = page.locator('div[role="status"]');
		await expect(notice).toBeVisible();

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click({ force: true });

		await expect(notice).toHaveClass(/scale-110/);
	});
});
