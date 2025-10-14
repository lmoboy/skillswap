<script lang="ts">
    import LoginForm from "$lib/components/auth/LoginForm.svelte";
    import { auth } from "$lib/stores/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    onMount(() => {
        // Redirect if already authenticated
        if ($auth.isAuthenticated) {
            goto("/");
        }
    });

    // Watch for auth changes and redirect if user logs in
    $effect(() => {
        if ($auth.isAuthenticated && !$auth.loading) {
            goto("/");
        }
    });
</script>

<div
    class="w-full h-full flex justify-center items-center bg-white text-gray-800 min-h-screen"
>
    <div class="w-full h-max max-w-md mx-auto p-4 sm:p-6">
        <div class="text-center mb-6 sm:mb-8">
            <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                Welcome Back
            </h2>
            <p class="text-sm sm:text-base text-gray-600">
                Sign in to your account to continue
            </p>
        </div>

        <LoginForm />
    </div>
</div>
