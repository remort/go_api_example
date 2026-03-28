import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api': {
        target: process.env.VITE_API_TARGET || 'http://0.0.0.0:8000',
        changeOrigin: true,
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            console.log(`Proxying: ${req.method} ${req.url} -> ${options.target}${req.url}`);
          });
        }
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: true
  }
})
