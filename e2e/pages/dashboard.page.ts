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
    const input = this.page.getByPlaceholder('请选择月份');
    await input.click();
    // Wait for picker dropdown to appear
    await this.page.locator('.ant-picker-dropdown:visible').waitFor();
    // In month picker, click the next month cell (skip current month, pick the one after)
    // The panel shows months as grid cells - click the right arrow to go to next year if needed,
    // or simply click a different month cell
    const cells = this.page.locator('.ant-picker-dropdown:visible .ant-picker-cell-inner');
    const count = await cells.count();
    if (count > 0) {
      // Click the second cell (next month) to change selection
      await cells.nth(Math.min(1, count - 1)).click();
    }
    await this.page.waitForLoadState('networkidle');
  }

  async goPrevMonth() {
    const input = this.page.getByPlaceholder('请选择月份');
    await input.click();
    // Wait for picker dropdown to appear
    await this.page.locator('.ant-picker-dropdown:visible').waitFor();
    // In month picker, select a different month - click the prev year button first
    const prevYearBtn = this.page.locator('.ant-picker-dropdown:visible .ant-picker-super-prev-btn');
    if (await prevYearBtn.isVisible()) {
      await prevYearBtn.click();
    }
    // Then click the last month cell (December of previous year)
    const cells = this.page.locator('.ant-picker-dropdown:visible .ant-picker-cell-inner');
    const count = await cells.count();
    if (count > 0) {
      await cells.nth(count - 1).click();
    }
    await this.page.waitForLoadState('networkidle');
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
