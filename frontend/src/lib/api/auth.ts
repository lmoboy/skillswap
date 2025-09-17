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
        id: string;
    };
    error?: string;
    status: string;
}

export async function login(credentials: { email: string; password: string }): Promise<LoginResponse> {
    try {
        auth.setLoading(true);
        auth.setStep("Fetching data...")
        const response = await fetch(`${API_BASE}/login`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(credentials),
            credentials: 'include', // Important for cookies
        });

        auth.setStep("Converting to json...")
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }
        auth.setStep("Setting user...")
        if (data) {
            auth.setUser({
                name: data.username || '',
                email: data.email || '',
                id: data.id || '',
            });

            auth.setStep("User set, checking auth...")
            await checkAuth();
        } else {
            throw new Error('Invalid response format from server');
        }
        auth.setStep("Login done")

        return {
            status: 'ok',
            user: {
                name: data.username || '',
                email: data.email || '',
                id: data.id || '',
            }
        };
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Login failed';
        auth.setError(errorMessage);
        throw error;
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
            credentials: 'include', // Important for cookie-based auth
        });
        auth.setStep("Removing token from db...")
        // Always clear the user state, even if the request fails
        auth.clearUser();
        auth.setStep("Checking response...")
        if (!response.ok) {
            const data = await response.json().catch(() => ({}));
            throw new Error(data.error || 'Logout failed');
        }
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Logout failed';
        auth.setStep('Logout error:' + errorMessage);
        // Still clear user state even if there's an error
        auth.clearUser();
        throw new Error(errorMessage);
    } finally {
        auth.setStep("Logout done");
        auth.setLoading(false);
    }
}

export async function checkAuth(): Promise<boolean> {
    try {
        auth.setLoading(true);
        const response = await fetch(`${API_BASE}/cookieUser`, {
            ...defaultOptions,
            method: 'GET',
            credentials: 'include', // Make sure to include credentials
        });
        auth.setStep("Fetching data from bd...");
        // console.log(response);
        const data = await response.json();
        auth.setStep("Converting data...");

        if (!response.ok) {
            auth.setStep("Error occured with response...");

            console.warn('Not authenticated:', data.error || 'No error details');
            auth.clearUser();
            return false;
        }
        // console.log(data);
        if (data) {
            auth.setStep("Set user data...");

            auth.setUser({
                name: data.username || '',
                email: data.email || '',
                id: data.id || '',
            });
            return true;
        }
        // console.log("we got this far for some reason");
        auth.setStep("Data failed, clearing...");

        auth.clearUser();
        return false;
    } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'Unknown error';
        auth.setStep('Auth check failed:' + errorMessage);
        auth.clearUser();
        return false;
    } finally {
        auth.setStep("Auth check done");
        auth.setLoading(false);
    }
}
