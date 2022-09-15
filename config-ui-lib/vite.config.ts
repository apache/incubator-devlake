import { resolve } from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    outDir: '_build',
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'config-ui',
    },
  },
});
