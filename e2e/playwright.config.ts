import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: true,
  retries: 0,
  workers: 1,
  reporter: [['html', { open: 'never' }]],
  use: {
    baseURL: 'http://localhost:8080',
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: {
        viewport: { width: 1280, height: 720 },
      },
    },
  ],
  webServer: {
    command: '../dist/warmisle.exe',
    cwd: '..',
    port: 8080,
    timeout: 30000,
    reuseExistingServer: true,
    env: {
      HC_DB_PATH: './e2e-data/test.db',
      HC_TEST_MODE: 'true',
      HC_PORT: '8080',
    },
  },
});
