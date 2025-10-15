import { test, expect } from '@playwright/test';

test.describe('Course Creation Flow', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to course creation page
		await page.goto('/course/add');
	});

	test('user can see course creation form elements', async ({ page }) => {
		// Check that main form elements are visible
		await expect(page.locator('[data-testid="course-title-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="course-description-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="skill-select"]')).toBeVisible();
		await expect(page.locator('[data-testid="duration-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="add-first-module-button"]')).toBeVisible();
	});

	test('user can fill basic course information', async ({ page }) => {
		// Fill in basic course information
		await page.fill('[data-testid="course-title-input"]', 'JavaScript Fundamentals');
		await page.fill('[data-testid="course-description-input"]', 'Learn JavaScript from scratch');
		await page.selectOption('[data-testid="skill-select"]', 'JavaScript');
		await page.fill('[data-testid="duration-input"]', '120');
		
		// Verify the values were filled
		await expect(page.locator('[data-testid="course-title-input"]')).toHaveValue('JavaScript Fundamentals');
		await expect(page.locator('[data-testid="course-description-input"]')).toHaveValue('Learn JavaScript from scratch');
		await expect(page.locator('[data-testid="skill-select"]')).toHaveValue('JavaScript');
		await expect(page.locator('[data-testid="duration-input"]')).toHaveValue('120');
	});

	test('user can add first module', async ({ page }) => {
		// Wait for the page to load
		await page.waitForLoadState('networkidle');
		
		// Click add first module button
		await page.click('[data-testid="add-first-module-button"]');
		
		// Wait for module form to appear
		await page.waitForSelector('[data-testid="module-item"]', { timeout: 5000 });
		
		// Check that module form appears
		await expect(page.locator('[data-testid="module-item"]')).toBeVisible();
	});

	test('user can add multiple modules', async ({ page }) => {
		// Wait for the page to load
		await page.waitForLoadState('networkidle');
		
		// Add first module
		await page.click('[data-testid="add-first-module-button"]');
		await page.waitForSelector('[data-testid="module-item"]', { timeout: 5000 });
		await expect(page.locator('[data-testid="module-item"]')).toBeVisible();
		
		// Add second module
		await page.click('[data-testid="add-module-button"]');
		
		// Should have 2 modules now
		const modules = page.locator('[data-testid="module-item"]');
		await expect(modules).toHaveCount(2);
	});

	test('user can remove modules', async ({ page }) => {
		// Wait for the page to load
		await page.waitForLoadState('networkidle');
		
		// Add a module first
		await page.click('[data-testid="add-first-module-button"]');
		await page.waitForSelector('[data-testid="module-item"]', { timeout: 5000 });
		await expect(page.locator('[data-testid="module-item"]')).toBeVisible();
		
		// Remove the module
		await page.click('[data-testid="remove-module-button"]');
		
		// Module should be removed
		await expect(page.locator('[data-testid="module-item"]')).toHaveCount(0);
	});

	test('user can fill module information', async ({ page }) => {
		// Wait for the page to load
		await page.waitForLoadState('networkidle');
		
		// Add a module
		await page.click('[data-testid="add-first-module-button"]');
		await page.waitForSelector('[data-testid="module-item"]', { timeout: 5000 });
		
		// Fill module information
		await page.fill('[data-testid="module-title-0"]', 'Introduction to JavaScript');
		await page.fill('[data-testid="module-description-0"]', 'Basic concepts and syntax');
		await page.fill('[data-testid="module-duration-0"]', '30');
		
		// Verify the values were filled
		await expect(page.locator('[data-testid="module-title-0"]')).toHaveValue('Introduction to JavaScript');
		await expect(page.locator('[data-testid="module-description-0"]')).toHaveValue('Basic concepts and syntax');
		await expect(page.locator('[data-testid="module-duration-0"]')).toHaveValue('30');
	});

	test('user can see submit button', async ({ page }) => {
		// Check that submit button is visible
		await expect(page.locator('[data-testid="submit-course-button"]')).toBeVisible();
	});

	test('form validation prevents empty submission', async ({ page }) => {
		// Try to submit empty form
		await page.click('[data-testid="submit-course-button"]');
		
		// Form should still be visible (validation should prevent submission)
		await expect(page.locator('[data-testid="course-title-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="course-description-input"]')).toBeVisible();
	});

	test('user can type in all input fields', async ({ page }) => {
		// Test title input
		await page.fill('[data-testid="course-title-input"]', 'My Course');
		await expect(page.locator('[data-testid="course-title-input"]')).toHaveValue('My Course');
		
		// Test description input
		await page.fill('[data-testid="course-description-input"]', 'Course description');
		await expect(page.locator('[data-testid="course-description-input"]')).toHaveValue('Course description');
		
		// Test duration input
		await page.fill('[data-testid="duration-input"]', '90');
		await expect(page.locator('[data-testid="duration-input"]')).toHaveValue('90');
		
		// Clear and test again
		await page.fill('[data-testid="course-title-input"]', '');
		await page.fill('[data-testid="course-description-input"]', '');
		await page.fill('[data-testid="duration-input"]', '');
		await expect(page.locator('[data-testid="course-title-input"]')).toHaveValue('');
		await expect(page.locator('[data-testid="course-description-input"]')).toHaveValue('');
		await expect(page.locator('[data-testid="duration-input"]')).toHaveValue('');
	});

	test('skill select dropdown works', async ({ page }) => {
		// Check that skill select is visible and has options
		const skillSelect = page.locator('[data-testid="skill-select"]');
		await expect(skillSelect).toBeVisible();
		
		// Select a skill
		await page.selectOption('[data-testid="skill-select"]', 'JavaScript');
		await expect(skillSelect).toHaveValue('JavaScript');
	});
});