import { error } from "@sveltejs/kit";

/**
 * Load user data for the page by fetching the API endpoint `/api/user?q=<id>`.
 *
 * @param {object} arg - The PageLoad input object.
 * @param {{ id: string }} arg.params - Route parameters; `id` is the user identifier used in the query.
 * @param {typeof fetch} arg.fetch - The fetch function for making requests.
 * @returns {Promise<any>} The user object returned by the API.
 * @throws {import('@sveltejs/kit').HttpError} Throws an error when the user cannot be retrieved.
 */
export async function load({ params, fetch: eventFetch }) {
  try {
    const response = await eventFetch(`/api/user?q=${params.id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      // Try to get error message from response
      let errorMessage = "Failed to load user";
      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch (e) {
        // If parsing fails, use default message
      }
      
      throw error(response.status, errorMessage);
    }

    const user = await response.json();
    
    if (!user) {
      throw error(404, "User not found");
    }

    return user;
  } catch (err) {
    // If it's already a SvelteKit error, rethrow it
    if (err && typeof err === 'object' && 'status' in err) {
      throw err;
    }
    // Otherwise, throw a 500 error
    console.error("Error loading user:", err);
    throw error(500, "Failed to load user profile");
  }
}
