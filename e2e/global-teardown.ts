import { type ChildProcess } from 'child_process';

declare global {
  var __e2e_servers__: ChildProcess[];
}

export default async function globalTeardown() {
  const servers = globalThis.__e2e_servers__ || [];
  for (const proc of servers) {
    try {
      proc.kill();
    } catch {}
  }
  if (servers.length > 0) {
    console.log(`[global-teardown] Stopped ${servers.length} server instances`);
  }
}
