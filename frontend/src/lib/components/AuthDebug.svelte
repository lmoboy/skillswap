<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { onMount, onDestroy } from "svelte";

    let authState: {
        user: { name: string; email: string } | null;
        isAuthenticated: boolean;
        loading: boolean;
        error: string | null;
    } = {
        user: null,
        isAuthenticated: false,
        loading: true,
        error: null,
    };

    const unsubscribe = auth.subscribe((state) => {
        authState = state;
    });

    onDestroy(() => {
        unsubscribe();
    });
</script>

<div
    class="fixed bottom-4 right-4 bg-white p-4 rounded-lg shadow-lg text-gray-600 text-sm max-w-xs"
>
    <h3 class="font-bold mb-2">Auth State</h3>
    <div class="space-y-2">
        <div class="flex justify-between">
            <span class="font-medium">Status:</span>
            <span class="font-mono">
                {authState.isAuthenticated
                    ? "✅ Authenticated"
                    : "❌ Not Authenticated"}
            </span>
        </div>

        {#if authState.loading}
            <div class="text-blue-600">Loading...</div>
        {/if}

        {#if authState.error}
            <div class="text-red-600">Error: {authState.error}</div>
        {/if}

        {#if authState.user}
            <div class="mt-2 pt-2 border-t">
                <div class="font-medium">User:</div>
                <div class="ml-2">
                    <div>Name: {authState.user.name}</div>
                    <div>Email: {authState.user.email}</div>
                </div>
            </div>
        {/if}
    </div>
</div>
