import { describe, it, expect } from 'vitest';
import {
	validateEmail,
	validateUsername,
	validateRequired,
	validateMaxLength,
	validateMinLength,
	isValidEmail,
	isValidUrl
} from '$lib/utils/validation';

describe('Validation Utilities', () => {
	describe('validateEmail', () => {
		it('returns null for valid email', () => {
			expect(validateEmail('test@example.com')).toBeNull();
			expect(validateEmail('user.name@domain.co.uk')).toBeNull();
			expect(validateEmail('test+tag@example.org')).toBeNull();
		});

		it('returns error for empty email', () => {
			expect(validateEmail('')).toBe('Email is required');
			expect(validateEmail('   ')).toBe('Email is required');
		});

		it('returns error for email that is too long', () => {
			const longEmail = 'a'.repeat(95) + '@example.com';
			expect(validateEmail(longEmail)).toBe('Email is too long (max 100 characters)');
		});

		it('returns error for invalid email format', () => {
			expect(validateEmail('invalid-email')).toBe('Please enter a valid email address');
			expect(validateEmail('test@')).toBe('Please enter a valid email address');
			expect(validateEmail('@example.com')).toBe('Please enter a valid email address');
			expect(validateEmail('test.example.com')).toBe('Please enter a valid email address');
		});
	});

	describe('validateUsername', () => {
		it('returns null for valid username', () => {
			expect(validateUsername('testuser')).toBeNull();
			expect(validateUsername('Daniel')).toBeNull();
		});

		it('returns error for empty username', () => {
			expect(validateUsername('')).toBe('Username is required');
			expect(validateUsername('   ')).toBe('Username is required');
		});

		it('returns error for username that is too long', () => {
			const longUsername = 'a'.repeat(51);
			expect(validateUsername(longUsername)).toBe('Username is too long (max 50 characters)');
		});
	});

	describe('validateRequired', () => {
		it('returns null for non-empty value', () => {
			expect(validateRequired('test', 'Field')).toBeNull();
			expect(validateRequired('  test  ', 'Field')).toBeNull();
		});

		it('returns error for empty value', () => {
			expect(validateRequired('', 'Field')).toBe('Field is required');
			expect(validateRequired('   ', 'Field')).toBe('Field is required');
		});

		it('uses custom field name in error message', () => {
			expect(validateRequired('', 'Username')).toBe('Username is required');
			expect(validateRequired('', 'Email Address')).toBe('Email Address is required');
		});
	});

	describe('validateMaxLength', () => {
		it('returns null for value within limit', () => {
			expect(validateMaxLength('test', 10, 'Field')).toBeNull();
			expect(validateMaxLength('test', 4, 'Field')).toBeNull();
		});

		it('returns error for value exceeding limit', () => {
			expect(validateMaxLength('test', 3, 'Field')).toBe('Field cannot exceed 3 characters');
			expect(validateMaxLength('very long text', 5, 'Field')).toBe('Field cannot exceed 5 characters');
		});
	});

	describe('validateMinLength', () => {
		it('returns null for value meeting minimum', () => {
			expect(validateMinLength('test', 3, 'Field')).toBeNull();
			expect(validateMinLength('test', 4, 'Field')).toBeNull();
		});

		it('returns error for value below minimum', () => {
			expect(validateMinLength('ab', 3, 'Field')).toBe('Field must be at least 3 characters');
			expect(validateMinLength('a', 2, 'Field')).toBe('Field must be at least 2 characters');
		});
	});

	describe('isValidEmail', () => {
		it('returns true for valid email formats', () => {
			expect(isValidEmail('test@example.com')).toBe(true);
			expect(isValidEmail('user.name@domain.co.uk')).toBe(true);
			expect(isValidEmail('test+tag@example.org')).toBe(true);
		});

		it('returns false for invalid email formats', () => {
			expect(isValidEmail('invalid-email')).toBe(false);
			expect(isValidEmail('test@')).toBe(false);
			expect(isValidEmail('@example.com')).toBe(false);
			expect(isValidEmail('test.example.com')).toBe(false);
		});
	});

	describe('isValidUrl', () => {
		it('returns true for valid URL formats', () => {
			expect(isValidUrl('https://example.com')).toBe(true);
			expect(isValidUrl('http://example.com')).toBe(true);
			expect(isValidUrl('https://www.example.com')).toBe(true);
		});

		it('returns false for invalid URL formats', () => {
			expect(isValidUrl('invalid-url')).toBe(false);
			expect(isValidUrl('just-text')).toBe(false);
			expect(isValidUrl('')).toBe(false);
		});
	});
});