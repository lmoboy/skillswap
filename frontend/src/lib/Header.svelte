<script lang="ts">
    import { Search, User, LogOut, Repeat } from "lucide-svelte";
    import { auth } from "$lib/stores/auth";
    import { logout } from "$lib/api/auth";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    let inputRef: HTMLInputElement;
    let showUserMenu = false;

    async function handleLogout() {
        try {
            await logout();
            goto("/auth/login");
        } catch (error) {
            console.error("Logout failed:", error);
        }
    }

    // Toggle user menu
    function toggleUserMenu() {
        showUserMenu = !showUserMenu;
    }

    // Close menu when clicking outside
    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement;
        if (!target.closest(".user-menu")) {
            showUserMenu = false;
        }
    }

    onMount(() => {
        // Add click outside listener
        document.addEventListener("click", handleClickOutside);
        return () => {
            document.removeEventListener("click", handleClickOutside);
        };
    });
</script>

<header class="bg-white border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
            <div class="flex items-center">
                <a href="/" class="text-xl font-semibold text-gray-800">
                    SkillSwap
                </a>
            </div>
            <div
                class="flex-1 flex items-center justify-center max-w-md h-full"
            >
                <div class="relative">
                    <div
                        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                    >
                        <Search class="h-4 w-4 text-gray-400" />
                    </div>
                    <input
                        bind:this={inputRef}
                        type="text"
                        placeholder="Search skills or teachers..."
                        class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                    />
                </div>
            </div>
            <nav class="flex items-center space-x-4">
                {#if $auth.loading && !$auth.isAuthenticated}
                    <div class="flex items-center justify-center w-5 h-5">
                        <div
                            class="w-4 h-4 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"
                        ></div>
                    </div>
                {:else if $auth.isAuthenticated && $auth.user}
                    <div class="relative user-menu">
                        <button
                            on:click={toggleUserMenu}
                            class="flex items-center space-x-2 text-sm text-gray-700 hover:bg-gray-100 rounded-full p-1 pr-3 transition-colors"
                            aria-label="User menu"
                            aria-expanded={showUserMenu}
                        >
                            <div
                                class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-600"
                            >
                                <User size={16} />
                            </div>
                            <span class="font-medium">{$auth.user.name}</span>
                        </button>

                        {#if showUserMenu}
                            <div
                                class="absolute right-0 mt-2 w-56 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 py-1 z-50"
                            >
                                <div class="px-4 py-2 border-b border-gray-100">
                                    <p
                                        class="text-sm font-medium text-gray-900"
                                    >
                                        {$auth.user.name}
                                    </p>
                                    <p class="text-xs text-gray-500 truncate">
                                        {$auth.user.email}
                                    </p>
                                </div>
                                <a
                                    href="/profile"
                                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                                >
                                    Your Profile
                                </a>
                                <a
                                    href="/settings"
                                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                                >
                                    Settings
                                </a>
                                <div
                                    class="border-t border-gray-100 my-1"
                                ></div>
                                <button
                                    on:click={handleLogout}
                                    class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center space-x-2"
                                >
                                    <LogOut size={14} class="text-gray-400" />
                                    <span>Sign out</span>
                                </button>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <a
                        href="/auth/login"
                        class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                    >
                        Sign in
                    </a>
                    <a
                        href="/auth/register"
                        class="ml-4 bg-gray-100 text-gray-800 hover:bg-gray-200 px-4 py-2 rounded-md text-sm font-medium transition-colors"
                    >
                        Sign up
                    </a>
                {/if}
            </nav>
        </div>
    </div>
</header>
