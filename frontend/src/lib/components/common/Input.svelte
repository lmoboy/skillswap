<script lang="ts">
    type InputType = "text" | "email" | "password" | "number" | "tel" | "url" | "search";

    interface Props {
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
    }

    const {
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

    // Generate stable ID to avoid re-renders
    const inputId = $derived(id || `input-${crypto.randomUUID()}`);

    // Optimized class computation
    const inputClasses = $derived(() => {
        const base = "w-full p-3 rounded-lg border bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 transition";
        const stateClasses = error 
            ? "border-red-300 focus:ring-red-500 focus:border-red-500"
            : "border-gray-300 focus:ring-blue-500 focus:border-blue-500";
        const disabledClass = disabled ? "opacity-50 cursor-not-allowed" : "";
        
        return [base, stateClasses, disabledClass, className]
            .filter(Boolean)
            .join(" ");
    });
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
