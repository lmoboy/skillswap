<script lang="ts">
    // Import components
    import Header from "$lib/components/layout/Header.svelte";
    import Footer from "$lib/components/layout/Footer.svelte";

    import "../app.css";
    import { auth } from "$lib/stores/auth";
    import Cookies from "$lib/components/ui/Cookies.svelte";

    // Get data from universal load function
    let { data } = $props();

    // Initialize auth store with data from load function
    if (data?.user) {
        auth.setUser({
            name: data.user.name,
            email: data.user.email,
            id: data.user.id,
            profile_picture: data.user.profile_picture || '',
        });
    } else {
        // No user found, set loading to false
        auth.setLoading(false);
    }
</script>

<div class="flex flex-col min-h-dvh relative">
    <Header class="shrink-0" />
    <main class="flex-1 w-full overflow-auto">
        <slot />
    </main>
    <Footer />
    <Cookies />
</div>
