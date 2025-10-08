<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { Play, Pause, Volume2, Maximize, Settings } from "lucide-svelte";

    let courseId = $page.params.id;
    let moduleId = $page.params.moduleId;
    let loading = $state(true);
    let error = $state<string | null>(null);
    let isPlaying = $state(false);
    let currentTime = $state(0);
    let duration = $state(0);
    let volume = $state(1);

    // Mock data - in a real app, this would come from an API
    let course = $state({
        title: "Sample Course",
        modules: [
            {
                id: "1",
                title: "Introduction",
                description: "Getting started with the course",
                video_url: "#",
                duration: "10:30",
            },
        ],
    });

    let currentModule = $state(course.modules[0]);

    onMount(() => {
        // Simulate loading
        setTimeout(() => {
            loading = false;
        }, 1000);
    });

    function togglePlay() {
        isPlaying = !isPlaying;
    }

    function handleTimeUpdate(event: Event) {
        const video = event.target as HTMLVideoElement;
        currentTime = video.currentTime;
    }

    function handleLoadedMetadata(event: Event) {
        const video = event.target as HTMLVideoElement;
        duration = video.duration;
    }

    function formatTime(seconds: number): string {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = Math.floor(seconds % 60);
        return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
    }
</script>

{#if loading}
    <div class="flex items-center justify-center min-h-screen bg-gray-50">
        <div
            class="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"
        ></div>
    </div>
{:else if error}
    <div class="flex items-center justify-center min-h-screen bg-gray-50">
        <div class="text-center">
            <h1 class="text-2xl font-bold text-gray-900 mb-2">Error</h1>
            <p class="text-gray-600 mb-4">{error}</p>
        </div>
    </div>
{:else}
    <div class="min-h-screen bg-gray-50 p-6">
        <!-- Main Video Container with Red Border -->
        <div class="max-w-6xl mx-auto">
            <div class="border-4 border-red-500 rounded-lg p-4 mb-6 bg-white">
                <!-- Video Player Section -->
                <div class="border-4 border-blue-500 rounded-lg p-4 mb-4">
                    <!-- Video Player -->
                    <div
                        class="relative bg-black rounded-lg overflow-hidden aspect-video"
                    >
                        <!-- Mock video player - replace with actual video element -->
                        <div
                            class="w-full h-full flex items-center justify-center bg-gray-900 text-white"
                        >
                            <div class="text-center">
                                <div class="text-6xl mb-4">🎥</div>
                                <p class="text-lg">Course Video Player</p>
                                <p class="text-sm text-gray-400 mt-2">
                                    Video URL: {currentModule.video_url}
                                </p>
                            </div>
                        </div>

                        <!-- Video Controls -->
                        <div
                            class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4"
                        >
                            <div
                                class="flex items-center justify-between text-white"
                            >
                                <div class="flex items-center space-x-4">
                                    <button
                                        class="hover:text-blue-400 transition-colors"
                                        onclick={togglePlay}
                                    >
                                        {#if isPlaying}
                                            <Pause class="w-6 h-6" />
                                        {:else}
                                            <Play class="w-6 h-6" />
                                        {/if}
                                    </button>

                                    <div class="flex items-center space-x-2">
                                        <Volume2 class="w-4 h-4" />
                                        <input
                                            type="range"
                                            min="0"
                                            max="1"
                                            step="0.1"
                                            bind:value={volume}
                                            class="w-16 accent-blue-500"
                                        />
                                    </div>

                                    <span class="text-sm font-mono">
                                        {formatTime(currentTime)} / {formatTime(
                                            duration,
                                        )}
                                    </span>
                                </div>

                                <div class="flex items-center space-x-2">
                                    <button
                                        class="hover:text-blue-400 transition-colors"
                                    >
                                        <Settings class="w-5 h-5" />
                                    </button>
                                    <button
                                        class="hover:text-blue-400 transition-colors"
                                    >
                                        <Maximize class="w-5 h-5" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Video Info with Green Border -->
                <div class="border-4 border-green-500 rounded-lg p-4">
                    <h1 class="text-2xl font-bold text-gray-900 mb-2">
                        {course.title}
                    </h1>
                    <h2 class="text-xl font-semibold text-gray-700 mb-2">
                        {currentModule.title}
                    </h2>
                    <p class="text-gray-600 mb-4">
                        {currentModule.description}
                    </p>
                    <div
                        class="flex items-center justify-between text-sm text-gray-500"
                    >
                        <span>Duration: {currentModule.duration}</span>
                        <span>Module {moduleId} of {course.modules.length}</span
                        >
                    </div>
                </div>
            </div>

            <!-- Sidebar with Purple Border -->
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div class="lg:col-span-2">
                    <!-- Course Content with Orange Border -->
                    <div
                        class="border-4 border-orange-500 rounded-lg p-4 bg-white"
                    >
                        <h3 class="text-lg font-bold text-gray-900 mb-4">
                            Course Modules
                        </h3>
                        <div class="space-y-3">
                            {#each course.modules as module, index}
                                <div
                                    class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                                >
                                    <div class="flex items-center space-x-3">
                                        <div
                                            class="w-8 h-8 bg-blue-500 text-white rounded-full flex items-center justify-center text-sm font-semibold"
                                        >
                                            {index + 1}
                                        </div>
                                        <div>
                                            <h4
                                                class="font-medium text-gray-900"
                                            >
                                                {module.title}
                                            </h4>
                                            <p class="text-sm text-gray-600">
                                                {module.description}
                                            </p>
                                        </div>
                                    </div>
                                    <div class="text-sm text-gray-500">
                                        {module.duration}
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>

                <!-- Sidebar Info with Purple Border -->
                <div class="lg:col-span-1">
                    <div
                        class="border-4 border-purple-500 rounded-lg p-4 bg-white"
                    >
                        <h3 class="text-lg font-bold text-gray-900 mb-4">
                            What's Next?
                        </h3>
                        <div class="space-y-3">
                            <div class="p-3 bg-blue-50 rounded-lg">
                                <h4 class="font-medium text-blue-900 mb-1">
                                    Practice Exercises
                                </h4>
                                <p class="text-sm text-blue-700">
                                    Complete the hands-on exercises to reinforce
                                    your learning.
                                </p>
                            </div>
                            <div class="p-3 bg-green-50 rounded-lg">
                                <h4 class="font-medium text-green-900 mb-1">
                                    Discussion Forum
                                </h4>
                                <p class="text-sm text-green-700">
                                    Ask questions and share insights with fellow
                                    students.
                                </p>
                            </div>
                            <div class="p-3 bg-yellow-50 rounded-lg">
                                <h4 class="font-medium text-yellow-900 mb-1">
                                    Additional Resources
                                </h4>
                                <p class="text-sm text-yellow-700">
                                    Download supplementary materials and
                                    references.
                                </p>
                            </div>
                        </div>

                        <!-- Note about third-party enhancements -->
                        <div class="mt-6 p-3 bg-gray-50 rounded-lg">
                            <p class="text-xs text-gray-600">
                                <strong>Note:</strong> For enhanced video calling
                                features during live sessions, consider third-party
                                tools like Zoom or Microsoft Teams.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
{/if}
