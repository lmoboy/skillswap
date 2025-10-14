<script lang="ts">
    // Import components
    import Header from "$lib/components/layout/Header.svelte";
    import Footer from "$lib/components/layout/Footer.svelte";

    import "../app.css";
    import { onMount } from "svelte";
    import { checkAuth } from "$lib/api/auth";
    import Cookies from "$lib/components/ui/Cookies.svelte";

    onMount(async () => {
        try {
            await checkAuth();
        } catch (error) {
            console.error("Error checking auth:", error);
        }
    });
</script>

<div class="flex flex-col min-h-dvh relative">
    <Header class="shrink-0" />
    <main class="flex-1 w-full overflow-auto">
        <slot />
    </main>
    <Footer />
    <Cookies />
</div>
