<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { onMount, onDestroy } from "svelte";

    type Messager = {
        user_name: string;
        user_id: string;
        user_image: string;
        chat_messages: {
            message_sender: string;
            message_content: string;
        }[];
    };

    let swappers = $state([]);
    let swappees = $state([]);
    let chosenMessager = $state(-1);

    onMount(async () => {
        let chats = fetch("/api/getChats?uid=" + $auth.user?.id).then((res) => {
            res.json().then((data) => {
                return data;
            });
        });
        console.log(chats);
        fetch("/api/getChatInfo?cid=" + chats[0].id).then((res) => {
            res.json().then((data) => {
                swappees = data;
            });
        });
    });

    // let swappers = $state([
    //     {
    //         name: "John",
    //         id: "412487213",
    //         lastMessage: "Hey, are you free to swap later?",
    //         image: "https://randomuser.me/api/portraits/men/33.jpg",
    //     },
    //     {
    //         name: $auth.user?.name,
    //         id: "123456789",
    //         lastMessage: "I'm available tomorrow morning.",
    //         image: "https://randomuser.me/api/portraits/women/44.jpg",
    //     },
    // ]);

    // let swappees = $state<Messager>({
    //     name: "Sara",
    //     id: "987654321",
    //     image: "https://randomuser.me/api/portraits/women/55.jpg",
    //     messages: [
    //         { sender: "other", text: "Hey there! Ready to trade?" },
    //         {
    //             sender: authState.user?.name,
    //             text: "Yep, what do you want to swap?",
    //         },
    //         { sender: "other", text: "I have a vintage watch." },
    //         {
    //             sender: authState.user?.name,
    //             text: "Oh, cool! I have a rare comic book.",
    //         },
    //     ],
    // });
</script>

<div class="h-screen w-full p-4 bg-gray-100 transition-colors duration-300">
    <div class="grid grid-cols-5 grid-rows-6 h-full w-full gap-4">
        <!-- <div
            class="flex flex-col col-span-1 row-span-6 bg-white p-4 gap-4 rounded-xl shadow-lg overflow-y-auto"
        >
            <h2 class="text-xl font-bold text-gray-800 dark:text-white">
                Inbox
            </h2>
            <div class="space-y-4 flex-grow">
                <span
                    class="text-sm font-semibold text-gray-500 dark:text-gray-400"
                    >Swappers</span
                >
                <div class="flex flex-col gap-3">
                    {#each swappers as swapper}
                        <div
                            class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200 cursor-pointer"
                        >
                            <img
                                src={swapper.image}
                                alt={swapper.name}
                                class="w-12 h-12 rounded-full ring-2 ring-gray-200 object-cover"
                            />
                            <div class="flex-grow min-w-0">
                                <span class="text-gray-900 font-medium truncate"
                                    >{swapper.name}</span
                                >
                                <p class="text-sm text-gray-600 truncate">
                                    {swapper.lastMessage}
                                </p>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>

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
                {#each messages as message}
                    <div
                        class="flex flex-col p-2 rounded-lg max-w-2/3 {message.sender ===
                        $auth.user?.name
                            ? 'bg-blue-500 text-white self-end'
                            : 'bg-gray-200 text-gray-800 self-start'}"
                    >
                        {message.text}
                    </div>
                {/each}
            </div>

            <div
                class="col-span-4 row-span-1 bg-white rounded-xl shadow-lg p-3 flex items-center gap-3"
            >
                <button
                    class="p-2 rounded-full bg-gray-200 text-gray-600 hover:bg-gray-300 transition-colors duration-200"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        ><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13.5"
                        /></svg
                    >
                </button>
                <button
                    class="p-2 rounded-full bg-gray-200 text-gray-600 hover:bg-gray-300 transition-colors duration-200"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        ><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        /></svg
                    >
                </button>
                <input
                    type="text"
                    placeholder="Send a message..."
                    class="flex-grow bg-gray-100 text-gray-900 p-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                />
                <button
                    class="p-3 rounded-lg bg-blue-500 text-white font-semibold hover:bg-blue-600 transition-colors duration-200"
                >
                    Send
                </button>
            </div>
        </div> -->
    </div>
</div>
