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

export interface SeedLedgersOptions {
  count?: number;
  startDate?: string;
  endDate?: string;
}

export interface SeedLedgersResult {
  code: number;
  message: string;
  data: {
    count: number;
    summary: { income: number; expense: number; balance: number };
    expense_category_count: number;
    income_category_count: number;
  };
}

export async function seedLedgers(
  token: string,
  options?: SeedLedgersOptions
): Promise<SeedLedgersResult> {
  const res = await fetch(`${BASE_URL}/api/test/seed-ledgers`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      count: options?.count ?? 35,
      start_date: options?.startDate,
      end_date: options?.endDate,
    }),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Seed ledgers failed: ${res.status} ${text}`);
  }

  return res.json();
}

export interface CreateMemberResult {
  code: number;
  message: string;
  data: {
    token: string;
    member_id: number;
  };
}

export async function createMember(
  options: { username?: string; password?: string; name?: string } = {}
): Promise<CreateMemberResult> {
  const res = await fetch(`${BASE_URL}/api/test/create-member`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: options.username ?? 'member1',
      password: options.password ?? 'test123',
      name: options.name ?? '成员一',
    }),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Create member failed: ${res.status} ${text}`);
  }

  return res.json();
}
