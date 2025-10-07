<script lang="ts">
    import { Search, User, LogOut, Repeat } from "lucide-svelte";
    import { auth } from "$lib/stores/auth";
    import { logout } from "$lib/api/auth";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import Debug from "./components/Debug.svelte";
    import { slide } from "svelte/transition";
    import { quintOut } from "svelte/easing";
    import AuthDebug from "./components/AuthDebug.svelte";

    let searchQuery = $state("");
    let searching = $state(false);

    let searchResult = <
        {
            user: {
                username: string;
                email: string;
                id: string;
            };
            skills_found: string;
        }[]
    >$state([]);
    let showUserMenu = $state(false);
    let showSearchDropdown = $state(false); // New state variable

    async function handleLogout() {
        try {
            await logout();
            goto("/auth/login");
        } catch (error) {
            console.error("Logout failed:", error);
        }
    }

    function toggleUserMenu() {
        showUserMenu = !showUserMenu;
    }

    // A reference to the search container element
    let searchContainer: HTMLElement;

    // Use a unified click handler for both user menu and search dropdown
    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement;

        // Close user menu if click is outside the user menu container
        if (!target.closest(".user-menu")) {
            showUserMenu = false;
        }

        // Close search dropdown if click is outside the search container
        if (searchContainer && !searchContainer.contains(target)) {
            showSearchDropdown = false;
        }
    }

    let timeoutId: number | undefined = undefined;

    function handleSearch() {
        // Show dropdown immediately when user starts typing
        showSearchDropdown = true;
        clearTimeout(timeoutId);

        // Don't search for empty queries
        if (searchQuery.length === 0) {
            searching = false;
            searchResult = [];
            return;
        }

        searching = true;
        timeoutId = window.setTimeout(async () => {
            try {
                const response = await fetch(`/api/search`, {
                    method: "POST",
                    credentials: "include",
                    body: JSON.stringify({
                        query: searchQuery,
                    }),
                });
                const data = await response.json();
                searchResult = data;
                console.log(searchResult);
            } catch (error) {
                console.error("Search failed:", error);
            }
            searching = false;
        }, 100); // 100ms debounce
    }

    onMount(() => {
        document.addEventListener("click", handleClickOutside);
        return () => {
            document.removeEventListener("click", handleClickOutside);
        };
    });
</script>

<header class="bg-white border-b border-gray-200">
    <!-- <AuthDebug /> -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
            <div class="flex items-center space-x-8">
                <a href="/" class="text-xl font-semibold text-gray-800">
                    SkillSwap
                </a>
                <nav class="hidden md:flex items-center space-x-4">
                    <a
                        href="/course"
                        class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                    >
                        Courses
                    </a>
                </nav>
            </div>

            <div
                class="flex-1 flex items-center justify-center max-w-md w-full relative"
            >
                <div class="relative w-full" bind:this={searchContainer}>
                    <div
                        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                    >
                        <Search class="h-4 w-4 text-gray-400" />
                    </div>
                    <input
                        bind:value={searchQuery}
                        oninput={handleSearch}
                        type="text"
                        placeholder="Search skills or teachers..."
                        class="block w-full pl-10 pr-3 py-2 border text-gray-800 border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                        onfocus={() => (showSearchDropdown = true)}
                    />
                </div>

                {#if showSearchDropdown}
                    <div
                        class="absolute h-fit z-10 top-full mt-2 left-0 right-0 bg-white rounded-md shadow-lg border border-gray-200 max-h-80 overflow-y-auto"
                        transition:slide={{ duration: 300, easing: quintOut }}
                    >
                        {#if searching}
                            <div
                                class="flex items-center justify-center w-full py-4"
                            >
                                <div
                                    class="w-5 h-5 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"
                                ></div>
                            </div>
                        {:else if searchResult != null && searchResult.length > 0}
                            <ul class="divide-y h-fit divide-gray-200">
                                {#each searchResult as result (result.user.id)}
                                    <li
                                        class="p-2 hover:bg-gray-100 transition-colors duration-200 h-fit"
                                    >
                                        <a
                                            href={`/profile/${result.user.id}`}
                                            class="block"
                                        >
                                            <div
                                                class="flex items-center space-x-3"
                                            >
                                                <img
                                                    class="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0"
                                                    src={`https://via.placeholder.com/${result.user.id}`}
                                                    alt="Profile Picture"
                                                />
                                                <div class="flex-1 min-w-0">
                                                    <p
                                                        class="text-sm font-medium text-gray-900 truncate"
                                                    >
                                                        {result.user.username}
                                                    </p>
                                                    <p
                                                        class="text-xs text-gray-500 truncate"
                                                    >
                                                        Skills: {result.skills_found ||
                                                            "N/A"}
                                                    </p>
                                                </div>
                                            </div>
                                        </a>
                                    </li>
                                {/each}
                                <a
                                    href="/search?q={searchQuery}"
                                    class="block p-2"
                                >
                                    <div
                                        class="flex items-center space-x-3 hover:bg-gray-100 transition-colors duration-200 h-fit"
                                    >
                                        <div class="flex-1 min-w-0">
                                            <p
                                                class="text-sm font-medium text-gray-900 truncate"
                                            >
                                                See all results
                                            </p>
                                        </div>
                                    </div>
                                </a>
                            </ul>
                        {:else}
                            <div class="p-4 text-center text-sm text-gray-500">
                                No results found.
                            </div>
                        {/if}
                    </div>
                {/if}
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
                            onclick={toggleUserMenu}
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
                                transition:slide={{ easing: quintOut }}
                                class="absolute right-0 top-full mt-2 w-56 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 py-1 z-50 origin-top"
                                class:scale-y-0={!showUserMenu}
                                class:opacity-0={!showUserMenu}
                            >
                                <div
                                    class="px-4 py-2 border-b border-gray-100 dark:border-gray-700"
                                >
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
                                    href={`/profile/${$auth.user.id}`}
                                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                                >
                                    Your Profile
                                </a>
                                <a
                                    href="/settings"
                                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                                >
                                    Settings
                                </a>
                                <div
                                    class="border-t border-gray-100 my-1"
                                ></div>
                                <button
                                    onclick={handleLogout}
                                    class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200 flex items-center space-x-2"
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
