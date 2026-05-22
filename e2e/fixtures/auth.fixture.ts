import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin } from './db.fixture';

type AuthFixtures = {
  authenticated: { page: Page; token: string };
};

export const test = base.extend<AuthFixtures>({
  authenticated: async ({ page, browser }, use) => {
    // 重置数据库并初始化管理员
    await resetDatabase();
    const { token } = await initAdmin();

    // 注入 token 到 localStorage
    await page.addInitScript((token) => {
      localStorage.setItem('token', token);
    }, token);

    await use({ page, token });
  },
});

export { expect } from '@playwright/test';
