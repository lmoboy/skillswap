<script lang="ts">
    import { onMount } from "svelte";
    import { Mic, PhoneOff, Share2, VideoOff } from "lucide-svelte";
    // import VideoDebug from "$lib/components/VideoDebug.svelte";
    // let remoteVideo: HTMLVideoElement;
    // let localVideo: HTMLVideoElement;
    // let localStream: MediaStream | null = null;
    let isMuted = false;
    let isVideoOff = false;
    let isSharingScreen = false;

    function toggleMute() {
        isMuted = !isMuted;
    }

    function toggleVideo() {
        isVideoOff = !isVideoOff;
    }

    async function toggleScreenShare() {
        isSharingScreen = !isSharingScreen;
    }

    function endCall() {
        console.log("End call implement later");
    }
</script>

<div class={`flex flex-col max-h-1/4`}>
    <main class="flex-grow flex items-center justify-center">
        <div class="w-full">
            <div class="relative bg-gray-900 overflow-hidden shadow-2xl">
                <div
                    class="absolute bottom-4 h-34 right-4 w-48 bg-black overflow-hidden border-2 border-white shadow-lg z-10"
                >
                    <!-- <video
                        bind:this={localVideo}
                        class="w-full h-full object-cover"
                        autoplay
                        playsinline
                        muted
                    ></video> -->
                </div>

                <div
                    class="aspect-video bg-black flex items-center justify-center"
                >
                    <!-- <video
                        bind:this={remoteVideo}
                        class="w-full h-full object-cover"
                        autoplay
                        playsinline
                        muted
                    ></video> -->
                    <div class="text-center text-white">
                        <div class="text-2xl font-semibold mb-2">
                            {isVideoOff ? "Your camera is off" : "Connected"}
                        </div>
                        <div class="text-gray-400">
                            {isMuted
                                ? "Microphone is muted"
                                : "Microphone is active"}
                        </div>
                    </div>
                </div>

                <div
                    class="bg-gray-800 p-4 flex items-center justify-center space-x-6"
                >
                    <button
                        onclick={toggleMute}
                        class={`p-3 rounded-full ${isMuted ? "bg-red-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
                        aria-label={isMuted ? "Unmute" : "Mute"}
                    >
                        <Mic size={20} class={isMuted ? "line-through" : ""} />
                    </button>

                    <button
                        onclick={toggleVideo}
                        class={`p-3 rounded-full ${isVideoOff ? "bg-red-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
                        aria-label={isVideoOff
                            ? "Turn on camera"
                            : "Turn off camera"}
                    >
                        <VideoOff
                            size={20}
                            class={isVideoOff ? "line-through" : ""}
                        />
                    </button>

                    <button
                        onclick={toggleScreenShare}
                        class={`p-3 rounded-full ${isSharingScreen ? "bg-blue-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
                        aria-label={isSharingScreen
                            ? "Stop sharing screen"
                            : "Share screen"}
                    >
                        <Share2 size={20} />
                    </button>

                    <button
                        onclick={endCall}
                        class="p-3 rounded-full bg-red-600 hover:bg-red-700 text-white transition-colors"
                        aria-label="End call"
                    >
                        <PhoneOff size={20} />
                    </button>
                </div>
            </div>
        </div>
    </main>
</div>
