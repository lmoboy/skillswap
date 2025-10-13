<script lang="ts">
    import LoginForm from "$lib/components/auth/LoginForm.svelte";
</script>

<div
    class="w-full h-full flex justify-center items-center bg-white text-gray-800 min-h-screen"
>
    <div class="w-full h-max max-w-md mx-auto p-4">
        <div class="text-center mb-8">
            <h2 class="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h2>
            <p class="text-gray-600">Sign in to your account to continue</p>
        </div>

        <LoginForm />
    </div>
</div>
