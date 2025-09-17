
import { error } from '@sveltejs/kit';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {

    const user = await fetch('https://localhost:8080/api/user?q=' + params.id, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            return response.json();
        }
    })







    if (!user) {
        throw error(404, 'Not found');
    }

    return {
        user: user,
    };
}