import { goto } from '$app/navigation';
import { browser } from '$app/environment';

interface ApiResponse<T = any> {
    error?: string;
    redirect?: string;
    previousPath?: string;
    data?: T;
}

/**
 * Enhanced fetch wrapper that handles authentication redirects automatically
 */
export async function apiFetch<T = any>(
    url: string,
    options: RequestInit = {}
): Promise<Response> {
    const response = await fetch(url, {
        credentials: 'include',
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    });

    // Check if response is a redirect from auth middleware
    if (response.status === 401 && browser) {
        try {
            const data: ApiResponse = await response.clone().json();
            if (data.redirect && data.previousPath) {
                // Store the intended destination
                sessionStorage.setItem('returnUrl', data.previousPath);
                // Redirect to login
                await goto(data.redirect);
                throw new Error('Authentication required - redirecting to login');
            }
        } catch (e) {
            // If parsing fails, just return the response
            if (e instanceof Error && e.message.includes('redirecting')) {
                throw e;
            }
        }
    }

    return response;
}

/**
 * Returns to the previously saved return URL or defaults to home
 */
export function getReturnUrl(defaultUrl: string = '/'): string {
    if (!browser) return defaultUrl;
    
    const returnUrl = sessionStorage.getItem('returnUrl');
    if (returnUrl) {
        sessionStorage.removeItem('returnUrl');
        return returnUrl;
    }
    return defaultUrl;
}

/**
 * Clears any stored return URL
 */
export function clearReturnUrl(): void {
    if (browser) {
        sessionStorage.removeItem('returnUrl');
    }
}

