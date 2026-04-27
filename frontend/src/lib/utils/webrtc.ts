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
      // Use a stable heuristic for politeness if possible, or random for now
      this.polite = Math.random() > 0.5
      console.log(`WebRTC Service initialized. Polite: ${this.polite}`)
   }

   public setPolite(polite: boolean) {
      this.polite = polite
      console.log(`WebRTC Politeness set to: ${this.polite}`)
   }

   public async startLocalStream(): Promise<MediaStream> {
      try {
         this.localStream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true,
         })
         return this.localStream
      } catch (error) {
         console.warn(
            'Could not start video, falling back to audio only:',
            error,
         )
         try {
            this.localStream = await navigator.mediaDevices.getUserMedia({
               video: false,
               audio: true,
            })
            return this.localStream
         } catch (audioError) {
            console.error('Error accessing audio devices:', audioError)
            throw audioError
         }
      }
   }

   public connect(wsUrl: string) {
      this.socket = new WebSocket(`${wsUrl}?room=${this.roomId}`)

      this.socket.onopen = () => {
         console.log('WebRTC signaling connected')
         this.onConnectionStateChange('connected')
      }

      this.socket.onmessage = async (event) => {
         try {
            const message: WebRTCMessage = JSON.parse(event.data)
            await this.handleSignalingMessage(message)
         } catch (err) {
            console.error('Error handling signaling message:', err)
         }
      }

      this.socket.onclose = () => {
         console.log('WebRTC signaling disconnected')
         this.onConnectionStateChange('disconnected')
         this.disconnect()
      }

      this.socket.onerror = (error) => {
         console.error('WebRTC signaling error:', error)
         this.onConnectionStateChange('error')
      }
   }

   private async handleSignalingMessage(message: WebRTCMessage) {
      if (!this.pc) {
         this.createPeerConnection()
      }

      const description = message.data

      try {
         if (message.type === 'offer') {
            const offerCollision =
               this.isMakingOffer || this.pc!.signalingState !== 'stable'
            this.isIgnoringOffer = !this.polite && offerCollision

            if (this.isIgnoringOffer) {
               console.log('Glare detected: ignoring offer (impolite)')
               return
            }

            console.log('Handling offer')
            await this.pc!.setRemoteDescription(
               new RTCSessionDescription(description),
            )
            this.isRemoteDescriptionSet = true

            if (this.localStream) {
               this.localStream.getTracks().forEach((track) => {
                  this.pc?.addTrack(track, this.localStream!)
               })
            }

            const answer = await this.pc!.createAnswer()
            await this.pc!.setLocalDescription(answer)

            this.socket?.send(
               JSON.stringify({
                  type: 'answer',
                  roomId: this.roomId,
                  data: this.pc!.localDescription,
               }),
            )
            await this.processQueuedCandidates()
         } else if (message.type === 'answer') {
            console.log('Handling answer')
            await this.pc!.setRemoteDescription(
               new RTCSessionDescription(description),
            )
            this.isRemoteDescriptionSet = true
            await this.processQueuedCandidates()
         } else if (message.type === 'candidate') {
            try {
               if (this.isRemoteDescriptionSet) {
                  await this.pc!.addIceCandidate(
                     new RTCIceCandidate(description),
                  )
               } else {
                  this.iceCandidateQueue.push(description)
               }
            } catch (err) {
               if (!this.isIgnoringOffer) {
                  console.error('Error adding ICE candidate:', err)
               }
            }
         }
      } catch (err) {
         console.error('Error in signaling state machine:', err)
      }
   }

   private createPeerConnection() {
      if (this.pc) return

      this.pc = new RTCPeerConnection(this.config)

      this.pc.onicecandidate = ({ candidate }) => {
         if (candidate && this.socket) {
            this.socket.send(
               JSON.stringify({
                  type: 'candidate',
                  roomId: this.roomId,
                  data: candidate,
               }),
            )
         }
      }

      this.pc.ontrack = (event) => {
         console.log('Received remote track:', event.track.kind)
         this.remoteStream = event.streams[0] || new MediaStream([event.track])
         this.onRemoteStream(this.remoteStream)
      }

      this.pc.onnegotiationneeded = async () => {
         try {
            this.isMakingOffer = true
            await this.pc!.setLocalDescription()
            this.socket?.send(
               JSON.stringify({
                  type: 'offer',
                  roomId: this.roomId,
                  data: this.pc!.localDescription,
               }),
            )
         } catch (err) {
            console.error('Error in onnegotiationneeded:', err)
         } finally {
            this.isMakingOffer = false
         }
      }

      this.pc.onconnectionstatechange = () => {
         if (this.pc) {
            this.onConnectionStateChange(this.pc.connectionState)
         }
      }
   }

   public async call() {
      if (!this.pc) {
         this.createPeerConnection()
      }

      if (this.localStream) {
         this.localStream.getTracks().forEach((track) => {
            this.pc?.addTrack(track, this.localStream!)
         })
      }
   }

   private async processQueuedCandidates() {
      if (!this.pc) return
      while (this.iceCandidateQueue.length > 0) {
         const candidate = this.iceCandidateQueue.shift()
         if (candidate) {
            try {
               await this.pc.addIceCandidate(new RTCIceCandidate(candidate))
            } catch (e) {
               console.error('Error adding queued ice candidate', e)
            }
         }
      }
   }

   public disconnect() {
      this.isRemoteDescriptionSet = false
      this.isMakingOffer = false
      this.isIgnoringOffer = false
      this.iceCandidateQueue = []
      if (this.localStream) {
         this.localStream.getTracks().forEach((track) => track.stop())
      }
      if (this.pc) {
         this.pc.onicecandidate = null
         this.pc.ontrack = null
         this.pc.onnegotiationneeded = null
         this.pc.onconnectionstatechange = null
         this.pc.close()
      }
      if (this.socket) {
         this.socket.close()
      }
      this.pc = null
      this.localStream = null
      this.remoteStream = null
      this.socket = null
   }
}
