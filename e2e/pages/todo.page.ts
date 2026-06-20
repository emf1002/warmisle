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
    await this.expectModalVisible();
  }

  async fillTitle(title: string) {
    await this.page.getByTestId('title-input').fill(title);
  }

  async fillDescription(desc: string) {
    await this.page.getByTestId('desc-input').fill(desc);
  }

  async selectPriority(priority: string) {
    await this.page.getByTestId('priority-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: priority }).click();
  }

  async submit() {
    await this.page.getByTestId('submit-btn').click();
    await this.expectModalHidden();
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
    await this.confirmModal();
  }

  async filterByStatus(status: string) {
    await this.page.getByTestId('status-filter').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: status }).click();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
  }

  /** 填写截止日期（格式 YYYY-MM-DD） */
  async fillDueDate(date: string) {
    const input = this.page.getByPlaceholder('请选择日期');
    await input.click();
    await input.fill(date);
    await this.page.keyboard.press('Enter');
  }

  /** 选择指派成员 */
  async selectAssignee(name: string) {
    await this.page.getByTestId('assignee-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: name }).first().click();
  }

  /** 认领第 n 条未指派的待办 */
  async claimTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('claim-btn').click();
  }

  /** 按指派成员筛选 */
  async filterByAssignee(name: string) {
    await this.page.getByTestId('assignee-filter').click();
    await this.page.locator('.ant-select-dropdown:visible .ant-select-item').filter({ hasText: name }).first().click({ timeout: 5000 });
  }

  /** 编辑第 n 条待办 */
  async editTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('edit-btn').click();
    await this.expectModalVisible();
  }

  /** 断言第 n 条待办的优先级标签 */
  async expectTodoPriority(index: number, priority: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('priority-tag')).toContainText(priority);
  }

  /** 断言第 n 条待办的截止日期显示 */
  async expectTodoDueDate(index: number, dateText: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('due-date')).toContainText(dateText);
  }

  /** 断言第 n 条待办的指派人 */
  async expectTodoAssignee(index: number, name: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('assignee-name')).toContainText(name);
  }

  /** 断言第 n 条待办的截止日期为红色（过期） */
  async expectOverdueHighlight(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('due-date')).toHaveClass(/overdue|red|error/);
  }

  /** 断言第 n 条待办已完成 */
  async expectTodoCompleted(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('todo-checkbox')).toBeChecked();
  }
}
