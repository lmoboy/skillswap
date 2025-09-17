<script lang="ts">
    import { page } from "$app/stores";
    import Debug from "$lib/components/Debug.svelte";
    import { auth } from "$lib/stores/auth";
    import { onDestroy, onMount } from "svelte";
    let data = $state(null);

    let authState: {
        user: { name: string; email: string; id: string } | null;
        isAuthenticated: boolean;
        loading: boolean;
        error: string | null;
        step: string | null;
    } = {
        user: null,
        isAuthenticated: false,
        loading: true,
        error: null,
        step: null,
    };

    const unsubscribe = auth.subscribe((state) => {
        authState = state;
        console.log(state);
    });

    onDestroy(() => {
        unsubscribe();
    });

    onMount(async () => {
        await fetch(
            "https://localhost:8080/api/user?q=" + authState?.user?.id,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            },
        ).then((response) => {
            if (response.ok) {
                return response.json();
            }
        });
    });
</script>

<div class="bg-gray-100 min-h-screen p-8">
    <Debug {data} />

    <div class="max-w-6xl mx-auto bg-white rounded-xl shadow-lg p-8 space-y-8">
        <header
            class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-8 pb-8 border-b border-gray-200"
        >
            <div class="relative w-32 h-32 flex-shrink-0">
                <img
                    src="https://via.placeholder.com/150"
                    alt="Profile Picture"
                    class="w-full h-full rounded-full object-cover border-4 border-white shadow-md"
                />
                <span
                    class="absolute bottom-0 right-0 w-8 h-8 bg-green-500 rounded-full border-2 border-white transform translate-x-1 translate-y-1 flex items-center justify-center text-sm text-white font-bold"
                >
                    ✓
                </span>
            </div>
            <div class="text-center md:text-left">
                <h1 class="text-4xl font-bold text-gray-900">John Doe</h1>
                <p class="text-gray-600 text-lg">Web Developer & Designer</p>
                <p class="text-sm text-gray-500 mt-2">📍 New York, NY</p>
            </div>
            <div class="flex-grow flex justify-center md:justify-end space-x-4">
                <a
                    href={`/swapping/${authState?.user?.id}`}
                    class="px-6 py-2 rounded-full bg-blue-600 text-white font-semibold hover:bg-blue-700 transition"
                >
                    Message
                </a>
                <button
                    class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                >
                    Edit Profile
                </button>
            </div>
        </header>

        <main class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <section class="md:col-span-2 space-y-6">
                <div class="bg-gray-50 rounded-lg p-6">
                    <h2 class="text-2xl font-bold text-gray-800 mb-4">
                        About Me
                    </h2>
                    <p class="text-gray-700 leading-relaxed">
                        A passionate web developer with over 5 years of
                        experience in building and designing modern web
                        applications. I specialize in front-end technologies
                        like React, Vue, and Svelte, with a strong foundation in
                        back-end development using Node.js and Python. I enjoy
                        solving complex problems and creating user-friendly
                        interfaces.
                    </p>
                </div>

                <div class="bg-gray-50 rounded-lg p-6">
                    <h2 class="text-2xl font-bold text-gray-800 mb-4">
                        Projects
                    </h2>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                        <div
                            class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition"
                        >
                            <h3 class="font-bold text-lg text-gray-800">
                                SkillSwap Platform
                            </h3>
                            <p class="text-sm text-gray-600 mt-1">
                                A platform for real-time knowledge exchange.
                                Technologies: Svelte, Tailwind CSS, MySQL.
                            </p>
                            <a
                                href="#"
                                class="text-blue-500 text-sm mt-2 block hover:underline"
                                >View Project →</a
                            >
                        </div>
                        <div
                            class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition"
                        >
                            <h3 class="font-bold text-lg text-gray-800">
                                E-commerce API
                            </h3>
                            <p class="text-sm text-gray-600 mt-1">
                                A robust RESTful API for a digital marketplace.
                                Technologies: Node.js, Express, MongoDB.
                            </p>
                            <a
                                href="#"
                                class="text-blue-500 text-sm mt-2 block hover:underline"
                                >View Project →</a
                            >
                        </div>
                    </div>
                </div>
            </section>

            <aside class="md:col-span-1 space-y-6">
                <div class="bg-gray-50 rounded-lg p-6">
                    <h2 class="text-2xl font-bold text-gray-800 mb-4">
                        Skills
                    </h2>
                    <div class="flex flex-wrap gap-2">
                        <span
                            class="bg-blue-200 text-blue-800 text-sm font-medium px-3 py-1 rounded-full"
                            >Svelte</span
                        >
                        <span
                            class="bg-green-200 text-green-800 text-sm font-medium px-3 py-1 rounded-full"
                            >Tailwind CSS</span
                        >
                        <span
                            class="bg-yellow-200 text-yellow-800 text-sm font-medium px-3 py-1 rounded-full"
                            >JavaScript</span
                        >
                        <span
                            class="bg-purple-200 text-purple-800 text-sm font-medium px-3 py-1 rounded-full"
                            >Node.js</span
                        >
                        <span
                            class="bg-red-200 text-red-800 text-sm font-medium px-3 py-1 rounded-full"
                            >SQL</span
                        >
                    </div>
                </div>

                <div class="bg-gray-50 rounded-lg p-6">
                    <h2 class="text-2xl font-bold text-gray-800 mb-4">
                        Contact
                    </h2>
                    <ul class="space-y-3 text-gray-700">
                        <li class="flex items-center space-x-2">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                                ><path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8m-2 2v6a2 2 0 01-2 2H7a2 2 0 01-2-2v-6"
                                ></path></svg
                            >
                            <span>john.doe@example.com</span>
                        </li>
                        <li class="flex items-center space-x-2">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="currentColor"
                                viewBox="0 0 24 24"
                                aria-hidden="true"
                                xmlns="http://www.w3.org/2000/svg"
                                ><path
                                    d="M12 0c-6.626 0-12 5.373-12 12 0 5.372 3.513 9.882 8.36 11.458.61.112.828-.266.828-.593 0-.292-.01-1.066-.016-2.09-.344.174-.951.378-1.15.52-.395.14-.852-.2-.654-.424.184-.207.96-.97.96-1.574 0-.426-.347-.723-.717-.988-2.651.135-4.275-1.21-4.275-3.328 0-1.206.883-2.188 2.5-2.188 1.158 0 1.968.604 2.183 1.25.115.344.385.642.75.922 1.054 1.115 1.763.978 2.093.748.107-.872.417-1.463 1.15-1.93-.85-.094-1.722-.338-1.722-1.854 0-.81.41-1.472 1.11-1.99.106-.12.48-.56.12-.56-.25 0-.46.3-.59.6-1.26.195-2.02.825-2.02 1.636 0 1.25.996 1.65 1.286 1.77.107.042.27.062.38.07.135.01.272.016.39.016 1.05 0 2.22-.38 2.22-2.188 0-2.887-2.67-3.95-4.42-3.95-4.1 0-7.46 3.36-7.46 7.46 0 4.1 3.36 7.46 7.46 7.46 4.1 0 7.46-3.36 7.46-7.46 0-4.1-3.36-7.46-7.46-7.46z"
                                ></path></svg
                            >
                            <span
                                ><a href="#" class="hover:underline"
                                    >github.com/johndoe</a
                                ></span
                            >
                        </li>
                    </ul>
                </div>
            </aside>
        </main>
    </div>
</div>
