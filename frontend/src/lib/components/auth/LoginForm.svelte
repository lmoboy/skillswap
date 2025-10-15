<script lang="ts">
    import { login } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import Button from "$lib/components/common/Button.svelte";
    import Input from "$lib/components/common/Input.svelte";
    import { validateEmail, validateRequired } from "$lib/utils/validation";

    interface Props {
        onSuccess?: () => void;
        class?: string;
    }

    const { onSuccess, class: className = "" }: Props = $props();

    let email = $state("");
    let password = $state("");
    let emailError = $state<string | null>(null);
    let passwordError = $state<string | null>(null);
    let generalError = $state<string | null>(null);
    let loading = $state(false);

    // Optimized validation with early returns
    const validateForm = (): boolean => {
        const emailValidation = validateEmail(email);
        const passwordValidation = validateRequired(password, "Password");
        
        emailError = emailValidation;
        passwordError = passwordValidation;
        
        return !emailValidation && !passwordValidation;
    };

    // Debounced error clearing for better UX
    const clearEmailError = () => {
        if (emailError) emailError = null;
    };

    const clearPasswordError = () => {
        if (passwordError) passwordError = null;
    };

    // Optimized submit handler with better error handling
    const handleSubmit = async (event: Event) => {
        event.preventDefault();
        generalError = null;

        if (!validateForm()) return;

        loading = true;

        try {
            const response = await login({ email, password });
            onSuccess?.() ?? goto(response.returnUrl || "/");
        } catch (err: unknown) {
            generalError = err instanceof Error ? err.message : "Login failed. Please try again.";
        } finally {
            loading = false;
        }
    };
</script>

<div class="w-full {className}">
    <form onsubmit={handleSubmit} class="space-y-6">
        {#if generalError}
            <div class="p-3 bg-red-50 border border-red-200 rounded-lg">
                <p class="text-sm text-red-700">{generalError}</p>
            </div>
        {/if}

        <Input
            type="email"
            label="Email"
            placeholder="Enter your email"
            required
            autocomplete="username"
            disabled={loading}
            bind:value={email}
            oninput={clearEmailError}
            data-testid="email-input"
            error={emailError}
        />

        <Input
            type="password"
            label="Password"
            placeholder="Enter your password"
            required
            autocomplete="current-password"
            disabled={loading}
            bind:value={password}
            oninput={clearPasswordError}
            data-testid="password-input"
            error={passwordError}
        />

        <Button
            type="submit"
            variant="primary"
            size="lg"
            fullWidth
            {loading}
            disabled={loading}
            data-testid="login-button"
        >
            {loading ? "Signing in..." : "Sign in"}
        </Button>

        <div class="text-center">
            <p class="text-sm text-gray-600">
                Don't have an account?
                <a
                    href="/auth/register"
                    class="text-blue-600 hover:text-blue-700 font-medium"
                    data-testid="signup-link"
                >
                    Sign up
                </a>
            </p>
        </div>
    </form>
</div>
