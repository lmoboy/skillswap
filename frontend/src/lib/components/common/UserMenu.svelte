<script lang="ts">
    import { LogOut, User } from "lucide-svelte";
    import { slide } from "svelte/transition";
    import { quintOut } from "svelte/easing";

    type UserData = {
        id: string | number;
        name: string;
        email: string;
        profile_picture?: string;
    };

    type Props = {
        user: UserData;
        onLogout?: () => void;
        class?: string;
    };

    let {
        user,
        onLogout,
        class: className = ""
    }: Props = $props();

    let showMenu = $state(false);
    let menuContainer: HTMLElement;

    function toggleMenu() {
        showMenu = !showMenu;
    }

    function closeMenu() {
        showMenu = false;
    }

    function handleLogout() {
        if (onLogout) {
            onLogout();
        }
        closeMenu();
    }
</script>

<svelte:window
    onclick={(e) => {
        if (menuContainer && !menuContainer.contains(e.target as Node)) {
            closeMenu();
        }
    }}
/>

<div class="relative user-menu {className}" bind:this={menuContainer}>
    <button
        onclick={toggleMenu}
        class="flex items-center space-x-2 text-sm text-gray-700 hover:bg-gray-100 rounded-full p-1 pr-3 transition-colors"
        aria-label="User menu"
        aria-expanded={showMenu}
    >
        {#if user.profile_picture}
            <img
                src={user.profile_picture}
                alt={user.name}
                class="w-8 h-8 rounded-full object-cover"
            />
        {:else}
            <div
                class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-600"
            >
                <User size={16} />
            </div>
        {/if}
        <span class="font-medium">{user.name}</span>
    </button>

    {#if showMenu}
        <div
            transition:slide={{ duration: 200, easing: quintOut }}
            class="absolute right-0 top-full mt-2 w-56 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 py-1 z-50"
        >
            <div class="px-4 py-2 border-b border-gray-100">
                <p class="text-sm font-medium text-gray-900">
                    {user.name}
                </p>
                <p class="text-xs text-gray-500 truncate">
                    {user.email}
                </p>
            </div>
            <a
                href={`/profile/${user.id}`}
                class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                onclick={closeMenu}
            >
                Your Profile
            </a>
            <a
                href="/settings"
                class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200"
                onclick={closeMenu}
            >
                Settings
            </a>
            <div class="border-t border-gray-100 my-1"></div>
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
