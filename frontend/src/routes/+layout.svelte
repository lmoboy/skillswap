<!-- <script lang="ts">
  import { onMount } from "svelte";
  let { children } = $props();
  import "../app.css";
  let darkMode = $state(false);
  let menuOpen = $state(false);

  function toggleDarkMode() {
    darkMode = !darkMode;

    localStorage.setItem("darkMode", darkMode.toString());
  }

  function toggleMenu() {
    menuOpen = !menuOpen;
  }
  // $effect(() => {
  //   if (localStorage.getItem("darkMode") === "true") {
  //     darkMode = true;
  //   } else {
  //     darkMode = false;
  //   }
  //   console.log("Dark mode is now:", darkMode);
  // });
  onMount(() => {
    if (localStorage.getItem("darkMode") === "true") {
      darkMode = true;
    } else {
      darkMode = false;
    }
    console.log("Dark mode is now:", darkMode);
  });
</script>

<div
  class="w-full min-h-screen flex flex-col items-center justify-center duration-300 relative overflow-hidden {darkMode
    ? 'dark bg-gray-900 text-white'
    : 'bg-gray-100 text-gray-900'}"
>
  <main class="p-8 max-w-xl text-center">
    {@render children?.()}
  </main>

  <button
    class="absolute top-4 left-4 p-2 rounded-full shadow-lg transition-colors duration-300
           {darkMode
      ? 'bg-gray-700 text-white hover:bg-gray-600'
      : 'bg-white text-gray-900 hover:bg-gray-200'}"
    onclick={toggleDarkMode}
  >
    {darkMode ? "☀️" : "🌙"}
  </button>

  <button
    class="absolute top-4 right-4 text-3xl transition-transform duration-300"
    onclick={toggleMenu}
  >
    ☰
  </button>

  <div
    class="fixed top-0 right-0 h-full w-64 bg-white dark:bg-gray-800 shadow-xl z-50 transform transition-transform duration-300
           {menuOpen ? 'translate-x-0' : 'translate-x-full'}"
  >
    <div class="p-6 flex flex-col h-full">
      <div class="flex justify-end mb-4">
        <button class="text-3xl" onclick={toggleMenu}> ✕ </button>
      </div>

      <h2 class="text-2xl font-bold mb-4">Menu</h2>
      <nav>
        <ul class="space-y-4">
          <li><a href="/" class="text-lg hover:text-blue-500">Home</a></li>
          <li>
            <a href="/video" class="text-lg hover:text-blue-500"
              >Live feed test</a
            >
          </li>
          <li>
            <a href="/contact" class="text-lg hover:text-blue-500">Contact</a>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</div> -->

<script lang="ts">
  import Header from "$lib/Header.svelte";
  import Footer from "$lib/Footer.svelte";
  import AuthManager from "$lib/authManager";
  // import { onMount } from "svelte";
  onMount(() => {
    const auth = new AuthManager();
    console.log(auth.getUserInfo());
    let authed = null;
    cookieStore.get("authentication").then((res) => {
      authed = res;
      if (authed == null) {
        window.location.href = "/auth/login";
      }
    });
  });
  // onMount(() => {
  //   if (!isAuthenticated()) {
  //     window.location.href = "/auth/login";
  //   }
  // });
  import "../app.css";
  import { onMount } from "svelte";
</script>

<div class="flex flex-col min-h-screen relative">
  <!-- Background covers the whole viewport -->
  <div
    class="fixed inset-0 -z-10 h-full w-full bg-black"
    style="pointer-events: none;"
  >
    <div
      class="absolute bottom-0 left-0 right-0 top-0 bg-[linear-gradient(to_right,#4f4f4f2e_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:14px_24px]"
    ></div>
    <div
      class="absolute left-1/2 top-[-10%] -translate-x-1/2 h-[1000px] w-[1000px] rounded-full bg-[radial-gradient(circle_400px_at_50%_300px,#fbfbfb36,#000)]"
    ></div>
  </div>
  <Header />
  <main class="flex-grow">
    <slot />
  </main>
  <Footer />
</div>
