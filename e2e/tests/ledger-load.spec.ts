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
    // 等待网络请求和 DOM 更新
    await page.waitForLoadState('networkidle', { timeout: 5000 }).catch(() => {});
    await page.waitForTimeout(500);

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

  test('日期范围筛选：选择子范围后只显示匹配记录', async ({ authenticated }) => {
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

    // 点击 RangePicker 的输入框打开日历
    await page.locator('.ant-picker-input input').first().click();
    // 等待日历面板出现
    await expect(page.locator('.ant-picker-dropdown:visible')).toBeVisible();

    // 使用 Ant Design 内置的 presets 点击"本月"来快速选择一个范围
    // 然后验证筛选后记录数减少
    // 但更简单的方式：直接在日历中点击两个日期
    // 选择起始日期：当月 3 日
    const cells = page.locator('.ant-picker-cell-inner');
    // 找到包含 "3" 的单元格（月视图中的日期数字）
    const day3 = cells.filter({ hasText: /^3$/ }).first();
    await day3.click();
    // 选择结束日期：当月 5 日
    const day5 = cells.filter({ hasText: /^5$/ }).first();
    await day5.click();

    // 等待日历关闭和数据重新加载
    await page.waitForTimeout(500);
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

    // 等待数据重新加载（debounce 300ms + 网络请求）
    await page.waitForTimeout(500);
    await page.waitForLoadState('networkidle');

    // 记录数应少于 20（只显示餐饮分类）
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeLessThan(20);
    expect(count).toBeGreaterThan(0);

    // 验证所有可见记录都是餐饮分类
    const categoryNames = page.locator('.item-cat-name');
    const catCount = await categoryNames.count();
    for (let i = 0; i < catCount; i++) {
      await expect(categoryNames.nth(i)).toHaveText('餐饮');
    }

    // 清除筛选后恢复全部
    await ledger.clearFilters();
    await page.waitForTimeout(500);
    await page.waitForLoadState('networkidle');
    await ledger.expectTotalItemCount(20);
  });
});
