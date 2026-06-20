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
    await this.expectModalVisible();
  }

  async fillTitle(title: string) {
    await this.page.getByTestId('title-input').fill(title);
  }

  async fillDescription(desc: string) {
    await this.page.getByTestId('desc-input').fill(desc);
  }

  async selectCategory(category: string) {
    await this.page.getByTestId('category-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: category }).click();
  }

  async fillAmount(amount: string) {
    const input = this.page.getByTestId('amount-input');
    await input.click();
    await input.fill(amount);
  }

  async submit() {
    const responsePromise = this.page.waitForResponse(
      (resp) => resp.url().includes('/api/wishes') && ['POST', 'PUT'].includes(resp.request().method()),
      { timeout: 5000 },
    ).catch(() => null);
    await super.submitModal();
    await responsePromise;
    await this.expectModalHidden();
  }

  async expectWishCount(count: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items).toHaveCount(count);
  }

  async switchType(type: string) {
    await this.page.getByTestId('type-switch').getByText(type, { exact: false }).click({ timeout: 5000 });
  }

  async filterByStatus(status: string) {
    await this.page.getByTestId('status-filter').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: status }).click();
  }

  async voteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('vote-btn').click();
  }

  /** 取消投票（再次点击投票按钮，触发已投票→取消投票流程） */
  async unvoteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
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
    await this.expectModalVisible();
    // Confirm the unvote
    await this.confirmModal();
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
    await this.page.locator('.ant-dropdown-menu-item').filter({ hasText: status }).click();
    await this.page.waitForLoadState('networkidle');
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
    await this.page.locator('.ant-dropdown-menu-item').filter({ hasText: '标记为放弃' }).click();
  }

  /** 将个人愿望提升为家庭愿望 */
  async promoteToFamily(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.locator('.ant-dropdown-menu-item').filter({ hasText: '提升为家庭愿望' }).click();
    // 提升成功可能弹出确认/提示弹窗，关闭后继续
    await this.page.waitForTimeout(500);
    await this.confirmModal().catch(() => {});
    await this.expectModalHidden().catch(() => {});
  }

  /** 删除愿望 */
  async deleteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.locator('.ant-dropdown-menu-item').filter({ hasText: '删除' }).click();
    await this.confirmModal();
  }

  /** 评论愿望 */
  async commentOnWish(index: number, text: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.locator('.ant-dropdown-menu-item').filter({ hasText: '评论' }).click();
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await expect(this.page.getByTestId('comment-list')).toContainText(text);
  }
}
