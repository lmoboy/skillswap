<script lang="ts">
    import RegisterForm from "$lib/components/auth/RegisterForm.svelte";
</script>

<div
    class="w-full h-full flex justify-center items-center bg-white text-gray-800 min-h-screen"
>
    <div class="w-full h-max max-w-md mx-auto p-4">
        <div class="text-center mb-8">
            <h2 class="text-3xl font-bold text-gray-900 mb-2">
                Create Account
            </h2>
            <p class="text-gray-600">Join our community and start learning</p>
        </div>

        <RegisterForm />
    </div>
</div>
