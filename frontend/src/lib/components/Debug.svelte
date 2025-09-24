<script>
    let watchables = $props();
    $effect(() => {
        watchables;
        // console.log(watchables);
    });
</script>

<div>
    <div
        class="fixed bottom-4 right-4 bg-white p-4 rounded-lg shadow-lg text-sm w-fit break-words max-w-xs"
    >
        <h3 class="font-bold mb-2 text-gray-700">Debug</h3>
        <code class="space-y-2">
            {#each Object.entries(watchables) as [key, value]}
                <div class="flex justify-between text-gray-700">
                    <span>{key} :</span>
                    <span> {JSON.stringify(value)}</span>
                </div>
                <br />
            {/each}
        </code>
    </div>
</div>
