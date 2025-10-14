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
        const file = input.files ? input.files[0] : null;
        
        if (file) {
            // Validate file type
            const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif'];
            if (!validTypes.includes(file.type)) {
                error = 'Invalid file type. Please upload PNG, JPG, or GIF image.';
                input.value = '';
                return;
            }
            
            // Validate file size (max 5MB)
            if (file.size > 5 * 1024 * 1024) {
                error = 'File size too large. Maximum size is 5MB.';
                input.value = '';
                return;
            }
            
            formData.preview_photo_file = file;
            preview = URL.createObjectURL(file);
            error = null;
        }
    }

    // Handler for module video file
    function handleModuleVideoChange(event: Event, index: number) {
        const input = event.target as HTMLInputElement;
        const file = input.files ? input.files[0] : null;
        
        if (file) {
            // Validate file type
            const validTypes = ['video/mp4', 'video/webm', 'video/avi', 'video/quicktime'];
            if (!validTypes.includes(file.type)) {
                error = `Invalid video file type for module ${index + 1}. Please upload MP4, WebM, or AVI video.`;
                input.value = '';
                return;
            }
            
            // Validate file size (max 500MB)
            if (file.size > 500 * 1024 * 1024) {
                error = `Video file for module ${index + 1} is too large. Maximum size is 500MB.`;
                input.value = '';
                return;
            }
            
            formData.modules[index].video_file = file;
            formData.modules = formData.modules; // Trigger reactivity
            error = null;
        }
    }

    // Handler for module thumbnail file
    function handleModuleThumbnailChange(event: Event, index: number) {
        const input = event.target as HTMLInputElement;
        const file = input.files ? input.files[0] : null;
        
        if (file) {
            // Validate file type
            const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif'];
            if (!validTypes.includes(file.type)) {
                error = `Invalid thumbnail type for module ${index + 1}. Please upload PNG, JPG, or GIF image.`;
                input.value = '';
                return;
            }
            
            // Validate file size (max 2MB)
            if (file.size > 2 * 1024 * 1024) {
                error = `Thumbnail for module ${index + 1} is too large. Maximum size is 2MB.`;
                input.value = '';
                return;
            }
            
            formData.modules[index].thumbnail_file = file;
            formData.modules = formData.modules; // Trigger reactivity
            error = null;
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

<div class="min-h-screen bg-gray-50 py-6 sm:py-12">
    <div class="max-w-6xl mx-auto px-3 sm:px-6 lg:px-8">
        <div class="mb-6 sm:mb-8">
            <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-2">
                Upload Your Course
            </h1>
            <p class="text-base sm:text-lg text-gray-500">
                Share your expertise by creating a comprehensive course with
                multiple modules and videos.
            </p>
        </div>

        <form
            onsubmit={handleSubmit}
            class="bg-white p-4 sm:p-6 lg:p-8 shadow-lg rounded-lg space-y-6 sm:space-y-8"
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
                    <div class="mt-1 space-y-3">
                        <label
                            for="preview-photo-upload"
                            class="cursor-pointer bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center justify-center sm:justify-start space-x-2 w-full sm:w-fit"
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
                            <div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
                                <div class="flex flex-col sm:flex-row gap-4">
                                    <div class="flex-1">
                                        <p class="text-sm font-medium text-gray-700 mb-1">
                                            {formData.preview_photo_file.name}
                                        </p>
                                        <p class="text-xs text-gray-500">
                                            {Math.round(formData.preview_photo_file.size / 1024)} KB
                                        </p>
                                    </div>
                                    {#if preview}
                                        <img
                                            src={preview}
                                            alt="Course preview"
                                            class="h-24 sm:h-32 w-auto rounded-lg border border-gray-300 mx-auto sm:mx-0"
                                        />
                                    {/if}
                                </div>
                            </div>
                        {:else}
                            <p class="text-xs text-gray-500">
                                Recommended size: 1200x630px. Accepted formats: PNG, JPG, GIF
                            </p>
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
                                            <div class="space-y-2">
                                                <label
                                                    for="module-video-{index}"
                                                    class="cursor-pointer bg-purple-500 hover:bg-purple-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2 w-full sm:w-fit justify-center sm:justify-start"
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
                                                    <div class="bg-green-50 border border-green-200 rounded-lg p-3">
                                                        <div class="flex items-start gap-2">
                                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                            </svg>
                                                            <div class="flex-1 min-w-0">
                                                                <p class="text-sm font-medium text-green-900 truncate">
                                                                    {module.video_file.name}
                                                                </p>
                                                                <p class="text-xs text-green-700 mt-1">
                                                                    {Math.round(module.video_file.size / 1024 / 1024)} MB
                                                                </p>
                                                            </div>
                                                        </div>
                                                    </div>
                                                {:else}
                                                    <p class="text-xs text-gray-500">
                                                        Accepted formats: MP4, WebM, AVI
                                                    </p>
                                                {/if}
                                            </div>
                                        </div>

                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700 mb-2"
                                                >Module Thumbnail (Optional)</label
                                            >
                                            <div class="space-y-2">
                                                <label
                                                    for="module-thumbnail-{index}"
                                                    class="cursor-pointer bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2 w-full sm:w-fit justify-center sm:justify-start"
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
                                                    <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
                                                        <div class="flex items-start gap-2">
                                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                                            </svg>
                                                            <div class="flex-1 min-w-0">
                                                                <p class="text-sm font-medium text-blue-900 truncate">
                                                                    {module.thumbnail_file.name}
                                                                </p>
                                                            </div>
                                                        </div>
                                                    </div>
                                                {:else}
                                                    <p class="text-xs text-gray-500">
                                                        Accepted formats: PNG, JPG, GIF
                                                    </p>
                                                {/if}
                                            </div>
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
