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
            return response.json();
        }
    }).catch(error => {
        console.log(error);
    });

    // console.log(user.contacts)
    user.contacts = JSON.parse(user.contacts);
    user.skills = JSON.parse(user.skills);
    user.projects = JSON.parse(user.projects);
    if (!user) {
        throw error(403, 'Error fetching user user');
    }

    return user;
}




