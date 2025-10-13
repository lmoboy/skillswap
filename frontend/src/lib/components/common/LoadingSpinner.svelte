<script lang="ts">
    type Props = {
        size?: 'sm' | 'md' | 'lg' | 'xl';
        color?: string;
        class?: string;
        text?: string;
    };

    let {
        size = 'md',
        color = 'text-blue-500',
        class: className = '',
        text = ''
    }: Props = $props();

    const sizeClasses = {
        sm: 'w-4 h-4 border-2',
        md: 'w-8 h-8 border-2',
        lg: 'w-12 h-12 border-3',
        xl: 'w-16 h-16 border-4'
    };

    const spinnerClasses = $derived(`
        ${sizeClasses[size]}
        border-current
        border-t-transparent
        rounded-full
        animate-spin
        ${color}
        ${className}
    `.trim().replace(/\s+/g, ' '));
</script>

<div class="flex flex-col items-center justify-center gap-3">
    <div class={spinnerClasses}></div>
    {#if text}
        <p class="text-sm text-gray-600">{text}</p>
    {/if}
</div>
