export interface WebRTCMessage {
   type: 'offer' | 'answer' | 'candidate' | 'join' | 'leave'
   roomId: string
   from?: string
   to?: string
   data: any
}

export class WebRTCService {
   private pc: RTCPeerConnection | null = null
   private localStream: MediaStream | null = null
   private remoteStream: MediaStream | null = null
   private socket: WebSocket | null = null
   private roomId: string
   private onRemoteStream: (stream: MediaStream) => void
   private onConnectionStateChange: (state: string) => void
   private iceCandidateQueue: RTCIceCandidateInit[] = []
   private isRemoteDescriptionSet = false
   private isMakingOffer = false
   private isIgnoringOffer = false
   private polite: boolean

   private config: RTCConfiguration = {
      iceServers: [
         { urls: 'stun:stun.l.google.com:19302' },
         {
            urls: [`turn:skillswap.online:3478`, `turns:skillswap.online:5349`],
            username: 'skillswap',
            credential: 'skillswap-turn-password',
         },
      ],
      iceTransportPolicy: 'all',
   }

   constructor(
      roomId: string,
      onRemoteStream: (stream: MediaStream) => void,
      onConnectionStateChange: (state: string) => void,
   ) {
      this.roomId = roomId
      this.onRemoteStream = onRemoteStream
      this.onConnectionStateChange = onConnectionStateChange
      this.polite = true
      console.log(`[WebRTC] Service initialized for room: ${this.roomId}`)
   }

   public setPolite(polite: boolean) {
      this.polite = polite
      console.log(`[WebRTC] Politeness set to: ${this.polite}`)
   }

