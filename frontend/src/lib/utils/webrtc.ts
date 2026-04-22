export interface WebRTCMessage {
   type: 'offer' | 'answer' | 'candidate' | 'join' | 'leave' | 'incoming-call'
   roomId: string
   from?: string
   to?: string
   data:
      | Record<string, unknown>
      | RTCSessionDescriptionInit
      | RTCIceCandidateInit
}

// STUN server URLs
const GEO_LOC_URL =
   'https://raw.githubusercontent.com/pradt2/always-online-stun/master/geoip_cache.txt'
const IPV4_URL =
   'https://raw.githubusercontent.com/pradt2/always-online-stun/master/valid_ipv4s.txt'
const GEO_USER_URL = 'https://geolocation-db.com/json/'

export class WebRTCService {
   private pc: RTCPeerConnection | null = null
   private localStream: MediaStream | null = null
   private remoteStream: MediaStream | null = null
   private socket: WebSocket | null = null
   private roomId: string
   private onRemoteStream: (stream: MediaStream) => void
   private onConnectionStateChange: (state: string) => void
   private onIncomingCall?: (offer: RTCSessionDescriptionInit) => void
   private iceCandidateQueue: RTCIceCandidateInit[] = []
   private isRemoteDescriptionSet = false

   private config: RTCConfiguration = {
      iceServers: [],
      iceCandidatePoolSize: 10,
      bundlePolicy: 'max-bundle',
      rtcpMuxPolicy: 'require',
   }

   // Method to get the closest STUN server based on geolocation
   private async getClosestStunServer(): Promise<string> {
      try {
         // Fetch geolocation data
         const geoLocsResponse = await fetch(GEO_LOC_URL)
         if (!geoLocsResponse.ok)
            throw new Error('Failed to fetch geo locations')
         const geoLocs = await geoLocsResponse.json()

         const geoUserResponse = await fetch(GEO_USER_URL)
         if (!geoUserResponse.ok)
            throw new Error('Failed to fetch user location')
         const { latitude, longitude } = await geoUserResponse.json()

         // Fetch available STUN servers
         const stunServersResponse = await fetch(IPV4_URL)
         if (!stunServersResponse.ok)
            throw new Error('Failed to fetch STUN servers')
         const stunServersText = await stunServersResponse.text()
         const stunServers = stunServersText.trim().split('\n')

         // Find the closest server
         const serverDistances: [string, number][] = []

         for (const addr of stunServers) {
            const ip = addr.split(':')[0]
            if (geoLocs[ip]) {
               const [stunLat, stunLon] = geoLocs[ip]
               const dist = Math.sqrt(
                  Math.pow(latitude - stunLat, 2) +
                     Math.pow(longitude - stunLon, 2),
               )
               serverDistances.push([addr, dist])
            }
         }

         if (serverDistances.length === 0) {
            return 'stun:stun.l.google.com:19302' // Fallback
         }

         // Sort by distance and get the closest
         serverDistances.sort((a, b) => a[1] - b[1])
         const closestAddr = serverDistances[0][0]

         return closestAddr
            ? `stun:${closestAddr}`
            : 'stun:stun.l.google.com:19302' // Fallback
      } catch (error) {
         console.warn(
            'Failed to get closest STUN server, using fallback:',
            error,
         )
         return 'stun:stun.l.google.com:19302' // Fallback
      }
   }

   constructor(
      roomId: string,
      onRemoteStream: (stream: MediaStream) => void,
      onConnectionStateChange: (state: string) => void,
      onIncomingCall?: (offer: RTCSessionDescriptionInit) => void,
   ) {
      this.roomId = roomId
      this.onRemoteStream = onRemoteStream
      this.onConnectionStateChange = onConnectionStateChange
      this.onIncomingCall = onIncomingCall

      // Initialize with async STUN server selection
      this.initializeStunServers()
   }

   // Initialize STUN servers asynchronously
   private async initializeStunServers() {
      const closestStun = await this.getClosestStunServer()
      this.config.iceServers = [
         { urls: closestStun },
         { urls: 'stun:stun1.l.google.com:19302' },
         { urls: 'stun:stun2.l.google.com:19302' },
         { urls: 'stun:stun3.l.google.com:19302' },
         { urls: 'stun:stun4.l.google.com:19302' },
      ]
   }

