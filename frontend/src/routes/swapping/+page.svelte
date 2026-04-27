<script lang="ts">
   import { onMount, onDestroy } from 'svelte'
   import { auth } from '$lib/stores/auth'
   import { goto } from '$app/navigation'
   import type {
      ChatWithMessages,
      Message,
      WebSocketChatMessage,
   } from '$lib/types/chat'
   import ChatList from '$lib/components/chat/ChatList.svelte'
   import ChatWindow from '$lib/components/chat/ChatWindow.svelte'
   import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte'
   import { WebRTCService } from '$lib/utils/webrtc'

   let socket: WebSocket | null = null
   let chats = $state<ChatWithMessages[]>([])
   let selectedChatIndex = $state<number>(-1)
   let selectedChat = $derived(
      selectedChatIndex >= 0 ? chats[selectedChatIndex] : null,
   )
   let loading = $state(true)
   let connectionStatus = $state('disconnected')
   let reconnectAttempts = 0
   let maxReconnectAttempts = 10
   let reconnectTimeout: number | null = null

   // WebRTC state
   let webrtcService = $state<WebRTCService | null>(null)
   let localStream = $state<MediaStream | null>(null)
   let remoteStream = $state<MediaStream | null>(null)
   let videoConnectionStatus = $state('disconnected')
   let localVideoElement: HTMLVideoElement | null = $state(null)
   let remoteVideoElement: HTMLVideoElement | null = $state(null)

   function getWebSocketUrl(path: string = '/api/chat'): string {
      if (typeof window !== 'undefined') {
         const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
         const host = window.location.hostname
         const port = window.location.port
         return `${protocol}//${host}${port ? `:${port}` : ''}${path}`
      }
      return `ws://localhost:5173${path}`
   }

   function initializeWebRTC() {
      if (!selectedChat || !$auth.user) return

      if (webrtcService) {
         webrtcService.disconnect()
      }

      webrtcService = new WebRTCService(
         `chat_${selectedChat.id}`,
         (stream) => {
            console.log('Remote stream received')
            remoteStream = stream
            if (remoteVideoElement) {
               remoteVideoElement.srcObject = stream
            }
         },
         (status) => {
            console.log('WebRTC status:', status)
            videoConnectionStatus = status
         },
      )

      // Use a stable heuristic for politeness: user with smaller ID is polite
      const currentUserId = Number($auth.user.id)
      const otherUserId =
         selectedChat.user1_id == currentUserId
            ? selectedChat.user2_id
            : selectedChat.user1_id
      webrtcService.setPolite(currentUserId < Number(otherUserId))

      const wsUrl = getWebSocketUrl('/api/video')
      webrtcService.connect(wsUrl)
   }

   async function startVideoCall() {
      if (!selectedChat) {
         console.warn('Cannot start call: no chat selected')
         return
      }

      if (!webrtcService) {
         initializeWebRTC()
      }

      if (!webrtcService) {
         console.error('Failed to initialize WebRTC service')
         return
      }

      try {
         const stream = await webrtcService.startLocalStream()
         localStream = stream
         if (localVideoElement) {
            localVideoElement.srcObject = stream
         }
         await webrtcService.call()
      } catch (err) {
         console.error('Failed to start video call:', err)
      }
   }

   function endVideoCall() {
      if (webrtcService) {
         webrtcService.disconnect()
         webrtcService = null
      }
      localStream = null
      remoteStream = null
      videoConnectionStatus = 'disconnected'
   }

   async function updateChat() {
      const uid = $auth.user?.id
      if (!uid) return

      try {
         const resp = await fetch(`/api/getChats?uid=${uid}`)
         const chatMetas = await resp.json()
         const chatPromises = chatMetas.map(async (cm: any) => {
            // console.log(cm)
            const res2 = await fetch(`/api/getChatInfo?cid=${cm.id}`)
            // console.log(res2)
            const body2 = await res2.json()
            console.log(body2)
            const msgs: Message[] = body2.messages ?? []

            return {
               ...cm,
               messages: msgs,
            } as ChatWithMessages
         })
         const chatsWithMsgs = await Promise.all(chatPromises)
         chats = chatsWithMsgs
         loading = false
      } catch (error) {
         console.error('Error loading chats:', error)
         loading = false
      }
   }

   function initializeWebSocket() {
      if (socket) {
         socket.close()
      }

      const wsUrl = getWebSocketUrl()
      console.log('Connecting to WebSocket:', wsUrl)

      try {
         loading = false
         socket = new WebSocket(wsUrl)

         socket.onopen = () => {
            console.log('WebSocket connected')
            connectionStatus = 'connected'
            reconnectAttempts = 0
            if (selectedChat) {
               socket?.send(
                  JSON.stringify({
                     type: 'update',
                     id: selectedChat.id,
                  }),
               )
            }
         }

         socket.onmessage = (event) => {
            try {
               const message: WebSocketChatMessage = JSON.parse(event.data)
               console.log('WebSocket message received:', message)

               if (message.type === 'new_message' && message.message) {
                  const chatId = message.chat_id
                  const chatIndex = chats.findIndex((c) => c.id === chatId)

                  if (chatIndex !== -1) {
                     chats[chatIndex].messages = [
                        ...chats[chatIndex].messages,
                        message.message,
                     ]
                  }
               } else if (message.type === 'update') {
                  updateChat()
               }
            } catch (error) {
               console.error('Error processing WebSocket message:', error)
            }
         }

         socket.onclose = (e) => {
            console.log('WebSocket closed:', e.code, e.reason)
            connectionStatus = 'disconnected'

            // Attempt to reconnect
            if (reconnectAttempts < maxReconnectAttempts) {
               reconnectAttempts++
               const delay = Math.min(
                  1000 * Math.pow(2, reconnectAttempts - 1),
                  30000,
               )
               console.log(
                  `Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`,
               )

               reconnectTimeout = window.setTimeout(() => {
                  initializeWebSocket()
               }, delay)
            } else {
               console.error('Max reconnection attempts reached')
            }
         }

         socket.onerror = (error) => {
            console.error('WebSocket error:', error)
            connectionStatus = 'error'
         }
      } catch (error) {
         console.error('Failed to create WebSocket:', error)
         connectionStatus = 'error'
      }
   }

   function selectChat(chatId: number, index: number) {
      selectedChatIndex = index
      if (socket && socket.readyState === WebSocket.OPEN) {
         socket.send(
            JSON.stringify({
               type: 'update',
               id: chatId,
            }),
         )
      }
   }

   function handleSendMessage(message: string) {
      if (!selectedChat || !socket || socket.readyState !== WebSocket.OPEN) {
         console.error('Cannot send message: chat or socket not ready')
         return
      }

      socket.send(
         JSON.stringify({
            type: 'post',
            id: selectedChat.id,
            user_id: $auth.user?.id,
            content: message,
         }),
      )
   }

   function handleAttachment() {
      console.log('Attachment clicked')
      // TODO: Implement file attachment
   }

   function getOtherUserInfo(chat: ChatWithMessages | null) {
      if (!chat) return { name: '', picture: '' }

      const isUser1 = chat.user1_id == $auth.user?.id
      return {
         name: isUser1 ? chat.user2_username : chat.user1_username,
         picture: isUser1
            ? chat.user2_profile_picture
            : chat.user1_profile_picture,
      }
   }

   onMount(() => {
      if (!$auth.isAuthenticated) {
         goto('/login')
      }
      updateChat()
      initializeWebSocket()
   })

   $effect(() => {
      if (selectedChat) {
         initializeWebRTC()
      }
   })

   onDestroy(() => {
      if (reconnectTimeout) {
         clearTimeout(reconnectTimeout)
      }
      if (socket) {
         reconnectAttempts = maxReconnectAttempts // Prevent reconnection
         socket.close()
      }
      endVideoCall()
   })
