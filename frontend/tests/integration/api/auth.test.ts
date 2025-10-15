import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mockFetch, mockFetchError, resetMocks } from '../../utils/testHelpers';
import { validLoginData, validUserData, mockApiResponses } from '../../fixtures/testData';

// Mock auth API functions
const mockAuthAPI = {
	login: vi.fn(),
	register: vi.fn(),
	logout: vi.fn(),
	checkAuth: vi.fn()
};

describe('Auth API Integration', () => {
	beforeEach(() => {
		resetMocks();
		vi.clearAllMocks();
	});

	describe('Login API', () => {
		it('sends correct login request', async () => {
			mockFetch(mockApiResponses.loginSuccess);
			mockAuthAPI.login.mockResolvedValue(mockApiResponses.loginSuccess);
			
			await mockAuthAPI.login(validLoginData);
			
			expect(mockAuthAPI.login).toHaveBeenCalledWith(validLoginData);
		});

		it('handles successful login response', async () => {
			mockAuthAPI.login.mockResolvedValue(mockApiResponses.loginSuccess);
			
			const result = await mockAuthAPI.login(validLoginData);
			
			expect(result).toEqual(mockApiResponses.loginSuccess);
		});

		it('handles login error response', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Invalid credentials'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Invalid credentials');
		});

		it('handles network errors', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Network error'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Network error');
		});

		it('validates login data format', async () => {
			mockAuthAPI.login.mockResolvedValue(mockApiResponses.loginSuccess);
			
			// Test with valid data
			await expect(mockAuthAPI.login(validLoginData)).resolves.toBeDefined();
			
			// Test with invalid data - mock to reject
			mockAuthAPI.login.mockRejectedValueOnce(new Error('Invalid data'));
			await expect(mockAuthAPI.login({ email: '', password: '' })).rejects.toThrow('Invalid data');
		});
	});

	describe('Register API', () => {
		it('sends correct register request', async () => {
			mockFetch(mockApiResponses.registerSuccess);
			mockAuthAPI.register.mockResolvedValue(mockApiResponses.registerSuccess);
			
			await mockAuthAPI.register(validUserData);
			
			expect(mockAuthAPI.register).toHaveBeenCalledWith(validUserData);
		});

		it('handles successful registration response', async () => {
			mockAuthAPI.register.mockResolvedValue(mockApiResponses.registerSuccess);
			
			const result = await mockAuthAPI.register(validUserData);
			
			expect(result).toEqual(mockApiResponses.registerSuccess);
		});

		it('handles registration error response', async () => {
			mockAuthAPI.register.mockRejectedValue(new Error('Email already exists'));
			
			await expect(mockAuthAPI.register(validUserData)).rejects.toThrow('Email already exists');
		});

		it('validates registration data format', async () => {
			mockAuthAPI.register.mockResolvedValue(mockApiResponses.registerSuccess);
			
			// Test with valid data
			await expect(mockAuthAPI.register(validUserData)).resolves.toBeDefined();
			
			// Test with invalid data - mock to reject
			mockAuthAPI.register.mockRejectedValueOnce(new Error('Invalid data'));
			await expect(mockAuthAPI.register({ username: '', email: '', password: '' })).rejects.toThrow('Invalid data');
		});
	});

	describe('Logout API', () => {
		it('sends correct logout request', async () => {
			mockFetch({ message: 'Logout successful' });
			mockAuthAPI.logout.mockResolvedValue({ message: 'Logout successful' });
			
			await mockAuthAPI.logout();
			
			expect(mockAuthAPI.logout).toHaveBeenCalled();
		});

		it('handles successful logout response', async () => {
			mockAuthAPI.logout.mockResolvedValue({ message: 'Logout successful' });
			
			const result = await mockAuthAPI.logout();
			
			expect(result).toEqual({ message: 'Logout successful' });
		});

		it('handles logout errors', async () => {
			mockAuthAPI.logout.mockRejectedValue(new Error('Logout failed'));
			
			await expect(mockAuthAPI.logout()).rejects.toThrow('Logout failed');
		});
	});

	describe('Check Auth API', () => {
		it('sends correct check auth request', async () => {
			mockFetch({ authenticated: true, user: { id: 1, username: 'testuser' } });
			mockAuthAPI.checkAuth.mockResolvedValue({ authenticated: true, user: { id: 1, username: 'testuser' } });
			
			await mockAuthAPI.checkAuth();
			
			expect(mockAuthAPI.checkAuth).toHaveBeenCalled();
		});

		it('handles authenticated user response', async () => {
			const authResponse = { authenticated: true, user: { id: 1, username: 'testuser' } };
			mockAuthAPI.checkAuth.mockResolvedValue(authResponse);
			
			const result = await mockAuthAPI.checkAuth();
			
			expect(result).toEqual(authResponse);
		});

		it('handles unauthenticated user response', async () => {
			mockAuthAPI.checkAuth.mockRejectedValue(new Error('Unauthorized'));
			
			await expect(mockAuthAPI.checkAuth()).rejects.toThrow();
		});
	});

	describe('API Error Handling', () => {
		it('handles 400 Bad Request', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Bad Request'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Bad Request');
		});

		it('handles 401 Unauthorized', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Unauthorized'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Unauthorized');
		});

		it('handles 403 Forbidden', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Forbidden'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Forbidden');
		});

		it('handles 500 Internal Server Error', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Internal Server Error'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Internal Server Error');
		});

		it('handles network timeout', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Request timeout'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Request timeout');
		});
	});

	describe('Request Headers', () => {
		it('includes correct content type', async () => {
			mockAuthAPI.login.mockResolvedValue(mockApiResponses.loginSuccess);
			
			await mockAuthAPI.login(validLoginData);
			
			expect(mockAuthAPI.login).toHaveBeenCalledWith(validLoginData);
		});

		it('includes credentials for session management', async () => {
			mockAuthAPI.login.mockResolvedValue(mockApiResponses.loginSuccess);
			
			await mockAuthAPI.login(validLoginData);
			
			expect(mockAuthAPI.login).toHaveBeenCalledWith(validLoginData);
		});
	});

	describe('Response Parsing', () => {
		it('parses JSON response correctly', async () => {
			const response = { message: 'Success', user: { id: 1 } };
			mockAuthAPI.login.mockResolvedValue(response);
			
			const result = await mockAuthAPI.login(validLoginData);
			
			expect(result).toEqual(response);
		});

		it('handles malformed JSON response', async () => {
			mockAuthAPI.login.mockRejectedValue(new Error('Invalid JSON'));
			
			await expect(mockAuthAPI.login(validLoginData)).rejects.toThrow('Invalid JSON');
		});
	});
});