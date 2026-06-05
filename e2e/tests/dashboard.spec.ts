// e2e/tests/dashboard.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { DashboardPage } from '../pages/dashboard.page';

test.describe('仪表盘', () => {
  test('页面加载显示统计区域', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
  });

  test('待办事项区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectUpcomingTodosVisible();
  });

  test('愿望趋势区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
  });

  test('论坛热门区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectForumHotVisible();
  });

  test('仪表盘页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.screenshot('dashboard-full.png');
  });

  test('月份切换', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    const monthText = await page.getByTestId('current-month').textContent();
    await dashboard.goPrevMonth();
    const newMonthText = await page.getByTestId('current-month').textContent();
    expect(newMonthText).not.toBe(monthText);
  });

  test('汇总数据正确性', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 20 });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    const incomeText = await page.getByTestId('summary-income').textContent();
    expect(incomeText).toMatch(/\+¥/);
    const expenseText = await page.getByTestId('summary-expense').textContent();
    expect(expenseText).toMatch(/-¥/);
  });

  test('支出分类饼图', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 20 });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectExpensePieChartVisible();
  });

  test('近期待办数据', async ({ authenticated }) => {
    const { page } = authenticated;
    await page.goto('/#/todo');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试待办');
    await page.getByTestId('submit-btn').click();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectUpcomingTodosVisible();
    await dashboard.expectUpcomingTodoCount(1);
  });

  test('愿望动态数据', async ({ authenticated }) => {
    const { page } = authenticated;
    await page.goto('/#/wish');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试愿望');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
  });

  test('论坛热点数据', async ({ authenticated }) => {
    const { page } = authenticated;
    await page.goto('/#/forum');
    await page.getByTestId('create-topic-btn').click();
    await page.getByTestId('topic-title').fill('测试话题');
    await page.getByTestId('topic-content').fill('内容');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectForumHotVisible();
  });

  test('移动端仪表盘', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    await dashboard.expectUpcomingTodosVisible();
    await dashboard.expectWishTrendsVisible();
    await dashboard.expectForumHotVisible();
  });
});
