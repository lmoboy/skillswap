<script lang="ts">
    import { register } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import Button from "$lib/components/common/Button.svelte";
    import {
        validateEmail,
        validatePassword,
        validateUsername,
    } from "$lib/utils/validation";

    type Props = {
        onSuccess?: () => void;
        class?: string;
    };

    let { onSuccess, class: className = "" }: Props = $props();

    let username = $state("");
    let email = $state("");
    let password = $state("");
    let usernameError = $state<string | null>(null);
    let emailError = $state<string | null>(null);
    let passwordError = $state<string | null>(null);
    let generalError = $state<string | null>(null);
    let loading = $state(false);

    function validateForm(): boolean {
        usernameError = validateUsername(username);
        emailError = validateEmail(email);
        passwordError = validatePassword(password);

        return !usernameError && !emailError && !passwordError;
    }

    async function handleSubmit(event: Event) {
        event.preventDefault();
        generalError = null;

        if (!validateForm()) {
            return;
        }

        loading = true;

        try {
            const response = await register({ username, email, password });

            if (onSuccess) {
                onSuccess();
            } else {
                // Redirect to the return URL or default to home
                const redirectUrl = response.returnUrl || "/";
                goto(redirectUrl);
            }
        } catch (err: unknown) {
            generalError =
                err instanceof Error
                    ? err.message
                    : "Registration failed. Please try again.";
        } finally {
            loading = false;
        }
    }

    function clearUsernameError() {
        usernameError = null;
    }

    function clearEmailError() {
        emailError = null;
    }

    function clearPasswordError() {
        passwordError = null;
    }
</script>

<div class="w-full {className}">
    <form onsubmit={handleSubmit} class="space-y-6">
        {#if generalError}
            <div class="p-3 bg-red-50 border border-red-200 rounded-lg">
                <p class="text-sm text-red-700">{generalError}</p>
            </div>
        {/if}

        <div class="w-full">
            <label
                for="username"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                Username
                <span class="text-red-500">*</span>
            </label>
            <input
                id="username"
                type="text"
                placeholder="Choose a unique username"
                required
                autocomplete="username"
                disabled={loading}
                maxlength={50}
                bind:value={username}
                oninput={clearUsernameError}
                data-testid="username-input"
                class="w-full p-3 rounded-lg border bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 transition border-gray-300 focus:ring-blue-500 focus:border-blue-500 {loading
                    ? 'opacity-50 cursor-not-allowed'
                    : ''}"
            />
            {#if usernameError}
                <p class="mt-1 text-sm text-red-600">
                    {usernameError}
                </p>
            {/if}
        </div>

        <div class="w-full">
            <label
                for="email"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                Email
                <span class="text-red-500">*</span>
            </label>
            <input
                id="email"
                type="email"
                placeholder="Enter your email address"
                required
                autocomplete="email"
                disabled={loading}
                maxlength={100}
                bind:value={email}
                oninput={clearEmailError}
                data-testid="email-input"
                class="w-full p-3 rounded-lg border bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 transition border-gray-300 focus:ring-blue-500 focus:border-blue-500 {loading
                    ? 'opacity-50 cursor-not-allowed'
                    : ''}"
            />
            {#if emailError}
                <p class="mt-1 text-sm text-red-600">
                    {emailError}
                </p>
            {/if}
        </div>

        <div class="w-full">
            <label
                for="password"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                Password
                <span class="text-red-500">*</span>
            </label>
            <input
                id="password"
                type="password"
                placeholder="Create a strong password (min 8 characters)"
                required
                autocomplete="new-password"
                disabled={loading}
                maxlength={50}
                bind:value={password}
                oninput={clearPasswordError}
                data-testid="password-input"
                class="w-full p-3 rounded-lg border bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 transition border-gray-300 focus:ring-blue-500 focus:border-blue-500 {loading
                    ? 'opacity-50 cursor-not-allowed'
                    : ''}"
            />
            {#if passwordError}
                <p class="mt-1 text-sm text-red-600">
                    {passwordError}
                </p>
            {/if}
        </div>

        <Button
            type="submit"
            variant="primary"
            size="lg"
            fullWidth
            {loading}
            disabled={loading}
            data-testid="register-button"
        >
            {loading ? "Creating Account..." : "Create Account"}
        </Button>

        <div class="text-center">
            <p class="text-sm text-gray-600">
                Already have an account?
                <a
                    href="/auth/login"
                    class="text-blue-600 hover:text-blue-700 font-medium"
                >
                    Sign in
                </a>
            </p>
        </div>
    </form>
</div>
