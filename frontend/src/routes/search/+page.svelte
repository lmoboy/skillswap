<script>
    let items = [
        {
            id: 1,
            name: "John Doe",
            about: "Hello my name is John Doe and I like programming.",
            joined: 2023,
            image: "https://picsum.photos/200?random=1",
        },
        {
            id: 2,
            name: "Jane Smith",
            about: "Passionate about graphic design and visual storytelling.",
            joined: 2022,
            image: "https://picsum.photos/200?random=2",
        },
        {
            id: 3,
            name: "Alex Johnson",
            about: "I enjoy learning new languages and cultures.",
            joined: 2024,
            image: "https://picsum.photos/200?random=3",
        },
        {
            id: 4,
            name: "Emily Davis",
            about: "Expert in public speaking and communication skills.",
            joined: 2021,
            image: "https://picsum.photos/200?random=4",
        },
        {
            id: 5,
            name: "Michael Brown",
            about: "Love all things related to data science and machine learning.",
            joined: 2023,
            image: "https://picsum.photos/200?random=5",
        },
        {
            id: 6,
            name: "Sarah Wilson",
            about: "I am a creative writer and aspiring novelist.",
            joined: 2022,
            image: "https://picsum.photos/200?random=6",
        },
        {
            id: 7,
            name: "Chris Lee",
            about: "Skilled in digital marketing and social media strategy.",
            joined: 2024,
            image: "https://picsum.photos/200?random=7",
        },
        {
            id: 8,
            name: "Olivia Martinez",
            about: "I enjoy playing musical instruments and composing songs.",
            joined: 2021,
            image: "https://picsum.photos/200?random=8",
        },
        {
            id: 9,
            name: "David Taylor",
            about: "Experienced in backend development with Go and Rust.",
            joined: 2023,
            image: "https://picsum.photos/200?random=9",
        },
        {
            id: 10,
            name: "Jessica Garcia",
            about: "My passion is building accessible and inclusive web applications.",
            joined: 2022,
            image: "https://picsum.photos/200?random=10",
        },
    ];
</script>

<div class={`flex flex-col h-screen bg-gray-50 text-gray-800`}>
    <div class="px-6 py-4 border-b border-gray-200">
        <h1 class="text-xl font-semibold">
            People you might like to swap with
        </h1>
    </div>

    <div class="flex-1 overflow-y-auto p-4 md:p-6">
        <section class="mb-8">
            <h2 class="text-lg font-medium mb-4">Best fit for the search!</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {#each items.slice(0, 3) as item}
                    <div
                        class="bg-white rounded-xl shadow-lg flex flex-col md:flex-row overflow-hidden"
                    >
                        <img
                            src={item.image}
                            alt={item.name}
                            class="w-full h-32 md:h-auto md:w-32 object-cover"
                        />
                        <div class="p-4 flex flex-col justify-between w-full">
                            <div>
                                <h3 class="text-md font-bold text-gray-900">
                                    {item.name}
                                </h3>
                                <p class="text-sm text-gray-500 line-clamp-2">
                                    {item.about}
                                </p>
                            </div>
                            <div
                                class="flex justify-between items-center mt-2 text-sm text-gray-500"
                            >
                                <span>Joined in: {item.joined}</span>
                                <button
                                    class="bg-amber-300 text-gray-900 rounded-full px-4 py-2 font-medium hover:bg-amber-400 transition-colors"
                                    >Swap</button
                                >
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </section>

        <hr class="my-6 border-t border-gray-200" />

        <section>
            <h2 class="text-lg font-medium mb-4">
                Explore all users if you had something else in mind
            </h2>
            <div
                class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
            >
                {#each items as item}
                    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
                        <img
                            src={item.image}
                            alt={item.name}
                            class="w-full h-40 object-cover"
                        />
                        <div class="p-4 flex flex-col">
                            <h3 class="text-lg font-semibold text-gray-900">
                                {item.name}
                            </h3>
                            <p class="text-sm text-gray-500 line-clamp-3">
                                {item.about}
                            </p>
                            <div
                                class="flex justify-between items-center mt-4 text-sm text-gray-500"
                            >
                                <span>Joined: {item.joined}</span>
                                <button
                                    class="bg-blue-500 text-white rounded-full px-4 py-2 font-medium hover:bg-blue-600 transition-colors"
                                    >Swap</button
                                >
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </section>
    </div>
</div>
