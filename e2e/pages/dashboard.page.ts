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
}
