import { defineConfig } from 'vite';
import path from 'path';
import vue from '@vitejs/plugin-vue';
import { VitePWA, VitePWAOptions } from 'vite-plugin-pwa';

const pwa: Partial<VitePWAOptions> = {
  registerType: 'autoUpdate',
  includeAssets: ['manifest.json', 'favicon.ico', 'robots.txt', 'apple-touch-icon.png'],
  workbox: {
    runtimeCaching: [
      {
        // All except http requests
        urlPattern: ({ request }) =>
          request.destination === 'document' ||
          request.destination === 'style' ||
          request.destination === 'script' ||
          request.destination === 'image' ||
          request.destination === 'font',
        handler: 'CacheFirst',
        options: {
          cacheName: 'assets-cache',
          expiration: {
            maxEntries: 1000,
            maxAgeSeconds: 60 * 60 * 24 * 30, // 30 days
          },
        },
      },
    ],
  },
};

const devProxy = {
  '/api': {
    target: 'http://localhost:2705',
    changeOrigin: true,
    rewrite: (path: string) => path.replace(/^\/api/, ''),
    configure: (proxy: any) => {
      proxy.on('proxyRes', (proxyRes: any) => {
        const cookies = proxyRes.headers['set-cookie'];
        if (cookies) {
          proxyRes.headers['set-cookie'] = (Array.isArray(cookies) ? cookies : [cookies]).map(
            (c: string) => c.replace(/;\s*domain=[^;]*/gi, ''),
          );
        }
      });
    },
  },
};

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), VitePWA(pwa)],
  define: {
    'process.env': process.env, // Optional, falls du `process.env` verwenden möchtest
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: devProxy,
  },
  preview: {
    proxy: devProxy,
  },
});
