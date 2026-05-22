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
    await expect(this.page.getByTestId('wish-modal')).toBeVisible();
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
}
