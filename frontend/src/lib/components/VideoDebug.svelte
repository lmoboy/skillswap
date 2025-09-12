<script lang="ts">
    import { onMount } from 'svelte';
    
    let videoState = {
        isConnected: false,
        loading: false,
        error: null as string | null,
        data: null as {
            bitrate: number;
            framerate: number;
            resolution: string;
        } | null
    };

    onMount(async () => {
        try {
            videoState.loading = true;
            const stream = await navigator.mediaDevices.getUserMedia({ 
                video: true, 
                audio: true 
            });
            videoState.isConnected = true;
            // You can add more stream processing here if needed
            stream.getTracks().forEach(track => {
                track.onended = () => {
                    videoState.isConnected = false;
                    videoState = videoState; // Trigger Svelte reactivity
                };
            });
        } catch (err) {
            videoState.error = err instanceof Error ? err.message : 'Failed to access media devices';
        } finally {
            videoState.loading = false;
        }
    });
</script>

<div
    class="fixed bottom-4 right-4 bg-white p-4 rounded-lg shadow-lg text-sm max-w-xs"
>
    <h3 class="font-bold mb-2">Video Debug</h3>
    <div class="space-y-2">
        <div class="flex justify-between">
            <span class="font-medium">Status:</span>
            <span class="font-mono">
                {videoState.isConnected ? "✅ Connected" : "❌ Not Connected"}
            </span>
        </div>

        {#if videoState.loading}
            <div class="text-blue-600">Loading...</div>
        {/if}

        {#if videoState.error}
            <div class="text-red-600">Error: {videoState.error}</div>
        {/if}

        {#if videoState.data}
            <div class="mt-2 pt-2 border-t">
                <div class="font-medium">User:</div>
                <div class="ml-2">
                    <div>Bitrate: {videoState.data.bitrate}</div>
                    <div>Framerate: {videoState.data.framerate}</div>
                    <div>Resolution: {videoState.data.resolution}</div>
                </div>
            </div>
        {/if}
    </div>
</div>
