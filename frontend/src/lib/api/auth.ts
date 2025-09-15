import { auth } from '$lib/stores/auth';

const API_BASE = '/api';

const defaultOptions: RequestInit = {
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
    },
};

async function handleResponse<T>(response: Response): Promise<T> {
    const data = await response.json();
    if (!response.ok) {
        throw new Error(data.message || 'Request failed');
    }
    return data;
}

export interface LoginResponse {
    user?: {
        name: string;
        email: string;
    };
    error?: string;
    status: string;
}

export async function login(credentials: { email: string; password: string }): Promise<LoginResponse> {
    try {
        auth.setLoading(true);
        const response = await fetch(`${API_BASE}/login`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(credentials),
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data || 'Login failed');
        }

        if (data.user) {
            auth.setUser({
                name: data.user.username || data.user.name || '',
                email: data.user.email || ''
            });
        }

        return {
            ...data,
            status: data.status || (data.user ? 'ok' : 'error'),
            error: data.error || (!data.user ? 'No user data received' : undefined)
        };
    } catch (error: unknown) {
        console.error('Login error:', error);
        const errorMessage = error instanceof Error ? error.message : 'Login failed';
        auth.setError(errorMessage);
        throw new Error(errorMessage);
    } finally {
        auth.setLoading(false);
    }
}

export async function logout(): Promise<void> {
    try {
        auth.setLoading(true);
        const response = await fetch(`${API_BASE}/logout`, {
            ...defaultOptions,
            method: 'POST',
        });

        if (!response.ok) {
            throw new Error('Logout failed');
        }

        // Clear user data regardless of response
        auth.clearUser();
    } catch (error: unknown) {
        console.error('Logout error:', error);
        const errorMessage = error instanceof Error ? error.message : 'Logout failed';
        auth.setError(errorMessage);
        throw new Error(errorMessage);
    } finally {
        auth.setLoading(false);
    }
}

export async function checkAuth(): Promise<boolean> {
    try {
        auth.setLoading(true);
        const response = await fetch(`${API_BASE}/cookieUser`, {
            ...defaultOptions,
            method: 'GET',
        });

        if (!response.ok) {
            console.warn('Not authenticated');
        }

        const data = await response.json();

        if (data && data.username) {
            auth.setUser({
                name: data.username,
                email: data.email || ''
            });
            return true;
        }

        auth.clearUser();
        return false;
    } catch (error: unknown) {
        console.error('Auth check failed:', error);
        auth.clearUser();
        return false;
    } finally {
        auth.setLoading(false);
    }
}
