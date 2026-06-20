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
    await this.expectModalVisible();
  }

  async selectType(type: string) {
    await this.page.getByTestId('type-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: type }).click();
  }

  async fillName(name: string) {
    await this.page.getByTestId('name-input').fill(name);
  }

  async submit() {
    await this.submitModal();
    await this.expectModalHidden();
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
    await this.confirmModal();
  }

  async changeCategoryType(type: string) {
    await this.page.getByTestId('type-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: type }).click();
  }

  async fillEditName(name: string) {
    await this.page.getByTestId('name-input').fill(name);
  }

  async submitEdit() {
    await this.submitModal();
    await this.expectModalHidden();
  }

  async expectDeleteDisabled(index: number) {
    const cards = this.page.getByTestId(/^category-card-/);
    await expect(cards.nth(index).getByTestId('delete-btn')).toBeDisabled();
  }
}
