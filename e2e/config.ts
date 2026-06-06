// Shared E2E config — used by playwright.config.ts, global-setup.ts, and test fixtures.
import path from 'path';
import fs from 'fs';

export const WORKERS = 2;
export const BASE_PORT = 8081;
export const PROJECTS = ['chromium', 'mobile'] as const;

const PORT_MAP_FILE = path.join(__dirname, '.e2e-port-map.json');

function readPortMap(): Record<string, number> {
  try {
    return JSON.parse(fs.readFileSync(PORT_MAP_FILE, 'utf-8'));
  } catch {
    return {};
  }
}

/** Get the API base URL for a project. Respects HC_TEST_PORT env var from run-test.bat. */
export function getBaseUrl(projectName: string): string {
  // 1. Explicit env var from run-test.bat (single-instance mode)
  if (process.env.HC_TEST_PORT) {
    return `http://localhost:${process.env.HC_TEST_PORT}`;
  }
  if (process.env.HC_TEST_BASE_URL) {
    return process.env.HC_TEST_BASE_URL;
  }
  // 2. Env var set by global-setup in main process (multi-instance mode)
  const envPort = process.env[`HC_PORT_${projectName}`];
  if (envPort) return `http://localhost:${envPort}`;
  // 3. Port map file fallback
  const map = readPortMap();
  if (map[projectName]) return `http://localhost:${map[projectName]}`;
  // 4. Default
  return `http://localhost:${BASE_PORT}`;
}
