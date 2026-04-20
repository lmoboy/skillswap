
export interface WebRTCMessage {
    type: 'offer' | 'answer' | 'candidate' | 'join' | 'leave';
    roomId: string;
    from?: string;
    to?: string;
    data: any;
}

export class WebRTCService {
    private pc: RTCPeerConnection | null = null;
    private localStream: MediaStream | null = null;
    private remoteStream: MediaStream | null = null;
    private socket: WebSocket | null = null;
    private roomId: string;
    private onRemoteStream: (stream: MediaStream) => void;
    private onConnectionStateChange: (state: string) => void;
    private iceCandidateQueue: RTCIceCandidateInit[] = [];
    private isRemoteDescriptionSet = false;

    private config: RTCConfiguration = {
        iceServers: [
            { urls: 'stun:stun.l.google.com:19302' },
            { urls: 'stun:stun1.l.google.com:19302' },
            { urls: 'stun:stun2.l.google.com:19302' },
        ]
    };

    constructor(roomId: string, onRemoteStream: (stream: MediaStream) => void, onConnectionStateChange: (state: string) => void) {
        this.roomId = roomId;
        this.onRemoteStream = onRemoteStream;
        this.onConnectionStateChange = onConnectionStateChange;
    }

    public async startLocalStream(): Promise<MediaStream> {
        try {
            // Try video + audio first
            this.localStream = await navigator.mediaDevices.getUserMedia({
                video: true,
                audio: true
            });
            return this.localStream;
        } catch (error) {
            console.warn('Could not start video, falling back to audio only:', error);
            try {
                // Fallback to audio only
                this.localStream = await navigator.mediaDevices.getUserMedia({
                    video: false,
                    audio: true
                });
                return this.localStream;
            } catch (audioError) {
                console.error('Error accessing audio devices:', audioError);
                throw audioError;
            }
        }
    }

    public connect(wsUrl: string) {
        this.socket = new WebSocket(`${wsUrl}?room=${this.roomId}`);

        this.socket.onopen = () => {
            console.log('WebRTC signaling connected');
            this.onConnectionStateChange('connected');
        };

        this.socket.onmessage = async (event) => {
            const message: WebRTCMessage = JSON.parse(event.data);
            await this.handleSignalingMessage(message);
        };

        this.socket.onclose = () => {
            console.log('WebRTC signaling disconnected');
            this.onConnectionStateChange('disconnected');
        };

        this.socket.onerror = (error) => {
            console.error('WebRTC signaling error:', error);
            this.onConnectionStateChange('error');
        };
    }

    private async handleSignalingMessage(message: WebRTCMessage) {
        switch (message.type) {
            case 'offer':
                await this.handleOffer(message.data);
                break;
            case 'answer':
                await this.handleAnswer(message.data);
                break;
            case 'candidate':
                await this.handleCandidate(message.data);
                break;
            default:
                console.warn('Unknown signaling message type:', message.type);
        }
    }

    private createPeerConnection() {
        this.pc = new RTCPeerConnection(this.config);

        this.pc.onicecandidate = (event) => {
            if (event.candidate && this.socket) {
                this.socket.send(JSON.stringify({
                    type: 'candidate',
                    roomId: this.roomId,
                    data: event.candidate
                }));
            }
        };

        this.pc.ontrack = (event) => {
            console.log('Received remote track:', event.track.kind);
            this.remoteStream = event.streams[0] || new MediaStream([event.track]);
            this.onRemoteStream(this.remoteStream);
        };

        this.pc.onconnectionstatechange = () => {
            if (this.pc) {
                this.onConnectionStateChange(this.pc.connectionState);
            }
        };

        if (this.localStream) {
            this.localStream.getTracks().forEach(track => {
                this.pc?.addTrack(track, this.localStream!);
            });
        }
    }

    public async call() {
        this.createPeerConnection();
        const offer = await this.pc!.createOffer();
        await this.pc!.setLocalDescription(offer);

        if (this.socket) {
            this.socket.send(JSON.stringify({
                type: 'offer',
                roomId: this.roomId,
                data: offer
            }));
        }
    }

    private async handleOffer(offer: RTCSessionDescriptionInit) {
        if (!this.pc) {
            this.createPeerConnection();
        }

        await this.pc!.setRemoteDescription(new RTCSessionDescription(offer));
        this.isRemoteDescriptionSet = true;
        await this.processQueuedCandidates();

        const answer = await this.pc!.createAnswer();
        await this.pc!.setLocalDescription(answer);

        if (this.socket) {
            this.socket.send(JSON.stringify({
                type: 'answer',
                roomId: this.roomId,
                data: answer
            }));
        }
    }

    private async handleAnswer(answer: RTCSessionDescriptionInit) {
        if (this.pc) {
            await this.pc.setRemoteDescription(new RTCSessionDescription(answer));
            this.isRemoteDescriptionSet = true;
            await this.processQueuedCandidates();
        }
    }

    private async handleCandidate(candidate: RTCIceCandidateInit) {
        if (this.pc && this.isRemoteDescriptionSet) {
            try {
                await this.pc.addIceCandidate(new RTCIceCandidate(candidate));
            } catch (e) {
                console.error('Error adding received ice candidate', e);
            }
        } else {
            this.iceCandidateQueue.push(candidate);
        }
    }

    private async processQueuedCandidates() {
        if (!this.pc) return;
        while (this.iceCandidateQueue.length > 0) {
            const candidate = this.iceCandidateQueue.shift();
            if (candidate) {
                try {
                    await this.pc.addIceCandidate(new RTCIceCandidate(candidate));
                } catch (e) {
                    console.error('Error adding queued ice candidate', e);
                }
            }
        }
    }

    public disconnect() {
        this.isRemoteDescriptionSet = false;
        this.iceCandidateQueue = [];
        if (this.localStream) {
            this.localStream.getTracks().forEach(track => track.stop());
        }
        if (this.pc) {
            this.pc.close();
        }
        if (this.socket) {
            this.socket.close();
        }
        this.pc = null;
        this.localStream = null;
        this.remoteStream = null;
        this.socket = null;
    }
}
