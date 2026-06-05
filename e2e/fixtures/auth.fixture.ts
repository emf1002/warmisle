import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, createMember, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';

const BASE_URL = 'http://localhost:8080';

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
    // 1. Reset DB + create admin (API only, no page navigation)
    await resetDatabase();
    const { token } = await initAdmin();

    // 2. Register addInitScript BEFORE any navigation.
    //    The test will make the first navigation to the app, which triggers
    //    a full page load. addInitScript runs before Vue, setting the token
    //    in localStorage. Pinia reads it on mount.
    //    IMPORTANT: Do NOT navigate to the app here — a same-origin hash
    //    change (e.g. /#/ → /#/ledger) does not reload the page, so the
    //    Pinia store keeps its stale (empty) token value.
    await page.addInitScript((t) => {
      localStorage.setItem('token', t);
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

    await page.addInitScript((t) => {
      localStorage.setItem('token', t);
    }, memberToken);

    await use({ page, token: memberToken });

    await context.close();
  },
});

export { expect } from '@playwright/test';
