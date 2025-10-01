import { error } from '@sveltejs/kit';
/**
 * Load user data for the page by fetching the API endpoint `/api/user?q=<id>`.
 *
 * @param {{ params: { id: string }, fetch: typeof fetch }} arg - The PageLoad input object.
 * @param {{ id: string }} arg.params - Route parameters; `id` is the user identifier used in the query.
 * @returns {Promise<any>} The user object returned by the API.
 * @throws {import('@sveltejs/kit').HttpError} Throws a 403 error when the user cannot be retrieved.
 */
export async function load({ params, fetch: eventFetch }) {

    const user = await eventFetch(`/api/user?q=${params.id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            return response.json();
        }
    }).catch(error => {
        console.log(error);
    });

    if (!user) {
        throw error(403, 'Error fetching user user');
    }

    return user;
}




