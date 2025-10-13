<script lang="ts">
    import { register } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import Input from "$lib/components/common/Input.svelte";
    import Button from "$lib/components/common/Button.svelte";
    import { validateEmail, validatePassword, validateUsername } from "$lib/utils/validation";

    type Props = {
        onSuccess?: () => void;
        class?: string;
    };

    let {
        onSuccess,
        class: className = ""
    }: Props = $props();

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
            await register({ username, email, password });

            if (onSuccess) {
                onSuccess();
            } else {
                goto("/auth/login?registered=true");
            }
        } catch (err: unknown) {
            generalError = err instanceof Error ? err.message : "Registration failed. Please try again.";
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

        <Input
            bind:value={username}
            type="text"
            label="Username"
            placeholder="Choose a unique username"
            error={usernameError}
            required
            autocomplete="username"
            disabled={loading}
            maxlength={50}
            oninput={clearUsernameError}
        />

        <Input
            bind:value={email}
            type="email"
            label="Email"
            placeholder="Enter your email address"
            error={emailError}
            required
            autocomplete="email"
            disabled={loading}
            maxlength={100}
            oninput={clearEmailError}
        />

        <Input
            bind:value={password}
            type="password"
            label="Password"
            placeholder="Create a strong password (min 8 characters)"
            error={passwordError}
            required
            autocomplete="new-password"
            disabled={loading}
            maxlength={50}
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
            {loading ? "Creating Account..." : "Create Account"}
        </Button>

        <div class="text-center">
            <p class="text-sm text-gray-600">
                Already have an account?
                <a href="/auth/login" class="text-blue-600 hover:text-blue-700 font-medium">
                    Sign in
                </a>
            </p>
        </div>
    </form>
</div>
