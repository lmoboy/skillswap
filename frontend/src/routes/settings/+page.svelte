<script>
    import { quintOut } from "svelte/easing";
    import { slide } from "svelte/transition";
    let preview = ""; // This should be a URL or data URL
    let name = "John Doe";
    let email = "john.doe@example.com";
    let preferences = "Daily email, weekly reports...";
    let password = "";
    let is2FAEnabled = false;

    /**
     * @param {{ target: { files: any[]; }; }} event
     */
    function handlePictureChange(event) {
        const file = event.target.files[0];
        if (file) {
            preview = URL.createObjectURL(file);
        }
    }

    function saveProfile() {
        // Handle form submission logic here
        console.log("Profile saved!");
    }
</script>

<div
    class="w-full min-h-screen flex items-center justify-center bg-white py-12"
>
    <section
        class="container mx-auto px-6 md:px-12 w-full max-w-4xl bg-white rounded-2xl shadow-xl p-6 sm:p-10 flex flex-col lg:flex-row gap-8"
    >
        <div
            class="flex-shrink-0 flex flex-col items-center gap-6 p-4 sm:p-8 border-b lg:border-r lg:border-b-0 border-gray-200 lg:w-1/3"
        >
            <img
                src={preview ||
                    "https://randomuser.me/api/portraits/men/33.jpg"}
                alt="Profile Preview"
                class="w-36 h-36 rounded-full border-4 border-gray-100 shadow-md transition-transform duration-300 hover:scale-105"
            />
            <div class="text-center">
                <h2 class="text-3xl font-bold text-gray-900 leading-tight">
                    {name}
                </h2>
                <p class="text-gray-600 mt-1">{email}</p>
            </div>
            <div class="w-full mt-4">
                <h3 class="font-semibold text-gray-700 mb-2">My Preferences</h3>
                <p class="text-gray-600 text-sm italic">
                    {preferences || "No preferences set yet."}
                </p>
            </div>
        </div>

        <div class="flex-grow p-4 sm:p-8">
            <h1 class="text-3xl font-bold mb-8 text-gray-900">
                Account Details
            </h1>
            <form onsubmit={saveProfile} class="space-y-6">
                <div>
                    <label class="block font-medium mb-2 text-gray-700"
                        >Change Profile Picture</label
                    >
                    <input
                        type="file"
                        accept="image/*"
                        onchange={handlePictureChange}
                        class="mt-1 block w-full text-sm text-gray-600 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-gray-200 file:text-gray-800 hover:file:bg-gray-300 transition"
                    />
                </div>

                <div class="grid sm:grid-cols-2 gap-6">
                    <div>
                        <label class="block font-medium mb-2 text-gray-700"
                            >Name</label
                        >
                        <input
                            type="text"
                            bind:value={name}
                            placeholder="Your name"
                            required
                            class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                        />
                    </div>
                    <div>
                        <label class="block font-medium mb-2 text-gray-700"
                            >Email Address</label
                        >
                        <input
                            type="email"
                            bind:value={email}
                            placeholder="Your email"
                            required
                            class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                        />
                    </div>
                </div>

                <div>
                    <label class="block font-medium mb-2 text-gray-700"
                        >New Password</label
                    >
                    <input
                        type="password"
                        bind:value={password}
                        placeholder="Leave blank to keep current password"
                        class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                    />
                </div>

                <div>
                    <label class="block font-medium mb-2 text-gray-700"
                        >Preferences</label
                    >
                    <textarea
                        bind:value={preferences}
                        placeholder="e.g., Daily email, weekly reports..."
                        class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                    ></textarea>
                </div>

                <div class="border-t border-gray-200 pt-6">
                    <h3 class="text-xl font-bold mb-4 text-gray-900">
                        Two-Factor Authentication (2FA)
                    </h3>
                    <div class="flex items-center justify-between">
                        <span class="text-gray-700 font-medium">Enable 2FA</span
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
                                class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-peach-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-peach-600"
                            ></div>
                        </label>
                    </div>
                    {#if is2FAEnabled}
                        <div
                            transition:slide={{ easing: quintOut }}
                            class="mt-4 p-4 bg-gray-100 rounded-lg text-sm text-gray-700 border border-gray-200"
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
                                    class="w-32 h-32 bg-white flex items-center justify-center rounded-md border border-gray-200"
                                >
                                    <span class="text-xs text-gray-500"
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
                        class="px-8 py-3 bg-gray-900 text-white rounded-lg font-medium hover:bg-gray-800 transition transform hover:scale-[1.02] shadow-lg"
                    >
                        Save Changes
                    </button>
                </div>
            </form>
        </div>
    </section>
</div>
