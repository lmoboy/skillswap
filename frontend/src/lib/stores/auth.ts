import { writable } from 'svelte/store';

export interface User {
    name: string;
    email: string;
}

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    loading: boolean;
    error: string | null;
}

const defaultState: AuthState = {
    user: null,
    isAuthenticated: false,
    loading: true,
    error: null
};

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>(defaultState);

    return {
        subscribe,
        setUser: (user: User | null) => {
            update(state => ({
                ...state,
                user,
                isAuthenticated: user !== null,
                loading: false,
                error: null
            }));
        },
        setLoading: (loading: boolean) => {
            update(state => ({
                ...state,
                loading
            }));
        },
        setError: (error: string | null) => {
            update(state => ({
                ...state,
                error,
                loading: false
            }));
        },
        clearUser: () => {
            set({
                ...defaultState,
                loading: false
            });
        },
        isAuthenticated: () => {
            let isAuth = false;
            update(state => {
                isAuth = state.isAuthenticated;
                return state;
            });
            return isAuth;
        }
    };
}

export const auth = createAuthStore();

export const { setUser, clearUser } = auth;
