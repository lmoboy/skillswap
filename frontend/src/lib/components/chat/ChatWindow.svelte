<script lang="ts">
    import type { Message } from "$lib/types/chat";
    import MessageInput from "./MessageInput.svelte";
    import { formatTime } from "$lib/utils/formatting";
    import LoadingSpinner from "$lib/components/common/LoadingSpinner.svelte";

    type Props = {
        messages?: Message[];
        currentUserId: number | string;
        otherUserName?: string;
        otherUserPicture?: string;
        loading?: boolean;
        onSendMessage?: (message: string) => void;
        onAttachment?: () => void;
        class?: string;
    };

    let {
        messages = [],
        currentUserId,
        otherUserName = "User",
        otherUserPicture = "",
        loading = false,
        onSendMessage,
        onAttachment,
        class: className = ""
    }: Props = $props();

    let messageContainer: HTMLElement;
    let newMessage = $state("");

    function handleSend(message: string) {
        if (onSendMessage) {
            onSendMessage(message);
        }
    }

    function scrollToBottom() {
        if (messageContainer) {
            messageContainer.scrollTop = messageContainer.scrollHeight;
        }
    }

    $effect(() => {
        if (messages.length > 0) {
            setTimeout(scrollToBottom, 100);
        }
    });

    function groupMessagesByDate(messages: Message[]) {
        const groups: { date: string; messages: Message[] }[] = [];
        let currentDate = "";

        messages.forEach((message) => {
            const messageDate = new Date(message.timestamp).toLocaleDateString();

            if (messageDate !== currentDate) {
                currentDate = messageDate;
                groups.push({ date: messageDate, messages: [message] });
            } else {
                groups[groups.length - 1].messages.push(message);
            }
        });

        return groups;
    }

    const messageGroups = $derived(groupMessagesByDate(messages));
</script>

<div class="flex flex-col h-full bg-white rounded-xl shadow-lg overflow-hidden {className}">
    <!-- Header -->
    <div class="flex-shrink-0 p-4 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center gap-3">
            <img
                src={otherUserPicture || '/api/profile/default/picture'}
                alt={otherUserName}
                class="w-10 h-10 rounded-full object-cover ring-2 ring-gray-200"
            />
            <div class="flex-1 min-w-0">
                <h3 class="text-lg font-semibold text-gray-900 truncate">
                    {otherUserName}
                </h3>
                <p class="text-sm text-gray-500">Active now</p>
            </div>
            <button
                class="p-2 hover:bg-gray-200 rounded-full transition-colors duration-200"
                aria-label="More options"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 text-gray-600"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                    />
                </svg>
            </button>
        </div>
    </div>

    <!-- Messages -->
    <div
        bind:this={messageContainer}
        class="flex-1 overflow-y-auto p-4 space-y-4"
    >
        {#if loading}
            <div class="flex items-center justify-center h-full">
                <LoadingSpinner size="lg" text="Loading messages..." />
            </div>
        {:else if messages.length === 0}
            <div class="flex flex-col items-center justify-center h-full text-center">
                <div class="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-10 w-10 text-gray-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
                        />
                    </svg>
                </div>
                <p class="text-gray-600 font-medium">No messages yet</p>
                <p class="text-gray-400 text-sm mt-1">Start the conversation by sending a message</p>
            </div>
        {:else}
            {#each messageGroups as group (group.date)}
                <!-- Date divider -->
                <div class="flex items-center justify-center my-4">
                    <div class="bg-gray-200 text-gray-600 text-xs font-medium px-3 py-1 rounded-full">
                        {group.date}
                    </div>
                </div>

                <!-- Messages for this date -->
                {#each group.messages as message (message.id || message.timestamp)}
                    {@const isCurrentUser = message.sender.id == currentUserId || message.sender.email == currentUserId}

                    <div class="flex {isCurrentUser ? 'justify-end' : 'justify-start'}">
                        <div class="flex flex-col max-w-[70%] {isCurrentUser ? 'items-end' : 'items-start'}">
                            {#if !isCurrentUser}
                                <div class="flex items-center gap-2 mb-1 px-2">
                                    <img
                                        src={message.sender.profile_picture || '/api/profile/default/picture'}
                                        alt={message.sender.username}
                                        class="w-6 h-6 rounded-full object-cover"
                                    />
                                    <span class="text-xs text-gray-500 font-medium">
                                        {message.sender.username}
                                    </span>
                                </div>
                            {/if}

                            <div
                                class="p-3 rounded-lg break-words {isCurrentUser
                                    ? 'bg-blue-500 text-white rounded-br-none'
                                    : 'bg-gray-200 text-gray-800 rounded-bl-none'}"
                            >
                                {message.content}
                            </div>

                            <span class="text-xs text-gray-400 mt-1 px-2">
                                {formatTime(message.timestamp)}
                            </span>
                        </div>
                    </div>
                {/each}
            {/each}
        {/if}
    </div>

    <!-- Input -->
    <MessageInput
        bind:value={newMessage}
        onSend={handleSend}
        onAttachment={onAttachment}
        disabled={loading}
    />
</div>
