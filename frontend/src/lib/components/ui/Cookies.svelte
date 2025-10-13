<script>
    import { Cookie } from "lucide-svelte";
    import { onMount } from "svelte";
    let show = $state(false);
    onMount(async () => {
        let agreedCookie = await cookieStore
            .get("cookieConsent")
            .then((res) => {
                return res?.value;
            });
        if (agreedCookie === "true") {
            show = false;
        } else {
            show = true;
        }
    });
    const handleCookie = () => {
        console.log("Cookie accepted");
        cookieStore.set("cookieConsent", "true");
        show = false;
    };
</script>

<div
    class={`fixed border-t-2 ${show ? "" : "hidden"} border-gray-200 flex flex-col bottom-0 bg-white p-4 rounded-lg shadow-lg text-sm min-w-screen`}
>
    <Cookie class="w-12 h-12 text-gray-500" />
    <p class="mb-2 text-sm text-gray-500">
        We use cookies to ensure you get the best experience on our website.
    </p>
    <p class="mb-2 text-sm text-gray-500">
        <a href="/privacy" class="text-blue-500 hover:underline">
            Learn more about our privacy policy
        </a>
    </p>
    <button
        class="px-4 py-2 bg-blue-500 text-white rounded-full w-fit hover:bg-blue-700 cursor-pointer transition"
        onclick={handleCookie}
    >
        Accept
    </button>
    <script>
        function acceptCookies() {
            const cookieConsent = document.querySelector(".cookie-consent");
            cookieConsent.remove();
        }
    </script>
</div>
