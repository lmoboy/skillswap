import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
   plugins: [tailwindcss(), sveltekit()],

   server: {
      proxy: {
         '/api': {
            target: 'http://0.0.0.0:8080',
            changeOrigin: true,
            secure: false,
            ws: true,
         },
         '/uploads': {
            target: 'http://0.0.0.0:8080',
            changeOrigin: true,
            secure: false,
            ws: true,
         },
      },
   },
})
