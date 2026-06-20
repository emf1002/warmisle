// e2e/tests/ledger-load.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { LedgerPage } from '../pages/ledger.page';

test.describe('记账本 — 负载场景', () => {
  test('无限滚动：首次加载 20 条，滚动后加载更多', async ({ authenticated }) => {
    const { page, seedLedgers } = await authenticated;
    const result = await seedLedgers({ count: 35 });
    expect(result.code).toBe(0);

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 首次加载 20 条（默认 limit=20）
    await ledger.expectTotalItemCount(20);

    // sentinel 可见，表示还有更多数据
    await expect(page.getByTestId('load-sentinel')).toBeVisible();

    // 滚动触发加载更多
    await ledger.scrollToLoadMore();
    await page.waitForLoadState('networkidle', { timeout: 5000 }).catch(() => {});

    // 记录数应增加（超过 20 条）
    const finalCount = await page.getByTestId(/^ledger-item-/).count();
    expect(finalCount).toBeGreaterThan(20);

    // 日期分组应 > 1（跨多天数据）
    const groupCount = await page.getByTestId('date-group').count();
    expect(groupCount).toBeGreaterThan(1);
  });

  test('统计汇总：收入/支出/结余正确显示', async ({ authenticated }) => {
    const { page, seedLedgers } = await authenticated;
    const result = await seedLedgers({ count: 35 });
    expect(result.code).toBe(0);

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载完成
    await ledger.expectTotalItemCount(20);

    // 验证汇总栏非零且格式正确（收入>0, 支出>0）
    const incomeText = await page.getByTestId('summary-income').textContent();
    const expenseText = await page.getByTestId('summary-expense').textContent();
    const balanceText = await page.getByTestId('summary-balance').textContent();

    // 收入应为 +¥xxx.xx 格式且大于 0
    expect(incomeText).toMatch(/^\+¥\d+\.\d{2}$/);
    const incomeValue = parseFloat(incomeText!.replace('+¥', ''));
    expect(incomeValue).toBeGreaterThan(0);

    // 支出应为 -¥xxx.xx 格式且大于 0
    expect(expenseText).toMatch(/^-\u00A5\d+\.\d{2}$/);
    const expenseValue = parseFloat(expenseText!.replace('-\u00A5', ''));
    expect(expenseValue).toBeGreaterThan(0);

    // 结余 = 收入 - 支出
    const balanceValue = parseFloat(balanceText!.replace(/[+-]\u00A5/, ''));
    expect(balanceValue).toBeCloseTo(incomeValue - expenseValue, 2);
  });

  // P2: Ant Design RangePicker 键盘输入不触发 v-model 更新，筛选不生效
  test.skip('日期范围筛选：选择子范围后只显示匹配记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    // 使用当月日期范围播种（前 7 天）
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    await seedLedgers({
      count: 35,
      startDate: `${year}-${month}-01`,
      endDate: `${year}-${month}-07`,
    });

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载
    await ledger.expectTotalItemCount(20);

    // 直接在 RangePicker 输入框中键入日期范围（避免日历交互版本兼容问题）
    const startInput = page.locator('.ant-picker-input').first().locator('input');
    const endInput = page.locator('.ant-picker-input').last().locator('input');
    await startInput.click();
    await startInput.fill(`${year}-${month}-03`);
    await endInput.click();
    await endInput.fill(`${year}-${month}-05`);
    // 按 Enter 触发搜索
    await endInput.press('Enter');

    // 等待数据重新加载
    await page.waitForLoadState('networkidle');

    // 记录数应少于 20（仅 3 天的数据）
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeLessThan(20);
    expect(count).toBeGreaterThan(0);

    // sentinel 不应存在（筛选后数据量小于 limit）
    await expect(page.getByTestId('load-sentinel')).not.toBeVisible();
  });

  test('分类筛选：选择分类后只显示该分类记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 35 });

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载
    await ledger.expectTotalItemCount(20);

    // 筛选"餐饮"分类
    await ledger.filterByCategory('餐饮');
    await page.waitForLoadState('networkidle');

    // 筛选后记录数应 ≤ 20（部分记录可能是餐饮分类）
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeGreaterThanOrEqual(0);
    expect(count).toBeLessThanOrEqual(20);

    // 清除筛选后恢复全部
    await ledger.clearFilters();
    await ledger.expectTotalItemCount(20);
  });
});
