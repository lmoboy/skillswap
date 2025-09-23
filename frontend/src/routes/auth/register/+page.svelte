<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { register } from "$lib/api/auth";

    let username = "";
    let email = "";
    let password = "";
    let error = "";
    let loading = false;

    async function handleSubmit(event: Event) {
        event.preventDefault();
        error = "";

        // Basic validation
        if (!username || !email || !password) {
            error = "All fields are required.";
            return;
        }

        if (username.length > 50) {
            error = "Username is too long (max 50 characters)";
            return;
        }

        if (email.length > 100) {
            error = "Email is too long (max 100 characters)";
            return;
        }

        if (password.length < 8) {
            error = "Password must be at least 8 characters long";
            return;
        }

        if (password.length > 50) {
            error = "Password is too long (max 50 characters)";
            return;
        }

        loading = true;

        try {
            await register({ username, email, password });
            // Redirect to login on success
            goto("/auth/login?registered=true");
        } catch (err: unknown) {
            error =
                err instanceof Error
                    ? err.message
                    : "Registration failed. Please try again.";
        } finally {
            loading = false;
        }
    }
</script>

<div
    class="w-full h-full flex justify-center items-center bg-white text-gray-800"
>
    <div class="w-full h-max max-w-md mx-auto p-4">
        <h2 class="text-center mb-6">Create Account</h2>

        {#if error}
            <div
                class="text-red-500 text-center mb-4 p-2 bg-red-50 rounded-lg border border-red-200"
            >
                {error}
            </div>
        {/if}

        <form on:submit={handleSubmit}>
            <label for="username" class="block text-center mb-2">Username</label
            >
            <input
                id="username"
                type="text"
                bind:value={username}
                required
                class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-100/60 text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition mb-4"
                autocomplete="username"
                disabled={loading}
                placeholder="Choose a unique username"
            />

            <label for="email" class="block text-center mb-2">Email</label>
            <input
                id="email"
                type="email"
                bind:value={email}
                required
                class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-100/60 text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition mb-4"
                autocomplete="email"
                disabled={loading}
                placeholder="Enter your email address"
            />

            <label for="password" class="block text-center mb-2">Password</label
            >
            <input
                id="password"
                type="password"
                bind:value={password}
                required
                class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-100/60 text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 transition mb-6"
                autocomplete="new-password"
                disabled={loading}
                placeholder="Create a strong password"
            />

            <button
                type="submit"
                disabled={loading}
                class="w-full p-3 rounded-lg bg-blue-500 hover:bg-blue-600 text-white font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed mb-4"
            >
                {#if loading}
                    Creating Account...
                {:else}
                    Create Account
                {/if}
            </button>

            <p class="text-center text-gray-600">
                Already have an account?
                <a
                    href="/auth/login"
                    class="text-blue-500 hover:text-blue-600 font-medium"
                    >Sign in</a
                >
            </p>
        </form>
    </div>
</div>
