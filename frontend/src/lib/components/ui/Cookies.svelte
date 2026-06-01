<script>
    import { onMount } from "svelte";
    let show = $state(false);

    onMount(() => {
        const consent = localStorage.getItem("cookieConsent");
        if (consent === "true") {
            show = false;
        } else {
            show = true;
        }
    });

    const handleCookie = () => {
        localStorage.setItem("cookieConsent", "true");
        show = false;
    };
</script>

{#if show}
    <div
        class="fixed bottom-0 left-0 right-0 bg-white border-t-2 border-gray-200 p-4 shadow-lg text-sm z-50"
    >
        <div class="max-w-6xl mx-auto flex flex-col sm:flex-row items-center gap-4">
            <p class="text-gray-600">
                We use cookies to ensure you get the best experience on our website.
                <a href="/privacy" class="text-blue-500 hover:underline ml-1">
                    Privacy policy
                </a>
            </p>
            <button
                class="px-4 py-2 bg-blue-500 text-white rounded-full hover:bg-blue-700 cursor-pointer transition flex-shrink-0"
                onclick={handleCookie}
            >
                Accept
            </button>
        </div>
    </div>
{/if}
