const BASE_URL = 'http://localhost:8080';

export async function resetDatabase(): Promise<void> {
  const res = await fetch(`${BASE_URL}/api/test/reset`, { method: 'POST' });
  if (!res.ok) {
    throw new Error(`Database reset failed: ${res.status}`);
  }
}

export async function initAdmin(): Promise<{ token: string }> {
  // 初始化管理员
  const initRes = await fetch(`${BASE_URL}/api/init/setup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      admin_name: '管理员',
      username: 'admin',
      password: 'test123',
    }),
  });

  if (!initRes.ok) {
    const text = await initRes.text();
    throw new Error(`Admin init failed: ${initRes.status} ${text}`);
  }

  const initData = await initRes.json();
  return { token: initData.data.token };
}
