// src/routes/+page.server.js
import { redirect } from '@sveltejs/kit';

export async function load({ fetch, cookies }) {
    const authCookie = cookies.get('authentication');

    if (!authCookie) {
        throw redirect(302, '/auth/login');
    }

    const response = await fetch("http://localhost:8080/api/cookieUser");
    const data = await response.json();

    if (data.error) {
        throw redirect(302, '/auth/login');
    }

    return {
        user: {
            name: data.username,
            email: data.email
        }
    };
}