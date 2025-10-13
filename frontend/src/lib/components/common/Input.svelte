<script lang="ts">
    type InputType =
        | "text"
        | "email"
        | "password"
        | "number"
        | "tel"
        | "url"
        | "search";

    type Props = {
        type?: InputType;
        placeholder?: string;
        label?: string;
        error?: string | null;
        disabled?: boolean;
        required?: boolean;
        id?: string;
        name?: string;
        autocomplete?: string;
        class?: string;
        min?: number;
        max?: number;
        maxlength?: number;
        oninput?: (event: Event) => void;
        onchange?: (event: Event) => void;
        onfocus?: (event: FocusEvent) => void;
        onblur?: (event: FocusEvent) => void;
    };

    let {
        type = "text",
        placeholder = "",
        label = "",
        error = null,
        disabled = false,
        required = false,
        id = "",
        name = "",
        autocomplete = "",
        class: className = "",
        min,
        max,
        maxlength,
        oninput,
        onchange,
        onfocus,
        onblur,
    }: Props = $props();

    let value = $state("");

    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;

    const baseClasses =
        "w-full p-3 rounded-lg border bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 transition";
    const normalClasses =
        "border-gray-300 focus:ring-blue-500 focus:border-blue-500";
    const errorClasses =
        "border-red-300 focus:ring-red-500 focus:border-red-500";
    const disabledClasses = "opacity-50 cursor-not-allowed";

    const inputClasses = $derived(
        `
        ${baseClasses}
        ${error ? errorClasses : normalClasses}
        ${disabled ? disabledClasses : ""}
        ${className}
    `
            .trim()
            .replace(/\s+/g, " "),
    );
</script>

<div class="w-full">
    {#if label}
        <label
            for={inputId}
            class="block text-sm font-medium text-gray-700 mb-2"
        >
            {label}
            {#if required}
                <span class="text-red-500">*</span>
            {/if}
        </label>
    {/if}

    <input
        id={inputId}
        {type}
        {name}
        {placeholder}
        {disabled}
        {required}
        autocomplete={autocomplete as any}
        {min}
        {max}
        {maxlength}
        bind:value
        class={inputClasses}
        {oninput}
        {onchange}
        {onfocus}
        {onblur}
    />

    {#if error}
        <p class="mt-1 text-sm text-red-600">
            {error}
        </p>
    {/if}
</div>
