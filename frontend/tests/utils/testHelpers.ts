import { vi } from 'vitest';
import { render, type RenderOptions } from '@testing-library/svelte';
import type { ComponentType } from 'svelte';

// Mock fetch responses
export const mockFetch = (response: any, status = 200) => {
	const mockResponse = {
		ok: status >= 200 && status < 300,
		status,
		json: vi.fn().mockResolvedValue(response),
		text: vi.fn().mockResolvedValue(JSON.stringify(response)),
		headers: new Headers(),
		clone: vi.fn().mockReturnThis()
	};
	
	vi.mocked(global.fetch).mockResolvedValue(mockResponse as any);
	return mockResponse;
};

// Mock fetch error
export const mockFetchError = (message = 'Network error') => {
	vi.mocked(global.fetch).mockRejectedValue(new Error(message));
};

// Reset all mocks
export const resetMocks = () => {
	vi.clearAllMocks();
	vi.mocked(global.fetch).mockClear();
};

// Custom render function with common providers
export const customRender = <T extends ComponentType>(
	component: T,
	props: any = {},
	options: RenderOptions = {}
) => {
	return render(component, { props, ...options });
};

// Wait for async operations
export const waitFor = (ms: number = 0) => 
	new Promise(resolve => setTimeout(resolve, ms));

// Simulate user interactions
export const simulateUserInput = async (element: HTMLElement, value: string) => {
	element.focus();
	(element as HTMLInputElement).value = value;
	element.dispatchEvent(new Event('input', { bubbles: true }));
	element.dispatchEvent(new Event('change', { bubbles: true }));
	await waitFor(10);
};

// Simulate form submission
export const simulateFormSubmit = async (form: HTMLFormElement) => {
	form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
	await waitFor(10);
};

// Check if element has specific classes
export const hasClass = (element: HTMLElement, className: string) => {
	return element.classList.contains(className);
};

// Check if element is visible
export const isVisible = (element: HTMLElement) => {
	const style = window.getComputedStyle(element);
	return style.display !== 'none' && style.visibility !== 'hidden' && style.opacity !== '0';
};

// Mock localStorage
export const mockLocalStorage = () => {
	const store: Record<string, string> = {};
	
	Object.defineProperty(window, 'localStorage', {
		value: {
			getItem: vi.fn((key: string) => store[key] || null),
			setItem: vi.fn((key: string, value: string) => {
				store[key] = value;
			}),
			removeItem: vi.fn((key: string) => {
				delete store[key];
			}),
			clear: vi.fn(() => {
				Object.keys(store).forEach(key => delete store[key]);
			}),
		},
		writable: true,
	});
	
	return store;
};

// Mock sessionStorage
export const mockSessionStorage = () => {
	const store: Record<string, string> = {};
	
	Object.defineProperty(window, 'sessionStorage', {
		value: {
			getItem: vi.fn((key: string) => store[key] || null),
			setItem: vi.fn((key: string, value: string) => {
				store[key] = value;
			}),
			removeItem: vi.fn((key: string) => {
				delete store[key];
			}),
			clear: vi.fn(() => {
				Object.keys(store).forEach(key => delete store[key]);
			}),
		},
		writable: true,
	});
	
	return store;
};

// Create mock file for file input testing
export const createMockFile = (name: string, type: string, size: number = 1024) => {
	const file = new File(['mock file content'], name, { type });
	Object.defineProperty(file, 'size', { value: size });
	return file;
};

// Mock FormData
export const mockFormData = (data: Record<string, any>) => {
	const formData = new FormData();
	Object.entries(data).forEach(([key, value]) => {
		if (value instanceof File) {
			formData.append(key, value);
		} else {
			formData.append(key, String(value));
		}
	});
	return formData;
};

// Test environment setup
export const setupTestEnvironment = () => {
	// Mock browser APIs
	Object.defineProperty(window, 'location', {
		value: {
			href: 'http://localhost:3000',
			origin: 'http://localhost:3000',
			pathname: '/',
			search: '',
			hash: '',
		},
		writable: true,
	});
	
	// Mock IntersectionObserver
	global.IntersectionObserver = vi.fn().mockImplementation(() => ({
		observe: vi.fn(),
		unobserve: vi.fn(),
		disconnect: vi.fn(),
	}));
	
	// Mock ResizeObserver
	global.ResizeObserver = vi.fn().mockImplementation(() => ({
		observe: vi.fn(),
		unobserve: vi.fn(),
		disconnect: vi.fn(),
	}));
};
