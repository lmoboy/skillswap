
class AuthData {
    token: string | null;

    constructor(token: string | null) {
        if (token == null) {
            token = getAuthToken();
        }
        this.token = token;
    }
}
function getAuthToken(): string | null {
    return localStorage.getItem("authentication");
}

function setAuthToken(token: string): void {
    localStorage.setItem("authentication", token);
}

function clearAuthToken(): void {
    localStorage.removeItem("authentication");
}

async function isAuthenticated() {
    if (!getAuthToken()) {
        return false;
    }
    await fetch('/api/checks').then(response => {
        if (response.status === 200) {
            return true;
        }
        clearAuthToken();
        return false;
    }).catch(() => {
        clearAuthToken();
        return false;
    });
}


const Auth = new AuthData(null);

export {
    isAuthenticated,
    AuthData,
}