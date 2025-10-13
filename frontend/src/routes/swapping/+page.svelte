<script lang="ts">
   import { onMount } from 'svelte'
   import { auth } from '$lib/stores/auth'
   import type {
      ChatWithMessages,
      Message,
      WebSocketChatMessage,
   } from '$lib/types/chat'
   import ChatList from '$lib/components/chat/ChatList.svelte'
   import ChatWindow from '$lib/components/chat/ChatWindow.svelte'
   import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte'
   import { ReconnectingWebSocket } from '$lib/utils/websocket-helper'

   let socket: any = null
   let chats = $state<ChatWithMessages[]>([])
   let selectedChatIndex = $state<number>(-1)
   let selectedChat = $derived(
      selectedChatIndex >= 0 ? chats[selectedChatIndex] : null,
   )
   let newMessage = $state('')
   let loading = $state(true)
   let connectionStatus = $state('disconnected')

   async function updateChat() {
      const uid = $auth.user?.id
      if (!uid) return

      try {
         const resp = await fetch(`/api/getChats?uid=${uid}`)
         const chatMetas = await resp.json()

         const chatPromises = chatMetas.map(async (cm: any) => {
            const res2 = await fetch(`/api/getChatInfo?cid=${cm.id}`)
            const body2 = await res2.json()
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

      socket = new ReconnectingWebSocket('api/chat')

      socket.onopen = () => {
         console.log('WebSocket connected')
         connectionStatus = 'connected'
         if (selectedChat) {
            socket?.send({
               type: 'update',
               id: selectedChat.id,
            })
         }
      }

      socket.onmessage = (event: { data: string; }) => {
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

      socket.onclose = (e: any) => {
         console.log('WebSocket closed:', e)
         connectionStatus = 'disconnected'
      }

      socket.onerror = (error: any) => {
         console.error('WebSocket error:', error)
         connectionStatus = 'error'
      }
   }

   function selectChat(chatId: number, index: number) {
      selectedChatIndex = index
      if (socket && socket.readyState === WebSocket.OPEN) {
         socket.send({
            type: 'update',
            id: chatId,
         })
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

      // Optimistically add message to UI
      // const optimisticMessage: Message = {
      //     sender: {
      //         id: $auth.user?.id || 0,
      //         username: $auth.user?.name || "",
      //         email: $auth.user?.email || "",
      //         profile_picture: $auth.user?.profile_picture || "",
      //     },
      //     content: message,
      //     timestamp: new Date().toISOString(),
      //     chat_id: selectedChat.id,
      // };

      // chats[selectedChatIndex].messages = [
      //     ...chats[selectedChatIndex].messages,
      //     optimisticMessage,
      // ];
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

      return () => {
         if (socket) {
            socket.close()
         }
      }
   })
</script>

<div class="h-screen w-full p-4 bg-gray-100">
   {#if loading}
      <div class="flex items-center justify-center h-full">
         <LoadingSpinner size="xl" text="Loading conversations..." />
      </div>
   {:else}
      <div class="grid grid-cols-1 lg:grid-cols-5 h-full gap-4">
         <!-- Chat List -->
         <div class="lg:col-span-1 h-full">
            <ChatList
               {chats}
               selectedChatId={selectedChat?.id || null}
               currentUserId={$auth.user?.id || 0}
               onSelectChat={selectChat}
            />
         </div>

         <!-- Chat Window + Video -->
         <div class="lg:col-span-4 flex flex-col h-full gap-4">
            <!-- Video Call Preview -->
            <div
               class="bg-gray-800 rounded-xl shadow-lg h-48 lg:h-64 flex-shrink-0 overflow-hidden"
            >
               <div
                  class="h-full w-full flex items-center justify-center text-white"
               >
                  <div class="text-center">
                     <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-16 w-16 mx-auto mb-4 opacity-50"
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
                     <p class="text-lg font-semibold opacity-75">
                        Video Call Preview
                     </p>
                     <p class="text-sm opacity-50 mt-1">
                        Video calling feature coming soon
                     </p>
                  </div>
               </div>
            </div>

            <!-- Chat Messages -->
            <div class="flex-1 min-h-0">
               {#if selectedChat}
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
                     <div class="text-center p-8">
                        <div
                           class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                           <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="h-12 w-12 text-gray-400"
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
                        <h3 class="text-xl font-semibold text-gray-900 mb-2">
                           Select a conversation
                        </h3>
                        <p class="text-gray-500">
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
