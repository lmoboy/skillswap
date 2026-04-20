import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
   plugins: [tailwindcss(), sveltekit()],

   server: {
      host: '0.0.0.0',
      port: 5173,
      allowedHosts: ['.ngrok-free.dev', 'undecided-groggily-bonsai.ngrok-free.dev'],
      hmr: {
         clientPort: 443,
      },
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
})
