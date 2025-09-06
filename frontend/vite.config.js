import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import dotenv from 'dotenv'

dotenv.config()

export default defineConfig(() => {
  const apiUrl = process.env.VITE_API_URL || 'http://api:8000'

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      host: true,        // слушает 0.0.0.0 внутри контейнера
      port: 5173,
      hmr: {
        host: 'localhost',   // так ты заходишь в браузере
        clientPort: 5173
      },
      watch: {
        usePolling: true,    // 👈 обязательно в Docker Desktop
        interval: 300        // интервал опроса (можно 300–500)
      },
      proxy: {
        '/api': {
          target: apiUrl,
          changeOrigin: true,
          secure: false,
        },
      },
    },
  }
})