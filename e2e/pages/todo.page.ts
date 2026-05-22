import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class TodoPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/todo');
  }

  async expectOnTodo() {
    await expect(this.page).toHaveURL(/\/#\/todo/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.getByTestId('todo-modal')).toBeVisible();
  }

  async fillTitle(title: string) {
    await this.page.getByTestId('title-input').fill(title);
  }

  async fillDescription(desc: string) {
    await this.page.getByTestId('desc-input').fill(desc);
  }

  async selectPriority(priority: string) {
    await this.page.getByTestId('priority-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: priority }).click();
  }

  async submit() {
    await this.page.getByTestId('submit-btn').click();
  }

  async toggleTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('todo-checkbox').click();
  }

  async expectTodoCount(count: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items).toHaveCount(count);
  }

  async expectTodoTitle(index: number, title: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index)).toContainText(title);
  }

  async deleteTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  async filterByStatus(status: string) {
    await this.page.getByTestId('status-filter').click();
    await this.page.locator('.ant-select-item-option', { hasText: status }).click();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
  }
}
