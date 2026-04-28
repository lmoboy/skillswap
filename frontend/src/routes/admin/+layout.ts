import type { LayoutLoad } from './$types';

// Admin layout load - check admin access server-side if possible
export const load: LayoutLoad = async ({ fetch }) => {
  // Try to check if user is admin
  // Note: This is a client-side check, real security is in API endpoints
  try {
    const res = await fetch('/api/admin/stats', { credentials: 'include' });
    if (res.status === 401 || res.status === 403) {
      return {
        allowed: false,
        error: 'Admin access required'
      };
    }
    return { allowed: true };
  } catch {
    return {
      allowed: false,
      error: 'Authentication required'
    };
  }
};
