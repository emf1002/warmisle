import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, createMember, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';

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
  authenticated: async ({ page, browser }, use) => {
    await resetDatabase();
    const { token } = await initAdmin();

    await page.addInitScript((token) => {
      localStorage.setItem('token', token);
    }, token);

    const boundSeedLedgers = (options?: SeedLedgersOptions) =>
      seedLedgers(token, options);

    await use({ page, token, seedLedgers: boundSeedLedgers });
  },

  memberContext: async ({ browser }, use) => {
    const memberResult = await createMember();
    const memberToken = memberResult.data.token;

    const context = await browser.newContext();
    const page = await context.newPage();

    await page.addInitScript((token) => {
      localStorage.setItem('token', token);
    }, memberToken);

    await use({ page, token: memberToken });

    await context.close();
  },
});

export { expect } from '@playwright/test';
