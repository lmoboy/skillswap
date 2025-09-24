<script lang="ts">
    import { page } from "$app/state";
    import Debug from "$lib/components/Debug.svelte";
    import { auth } from "$lib/stores/auth.js";
    import { User, Mail, Twitter, Linkedin, MailIcon } from "lucide-svelte";

    let { data } = $props();

    let user = data;
    let id = page.params.id;

    $effect(() => {
        id;
        user = data;
    });
</script>

<div class="bg-gray-100 min-h-screen p-8">
    <!-- <Debug {data} /> -->
    {#if !data}
        <div class="flex items-center justify-center">
            <div
                class="spinner-border animate-spin text-gray-500"
                role="status"
            >
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    {:else}
        <div
            class="max-w-6xl mx-auto bg-white rounded-xl shadow-lg p-8 space-y-8"
        >
            <header
                class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-8 pb-8 border-b border-gray-200"
            >
                <div class="relative w-32 h-32 flex-shrink-0">
                    <img
                        src={`/api/profile/${id}/picture`}
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
                    <h1 class="text-4xl font-bold text-gray-900">
                        {user.username}
                    </h1>
                    <p class="text-gray-600 text-lg">{user.profession}</p>
                    <p class="text-sm text-gray-500 mt-2">📍 {user.location}</p>
                </div>
                <div
                    class="flex-grow flex justify-center md:justify-end space-x-4"
                >
                    <a
                        href={`/swapping/${id}`}
                        class="px-6 py-2 rounded-full bg-blue-600 text-white font-semibold hover:bg-blue-700 transition"
                    >
                        Message
                    </a>
                    {#if user.id == $auth?.user?.id}
                        <button
                            class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                        >
                            Edit Profile
                        </button>
                    {/if}
                </div>
            </header>

            <main class="grid grid-cols-1 md:grid-cols-3 gap-8">
                <section class="md:col-span-2 space-y-6">
                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            About Me
                        </h2>
                        <p class="text-gray-700 leading-relaxed">
                            {user.aboutme}
                        </p>
                    </div>

                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            Projects
                        </h2>
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                            {#if user.projects != null && user.projects.length > 0}
                                {#each user.projects as { name, description, link }}
                                    <div
                                        class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition"
                                    >
                                        <h3
                                            class="font-bold text-lg text-gray-800"
                                        >
                                            {name}
                                        </h3>
                                        <p class="text-sm text-gray-600 mt-1">
                                            {description}
                                        </p>
                                        <a
                                            href={`${link}`}
                                            class="text-blue-500 text-sm mt-2 block hover:underline"
                                            >View Project →</a
                                        >
                                    </div>
                                {/each}
                            {:else}
                                <p class="text-gray-600 w-full justify-center">
                                    No projects found.
                                </p>
                            {/if}
                        </div>
                    </div>
                </section>

                <aside class="md:col-span-1 space-y-6">
                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            Skills
                        </h2>
                        <div class="flex flex-wrap gap-2">
                            {#each user.skills_coding as skill}
                                <span
                                    class="bg-blue-200 text-blue-800 text-sm font-medium px-3 py-1 rounded-full"
                                    >{skill.name}&nbsp;{skill.verified
                                        ? "✓"
                                        : ""}</span
                                >
                            {/each}
                        </div>
                    </div>

                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            Contact
                        </h2>
                        <ul class="space-y-4 text-gray-600">
                            <li class="flex items-center space-x-2">
                                <User class="w-5 h-5 text-gray-500" />
                                <span>{user.contacts.name ?? "Unknown"}</span>
                            </li>

                            <li class="flex items-center space-x-2">
                                <MailIcon class="w-5 h-5 text-gray-500" />
                                <span class="text-wrap"
                                    >{user.contacts.email ?? "Unknown"}</span
                                >
                            </li>

                            <li class="flex items-center space-x-2">
                                <Twitter class="w-5 h-5 text-gray-500" />
                                <span>
                                    <a
                                        href="https://twitter.com/{user.contacts
                                            .twitter
                                            ? user.contacts.twitter.substring(1)
                                            : ''}"
                                        class="text-blue-500 hover:underline {user
                                            .contacts.twitter
                                            ? ''
                                            : 'text-gray-400 disabled'}"
                                    >
                                        {user.contacts.twitter ?? "Unknown"}
                                    </a>
                                </span>
                            </li>

                            <li class="flex items-center space-x-2">
                                <Linkedin class="w-5 h-5 text-gray-500" />
                                <span>
                                    <a
                                        href="https://{user.contacts.linkedin
                                            ? user.contacts.linkedin
                                            : ''}"
                                        class="text-blue-500 hover:underline {user
                                            .contacts.linkedin
                                            ? ''
                                            : 'text-gray-400 disabled'}"
                                    >
                                        {user.contacts.linkedin ?? "Unknown"}
                                    </a>
                                </span>
                            </li>
                        </ul>
                    </div>
                </aside>
            </main>
        </div>
    {/if}
</div>
