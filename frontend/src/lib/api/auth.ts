import { auth } from '$lib/stores/auth';
import { getReturnUrl, clearReturnUrl } from '$lib/utils/apiClient';

const API_BASE = '/api';

// Noklusējuma pieprasījuma opcijas (angļu v. default options)
const defaultOptions: RequestInit = {
    credentials: 'include', // Iekļaut akreditācijas datus, piemēram, sīkdatnes (angļu v. cookies)
    headers: {
        'Content-Type': 'application/json', // Norādīt satura tipu kā JSON
    },
};

/**
 * Apstrādā atbildi (angļu v. response) no API un parsē JSON datus (angļu v. JSON data).
 * @template T Atgriežamā tipa norāde (angļu v. generic type).
 * @param {Response} response Atbildes objekts (angļu v. response object) no fetch.
 * @returns {Promise<T>} Promise ar datiem (angļu v. data) vai kļūdu (angļu v. error).
 */
async function handleResponse<T>(response: Response): Promise<T> {
    const data = await response.json(); // Parsē atbildi kā JSON
    if (!response.ok) {
        // Ja atbilde nav veiksmīga (piem., 4xx vai 5xx statusa kods), izmest kļūdu (angļu v. throw an error)
        throw new Error(data.message || 'Pieprasījums neizdevās (angļu v. Request failed)');
    }
    return data;
}

// Interfeiss (angļu v. interface) veiksmīgai pieteikšanās (angļu v. login) atbildei
export interface LoginResponse {
    returnUrl: any;
    user?: {
        name: string;
        email: string;
        id: string;
    };
    error?: string;
    status: string;
}

/**
 * Funkcija lietotāja pieteikšanai (angļu v. login).
 * @param {object} credentials Lietotāja e-pasts (angļu v. email) un parole (angļu v. password).
 * @returns {Promise<LoginResponse>} Promise ar pieteikšanās atbildes datiem (angļu v. login response data).
 */
export async function login(credentials: { email: string; password: string }): Promise<LoginResponse> {
    try {
        auth.setLoading(true); // Iestatīt ielādes statusu (angļu v. loading status) uz 'true'
        auth.setStep("Datu ielāde... (angļu v. Fetching data...)") // Iestatīt pašreizējo darbības soli
        // Veikt POST pieprasījumu (angļu v. POST request) uz pieteikšanās galapunktu (angļu v. login endpoint)
        const response = await fetch(`${API_BASE}/login`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(credentials), // Pārvērst datus (angļu v. data) par JSON virkni (angļu v. JSON string)
            credentials: 'include', // Svarīgi sīkdatnēm (angļu v. cookies)
        });

        auth.setStep("Pārvēršanās uz json... (angļu v. Converting to json...)")
        const data = await response.json(); // Parsēt atbildi kā JSON

        if (!response.ok) {
            throw new Error(data.error || 'Pieteikšanās neizdevās (angļu v. Login failed)');
        }

        auth.setStep("Lietotājs iestatīts, pārbauda autentifikāciju... (angļu v. User set, checking auth...)")
        // Login endpoint only returns status, get user data from checkAuth
        const isAuthenticated = await checkAuth();
        if (!isAuthenticated) {
            throw new Error('Autentifikācija neizdevās pēc pieteikšanās (angļu v. Authentication failed after login)');
        }

        auth.setStep("Pieteikšanās pabeigta (angļu v. Login done)")

        // Wait for user data to be set by checkAuth
        const currentUser = await auth.waitForUser();
        return {
            status: 'ok',
            user: {
                name: currentUser.name || '',
                email: currentUser.email || '',
                id: currentUser.id.toString() || '',
            }
        };
    } catch (error) {
        // Apstrādāt kļūdas, kas rodas pieteikšanās laikā
        const errorMessage = error instanceof Error ? error.message : 'Pieteikšanās neizdevās (angļu v. Login failed)';
        auth.setError(errorMessage); // Iestatīt kļūdas ziņojumu
        throw error; // Pārsūtīt kļūdu tālāk
    } finally {
        auth.setLoading(false); // Vienmēr iestatīt ielādes statusu uz 'false'
    }
}


