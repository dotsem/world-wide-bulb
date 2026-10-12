import { test, expect } from '@playwright/test';

test.describe('Index Route and Theme Transitions', () => {
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

	test('renders document title and initial off theme', async ({ page }) => {
		await expect(page).toHaveTitle('World Wide Bulb');

		const html = page.locator('html');
		await expect(html).toHaveAttribute('data-theme', 'off');

		const favicon = page.locator('link[rel="icon"]');
		await expect(favicon).toHaveAttribute('href', /bulb_off/);

		const layoutContainer = page.locator('div[data-theme]');
		await expect(layoutContainer).toHaveAttribute('data-theme', 'off');
	});

	test('switches theme, favicon, and background color upon toggling bulb', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ id: 'toggle-theme-1', state: true, cooldown_ms: 0 })
			});
		});

		const html = page.locator('html');
		const bulbButton = page.locator('button[aria-label="Turn on lightbulb"]');
		await bulbButton.click();

		await expect(html).toHaveAttribute('data-theme', 'on');

		const favicon = page.locator('link[rel="icon"]');
		await expect(favicon).toHaveAttribute('href', /bulb_on/);

		const layoutContainer = page.locator('div[data-theme]');
		await expect(layoutContainer).toHaveAttribute('data-theme', 'on');

		await expect(page.locator('body')).toHaveCSS('background-color', 'rgb(28, 24, 22)');
	});
});
