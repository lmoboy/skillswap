<script>
    import { goto } from "$app/navigation";
    import { ArrowLeft, ArrowRight, CircleDashed, Upload } from "lucide-svelte";
    let curSlide = 0;
    let fetching = false;
    let state = "intro";
    let username = "";
    let email = "";
    let password = "";
    let passwordr = "";
    let error = "";
    let success = "";

    const nextSlide = () => {
        state = "outro";
        setTimeout(() => {
            state = "intro";
            fetching = false;
            curSlide = (curSlide + 1) % 4;
        }, 500);
    };
    const prevSlide = () => {
        state = "outro";
        setTimeout(() => {
            state = "intro";
            fetching = false;
            curSlide = (curSlide - 1 + 4) % 4;
        }, 500);
    };
    const checkEmail = () => {
        fetching = true;
        fetch("http://localhost:8080/api/isEmailUsed", {
            method: "POST",
            body: JSON.stringify({ email: email }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    error = data.error;
                    nextSlide();
                } else if (data.status == "ok") {
                    error = "";
                    nextSlide();
                }
            });
    };
    const checkPassword = () => {
        if (password.includes(" ")) {
        }
        if (password.length >= 8 || password.length <= 32) {
            error = "";
            nextSlide();
        }
    };
    const checkUsername = () => {
        fetching = true;
        fetch("http://localhost:8080/api/isNameUsed", {
            method: "POST",
            body: JSON.stringify({ name: name }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    error = "name is already used!";
                    nextSlide();
                } else if (data.status == "ok") {
                    error = "";
                    nextSlide();
                }
            });
    };
    const handleSubmit = async () => {
        error = "";
        success = "";

        if (!email || !password || !passwordr || !username) {
            error = "All fields are required!";
            return;
        }

        if (password !== passwordr) {
            error = "Passwords do not match!";
            return;
        }

        if (username.length > 50) {
            error = "Username cannot be longer than 50 characters!";
            return;
        }

        if (email.length > 100) {
            error = "Email is too long!";
            return;
        }

        if (password.length > 50) {
            error = "Password cannot be longer than 50 characters!";
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/api/register", {
                method: "POST",
                body: JSON.stringify({ username, email, password }),
            });

            const data = await response.json();

            if (!response.ok) {
                error = data.error || "Registration failed. Please try again.";
                return;
            }

            success = "Registration successful! Redirecting to login...";
            setTimeout(() => {
                goto("/auth/login");
            }, 2000);
        } catch (err) {
            console.error("Network error during registration:", err);
            error = "Network error. Please check your connection.";
        }
    };
</script>

<div class="flex flex-col items-center justify-center h-screen">
    <div
        class="bg-[#202020] p-8 rounded-lg flex-col shadow-lg h-2/6 w-lg flex justify-center align-middle items-center"
    >
        <div class="flex justify-center align-middle items-center flex-col">
            {#if curSlide === 0}
                <div
                    class={` ${state} font-bold text-2xl text-white mb-4 flex justify-center align-middle items-center text-center`}
                >
                    Ready to join a community of like-minded individuals? Sign
                    up now and unlock a world of opportunities!
                </div>
                <button
                    disabled={state === "outro"}
                    on:click={nextSlide}
                    class={`${state} bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-full`}
                >
                    <ArrowRight class="w-6 h-6 p-0 m-0" />
                </button>
            {:else if curSlide === 1}
                <div
                    class={`font-bold text-2xl text-white mb-4 flex justify-center align-middle items-center text-center ${state}`}
                >
                    Choose a unique username.
                </div>
                <div class={`${state} flex flex-row items-center gap-2`}>
                    <input
                        type="text"
                        placeholder="Username"
                        bind:value={username}
                        class={`bg-gray-700 px-4 py-2 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500`}
                    />
                    <button
                        on:click={nextSlide}
                        class={`${state} bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-full`}
                    >
                        <ArrowRight class="w-6 h-6 p-0 m-0" />
                    </button>
                </div>
            {:else if curSlide === 2}
                <div
                    class={`${state} font-bold text-2xl text-white mb-4 flex justify-center align-middle items-center text-center`}
                >
                    Enter your email address.
                </div>
                <div class={`${state} flex flex-row items-center gap-2`}>
                    <input
                        type="email"
                        placeholder="Email"
                        bind:value={email}
                        class={`${state} bg-gray-700 px-4 py-2 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500`}
                    />
                    <button
                        on:click={nextSlide}
                        class={`${state} bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-full`}
                    >
                        <ArrowRight class="w-6 h-6 p-0 m-0" />
                    </button>
                </div>
            {:else if curSlide === 3}
                <form
                    on:submit|preventDefault={handleSubmit}
                    class={`${state} flex flex-col items-center gap-4`}
                >
                    <div class="font-bold text-2xl text-white mb-2 text-center">
                        Create a strong password.
                    </div>
                    <input
                        type="password"
                        placeholder="Password (min 8 chars)"
                        bind:value={password}
                        class="bg-gray-700 px-4 py-2 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <input
                        type="password"
                        placeholder="Confirm Password"
                        bind:value={passwordr}
                        class="bg-gray-700 px-4 py-2 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <button
                        type="submit"
                        disabled={fetching}
                        class={`${fetching ? "animate-spin" : ""} bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-full`}
                    >
                        {#if fetching}
                            <svg
                                class="animate-spin h-6 w-6 text-white"
                                viewBox="0 0 24 24"
                            >
                                <circle
                                    class="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    stroke-width="4"
                                ></circle>
                                <path
                                    class="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                ></path>
                            </svg>
                        {:else}
                            <Upload class="w-6 h-6 p-0 m-0" />
                        {/if}
                    </button>
                </form>
            {/if}

            {#if error}
                <div class="text-red-500 mt-4 text-center">{error}</div>
            {:else if success}
                <div class="text-green-500 mt-4 text-center">{success}</div>
            {/if}
        </div>
    </div>
</div>

<style>
    @keyframes fade-in-from-side {
        from {
            transform: translateX(20px);
            opacity: 0;
        }
        to {
            transform: translateX(0px);
            opacity: 1;
        }
    }

    @keyframes fade-out-to-side {
        from {
            transform: translateX(0px);
            opacity: 1;
        }
        to {
            transform: translateX(-20px);
            opacity: 0;
        }
    }
    .intro {
        animation: fade-in-from-side 0.5s ease-in-out;
    }
    .outro {
        animation: fade-out-to-side 0.6s ease-in-out;
    }
</style>
