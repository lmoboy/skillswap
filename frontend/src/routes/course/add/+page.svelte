<script lang="ts">
    import { onMount } from "svelte";
    import { Upload, X, CheckCircle, Plus, Trash2, Video } from "lucide-svelte";
    import { goto } from "$app/navigation";

    type Module = {
        title: string;
        description: string;
        video_file: File | null;
        thumbnail_file: File | null;
        video_duration: number | null;
        order_index: number;
    };

    // Define the form data structure using Svelte 5 $state
    let formData = $state<{
        title: string;
        description: string;
        skill_name: string;
        duration_minutes: number | null;
        preview_photo_file: File | null;
        modules: Module[];
    }>({
        title: "",
        description: "",
        skill_name: "",
        duration_minutes: null,
        preview_photo_file: null,
        modules: [],
    });

    let preview = $state("");
    let loading = $state(false);
    let error = $state<string | null>(null);
    let success = $state(false);
    let availableSkills = $state<
        { id: number; name: string; description: string }[]
    >([]);

    onMount(async () => {
        fetch("/api/getSkills")
            .then((res) => {
                return res.json();
            })
            .then((res) => {
                availableSkills = res;
            });
    });

    // Add a new module
    function addModule() {
        formData.modules.push({
            title: "",
            description: "",
            video_file: null,
            thumbnail_file: null,
            video_duration: null,
            order_index: formData.modules.length,
        });
        formData.modules = formData.modules; // Trigger reactivity
    }

    // Remove a module
    function removeModule(index: number) {
        formData.modules.splice(index, 1);
        // Re-index remaining modules
        formData.modules.forEach((module, idx) => {
            module.order_index = idx;
        });
        formData.modules = formData.modules; // Trigger reactivity
    }

    // Move module up
    function moveModuleUp(index: number) {
        if (index > 0) {
            const temp = formData.modules[index];
            formData.modules[index] = formData.modules[index - 1];
            formData.modules[index - 1] = temp;
            // Re-index
            formData.modules.forEach((module, idx) => {
                module.order_index = idx;
            });
            formData.modules = formData.modules;
        }
    }

    // Move module down
    function moveModuleDown(index: number) {
        if (index < formData.modules.length - 1) {
            const temp = formData.modules[index];
            formData.modules[index] = formData.modules[index + 1];
            formData.modules[index + 1] = temp;
            // Re-index
            formData.modules.forEach((module, idx) => {
                module.order_index = idx;
            });
            formData.modules = formData.modules;
        }
    }

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault();
        loading = true;
        error = null;
        success = false;

        const data = new FormData();

        // Append course basic fields
        data.append("title", formData.title);
        data.append("description", formData.description);
        data.append("skill_name", formData.skill_name);
        if (formData.duration_minutes !== null) {
            data.append("duration_minutes", String(formData.duration_minutes));
        }

        // Append the preview photo file
        if (formData.preview_photo_file) {
            data.append("preview_photo", formData.preview_photo_file);
        }

        // Prepare modules metadata (without files)
        const modulesMetadata = formData.modules.map((module) => ({
            title: module.title,
            description: module.description,
            order_index: module.order_index,
            video_duration: module.video_duration || 0,
        }));

        // Append modules metadata as JSON
        data.append("modules", JSON.stringify(modulesMetadata));

        // Append module video files
        formData.modules.forEach((module, index) => {
            if (module.video_file) {
                data.append(`module_${index}_video`, module.video_file);
            }
            if (module.thumbnail_file) {
                data.append(`module_${index}_thumbnail`, module.thumbnail_file);
            }
        });

        try {
            const response = await fetch("/api/course/add", {
                method: "POST",
                body: data,
                credentials: "include",
            });

            if (!response.ok) {
                const responseData = await response
                    .json()
                    .catch(() => ({ message: "Server error" }));
                throw new Error(
                    responseData.message || "Failed to upload course.",
                );
            }

            // Success handling
            success = true;
            setTimeout(() => {
                goto("/course");
            }, 1500);
        } catch (err) {
            console.error("Course upload error:", err);
            error =
                err instanceof Error
                    ? err.message
                    : "An unexpected error occurred.";
        } finally {
            loading = false;
        }
    }

    // Handler for course preview photo
    function handlePreviewPhotoChange(event: Event) {
        const input = event.target as HTMLInputElement;
        formData.preview_photo_file = input.files ? input.files[0] : null;

        if (formData.preview_photo_file) {
            preview = URL.createObjectURL(formData.preview_photo_file);
        }
    }

    // Handler for module video file
    function handleModuleVideoChange(event: Event, index: number) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            formData.modules[index].video_file = input.files[0];
            formData.modules = formData.modules; // Trigger reactivity
        }
    }

    // Handler for module thumbnail file
    function handleModuleThumbnailChange(event: Event, index: number) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files[0]) {
            formData.modules[index].thumbnail_file = input.files[0];
            formData.modules = formData.modules; // Trigger reactivity
        }
    }

    // Validation for form submission
    const isFormValid = $derived(() => {
        const basicFieldsValid =
            formData.title.trim() !== "" &&
            formData.description.trim() !== "" &&
            formData.skill_name.trim() !== "" &&
            formData.duration_minutes !== null &&
            formData.duration_minutes > 0 &&
            formData.preview_photo_file !== null;

        const modulesValid =
            formData.modules.length > 0 &&
            formData.modules.every(
                (module) =>
                    module.title.trim() !== "" &&
                    module.description.trim() !== "" &&
                    module.video_file !== null &&
                    module.video_duration !== null &&
                    module.video_duration > 0,
            );

        return basicFieldsValid && modulesValid;
    });
