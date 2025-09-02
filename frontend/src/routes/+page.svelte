<script>
    import { onMount } from "svelte";

    var input = $state("");
    var room = $state("EU");
    var messages = $state([{ Name: "", Message: "", isCurrentUser: false }]);
    var name = $state("");

    import Message from "../components/messages/+message.svelte";
    import { json } from "@sveltejs/kit";
    $effect(() => {
        console.log(input);
    });

    var socket = new WebSocket("ws://localhost:8080/chat");

    socket.onopen = () => {
        console.log("Socket is open");
    };

    socket.onmessage = (e) => {
        messages.push(JSON.parse(e.data));
    };

    socket.onclose = () => {
        console.log("Socket is closed");
    };

    // onMount(() => {
    //     var average = $state(0);
    //     setInterval(() => {
    //         var sendTime = Date.now();

    //         fetch("http://localhost:8080/ping").then((response) => {
    //             if (response.ok) {
    //                 average = Date.now() - sendTime;
    //                 ping = Math.floor((average + ping) / 2);
    //             }
    //         });
    //     }, 500);
    // });

    let newMessage = $state("");
    function sendMessage() {
        if (newMessage.trim() !== "") {
            socket.send(JSON.stringify({ Name: name, Message: newMessage }));
            newMessage = "";
        }
    }
</script>

<div>WELCOME TO THE ROOM {room}</div>
<input type="text" name="" id="" bind:value={input} />

<div
    class="flex flex-col h-[70vh] w-3xl max-w-lg mx-auto bg-gray-800 rounded-lg shadow-xl overflow-hidden border border-gray-700"
>
    <div
        class="flex-shrink-0 p-4 bg-gray-900 text-white font-semibold text-lg border-b border-gray-700"
    >
        <input type="text" placeholder="name:" bind:value={name} />
    </div>

    <div class="flex-1 p-4 overflow-y-auto w-full space-y-4">
        {#each messages as message}
            {#if message.Message.trim() !== ""}
                <Message
                    user={message.Name}
                    message={message.Message}
                    isCurrentUser={message.isCurrentUser}
                />
            {/if}
        {/each}
    </div>

    <div class="flex-shrink-0 p-4 bg-gray-900 border-t border-gray-700">
        <form onsubmit={sendMessage} class="flex space-x-2">
            <input
                type="text"
                placeholder="Write a message..."
                class="flex-1 p-3 rounded-full bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                bind:value={newMessage}
            />
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <button
                type="submit"
                class="bg-blue-600 text-white p-3 rounded-full hover:bg-blue-700 transition-colors"
            >
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
                        d="M5 12h14M12 5l7 7-7 7"
                    />
                </svg>
            </button>
        </form>
    </div>
</div>
