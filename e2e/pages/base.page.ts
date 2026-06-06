import { type Page, type Locator, expect } from '@playwright/test';

export class BasePage {
  constructor(protected page: Page) {}

  async navigate(path: string) {
    const target = `/#${path}`;
    // If already on the same URL, goto() is a no-op for hash routing.
    // Force a reload so the component re-fetches fresh data.
    if (this.page.url().endsWith(target)) {
      await this.page.reload({ waitUntil: 'domcontentloaded' });
    } else {
      await this.page.goto(target, { waitUntil: 'domcontentloaded' });
    }
    await this.page.waitForLoadState('networkidle');
  }

  async waitForPage() {
    await this.page.waitForLoadState('networkidle');
  }

  async screenshot(name: string) {
    await expect(this.page).toHaveScreenshot(name, {
      fullPage: true,
      maxDiffPixelRatio: 0.01,
    });
  }

  async screenshotComponent(selector: string, name: string) {
    await expect(this.page.locator(selector)).toHaveScreenshot(name, {
      maxDiffPixelRatio: 0.01,
    });
  }

  async expectToast(message: string) {
    await expect(this.page.locator('.ant-message')).toContainText(message);
  }
}
