<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { goto } from "$app/navigation";
    import { login } from "$lib/api/auth";
    import { auth } from "$lib/stores/auth";

    let email = "";
    let password = "";
    let error = "";
    let loading = false;

    // Subscribe to auth changes
    const unsubscribe = auth.subscribe((state) => {
        // If user becomes authenticated, redirect to home
        if (state.isAuthenticated) {
            goto("/");
        }
        // Handle loading state
        loading = state.loading;

        // Handle errors
        if (state.error) {
            error = state.error;
        }
    });

    // Clean up subscription on component destroy
    onDestroy(() => {
        unsubscribe();
    });

    async function handleSubmit() {
        // Reset error state
        error = "";

        // Basic validation
        if (!email || !password) {
            error = "Email and password are required.";
            return;
        }
        if (email.length > 100) {
            error = "Email is too long!";
            return;
        }
        if (password.length > 50) {
            error = "Password cannot be longer than 50 characters!";
            return;
        }

        try {
            // Call the login API
            await login({ email, password });
            // The auth store subscription will handle the redirect on success
        } catch (err: unknown) {
            // Error is already handled by the auth store
            console.error("Login error:", err);
        }
    }
</script>

<div class="login-container">
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
            class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-900/60 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
            autocomplete="username"
            disabled={loading}
        />

        <label for="password">Password</label>
        <input
            id="password"
            type="password"
            bind:value={password}
            required
            autocomplete="current-password"
            disabled={loading}
        />

        <button type="submit" disabled={loading}>
            {#if loading}
                Signing in...
            {:else}
                Sign in
            {/if}
        </button>
        <a href="/auth/register">Don't have an account?</a>
    </form>
</div>
