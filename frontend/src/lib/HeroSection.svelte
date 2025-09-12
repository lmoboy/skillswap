<script>
    import {
        Mic,
        MessageSquare,
        PhoneOff,
        Share2,
        VideoOff,
    } from "lucide-svelte";
</script>

<section class="w-full min-h-screen flex items-center justify-center bg-white">
    <div class="container mx-auto px-6 md:px-12">
        <div
            class="flex flex-col lg:flex-row items-center justify-between gap-12"
        >
            <div class="text-center lg:text-left lg:w-1/2">
                <div
                    class="inline-block px-4 py-2 bg-gray-100 rounded-full mb-6"
                >
                    <span class="text-gray-700 text-sm font-medium"
                        >Knowledge Exchange Platform</span
                    >
                </div>
                <h1
                    class="text-4xl md:text-5xl font-bold text-gray-900 mb-6 leading-tight"
                >
                    Connect & Learn <span class="text-peach-600"
                        >Face to Face</span
                    >
                </h1>
                <p class="text-lg text-gray-600 mb-8 max-w-lg mx-auto lg:mx-0">
                    Share your expertise and learn from others through live
                    video calls. Build meaningful connections while exchanging
                    knowledge in real-time for free*.
                </p>
                <p class="text-sm text-gray-600 mb-8 max-w-lg mx-auto lg:mx-0">
                    *Our platform leverages additional services to enhance the quality of the video calls, which are offered as a subscription service.
                </p>
                <div
                    class="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start"
                >
                    <a
                        href="/signup"
                        class="px-8 py-4 bg-gray-900 hover:bg-gray-800 text-white font-medium rounded-lg transition-all transform hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
                    >
                        Start Your First Call
                    </a>
                    <a
                        href="#how-it-works"
                        class="px-8 py-4 border border-gray-300 hover:border-gray-400 text-gray-700 font-medium rounded-lg transition-all transform hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-300"
                    >
                        How It Works
                    </a>
                </div>
                <div
                    class="mt-10 flex flex-col sm:flex-row items-center justify-center lg:justify-start space-y-4 sm:space-y-0 sm:space-x-8"
                >
                    <div class="flex -space-x-3">
                        <div
                            class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span class="text-gray-600 font-medium text-sm"
                                >JD</span
                            >
                        </div>
                        <div
                            class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span class="text-gray-600 font-medium text-sm"
                                >AM</span
                            >
                        </div>
                        <div
                            class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span class="text-gray-600 font-medium text-sm"
                                >SK</span
                            >
                        </div>
                    </div>
                    <div class="text-sm text-gray-600">
                        <span class="font-medium text-gray-900">10,000+</span> knowledge
                        exchanges completed
                    </div>
                </div>
            </div>

            <div class="lg:w-1/2 relative mt-12 lg:mt-0">
                <div
                    class="relative bg-white rounded-xl shadow-lg overflow-hidden border border-gray-100"
                >
                    <div class="h-3 bg-gray-100 flex items-center px-4">
                        <div
                            class="w-2 h-2 rounded-full bg-gray-300 mr-2"
                        ></div>
                        <div
                            class="w-2 h-2 rounded-full bg-gray-300 mr-2"
                        ></div>
                        <div class="w-2 h-2 rounded-full bg-gray-300"></div>
                    </div>

                    <div
                        class="relative aspect-video bg-gray-50 flex items-center justify-center"
                    >
                        <div class="absolute inset-0">
                            <img
                                src="https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1376&q=80"
                                alt="Person in video call"
                                class="w-full h-full object-cover"
                            />
                        </div>

                        <div
                            class="absolute bottom-4 right-4 w-32 h-24 rounded-lg overflow-hidden border-2 border-white shadow-md"
                        >
                            <div
                                class="w-full h-full bg-gray-200 flex items-center justify-center"
                            >
                                <span class="text-gray-500 text-sm">You</span>
                            </div>
                        </div>

                        <div
                            class="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex items-center space-x-4 bg-white/90 px-4 py-2 rounded-full shadow-lg"
                        >
                            <button
                                class="p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <Mic />
                            </button>
                            <button
                                class="p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <VideoOff />
                            </button>
                            <button
                                class="p-3 bg-peach-500 text-white rounded-full hover:bg-peach-600 transition-colors"
                            >
                                <PhoneOff />
                            </button>
                            <button
                                class="p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <MessageSquare />
                            </button>
                            <button
                                class="p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <Share2 />
                            </button>
                        </div>
                    </div>
                </div>

                <div
                    class="absolute -z-10 w-32 h-32 bg-peach-100 rounded-full -top-8 -right-8 opacity-70"
                ></div>
                <div
                    class="absolute -z-10 w-24 h-24 bg-gray-100 rounded-full -bottom-6 -left-6"
                ></div>
            </div>
        </div>
    </div>
</section>
