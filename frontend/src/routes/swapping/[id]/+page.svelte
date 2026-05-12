<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import { fade, fly } from "svelte/transition";
    import { ArrowRightLeft, ShieldAlert, CheckCircle2, XCircle, Loader2 } from "lucide-svelte";
    import type { User } from "$lib/types/user";

    let targetUser = $state<User | null>(null);
    let loading = $state(true);
    let error = $state("");
    let creating = $state(false);
    let swapSuccess = $state(false);

    onMount(async () => {
        if (!$auth.isAuthenticated) {
            goto("/auth/login");
            return;
        }

        const targetUserId = $page.params.id;
        if (!targetUserId) {
            error = "No user specified for swapping.";
            loading = false;
            return;
        }

        try {
            const response = await fetch(`/api/getUserInfo?q=${targetUserId}`);
            if (response.ok) {
                targetUser = await response.json();
            } else {
                error = "Could not find user information.";
            }
        } catch (err) {
            console.error("Failed to fetch target user:", err);
            error = "Network error while fetching user info.";
        } finally {
            loading = false;
        }
        
        try {
            const response = await fetch(
                `/api/createChat?u1=${$auth.user.id}&u2=${targetUser.id}`,
            );

            if (response.ok) {
                const result = await response.json();
                swapSuccess = true;
                // Brief delay to show success state
                setTimeout(() => {
                    goto("/swapping");
                }, 1500);
            } else {
                const errorData = await response.json();
                error = errorData.error || "Failed to initiate swap.";
                creating = false;
            }
        } catch (err) {
            console.error("Error creating chat:", err);
            error = "A connection error occurred. Please try again.";
            creating = false;
        }
    });

    async function handleConfirmSwap() {
        if (!targetUser || !$auth.user) return;
        
        creating = true;
        try {
            const response = await fetch(
                `/api/createChat?u1=${$auth.user.id}&u2=${targetUser.id}`,
            );

            if (response.ok) {
                const result = await response.json();
                swapSuccess = true;
                // Brief delay to show success state
                setTimeout(() => {
                    goto("/swapping");
                }, 1500);
            } else {
                const errorData = await response.json();
                error = errorData.error || "Failed to initiate swap.";
                creating = false;
            }
        } catch (err) {
            console.error("Error creating chat:", err);
            error = "A connection error occurred. Please try again.";
            creating = false;
        }
    }

    function handleCancel() {
        goto("/search");
    }
</script>

<svelte:head>
    <title>Confirm Swap - SkillSwap</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
        {#if error}
            <div class="text-red-600">
                <p class="font-semibold">Error creating chat:</p>
                <p class="mt-2">{error}</p>
            </div>
        {:else}
            <div
                class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"
            ></div>
            <p class="mt-4 text-gray-600">Creating chat...</p>
        {/if}
    </div>
</div>