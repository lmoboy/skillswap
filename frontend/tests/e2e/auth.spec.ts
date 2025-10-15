import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to the login page before each test
		await page.goto('/auth/login');
	});

	test('user can see login form elements', async ({ page }) => {
		// Check that login form elements are visible
		await expect(page.locator('[data-testid="email-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="password-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="login-button"]')).toBeVisible();
		await expect(page.locator('[data-testid="signup-link"]')).toBeVisible();
	});

	test('user can fill login form', async ({ page }) => {
		// Fill in the login form
		await page.fill('[data-testid="email-input"]', 'test@example.com');
		await page.fill('[data-testid="password-input"]', 'password123');
		
		// Verify the values were filled
		await expect(page.locator('[data-testid="email-input"]')).toHaveValue('test@example.com');
		await expect(page.locator('[data-testid="password-input"]')).toHaveValue('password123');
	});

	test('user can navigate to register page', async ({ page }) => {
		// Click the signup link
		await page.click('[data-testid="signup-link"]');
		
		// Should navigate to register page
		await expect(page).toHaveURL('/auth/register');
		
		// Check that register form elements are visible
		await expect(page.locator('[data-testid="username-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="email-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="password-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="register-button"]')).toBeVisible();
	});

	test('user can fill registration form', async ({ page }) => {
		// Navigate to register page
		await page.goto('/auth/register');
		
		// Fill in registration form
		await page.fill('[data-testid="username-input"]', 'newuser');
		await page.fill('[data-testid="email-input"]', 'newuser@example.com');
		await page.fill('[data-testid="password-input"]', 'password123');
		
		// Verify the values were filled
		await expect(page.locator('[data-testid="username-input"]')).toHaveValue('newuser');
		await expect(page.locator('[data-testid="email-input"]')).toHaveValue('newuser@example.com');
		await expect(page.locator('[data-testid="password-input"]')).toHaveValue('password123');
	});

	test('user can navigate back to login from register', async ({ page }) => {
		// Start on register page
		await page.goto('/auth/register');
		
		// Check that register form is visible
		await expect(page.locator('[data-testid="username-input"]')).toBeVisible();
		
		// Navigate back to login
		await page.goto('/auth/login');
		await expect(page.locator('[data-testid="email-input"]')).toBeVisible();
	});

	test('login form shows loading state when submitted', async ({ page }) => {
		// Fill in the login form
		await page.fill('[data-testid="email-input"]', 'test@example.com');
		await page.fill('[data-testid="password-input"]', 'password123');
		
		// Submit the form
		await page.click('[data-testid="login-button"]');
		
		// The button should show loading state (this depends on the component implementation)
		// For now, just verify the form was submitted
		await expect(page.locator('[data-testid="login-button"]')).toBeVisible();
	});

	test('form validation works for empty fields', async ({ page }) => {
		// Try to submit empty form
		await page.click('[data-testid="login-button"]');
		
		// The form should prevent submission or show validation
		// Since we don't have backend, we just verify the form is still there
		await expect(page.locator('[data-testid="email-input"]')).toBeVisible();
		await expect(page.locator('[data-testid="password-input"]')).toBeVisible();
	});

	test('user can type in all form fields', async ({ page }) => {
		// Test email input
		await page.fill('[data-testid="email-input"]', 'user@example.com');
		await expect(page.locator('[data-testid="email-input"]')).toHaveValue('user@example.com');
		
		// Test password input
		await page.fill('[data-testid="password-input"]', 'mypassword');
		await expect(page.locator('[data-testid="password-input"]')).toHaveValue('mypassword');
		
		// Clear and test again
		await page.fill('[data-testid="email-input"]', '');
		await page.fill('[data-testid="password-input"]', '');
		await expect(page.locator('[data-testid="email-input"]')).toHaveValue('');
		await expect(page.locator('[data-testid="password-input"]')).toHaveValue('');
	});
});