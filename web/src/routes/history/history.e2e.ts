import { test, expect } from '@playwright/test';

test.describe('History Page and HistoryTile UI Components', () => {
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
	});

	test('renders empty state when no toggles exist', async ({ page }) => {
		await page.route('**/api/v1/history*', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({ toggles: [], next_cursor: null, has_more: false })
			});
		});

		await page.goto('/history');

		await expect(page.locator('h1', { hasText: 'Toggle History' })).toBeVisible();
		await expect(page.locator('text=No one toggled the bulb, I guess...')).toBeVisible();
	});

	test('renders history tiles with milestone badge and user reasons', async ({ page }) => {
		await page.route('**/api/v1/history*', async (route) => {
			await route.fulfill({
				status: 200,
				contentType: 'application/json',
				body: JSON.stringify({
					toggles: [
						{
							id: 100,
							state: true,
							reason: 'Centennial celebration!',
							created_at: new Date().toISOString()
						},
						{
							id: 99,
							state: false,
							reason: '',
							created_at: new Date().toISOString()
						}
					],
					next_cursor: null,
					has_more: false
				})
			});
		});

		await page.goto('/history');

		await expect(page.locator('text=Milestone #100')).toBeVisible();
		await expect(page.getByText('Centennial celebration')).toBeVisible();
		await expect(page.locator('text=No reason specified')).toBeVisible();
		await expect(page.getByText('#100', { exact: true })).toBeVisible();
		await expect(page.getByText('#99', { exact: true })).toBeVisible();
	});

	test('renders error state and retries on button click', async ({ page }) => {
		let callCount = 0;
		await page.route('**/api/v1/history*', async (route) => {
			callCount += 1;
			if (callCount === 1) {
				await route.fulfill({
					status: 500,
					contentType: 'application/json',
					body: JSON.stringify({ error: 'Failed to fetch history' })
				});
			} else {
				await route.fulfill({
					status: 200,
					contentType: 'application/json',
					body: JSON.stringify({
						toggles: [
							{
								id: 1,
								state: true,
								reason: 'First toggle',
								created_at: new Date().toISOString()
							}
						],
						next_cursor: null,
						has_more: false
					})
				});
			}
		});

		await page.goto('/history');

		const tryAgainBtn = page.locator('button', { hasText: 'Try Again' });
		await expect(tryAgainBtn).toBeVisible();

		await tryAgainBtn.click();
		await expect(page.locator('text=First toggle')).toBeVisible();
	});
});
