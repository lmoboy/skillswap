<script>
    var input = $state("");
    let newMessage = $state("");
    var room = $state("EU");
    var messages = $state([{ Name: "", Message: "", isCurrentUser: false }]);
    var name = $state("");
    import { page } from "$app/state";

    $effect(() => {
        console.log(input);
    });

    var socket = new WebSocket("ws://localhost:8080/chat");

    socket.onopen = () => {
        console.log("Socket is open");
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log(event.data);
        console.log(message.name);
        messages.push({
            Name: message.name,
            Message: message.message,
            isCurrentUser: message.name === name,
        });
        console.log(messages);
    };

    socket.onclose = (event) => {
        console.log("Connection closed:", event.code, event.reason);
        if (event.code === 1001) {
            setTimeout(() => {
                socket = new WebSocket("ws://localhost:8080/chat");
            }, 1000);
        }
    };

    socket.onerror = (event) => {
        console.log("Socket closed: ", event);
    };

    function sendMessage() {
        socket.send(JSON.stringify({ Name: name, Message: newMessage }));
        newMessage = "";
    }
</script>

<div
    class="flex flex-col h-full w-ful mx-auto bg-gray-800 rounded-lg shadow-xl overflow-hidden border border-gray-700"
>
    <div
        class="flex-shrink-0 p-4 bg-gray-900 text-white font-semibold text-lg border-b border-gray-700"
    >
        <input type="text" placeholder="name:" bind:value={name} />
        <input type="text" placeholder="room:" bind:value={room} />
    </div>

    <div class="flex-1 p-4 overflow-y-auto w-full space-y-4">
        {#each messages as message}
            <div class="flex {message.isCurrentUser ? 'justify-end' : ''}">
                {message.Name}: {message.Message}
            </div>
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
