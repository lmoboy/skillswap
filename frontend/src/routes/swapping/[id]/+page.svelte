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

<div class="min-h-[calc(100vh-64px)] bg-gray-50 flex items-center justify-center p-4">
    {#if loading}
        <div class="text-center" in:fade>
            <Loader2 class="w-12 h-12 text-peach-600 animate-spin mx-auto mb-4" />
            <p class="text-gray-600 font-medium">Loading swap details...</p>
        </div>
    {:else if error}
        <div class="max-w-md w-full bg-white rounded-3xl shadow-xl p-8 text-center" in:fly={{ y: 20 }}>
            <div class="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <XCircle class="w-10 h-10 text-red-600" />
            </div>
            <h2 class="text-2xl font-bold text-gray-900 mb-2">Oops!</h2>
            <p class="text-gray-600 mb-8">{error}</p>
            <button 
                onclick={handleCancel}
                class="w-full py-4 bg-gray-900 text-white rounded-2xl font-bold hover:bg-gray-800 transition-all"
            >
                Go Back
            </button>
        </div>
    {:else if targetUser}
        <div class="max-w-lg w-full" in:fly={{ y: 20 }}>
            <div class="bg-white rounded-[2.5rem] shadow-2xl shadow-peach-200/20 overflow-hidden border border-gray-100">
                <!-- Top Header Decor -->
                <div class="h-32 bg-gradient-to-br from-peach-500 to-peach-600 relative">
                    <div class="absolute inset-0 opacity-20" style="background-image: radial-gradient(circle at 2px 2px, white 1px, transparent 0); background-size: 24px 24px;"></div>
                    <div class="absolute -bottom-12 left-1/2 -translate-x-1/2">
                        <div class="relative">
                            <img 
                                src={targetUser.profile_picture === "noPicture" ? "/default-avatar.svg" : `/api/profile/${targetUser.id}/picture`}
                                alt={targetUser.username}
                                class="w-24 h-24 rounded-3xl border-4 border-white object-cover bg-white shadow-lg"
                            />
                            {#if swapSuccess}
                                <div class="absolute -right-2 -bottom-2 bg-green-500 text-white p-1.5 rounded-full shadow-lg" in:scale>
                                    <CheckCircle2 class="w-5 h-5" />
                                </div>
                            {/if}
                        </div>
                    </div>
                </div>

                <div class="pt-16 pb-10 px-8 text-center">
                    <h2 class="text-2xl font-bold text-gray-900 mb-1">{targetUser.username}</h2>
                    <p class="text-peach-600 font-semibold mb-8">{targetUser.profession || 'Skill Swapper'}</p>

                    <div class="bg-peach-50/50 rounded-3xl p-6 mb-8 border border-peach-100/50">
                        <div class="flex items-center justify-center gap-3 mb-4 text-peach-700 font-bold">
                            <ArrowRightLeft class="w-5 h-5" />
                            <span>Swap Proposal</span>
                        </div>
                        
                        <p class="text-gray-600 text-sm leading-relaxed mb-4">
                            Initiating a swap with <strong>{targetUser.username}</strong> will open a direct chat and exchange connection details.
                        </p>

                        {#if ($auth.user?.swaps ?? 0) <= 0}
                            <div class="flex items-start gap-3 bg-red-50 p-4 rounded-2xl text-left border border-red-100">
                                <ShieldAlert class="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
                                <div>
                                    <p class="text-red-700 font-bold text-sm">No Swaps Available</p>
                                    <p class="text-red-600/80 text-xs">You need at least 1 swap credit to start a new connection.</p>
                                </div>
                            </div>
                        {:else}
                            <div class="flex items-center justify-between bg-white/60 p-4 rounded-2xl border border-peach-100/30">
                                <div class="text-left">
                                    <p class="text-gray-400 text-[10px] font-bold uppercase tracking-wider">Your Balance</p>
                                    <p class="text-gray-900 font-bold">{$auth.user?.swaps ?? 0} Credits</p>
                                </div>
                                <div class="text-right">
                                    <p class="text-gray-400 text-[10px] font-bold uppercase tracking-wider">Cost</p>
                                    <p class="text-peach-600 font-bold">-1 Credit</p>
                                </div>
                            </div>
                        {/if}
                    </div>

                    <div class="flex flex-col gap-3">
                        {#if swapSuccess}
                            <div class="w-full py-4 bg-green-500 text-white rounded-2xl font-bold flex items-center justify-center gap-2" in:fade>
                                <CheckCircle2 class="w-5 h-5" />
                                Connection Secured!
                            </div>
                        {:else}
                            <button 
                                onclick={handleConfirmSwap}
                                disabled={creating || ($auth.user?.swaps ?? 0) <= 0}
                                class="w-full py-4 bg-peach-600 text-white rounded-2xl font-bold hover:bg-peach-700 transition-all shadow-lg shadow-peach-600/20 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                            >
                                {#if creating}
                                    <Loader2 class="w-5 h-5 animate-spin" />
                                    Processing...
                                {:else}
                                    Confirm & Start Swap
                                {/if}
                            </button>
                            
                            <button 
                                onclick={handleCancel}
                                disabled={creating}
                                class="w-full py-4 text-gray-500 font-semibold hover:text-gray-800 transition-colors disabled:opacity-50"
                            >
                                Not now, thanks
                            </button>
                        {/if}
                    </div>
                </div>
            </div>
            
            <p class="mt-8 text-center text-gray-400 text-xs px-8">
                By confirming, you agree to our community guidelines. Swaps are non-refundable once the connection is established.
            </p>
        </div>
    {/if}
</div>

<style>
    :global(body) {
        background-color: #f9fafb;
    }
    
    .shadow-peach-200\/20 {
        --tw-shadow-color: rgba(242, 100, 68, 0.2);
    }
    
    .bg-peach-500 { background-color: #f26444; }
    .bg-peach-600 { background-color: #e85a3a; }
    .text-peach-600 { color: #f26444; }
    .text-peach-700 { color: #d44d2e; }
    .bg-peach-50 { background-color: #fff7f5; }
    .border-peach-100 { border-color: #ffedea; }
</style>
