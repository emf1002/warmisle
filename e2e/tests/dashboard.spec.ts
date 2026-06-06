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
    // Create a todo so the section has data
    await page.goto('/#/todo');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试待办');
    await page.getByTestId('submit-btn').click();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectUpcomingTodosVisible();
  });

  test('愿望趋势区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    // Create a wish so the section has data
    await page.goto('/#/wish');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试愿望');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
  });

  test('论坛热门区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    // Create a topic so the section has data
    await page.goto('/#/forum');
    await page.getByTestId('create-topic-btn').click();
    await page.getByTestId('topic-title').fill('测试话题');
    await page.getByTestId('topic-content').fill('内容');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
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
    expect(incomeText).toMatch(/¥[\d,.]+/);
    const expenseText = await page.getByTestId('summary-expense').textContent();
    expect(expenseText).toMatch(/¥[\d,.]+/);
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

  test('移动端仪表盘视觉回归', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    await page.setViewportSize({ width: 390, height: 844 });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.screenshot('dashboard-mobile.png');
  });

  test('移动端仪表盘', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page, token } = authenticated;
    await page.setViewportSize({ width: 390, height: 844 });
    // Create test data via API to avoid mobile modal timing issues
    const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };
    await page.request.post('http://localhost:8080/api/todos', { headers, data: { title: '测试待办' } });
    await page.request.post('http://localhost:8080/api/wishes', { headers, data: { title: '测试愿望', type: 'personal' } });
    await page.request.post('http://localhost:8080/api/posts', { headers, data: { content: '测试动态', type: 'post' } });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    await dashboard.expectUpcomingTodosVisible();
    await dashboard.expectWishTrendsVisible();
    await dashboard.expectForumHotVisible();
  });
});
