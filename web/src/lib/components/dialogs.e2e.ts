import { test, expect } from '@playwright/test';

test.describe('Dialog Components (InfoDialog, HookDialog, Dialog UI)', () => {
	test.beforeEach(async ({ page }) => {
		await page.route('**/api/v1/state', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ state: false, viewers: 1, cooldown_ms: 0 })
			});
		});
	});

	test('auto-opens InfoDialog on first visit when intro is not seen', async ({ page }) => {
		await page.goto('/');

		const dialogTitle = page.locator('h2', { hasText: 'World Wide Bulb' });
		await expect(dialogTitle).toBeVisible();

		const coolBtn = page.locator('button', { hasText: 'Cool!' });
		await coolBtn.click();

		await expect(dialogTitle).toHaveCount(0);
		const seenIntro = await page.evaluate(() => localStorage.getItem('wwb_seen_intro'));
		expect(seenIntro).toBe('true');
	});

	test('closes modal upon pressing Escape key', async ({ page }) => {
		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});
		await page.goto('/');

		const infoBtn = page.locator('button[aria-label="About World Wide Bulb"]');
		await infoBtn.click();

		const dialog = page.locator('div[role="dialog"]');
		await expect(dialog).toBeVisible();

		await page.keyboard.press('Escape');
		await expect(dialog).toHaveCount(0);
	});

	test('closes modal upon clicking backdrop', async ({ page }) => {
		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});
		await page.goto('/');

		const infoBtn = page.locator('button[aria-label="About World Wide Bulb"]');
		await infoBtn.click();

		const backdrop = page.locator('div[role="presentation"]');
		await expect(backdrop).toBeVisible();

		await backdrop.click({ position: { x: 10, y: 10 } });
		await expect(page.locator('div[role="dialog"]')).toHaveCount(0);
	});

	test('displays HookDialog stream command and navigates to docs', async ({ page }) => {
		await page.addInitScript(() => {
			localStorage.setItem('wwb_seen_intro', 'true');
		});
		await page.goto('/');

		const hookBtn = page.locator('button[aria-label="Hook into the lamp"]');
		await hookBtn.click();

		const hookDialog = page.locator('div[role="dialog"]');
		await expect(hookDialog).toBeVisible();
		await expect(hookDialog).toContainText('curl -N');

		const docsLink = hookDialog.locator('a[href^="/docs"]');
		await expect(docsLink).toBeVisible();
		await docsLink.click();

		await expect(page).toHaveURL(/\/docs/);
	});
});
