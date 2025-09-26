<script>
    // @ts-nocheck

    import { auth } from "$lib/stores/auth";
    import { onMount } from "svelte";

    onMount(() => {
        document.location.replace(`/swapping`);
    });
</script>
