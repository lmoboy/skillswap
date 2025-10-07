<script>
    // @ts-nocheck

    import { page } from "$app/stores";
    import { Frown } from "lucide-svelte";
    import { onMount } from "svelte";
    import CourseCard from "$lib/components/CourseCard.svelte";

    let users = $state([]);
    let courses = $state([]);
    let loading = $state(true);

    onMount(async () => {
        const searchQuery = $page.url.searchParams.get("q");

        try {
            // Fetch users
            const userResponse = await fetch(`/api/fullSearch`, {
                method: "POST",
                credentials: "include",
                body: JSON.stringify({
                    query: searchQuery,
                }),
            });
            const userData = await userResponse.json();
            users = userData || [];

            // Fetch courses
            const courseResponse = await fetch(`/api/searchCourses`, {
                method: "POST",
                credentials: "include",
                body: JSON.stringify({
                    query: searchQuery,
                }),
            });
            const courseData = await courseResponse.json();
            courses = courseData || [];
            $inspect(courses);

            loading = false;
            console.log("Users:", users, "Courses:", courses);
        } catch (error) {
            console.error("Search failed:", error);
            loading = false;
        }
    });
</script>

<div class={`flex flex-col h-fit min-h-screen bg-gray-50 text-gray-800`}>
    <div class="px-6 py-4 border-b border-gray-200">
        <h1 class="text-xl font-semibold">Search Results</h1>
    </div>
    {#if loading}
        <p
            class={`w-full h-full text-center align-middle col-span-full flex justify-center items-center       `}
        >
            Loading...
        </p>
    {:else}
        <div class="flex-1 overflow-y-auto p-4 md:p-6">
            <!-- Courses Section -->

            <!-- Users Section -->
            <section class="mb-8">
                <h2 class="text-lg font-medium mb-4">
                    People you might like to swap with
                </h2>
                <div class="flex w-full flex-nowrap overflow-x-auto gap-4">
                    {#if users != null && users.length > 0}
                        {#each users as user}
                            <div
                                class="bg-white rounded-xl shadow-lg flex flex-col w-[320px] md:flex-row overflow-hidden mr-4 flex-shrink-0"
                            >
                                <img
                                    src={user.user.profile_picture}
                                    alt={user.user.username}
                                    class="w-full h-32 md:h-auto md:w-32 object-cover"
                                />
                                <div
                                    class="p-4 flex flex-col justify-between w-full"
                                >
                                    <div>
                                        <h3
                                            class="text-md font-bold text-gray-900"
                                        >
                                            {user.user.username}
                                        </h3>
                                        <p
                                            class="text-sm text-gray-500 line-clamp-2"
                                        >
                                            {user.user.aboutme}
                                        </p>
                                    </div>
                                    <div
                                        class="flex justify-between items-center mt-2 text-sm text-gray-500"
                                    >
                                        <span
                                            >Joined in: {user.user
                                                .created_at}</span
                                        >
                                        <a
                                            href="/swapping/{user.user.id}"
                                            class="bg-amber-300 text-gray-900 rounded-full px-4 py-2 font-medium hover:bg-amber-400 transition-colors"
                                            >Swap</a
                                        >
                                    </div>
                                </div>
                            </div>
                        {/each}
                    {:else}
                        <p
                            class={`w-full h-full text-center align-middle col-span-full flex justify-center items-center flex-col`}
                        >
                            <Frown size={64} />
                            No users found.
                        </p>
                    {/if}
                </div>
            </section>

            <hr class="my-6 border-t border-gray-200" />

            <section>
                <h2 class="text-lg font-medium mb-4">
                    Explore all users if you had something else in mind
                </h2>
                <div class="flex flex-wrap w-full h-full">
                    {#if courses != null && courses.length > 0}
                        <section class="mb-8">
                            <h2 class="text-lg font-medium mb-4">
                                Courses matching your search
                            </h2>
                            <div class="flex flex-wrap gap-6 w-full h-full">
                                {#each courses as course}
                                    <CourseCard {course} />
                                {/each}
                            </div>
                        </section>

                        <hr class="my-6 border-t border-gray-200" />
                    {/if}
                </div>
            </section>
        </div>
    {/if}
</div>
