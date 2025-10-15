import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
   plugins: [tailwindcss(), sveltekit()],

   server: {
      proxy: {
         '/api': {
            target: 'http://localhost:8080',
            changeOrigin: true,
            secure: false,
            ws: true,
         },
         '/uploads': {
            target: 'http://localhost:8080',
            changeOrigin: true,
            secure: false,
            ws: true,
         },
      },
   },

   build: {
      target: 'esnext',
      minify: 'esbuild',
      sourcemap: false,
      rollupOptions: {
         output: {
            manualChunks: {
               vendor: ['svelte', '@sveltejs/kit'],
               icons: ['lucide-svelte']
            }
         }
      }
   },

   optimizeDeps: {
      include: ['lucide-svelte', '@sveltejs/kit']
   }
})
