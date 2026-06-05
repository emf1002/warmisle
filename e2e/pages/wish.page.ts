import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class WishPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/wish');
  }

  async expectOnWish() {
    await expect(this.page).toHaveURL(/\/#\/wish/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async fillTitle(title: string) {
    await this.page.getByTestId('title-input').fill(title);
  }

  async fillDescription(desc: string) {
    await this.page.getByTestId('desc-input').fill(desc);
  }

  async selectCategory(category: string) {
    await this.page.getByTestId('category-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: category }).click();
  }

  async fillAmount(amount: string) {
    const input = this.page.locator('.ant-modal:visible').getByRole('spinbutton');
    await input.click();
    await input.fill(amount);
  }

  async submit() {
    const responsePromise = this.page.waitForResponse(
      (resp) => resp.url().includes('/api/wishes') && ['POST', 'PUT'].includes(resp.request().method()),
      { timeout: 5000 },
    ).catch(() => null);
    await this.page.locator('.ant-modal-footer .ant-btn-primary').click();
    await responsePromise;
    // Wait for modal close animation to finish
    await this.page.locator('.ant-modal-wrap:visible').waitFor({ state: 'hidden', timeout: 5000 }).catch(() => {});
  }

  async expectWishCount(count: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items).toHaveCount(count);
  }

  async switchType(type: string) {
    await this.page.getByTestId('type-switch').locator(`text=${type}`).click();
  }

  async filterByStatus(status: string) {
    await this.page.getByTestId('status-filter').click();
    await this.page.locator('.ant-select-item-option', { hasText: status }).click();
  }

  async voteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('vote-btn').click();
  }

  /** 取消投票（再次点击投票按钮，触发已投票→取消投票流程） */
  async unvoteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    // Set up listener BEFORE click: vote fails → modal appears → confirm triggers DELETE /api/wishes/:id/vote → fetchWishes
    const unvotePromise = this.page.waitForResponse(
      (resp) => resp.url().includes('/vote') && resp.request().method() === 'DELETE',
      { timeout: 5000 },
    ).catch(() => null);
    await items.nth(index).getByTestId('vote-btn').click();
    // Wait for the "已投票" error response to arrive
    await this.page.waitForResponse(
      (resp) => resp.url().includes('/vote') && resp.request().method() === 'POST',
      { timeout: 5000 },
    ).catch(() => null);
    // Wait for the unvote confirmation modal
    await this.page.waitForSelector('.ant-modal-confirm:visible', { timeout: 5000 }).catch(() => {});
    // Confirm the unvote
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
    // Wait for the unvote API call to complete
    await unvotePromise;
    // Wait for list refresh to complete
    await this.page.waitForResponse(
      (resp) => resp.url().includes('/api/wishes') && resp.request().method() === 'GET',
      { timeout: 5000 },
    ).catch(() => null);
  }

  /** 断言投票人数 */
  async expectVoteCount(index: number, count: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items.nth(index).getByTestId('vote-count')).toContainText(String(count));
  }

  /** 变更愿望状态（管理员操作） */
  async changeWishStatus(index: number, status: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.getByRole('menuitem', { name: status }).click();
    await this.page.waitForTimeout(300);
  }

  /** 断言愿望状态 */
  async expectStatus(index: number, status: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items.nth(index).getByTestId('status-tag')).toContainText(status);
  }

  /** 创建者放弃愿望 */
  async abandonWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.getByRole('menuitem', { name: '标记为放弃' }).click();
  }

  /** 删除愿望 */
  async deleteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.getByRole('menuitem', { name: '删除' }).click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
  }

  /** 评论愿望 */
  async commentOnWish(index: number, text: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.getByRole('menuitem', { name: '评论' }).click();
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await this.page.waitForTimeout(300);
  }
}
