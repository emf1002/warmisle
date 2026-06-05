// e2e/pages/dashboard.page.ts
import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class DashboardPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/');
  }

  async expectOnDashboard() {
    await expect(this.page).toHaveURL(/\/#\/$/);
  }

  async expectSummaryVisible() {
    await expect(this.page.getByTestId('summary-grid')).toBeVisible();
  }

  async goNextMonth() {
    await this.page.getByTestId('month-next').click();
  }

  async goPrevMonth() {
    await this.page.getByTestId('month-prev').click();
  }

  async expectUpcomingTodosVisible() {
    await expect(this.page.getByTestId('upcoming-todos')).toBeVisible();
  }

  async expectWishTrendsVisible() {
    await expect(this.page.getByTestId('wish-trends')).toBeVisible();
  }

  async expectForumHotVisible() {
    await expect(this.page.getByTestId('forum-hot')).toBeVisible();
  }

  /** 断言汇总数据包含指定值 */
  async expectSummaryData(values: { income?: string; expense?: string; balance?: string }) {
    if (values.income) {
      await expect(this.page.getByTestId('summary-income')).toContainText(values.income);
    }
    if (values.expense) {
      await expect(this.page.getByTestId('summary-expense')).toContainText(values.expense);
    }
    if (values.balance) {
      await expect(this.page.getByTestId('summary-balance')).toContainText(values.balance);
    }
  }

  /** 断言支出饼图可见 */
  async expectExpensePieChartVisible() {
    await expect(this.page.getByTestId('expense-chart')).toBeVisible();
  }

  /** 断言近期待办区域有待办项 */
  async expectUpcomingTodoCount(count: number) {
    const items = this.page.getByTestId('upcoming-todos').getByTestId(/^todo-link-/);
    await expect(items).toHaveCount(count);
  }

  /** 断言愿望动态区域有愿望项 */
  async expectWishTrendCount(count: number) {
    const items = this.page.getByTestId('wish-trends').getByTestId(/^wish-link-/);
    await expect(items).toHaveCount(count);
  }

  /** 断言论坛热门区域有话题项 */
  async expectForumHotCount(count: number) {
    const items = this.page.getByTestId('forum-hot').getByTestId(/^topic-link-/);
    await expect(items).toHaveCount(count);
  }
}
