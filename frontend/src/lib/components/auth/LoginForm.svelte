<script lang="ts">
    import { login } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import Input from "$lib/components/common/Input.svelte";
    import Button from "$lib/components/common/Button.svelte";
    import { validateEmail, validateRequired } from "$lib/utils/validation";

    type Props = {
        onSuccess?: () => void;
        class?: string;
    };

    let {
        onSuccess,
        class: className = ""
    }: Props = $props();

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
            generalError = err instanceof Error ? err.message : "Login failed. Please try again.";
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

        <Input
            bind:value={email}
            type="email"
            label="Email"
            placeholder="Enter your email"
            error={emailError}
            required
            autocomplete="username"
            disabled={loading}
            oninput={clearEmailError}
        />

        <Input
            bind:value={password}
            type="password"
            label="Password"
            placeholder="Enter your password"
            error={passwordError}
            required
            autocomplete="current-password"
            disabled={loading}
            oninput={clearPasswordError}
        />

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
                <a href="/auth/register" class="text-blue-600 hover:text-blue-700 font-medium">
                    Sign up
                </a>
            </p>
        </div>
    </form>
</div>
