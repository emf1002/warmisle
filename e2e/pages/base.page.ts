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
    await expect(this.page.locator('.ant-message').first()).toContainText(message);
  }

  /** 断言弹窗可见 */
  async expectModalVisible() {
    await expect(this.page.locator('.ant-modal-wrap:visible').first()).toBeVisible();
  }

  /** 断言弹窗已关闭 */
  async expectModalHidden() {
    await expect(this.page.locator('.ant-modal-wrap:visible')).not.toBeVisible();
  }

  /** 确认弹窗的确认按钮（如删除确认） */
  async confirmModal() {
    await this.page.getByTestId('modal-confirm-btn').click();
  }

  /** 提交弹窗表单 */
  async submitModal() {
    await this.page.getByTestId('modal-submit-btn').click();
  }
}
