<script lang="ts">
    import { onMount } from 'svelte';

    interface Props {
        enabled?: boolean;
    }

    const { enabled = false }: Props = $props();

    onMount(() => {
        if (!enabled || typeof window === 'undefined') return;

        // Monitor Core Web Vitals
        const observer = new PerformanceObserver((list) => {
            for (const entry of list.getEntries()) {
                if (entry.entryType === 'largest-contentful-paint') {
                    console.log('LCP:', entry.startTime);
                } else if (entry.entryType === 'first-input') {
                    console.log('FID:', entry.processingStart - entry.startTime);
                } else if (entry.entryType === 'layout-shift') {
                    console.log('CLS:', entry.value);
                }
            }
        });

        observer.observe({ entryTypes: ['largest-contentful-paint', 'first-input', 'layout-shift'] });

        // Monitor bundle size
        if ('performance' in window) {
            const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
            console.log('Page load time:', navigation.loadEventEnd - navigation.loadEventStart);
        }

        return () => observer.disconnect();
    });
</script>

<!-- This component doesn't render anything, it just monitors performance -->
