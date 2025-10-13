<script lang="ts">
    import { Search } from "lucide-svelte";

    type SearchResult = {
        user: {
            username: string;
            email: string;
            id: string;
        };
        skills_found: string;
    };

    type Props = {
        placeholder?: string;
        onSearch?: (query: string) => void;
        showDropdown?: boolean;
        results?: SearchResult[];
        loading?: boolean;
        class?: string;
    };

    let {
        placeholder = "Search skills or teachers...",
        onSearch,
        showDropdown = false,
        results = [],
        loading = false,
        class: className = "",
    }: Props = $props();

    let searchQuery = $state("");
    let searchContainer: HTMLElement;
    let showResults = $state(false);

    function handleInput() {
        showResults = true;
        if (onSearch) {
            onSearch(searchQuery);
        }
    }

    function handleFocus() {
        if (searchQuery.length > 0) {
            showResults = true;
        }
    }

    function closeDropdown() {
        showResults = false;
    }

    // Export methods if needed
    export function clear() {
        searchQuery = "";
        showResults = false;
    }

    export function getValue() {
        return searchQuery;
    }
</script>

<svelte:window
    onclick={(e) => {
        const target = e.target as Node;
        if (searchContainer && !searchContainer.contains(target)) {
            closeDropdown();
        }
    }}
/>

<div class="relative w-full {className}" bind:this={searchContainer}>
    <div class="relative">
        <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
        >
            <Search class="h-4 w-4 text-gray-400" />
        </div>
        <input
            bind:value={searchQuery}
            oninput={handleInput}
            onfocus={handleFocus}
            type="text"
            {placeholder}
            class="block w-full pl-10 pr-3 py-2 border text-gray-800 border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        />
    </div>

    {#if showDropdown && showResults && searchQuery.length > 0}
        <div
            class="absolute z-10 top-full mt-2 left-0 right-0 bg-white rounded-md shadow-lg border border-gray-200 max-h-80 overflow-y-auto"
        >
            {#if loading}
                <div class="flex items-center justify-center w-full py-4">
                    <div
                        class="w-5 h-5 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"
                    ></div>
                </div>
            {:else if results && results.length > 0}
                <ul class="divide-y divide-gray-200">
                    {#each results as result (result.user.id)}
                        <li
                            class="p-2 hover:bg-gray-100 transition-colors duration-200"
                        >
                            <a
                                href={`/profile/${result.user.id}`}
                                class="block"
                            >
                                <div class="flex items-center space-x-3">
                                    <img
                                        class="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0"
                                        src={`/api/profile/${result.user.id}/picture`}
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
                    <li>
                        <a
                            href="/search?q={searchQuery}"
                            class="block p-2 hover:bg-gray-100"
                        >
                            <div class="flex items-center space-x-3">
                                <div class="flex-1 min-w-0">
                                    <p
                                        class="text-sm font-medium text-gray-900 truncate"
                                    >
                                        See all results
                                    </p>
                                </div>
                            </div>
                        </a>
                    </li>
                </ul>
            {:else}
                <div class="p-4 text-center text-sm text-gray-500">
                    No results found.
                </div>
            {/if}
        </div>
    {/if}
</div>
