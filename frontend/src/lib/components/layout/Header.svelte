<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { logout } from "$lib/api/auth";
    import { goto } from "$app/navigation";
    import SearchBar from "$lib/components/common/SearchBar.svelte";
    import UserMenu from "$lib/components/common/UserMenu.svelte";

    let searching = $state(false);
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

<header class="bg-white border-b border-gray-200">
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
                class="flex-1 flex items-center justify-center max-w-md w-full"
            >
                <SearchBar
                    showDropdown={true}
                    results={searchResult}
                    loading={searching}
                    onSearch={handleSearch}
                />
            </div>

            <nav class="flex items-center space-x-4">
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
                        class="ml-4 bg-gray-100 text-gray-800 hover:bg-gray-200 px-4 py-2 rounded-md text-sm font-medium transition-colors"
                    >
                        Sign up
                    </a>
                {/if}
            </nav>
        </div>
    </div>
</header>
