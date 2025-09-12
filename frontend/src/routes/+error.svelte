<script>
    // You can pass a status code and message as props if needed
    export let status = 404;
    export let message = "Page Not Found";

    // Random learning quotes to show on the error page
    const learningQuotes = [
        "The beautiful thing about learning is that no one can take it away from you. - B.B. King",
        "Learning never exhausts the mind. - Leonardo da Vinci",
        "The more that you read, the more things you will know. The more that you learn, the more places you'll go. - Dr. Seuss",
        "Live as if you were to die tomorrow. Learn as if you were to live forever. - Mahatma Gandhi",
        "The capacity to learn is a gift; the ability to learn is a skill; the willingness to learn is a choice. - Brian Herbert",
    ];

    // Get a random quote
    const randomQuote =
        learningQuotes[Math.floor(Math.random() * learningQuotes.length)];
</script>

<div
    class="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4 sm:px-6 lg:px-8"
>
    <div class="max-w-md w-full space-y-8 text-center">
        <!-- Error Code -->
        <div class="text-9xl font-bold text-gray-300">
            {status}
        </div>

        <!-- Error Message -->
        <h1 class="text-3xl font-bold text-gray-900">
            {message}
        </h1>

        <!-- Random Learning Quote -->
        <div class="mt-6 p-4 bg-white rounded-lg border border-gray-100">
            <p class="text-gray-600 italic">"{randomQuote.split(" - ")[0]}"</p>
            {#if randomQuote.includes(" - ")}
                <p class="text-gray-500 text-sm mt-2">
                    — {randomQuote.split(" - ")[1]}
                </p>
            {/if}
        </div>

        <!-- Action Buttons -->
        <div class="mt-8 flex flex-col sm:flex-row justify-center gap-4">
            <a
                href="/"
                class="w-full sm:w-auto px-6 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors"
            >
                Go back home
            </a>
            <a
                href="/browse"
                class="w-full sm:w-auto px-6 py-3 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors"
            >
                Browse skills
            </a>
        </div>
    </div>
</div>
