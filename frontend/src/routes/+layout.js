import { checkAuth } from "$lib/api/auth";
import { redirect } from "@sveltejs/kit";
/**
 * Universal load function that runs on both server and client
 * to ensure auth state is available during SSR
 */

 const securePath = ['/course', '/settings', '/swapping'];
 
export async function load({ url }) {
    const pathname = url.pathname
    // Check authentication status
    // This will run on both server and client, ensuring auth is always checked
    if (securePath.some(path => pathname.startsWith(path))) {
        try {
          const isUser = await checkAuth();
          if (!isUser) {
            // Redirect to login if not authenticated
            throw redirect(303, `/auth/login`);
          }
        } catch (error) {
          console.error('Error checking auth in layout load:', error);
          // Redirect to login on error as well
          throw redirect(303, `/auth/login`);
        }
      }
    
    // Return empty object - the auth state is managed by the store
    return {};
}

