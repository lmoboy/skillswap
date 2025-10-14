<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { logout } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import SearchBar from "$lib/components/common/SearchBar.svelte";
    import UserMenu from "$lib/components/common/UserMenu.svelte";
    import { Menu, X } from "lucide-svelte";

    type Props = {
        class?: string;
    };

    let { class: className = "" }: Props = $props();

    let searching = $state(false);
    let mobileMenuOpen = $state(false);
    let searchResult = $state<
        {
            user: {
                username: string;
                email: string;
                id: string;
            };
            skills_found: string;
        }[]
    >([]);

    let timeoutId: number | undefined = undefined;

    async function handleLogout() {
        try {
            await logout();
            goto("/auth/login");
        } catch (error) {
            console.error("Logout failed:", error);
        }
    }

    function handleSearch(query: string) {
        clearTimeout(timeoutId);

        if (query.length === 0) {
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
                        query: query,
                    }),
                });
                const data = await response.json();
                searchResult = data;
            } catch (error) {
                console.error("Search failed:", error);
            }
            searching = false;
        }, 100);
    }
</script>

<header class="bg-white border-b border-gray-200 {className}">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
            <!-- Logo -->
            <div class="flex items-center">
                <a href="/" class="text-xl font-semibold text-gray-800">
                    SkillSwap
                </a>
            </div>

            <!-- Desktop Navigation -->
            <nav class="hidden md:flex items-center space-x-4">
                <a
                    href="/course"
                    class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                >
                    Courses
                </a>
            </nav>
            <nav class="hidden md:flex items-center space-x-4">
                <a
                    href="/swapping"
                    class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                >
                    Swapping
                </a>
            </nav>
            <!-- Desktop Search -->
            <div
                class="hidden md:flex flex-1 items-center justify-center max-w-md mx-4"
            >
                <SearchBar
                    showDropdown={true}
                    results={searchResult}
                    loading={searching}
                    onSearch={handleSearch}
                />
            </div>

            <!-- Desktop Auth -->
            <nav class="hidden md:flex items-center space-x-4">
                {#if $auth.loading && !$auth.isAuthenticated}
                    <div class="flex items-center justify-center w-5 h-5">
                        <div
                            class="w-4 h-4 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"
                        ></div>
                    </div>
                {:else if $auth.isAuthenticated && $auth.user}
                    <UserMenu user={$auth.user} onLogout={handleLogout} />
                {:else}
                    <a
                        href="/auth/login"
                        class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                    >
                        Sign in
                    </a>
                    <a
                        href="/auth/register"
                        class="bg-gray-100 text-gray-800 hover:bg-gray-200 px-4 py-2 rounded-md text-sm font-medium transition-colors"
                    >
                        Sign up
                    </a>
                {/if}
            </nav>

            <!-- Mobile menu button -->
            <div class="flex md:hidden">
                <button
                    onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
                    class="inline-flex items-center justify-center p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100 focus:outline-none"
                >
                    {#if mobileMenuOpen}
                        <X class="h-6 w-6" />
                    {:else}
                        <Menu class="h-6 w-6" />
                    {/if}
                </button>
            </div>
        </div>

        <!-- Mobile menu -->
        {#if mobileMenuOpen}
            <div class="md:hidden border-t border-gray-200 py-4 space-y-4">
                <!-- Mobile Search -->
                <div class="px-2">
                    <SearchBar
                        showDropdown={true}
                        results={searchResult}
                        loading={searching}
                        onSearch={handleSearch}
                    />
                </div>

                <!-- Mobile Navigation -->
                <nav class="flex flex-col space-y-2">
                    <a
                        href="/course"
                        class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                        onclick={() => (mobileMenuOpen = false)}
                    >
                        Courses
                    </a>
                </nav>

                <!-- Mobile Auth -->
                <div
                    class="border-t border-gray-200 pt-4 flex flex-col space-y-2"
                >
                    {#if $auth.loading && !$auth.isAuthenticated}
                        <div
                            class="flex items-center justify-center w-full py-2"
                        >
                            <div
                                class="w-4 h-4 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"
                            ></div>
                        </div>
                    {:else if $auth.isAuthenticated && $auth.user}
                        <a
                            href={`/profile/${$auth.user.id}`}
                            class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                            onclick={() => (mobileMenuOpen = false)}
                        >
                            Profile
                        </a>
                        <a
                            href="/settings"
                            class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                            onclick={() => (mobileMenuOpen = false)}
                        >
                            Settings
                        </a>
                        <button
                            onclick={() => {
                                handleLogout();
                                mobileMenuOpen = false;
                            }}
                            class="text-left text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors"
                        >
                            Sign out
                        </button>
                    {:else}
                        <a
                            href="/auth/login"
                            class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium hover:bg-gray-50 rounded-md transition-colors text-center"
                            onclick={() => (mobileMenuOpen = false)}
                        >
                            Sign in
                        </a>
                        <a
                            href="/auth/register"
                            class="bg-gray-900 text-white hover:bg-gray-800 px-4 py-2 rounded-md text-sm font-medium transition-colors text-center"
                            onclick={() => (mobileMenuOpen = false)}
                        >
                            Sign up
                        </a>
                    {/if}
                </div>
            </div>
        {/if}
    </div>
</header>
