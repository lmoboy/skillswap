<script lang="ts">
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
  class="w-screen min-h-screen flex flex-col items-center justify-center duration-300 relative overflow-hidden {darkMode
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
</div>
