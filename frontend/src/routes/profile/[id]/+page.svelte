<script lang="ts">
    import { page } from "$app/state";
    import Debug from "$lib/components/Debug.svelte";
    import { auth } from "$lib/stores/auth.js";
    import {
        User,
        Twitter,
        Linkedin,
        MailIcon,
        X,
        ArrowRight,
    } from "lucide-svelte";
    import { onMount } from "svelte";

    let { data } = $props();

    let newProjectDescription = $state("");
    let newProjectLink = $state("");
    let newProjectName = $state("");

    let newSkillName = $state("");

    let newContactName = $state("");
    let newContactLink = $state("");
    let newContactIcon = $state("email");

    let editing = $state(false);
    let user = $state(data);
    let availableSkills = $state([]);
    const original = user;
    let id = page.params.id;

    $effect(() => {
        id;
        user = data;
    });

    const handleCancel = () => {
        user = original;

        editing = false;
    };

    const handleUpdate = () => {
        fetch("/api/updateUser", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(user),
        })
            .then((res) => {
                return res.json();
            })
            .then((res) => {
                console.log(res);
            });
    };

    const addProject = () => {
        user.projects.push({
            name: newProjectName,
            description: newProjectDescription,
            link: newProjectLink,
        });
    };

    const addSkill = () => {
        user.skills.push({
            name: newSkillName,
            verified: false,
        });
    };

    const addContact = () => {
        user.contacts.push({
            name: newContactName,
            link: newContactLink,
            icon: newContactIcon,
        });
    };

    onMount(async () => {
        fetch("/api/getSkills")
            .then((res) => {
                return res.json();
            })
            .then((res) => {
                availableSkills = res;
            });
    });


    function removeContact(id: any): any {
        user.contacts = user.contacts.filter((contact: { id: any; }) => contact.id !== id);
    }
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
                        src={user.profile_picture === "noPicture"
                            ? "/default-avatar.png"
                            : `/api/profile/${id}/picture`}
                        alt={`Profile picture of ${user.username}`}
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
                    {#if user.profession}
                        {#if editing}
                            <input
                                type="text"
                                bind:value={user.profession}
                                class="text-gray-600 text-lg"
                            />
                        {:else}
                            <p class="text-gray-600 text-lg">
                                {user.profession}
                            </p>
                        {/if}
                    {/if}
                    {#if editing}
                        <p class="text-sm text-gray-500 mt-2">
                            📍
                            <input
                                type="text"
                                bind:value={user.location}
                                class="text-sm text-gray-500 mt-2"
                            />
                        </p>
                    {:else}
                        <p class="text-sm text-gray-500 mt-2">
                            📍 {user.location}
                        </p>
                    {/if}
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
                        {#if editing}
                            <button
                                onclick={handleUpdate}
                                class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                            >
                                Save
                            </button>
                            <button
                                onclick={handleCancel}
                                class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                            >
                                Cancel
                            </button>
                        {:else}
                            <button
                                onclick={() => (editing = !editing)}
                                class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                            >
                                Edit Profile
                            </button>
                        {/if}
                    {/if}
                </div>
            </header>

            <main class="grid grid-cols-1 md:grid-cols-3 gap-8">
                <section class="md:col-span-2 space-y-6">
                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            About Me
                        </h2>
                        {#if editing}
                            <p class="text-gray-700 leading-relaxed">
                                <input
                                    type="text"
                                    bind:value={user.aboutme}
                                    class="text-gray-600 w-full h-full text-lg"
                                />
                            </p>
                        {:else}
                            <p class="text-gray-700 leading-relaxed">
                                {user.aboutme}
                            </p>
                        {/if}
                    </div>

                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            Projects
                        </h2>
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                            {#if user.projects != null && user.projects.length > 0}
                                {#if editing}
                                    <div
                                        class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition"
                                    >
                                        <h3
                                            class="font-bold text-lg text-gray-800"
                                        >
                                            New Project
                                        </h3>
                                        <p class="text-sm text-gray-600 mt-1">
                                            Add a new project to your profile.
                                        </p>
                                        <input
                                            type="text"
                                            bind:value={newProjectName}
                                            placeholder="Project Name"
                                            class="w-full border border-gray-300 rounded-md px-3 py-2 mt-2"
                                        />
                                        <input
                                            type="text"
                                            bind:value={newProjectDescription}
                                            placeholder="Project Description"
                                            class="w-full border border-gray-300 rounded-md px-3 py-2 mt-2"
                                        />
                                        <input
                                            type="text"
                                            bind:value={newProjectLink}
                                            placeholder="Project Link"
                                            class="w-full border border-gray-300 rounded-md px-3 py-2 mt-2"
                                        />
                                        <button
                                            onclick={addProject}
                                            class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                                        >
                                            Add Project
                                        </button>
                                    </div>
                                    {#each user.projects as { name, description, link }}
                                        <div
                                            class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm hover:shadow-md transition"
                                        >
                                            {#if editing}
                                                <button
                                                    class="bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center"
                                                >
                                                    <X />
                                                </button>
                                            {/if}
                                            <h3
                                                class="font-bold text-lg text-gray-800"
                                            >
                                                {name}
                                            </h3>
                                            <p
                                                class="text-sm text-gray-600 mt-1"
                                            >
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
                                    <p
                                        class="text-gray-600 w-full justify-center"
                                    >
                                        No projects found.
                                    </p>
                                {/if}
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
                            {#if user.skills && user.skills.length > 0}
                                {#if editing}
                                    <div class="relative">
                                        <select
                                            class="block appearance-none w-full bg-gray-200 border border-gray-200 text-gray-700 py-3 px-4 pr-8 leading-tight focus:outline-none focus:shadow-outline"
                                            bind:value={newSkillName}
                                        >
                                            <option value=""
                                                >Select a skill to add</option
                                            >
                                            {#each availableSkills as skill}
                                                <option value={skill.name}
                                                    >{skill.name}</option
                                                >
                                            {/each}
                                        </select>
                                        <div
                                            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700"
                                        >
                                            <ArrowRight />
                                        </div>
                                        <button
                                            onclick={addSkill}
                                            class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                                        >
                                            Add Skill
                                        </button>
                                    </div>
                                {/if}
                                {#each user.skills as skill}
                                    <span
                                        class="bg-blue-200 text-blue-800 text-sm font-medium px-3 py-1 rounded-full"
                                        >{skill.name}&nbsp;{skill.verified
                                            ? "✓"
                                            : ""}</span
                                    >
                                {/each}
                            {:else}
                                <p class="text-gray-600">No skills found.</p>
                            {/if}
                        </div>
                    </div>

                    <div class="bg-gray-50 rounded-lg p-6">
                        <h2 class="text-2xl font-bold text-gray-800 mb-4">
                            Contact
                        </h2>
                        <ul class="space-y-4 text-gray-600">
                            {#if editing}
                                <div class="flex items-center space-x-2">
                                    <input
                                        type="text"
                                        bind:value={newContactName}
                                        class="text-gray-600 w-full h-full text-lg"
                                    />
                                    <input
                                        type="text"
                                        bind:value={newContactLink}
                                        class="text-gray-600 w-full h-full text-lg"
                                    />
                                    <input
                                        type="text"
                                        bind:value={newContactIcon}
                                        class="text-gray-600 w-full h-full text-lg"
                                    />
                                    <button
                                        onclick={addContact}
                                        class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                                    >
                                        Add Contact
                                    </button>
                                </div>
                            {/if}
                            {#if user.contacts && user.contacts.length > 0}
                                {#each user.contacts as contact}
                                    <li class="flex items-center space-x-2">
                                        {#if editing}
                                            <button
                                                onclick={() =>
                                                    removeContact(contact.id)}
                                                class="px-6 py-2 rounded-full border border-gray-300 text-gray-700 font-semibold hover:bg-gray-100 transition"
                                            >
                                                Remove
                                            </button>
                                        {/if}
                                        {#if contact.icon === "email"}
                                            <MailIcon
                                                class="w-5 h-5 text-gray-500"
                                            />
                                            <span
                                                >{contact.name}: {contact.link}</span
                                            >
                                        {:else if contact.icon === "twitter"}
                                            <Twitter
                                                class="w-5 h-5 text-gray-500"
                                            />
                                            <a
                                                href={contact.link}
                                                class="text-blue-500 hover:underline"
                                            >
                                                {contact.name}
                                            </a>
                                        {:else if contact.icon === "linkedin"}
                                            <Linkedin
                                                class="w-5 h-5 text-gray-500"
                                            />
                                            <a
                                                href={contact.link}
                                                class="text-blue-500 hover:underline"
                                            >
                                                {contact.name}
                                            </a>
                                        {:else}
                                            <User
                                                class="w-5 h-5 text-gray-500"
                                            />
                                            <a
                                                href={contact.link}
                                                class="text-blue-500 hover:underline"
                                            >
                                                {contact.name}
                                            </a>
                                        {/if}
                                    </li>
                                {/each}
                            {:else}
                                <li class="text-gray-500">
                                    No contact information available
                                </li>
                            {/if}
                        </ul>
                    </div>
                </aside>
            </main>
        </div>
    {/if}
</div>
