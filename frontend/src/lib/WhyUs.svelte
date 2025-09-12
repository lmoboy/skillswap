<script>
    import { RefreshCw, Users, Video, Clock, TrendingUp, BookOpen } from "lucide-svelte";
</script>

<section id="how-it-works" class="w-full py-20 bg-gray-50">
    <div class="container mx-auto px-6">
        <div class="text-center max-w-3xl mx-auto mb-16">
            <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                How <span class="text-gray-600">It Works</span>
            </h2>
            <p class="text-gray-600">
                Our unique fair exchange system ensures a balanced community where everyone contributes and benefits
            </p>
        </div>
        
        <div class="mb-16">
            <h3 class="text-2xl font-bold text-gray-900 mb-6 text-center">Our Fair Exchange System</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
                <div class="bg-white p-8 rounded-xl border border-gray-100 hover:shadow-md transition-shadow duration-300">
                    <div class="w-12 h-12 rounded-lg bg-gray-50 flex items-center justify-center mb-6">
                        <Users class="w-6 h-6 text-gray-600" />
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">Request 3 Swaps</h3>
                    <p class="text-gray-600 text-sm">Find and request up to 3 skill swaps with other members based on your learning goals and expertise.</p>
                </div>
                
                <div class="bg-white p-8 rounded-xl border border-gray-100 hover:shadow-md transition-shadow duration-300">
                    <div class="w-12 h-12 rounded-lg bg-gray-50 flex items-center justify-center mb-6">
                        <RefreshCw class="w-6 h-6 text-gray-600" />
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">Give Back to the Community</h3>
                    <p class="text-gray-600 text-sm">After 3 received swaps, you'll be prompted to give back by teaching a skill to another community member.</p>
                </div>
                
                <div class="bg-white p-8 rounded-xl border border-gray-100 hover:shadow-md transition-shadow duration-300">
                    <div class="w-12 h-12 rounded-lg bg-gray-50 flex items-center justify-center mb-6">
                        <TrendingUp class="w-6 h-6 text-gray-600" />
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">Maintain Balance</h3>
                    <p class="text-gray-600 text-sm">This 3:1 ratio ensures a sustainable ecosystem where everyone contributes and benefits equally.</p>
                </div>
            </div>
        </div>
        
        <div>
            <h3 class="text-2xl font-bold text-gray-900 mb-6 text-center">Pre-recorded Classes</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="bg-white p-8 rounded-xl border border-gray-100 hover:shadow-md transition-shadow duration-300">
                    <div class="w-12 h-12 rounded-lg bg-gray-50 flex items-center justify-center mb-6">
                        <Video class="w-6 h-6 text-gray-600" />
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">Top-Rated Content</h3>
                    <p class="text-gray-600 text-sm mb-4">Access our curated library of pre-recorded classes, ranked and reviewed by our community. Learn at your own pace from the best instructors.</p>
                    <a href="/classes" class="text-sm font-medium text-gray-900 hover:text-gray-600 transition-colors">
                        Browse Classes →
                    </a>
                </div>
                
                <div class="bg-white p-8 rounded-xl border border-gray-100 hover:shadow-md transition-shadow duration-300">
                    <div class="w-12 h-12 rounded-lg bg-gray-50 flex items-center justify-center mb-6">
                        <BookOpen class="w-6 h-6 text-gray-600" />
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900 mb-2">Contribute Your Knowledge</h3>
                    <p class="text-gray-600 text-sm mb-4">Share your expertise by creating pre-recorded classes. Earn recognition and help others learn while building your personal brand.</p>
                    <a href="/create-class" class="text-sm font-medium text-gray-900 hover:text-gray-600 transition-colors">
                        Start Teaching →
                    </a>
                </div>
            </div>
        </div>
        
        <div class="mt-16 grid grid-cols-2 md:grid-cols-4 gap-6">
            <div class="text-center p-6 bg-white rounded-xl border border-gray-100 hover:shadow-md transition-colors">
                <div class="text-2xl font-bold text-gray-600 mb-2">10,000+</div>
                <p class="text-sm text-gray-600">Skills Shared</p>
            </div>
            <div class="text-center p-6 bg-white rounded-xl border border-gray-100 hover:shadow-md transition-colors">
                <div class="text-2xl font-bold text-gray-600 mb-2">5,000+</div>
                <p class="text-sm text-gray-600">Fair Exchanges</p>
            </div>
            <div class="text-center p-6 bg-white rounded-xl border border-gray-100 hover:shadow-md transition-colors">
                <div class="text-2xl font-bold text-gray-600 mb-2">500+</div>
                <p class="text-sm text-gray-600">Pre-recorded Classes</p>
            </div>
            <div class="text-center p-6 bg-white rounded-xl border border-gray-100 hover:shadow-md transition-colors">
                <div class="text-2xl font-bold text-gray-600 mb-2">98%</div>
                <p class="text-sm text-gray-600">Satisfaction Rate</p>
            </div>
        </div>
        
        <div class="mt-16 text-center">
            <h3 class="text-xl font-semibold text-gray-900 mb-6">Ready to start your knowledge exchange journey?</h3>
            <a 
                href="/signup" 
                class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-lg text-white bg-gray-900 hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-900 transition-all transform hover:scale-[1.02]"
            >
                Join Now for Free
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
            </a>
        </div>
    </div>
</section>
