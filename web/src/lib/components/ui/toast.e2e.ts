import { test, expect } from '@playwright/test';

test.describe('Toast Notification Component', () => {
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

	test('surfaces exception toast when backend returns an error on toggle', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 429,
				contentType: 'application/json',
				body: JSON.stringify({ error: 'Rate limit exceeded: please wait', cooldown_ms: 5000 })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await expect(bulbButton).toBeVisible();
		await bulbButton.click();

		const toastCard = page.locator('[data-testid="toast-card"]');
		await expect(toastCard).toBeVisible();
		await expect(toastCard).toContainText('Rate limit exceeded: please wait');
	});

	test('dismisses toast when clicking close button', async ({ page }) => {
		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 500,
				contentType: 'application/json',
				body: JSON.stringify({ error: 'Internal Server Error' })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click();

		const toastCard = page.locator('[data-testid="toast-card"]');
		await expect(toastCard).toBeVisible();

		const dismissBtn = page.locator('[data-testid="toast-dismiss"]');
		await dismissBtn.click();

		await expect(toastCard).toHaveCount(0);
	});

	test('stacks up to 5 toasts and evicts oldest on overflow', async ({ page }) => {
		let count = 0;
		await page.route('**/api/v1/toggle', async (route) => {
			count += 1;
			await route.fulfill({
				status: 400,
				contentType: 'application/json',
				body: JSON.stringify({ error: `Error message ${count}` })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');

		// Trigger 6 errors via bulb interaction
		for (let i = 0; i < 6; i++) {
			await bulbButton.click();
			await page.waitForTimeout(50);
		}

		const toasts = page.locator('[data-testid="toast-card"]');
		await expect(toasts).toHaveCount(5);

		// The first message (Error message 1) should be evicted
		await expect(toasts.first()).not.toContainText('Error message 1');
		await expect(toasts.last()).toContainText('Error message 6');
	});

	test('renders centered toast on small viewport', async ({ page }) => {
		await page.setViewportSize({ width: 375, height: 667 });

		await page.route('**/api/v1/toggle', async (route) => {
			await route.fulfill({
				status: 429,
				contentType: 'application/json',
				body: JSON.stringify({ error: 'Mobile rate limit error' })
			});
		});

		const bulbButton = page.locator('button[aria-label*="Turn"]');
		await bulbButton.click();

		const toastAside = page.locator('aside[aria-label="Notifications"]');
		const toastCard = page.locator('[data-testid="toast-card"]');

		await expect(toastCard).toBeVisible();
		await expect(toastCard).toBeInViewport();
		await expect(toastAside).toHaveClass(/inset-x-0/);
		await expect(toastAside).toHaveClass(/mx-auto/);
	});
});
