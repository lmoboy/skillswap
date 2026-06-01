import { checkAuth } from "$lib/api/auth";
import { redirect } from "@sveltejs/kit";
/**
 * Universal load function that runs on both server and client
 * to ensure auth state is available during SSR
 */

 const securePath = ['/course', '/settings', '/swapping'];
 
export async function load({ url, cookies }) {
    const pathname = url.pathname
    // Check authentication status
    if (securePath.some(path => pathname.startsWith(path))) {
        try {
          const isUser = await checkAuth(cookies);
          if (!isUser) {
            throw redirect(303, `/auth/login`);
          }
        } catch (error) {
          // Don't re-throw redirects from checkAuth
          if (error && typeof error === 'object' && 'status' in error && 'location' in error) {
            throw error;
          }
          console.error('Error checking auth in layout load:', error);
          throw redirect(303, `/auth/login`);
        }
      }
    
    // Return empty object - the auth state is managed by the store
    return {};
}