   private sendSignalingMessage(type: string, data: any) {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
         console.log(`[WebRTC] Sending signaling message: ${type}`)
         this.socket.send(
            JSON.stringify({
               type,
               roomId: this.roomId,
               data,
            }),
         )
      } else {
         console.warn(
            `[WebRTC] Cannot send signaling message ${type}: socket not open (state: ${this.socket?.readyState})`,
         )
      }
   }

   public async startLocalStream(): Promise<MediaStream> {
      console.log('[WebRTC] Requesting local stream...')
      try {
         this.localStream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true,
         })
         console.log('[WebRTC] Local stream acquired')

         if (this.pc && this.pc.signalingState !== 'closed') {
            console.log('[WebRTC] PC exists, adding tracks from new stream')
            this.addTracksToPC()
         }

         return this.localStream
      } catch (error) {
         console.warn(
            '[WebRTC] Could not start video, falling back to audio only:',
            error,
         )
         try {
            this.localStream = await navigator.mediaDevices.getUserMedia({
               video: false,
               audio: true,
            })
            console.log('[WebRTC] Local audio stream acquired (fallback)')
            return this.localStream
         } catch (audioError) {
            console.error('[WebRTC] Error accessing audio devices:', audioError)
            throw audioError
         }
      }
   }

   private addTracksToPC() {
      if (!this.pc || !this.localStream) return

      const senders = this.pc.getSenders()
      this.localStream.getTracks().forEach((track) => {
         const alreadyExists = senders.some((s) => s.track === track)
         if (!alreadyExists) {
            console.log(`[WebRTC] Adding track: ${track.kind}`)
            this.pc?.addTrack(track, this.localStream!)
         }
      })
   }

   public connect(wsUrl: string) {
      console.log(`[WebRTC] Connecting to signaling: ${wsUrl}`)
      this.socket = new WebSocket(`${wsUrl}?room=${this.roomId}`)

      this.socket.onopen = () => {
         console.log('[WebRTC] Signaling socket opened')
         this.onConnectionStateChange('signaling:connected')
      }

      this.socket.onmessage = async (event) => {
         try {
            const message: WebRTCMessage = JSON.parse(event.data)
            console.log(
               `[WebRTC] Received signaling message: ${message.type} from ${message.from}`,
            )
            await this.handleSignalingMessage(message)
         } catch (err) {
            console.error('[WebRTC] Error handling signaling message:', err)
         }
      }

      this.socket.onclose = (event) => {
         console.log(
            `[WebRTC] Signaling socket closed: ${event.code} ${event.reason}`,
         )
         this.onConnectionStateChange('signaling:disconnected')
      }

      this.socket.onerror = (error) => {
         console.error('[WebRTC] Signaling socket error:', error)
         this.onConnectionStateChange('signaling:error')
      }
   }

   private async handleSignalingMessage(message: WebRTCMessage) {
      if (this.pc && this.pc.signalingState === 'closed') {
         console.warn('[WebRTC] Received signaling message but PC is closed')
         return
      }

      if (!this.pc) {
         console.log(
            "[WebRTC] PC doesn't exist, creating one for incoming message",
         )
         this.createPeerConnection()
      }

      const description = message.data

      try {
         if (message.type === 'offer') {
            const offerCollision =
               this.isMakingOffer || this.pc!.signalingState !== 'stable'
            this.isIgnoringOffer = !this.polite && offerCollision

            if (this.isIgnoringOffer) {
               console.log('[WebRTC] Glare detected: ignoring offer (impolite)')
               return
            }

            console.log('[WebRTC] Setting remote description (offer)')
            await this.pc!.setRemoteDescription(
               new RTCSessionDescription(description),
            )
            this.isRemoteDescriptionSet = true

            this.addTracksToPC()

            console.log('[WebRTC] Creating answer')
            const answer = await this.pc!.createAnswer()
            console.log('[WebRTC] Setting local description (answer)')
            await this.pc!.setLocalDescription(answer)

            this.sendSignalingMessage('answer', this.pc!.localDescription)
            await this.processQueuedCandidates()
         } else if (message.type === 'answer') {
            console.log('[WebRTC] Setting remote description (answer)')
            await this.pc!.setRemoteDescription(
               new RTCSessionDescription(description),
            )
            this.isRemoteDescriptionSet = true
            await this.processQueuedCandidates()
         } else if (message.type === 'candidate') {
            try {
               if (this.isRemoteDescriptionSet) {
                  console.log('[WebRTC] Adding ICE candidate')
                  await this.pc!.addIceCandidate(
                     new RTCIceCandidate(description),
                  )
               } else {
                  console.log(
                     '[WebRTC] Queuing ICE candidate (remote description not set)',
                  )
                  this.iceCandidateQueue.push(description)
               }
            } catch (err) {
               if (!this.isIgnoringOffer) {
                  console.error('[WebRTC] Error adding ICE candidate:', err)
               }
            }
         }
      } catch (err) {
         console.error('[WebRTC] Error in signaling state machine:', err)
      }
   }

   private createPeerConnection() {
      if (this.pc && this.pc.signalingState !== 'closed') {
         console.log(
            '[WebRTC] createPeerConnection called but PC already exists and is not closed',
         )
         return
      }

      console.log('[WebRTC] Creating RTCPeerConnection...')
      this.pc = new RTCPeerConnection(this.config)
      console.log(
         `[WebRTC] PC created. Initial state: ${this.pc.connectionState}`,
      )

      this.pc.onicecandidate = ({ candidate }) => {
         if (candidate) {
            console.log(
               `[WebRTC] Local ICE candidate generated: ${candidate.candidate.substring(0, 30)}...`,
            )
            if (this.pc && this.pc.signalingState !== 'closed') {
               this.sendSignalingMessage('candidate', candidate)
            }
         } else {
            console.log('[WebRTC] ICE gathering complete (null candidate)')
         }
      }

      this.pc.ontrack = (event) => {
         console.log(`[WebRTC] Received remote track: ${event.track.kind}`)
         this.remoteStream = event.streams[0] || new MediaStream([event.track])
         this.onRemoteStream(this.remoteStream)
      }

      this.pc.onnegotiationneeded = async () => {
         console.log('[WebRTC] Negotiation needed')
         if (!this.pc || this.pc.signalingState === 'closed') return
         try {
            this.isMakingOffer = true
            console.log('[WebRTC] Creating offer')
            const offer = await this.pc.createOffer()
            console.log('[WebRTC] Setting local description (offer)')
            await this.pc.setLocalDescription(offer)
            this.sendSignalingMessage('offer', this.pc.localDescription)
         } catch (err) {
            console.error('[WebRTC] Error in onnegotiationneeded:', err)
         } finally {
            this.isMakingOffer = false
         }
      }

      this.pc.onconnectionstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] PeerConnection state: ${this.pc.connectionState}`,
            )
            this.onConnectionStateChange(this.pc.connectionState)
         }
      }

      this.pc.oniceconnectionstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] ICE Connection state: ${this.pc.iceConnectionState}`,
            )
         }
      }

      this.pc.onicegatheringstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] ICE Gathering state: ${this.pc.iceGatheringState}`,
            )
         }
      }

      this.pc.onsignalingstatechange = () => {
         if (this.pc) {
            console.log(`[WebRTC] Signaling state: ${this.pc.signalingState}`)
         }
      }
   }

   public async call() {
      console.log('[WebRTC] call() initiated')
      if (!this.pc || this.pc.signalingState === 'closed') {
         this.createPeerConnection()
      }

      this.addTracksToPC()

      // If no tracks were added but we want to force negotiation (e.g. for data channel)
      // or if onnegotiationneeded didn't fire for some reason
      if (this.pc!.signalingState === 'stable' && !this.isMakingOffer) {
         console.log('[WebRTC] Manually triggering negotiation')
         this.pc!.onnegotiationneeded?.(new Event('negotiationneeded'))
      }
   }

   private async processQueuedCandidates() {
      if (!this.pc || this.pc.signalingState === 'closed') return
      console.log(
         `[WebRTC] Processing ${this.iceCandidateQueue.length} queued candidates`,
      )
      while (this.iceCandidateQueue.length > 0) {
         const candidate = this.iceCandidateQueue.shift()
         if (candidate) {
            try {
               await this.pc.addIceCandidate(new RTCIceCandidate(candidate))
            } catch (e) {
               console.error('[WebRTC] Error adding queued ice candidate', e)
            }
         }
      }
   }

   public disconnect() {
      console.log('[WebRTC] Disconnecting service...')
      this.isRemoteDescriptionSet = false
      this.isMakingOffer = false
      this.isIgnoringOffer = false
      this.iceCandidateQueue = []

      if (this.localStream) {
         console.log('[WebRTC] Stopping local stream tracks')
         this.localStream.getTracks().forEach((track) => track.stop())
         this.localStream = null
      }

      if (this.pc) {
         this.pc.onicecandidate = null
         this.pc.ontrack = null
         this.pc.onnegotiationneeded = null
         this.pc.onconnectionstatechange = null
         this.pc.oniceconnectionstatechange = null
         this.pc.onicegatheringstatechange = null
         this.pc.onsignalingstatechange = null

         if (this.pc.signalingState !== 'closed') {
            this.pc.close()
         }
         this.pc = null
         console.log('[WebRTC] PeerConnection closed')
      }

      if (this.socket) {
         this.socket.onopen = null
         this.socket.onmessage = null
         this.socket.onclose = null
         this.socket.onerror = null

         if (
            this.socket.readyState === WebSocket.OPEN ||
            this.socket.readyState === WebSocket.CONNECTING
         ) {
            this.socket.close()
         }
         this.socket = null
         console.log('[WebRTC] Signaling socket closed')
      }

      this.remoteStream = null
   }
}
