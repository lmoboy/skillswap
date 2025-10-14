<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { Star, Clock, Users, BookOpen, Play } from "lucide-svelte";
    import type { CourseDetail } from "$lib/types/course";

    let course = $state<CourseDetail | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);
    let currentModule = $state(0);

    onMount(async () => {
        const courseId = $page.params.id;
        try {
            const response = await fetch(`/api/course?id=${courseId}`);
            if (!response.ok) throw new Error("Failed to fetch course");
            course = await response.json();
        } catch (err) {
            error = "Failed to load course";
            console.error("Error loading course:", err);
        } finally {
            loading = false;
        }
    });

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, "0")}`;
    }
</script>

{#if loading}
    <div class="flex items-center justify-center min-h-screen">
        <div
            class="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"
        ></div>
    </div>
{:else if error}
    <div class="flex items-center justify-center min-h-screen">
        <p class="text-red-600">{error}</p>
    </div>
{:else if course}
    <div class="min-h-screen bg-gray-100">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 py-4 sm:py-6 md:py-8">
            <h1
                class="text-xl sm:text-2xl md:text-3xl text-gray-700 font-bold mb-2"
            >
                {course.title}
            </h1>
            <p class="text-sm sm:text-base text-gray-600 mb-4 sm:mb-6">
                {course.description}
            </p>

            <div class="grid lg:grid-cols-3 gap-4 sm:gap-6">
                <div class="lg:col-span-2">
                    {#if course.modules && course.modules.length > 0 && course.modules[currentModule].video_url}
                        <div
                            class="bg-black rounded-lg overflow-hidden mb-3 sm:mb-4"
                        >
                            <video
                                controls
                                class="w-full aspect-video"
                                src="/api/course/video?path={course.modules[
                                    currentModule
                                ].video_url}"
                                poster={course.modules[currentModule]
                                    .thumbnail_url || course.thumbnail_url}
                            >
                                <track kind="captions" />
                            </video>
                        </div>
                        <div class="bg-white rounded-lg p-4 sm:p-6 shadow">
                            <h2
                                class="text-lg sm:text-xl md:text-2xl text-gray-700 font-bold mb-2"
                            >
                                {course.modules[currentModule].title}
                            </h2>
                            <p class="text-sm sm:text-base text-gray-700">
                                {course.modules[currentModule].description}
                            </p>
                        </div>
                    {:else}
                        <div
                            class="bg-gray-200 rounded-lg p-8 sm:p-12 text-center"
                        >
                            <p class="text-sm sm:text-base text-gray-600">
                                No video available
                            </p>
                        </div>
                    {/if}
                </div>

                <div class="lg:col-span-1">
                    <div
                        class="bg-white rounded-lg shadow p-4 sm:p-6 mb-4 sm:mb-6"
                    >
                        <div class="flex items-center gap-2 mb-3 sm:mb-4">
                            <Star
                                class="w-4 h-4 sm:w-5 sm:h-5 text-yellow-500 fill-current"
                            />
                            <span
                                class="text-sm sm:text-base text-gray-700 font-bold"
                                >{course.average_rating.toFixed(1)}</span
                            >
                            <span class="text-xs sm:text-sm text-gray-600"
                                >({course.review_count} reviews)</span
                            >
                        </div>
                        <div
                            class="flex items-center gap-2 text-gray-600 mb-2 text-sm sm:text-base"
                        >
                            <Clock class="w-3 h-3 sm:w-4 sm:h-4" />
                            <span>{course.duration_hours} hours</span>
                        </div>
                        <div
                            class="flex items-center gap-2 text-gray-600 text-sm sm:text-base"
                        >
                            <Users class="w-3 h-3 sm:w-4 sm:h-4" />
                            <span
                                >{course.current_students}/{course.max_students}
                                students</span
                            >
                        </div>
                    </div>

                    {#if course.modules && course.modules.length > 0}
                        <div class="bg-white rounded-lg shadow">
                            <div class="p-3 sm:p-4 border-b">
                                <div class="flex items-center gap-2">
                                    <BookOpen class="w-4 h-4 sm:w-5 sm:h-5" />
                                    <h3
                                        class="text-sm sm:text-base text-gray-700 font-bold"
                                    >
                                        Course Content
                                    </h3>
                                </div>
                            </div>
                            <div class="max-h-72 sm:max-h-96 overflow-y-auto">
                                {#each course.modules as module, index}
                                    <button
                                        onclick={() => (currentModule = index)}
                                        class="w-full text-left p-3 sm:p-4 border-b hover:bg-gray-50 transition-colors {currentModule ===
                                        index
                                            ? 'bg-blue-50'
                                            : ''}"
                                    >
                                        <div
                                            class="flex items-start gap-2 sm:gap-3"
                                        >
                                            <div
                                                class="flex-shrink-0 w-7 h-7 sm:w-8 sm:h-8 bg-gray-200 rounded flex items-center justify-center"
                                            >
                                                {#if currentModule === index}
                                                    <Play
                                                        class="w-3 h-3 sm:w-4 sm:h-4 text-blue-600"
                                                    />
                                                {:else}
                                                    <span
                                                        class="text-xs sm:text-sm font-medium"
                                                        >{index + 1}</span
                                                    >
                                                {/if}
                                            </div>
                                            <div class="flex-1 min-w-0">
                                                <p
                                                    class="font-medium text-xs sm:text-sm {currentModule ===
                                                    index
                                                        ? 'text-blue-600'
                                                        : 'text-gray-900'}"
                                                >
                                                    {module.title}
                                                </p>
                                                {#if module.video_duration > 0}
                                                    <p
                                                        class="text-xs text-gray-500 mt-1"
                                                    >
                                                        {formatDuration(
                                                            module.video_duration,
                                                        )}
                                                    </p>
                                                {/if}
                                            </div>
                                        </div>
                                    </button>
                                {/each}
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
{/if}