   public async startLocalStream(
      constraints?: MediaStreamConstraints,
   ): Promise<MediaStream> {
      try {
         // Use provided constraints or default to video and audio
         const streamConstraints = constraints || {
            video: {
               width: { ideal: 1280 },
               height: { ideal: 720 },
               frameRate: { ideal: 30 },
            },
            audio: true,
         }

         this.localStream =
            await navigator.mediaDevices.getUserMedia(streamConstraints)
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
         } catch (audioError) {
            console.error('Error accessing audio devices:', audioError)
            throw audioError
         }
      }

      // If peer connection already exists, add tracks now
      if (this.pc && this.localStream) {
         const senders = this.pc.getSenders()
         this.localStream.getTracks().forEach((track) => {
            const alreadyAdded = senders.find((s) => s.track === track)
            if (!alreadyAdded) {
               this.pc?.addTrack(track, this.localStream!)
            }
         })
      }

      return this.localStream
   }

   public connect(wsUrl: string) {
      console.log(`[WebRTC] Connecting to signaling: ${wsUrl}`)
      this.socket = new WebSocket(`${wsUrl}?room=${this.roomId}`)

      this.socket.onopen = () => {
         console.log('[WebRTC] Signaling connected')
         this.onConnectionStateChange('connected')
         this.sendMessage('join', {})
      }

      this.socket.onmessage = async (event) => {
         try {
            const message: WebRTCMessage = JSON.parse(event.data)
            console.log(`[WebRTC] Received message: ${message.type}`, message)
            await this.handleSignalingMessage(message)
         } catch (error) {
            console.error('[WebRTC] Error parsing message:', error)
         }
      }

      this.socket.onclose = (e) => {
         console.log('[WebRTC] Signaling closed:', e.code, e.reason)
         this.onConnectionStateChange('disconnected')
         // Attempt to reconnect if not intentionally disconnected
         if (e.code !== 1000) {
            // Normal closure
            this.scheduleReconnect(wsUrl)
         }
      }

      this.socket.onerror = (error) => {
         console.error('[WebRTC] Signaling error:', error)
         this.onConnectionStateChange('error')
      }
   }

   private reconnectTimeoutId: number | null = null

   private scheduleReconnect(wsUrl: string) {
      if (this.reconnectTimeoutId) {
         clearTimeout(this.reconnectTimeoutId)
      }

      this.reconnectTimeoutId = window.setTimeout(() => {
         console.log('[WebRTC] Attempting to reconnect...')
         this.connect(wsUrl)
      }, 5000) // Retry after 5 seconds
   }

   private sendMessage(
      type: WebRTCMessage['type'],
      data:
         | Record<string, unknown>
         | RTCSessionDescriptionInit
         | RTCIceCandidateInit,
   ) {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
         console.log(`[WebRTC] Sending message: ${type}`)
         this.socket.send(
            JSON.stringify({
               type,
               roomId: this.roomId,
               data,
            }),
         )
      } else {
         console.warn(`[WebRTC] Cannot send ${type}: Socket not open`)
      }
   }

   private async handleSignalingMessage(message: WebRTCMessage) {
      switch (message.type) {
         case 'join':
            console.log('[WebRTC] Peer joined the room')
            break
         case 'offer':
            console.log('[WebRTC] Handling offer')
            if (this.onIncomingCall) {
               this.onIncomingCall(message.data as RTCSessionDescriptionInit)
            } else {
               await this.handleOffer(message.data as RTCSessionDescriptionInit)
            }
            break
         case 'answer':
            console.log('[WebRTC] Handling answer')
            await this.handleAnswer(message.data as RTCSessionDescriptionInit)
            break
         case 'candidate':
            console.log('[WebRTC] Handling ICE candidate')
            await this.handleCandidate(message.data as RTCIceCandidateInit)
            break
         case 'leave':
            console.log('[WebRTC] Peer left the room')
            this.onConnectionStateChange('disconnected')
            break
         default:
            console.warn('[WebRTC] Unknown message type:', message.type)
      }
   }

   private async createPeerConnection() {
      if (this.pc) {
         console.log('[WebRTC] PeerConnection already exists, reusing')
         return
      }

      // Ensure STUN servers are initialized
      if (this.config.iceServers.length === 0) {
         console.log('[WebRTC] Initializing STUN servers...')
         await this.initializeStunServers()
      }

      console.log('[WebRTC] Creating RTCPeerConnection')
      this.pc = new RTCPeerConnection(this.config)

      this.pc.onicecandidate = (event) => {
         if (event.candidate) {
            console.log('[WebRTC] Local ICE candidate found')
            this.sendMessage(
               'candidate',
               event.candidate as RTCIceCandidateInit,
            )
         } else {
            console.log('[WebRTC] ICE gathering complete')
         }
      }

      this.pc.ontrack = (event) => {
         console.log(`[WebRTC] Remote track received: ${event.track.kind}`)
         this.remoteStream = event.streams[0] || new MediaStream([event.track])
         this.onRemoteStream(this.remoteStream)
      }

      this.pc.oniceconnectionstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] ICE State Change: ${this.pc.iceConnectionState}`,
            )
            this.onConnectionStateChange(this.pc.iceConnectionState)
         }
      }

      this.pc.onconnectionstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] Connection State Change: ${this.pc.connectionState}`,
            )
            this.onConnectionStateChange(this.pc.connectionState)
         }
      }

      this.pc.onicegatheringstatechange = () => {
         if (this.pc) {
            console.log(
               `[WebRTC] ICE Gathering State: ${this.pc.iceGatheringState}`,
            )
         }
      }

      if (this.localStream) {
         console.log('[WebRTC] Adding local tracks to PeerConnection')
         this.localStream.getTracks().forEach((track) => {
            this.pc?.addTrack(track, this.localStream!)
         })
      }
   }

   public async call() {
      console.log('[WebRTC] Initiating call')
      await this.createPeerConnection()
      const offer = await this.pc!.createOffer()
      console.log('[WebRTC] Offer created, setting local description')
      await this.pc!.setLocalDescription(offer)
      this.sendMessage('offer', offer as RTCSessionDescriptionInit)
   }

   private async handleOffer(offer: RTCSessionDescriptionInit) {
      console.log('[WebRTC] Setting remote description (offer)')
      await this.createPeerConnection()

      try {
         if (!this.pc) {
            throw new Error('PeerConnection not initialized')
         }

         await this.pc.setRemoteDescription(
            new RTCSessionDescription(offer as RTCSessionDescriptionInit),
         )
         this.isRemoteDescriptionSet = true
         console.log(
            '[WebRTC] Remote description set, processing queued candidates',
         )
         await this.processQueuedCandidates()

         const answer = await this.pc.createAnswer()
         console.log('[WebRTC] Answer created, setting local description')
         await this.pc.setLocalDescription(answer)
         this.sendMessage('answer', answer as RTCSessionDescriptionInit)
      } catch (err) {
         console.error('[WebRTC] Error in handleOffer:', err)
         this.onConnectionStateChange('error')
      }
   }

   private async handleAnswer(answer: RTCSessionDescriptionInit) {
      if (this.pc) {
         try {
            console.log('[WebRTC] Setting remote description (answer)')
            await this.pc.setRemoteDescription(
               new RTCSessionDescription(answer),
            )
            this.isRemoteDescriptionSet = true
            await this.processQueuedCandidates()
         } catch (err) {
            console.error('[WebRTC] Error in handleAnswer:', err)
            this.onConnectionStateChange('error')
         }
      } else {
         console.error('[WebRTC] Received answer but no PeerConnection exists')
         this.onConnectionStateChange('error')
      }
   }

   private async handleCandidate(candidate: RTCIceCandidateInit) {
      if (this.pc && this.isRemoteDescriptionSet) {
         try {
            await this.pc.addIceCandidate(new RTCIceCandidate(candidate))
         } catch (e) {
            console.error('Error adding received ice candidate', e)
         }
      } else {
         this.iceCandidateQueue.push(candidate)
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
      this.iceCandidateQueue = []
      if (this.localStream) {
         this.localStream.getTracks().forEach((track) => track.stop())
      }
      if (this.pc) {
         this.pc.close()
      }
      if (this.socket) {
         // Send leave message before closing
         this.sendMessage('leave', {})
         this.socket.close()
      }
      // Clear any pending reconnect attempts
      if (this.reconnectTimeoutId) {
         clearTimeout(this.reconnectTimeoutId)
         this.reconnectTimeoutId = null
      }
      this.pc = null
      this.localStream = null
      this.remoteStream = null
      this.socket = null
   }
}
