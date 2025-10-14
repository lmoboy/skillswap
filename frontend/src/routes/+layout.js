import { checkAuth } from "$lib/api/auth";

/**
 * Universal load function that runs on both server and client
 * to ensure auth state is available during SSR
 */
export async function load({ fetch }) {
    // Check authentication status
    // This will run on both server and client, ensuring auth is always checked
    try {
        await checkAuth();
    } catch (error) {
        console.error("Error checking auth in layout load:", error);
    }
    
    // Return empty object - the auth state is managed by the store
    return {};
}

