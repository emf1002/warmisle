import path from 'path';
import { defineConfig, devices } from '@playwright/test';
import { WORKERS, BASE_PORT, getBaseUrl } from './config';

const e2eDir = __dirname;

// When HC_TEST_PORT is set (by run-test.bat), use single-instance mode
// and isolate reports by port to avoid conflicts in parallel runs
const isSingleInstance = !!process.env.HC_TEST_PORT;
const portSuffix = process.env.HC_TEST_PORT || 'default';
const reportsDir = path.join(e2eDir, 'reports', portSuffix);
const effectiveBaseUrl = process.env.HC_TEST_BASE_URL || getBaseUrl('chromium');

export default defineConfig({
  testDir: './tests',
  // In single-instance mode (run-test.bat), skip global setup/teardown
  // since the server is managed externally
  ...(isSingleInstance ? {} : {
    globalSetup: './global-setup.ts',
    globalTeardown: './global-teardown.ts',
  }),
  snapshotDir: path.join(reportsDir, 'snapshots'),
  outputDir: path.join(reportsDir, 'output'),
  fullyParallel: true,
  forbidOnly: true,
  retries: 0,
  workers: isSingleInstance ? 1 : WORKERS,
  reporter: [
    ['list'],
    ['html', { open: 'never', outputFolder: path.join(reportsDir, 'html') }],
    ['json', { outputFile: path.join(reportsDir, 'results.json') }],
  ],
  use: {
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
    actionTimeout: 5000,
    navigationTimeout: 30000,
  },
  projects: [
    {
      name: 'chromium',
      use: {
        viewport: { width: 1280, height: 720 },
        baseURL: effectiveBaseUrl,
      },
    },
    {
      name: 'mobile',
      use: {
        ...devices['iPhone 13'],
        baseURL: effectiveBaseUrl,
      },
      grep: /@mobile/,
    },
  ],
});
