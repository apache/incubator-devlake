import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  testMatch: '*.spec.ts',
  timeout: 180000,
  expect: {
    timeout: 10000,
  },
  use: {
    baseURL: 'http://localhost:4000',
    screenshot: 'on',
    trace: 'on-first-retry',
  },
  reporter: [['html', { open: 'never' }], ['list']],
  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium', viewport: { width: 1440, height: 900 } },
    },
  ],
});
