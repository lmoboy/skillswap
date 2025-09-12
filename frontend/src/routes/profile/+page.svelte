<script>
    // @ts-nocheck
    import "../../app.css";
    import { quintOut } from "svelte/easing";
    import { slide } from "svelte/transition";

    let data = $props();
    let name = $state(data.name);
    let email = $state(data.email);
    let preferences = $state("Notifications: daily digest, email newsletters.");
    let password = $state("");
    let preview = "https://randomuser.me/api/portraits/men/33.jpg"; // Default profile picture
    let is2FAEnabled = $state(false);

    // function handlePictureChange(event) {
    //     const file = event.target.files[0];
    //     if (file) {
    //         preview = URL.createObjectURL(file);
    //     }
    // }

    // function saveProfile() {
    //     if (password) {
    //         console.log("Password changed.");
    //     }
    // }
</script>

<div
    class="min-h-screen bg-black/90 p-6 sm:p-10 flex items-center justify-center"
>
    <section
        class="w-full max-w-4xl rounded-2xl shadow-2xl shadow-blue-900/30 p-6 sm:p-10 flex flex-col lg:flex-row gap-8 border border-blue-400/30 backdrop-blur-lg"
    >
        <div
            class="flex-shrink-0 flex flex-col items-center gap-6 p-4 sm:p-8 border-b lg:border-r lg:border-b-0 border-blue-400/20 lg:w-1/3"
        >
            <img
                src={preview}
                alt="Profile Preview"
                class="w-36 h-36 rounded-full border-4 border-blue-500 shadow-lg transition-transform duration-300 hover:scale-105 bg-white/20"
            />
            <div class="text-center">
                <h2 class="text-3xl font-extrabold text-gray-100 drop-shadow">
                    {name}
                </h2>
                <p class="text-gray-400 mt-1">{email}</p>
            </div>
            <div class="w-full mt-4">
                <h3 class="font-bold text-blue-300 mb-2">My Preferences</h3>
                <p class="text-gray-300 text-sm italic">
                    {preferences || "No preferences set yet."}
                </p>
            </div>
        </div>

        <div class="flex-grow p-4 sm:p-8">
            <h1 class="text-3xl font-bold mb-8 text-blue-100 drop-shadow">
                Account Details
            </h1>
            <form on:submit|preventDefault={saveProfile} class="space-y-6">
                <div>
                    <label class="block font-semibold mb-2 text-blue-200"
                        >Change Profile Picture</label
                    >
                    <input
                        type="file"
                        accept="image/*"
                        on:change={handlePictureChange}
                        class="mt-1 block w-full text-sm text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-100 file:text-blue-700 hover:file:bg-blue-200 transition"
                    />
                </div>

                <div class="grid sm:grid-cols-2 gap-6">
                    <div>
                        <label class="block font-semibold mb-2 text-blue-200"
                            >Name</label
                        >
                        <input
                            type="text"
                            bind:value={name}
                            placeholder="Your name"
                            required
                            class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-900/60 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                        />
                    </div>
                    <div>
                        <label class="block font-semibold mb-2 text-blue-200"
                            >Email Address</label
                        >
                        <input
                            type="email"
                            bind:value={email}
                            placeholder="Your email"
                            required
                            class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-900/60 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                        />
                    </div>
                </div>

                <div>
                    <label class="block font-semibold mb-2 text-blue-200"
                        >New Password</label
                    >
                    <input
                        type="password"
                        bind:value={password}
                        placeholder="Leave blank to keep current password"
                        class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-900/60 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                    />
                </div>

                <div>
                    <label class="block font-semibold mb-2 text-blue-200"
                        >Preferences</label
                    >
                    <textarea
                        bind:value={preferences}
                        placeholder="e.g., Daily email, weekly reports..."
                        class="w-full p-3 rounded-lg border border-blue-400/30 bg-gray-900/60 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                    ></textarea>
                </div>

                <div class="border-t border-blue-400/20 pt-6">
                    <h3 class="text-xl font-bold mb-4 text-blue-100">
                        Two-Factor Authentication (2FA)
                    </h3>
                    <div class="flex items-center justify-between">
                        <span class="text-blue-200 font-medium">Enable 2FA</span
                        >
                        <label
                            class="relative inline-flex items-center cursor-pointer"
                        >
                            <input
                                type="checkbox"
                                bind:checked={is2FAEnabled}
                                class="sr-only peer"
                            />
                            <div
                                class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
                            ></div>
                        </label>
                    </div>
                    {#if is2FAEnabled}
                        <div
                            transition:slide|local={{ easing: quintOut }}
                            class="mt-4 p-4 bg-gray-800/80 rounded-lg text-sm text-blue-100 border border-blue-400/20"
                        >
                            <p>
                                <strong>Setup Instructions:</strong>
                                <br />
                                1. Install an authenticator app (e.g., Google Authenticator,
                                Authy).
                                <br />
                                2. Scan the QR code below or enter the key manually.
                            </p>
                            <div class="mt-4 flex items-center justify-center">
                                <div
                                    class="w-32 h-32 bg-gray-700 flex items-center justify-center rounded-md border border-blue-400/30"
                                >
                                    <span class="text-xs text-blue-300"
                                        >QR Code Here</span
                                    >
                                </div>
                            </div>
                        </div>
                    {/if}
                </div>

                <div class="flex justify-end pt-6">
                    <button
                        type="submit"
                        class="px-8 py-3 bg-blue-600 text-white rounded-full font-bold hover:bg-blue-700 transition transform hover:scale-105 shadow-lg"
                    >
                        Save Changes
                    </button>
                </div>
            </form>
        </div>
    </section>
</div>
