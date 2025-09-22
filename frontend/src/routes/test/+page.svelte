<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { Mic, PhoneOff, Share2, VideoOff, Send } from "lucide-svelte";
    import VideoDebug from "$lib/components/VideoDebug.svelte";

    let localVideo: HTMLVideoElement;
    let remoteVideo: HTMLVideoElement;

    let localStream: MediaStream | null = null;
    let remoteStream: MediaStream | null = null;
    let peerConnection: RTCPeerConnection | null = null;
    let pendingCandidates: RTCIceCandidateInit[] = [];
    let socket: WebSocket | null = null;
    let isConnected = false;
    let isStreaming = false;

    let isMuted = false;
    let isVideoOff = false;
    let isSharingScreen = false;
    let isSubmitting = false;
    let errorMessage = "";
    let statusMessage = "Initializing...";

    let formData = {
        message: "",
        rating: 5,
        isPublic: true,
    };
    const rtcConfig: RTCConfiguration = {
        iceServers: [
            // { urls: "stun:stun.l.google.com:19302" },
            {
                urls: "turn:localhost:3499",
                username: "admin",
                credential: "admin",
            },
        ],
    };

    const initWebSocket = () => {
        try {
            const wsProtocol =
                window.location.protocol === "https:" ? "wss:" : "ws:";
            const wsUrl = `${wsProtocol}//${window.location.host}/api/video`;

            console.log("Connecting to WebSocket:", wsUrl);
            socket = new WebSocket(wsUrl);

            socket.onopen = () => {
                console.log("WebSocket connected successfully");
                isConnected = true;
                statusMessage = "Connected to server";
                errorMessage = "";
                startWebRTC();
            };

            socket.onclose = (event) => {
                console.log(
                    "WebSocket disconnected:",
                    event.code,
                    event.reason,
                );
                isConnected = false;
                statusMessage = "Disconnected from server";
                if (event.code !== 1000) {
                    errorMessage = `Connection closed: ${event.reason || "Unknown reason"}`;
                }
                stopWebRTC();
                // optionally try reconnect
                if (!event.wasClean) {
                    console.log("Attempting to reconnect in 3 seconds...");
                    setTimeout(initWebSocket, 3000);
                }
            };

            socket.onerror = (error) => {
                console.error("WebSocket error:", error);
                errorMessage = "Connection error. Trying to reconnect...";
            };

            socket.onmessage = async (event) => {
                try {
                    const message = JSON.parse(event.data);
                    console.log("Received WebSocket message:", message);

                    if (!peerConnection) {
                        console.error(
                            "Received message but no peer connection exists",
                        );
                        return;
                    }
                    switch (message.type) {
                        case "answer":
                            const answer = new RTCSessionDescription(
                                message.data,
                            );
                            await peerConnection.setRemoteDescription(answer);
                            console.log("Remote description set (answer)");
                            // Now drain pending
                            for (const cand of pendingCandidates) {
                                try {
                                    await peerConnection.addIceCandidate(
                                        new RTCIceCandidate(cand),
                                    );
                                } catch (e) {
                                    console.error(
                                        "Error adding pending ICE candidate",
                                        e,
                                    );
                                }
                            }
                            pendingCandidates = [];
                            break;
                        case "ice-candidate":
                            const candData = message.data;
                            const iceCandidateInit = new RTCIceCandidate(
                                candData,
                            );
                            if (
                                !peerConnection.remoteDescription ||
                                !peerConnection.remoteDescription.type
                            ) {
                                console.log(
                                    "RemoteDescription not set yet, queuing ICE candidate",
                                    iceCandidateInit,
                                );
                                pendingCandidates.push(iceCandidateInit);
                            } else {
                                try {
                                    await peerConnection.addIceCandidate(
                                        iceCandidateInit,
                                    );
                                } catch (e) {
                                    console.error(
                                        "Error adding ICE candidate:",
                                        e,
                                    );
                                }
                            }
                            break;
                    }
                } catch (err) {
                    console.error(
                        "Error processing WebSocket message:",
                        err,
                        event.data,
                    );
                }
            };
        } catch (err) {
            console.error("Failed to initialize WebSocket:", err);
            errorMessage = "Failed to connect to server";
            setTimeout(initWebSocket, 3000);
        }
    };

    const startWebRTC = async () => {
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.error("Cannot start WebRTC: WebSocket is not connected");
            errorMessage = "Not connected to server";
            return;
        }

        try {
            console.log("Starting WebRTC connection...");
            peerConnection = new RTCPeerConnection(rtcConfig);

            peerConnection.onicecandidate = (event) => {
                if (event.candidate) {
                    console.log("New ICE candidate:", event.candidate);
                    if (socket?.readyState === WebSocket.OPEN) {
                        socket.send(
                            JSON.stringify({
                                type: "ice-candidate",
                                data: event.candidate.toJSON(),
                            }),
                        );
                        console.log("Sent ICE candidate to server");
                    }
                } else {
                    console.log("ICE candidate gathering finished");
                }
            };

            peerConnection.ontrack = (event) => {
                console.log("ontrack", event);
                const streamFromEvent = event.streams && event.streams[0];
                if (streamFromEvent) {
                    remoteStream = streamFromEvent;
                    if (remoteVideo) {
                        remoteVideo.srcObject = streamFromEvent;
                        remoteVideo.play().catch((e) => console.error(e));
                    }
                } else {
                    // fallback to adding individual tracks:
                    if (!remoteStream) {
                        remoteStream = new MediaStream();
                        if (remoteVideo) {
                            remoteVideo.srcObject = remoteStream;
                            remoteVideo.play().catch((e) => console.error(e));
                        }
                    }
                    remoteStream.addTrack(event.track);
                }
                statusMessage = "Video stream received";
            };

            peerConnection.onconnectionstatechange = () => {
                if (!peerConnection) return;
                const state = peerConnection.connectionState;
                console.log("Connection state:", state);
                statusMessage = `Connection: ${state}`;
                if (state === "connected") {
                    isStreaming = true;
                    console.log("WebRTC connection established");
                } else if (
                    ["disconnected", "failed", "closed"].includes(state)
                ) {
                    console.warn("WebRTC connection lost or failed", state);
                    isStreaming = false;
                }
            };

            // add local stream tracks
            if (localStream) {
                console.log("Adding local stream tracks");
                localStream.getTracks().forEach((track) => {
                    if (peerConnection) {
                        peerConnection.addTrack(
                            track,
                            localStream as MediaStream,
                        );
                        console.log(`Added local ${track.kind} track`);
                    }
                });
            } else {
                console.warn("No local stream available when starting WebRTC");
            }

            // create offer
            console.log("Creating offer...");
            const offer = await peerConnection.createOffer({
                offerToReceiveAudio: true,
                offerToReceiveVideo: true,
            });
            await peerConnection.setLocalDescription(offer);

            // send offer via websocket
            if (socket?.readyState === WebSocket.OPEN) {
                socket.send(
                    JSON.stringify({
                        type: "offer",
                        data: peerConnection.localDescription,
                    }),
                );
                console.log("Sent offer to server");
                statusMessage = "Sent offer, waiting for answer...";
            } else {
                throw new Error("WebSocket closed before sending offer");
            }
        } catch (error: any) {
            console.error("Error in startWebRTC:", error);
            errorMessage = `Connection error: ${error.message}`;
            stopWebRTC();
        }
    };

    const stopWebRTC = () => {
        if (peerConnection) {
            peerConnection.getSenders().forEach((sender) => {
                // optionally stop tracks
                if (sender.track) {
                    sender.track.stop();
                }
            });
            peerConnection.ontrack = null;
            peerConnection.onicecandidate = null;
            peerConnection.onconnectionstatechange = null;
            peerConnection.close();
            peerConnection = null;
        }
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.close();
        }
        if (remoteStream) {
            remoteStream.getTracks().forEach((t) => t.stop());
            remoteStream = null;
        }
    };

    const toggleVideo = () => {
        if (localStream) {
            const videoTracks = localStream.getVideoTracks();
            if (videoTracks.length > 0) {
                isVideoOff = !isVideoOff;
                videoTracks[0].enabled = !isVideoOff;
            }
        }
    };

    const toggleAudio = () => {
        if (localStream) {
            const audioTracks = localStream.getAudioTracks();
            if (audioTracks.length > 0) {
                isMuted = !isMuted;
                audioTracks[0].enabled = !isMuted;
            }
        }
    };

    onMount(async (): Promise<any> => {
        let cleanupComplete = false;
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                video: {
                    width: { ideal: 1280 },
                    height: { ideal: 720 },
                    frameRate: { ideal: 30, max: 30 },
                },
                audio: true,
            });
            localStream = stream;
            if (localVideo) {
                localVideo.srcObject = stream;
                await localVideo.play();
                statusMessage = "Camera ready";
            }
            initWebSocket();
        } catch (err) {
            console.error("Error initializing camera:", err);
            errorMessage =
                "Failed to access camera/microphone. Please check permissions.";
            return () => {};
        }

        return (): void => {
            if (cleanupComplete) return;
            cleanupComplete = true;
            stopWebRTC();
            if (localStream) {
                localStream.getTracks().forEach((track) => track.stop());
                localStream = null;
            }
            if (remoteStream) {
                remoteStream.getTracks().forEach((track) => track.stop());
                remoteStream = null;
            }
            if (remoteVideo?.srcObject) {
                const tracks = (
                    remoteVideo.srcObject as MediaStream
                ).getTracks();
                tracks.forEach((t) => t.stop());
                remoteVideo.srcObject = null;
            }
        };
    });

    async function toggleScreenShare() {
        if (isSharingScreen) {
            // revert to camera
            if (localStream) {
                const videoTracks = localStream.getVideoTracks();
                videoTracks.forEach((t) => t.stop());
                localStream = null;
            }
            // re-get camera
            try {
                const stream = await navigator.mediaDevices.getUserMedia({
                    video: true,
                    audio: false,
                });
                localStream = stream;
                if (localVideo) {
                    localVideo.srcObject = stream;
                    await localVideo.play();
                }
                // re-add tracks if peerConnection already exists?
                // for simplicity restart connection
                if (isConnected) {
                    stopWebRTC();
                    startWebRTC();
                }
            } catch (err) {
                console.error("Error switching back to camera:", err);
            }
        } else {
            try {
                const screenStream =
                    await navigator.mediaDevices.getDisplayMedia({
                        video: true,
                        audio: false,
                    });
                if (localStream) {
                    // stop old video track
                    const videoTracks = localStream.getVideoTracks();
                    videoTracks.forEach((t) => t.stop());
                }
                const newStream = new MediaStream();
                // keep audio tracks
                if (localStream) {
                    localStream.getAudioTracks().forEach((t) => {
                        newStream.addTrack(t);
                    });
                }
                const screenTrack = screenStream.getVideoTracks()[0];
                if (screenTrack) {
                    newStream.addTrack(screenTrack);
                    if (localVideo) {
                        localVideo.srcObject = newStream;
                        await localVideo.play();
                    }
                    screenTrack.onended = () => {
                        if (isSharingScreen) {
                            toggleScreenShare();
                        }
                    };
                }
                localStream = newStream;
                if (isConnected) {
                    stopWebRTC();
                    startWebRTC();
                }
            } catch (err) {
                console.error("Error sharing screen:", err);
            }
        }
        isSharingScreen = !isSharingScreen;
    }

    function toggleMute() {
        if (localStream) {
            const audioTracks = localStream.getAudioTracks();
            if (audioTracks.length > 0) {
                isMuted = !isMuted;
                audioTracks[0].enabled = !isMuted;
            }
        }
    }

    async function handleSubmit() {
        if (!formData.message.trim()) {
            console.error("Please enter a message");
            return;
        }
        isSubmitting = true;
        try {
            const response = await fetch("/api/test", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    ...formData,
                    timestamp: new Date().toISOString(),
                    videoActive: !isVideoOff,
                    audioActive: !isMuted,
                }),
            });
            if (!response.ok) {
                throw new Error("Failed to submit data");
            }
            const result = await response.json();
            console.log("Data submitted successfully!", result);
            formData.message = "";
        } catch (error) {
            console.error("Error submitting data:", error);
        } finally {
            isSubmitting = false;
        }
    }

    function endCall() {
        if (localStream) {
            localStream.getTracks().forEach((track) => track.stop());
        }
        window.location.href = "/";
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

                    <button
                        on:click={toggleScreenShare}
                        class={`p-3 rounded-full ${isSharingScreen ? "bg-blue-500" : "bg-gray-700 hover:bg-gray-600"} text-white transition-colors`}
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
