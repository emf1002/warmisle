import { spawn, type ChildProcess } from 'child_process';
import path from 'path';
import fs from 'fs';
import { WORKERS, BASE_PORT, PROJECTS } from './config';

const PORT_MAP_FILE = path.join(__dirname, '.e2e-port-map.json');

const BINARY = path.resolve(__dirname, '..', 'dist', 'warmisle.exe');
const E2E_DIR = __dirname;

declare global {
  var __e2e_servers__: ChildProcess[];
}

async function waitForReady(port: number, timeoutMs = 15000): Promise<void> {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const res = await fetch(`http://localhost:${port}/api/init/check`);
      if (res.ok) return;
    } catch {}
    await new Promise(r => setTimeout(r, 200));
  }
  throw new Error(`Server on port ${port} not ready within ${timeoutMs}ms`);
}

export default async function globalSetup() {
  globalThis.__e2e_servers__ = [];

  const cleanup = () => {
    for (const proc of globalThis.__e2e_servers__) {
      try { proc.kill(); } catch {}
    }
  };
  process.on('exit', cleanup);
  process.on('SIGINT', () => { cleanup(); process.exit(1); });
  process.on('SIGTERM', () => { cleanup(); process.exit(1); });

  // Spawn one server per project, each with its own DB and port
  const ready: Promise<void>[] = [];
  for (let i = 0; i < PROJECTS.length; i++) {
    const projectName = PROJECTS[i];
    const port = BASE_PORT + i;
    const dbPath = path.join(E2E_DIR, 'e2e-data', `test-${projectName}.db`);

    // Clean stale DB files
    for (const suffix of ['', '-wal', '-shm']) {
      try { fs.unlinkSync(dbPath + suffix); } catch {}
    }

    // Set env so test workers can discover their project's port
    process.env[`HC_PORT_${projectName}`] = String(port);

    const proc = spawn(BINARY, [], {
      cwd: E2E_DIR,
      env: {
        ...process.env,
        HC_PORT: String(port),
        HC_DB_PATH: dbPath,
        HC_TEST_MODE: 'true',
      },
      stdio: 'pipe',
    });

    proc.on('error', (err) => {
      console.error(`[global-setup] Server "${projectName}" (port ${port}) error:`, err.message);
    });
    proc.on('exit', (code) => {
      console.log(`[global-setup] Server "${projectName}" (port ${port}) exited with code ${code}`);
    });

    globalThis.__e2e_servers__.push(proc);
    ready.push(waitForReady(port));
  }

  await Promise.all(ready);

  // Write port mapping to file as fallback for worker threads
  const portMap: Record<string, number> = {};
  PROJECTS.forEach((p, i) => { portMap[p] = BASE_PORT + i; });
  fs.writeFileSync(PORT_MAP_FILE, JSON.stringify(portMap));

  console.log(`[global-setup] ${PROJECTS.length} server instances ready: ${PROJECTS.map((p, i) => `${p}->:${BASE_PORT + i}`).join(', ')}`);
}
