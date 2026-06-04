import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Saat dev, request ke /api diteruskan (proxy) ke backend Go di port 8080
// sehingga frontend cukup memanggil "/api/..." tanpa masalah CORS.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
