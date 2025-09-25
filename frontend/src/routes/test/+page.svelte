<script>
// @ts-nocheck

    import { Mic, PhoneOff, Share2, VideoOff, Send } from "lucide-svelte";
    import VideoDebug from "$lib/components/VideoDebug.svelte";
    import { auth } from "$lib/stores/auth";

    let localStream;
    let peerConnection;
    let isMuted = false;
    let isVideoStopped = false;

    async function joinSession() {
        const name = $auth.user.name;

        
        peerConnection = new RTCPeerConnection({
            iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
        });

        localStream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true,
        });
        localStream
            .getTracks()
            .forEach((track) => peerConnection.addTrack(track, localStream));

        const localVideo = document.createElement("video");
        localVideo.srcObject = localStream;
        localVideo.autoplay = true;
        localVideo.muted = true;

        const ws = new WebSocket(`/api/video`);

        ws.onopen = () => {
            console.log("Connected to the signaling server");
            ws.send(JSON.stringify({ type: "join", name: name }));
        };

        ws.onmessage = async (message) => {
            const data = JSON.parse(message.data);
            switch (data.type) {
                case "offer":
                    await peerConnection.setRemoteDescription(
                        new RTCSessionDescription(data.offer),
                    );
                    const answer = await peerConnection.createAnswer();
                    await peerConnection.setLocalDescription(answer);
                    ws.send(JSON.stringify({ type: "answer", answer: answer }));
                    break;
                case "answer":
                    await peerConnection.setRemoteDescription(
                        new RTCSessionDescription(data.answer),
                    );
                    break;
                case "candidate":
                    await peerConnection.addIceCandidate(
                        new RTCIceCandidate(data.candidate),
                    );
                    break;
                default:
                    break;
            }
        };

        peerConnection.onicecandidate = (event) => {
            if (event.candidate) {
                ws.send(
                    JSON.stringify({
                        type: "candidate",
                        candidate: event.candidate,
                    }),
                );
            }
        };

        peerConnection.ontrack = (event) => {
            addRemoteStream(event.streams[0]);
        };
    }

    function toggleMute() {
        localStream
            .getAudioTracks()
            .forEach(
                (track) =>
                    (track.enabled = !track.enabled),
            );
        isMuted = !isMuted;
    }

    function toggleVideo() {
        localStream.getVideoTracks().forEach(
                (track) =>
                    (track.enabled = !track.enabled),
            );
        isVideoStopped = !isVideoStopped;
    }

    function addRemoteStream(stream) {
        const remoteVideo = document.createElement("video");
        remoteVideo.srcObject = stream;
        remoteVideo.autoplay = true;
    }
</script>

<div class="min-h-screen bg-gray-50 flex flex-col">
    <main class="flex-grow flex items-center justify-center p-4">
        <div class="w-full max-w-6xl">
            <div
                class="relative bg-gray-900 rounded-xl overflow-hidden shadow-2xl"
            >
                <div
                    class="absolute bottom-4 right-4 w-48 h-36 bg-black rounded-lg overflow-hidden border-2 border-white shadow-lg z-10"
                >
                    <video
                        bind:this={localVideo}
                        class="w-full h-full object-cover"
                        autoplay
                        playsinline
                        muted
                    ></video>
                </div>

                <div
                    class="aspect-video bg-black flex items-center justify-center"
                >
                    <video
                        bind:this={remoteVideo}
                        class="w-full h-full object-cover"
                        autoplay
                        playsinline
                        muted={false}
                    ></video>
                    <!-- {#if !remoteVideo?.srcObject}
                        <div
                            class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 text-white"
                        >
                            <div class="text-center">
                                <div class="text-xl font-semibold mb-2">
                                    Waiting for video stream...
                                </div>
                                <div class="text-sm">
                                    Status: {statusMessage}
                                </div>
                                {#if errorMessage}
                                    <div class="text-red-400 mt-2 text-sm">
                                        {errorMessage}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    {/if} -->
                </div>

                <div
                    class="bg-gray-800 p-4 flex items-center justify-center space-x-6"
                >
                    <button
                        on:click={toggleMute}
                        class={`p-3 rounded-full ${isMuted ? "bg-red-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
                        aria-label={isMuted ? "Unmute" : "Mute"}
                    >
                        <Mic size={20} class={isMuted ? "line-through" : ""} />
                    </button>

                    <button
                        on:click={toggleVideo}
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

                    <!-- <button
                        on:click={toggleScreenShare}
                        class={`p-3 rounded-full ${isSharingScreen ? "bg-blue-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
                        aria-label={isSharingScreen
                            ? "Stop sharing screen"
                            : "Share screen"}
                    >
                        <Share2 size={20} />
                    </button> -->

                    <button
                        on:click={endCall}
                        class="p-3 rounded-full bg-red-600 hover:bg-red-700 text-white transition-colors"
                        aria-label="End call"
                    >
                        <PhoneOff size={20} />
                    </button>
                </div>
            </div>
        </div>
    </main>

    <div class="p-4 border-t border-gray-200 bg-white">
        <details class="group" open>
            <summary
                class="text-sm font-medium text-gray-700 cursor-pointer hover:text-gray-900 flex items-center"
            >
                <svg
                    class="w-4 h-4 mr-2 transition-transform group-open:rotate-90"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 5l7 7-7 7"
                    />
                </svg>
                Debug Information
            </summary>
            <div class="mt-4 p-4 bg-gray-50 rounded-lg">
                <VideoDebug />
            </div>
        </details>
    </div>
</div>
