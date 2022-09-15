import { resolve } from 'path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],

  root: resolve(__dirname, 'website'),

  build: {
    outDir: resolve(__dirname, '_website'),
  },
});