export async function register(credentials: { username: string, email: string, password: string }) {
    try {
        auth.setLoading(true); // Iestatīt ielādes statusu (angļu v. loading status) uz 'true'
        auth.setStep("Datu ielāde... (angļu v. Fetching data...)") // Iestatīt pašreizējo darbības soli
        // Veikt POST pieprasījumu (angļu v. POST request) uz pieteikšanās galapunktu (angļu v. login endpoint)
        const response = await fetch(`${API_BASE}/register`, {
            ...defaultOptions,
            method: 'POST',
            body: JSON.stringify(credentials), // Pārvērst datus (angļu v. data) par JSON virkni (angļu v. JSON string)
            credentials: 'include', // Svarīgi sīkdatnēm (angļu v. cookies)
        });

        auth.setStep("Pārvēršanās uz json... (angļu v. Converting to json...)")
        const data = await response.json(); // Parsēt atbildi kā JSON

        if (!response.ok) {
            // Ja atbilde nav veiksmīga, izmest kļūdu
            throw new Error(data.error || 'Pieteikšanās neizdevās (angļu v. Login failed)');
        }
        auth.setStep("Lietotāja iestatīšana... (angļu v. Setting user...)")
        if (data) {
            // Ja ir saņemti dati, iestatīt lietotāju (angļu v. user)
            auth.setUser({
                name: data.username || '',
                email: data.email || '',
                id: data.id || '',
            });

            auth.setStep("Lietotājs iestatīts, pārbauda autentifikāciju... (angļu v. User set, checking auth...)")
            await checkAuth(); // Pārbaudīt lietotāja autentifikāciju (angļu v. check auth)
        } else {
            throw new Error('Nederīgs atbildes formāts no servera (angļu v. Invalid response format from server)');
        }
        auth.setStep("Pieteikšanās pabeigta (angļu v. Login done)")

        return {
            status: 'ok',
            user: {
                name: data.username || '',
                email: data.email || '',
                id: data.id || '',
            },
            returnUrl: getReturnUrl('/')
        };
    } catch (error) {
        // Apstrādāt kļūdas, kas rodas pieteikšanās laikā
        const errorMessage = error instanceof Error ? error.message : 'Pieteikšanās neizdevās (angļu v. Login failed)';
        auth.setError(errorMessage); // Iestatīt kļūdas ziņojumu
        throw error; // Pārsūtīt kļūdu tālāk
    } finally {
        auth.setLoading(false); // Vienmēr iestatīt ielādes statusu uz 'false'
    }

}



/**
 * Funkcija lietotāja izrakstīšanai (angļu v. logout).
 * @returns {Promise<void>}
 */
export async function logout(): Promise<void> {
    try {
        auth.setLoading(true); // Sākt ielādi
        const response = await fetch(`${API_BASE}/logout`, {
            ...defaultOptions,
            method: 'POST',
            credentials: 'include', // Svarīgi autentifikācijai, kas balstīta uz sīkdatnēm (angļu v. cookie-based auth)
        });
        auth.setStep("Žetona noņemšana no db... (angļu v. Removing token from db...)")
        // Vienmēr notīrīt lietotāja stāvokli (angļu v. user state), pat ja pieprasījums neizdodas
        auth.clearUser();
        auth.setStep("Atbildes pārbaude... (angļu v. Checking response...)")
        if (!response.ok) {
            const data = await response.json().catch(() => ({}));
            throw new Error(data.error || 'Izrakstīšanās neizdevās (angļu v. Logout failed)');
        }
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Izrakstīšanās neizdevās (angļu v. Logout failed)';
        auth.setStep('Izrakstīšanās kļūda: ' + errorMessage);
        // Joprojām notīrīt lietotāja stāvokli, ja rodas kļūda
        auth.clearUser();
        throw new Error(errorMessage);
    } finally {
        auth.setStep("Izrakstīšanās pabeigta (angļu v. Logout done)");
        auth.setLoading(false); // Beigt ielādi
    }
}

/**
 * Pārbauda, vai lietotājs ir autentificēts, pamatojoties uz sīkdatnēm (angļu v. cookies).
 * When called with cookies (SSR), skips store operations and returns boolean directly.
 * @returns {Promise<boolean>} 'true' ja autentificēts, 'false' ja nav.
 */
export async function checkAuth(cookies?: { get: (name: string) => string | undefined }): Promise<boolean> {
    const isSSR = !!cookies;

    try {
        if (!isSSR) auth.setLoading(true);

        const fetchHeaders: Record<string, string> = {
            'Content-Type': 'application/json',
        };

        // Forward cookies for SSR context
        if (cookies) {
            const sessionCookie = cookies.get('authentication');
            if (sessionCookie) {
                fetchHeaders['Cookie'] = `authentication=${sessionCookie}`;
            }
        }

        const response = await fetch(`${API_BASE}/cookieUser`, {
            method: 'GET',
            headers: fetchHeaders,
            credentials: 'include',
        });

        if (!isSSR) auth.setStep("Datu ielāde no bd...");

        const data = await response.json();

        if (!isSSR) auth.setStep("Datu konvertēšana...");

        if (!response.ok) {
            if (!isSSR) {
                console.warn('Nav autentificēts:', data.error || 'Nav kļūdas detaļu');
                auth.clearUser();
            }
            return false;
        }

        if (data) {
            if (!isSSR) {
                auth.setUser({
                    name: data.user || '',
                    email: data.email || '',
                    id: data.id || '',
                    profile_picture: data.profile_picture ? data.profile_picture : '',
                    is_admin: data.is_admin || false,
                    swaps: data.swaps || 0,
                });
            }
            return true;
        }

        if (!isSSR) auth.clearUser();
        return false;
    } catch (error: unknown) {
        if (!isSSR) {
            const errorMessage = error instanceof Error ? error.message : 'Nezināma kļūda';
            auth.setStep('Autentifikācijas pārbaude neizdevās:' + errorMessage);
            auth.clearUser();
        }
        return false;
    } finally {
        if (!isSSR) {
            auth.setStep("Autentifikācijas pārbaude pabeigta");
            auth.setLoading(false);
        }
    }
}

// Export utility functions for handling return URLs
export { getReturnUrl, clearReturnUrl } from '$lib/utils/apiClient';