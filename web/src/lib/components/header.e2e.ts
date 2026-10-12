import { test, expect } from '@playwright/test';

test.describe('Header and ViewerBadge UI Components', () => {
	test.beforeEach(async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 7, cooldown_ms: 0 })
			});
		});

		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});

		await page.goto('/');
	});

	test('renders viewer badge with count and tooltip', async ({ page }) => {
		const badge = page.locator('div[title*="currently watching"]');
		await expect(badge).toBeVisible();
		await expect(badge).toContainText('7');
	});

	test('opens HookDialog when webhook button is clicked', async ({ page }) => {
		const hookBtn = page.locator('button[aria-label="Hook into the lamp"]');
		await expect(hookBtn).toBeVisible();
		await hookBtn.click();

		const dialogTitle = page.locator('h2', { hasText: 'Hook into the Lamp' });
		await expect(dialogTitle).toBeVisible();
	});

	test('opens InfoDialog when about button is clicked', async ({ page }) => {
		const infoBtn = page.locator('button[aria-label="About World Wide Bulb"]');
		await expect(infoBtn).toBeVisible();
		await infoBtn.click();

		const dialogTitle = page.locator('h2', { hasText: 'World Wide Bulb' });
		await expect(dialogTitle).toBeVisible();
	});
});
