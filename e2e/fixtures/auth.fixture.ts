import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';

type AuthFixtures = {
  authenticated: {
    page: Page;
    token: string;
    seedLedgers: (options?: SeedLedgersOptions) => Promise<SeedLedgersResult>;
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
});

export { expect } from '@playwright/test';
