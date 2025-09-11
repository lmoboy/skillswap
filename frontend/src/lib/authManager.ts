class AuthManager {
    private name: string;
    private email: string;
    private picture: string;

    constructor() {
        this.email = ''
        this.name = ''
        this.picture = ''
        this.getUserFromCookie()
    }

    public getUserInfo() {
        return {
            name: this.name,
            email: this.email,
            picture: this.picture,
        };
    }

    public async getUserFromCookie() {
        const response = await fetch("http://localhost:8080/api/cookieUser", {
            credentials: "include",
        });
        const data = await response.json();
        this.name = data.name;
        this.email = data.email;
        this.picture = data.picture;
    }
}

const Auth = new AuthManager();

export default AuthManager;
