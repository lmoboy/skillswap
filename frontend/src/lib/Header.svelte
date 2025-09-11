<script lang="ts">
    import { Search } from "lucide-svelte";
    import { onMount } from "svelte";
    // import AuthManager from "./authManager";
    // let user = $state(new AuthManager().getUserInfo());

    let hover = false;
    let inputRef: HTMLInputElement;
    function handleSearch(): any {
        console.log("big black chicken");
    }
</script>

<header class="relative shadow-lg py-4 px-8 bg-transparent">
    <div class="absolute inset-0 bg-black bg-opacity-60"></div>
    <div
        class="relative z-10 flex justify-between items-center max-w-5xl mx-auto"
    >
        <div class="text-2xl font-extrabold text-white tracking-wide">
            <a href="/">SkillSwap</a>
        </div>
        <button
            type="button"
            aria-label="Focus search input"
            on:click={() => inputRef && inputRef.focus()}
            on:keydown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    inputRef && inputRef.focus();
                }
            }}
            class="flex items-center gap-4 border border-white/30 bg-white/10 backdrop-blur-lg rounded-full px-4 py-2 transition-all duration-300 cursor-text"
        >
            <Search class="w-5 h-5 text-white" />
            <input
                bind:this={inputRef}
                on:submit={handleSearch()}
                type="text"
                placeholder="Search skills or teachers..."
                class="transition-all outline-none p-0 m-0 w-0 left-0 focus:w-fit text-white placeholder-white"
            />
        </button>
        <nav class="flex gap-8">
            <!-- {#if !user.name || !user.email}
                <a
                    href="/auth/login"
                    class="text-gray-200 hover:text-white font-semibold transition"
                    >New here?</a
                >
            {:else}
                <a
                    href="/profile"
                    class="text-gray-200 hover:text-white font-semibold transition"
                    >{Auth.getUserInfo().name}</a
                >
            {/if} -->
        </nav>
    </div>
</header>
