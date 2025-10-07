<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import {
        Star,
        Clock,
        Users,
        TrendingUp,
        BookOpen,
        Award,
        ChevronRight,
        User,
    } from "lucide-svelte";
    import type { CourseDetail } from "$lib/types/course";

    let course = $state<CourseDetail | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        const courseId = $page.params.id;

        try {
            const response = await fetch(`/api/course?id=${courseId}`, {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch course");
            }

            const data = await response.json();
            course = data;
            loading = false;
        } catch (err) {
            console.error("Error fetching course:", err);
            error = "Failed to load course details";
            loading = false;
        }
    });

    function getDifficultyColor(level: string): string {
        switch (level) {
            case "Beginner":
                return "bg-green-100 text-green-800";
            case "Intermediate":
                return "bg-blue-100 text-blue-800";
            case "Advanced":
                return "bg-orange-100 text-orange-800";
            case "Expert":
                return "bg-red-100 text-red-800";
            default:
                return "bg-gray-100 text-gray-800";
        }
    }

    function formatPrice(price: number): string {
        return price === 0 ? "Free" : `$${price.toFixed(2)}`;
    }

    function formatDate(dateString: string): string {
        if (!dateString) return "";
        const date = new Date(dateString);
        return date.toLocaleDateString("en-US", {
            year: "numeric",
            month: "long",
            day: "numeric",
        });
    }

    function renderStars(rating: number) {
        return Array.from({ length: 5 }, (_, i) => i < Math.floor(rating));
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
            <h1 class="text-2xl font-bold text-gray-900 mb-2">
                Oops! Something went wrong
            </h1>
            <p class="text-gray-600 mb-4">{error}</p>
            <a
                href="/course"
                class="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
                Browse Courses
            </a>
        </div>
    </div>
{:else if course}
    <div class="min-h-screen bg-gray-50">
        <!-- Hero Section -->
        <div class="text-gray-800">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <div class="grid md:grid-cols-3 gap-8">
                    <div class="md:col-span-2">
                        <div class="mb-4">
                            <span
                                class={`px-3 py-1 rounded-full text-sm font-semibold ${getDifficultyColor(course.difficulty_level)}`}
                            >
                                {course.difficulty_level}
                            </span>
                        </div>
                        <h1 class="text-4xl font-bold mb-4">{course.title}</h1>
                        <p class="text-lg text-blue-900 mb-6">
                            {course.description}
                        </p>

                        <div class="flex items-center gap-6 text-sm">
                            <div class="flex items-center gap-2">
                                <User class="w-5 h-5" />
                                <span>Instructor: {course.instructor_name}</span
                                >
                            </div>
                            <div class="flex items-center gap-2">
                                <Star
                                    class="w-5 h-5 fill-yellow-400 text-yellow-400"
                                />
                                <span class="font-medium"
                                    >{course.average_rating > 0
                                        ? course.average_rating.toFixed(1)
                                        : "New"}</span
                                >
                                {#if course.review_count > 0}
                                    <span>({course.review_count} reviews)</span>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <div class="md:col-span-1">
                        <div
                            class="bg-white rounded-xl shadow-lg p-6 text-gray-900"
                        >
                            <div class="mb-4">
                                <img
                                    src={course.thumbnail_url}
                                    alt={course.title}
                                    class="w-full h-40 object-cover rounded-lg"
                                />
                            </div>
                            <div class="text-3xl font-bold mb-4 text-center">
                                {formatPrice(course.price)}
                            </div>
                            <button
                                class="w-full bg-blue-600 text-white py-3 rounded-xl font-semibold hover:bg-blue-700 transition-colors mb-3"
                            >
                                Enroll Now
                            </button>
                            <button
                                class="w-full border-2 border-blue-600 text-blue-600 py-3 rounded-lg font-semibold hover:bg-blue-50 transition-colors"
                            >
                                Add to Wishlist
                            </button>

                            <div class="mt-6 space-y-3 text-sm">
                                <div class="flex items-center justify-between">
                                    <span class="text-gray-600">Duration</span>
                                    <span
                                        class="font-medium flex items-center gap-1"
                                    >
                                        <Clock class="w-4 h-4" />
                                        {course.duration_hours} hours
                                    </span>
                                </div>
                                <div class="flex items-center justify-between">
                                    <span class="text-gray-600">Students</span>
                                    <span
                                        class="font-medium flex items-center gap-1"
                                    >
                                        <Users class="w-4 h-4" />
                                        {course.current_students}/{course.max_students}
                                    </span>
                                </div>
                                <div class="flex items-center justify-between">
                                    <span class="text-gray-600">Skill</span>
                                    <span
                                        class="font-medium flex items-center gap-1"
                                    >
                                        <TrendingUp class="w-4 h-4" />
                                        {course.skill_name}
                                    </span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
            <div class="grid md:grid-cols-3 gap-8">
                <div class="md:col-span-2 space-y-8">
                    <!-- Course Modules -->
                    {#if course.modules && course.modules.length > 0}
                        <div class="bg-white rounded-xl shadow-md p-6">
                            <h2
                                class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2"
                            >
                                <BookOpen class="w-6 h-6 text-blue-600" />
                                Course Curriculum
                            </h2>
                            <div class="space-y-3">
                                {#each course.modules as module, index}
                                    <div
                                        class="border border-gray-200 rounded-lg p-4 hover:border-blue-300 transition-colors"
                                    >
                                        <div
                                            class="flex items-start justify-between"
                                        >
                                            <div class="flex-1">
                                                <div
                                                    class="flex items-center gap-3 mb-2"
                                                >
                                                    <span
                                                        class="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full text-sm font-semibold"
                                                    >
                                                        {index + 1}
                                                    </span>
                                                    <h3
                                                        class="font-semibold text-gray-900"
                                                    >
                                                        {module.title}
                                                    </h3>
                                                </div>
                                                <p
                                                    class="text-sm text-gray-600 ml-11"
                                                >
                                                    {module.description}
                                                </p>
                                            </div>
                                            <ChevronRight
                                                class="w-5 h-5 text-gray-400"
                                            />
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {/if}

                    <!-- Reviews Section -->
                    {#if course.reviews && course.reviews.length > 0}
                        <div class="bg-white rounded-xl shadow-md p-6">
                            <h2
                                class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2"
                            >
                                <Star
                                    class="w-6 h-6 text-yellow-400 fill-yellow-400"
                                />
                                Student Reviews
                            </h2>
                            <div class="space-y-6">
                                {#each course.reviews as review}
                                    <div
                                        class="border-b border-gray-200 pb-6 last:border-0"
                                    >
                                        <div
                                            class="flex items-start justify-between mb-3"
                                        >
                                            <div>
                                                <h4
                                                    class="font-semibold text-gray-900"
                                                >
                                                    {review.student_name}
                                                </h4>
                                                <div
                                                    class="flex items-center gap-2 mt-1"
                                                >
                                                    <div class="flex gap-1">
                                                        {#each renderStars(review.rating) as filled}
                                                            <Star
                                                                class={`w-4 h-4 ${filled ? "fill-yellow-400 text-yellow-400" : "text-gray-300"}`}
                                                            />
                                                        {/each}
                                                    </div>
                                                    <span
                                                        class="text-sm text-gray-500"
                                                    >
                                                        {formatDate(
                                                            review.created_at,
                                                        )}
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                        <p class="text-gray-700">
                                            {review.review_text}
                                        </p>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Sidebar -->
                <div class="md:col-span-1">
                    <div class="bg-white rounded-xl shadow-md p-6 sticky top-6">
                        <h3 class="text-lg font-bold text-gray-900 mb-4">
                            What you'll learn
                        </h3>
                        <ul class="space-y-3">
                            <li class="flex items-start gap-2">
                                <Award
                                    class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5"
                                />
                                <span class="text-sm text-gray-700"
                                    >Master {course.skill_name} fundamentals</span
                                >
                            </li>
                            <li class="flex items-start gap-2">
                                <Award
                                    class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5"
                                />
                                <span class="text-sm text-gray-700"
                                    >Build real-world projects</span
                                >
                            </li>
                            <li class="flex items-start gap-2">
                                <Award
                                    class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5"
                                />
                                <span class="text-sm text-gray-700"
                                    >Get hands-on experience</span
                                >
                            </li>
                            <li class="flex items-start gap-2">
                                <Award
                                    class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5"
                                />
                                <span class="text-sm text-gray-700"
                                    >Receive instructor feedback</span
                                >
                            </li>
                        </ul>

                        <div class="mt-6 pt-6 border-t border-gray-200">
                            <h3 class="text-lg font-bold text-gray-900 mb-4">
                                Instructor
                            </h3>
                            <a
                                href={`/profile/${course.instructor_id}`}
                                class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 transition-colors"
                            >
                                <div
                                    class="w-12 h-12 rounded-full flex items-center justify-center text-white font-bold"
                                >
                                    {course.instructor_name
                                        .charAt(0)
                                        .toUpperCase()}
                                </div>
                                <div>
                                    <p class="font-semibold text-gray-900">
                                        {course.instructor_name}
                                    </p>
                                    <p class="text-sm text-gray-600">
                                        View Profile
                                    </p>
                                </div>
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
{/if}
