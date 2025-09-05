<script>
    let username = "";
    let email = "";
    let password = "";
    let passwordr = "";
    let error = "";

    const handleSubmit = async () => {
        error = "";
        // Dummy register logic
        if (email && password && passwordr && username) {
            if (username.length > 50) {
                error = "Username cannot be longer than 50 characters!";
                return;
            }
            if (email.length > 100) {
                error = "Email too long!";
                return;
            }
            if (password.length > 50) {
                error = "Password cannot be longer than 50 characters!";
                return;
            }
            if (password != passwordr) {
                error = "Passwords must match!";
                return;
            }
            // Redirect or show success
            fetch("http://localhost:8080/api/register", {
                method: "POST",
                headers: {
                    "Content-Type": "multipart/form-data",
                },
                body: JSON.stringify({ username, email, password }),
            })
                .then((response) => {
                    if (response.ok) {
                        alert("Registration successful!");
                    } else {
                        error = "Registration failed";
                    }
                })
                .catch(() => {
                    error = "Network error";
                });
        } else {
            error = "You must fill out all the fields";
        }
    };
</script>

<div class="login-container">
    <h2>Sign up</h2>
    {#if error}
        <div class="error">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleSubmit}>
        <label for="username">Username</label>
        <input id="username" type="username" bind:value={username} required />

        <label for="email">Email</label>
        <input id="email" type="email" bind:value={email} required />

        <label for="password">Password</label>
        <input id="password" type="password" bind:value={password} required />
        <label for="password">Password repeat</label>
        <input id="password" type="password" bind:value={passwordr} required />

        <button type="submit">Sign up</button>
        <a href="/auth/login">Already have an account?</a>
    </form>
</div>
