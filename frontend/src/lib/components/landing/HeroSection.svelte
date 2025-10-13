<script>
    import {
        Mic,
        MessageSquare,
        PhoneOff,
        Share2,
        VideoOff,
    } from "lucide-svelte";
    import { auth } from "$lib/stores/auth";
</script>

<section
    class="w-full min-h-screen flex items-center justify-center bg-white py-12 sm:py-16"
>
    <div class="container mx-auto px-4 sm:px-6 md:px-12">
        <div
            class="flex flex-col lg:flex-row items-center justify-between gap-8 lg:gap-12"
        >
            <div class="text-center lg:text-left lg:w-1/2">
                <div
                    class="inline-block px-3 py-1.5 sm:px-4 sm:py-2 bg-gray-100 rounded-full mb-4 sm:mb-6"
                >
                    <span class="text-gray-700 text-xs sm:text-sm font-medium"
                        >Knowledge Exchange Platform</span
                    >
                </div>
                <h1
                    class="text-3xl sm:text-4xl md:text-5xl font-bold text-gray-900 mb-4 sm:mb-6 leading-tight"
                >
                    Connect & Learn <span class="text-peach-600"
                        >Face to Face</span
                    >
                </h1>
                <p
                    class="text-base sm:text-lg text-gray-600 mb-4 sm:mb-6 max-w-lg mx-auto lg:mx-0"
                >
                    Share your expertise and learn from others through live
                    video calls. Build meaningful connections while exchanging
                    knowledge in real-time for free*.
                </p>
                <p
                    class="text-xs sm:text-sm text-gray-600 mb-6 sm:mb-8 max-w-lg mx-auto lg:mx-0"
                >
                    *Our platform leverages additional services to enhance the
                    quality of the video calls, which are offered as a
                    subscription service.
                </p>
                <div
                    class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center lg:justify-start"
                >
                    <a
                        href={$auth.isAuthenticated ? "/swapping" : "/signup"}
                        class="px-6 py-3 sm:px-8 sm:py-4 bg-gray-900 hover:bg-gray-800 text-white text-sm sm:text-base font-medium rounded-lg transition-all transform hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900"
                    >
                        {$auth.isAuthenticated
                            ? "Continue swapping!"
                            : "Start Your First Call!"}
                    </a>
                    <a
                        href="#how-it-works"
                        class="px-6 py-3 sm:px-8 sm:py-4 border border-gray-300 hover:border-gray-400 text-gray-700 text-sm sm:text-base font-medium rounded-lg transition-all transform hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-300"
                    >
                        How It Works
                    </a>
                </div>
                <div
                    class="mt-8 sm:mt-10 flex flex-col sm:flex-row items-center justify-center lg:justify-start space-y-4 sm:space-y-0 sm:space-x-8"
                >
                    <div class="flex -space-x-2 sm:-space-x-3">
                        <div
                            class="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span
                                class="text-gray-600 font-medium text-xs sm:text-sm"
                                >JD</span
                            >
                        </div>
                        <div
                            class="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span
                                class="text-gray-600 font-medium text-xs sm:text-sm"
                                >AM</span
                            >
                        </div>
                        <div
                            class="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-gray-200 flex items-center justify-center border-2 border-white"
                        >
                            <span
                                class="text-gray-600 font-medium text-xs sm:text-sm"
                                >SK</span
                            >
                        </div>
                    </div>
                    <div
                        class="text-xs sm:text-sm text-gray-600 text-center sm:text-left"
                    >
                        <span class="font-medium text-gray-900">100+</span>
                        knowledgeable speakers and over
                        <span class="font-medium text-gray-900">1000+</span> exchanges!
                    </div>
                </div>
            </div>

            <div
                class="lg:w-1/2 relative mt-8 lg:mt-0 w-full max-w-lg lg:max-w-none mx-auto"
            >
                <div
                    class="relative bg-white rounded-lg sm:rounded-xl shadow-lg overflow-hidden border border-gray-100"
                >
                    <div
                        class="h-2 sm:h-3 bg-gray-100 flex items-center px-2 sm:px-4"
                    >
                        <div
                            class="w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full bg-gray-300 mr-1.5 sm:mr-2"
                        ></div>
                        <div
                            class="w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full bg-gray-300 mr-1.5 sm:mr-2"
                        ></div>
                        <div
                            class="w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full bg-gray-300"
                        ></div>
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
                            class="absolute bottom-2 right-2 sm:bottom-4 sm:right-4 w-20 h-16 sm:w-32 sm:h-24 rounded-md sm:rounded-lg overflow-hidden border-2 border-white shadow-md"
                        >
                            <div
                                class="w-full h-full bg-gray-200 flex items-center justify-center"
                            >
                                <span class="text-gray-500 text-xs sm:text-sm"
                                    >You</span
                                >
                            </div>
                        </div>

                        <div
                            class="absolute bottom-2 sm:bottom-4 left-1/2 transform -translate-x-1/2 flex items-center space-x-2 sm:space-x-4 bg-white/90 px-2 py-1.5 sm:px-4 sm:py-2 rounded-full shadow-lg"
                        >
                            <button
                                class="p-1.5 sm:p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <Mic class="w-4 h-4 sm:w-5 sm:h-5" />
                            </button>
                            <button
                                class="p-1.5 sm:p-2 text-gray-600 hover:bg-gray-100 rounded-full"
                            >
                                <VideoOff class="w-4 h-4 sm:w-5 sm:h-5" />
                            </button>
                            <button
                                class="p-2 sm:p-3 bg-peach-500 text-white rounded-full hover:bg-peach-600 transition-colors"
                            >
                                <PhoneOff class="w-4 h-4 sm:w-5 sm:h-5" />
                            </button>
                            <button
                                class="p-1.5 sm:p-2 text-gray-600 hover:bg-gray-100 rounded-full hidden sm:block"
                            >
                                <MessageSquare class="w-4 h-4 sm:w-5 sm:h-5" />
                            </button>
                            <button
                                class="p-1.5 sm:p-2 text-gray-600 hover:bg-gray-100 rounded-full hidden sm:block"
                            >
                                <Share2 class="w-4 h-4 sm:w-5 sm:h-5" />
                            </button>
                        </div>
                    </div>
                </div>

                <div
                    class="absolute -z-10 w-20 h-20 sm:w-32 sm:h-32 bg-peach-100 rounded-full -top-4 -right-4 sm:-top-8 sm:-right-8 opacity-70"
                ></div>
                <div
                    class="absolute -z-10 w-16 h-16 sm:w-24 sm:h-24 bg-gray-100 rounded-full -bottom-3 -left-3 sm:-bottom-6 sm:-left-6"
                ></div>
            </div>
        </div>
    </div>
</section>
