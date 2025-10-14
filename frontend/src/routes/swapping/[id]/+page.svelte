<script>
    // @ts-nocheck
    import { auth } from "$lib/stores/auth";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    let err = $state("");
    let takingLong = $state(false);

    onMount(async () => {
        // Get the user ID from route parameters
        const targetUserId = $page.params.id;
        setTimeout(() => {
            takingLong = true;
        }, 5000);
        if (!targetUserId) {
            console.error("No target user ID found in route parameters");
            await goto("/swapping");
            return;
        }

        console.log(
            `Creating chat between user ${$auth.user.id} and target user ${targetUserId}`,
        );

        try {
            const response = await fetch(
                `/api/createChat?u1=${$auth.user.id}&u2=${targetUserId}`,
            );

            if (response.ok) {
                const result = await response.json();
                console.log("Chat created successfully:", result);
                // Chat created successfully, redirect to swapping
                await goto("/swapping");
            } else {
                const errorText = await response.text();
                console.error(
                    "Failed to create chat:",
                    response.status,
                    errorText,
                );
                // Handle error - could show an error message or redirect to error page
                await goto("/swapping");
            }
        } catch (error) {
            err = error;

            // Handle network or other errors
            console.error("Error creating chat:", error);
            await goto("/swapping");
        }
    });
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
        {#if err}
            <div class="text-red-600">
                <p class="font-semibold">Error creating chat:</p>
                <p class="mt-2">{err}</p>
            </div>
        {:else}
            <div
                class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"
            ></div>
            <p class="mt-4 text-gray-600">Creating chat...</p>
            {#if takingLong}
                <p in:fade class="mt-4 text-gray-600">
                    It would appear that the request is taking longer than
                    expected.
                </p>
            {/if}
        {/if}
    </div>
</div>
