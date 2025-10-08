<script>
    // @ts-nocheck

    import { onMount } from "svelte";
    import { Mic, PhoneOff, Share2, VideoOff } from "lucide-svelte";
    import VideoDebug from "$lib/components/VideoDebug.svelte";
    import { auth } from "$lib/stores/auth";

    // Component-level variables bound to DOM elements
    let localVideo; // Will be bound to the local video element
    let remoteVideo; // Will be bound to the remote video element

    // WebRTC and state variables
    let localStream;
    let peerConnection;
    let ws;
    let roomId = "test-room"; // Default room for testing
    let isInitiator = false;

    // UI State
    let callStarted = false;
    let isMuted = false;
    let isVideoOff = false; // Used for UI toggle and track disabling
    let isSharingScreen = false; // Placeholder for future screen sharing logic

    // --- WebRTC Functions ---

    /**
     * Initializes the media stream, sets up the peer connection,
     * connects to the signaling server, and starts the session.
     */
    async function joinSession() {
        if (callStarted) return;
        if (!$auth.user) joinSession();
        const name = $auth.user.name || "Anonymous";

        try {
            // 1. Get local media (Video and Audio)
            localStream = await navigator.mediaDevices.getUserMedia({
                video: true,
                audio: true,
            });

            // Set the local video source
            if (localVideo) {
                localVideo.srcObject = localStream;
            }

            // 2. Setup RTCPeerConnection
            peerConnection = new RTCPeerConnection({
                iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
            });

            // Add local tracks to the peer connection
            localStream
                .getTracks()
                .forEach((track) =>
                    peerConnection.addTrack(track, localStream),
                );

            // 3. Setup WebSocket Signaling
            ws = new WebSocket(`/api/video?room=${roomId}`);

            ws.onopen = () => {
                console.log("Connected to the signaling server.");

                // If we're the initiator (first in room), create offer
                // For demo purposes, we'll assume the first user is the initiator
                // In a real app, you'd coordinate this via the server
                setTimeout(() => {
                    if (isInitiator) {
                        console.log("Creating OFFER as initiator");
                        createOffer();
                    }
                }, 1000);
            };

            // 4. Handle ICE Candidates
            peerConnection.onicecandidate = (event) => {
                if (event.candidate && ws.readyState === WebSocket.OPEN) {
                    console.log("Sending ICE Candidate");
                    ws.send(
                        JSON.stringify({
                            type: "candidate",
                            data: event.candidate,
                        }),
                    );
                }
            };

            // 5. Handle Remote Tracks (Receiving Stream)
            peerConnection.ontrack = (event) => {
                console.log("Remote track received.");

                if (remoteVideo && event.streams && event.streams[0]) {
                    console.log("Remote stream received.");
                    remoteVideo.srcObject = event.streams[0];
                }
            };

            // 6. Handle Signaling Messages
            ws.onmessage = async (message) => {
                const data = JSON.parse(message.data);
                if (!data || !data.type) return;

                console.log("Received signaling message:", data.type, data);

                switch (data.type) {
                    case "offer":
                        console.log("Received OFFER");
                        isInitiator = false;
                        await peerConnection.setRemoteDescription(
                            new RTCSessionDescription(data.data),
                        );
                        // Create and send ANSWER
                        const answer = await peerConnection.createAnswer();
                        await peerConnection.setLocalDescription(answer);
                        console.log("Sending ANSWER");
                        ws.send(
                            JSON.stringify({
                                type: "answer",
                                data: answer,
                            }),
                        );
                        break;
                    case "answer":
                        console.log("Received ANSWER");
                        await peerConnection.setRemoteDescription(
                            new RTCSessionDescription(data.data),
                        );
                        break;
                    case "candidate":
                        console.log("Received ICE Candidate");
                        // Ensure remote description is set before adding candidate
                        if (peerConnection.remoteDescription) {
                            await peerConnection.addIceCandidate(
                                new RTCIceCandidate(data.data),
                            );
                        } else {
                            // Candidate received before offer/answer, typically stored and added later
                            console.warn(
                                "ICE Candidate received before Remote Description was set.",
                            );
                        }
                        break;
                    default:
                        console.log("Unknown message type:", data.type);
                        break;
                }
            };

            callStarted = true;
        } catch (error) {
            console.error("Error joining session:", error);
            alert(`Failed to start video session: ${error.message}`);
            // Clean up if setup fails
            endCall();
        }
    }

    /**
     * Creates and sends an SDP offer to initiate the WebRTC connection.
     */
    async function createOffer() {
        try {
            const offer = await peerConnection.createOffer();
            await peerConnection.setLocalDescription(offer);
            console.log("Sending OFFER");
            ws.send(
                JSON.stringify({
                    type: "offer",
                    data: offer,
                }),
            );
        } catch (error) {
            console.error("Error creating offer:", error);
        }
    }

    /**
     * Toggles the enabled state of the local audio track.
     */
    function toggleMute() {
        if (!localStream) return;

        const audioTrack = localStream.getAudioTracks()[0];
        if (audioTrack) {
            audioTrack.enabled = !audioTrack.enabled;
            isMuted = !isMuted;
        }
    }

    /**
     * Toggles the enabled state of the local video track.
     * Note: Changed variable from isVideoStopped to isVideoOff for consistency.
     */
    function toggleVideo() {
        if (!localStream) return;

        const videoTrack = localStream.getVideoTracks()[0];
        if (videoTrack) {
            videoTrack.enabled = !videoTrack.enabled;
            isVideoOff = !isVideoOff;
        }
    }

    /**
     * Stops the media streams and closes the connections.
     */
    function endCall() {
        if (localStream) {
            localStream.getTracks().forEach((track) => track.stop());
            localStream = null;
        }
        if (peerConnection) {
            peerConnection.close();
            peerConnection = null;
        }
        if (ws) {
            ws.close();
            ws = null;
        }

        // Clear video elements
        if (localVideo) localVideo.srcObject = null;
        if (remoteVideo) remoteVideo.srcObject = null;

        callStarted = false;
        console.log("Call ended and connections closed.");
    }

    // Placeholder function for screen sharing
    async function toggleScreenShare() {
        if (isSharingScreen) {
            // Stop sharing (re-add camera track)
            console.log("Stopping screen share...");
            // TODO: Logic to replace screen track with camera track
            isSharingScreen = false;
        } else {
            // Start sharing
            try {
                const screenStream =
                    await navigator.mediaDevices.getDisplayMedia({
                        video: true,
                        audio: true,
                    });
                // TODO: Logic to replace camera track with screen track
                console.log("Starting screen share...");
                isSharingScreen = true;
            } catch (err) {
                console.error("Error starting screen share:", err);
            }
        }
    }

    // Automatically join the session when the component mounts
    onMount(() => {
        // Only attempt to join if the user is authenticated (auth store check)
        if (
            auth
                .waitForUser()
                .then(() => true)
                .catch(() => false)
        ) {
            joinSession();
        } else {
            // Handle case where user is not logged in or auth is not ready
            console.warn("User not authenticated. Cannot join session.");
        }

        // Clean up on component destruction
        return () => {
            endCall();
        };
    });
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
                    {#if !localStream || isVideoOff}
                        <div
                            class="absolute inset-0 flex items-center justify-center bg-gray-900 bg-opacity-90"
                        >
                            <span class="text-white text-sm">Camera Off</span>
                        </div>
                    {/if}
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

                    {#if !remoteVideo?.srcObject}
                        <div
                            class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 text-white"
                        >
                            <div class="text-center">
                                <div class="text-xl font-semibold mb-2">
                                    {#if callStarted}
                                        Waiting for partner to connect...
                                    {:else}
                                        Connecting...
                                    {/if}
                                </div>
                            </div>
                        </div>
                    {/if}
                </div>

                <div
                    class="bg-gray-800 p-4 flex items-center justify-center space-x-6"
                >
                    <button
                        on:click={toggleMute}
                        disabled={!callStarted}
                        class={`p-3 rounded-full ${isMuted ? "bg-red-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors disabled:opacity-50`}
                        aria-label={isMuted ? "Unmute" : "Mute"}
                    >
                        <Mic size={20} class={isMuted ? "line-through" : ""} />
                    </button>

                    <button
                        on:click={toggleVideo}
                        disabled={!callStarted}
                        class={`p-3 rounded-full ${isVideoOff ? "bg-red-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors disabled:opacity-50`}
                        aria-label={isVideoOff
                            ? "Turn on camera"
                            : "Turn off camera"}
                    >
                        <VideoOff
                            size={20}
                            class={isVideoOff ? "" : "stroke-white"}
                        />
                    </button>

                    <button
                        on:click={toggleScreenShare}
                        disabled={!callStarted}
                        class={`p-3 rounded-full ${isSharingScreen ? "bg-blue-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors disabled:opacity-50`}
                        aria-label={isSharingScreen
                            ? "Stop sharing screen"
                            : "Share screen"}
                    >
                        <Share2 size={20} />
                    </button>

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
        <details class="group">
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
