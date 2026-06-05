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
    await this.page.getByTestId('amount-input').click();
    await this.page.getByTestId('amount-input').locator('input').fill(amount);
  }

  async submit() {
    await this.page.locator('.ant-modal-footer .ant-btn-primary').click();
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

  /** 取消投票（再次点击投票按钮） */
  async unvoteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('vote-btn').click();
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
    await this.page.locator('.ant-dropdown', { hasText: status }).click();
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
    await items.nth(index).getByTestId('abandon-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
  }

  /** 删除愿望 */
  async deleteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
  }

  /** 评论愿望 */
  async commentOnWish(index: number, text: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('comment-btn').click();
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await this.page.waitForTimeout(300);
  }
}
