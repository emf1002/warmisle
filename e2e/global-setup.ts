import { rmSync, existsSync } from 'fs';
import path from 'path';

const e2eDir = __dirname;
const reportsDir = path.join(e2eDir, 'reports');

const preserve = new Set([path.join(reportsDir, 'snapshots')]);

export default async function globalSetup() {
  if (!existsSync(reportsDir)) return;

  const { readdirSync, statSync } = await import('fs');
  for (const entry of readdirSync(reportsDir)) {
    const full = path.join(reportsDir, entry);
    if (preserve.has(full)) continue;
    rmSync(full, { recursive: true, force: true });
  }
}
