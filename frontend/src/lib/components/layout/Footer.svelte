<script>
    import { Facebook, Linkedin, Mail, MapPin, Twitter } from "lucide-svelte";
</script>

<footer class="bg-white border-t border-gray-200 mt-auto shrink-0">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10 md:py-12">
        <div
            class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8"
        >
            <!-- Company Info -->
            <div class="space-y-3 sm:space-y-4">
                <h3 class="text-base sm:text-lg font-semibold text-gray-900">
                    SkillSwap
                </h3>
                <p class="text-xs sm:text-sm text-gray-600">
                    Connecting knowledge seekers with skilled professionals
                    through a seamless learning experience.
                </p>
                <div class="flex space-x-4">
                    <a
                        href="https://facebook.com"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="text-gray-500 hover:text-gray-700"
                    >
                        <span class="sr-only">Facebook</span>
                        <Facebook class="w-5 h-5" />
                    </a>
                    <a
                        href="https://twitter.com"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="text-gray-500 hover:text-gray-700"
                    >
                        <span class="sr-only">X (formerly twitter)</span>
                        <Twitter class="w-5 h-5" />
                    </a>
                    <a
                        href="https://linkedin.com"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="text-gray-500 hover:text-gray-700"
                    >
                        <span class="sr-only">LinkedIn</span>
                        <Linkedin class="w-5 h-5" />
                    </a>
                </div>
            </div>

            <!-- Quick Links -->
            <div>
                <h3
                    class="text-xs sm:text-sm font-semibold text-gray-900 uppercase tracking-wider"
                >
                    Quick Links
                </h3>
                <ul class="mt-3 sm:mt-4 space-y-2">
                    <li>
                        <a
                            href="/discover"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Discover Skills</a
                        >
                    </li>
                    <li>
                        <a
                            href="/teach"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Become a Teacher</a
                        >
                    </li>
                    <li>
                        <a
                            href="/about"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >About Us</a
                        >
                    </li>
                    <li>
                        <a
                            href="/blog"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Blog</a
                        >
                    </li>
                </ul>
            </div>

            <!-- Support -->
            <div>
                <h3
                    class="text-xs sm:text-sm font-semibold text-gray-900 uppercase tracking-wider"
                >
                    Support
                </h3>
                <ul class="mt-3 sm:mt-4 space-y-2">
                    <li>
                        <a
                            href="/help"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Help Center</a
                        >
                    </li>
                    <li>
                        <a
                            href="/safety"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Safety Information</a
                        >
                    </li>
                    <li>
                        <a
                            href="/privacy"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Privacy Policy</a
                        >
                    </li>
                    <li>
                        <a
                            href="/terms"
                            class="text-xs sm:text-sm text-gray-600 hover:text-gray-900"
                            >Terms of Service</a
                        >
                    </li>
                </ul>
            </div>

            <!-- Contact -->
            <div>
                <h3
                    class="text-xs sm:text-sm font-semibold text-gray-900 uppercase tracking-wider"
                >
                    Contact Us
                </h3>
                <ul class="mt-3 sm:mt-4 space-y-2">
                    <li class="flex items-start">
                        <Mail
                            class="w-4 h-4 sm:w-5 sm:h-5 mr-2 flex-shrink-0"
                        />
                        <span class="text-xs sm:text-sm text-gray-600"
                            >contact@skillswap.com</span
                        >
                    </li>
                    <li class="flex items-start">
                        <MapPin
                            class="w-4 h-4 sm:w-5 sm:h-5 mr-2 flex-shrink-0"
                        />
                        <span class="text-xs sm:text-sm text-gray-600"
                            >123 Learning St.<br />San Francisco, CA 94103</span
                        >
                    </li>
                </ul>
            </div>
        </div>

        <!-- Copyright -->
        <div
            class="mt-8 sm:mt-10 md:mt-12 pt-6 sm:pt-8 border-t border-gray-200"
        >
            <p class="text-xs sm:text-sm text-gray-500 text-center">
                {new Date().getFullYear()} SkillSwap&trade;
            </p>
        </div>
    </div>
</footer>
