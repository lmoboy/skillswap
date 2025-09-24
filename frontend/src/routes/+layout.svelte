<script lang="ts">
  // Import components
  import Header from "$lib/Header.svelte";
  import Footer from "$lib/Footer.svelte";

  import "../app.css";
  import { onMount } from "svelte";
  import { checkAuth } from "$lib/api/auth";
  import Cookies from "$lib/Cookies.svelte";
  // import AuthDebug from "$lib/components/AuthDebug.svelte";

  onMount(async () => {
    try {
      await checkAuth();
    } catch (error) {
      console.error("Error checking auth:", error);
    }
  });
</script>

<div class="flex flex-col min-h-screen relative">
  <!-- <AuthDebug /> -->

  <Header />
  <main class="flex-grow">
    <slot />
  </main>
  <Footer />
  <Cookies />
</div>
