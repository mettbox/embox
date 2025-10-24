import { defineConfig } from 'vite';
import path from 'path';
import vue from '@vitejs/plugin-vue';
import { VitePWA, VitePWAOptions } from 'vite-plugin-pwa';

const pwa: Partial<VitePWAOptions> = {
  registerType: 'autoUpdate',
  includeAssets: ['manifest.json', 'favicon.svg', 'robots.txt', 'apple-touch-icon.png'],
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
            maxEntries: 200,
            maxAgeSeconds: 60 * 60 * 24 * 30, // 30 days
          },
        },
      },
    ],
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
});
