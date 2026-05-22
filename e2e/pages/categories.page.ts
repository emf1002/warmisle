import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class CategoriesPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/categories');
  }

  async expectOnCategories() {
    await expect(this.page).toHaveURL(/\/#\/categories/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.getByTestId('category-modal')).toBeVisible();
  }

  async selectType(type: string) {
    await this.page.getByTestId('type-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: type }).click();
  }

  async fillName(name: string) {
    await this.page.getByTestId('name-input').fill(name);
  }

  async submit() {
    await this.page.locator('.ant-modal-footer .ant-btn-primary').click();
    // 等待弹窗关闭，确认提交成功
    await expect(this.page.getByTestId('category-modal')).not.toBeVisible();
  }

  async expectExpenseCategoryCount(count: number) {
    const cards = this.page.getByTestId('expense-categories').getByTestId(/^category-card-/);
    await expect(cards).toHaveCount(count);
  }

  async expectIncomeCategoryCount(count: number) {
    const cards = this.page.getByTestId('income-categories').getByTestId(/^category-card-/);
    await expect(cards).toHaveCount(count);
  }

  async editCategory(index: number) {
    const cards = this.page.getByTestId(/^category-card-/);
    await cards.nth(index).getByTestId('edit-btn').click();
  }

  async deleteCategory(index: number) {
    const cards = this.page.getByTestId(/^category-card-/);
    await cards.nth(index).getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }
}
