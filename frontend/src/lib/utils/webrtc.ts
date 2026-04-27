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
      // Default to true, will be set specifically by the UI
      this.polite = true
      console.log(`WebRTC Service initialized. Room: ${this.roomId}`)
   }

   public setPolite(polite: boolean) {
      this.polite = polite
      console.log(`WebRTC Politeness set to: ${this.polite}`)
   }

   private sendSignalingMessage(type: string, data: any) {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
         this.socket.send(
            JSON.stringify({
               type,
               roomId: this.roomId,
               data,
            }),
         )
      } else {
         console.warn(`Cannot send signaling message ${type}: socket not open`)
      }
   }

   public async startLocalStream(): Promise<MediaStream> {
      try {
         this.localStream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true,
         })

         // If we already have a peer connection, add tracks now
         if (this.pc && this.pc.signalingState !== 'closed') {
            const senders = this.pc.getSenders()
            this.localStream.getTracks().forEach((track) => {
               const alreadyExists = senders.some((s) => s.track === track)
               if (!alreadyExists) {
                  this.pc?.addTrack(track, this.localStream!)
               }
            })
         }

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
         this.onConnectionStateChange('signaling:connected')
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
         this.onConnectionStateChange('signaling:disconnected')
      }

      this.socket.onerror = (error) => {
         console.error('WebRTC signaling error:', error)
         this.onConnectionStateChange('signaling:error')
      }
   }

   private async handleSignalingMessage(message: WebRTCMessage) {
      // Guard: if PC is closed, ignore signaling
      if (this.pc && this.pc.signalingState === 'closed') return

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

            // Add local tracks if they exist and haven't been added yet
            if (this.localStream) {
               const senders = this.pc!.getSenders()
               this.localStream.getTracks().forEach((track) => {
                  const alreadyExists = senders.some((s) => s.track === track)
                  if (!alreadyExists) {
                     this.pc?.addTrack(track, this.localStream!)
                  }
               })
            }

            const answer = await this.pc!.createAnswer()
            await this.pc!.setLocalDescription(answer)

            this.sendSignalingMessage('answer', this.pc!.localDescription)
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
      if (this.pc && this.pc.signalingState !== 'closed') return

      this.pc = new RTCPeerConnection(this.config)

      this.pc.onicecandidate = ({ candidate }) => {
         if (candidate && this.pc && this.pc.signalingState !== 'closed') {
            this.sendSignalingMessage('candidate', candidate)
         }
      }

      this.pc.ontrack = (event) => {
         console.log('Received remote track:', event.track.kind)
         this.remoteStream = event.streams[0] || new MediaStream([event.track])
         this.onRemoteStream(this.remoteStream)
      }

      this.pc.onnegotiationneeded = async () => {
         if (!this.pc || this.pc.signalingState === 'closed') return
         try {
            this.isMakingOffer = true
            await this.pc.setLocalDescription()
            this.sendSignalingMessage('offer', this.pc.localDescription)
         } catch (err) {
            console.error('Error in onnegotiationneeded:', err)
         } finally {
            this.isMakingOffer = false
         }
      }

      this.pc.onconnectionstatechange = () => {
         if (this.pc) {
            console.log(`PeerConnection state: ${this.pc.connectionState}`)
            this.onConnectionStateChange(this.pc.connectionState)
         }
      }

      this.pc.oniceconnectionstatechange = () => {
         if (this.pc) {
            console.log(`ICE Connection state: ${this.pc.iceConnectionState}`)
         }
      }
   }

   public async call() {
      if (!this.pc || this.pc.signalingState === 'closed') {
         this.createPeerConnection()
      }

      if (this.localStream) {
         const senders = this.pc!.getSenders()
         this.localStream.getTracks().forEach((track) => {
            const alreadyExists = senders.some((s) => s.track === track)
            if (!alreadyExists) {
               this.pc?.addTrack(track, this.localStream!)
            }
         })
      }
   }

   private async processQueuedCandidates() {
      if (!this.pc || this.pc.signalingState === 'closed') return
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
         this.localStream = null
      }

      if (this.pc) {
         this.pc.onicecandidate = null
         this.pc.ontrack = null
         this.pc.onnegotiationneeded = null
         this.pc.onconnectionstatechange = null
         this.pc.oniceconnectionstatechange = null

         if (this.pc.signalingState !== 'closed') {
            this.pc.close()
         }
         this.pc = null
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
      }

      this.remoteStream = null
   }
}
