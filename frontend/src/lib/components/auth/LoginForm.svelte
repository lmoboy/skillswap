<script lang="ts">
    import { login } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import Button from "$lib/components/common/Button.svelte";
    import { validateEmail, validateRequired } from "$lib/utils/validation";

    type Props = {
        onSuccess?: () => void;
        class?: string;
    };

    let { onSuccess, class: className = "" }: Props = $props();

    let email = $state("");
    let password = $state("");
    let emailError = $state<string | null>(null);
    let passwordError = $state<string | null>(null);
    let generalError = $state<string | null>(null);
    let loading = $state(false);

    function validateForm(): boolean {
        emailError = validateEmail(email);
        passwordError = validateRequired(password, "Password");

        return !emailError && !passwordError;
    }

    async function handleSubmit(event: Event) {
        event.preventDefault();
        generalError = null;

        if (!validateForm()) {
            return;
        }

        loading = true;

        try {
            await login({ email, password });

            if (onSuccess) {
                onSuccess();
            } else {
                goto("/");
            }
        } catch (err: unknown) {
            generalError =
                err instanceof Error
                    ? err.message
                    : "Login failed. Please try again.";
        } finally {
            loading = false;
        }
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
                for="email"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                Email
                <span class="text-red-500">*</span>
            </label>
            <input
                id="email"
                type="email"
                placeholder="Enter your email"
                required
                autocomplete="username"
                disabled={loading}
                bind:value={email}
                oninput={clearEmailError}
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
                placeholder="Enter your password"
                required
                autocomplete="current-password"
                disabled={loading}
                bind:value={password}
                oninput={clearPasswordError}
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
        >
            {loading ? "Signing in..." : "Sign in"}
        </Button>

        <div class="text-center">
            <p class="text-sm text-gray-600">
                Don't have an account?
                <a
                    href="/auth/register"
                    class="text-blue-600 hover:text-blue-700 font-medium"
                >
                    Sign up
                </a>
            </p>
        </div>
    </form>
</div>
