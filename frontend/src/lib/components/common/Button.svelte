<script lang="ts">
    type ButtonVariant =
        | "primary"
        | "secondary"
        | "danger"
        | "ghost"
        | "outline";
    type ButtonSize = "sm" | "md" | "lg";

    type Props = {
        variant?: ButtonVariant;
        size?: ButtonSize;
        disabled?: boolean;
        loading?: boolean;
        type?: "button" | "submit" | "reset";
        fullWidth?: boolean;
        class?: string;
        onclick?: (event: MouseEvent) => void;
        children?: any;
        "data-testid"?: string;
    };

    let {
        variant = "primary",
        size = "md",
        disabled = false,
        loading = false,
        type = "button",
        fullWidth = false,
        class: className = "",
        onclick,
        children,
        "data-testid": dataTestId,
    }: Props = $props();

    const variantClasses: Record<ButtonVariant, string> = {
        primary: "bg-blue-600 hover:bg-blue-700 text-white border-transparent",
        secondary:
            "bg-gray-200 hover:bg-gray-300 text-gray-800 border-transparent",
        danger: "bg-red-600 hover:bg-red-700 text-white border-transparent",
        ghost: "bg-transparent hover:bg-gray-100 text-gray-700 border-transparent",
        outline:
            "bg-transparent hover:bg-gray-50 text-gray-700 border-gray-300",
    };

    const sizeClasses: Record<ButtonSize, string> = {
        sm: "px-3 py-1.5 text-sm",
        md: "px-4 py-2 text-base",
        lg: "px-6 py-3 text-lg",
    };

    const baseClasses =
        "font-medium rounded-lg border transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500";
    const disabledClasses = "opacity-50 cursor-not-allowed";
    const fullWidthClass = fullWidth ? "w-full" : "";

    const buttonClasses = $derived(
        `
        ${baseClasses}
        ${variantClasses[variant]}
        ${sizeClasses[size]}
        ${fullWidthClass}
        ${disabled || loading ? disabledClasses : ""}
        ${className}
    `
            .trim()
            .replace(/\s+/g, " "),
    );
</script>

<button {type} class={buttonClasses} disabled={disabled || loading} {onclick} data-testid={dataTestId}>
    {#if loading}
        <div class="flex items-center justify-center gap-2">
            <div
                class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"
            ></div>
            <span>Loading...</span>
        </div>
    {:else}
        {@render children?.()}
    {/if}
</button>
