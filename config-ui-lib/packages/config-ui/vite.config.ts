import { resolve } from 'path';

import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],

  build: {
    outDir: '_build',
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'config-ui',
    },
    rollupOptions: {
      external: ['react'],
      output: {
        globals: {
          vue: 'React',
        },
      },
    },
  },
});
