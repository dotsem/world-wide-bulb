import { test, expect } from '@playwright/test';

test.describe('Navigation and ToolBar UI Components', () => {
	test.beforeEach(async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 0 })
			});
		});

		await page.route('**/api/v1/history*', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ toggles: [], next_cursor: null, has_more: false })
			});
		});

		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});

		await page.goto('/');
	});

	test('renders navigation toolbar with Bulb, History, and Docs links', async ({ page }) => {
		const nav = page.locator('nav.fixed.bottom-6');
		await expect(nav).toBeVisible();

		const bulbLink = nav.locator('a[href="/"]');
		const historyLink = nav.locator('a[href="/history"]');
		const docsLink = nav.locator('a[href="/docs"]');

		await expect(bulbLink).toBeVisible();
		await expect(historyLink).toBeVisible();
		await expect(docsLink).toBeVisible();
	});

	test('navigates between views when toolbar items are clicked', async ({ page }) => {
		const nav = page.locator('nav.fixed.bottom-6');
		const historyLink = nav.locator('a[href="/history"]');
		const docsLink = nav.locator('a[href="/docs"]');
		const bulbLink = nav.locator('a[href="/"]');

		await historyLink.click();
		await expect(page).toHaveURL(/\/history/);

		await docsLink.click();
		await expect(page).toHaveURL(/\/docs/);

		await bulbLink.click();
		await expect(page).toHaveURL(/\//);
	});
});
