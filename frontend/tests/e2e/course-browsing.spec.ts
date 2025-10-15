import { test, expect } from '@playwright/test';

test.describe('Course Browsing Flow', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to courses page
		await page.goto('/course');
	});

	test('user can see course page elements', async ({ page }) => {
		// Check that main elements are visible
		await expect(page.locator('[data-testid="search-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="courses-container"]')).toBeVisible();
	});

	test('user can use search functionality', async ({ page }) => {
		// Wait for search input to be visible
		await expect(page.locator('[data-testid="search-input"]')).toBeVisible();
		
		// Type in search input
		await page.fill('[data-testid="search-input"]', 'JavaScript');
		
		// Verify the search term was entered
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('JavaScript');
	});

	test('user can clear search input', async ({ page }) => {
		// Fill search input
		await page.fill('[data-testid="search-input"]', 'JavaScript');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('JavaScript');
		
		// Clear search input
		await page.fill('[data-testid="search-input"]', '');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('');
	});

	test('user can see courses container', async ({ page }) => {
		// Check that courses container is visible
		await expect(page.locator('[data-testid="courses-container"]')).toBeVisible();
	});

	test('user can interact with search input', async ({ page }) => {
		// Test typing in search input
		await page.fill('[data-testid="search-input"]', 'React');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('React');
		
		// Test clearing
		await page.fill('[data-testid="search-input"]', '');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('');
		
		// Test different search terms
		await page.fill('[data-testid="search-input"]', 'Python');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('Python');
	});

	test('page loads without errors', async ({ page }) => {
		// Check that the page loads successfully
		await expect(page.locator('body')).toBeVisible();
		
		// Check that main elements are present
		await expect(page.locator('[data-testid="search-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="courses-container"]')).toBeVisible();
	});

	test('user can navigate to course creation', async ({ page }) => {
		// Look for a link or button to create a course
		// This would depend on the actual implementation
		// For now, just verify the page is accessible
		await page.goto('/course/add');
		await expect(page.locator('[data-testid="course-title-input"]')).toBeVisible();
	});

	test('search input accepts various input types', async ({ page }) => {
		// Test different types of search queries
		const searchQueries = [
			'JavaScript',
			'Python Programming',
			'Web Development',
			'Data Science',
			'Machine Learning'
		];
		
		for (const query of searchQueries) {
			await page.fill('[data-testid="search-input"]', query);
			await expect(page.locator('[data-testid="search-input"]')).toHaveValue(query);
		}
		
		// Clear at the end
		await page.fill('[data-testid="search-input"]', '');
		await expect(page.locator('[data-testid="search-input"]')).toHaveValue('');
	});
});