</script>

<div class="min-h-screen bg-gray-50 py-12">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="mb-8">
            <h1 class="text-4xl font-bold text-gray-800 mb-2">
                Upload Your Course
            </h1>
            <p class="text-lg text-gray-500">
                Share your expertise by creating a comprehensive course with
                multiple modules and videos.
            </p>
        </div>

        <form
            onsubmit={handleSubmit}
            class="bg-white p-8 shadow-lg rounded-lg space-y-8"
        >
            <!-- Course Information Section -->
            <div class="space-y-4">
                <h2
                    class="text-2xl font-semibold text-gray-900 border-b pb-2 mb-4"
                >
                    Course Information
                </h2>

                <div>
                    <label
                        for="title"
                        class="block text-sm font-medium text-gray-700"
                        >Course Title *</label
                    >
                    <input
                        type="text"
                        id="title"
                        bind:value={formData.title}
                        required
                        placeholder="e.g., Complete Web Development Bootcamp"
                        class="mt-1 block w-full border border-gray-300 rounded-md text-gray-900 shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                    />
                </div>

                <div>
                    <label
                        for="description"
                        class="block text-sm font-medium text-gray-700"
                        >Course Description *</label
                    >
                    <textarea
                        id="description"
                        bind:value={formData.description}
                        rows="4"
                        required
                        placeholder="Describe what students will learn in this course..."
                        class="mt-1 block w-full border border-gray-300 text-gray-900 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                    ></textarea>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label
                            for="skill_name"
                            class="block text-sm font-medium text-gray-700"
                            >Skill Category *</label
                        >
                        <select
                            id="skill_name"
                            bind:value={formData.skill_name}
                            required
                            class="mt-1 block w-full border text-gray-900 border-gray-300 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">Select a skill...</option>
                            {#each availableSkills as skill}
                                <option value={skill.name}>{skill.name}</option>
                            {/each}
                        </select>
                    </div>

                    <div>
                        <label
                            for="duration"
                            class="block text-sm font-medium text-gray-700"
                            >Total Duration (minutes) *</label
                        >
                        <input
                            type="number"
                            id="duration"
                            bind:value={formData.duration_minutes}
                            min="1"
                            required
                            placeholder="e.g., 120"
                            class="mt-1 block w-full border text-gray-900 border-gray-300 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                        />
                    </div>
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2"
                        >Course Thumbnail *</label
                    >
                    <div class="mt-1 flex items-center space-x-4">
                        <label
                            for="preview-photo-upload"
                            class="cursor-pointer bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2"
                        >
                            <Upload class="h-4 w-4" />
                            <span
                                >{formData.preview_photo_file
                                    ? "Change Photo"
                                    : "Upload Photo"}</span
                            >
                            <input
                                id="preview-photo-upload"
                                type="file"
                                accept="image/*"
                                onchange={handlePreviewPhotoChange}
                                required={!formData.preview_photo_file}
                                class="sr-only"
                            />
                        </label>
                        {#if formData.preview_photo_file}
                            <div class="flex-1">
                                <span class="text-sm text-gray-600 block mb-2">
                                    {formData.preview_photo_file.name}
                                </span>
                                {#if preview}
                                    <img
                                        src={preview}
                                        alt="Course preview"
                                        class="h-32 w-auto rounded-lg border border-gray-300"
                                    />
                                {/if}
                            </div>
                        {/if}
                    </div>
                </div>
            </div>

            <hr class="my-8" />

            <!-- Modules Section -->
            <div class="space-y-6">
                <div class="flex justify-between items-center">
                    <div>
                        <h2 class="text-2xl font-semibold text-gray-900">
                            Course Modules
                        </h2>
                        <p class="text-sm text-gray-500 mt-1">
                            Add videos and lessons that make up your course
                        </p>
                    </div>
                    <button
                        type="button"
                        onclick={addModule}
                        class="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2"
                    >
                        <Plus class="h-4 w-4" />
                        <span>Add Module</span>
                    </button>
                </div>

                {#if formData.modules.length === 0}
                    <div
                        class="border-2 border-dashed border-gray-300 rounded-lg p-12 text-center"
                    >
                        <Video class="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <p class="text-gray-500 mb-4">
                            No modules added yet. Start by adding your first
                            module.
                        </p>
                        <button
                            type="button"
                            onclick={addModule}
                            class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-6 rounded transition-all"
                        >
                            Add First Module
                        </button>
                    </div>
                {:else}
                    <div class="space-y-4">
                        {#each formData.modules as module, index (index)}
                            <div
                                class="border border-gray-300 rounded-lg p-6 bg-gray-50"
                            >
                                <div
                                    class="flex justify-between items-start mb-4"
                                >
                                    <h3
                                        class="text-lg font-semibold text-gray-800"
                                    >
                                        Module {index + 1}
                                    </h3>
                                    <div class="flex space-x-2">
                                        {#if index > 0}
                                            <button
                                                type="button"
                                                onclick={() =>
                                                    moveModuleUp(index)}
                                                class="text-gray-500 hover:text-gray-700 p-1"
                                                title="Move up"
                                            >
                                                ↑
                                            </button>
                                        {/if}
                                        {#if index < formData.modules.length - 1}
                                            <button
                                                type="button"
                                                onclick={() =>
                                                    moveModuleDown(index)}
                                                class="text-gray-500 hover:text-gray-700 p-1"
                                                title="Move down"
                                            >
                                                ↓
                                            </button>
                                        {/if}
                                        <button
                                            type="button"
                                            onclick={() => removeModule(index)}
                                            class="text-red-500 hover:text-red-700 p-1"
                                            title="Remove module"
                                        >
                                            <Trash2 class="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>

                                <div class="space-y-4">
                                    <div>
                                        <label
                                            for="module-title-{index}"
                                            class="block text-sm font-medium text-gray-700"
                                            >Module Title *</label
                                        >
                                        <input
                                            type="text"
                                            id="module-title-{index}"
                                            bind:value={module.title}
                                            required
                                            placeholder="e.g., Introduction to HTML"
                                            class="mt-1 block w-full border border-gray-300 rounded-md text-gray-900 shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                                        />
                                    </div>

                                    <div>
                                        <label
                                            for="module-description-{index}"
                                            class="block text-sm font-medium text-gray-700"
                                            >Module Description *</label
                                        >
                                        <textarea
                                            id="module-description-{index}"
                                            bind:value={module.description}
                                            rows="2"
                                            required
                                            placeholder="Describe what this module covers..."
                                            class="mt-1 block w-full border border-gray-300 text-gray-900 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                                        ></textarea>
                                    </div>

                                    <div>
                                        <label
                                            for="module-duration-{index}"
                                            class="block text-sm font-medium text-gray-700"
                                            >Video Duration (seconds) *</label
                                        >
                                        <input
                                            type="number"
                                            id="module-duration-{index}"
                                            bind:value={module.video_duration}
                                            min="1"
                                            required
                                            placeholder="e.g., 600 (10 minutes)"
                                            class="mt-1 block w-full border text-gray-900 border-gray-300 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                                        />
                                    </div>

                                    <div
                                        class="grid grid-cols-1 md:grid-cols-2 gap-4"
                                    >
                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700 mb-2"
                                                >Module Video *</label
                                            >
                                            <label
                                                for="module-video-{index}"
                                                class="cursor-pointer bg-purple-500 hover:bg-purple-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2 w-fit"
                                            >
                                                <Upload class="h-4 w-4" />
                                                <span
                                                    >{module.video_file
                                                        ? "Change Video"
                                                        : "Upload Video"}</span
                                                >
                                                <input
                                                    id="module-video-{index}"
                                                    type="file"
                                                    accept="video/*"
                                                    onchange={(e) =>
                                                        handleModuleVideoChange(
                                                            e,
                                                            index,
                                                        )}
                                                    required={!module.video_file}
                                                    class="sr-only"
                                                />
                                            </label>
                                            {#if module.video_file}
                                                <p
                                                    class="text-xs text-gray-600 mt-2"
                                                >
                                                    {module.video_file.name}
                                                    ({Math.round(
                                                        module.video_file.size /
                                                            1024 /
                                                            1024,
                                                    )} MB)
                                                </p>
                                            {/if}
                                        </div>

                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700 mb-2"
                                                >Module Thumbnail (Optional)</label
                                            >
                                            <label
                                                for="module-thumbnail-{index}"
                                                class="cursor-pointer bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2 w-fit"
                                            >
                                                <Upload class="h-4 w-4" />
                                                <span
                                                    >{module.thumbnail_file
                                                        ? "Change Thumbnail"
                                                        : "Upload Thumbnail"}</span
                                                >
                                                <input
                                                    id="module-thumbnail-{index}"
                                                    type="file"
                                                    accept="image/*"
                                                    onchange={(e) =>
                                                        handleModuleThumbnailChange(
                                                            e,
                                                            index,
                                                        )}
                                                    class="sr-only"
                                                />
                                            </label>
                                            {#if module.thumbnail_file}
                                                <p
                                                    class="text-xs text-gray-600 mt-2"
                                                >
                                                    {module.thumbnail_file.name}
                                                </p>
                                            {/if}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>

            <hr class="my-8" />

            <!-- Submit Section -->
            <div>
                {#if error}
                    <div
                        class="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded-lg flex items-center space-x-2"
                        role="alert"
                    >
                        <X class="h-4 w-4" />
                        <span><strong>Error:</strong> {error}</span>
                    </div>
                {/if}

                {#if success}
                    <div
                        class="p-3 mb-4 text-sm text-green-700 bg-green-100 rounded-lg flex items-center space-x-2"
                        role="alert"
                    >
                        <CheckCircle class="h-4 w-4" />
                        <span
                            ><strong>Success!</strong> Your course has been uploaded.
                            Redirecting...</span
                        >
                    </div>
                {/if}

                <button
                    type="submit"
                    disabled={loading || !isFormValid() || success}
                    class="w-full py-3 px-4 border border-transparent rounded-md shadow-sm text-lg font-medium text-white transition-colors"
                    class:bg-blue-500={isFormValid() && !loading && !success}
                    class:hover:bg-blue-600={isFormValid() &&
                        !loading &&
                        !success}
                    class:bg-gray-400={!isFormValid() || loading || success}
                    class:cursor-not-allowed={!isFormValid() ||
                        loading ||
                        success}
                >
                    {#if loading}
                        <div class="flex items-center justify-center">
                            <div
                                class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
                            ></div>
                            Uploading Course...
                        </div>
                    {:else}
                        Submit Course
                    {/if}
                </button>

                {#if !isFormValid()}
                    <p class="text-sm text-gray-500 text-center mt-2">
                        Please fill in all required fields and add at least one
                        module with a video.
                    </p>
                {/if}
            </div>
        </form>
    </div>
</div>
