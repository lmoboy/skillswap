<script lang="ts">
    import { onMount } from "svelte";
    import CourseCard from "$lib/components/CourseCard.svelte";
    import type { Course } from "$lib/types/course";
    import { Search, Filter } from "lucide-svelte";

    let courses = $state<Course[]>([]);
    let loading = $state(true);
    let searchQuery = $state("");
    let filteredCourses = $derived(
        searchQuery
            ? courses.filter(
                  (course) =>
                      course.title
                          .toLowerCase()
                          .includes(searchQuery.toLowerCase()) ||
                      course.description
                          .toLowerCase()
                          .includes(searchQuery.toLowerCase()) ||
                      course.skill_name
                          .toLowerCase()
                          .includes(searchQuery.toLowerCase()),
              )
            : courses,
    );

    onMount(async () => {
        try {
            const response = await fetch("/api/courses", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to fetch courses");
            }

            const data = await response.json();
            courses = data || [];
            loading = false;
        } catch (error) {
            console.error("Error fetching courses:", error);
            loading = false;
        }
    });
</script>

<div class="min-h-screen bg-gray-50">
    <div class=" text-white py-16">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-4">
                Explore Our Courses
            </h1>
            <p class="text-lg text-gray-500 mb-8">
                Learn new skills from expert instructors and grow your career
            </p>

            <div class="max-w-2xl">
                <div class="relative">
                    <div
                        class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"
                    >
                        <Search class="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        bind:value={searchQuery}
                        type="text"
                        placeholder="Search courses by title, skill, or description..."
                        class="block w-full pl-12 pr-4 py-4 border border-transparent rounded-lg text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-white focus:border-transparent"
                    />
                </div>
            </div>
        </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {#if loading}
            <div class="flex items-center justify-center py-20">
                <div
                    class="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"
                ></div>
            </div>
        {:else if filteredCourses.length === 0}
            <div class="text-center py-20">
                <h2 class="text-2xl font-bold text-gray-900 mb-2">
                    No courses found
                </h2>
                <p class="text-gray-600">
                    {searchQuery
                        ? "Try adjusting your search terms"
                        : "Check back later for new courses"}
                </p>
            </div>
        {:else}
            <div class="mb-6">
                <h2 class="text-2xl font-bold text-gray-900">
                    {filteredCourses.length} Course{filteredCourses.length !== 1
                        ? "s"
                        : ""} Available
                </h2>
            </div>

            <div
                class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
            >
                {#each filteredCourses as course (course.id)}
                    <CourseCard {course} />
                {/each}
            </div>
        {/if}
    </div>
</div>
