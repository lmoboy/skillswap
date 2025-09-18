<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { Mic, PhoneOff, Share2, VideoOff, Send } from "lucide-svelte";
    import VideoDebug from "$lib/components/VideoDebug.svelte";

    
    let localVideo: HTMLVideoElement;
    let remoteVideo: HTMLVideoElement;

    
    let localStream: MediaStream | null = null;
    let remoteStream: MediaStream | null = null;
    let peerConnection: RTCPeerConnection | null = null;
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

    
    const rtcConfig = {
        iceServers: [
            { urls: "stun:stun.l.google.com:19302" },
            
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
                            console.log("Received answer:", message.data);
                            try {
                                const answer = new RTCSessionDescription(
                                    message.data,
                                );
                                await peerConnection.setRemoteDescription(
                                    answer,
                                );
                                console.log(
                                    "Successfully set remote description (answer)",
                                );
                            } catch (err) {
                                console.error(
                                    "Error setting remote description:",
                                    err,
                                );
                                errorMessage = "Failed to establish connection";
                            }
                            break;

                        case "ice-candidate":
                            if (message.data) {
                                try {
                                    console.log(
                                        "Adding ICE candidate:",
                                        message.data,
                                    );
                                    await peerConnection.addIceCandidate(
                                        new RTCIceCandidate(message.data),
                                    );
                                } catch (err) {
                                    console.error(
                                        "Error adding ICE candidate:",
                                        err,
                                    );
                                }
                            }
                            break;

                        default:
                            console.warn("Unknown message type:", message.type);
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

    
    const preferCodec = (sdp: string, codec: string) => {
        const lines = sdp.split("\n");
        const mLineIndices: number[] = [];

        
        for (let i = 0; i < lines.length; i++) {
            if (
                lines[i].startsWith("m=") &&
                (lines[i].includes("audio") || lines[i].includes("video"))
            ) {
                mLineIndices.push(i);
            }
        }

        
        for (const mLineIndex of mLineIndices) {
            const mLine = lines[mLineIndex];
            const isVideo = mLine.includes("video");

            
            const codecLines: { line: string; priority: number }[] = [];

            for (let i = mLineIndex + 1; i < lines.length; i++) {
                if (lines[i].startsWith("a=rtpmap")) {
                    const payloadType = lines[i].split(" ")[0].split(":")[1];
                    const codecName = lines[i].toLowerCase();
                    let priority = 0;

                    
                    if (isVideo && codecName.includes("h264")) {
                        priority = 1000;
                    } else if (!isVideo && codecName.includes("opus")) {
                        priority = 1000;
                    } else if (codecName.includes("vp8")) {
                        priority = 900;
                    } else if (codecName.includes("vp9")) {
                        priority = 800;
                    }

                    codecLines.push({
                        line: `a=rtpmap:${payloadType}`,
                        priority: priority,
                    });
                } else if (lines[i] === "") {
                    break; 
                }
            }

            
            codecLines.sort((a, b) => b.priority - a.priority);

            
            if (codecLines.length > 0) {
                const newLines: string[] = [];
                const processedPayloads = new Set<string>();

                
                newLines.push(mLine);

                
                for (const codec of codecLines) {
                    if (!processedPayloads.has(codec.line.split(" ")[0])) {
                        newLines.push(codec.line);
                        processedPayloads.add(codec.line.split(" ")[0]);
                    }
                }

                
                for (let i = mLineIndex + 1; i < lines.length; i++) {
                    if (lines[i] === "") {
                        newLines.push("");
                        break;
                    }
                    if (!lines[i].startsWith("a=rtpmap")) {
                        newLines.push(lines[i]);
                    }
                }

                
                lines.splice(mLineIndex, newLines.length, ...newLines);
            }
        }

        return lines.join("\n");
    };

    
    const startWebRTC = async () => {
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.error("Cannot start WebRTC: WebSocket is not connected");
            errorMessage = "Not connected to server";
            return;
        }

        try {
            console.log("Starting WebRTC connection...");

            
            peerConnection = new RTCPeerConnection({
                ...rtcConfig,
                iceTransportPolicy: "all",
                bundlePolicy: "max-bundle",
                rtcpMuxPolicy: "require",
                iceCandidatePoolSize: 10,
            });

            console.log("Created new RTCPeerConnection");

            
            peerConnection.onicecandidate = (event) => {
                if (event.candidate) {
                    console.log("New ICE candidate:", event.candidate);
                    if (socket?.readyState === WebSocket.OPEN) {
                        const candidate = event.candidate.toJSON();
                        socket.send(
                            JSON.stringify({
                                type: "ice-candidate",
                                data: candidate,
                            }),
                        );
                        console.log("Sent ICE candidate to server");
                    }
                } else {
                    console.log("All ICE candidates have been sent");
                }
            };

            
            peerConnection.oniceconnectionstatechange = () => {
                if (!peerConnection) return;

                console.log(
                    "ICE connection state:",
                    peerConnection.iceConnectionState,
                );
                statusMessage = `ICE: ${peerConnection.iceConnectionState}`;

                
                if (
                    peerConnection.iceConnectionState === "failed" ||
                    peerConnection.iceConnectionState === "disconnected"
                ) {
                    console.warn(
                        "ICE connection failed or disconnected, restarting ICE...",
                    );
                    if (peerConnection.iceConnectionState !== "closed") {
                        peerConnection.restartIce();
                    }
                }

                
                if (
                    peerConnection.iceConnectionState === "closed" &&
                    isConnected
                ) {
                    console.log(
                        "ICE connection closed, attempting to reconnect...",
                    );
                    setTimeout(startWebRTC, 2000);
                }
            };

            
            peerConnection.ontrack = (event) => {
                console.log(
                    "Received remote track:",
                    event.track.kind,
                    "with id:",
                    event.track.id,
                );

                if (!remoteStream) {
                    remoteStream = new MediaStream();
                    console.log("Created new remote MediaStream");

                    if (remoteVideo) {
                        remoteVideo.srcObject = remoteStream;
                        remoteVideo.play().catch((err) => {
                            console.error("Error playing remote video:", err);
                            errorMessage = "Error playing video";
                        });
                    }
                }

                
                remoteStream.addTrack(event.track);
                console.log(
                    "Added track to remote stream. Track count:",
                    remoteStream.getTracks().length,
                );

                
                event.track.onended = () => {
                    console.log("Track ended:", event.track.id);
                    if (remoteStream) {
                        remoteStream.removeTrack(event.track);
                    }
                };

                statusMessage = "Connected";
            };

            
            peerConnection.onconnectionstatechange = () => {
                if (!peerConnection) return;

                const state = peerConnection.connectionState;
                console.log("Connection state changed to:", state);
                statusMessage = `Connection: ${state}`;

                if (state === "connected") {
                    isStreaming = true;
                    console.log("WebRTC connection established successfully");
                } else if (
                    ["disconnected", "failed", "closed"].includes(state)
                ) {
                    console.warn("WebRTC connection lost or failed");
                    isStreaming = false;

                    
                    if (isConnected) {
                        console.log("Attempting to reconnect WebRTC...");
                        setTimeout(startWebRTC, 2000);
                    }
                }
            };

            
            if (localStream) {
                console.log("Adding local stream tracks to peer connection");
                for (const track of localStream.getTracks()) {
                    if (peerConnection) {
                        console.log(`Adding local ${track.kind} track`);
                        peerConnection.addTrack(track, localStream);
                    }
                }
            } else {
                console.warn("No local stream available when starting WebRTC");
            }

            
            console.log("Creating offer...");
            const offer = await peerConnection.createOffer({
                offerToReceiveAudio: true,
                offerToReceiveVideo: true,
                voiceActivityDetection: false,
                iceRestart: false,
            });

            
            if (offer.sdp) {
                offer.sdp = preferCodec(offer.sdp, "H264");
            }

            console.log("Setting local description...");
            await peerConnection.setLocalDescription(offer);

            if (socket?.readyState === WebSocket.OPEN) {
                const offerToSend = peerConnection.localDescription?.toJSON();
                if (offerToSend) {
                    socket.send(
                        JSON.stringify({
                            type: "offer",
                            data: offerToSend,
                        }),
                    );
                    console.log("Sent offer to server");
                    statusMessage = "Sent offer, waiting for response...";
                } else {
                    throw new Error("Failed to create offer");
                }
            } else {
                throw new Error("WebSocket is not connected");
            }
        } catch (error) {
            console.error("Error in startWebRTC:", error);
            errorMessage = `Connection error: ${error.message}`;
            stopWebRTC();

            
            if (isConnected) {
                console.log("Retrying WebRTC connection in 2 seconds...");
                setTimeout(startWebRTC, 2000);
            }
        }
    };

    
    const stopWebRTC = () => {
        if (peerConnection) {
            peerConnection.ontrack = null;
            peerConnection.onicecandidate = null;
            peerConnection.oniceconnectionstatechange = null;

            
            if (remoteStream) {
                remoteStream.getTracks().forEach((track) => track.stop());
                remoteStream = null;
            }

            
            peerConnection.close();
            peerConnection = null;
        }

        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.close();
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

    
    onMount(async (): Promise<() => void> => {
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
                tracks.forEach((track) => track.stop());
                remoteVideo.srcObject = null;
            }
        };
    });

    
    async function toggleScreenShare() {
        if (isSharingScreen) {
            
            if (localStream) {
                const videoTrack = localStream.getVideoTracks()[0];
                if (videoTrack) {
                    videoTrack.stop();
                }

                try {
                    const stream = await navigator.mediaDevices.getUserMedia({
                        video: true,
                        audio: false,
                    });

                    const newVideoTrack = stream.getVideoTracks()[0];
                    if (newVideoTrack) {
                        localStream.addTrack(newVideoTrack);
                        if (localVideo) {
                            localVideo.srcObject = localStream;
                        }
                    }
                } catch (err) {
                    console.error("Error switching back to camera:", err);
                }
            }
        } else {
            
            try {
                const screenStream =
                    await navigator.mediaDevices.getDisplayMedia({
                        video: true,
                        audio: false,
                    });

                if (localStream) {
                    const videoTrack = localStream.getVideoTracks()[0];
                    if (videoTrack) {
                        videoTrack.stop();
                        localStream.removeTrack(videoTrack);
                    }

                    const screenTrack = screenStream.getVideoTracks()[0];
                    if (screenTrack) {
                        localStream.addTrack(screenTrack);
                        if (localVideo) {
                            localVideo.srcObject = localStream;
                        }

                        
                        screenTrack.onended = () => {
                            if (isSharingScreen) {
                                toggleScreenShare();
                            }
                        };
                    }
                }
            } catch (err) {
                console.error("Error sharing screen:", err);
                return;
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
            console.log("Data submitted successfully!");
            console.log("Server response:", result);

            formData.message = "";
        } catch (error) {
            console.error("Error submitting data:", error);
            console.error("Failed to submit data");
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
                        muted
                        controls
                    ></video>
                    {#if !remoteVideo?.srcObject}
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
                    {/if}
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
