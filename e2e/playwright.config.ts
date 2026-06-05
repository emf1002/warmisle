import path from 'path';
import { defineConfig, devices } from '@playwright/test';

const rootDir = path.resolve(__dirname, '..');
const e2eDir = __dirname;
const reportsDir = path.join(e2eDir, 'reports');
const binaryPath = path.join(rootDir, 'dist', 'warmisle.exe');

export default defineConfig({
  testDir: './tests',
  globalSetup: './global-setup.ts',
  snapshotDir: path.join(reportsDir, 'snapshots'),
  outputDir: path.join(reportsDir, 'output'),
  fullyParallel: false,
  forbidOnly: true,
  retries: 0,
  workers: 1,
  reporter: [
    ['list'],
    ['html', { open: 'never', outputFolder: path.join(reportsDir, 'html') }],
    ['json', { outputFile: path.join(reportsDir, 'results.json') }],
  ],
  use: {
    baseURL: 'http://localhost:8080',
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
    actionTimeout: 5000,
    navigationTimeout: 5000,
  },
  projects: [
    {
      name: 'chromium',
      use: {
        viewport: { width: 1280, height: 720 },
      },
    },
    {
      name: 'mobile',
      use: {
        ...devices['iPhone 13'],
      },
      grep: /@mobile/,
    },
  ],
  webServer: {
    command: binaryPath,
    cwd: e2eDir,
    port: 8080,
    timeout: 30000,
    reuseExistingServer: false,
    env: {
      HC_DB_PATH: path.join(e2eDir, 'e2e-data', 'test.db'),
      HC_TEST_MODE: 'true',
      HC_PORT: '8080',
    },
  },
});
