<script lang="ts">
    import { Star, Clock, Users, TrendingUp } from "lucide-svelte";
    import type { Course } from "$lib/types/course";

    type Props = {
        course: Course;
    };

    let { course }: Props = $props();

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
</script>

<a
    href={`/course/${course.id}`}
    class="block bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden group"
>
    <div class="relative">
        <img
            src={course.thumbnail_url}
            alt={course.title}
            class="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-300"
        />
        <div
            class="absolute top-3 right-3 bg-white px-3 py-1 rounded-full text-sm font-bold text-gray-900"
        >
            {formatPrice(course.price)}
        </div>
        <div class="absolute top-3 left-3">
            <span
                class={`px-3 py-1 rounded-full text-xs font-semibold ${getDifficultyColor(course.difficulty_level)}`}
            >
                {course.difficulty_level}
            </span>
        </div>
    </div>

    <div class="p-5">
        <div class="mb-2">
            <h3
                class="text-lg font-bold text-gray-900 line-clamp-2 group-hover:text-blue-600 transition-colors"
            >
                {course.title}
            </h3>
            <p class="text-sm text-gray-600 mt-1">by {course.instructor_name}</p>
        </div>

        <p class="text-sm text-gray-600 line-clamp-2 mb-4">
            {course.description}
        </p>

        <div class="flex items-center gap-4 text-sm text-gray-500 mb-4">
            <div class="flex items-center gap-1">
                <Star class="w-4 h-4 fill-yellow-400 text-yellow-400" />
                <span class="font-medium text-gray-700"
                    >{course.average_rating > 0
                        ? course.average_rating.toFixed(1)
                        : "New"}</span
                >
                {#if course.review_count > 0}
                    <span class="text-xs">({course.review_count})</span>
                {/if}
            </div>

            <div class="flex items-center gap-1">
                <Clock class="w-4 h-4" />
                <span>{course.duration_hours}h</span>
            </div>

            <div class="flex items-center gap-1">
                <Users class="w-4 h-4" />
                <span>{course.current_students}/{course.max_students}</span>
            </div>
        </div>

        <div class="flex items-center justify-between">
            <span
                class="inline-flex items-center gap-1 text-xs font-medium text-blue-600 bg-blue-50 px-3 py-1 rounded-full"
            >
                <TrendingUp class="w-3 h-3" />
                {course.skill_name}
            </span>

            <button
                class="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
            >
                View Course
            </button>
        </div>
    </div>
</a>