</script>

<div
   class="h-[calc(100dvh-64px)] w-full p-2 sm:p-4 bg-gray-100 overflow-hidden"
>
   {#if loading}
      <div class="flex items-center justify-center h-full">
         <LoadingSpinner size="xl" text="Loading conversations..." />
      </div>
   {:else}
      <div
         class="grid grid-cols-1 lg:grid-cols-5 h-full gap-2 sm:gap-4 overflow-hidden"
      >
         <!-- Chat List - Hidden on mobile when chat selected -->
         <div
            class="lg:col-span-1 h-full min-h-0 {selectedChat
               ? 'hidden lg:block'
               : ''}"
         >
            <ChatList
               {chats}
               selectedChatId={selectedChat?.id || null}
               currentUserId={$auth.user?.id || 0}
               onSelectChat={selectChat}
            />
         </div>

         <!-- Chat Window + Video -->
         <div
            class="lg:col-span-4 flex flex-col h-full min-h-0 gap-2 sm:gap-4 {selectedChat
               ? ''
               : 'hidden lg:flex'}"
         >
            <!-- Video Call Container -->
            <div
               class="bg-gray-900 rounded-xl shadow-lg flex-1 min-h-[200px] lg:min-h-[300px] overflow-hidden relative"
            >
               {#if remoteStream}
                  <video
                     bind:this={remoteVideoElement}
                     autoplay
                     playsinline
                     class="absolute inset-0 w-full h-full object-contain bg-gray-900"
                  >
                     <track kind="captions" />
                  </video>
               {:else}
                  <div
                     class="absolute inset-0 flex items-center justify-center text-white bg-gradient-to-br from-gray-800 to-gray-900"
                  >
                     <div class="text-center px-4">
                        <div
                           class="w-12 h-12 sm:w-16 sm:h-16 mx-auto mb-2 rounded-full bg-gray-700 flex items-center justify-center"
                        >
                           <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-6 w-6 sm:h-8 sm:w-8 text-gray-400"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                           >
                              <path
                                 stroke-linecap="round"
                                 stroke-linejoin="round"
                                 stroke-width="2"
                                 d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
                              />
                           </svg>
                        </div>
                        <p class="text-xs sm:text-base font-medium opacity-90">
                           {selectedChat
                              ? `Waiting for ${getOtherUserInfo(selectedChat).name}...`
                              : 'Select a chat to start video'}
                        </p>
                     </div>
                  </div>
               {/if}
               {#if localStream}
                  <div
                     class="absolute bottom-4 right-4 w-24 sm:w-32 lg:w-40 aspect-video bg-black rounded-lg border-2 border-white/50 overflow-hidden shadow-xl z-10"
                  >
                     <video
                        bind:this={localVideoElement}
                        autoplay
                        playsinline
                        muted
                        class="h-full w-full object-cover mirror"
                     >
                        <track kind="captions" />
                     </video>
                  </div>
               {/if}

               <!-- Video Controls Overlay -->
               <div
                  class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-3 z-20"
               >
                  {#if !localStream}
                     <button
                        onclick={startVideoCall}
                        class="bg-green-600 hover:bg-green-700 active:bg-green-800 text-white px-4 py-2 sm:px-5 sm:py-2.5 rounded-full flex items-center gap-2 transition-all shadow-lg text-xs sm:text-sm font-medium hover:scale-105 active:scale-95"
                     >
                        <svg
                           xmlns="http://www.w3.org/2000/svg"
                           class="h-4 w-4"
                           viewBox="0 0 20 20"
                           fill="currentColor"
                        >
                           <path
                              d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z"
                           />
                        </svg>
                        <span>Start Call</span>
                     </button>
                  {:else}
                     <button
                        onclick={endVideoCall}
                        class="bg-red-600 hover:bg-red-700 active:bg-red-800 text-white px-4 py-2 sm:px-5 sm:py-2.5 rounded-full flex items-center gap-2 transition-all shadow-lg text-xs sm:text-sm font-medium hover:scale-105 active:scale-95"
                     >
                        <svg
                           xmlns="http://www.w3.org/2000/svg"
                           class="h-4 w-4"
                           viewBox="0 0 20 20"
                           fill="currentColor"
                        >
                           <path
                              fill-rule="evenodd"
                              d="M3 5a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-10a2 2 0 01-2-2V5zm11 1H6v8l4-2 4 2V6z"
                              clip-rule="evenodd"
                           />
                        </svg>
                        <span>End Call</span>
                     </button>
                  {/if}
               </div>

               <!-- Connection Status Badge -->
               <div class="absolute top-4 right-4 z-20">
                  <span
                     class="px-2 py-1 rounded-full text-[10px] font-medium capitalize shadow-lg {videoConnectionStatus ===
                     'connected'
                        ? 'bg-green-500/90 text-white backdrop-blur-sm'
                        : videoConnectionStatus === 'error'
                          ? 'bg-red-500/90 text-white backdrop-blur-sm'
                          : 'bg-yellow-500/90 text-white backdrop-blur-sm'}"
                  >
                     {videoConnectionStatus}
                  </span>
               </div>
            </div>

            <!-- Chat Messages -->
            <div class="flex-1 min-h-0 relative">
               {#if selectedChat}
                  <!-- Back button for mobile -->
                  <button
                     aria-label="Go back"
                     onclick={() => (selectedChatIndex = -1)}
                     class="lg:hidden absolute top-2 left-2 z-10 bg-white rounded-full p-2 shadow-md"
                  >
                     <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                     >
                        <path
                           stroke-linecap="round"
                           stroke-linejoin="round"
                           stroke-width="2"
                           d="M15 19l-7-7 7-7"
                        />
                     </svg>
                  </button>

                  {@const otherUser = getOtherUserInfo(selectedChat)}
                  <ChatWindow
                     messages={selectedChat.messages}
                     currentUserId={$auth.user?.id || 0}
                     otherUserName={otherUser.name}
                     otherUserPicture={otherUser.picture}
                     onSendMessage={handleSendMessage}
                     onAttachment={handleAttachment}
                     class="h-full"
                  />
               {:else}
                  <div
                     class="h-full bg-white rounded-xl shadow-lg flex items-center justify-center"
                  >
                     <div class="text-center p-4 sm:p-8">
                        <div
                           class="w-16 h-16 sm:w-20 sm:h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                           <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-8 w-8 sm:h-10 sm:w-10 text-gray-400"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                           >
                              <path
                                 stroke-linecap="round"
                                 stroke-linejoin="round"
                                 stroke-width="2"
                                 d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                              />
                           </svg>
                        </div>
                        <h3
                           class="text-lg sm:text-xl font-semibold text-gray-900 mb-2"
                        >
                           Select a conversation
                        </h3>
                        <p class="text-sm sm:text-base text-gray-500">
                           Choose a conversation from the list to start chatting
                        </p>
                     </div>
                  </div>
               {/if}
            </div>
         </div>
      </div>
   {/if}
</div>
