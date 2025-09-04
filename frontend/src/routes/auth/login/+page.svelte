<script>
    let email = "";
    let password = "";
    let error = "";

    const handleSubmit = async () => {
        error = "";
        // Dummy Sign up logic

        if (email && password) {
            if (email.length > 100) {
                error = "Email too long!";
                return;
            }
            if (password.length > 50) {
                error = "Password cannot be longer than 50 characters!";
                return;
            }
            fetch("http://localhost:8080/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "multipart/form-data",
                },
                body: JSON.stringify({ email, password }),
                credentials: "include",
            });
        } else {
            error = "Invalid email or password";
        }
    };
</script>

<div class="Signup-container">
    <h2>Sign in</h2>
    {#if error}
        <div class="error">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleSubmit}>
        <label for="email">Email</label>
        <input
            id="email"
            type="email"
            bind:value={email}
            required
            autocomplete="username"
        />

        <label for="password">Password</label>
        <input
            id="password"
            type="password"
            bind:value={password}
            required
            autocomplete="current-password"
        />

        <button type="submit">Sign in</button>
        <a href="/auth/register">Don't have an account?</a>
    </form>
</div>

<style>
    .Signup-container {
        max-width: 350px;
        margin: 2rem auto;
        padding: 2rem;
        border: 1px solid #eee;
        border-radius: 8px;
        background: #fafafa;
    }
    label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    input {
        width: 100%;
        padding: 0.5rem;
        margin-bottom: 1rem;
        border-radius: 4px;
        border: 1px solid #ccc;
    }
    button {
        width: 100%;
        padding: 0.7rem;
        background: #0070f3;
        color: white;
        border: none;
        border-radius: 4px;
        font-weight: bold;
        cursor: pointer;
    }
    .error {
        color: #d00;
        margin-bottom: 1rem;
        text-align: center;
    }
</style>
