<script lang="ts">
    type BadgeVariant =
        | "default"
        | "primary"
        | "success"
        | "warning"
        | "danger"
        | "info";
    type BadgeSize = "sm" | "md" | "lg";

    type Props = {
        variant?: BadgeVariant;
        size?: BadgeSize;
        class?: string;
        children?: any;
    };

    let {
        variant = "default",
        size = "md",
        class: className = "",
        children,
    }: Props = $props();

    const variantClasses: Record<BadgeVariant, string> = {
        default: "bg-gray-100 text-gray-800",
        primary: "bg-blue-100 text-blue-800",
        success: "bg-green-100 text-green-800",
        warning: "bg-orange-100 text-orange-800",
        danger: "bg-red-100 text-red-800",
        info: "bg-purple-100 text-purple-800",
    };

    const sizeClasses: Record<BadgeSize, string> = {
        sm: "px-2 py-0.5 text-xs",
        md: "px-3 py-1 text-sm",
        lg: "px-4 py-1.5 text-base",
    };

    const baseClasses = "inline-flex items-center font-semibold rounded-full";

    const badgeClasses = $derived(
        `
        ${baseClasses}
        ${variantClasses[variant]}
        ${sizeClasses[size]}
        ${className}
    `
            .trim()
            .replace(/\s+/g, " "),
    );
</script>

<span class={badgeClasses}>
    <slot />
</span>
