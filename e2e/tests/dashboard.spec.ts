// e2e/tests/dashboard.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { DashboardPage } from '../pages/dashboard.page';
import { WishPage } from '../pages/wish.page';
import { ForumPage } from '../pages/forum.page';

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
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('测试愿望');
    await wish.submit();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
  });

  test('论坛热门区域可见', async ({ authenticated }) => {
    const { page } = authenticated;
    // Create a topic so the section has data
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试话题');
    await forum.fillTopicContent('内容');
    await forum.submitModal();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectForumHotVisible();
  });


  test('月份切换', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    const monthText = await page.getByPlaceholder('请选择月份').inputValue();
    await dashboard.goPrevMonth();
    const newMonthText = await page.getByPlaceholder('请选择月份').inputValue();
    expect(newMonthText).not.toBe(monthText);
  });

  test('无数据时显示零值', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    // 新系统无数据时，汇总应显示 ¥0.00
    await expect(page.getByTestId('summary-income')).toContainText('¥0.00');
    await expect(page.getByTestId('summary-expense')).toContainText('¥0.00');
    await expect(page.getByTestId('summary-balance')).toContainText('¥0.00');
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
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('测试愿望');
    await wish.submit();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
    await dashboard.expectWishTrendCount(1);
  });

  test('论坛热点数据', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试话题');
    await forum.fillTopicContent('内容');
    await forum.submitModal();
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectForumHotVisible();
    await dashboard.expectForumHotCount(1);
  });


  test('移动端仪表盘', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page, token } = authenticated;
    await page.setViewportSize({ width: 390, height: 844 });
    // Create test data via API to avoid mobile modal timing issues
    const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' };
    await page.request.post('/api/todos', { headers, data: { title: '测试待办' } });
    await page.request.post('/api/wishes', { headers, data: { title: '测试愿望', type: 'personal' } });
    await page.request.post('/api/posts', { headers, data: { content: '测试动态', type: 'post' } });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    await dashboard.expectUpcomingTodosVisible();
    await dashboard.expectWishTrendsVisible();
    await dashboard.expectForumHotVisible();
  });
});
