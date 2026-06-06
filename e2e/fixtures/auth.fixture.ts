import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, createMember, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';
import { getBaseUrl } from '../config';

type AuthFixtures = {
  authenticated: {
    page: Page;
    token: string;
    seedLedgers: (options?: SeedLedgersOptions) => Promise<SeedLedgersResult>;
  };
  memberContext: {
    page: Page;
    token: string;
  };
};

export const test = base.extend<AuthFixtures>({
  authenticated: async ({ page, browser }, use, testInfo) => {
    process.env.HC_TEST_BASE_URL = getBaseUrl(testInfo.project.name);
    await resetDatabase();
    const { token } = await initAdmin();

    await page.addInitScript((t) => {
      localStorage.setItem('token', t);
    }, token);

    const boundSeedLedgers = (options?: SeedLedgersOptions) =>
      seedLedgers(token, options);

    await use({ page, token, seedLedgers: boundSeedLedgers });
  },

  memberContext: async ({ browser }, use, testInfo) => {
    process.env.HC_TEST_BASE_URL = getBaseUrl(testInfo.project.name);
    const memberResult = await createMember();
    const memberToken = memberResult.data.token;

    const context = await browser.newContext();
    const page = await context.newPage();

    await page.addInitScript((t) => {
      localStorage.setItem('token', t);
    }, memberToken);

    await use({ page, token: memberToken });

    await context.close();
  },
});

export { expect } from '@playwright/test';
