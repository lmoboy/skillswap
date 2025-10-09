<script lang="ts">
    import { onMount } from "svelte";
    import type { Course } from "$lib/types/course"; // Assuming you have this type defined
    import { Upload, X, CheckCircle } from "lucide-svelte";
    import { goto } from "$app/navigation";

    // Define the form data structure using Svelte 5 $state
    let formData = $state<{
        title: string;
        description: string;
        skill_name: string;
        duration_minutes: number | null; // Changed to number for duration
        preview_photo_file: File | null;
        files: File[]; // For multiple file uploads
    }>({
        title: "",
        description: "",
        skill_name: "",
        duration_minutes: null,
        preview_photo_file: null,
        files: [],
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

    // Form submission handler
    async function handleSubmit(event: Event) {
        event.preventDefault();
        loading = true;
        error = null;
        success = false;

        const data = new FormData();

        // Append text fields
        data.append("title", formData.title);
        data.append("description", formData.description);
        data.append("skill_name", formData.skill_name);
        // Ensure duration is appended as a string if it's not null
        if (formData.duration_minutes !== null) {
            data.append("duration_minutes", String(formData.duration_minutes));
        }

        // Append the preview photo file
        if (formData.preview_photo_file) {
            data.append("preview_photo", formData.preview_photo_file);
        }

        // Append the multiple course files
        formData.files.forEach((file, index) => {
            data.append(`course_files`, file, file.name); // Backend expects 'course_files' for each file
        });

        try {
            // Replace with your actual backend endpoint for course submission
            const response = await fetch("/api/course/add", {
                method: "POST",
                body: data,
                credentials: "include",
            });

            if (!response.ok) {
                // Attempt to read a more specific error message from the backend
                const responseData = await response
                    .json()
                    .catch(() => ({ message: "Server error" }));
                throw new Error(
                    responseData.message || "Failed to upload course.",
                );
            }

            // Success handling
            success = true;
            // Optionally redirect after a brief delay
            setTimeout(() => {
                goto("/"); // Redirect to the course list or dashboard
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

    // Handlers for file inputs

    function handlePreviewPhotoChange(event: Event) {
        const input = event.target as HTMLInputElement;

        formData.preview_photo_file = input.files ? input.files[0] : null;

        if (formData.preview_photo_file) {
            preview = URL.createObjectURL(formData.preview_photo_file);
        }
    }

    function handleFilesChange(event: Event) {
        const input = event.target as HTMLInputElement;
        // Convert FileList to Array and update the files state
        formData.files = Array.from(input.files || []);
    }

    function removeFile(index: number) {
        formData.files.splice(index, 1);
        // Important: Svelte 5 state requires reassigning the array to trigger reactivity
        formData.files = formData.files;
    }

    // Validation for form submission (basic example)
    const isFormValid = $derived(() => {
        return (
            formData.title.trim() !== "" &&
            formData.description.trim() !== "" &&
            formData.skill_name.trim() !== "" &&
            formData.duration_minutes !== null &&
            formData.duration_minutes > 0 &&
            formData.preview_photo_file !== null &&
            formData.files.length > 0
        );
    });
</script>

<div class="min-h-screen bg-gray-50 py-12">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <h1 class="text-4xl font-bold text-gray-800 mb-2">
            Upload Your Course
        </h1>
        <p class="text-lg text-gray-500 mb-8">
            Share your expertise and upload course materials for others to
            learn.
        </p>

        <form
            onsubmit={handleSubmit}
            class="bg-white p-8 shadow-lg rounded-lg space-y-6"
        >
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
                        >Course Title</label
                    >
                    <input
                        type="text"
                        id="title"
                        bind:value={formData.title}
                        required
                        class="mt-1 block w-full border border-gray-300 rounded-md text-gray-500 shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                    />
                </div>

                <div>
                    <label
                        for="description"
                        class="block text-sm font-medium text-gray-700"
                        >Description</label
                    >
                    <textarea
                        id="description"
                        bind:value={formData.description}
                        rows="4"
                        required
                        class="mt-1 block w-full border border-gray-300 text-gray-500 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                    ></textarea>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label
                            for="skill_name"
                            class="block text-sm font-medium text-gray-700"
                            >Skill Taught</label
                        >
                        <select
                            id="skill_name"
                            bind:value={formData.skill_name}
                            required
                            class="mt-1 block w-full border text-gray-500 border-gray-300 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                        >
                            {#each availableSkills as skill}
                                <option value={skill.name}>{skill.name}</option>
                            {/each}
                        </select>
                    </div>

                    <div>
                        <label
                            for="duration"
                            class="block text-sm font-medium text-gray-700"
                            >Duration (in minutes)</label
                        >
                        <input
                            type="number"
                            id="duration"
                            bind:value={formData.duration_minutes}
                            min="1"
                            required
                            class="mt-1 block w-full border text-gray-500 border-gray-300 rounded-md shadow-sm p-3 focus:ring-blue-500 focus:border-blue-500"
                        />
                    </div>
                </div>
            </div>

            ---

            <div class="space-y-4">
                <h2
                    class="text-2xl font-semibold text-gray-900 border-b pb-2 mb-4"
                >
                    Course Media & Files
                </h2>

                <div>
                    <label class="block text-sm font-medium text-gray-700"
                        >Preview Photo (Thumbnail)</label
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
                            <span
                                class="text-sm text-gray-600 truncate max-w-xs"
                            >
                                **Selected:** {formData.preview_photo_file.name}
                                <img src={preview} />
                            </span>
                        {/if}
                    </div>
                </div>

                <div>
                    <label class="block text-sm font-medium text-gray-700"
                        >Course Files (Videos, PDFs, etc.)</label
                    >
                    <div class="mt-1">
                        <label
                            for="course-files-upload"
                            class="cursor-pointer bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded transition-all flex items-center space-x-2 w-fit"
                        >
                            <Upload class="h-4 w-4" />
                            <span>Select Course Files</span>
                            <input
                                id="course-files-upload"
                                type="file"
                                multiple
                                onchange={handleFilesChange}
                                required={formData.files.length === 0}
                                class="sr-only"
                            />
                        </label>

                        {#if formData.files.length > 0}
                            <div
                                class="mt-4 space-y-2 max-h-40 overflow-y-auto p-2 border rounded-md"
                            >
                                {#each formData.files as file, index (file.name + index)}
                                    <div
                                        class="flex justify-between items-center text-sm bg-gray-100 p-2 rounded"
                                    >
                                        <span
                                            class="truncate max-w-[calc(100%-40px)] text-gray-800"
                                            >{file.name} ({Math.round(
                                                file.size / 1024 / 1024,
                                            )} MB)</span
                                        >
                                        <button
                                            type="button"
                                            onclick={() => removeFile(index)}
                                            class="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-100"
                                            aria-label="Remove file"
                                        >
                                            <X class="h-4 w-4" />
                                        </button>
                                    </div>
                                {/each}
                            </div>
                        {:else}
                            <p class="mt-2 text-sm text-gray-500">
                                No files selected. Please upload your course
                                content.
                            </p>
                        {/if}
                    </div>
                </div>
            </div>

            ---

            <div>
                {#if error}
                    <div
                        class="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded-lg flex items-center space-x-2"
                        role="alert"
                    >
                        <X class="h-4 w-4" />
                        <span>**Error:** {error}</span>
                    </div>
                {/if}

                {#if success}
                    <div
                        class="p-3 mb-4 text-sm text-green-700 bg-green-100 rounded-lg flex items-center space-x-2"
                        role="alert"
                    >
                        <CheckCircle class="h-4 w-4" />
                        <span
                            >**Success!** Your course has been uploaded.
                            Redirecting...</span
                        >
                    </div>
                {/if}

                <button
                    type="submit"
                    disabled={loading || !isFormValid || success}
                    class="w-full py-3 px-4 border border-transparent rounded-md shadow-sm text-lg font-medium text-white transition-colors"
                    class:bg-blue-500={isFormValid() && !loading && !success}
                    class:hover:bg-blue-600={isFormValid() &&
                        !loading &&
                        !success}
                    class:bg-gray-400={!isFormValid || loading || success}
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
            </div>
        </form>
    </div>
</div>
