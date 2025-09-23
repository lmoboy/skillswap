import { error } from '@sveltejs/kit';

/** @type {import('./$types').PageLoad} */
export async function load({ params, fetch: eventFetch }) {

    const user = await eventFetch(`/api/user?q=${params.id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            // console.log(response.json());
            return response.json();
        }
    })
    console.log(user)
    // user.contacts = JSON.parse(user.contacts);
    // user.projects = JSON.parse(user.projects);
    if (!user) {
        throw error(403, 'Error fetching user user');
    }

    return user;
}




