import { describe, it, expect, beforeEach, vi } from 'vitest';
import { get } from 'svelte/store';
import { auth } from '$lib/stores/auth';

describe('Auth Store', () => {
	beforeEach(() => {
		// Reset auth store before each test
		auth.clearUser();
	});

	describe('Initial State', () => {
		it('has correct initial values', () => {
			const authState = get(auth);
			
			expect(authState.user).toBeNull();
			expect(authState.isAuthenticated).toBe(false);
			expect(authState.loading).toBe(false);
			expect(authState.error).toBeNull();
			expect(authState.step).toBeNull();
		});
	});

	describe('User Management', () => {
		it('sets user correctly', () => {
			const user = { id: 1, name: 'testuser', email: 'test@example.com' };
			auth.setUser(user);
			
			const authState = get(auth);
			expect(authState.user).toEqual(user);
			expect(authState.isAuthenticated).toBe(true);
		});

		it('clears user correctly', () => {
			const user = { id: 1, name: 'testuser', email: 'test@example.com' };
			auth.setUser(user);
			auth.clearUser();
			
			const authState = get(auth);
			expect(authState.user).toBeNull();
			expect(authState.isAuthenticated).toBe(false);
		});
	});

	describe('Loading State', () => {
		it('sets loading state correctly', () => {
			auth.setLoading(true);
			
			const authState = get(auth);
			expect(authState.loading).toBe(true);
		});

		it('clears loading state correctly', () => {
			auth.setLoading(true);
			auth.setLoading(false);
			
			const authState = get(auth);
			expect(authState.loading).toBe(false);
		});
	});

	describe('Error Handling', () => {
		it('sets error message correctly', () => {
			const errorMessage = 'Login failed';
			auth.setError(errorMessage);
			
			const authState = get(auth);
			expect(authState.error).toBe(errorMessage);
		});

		it('clears error message correctly', () => {
			auth.setError('Some error');
			auth.setError(null);
			
			const authState = get(auth);
			expect(authState.error).toBeNull();
		});
	});

	describe('Step Management', () => {
		it('sets step correctly', () => {
			const step = 'Logging in...';
			auth.setStep(step);
			
			const authState = get(auth);
			expect(authState.step).toBe(step);
		});
	});

	describe('Store Subscriptions', () => {
		it('notifies subscribers of changes', () => {
			const callback = vi.fn();
			const unsubscribe = auth.subscribe(callback);
			
			// Initial call
			expect(callback).toHaveBeenCalledTimes(1);
			
			// Change user
			auth.setUser({ id: 1, name: 'test', email: 'test@example.com' });
			expect(callback).toHaveBeenCalledTimes(2);
			
			// Change loading
			auth.setLoading(true);
			expect(callback).toHaveBeenCalledTimes(3);
			
			unsubscribe();
		});

		it('stops notifying after unsubscribe', () => {
			const callback = vi.fn();
			const unsubscribe = auth.subscribe(callback);
			
			// Initial call
			expect(callback).toHaveBeenCalledTimes(1);
			
			unsubscribe();
			
			// Should not notify after unsubscribe
			auth.setUser({ id: 1, name: 'test', email: 'test@example.com' });
			expect(callback).toHaveBeenCalledTimes(1);
		});
	});

	describe('Authentication State', () => {
		it('is authenticated when user exists', () => {
			const user = { id: 1, name: 'testuser', email: 'test@example.com' };
			auth.setUser(user);
			
			const authState = get(auth);
			expect(authState.isAuthenticated).toBe(true);
		});

		it('is not authenticated when user is null', () => {
			auth.setUser(null);
			
			const authState = get(auth);
			expect(authState.isAuthenticated).toBe(false);
		});
	});

	describe('State Combinations', () => {
		it('handles loading with user', () => {
			const user = { id: 1, name: 'testuser', email: 'test@example.com' };
			auth.setUser(user);
			auth.setLoading(true);
			
			const authState = get(auth);
			expect(authState.user).toEqual(user);
			expect(authState.loading).toBe(true);
		});

		it('handles error with loading', () => {
			auth.setLoading(true);
			auth.setError('Network error');
			
			const authState = get(auth);
			expect(authState.loading).toBe(false); // setError sets loading to false
			expect(authState.error).toBe('Network error');
		});
	});
});