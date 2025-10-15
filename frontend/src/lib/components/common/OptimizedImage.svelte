  <script lang="ts">
    interface Props {
        src: string;
        alt: string;
        width?: number;
        height?: number;
        class?: string;
        loading?: "lazy" | "eager";
        priority?: boolean;
        sizes?: string;
        quality?: number;
    }

    const {
        src,
        alt,
        width,
        height,
        class: className = "",
        loading = "lazy",
        priority = false,
        sizes = "100vw",
        quality = 80
    }: Props = $props();

    // Generate optimized image URL with WebP support
    const optimizedSrc = $derived(() => {
        if (!src) return "";
        
        // If it's already an optimized URL or external, return as-is
        if (src.startsWith('http') || src.includes('data:')) {
            return src;
        }
        
        // For local images, you could implement image optimization here
        // This is a placeholder for future image optimization service
        return src;
    });

    // Generate srcset for responsive images
    const srcSet = $derived(() => {
        if (!src || src.startsWith('http') || src.includes('data:')) {
            return undefined;
        }
        
        // Generate different sizes for responsive images
        const sizes = [320, 640, 1024, 1280];
        return sizes
            .map(size => `${src}?w=${size}&q=${quality} ${size}w`)
            .join(', ');
    });
</script>

<img
    src={optimizedSrc}
    {alt}
    {width}
    {height}
    {loading}
    fetchpriority={priority ? "high" : "low"}
    decoding="async"
    class={className}
    sizes={sizes}
    srcset={srcSet}
    onerror="this.style.display='none'"
/>
