<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "$lib/stores/auth";
    import Debug from "$lib/components/Debug.svelte";

    let socket = new WebSocket("ws://localhost:8080/api/chat");

    type User = {
        email: string | undefined;
        id: number;
        username: string;
        profile_picture: string;
    };

    type Message = {
        sender: User;
        content: string;
        timestamp: string;
    };

    type ChatMeta = {
        id: number;
        Initiator: number;
        Responder: number;
        Created_at: string;
    };

    type ChatWithMessages = ChatMeta & {
        messages: Message[];
    };

    let newMessage = $state("");
    let chats: ChatWithMessages[] = $state([]);
    let selectedChatIndex: number = $state(-1);

    async function updateChat() {
        const uid = $auth.user?.id;
        if (!uid) return;

        const resp = await fetch(`/api/getChats?uid=${uid}`);
        const chatMetas: ChatMeta[] = await resp.json();

        const chatPromises = chatMetas.map(async (cm) => {
            console.log(cm);
            const res2 = await fetch(`/api/getChatInfo?cid=${cm.id}`);
            const body2 = await res2.json();
            const msgs: Message[] = body2.messages ?? [];
            return {
                ...cm,
                messages: msgs,
            } as ChatWithMessages;
        });

        const chatsWithMsgs = await Promise.all(chatPromises);
        chats = chatsWithMsgs;
    }

    onMount(async () => {
        updateChat();

        if (chats.length > 0) {
            selectedChatIndex = 0;
        }
    });

    function selectChat(i: number) {
        selectedChatIndex = i;
        socket.send(JSON.stringify({ type: "update", id: i }));
    }

    function handleMessage() {
        socket.send(
            JSON.stringify({
                type: "update",
                id: selectedChatIndex,
            }),
        );
        newMessage = "";
    }

    socket.onopen = () => {
        socket.send(JSON.stringify({ type: "update", id: selectedChatIndex }));
        updateChat();
    };
    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log(message);
        if (message.type === "update") {
            updateChat();
        }
    };
    socket.onclose = (e) => {};
    socket.onerror = () => {};
</script>

<div class="h-screen w-full p-4 bg-gray-100 transition-colors duration-300">
    <!-- <Debug {chats} /> -->

    <div class="grid grid-cols-5 grid-rows-6 h-full w-full gap-4">
        <div
            class="flex flex-col col-span-1 row-span-6 bg-white p-4 gap-4 rounded-xl shadow-lg overflow-y-auto"
        >
            <h2 class="text-xl font-bold text-gray-800 dark:text-white">
                Inbox
            </h2>
            <div class="space-y-4 flex-grow">
                <span
                    class="text-sm font-semibold text-gray-500 dark:text-gray-400"
                    >Chats</span
                >
                <div class="flex flex-col gap-3">
                    {#each chats as chat, i}
                        <!-- svelte-ignore event_directive_deprecated -->
                        <div
                            class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200 cursor-pointer"
                            onclick={() => selectChat(i)}
                        >
                            {#if chat.messages.length > 0}
                                {#key chat.messages[chat.messages.length - 1]}
                                    <img
                                        src={chat.messages[
                                            chat.messages.length - 1
                                        ].sender.profile_picture}
                                        alt={chat.messages[
                                            chat.messages.length - 1
                                        ].sender.username}
                                        class="w-12 h-12 rounded-full ring-2 ring-gray-200 object-cover"
                                    />
                                    <div class="flex-grow min-w-0">
                                        <span
                                            class="text-gray-900 font-medium truncate"
                                        >
                                            {chat.messages[
                                                chat.messages.length - 1
                                            ].sender.username}
                                        </span>
                                        <p
                                            class="text-sm text-gray-600 truncate"
                                        >
                                            {chat.messages[
                                                chat.messages.length - 1
                                            ].content}
                                        </p>
                                    </div>
                                {/key}
                            {:else}
                                <div
                                    class="w-12 h-12 bg-gray-200 rounded-full"
                                ></div>
                                <div class="flex-grow min-w-0">
                                    <span
                                        class="text-gray-900 font-medium truncate"
                                        >No messages</span
                                    >
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        </div>

        <!-- Right area: video + messages + send -->
        <div class="col-span-4 row-span-6 flex flex-col gap-4">
            <div
                class="col-span-4 row-span-2 bg-gray-300 rounded-xl shadow-lg flex-grow overflow-hidden"
            >
                <div
                    class="h-full w-full flex items-center justify-center text-gray-600 font-bold text-2xl"
                >
                    Video Call Preview
                </div>
            </div>

            <div
                class="col-span-4 row-span-3 bg-white rounded-xl p-4 shadow-lg overflow-y-auto flex flex-col-reverse gap-3"
            >
                {#if selectedChatIndex >= 0}
                    {#each chats[selectedChatIndex].messages as message}
                        <div
                            class="flex flex-col p-2 rounded-lg max-w-2/3 {message
                                .sender.email === $auth.user?.email
                                ? 'bg-blue-500 text-white self-end'
                                : 'bg-gray-200 text-gray-800 self-start'}"
                        >
                            {message.content}
                        </div>
                    {/each}
                {:else}
                    <div class="text-gray-500 italic">No messages to show</div>
                {/if}
            </div>

            <div
                class="col-span-4 row-span-1 bg-white rounded-xl shadow-lg p-3 flex items-center gap-3"
            >
                <button
                    class="p-2 rounded-full bg-gray-200 text-gray-600 hover:bg-gray-300 transition-colors duration-200"
                >
                    <!-- icon -->
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0
               00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13.5"
                        />
                    </svg>
                </button>
                <button
                    class="p-2 rounded-full bg-gray-200 text-gray-600 hover:bg-gray-300 transition-colors duration-200"
                >
                    <!-- icon -->
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15
              10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </button>
                <input
                    type="text"
                    placeholder="Send a message..."
                    onsubmit={handleMessage}
                    bind:value={newMessage}
                    class="flex-grow bg-gray-100 text-gray-900 p-3 rounded-lg
                 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                />
                <button
                    class="p-3 rounded-lg bg-blue-500 text-white font-semibold hover:bg-blue-600 transition-colors duration-200"
                    onclick={handleMessage}
                >
                    Send
                </button>
            </div>
        </div>
    </div>
</div>
