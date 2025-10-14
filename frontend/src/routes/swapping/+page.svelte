<script lang="ts">
   import { onMount, onDestroy } from 'svelte'
   import { auth } from '$lib/stores/auth'
   import type {
      ChatWithMessages,
      Message,
      WebSocketChatMessage,
   } from '$lib/types/chat'
   import ChatList from '$lib/components/chat/ChatList.svelte'
   import ChatWindow from '$lib/components/chat/ChatWindow.svelte'
   import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte'

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

   function getWebSocketUrl(): string {
      if (typeof window !== 'undefined') {
         const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
         const host = window.location.hostname
         const port = '8080'
         return `${protocol}//${host}:${port}/api/chat`
      }
      return 'ws://localhost:8080/api/chat'
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
      updateChat()
      initializeWebSocket()
   })

   onDestroy(() => {
      if (reconnectTimeout) {
         clearTimeout(reconnectTimeout)
      }
      if (socket) {
         reconnectAttempts = maxReconnectAttempts // Prevent reconnection
         socket.close()
      }
   })
</script>

<div class="h-screen w-full p-2 sm:p-4 bg-gray-100">
   {#if loading}
      <div class="flex items-center justify-center h-full">
         <LoadingSpinner size="xl" text="Loading conversations..." />
      </div>
   {:else}
      <div class="grid grid-cols-1 lg:grid-cols-5 h-full gap-2 sm:gap-4">
         <!-- Chat List - Hidden on mobile when chat selected -->
         <div class="lg:col-span-1 h-full {selectedChat ? 'hidden lg:block' : ''}">
            <ChatList
               {chats}
               selectedChatId={selectedChat?.id || null}
               currentUserId={$auth.user?.id || 0}
               onSelectChat={selectChat}
            />
         </div>

         <!-- Chat Window + Video -->
         <div class="lg:col-span-4 flex flex-col h-full gap-2 sm:gap-4 {selectedChat ? '' : 'hidden lg:flex'}">
            <!-- Video Call Preview - Hidden on mobile -->
            <div
               class="bg-gray-800 rounded-xl shadow-lg h-32 sm:h-48 lg:h-64 flex-shrink-0 overflow-hidden relative hidden sm:block"
            >
               <div
                  class="h-full w-full flex items-center justify-center text-white"
               >
                  <div class="text-center px-4">
                     <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-8 w-8 sm:h-16 sm:w-16 mx-auto mb-2 sm:mb-4 opacity-50"
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
                     <p class="text-sm sm:text-lg font-semibold opacity-75">
                        Video Call Preview
                     </p>
                     <p class="text-xs sm:text-sm opacity-50 mt-1">
                        Video calling feature coming soon
                     </p>
                  </div>
               </div>
               <!-- Connection Status Badge -->
               <div class="absolute top-2 sm:top-4 right-2 sm:right-4">
                  <span
                     class="px-2 sm:px-3 py-0.5 sm:py-1 rounded-full text-xs font-medium {connectionStatus ===
                     'connected'
                        ? 'bg-green-500 text-white'
                        : connectionStatus === 'error'
                          ? 'bg-red-500 text-white'
                          : 'bg-yellow-500 text-white'}"
                  >
                     {connectionStatus}
                  </span>
               </div>
            </div>

            <!-- Chat Messages -->
            <div class="flex-1 min-h-0 relative">
               {#if selectedChat}
                  <!-- Back button for mobile -->
                  <button
                     onclick={() => selectedChatIndex = -1}
                     class="lg:hidden absolute top-2 left-2 z-10 bg-white rounded-full p-2 shadow-md"
                  >
                     <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
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
                           class="w-16 h-16 sm:w-24 sm:h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                           <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-8 w-8 sm:h-12 sm:w-12 text-gray-400"
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
                        <h3 class="text-lg sm:text-xl font-semibold text-gray-900 mb-2">
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
