// e2e/tests/ledger-load.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { LedgerPage } from '../pages/ledger.page';

test.describe('记账本 — 负载场景', () => {
  test('无限滚动：首次加载 20 条，滚动后加载全部 35 条', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
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

    // 等待全部 35 条加载完成
    await ledger.expectTotalItemCount(35);

    // sentinel 消失
    await expect(page.getByTestId('load-sentinel')).not.toBeVisible();

    // 7 天 = 7 个日期分组
    await ledger.expectDateGroupCount(7);
  });

  test('统计汇总：收入/支出/结余与 seed 数据一致', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 35 });
    const { summary } = result.data;

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载完成
    await ledger.expectTotalItemCount(20);

    // 断言汇总值（页面首次加载即包含 summary，基于全量数据计算）
    const expectedIncome = (summary.income / 100).toFixed(2);
    const expectedExpense = (summary.expense / 100).toFixed(2);
    const expectedBalance = (summary.balance / 100).toFixed(2);

    await ledger.expectSummary({
      income: expectedIncome,
      expense: expectedExpense,
    });

    // 结余符号取决于正负
    if (summary.balance >= 0) {
      await ledger.expectSummary({ balance: `+${expectedBalance}` });
    } else {
      await ledger.expectSummary({ balance: `-${Math.abs(Number(expectedBalance)).toFixed(2)}` });
    }

    // 滚动加载全部后，汇总值不变
    await ledger.scrollToLoadMore();
    await ledger.expectTotalItemCount(35);
    await ledger.expectSummary({
      income: expectedIncome,
      expense: expectedExpense,
    });
  });

  test('日期范围筛选：选择子范围后只显示匹配记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({
      count: 35,
      startDate: '2026-06-01',
      endDate: '2026-06-07',
    });

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待全量数据加载
    await ledger.scrollToLoadMore();
    await ledger.expectTotalItemCount(35);

    // 使用 RangePicker 选择 6 月 3 日 ~ 6 月 5 日（3 天）
    const rangePicker = page.getByTestId('date-range-picker');
    await rangePicker.click();

    // 点击 6 月 3 日
    await page.locator('.ant-picker-cell[title="2026-06-03"]').click();
    // 点击 6 月 5 日
    await page.locator('.ant-picker-cell[title="2026-06-05"]').click();

    // 等待数据重新加载
    await page.waitForLoadState('networkidle');

    // 记录数应少于 35
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeLessThan(35);
    expect(count).toBeGreaterThan(0);

    // sentinel 不应存在（筛选后数据量小于 limit）
    await expect(page.getByTestId('load-sentinel')).not.toBeVisible();
  });

  test('分类筛选：选择分类后只显示该分类记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 35 });
    const fullSummary = result.data.summary;

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载
    await ledger.expectTotalItemCount(20);

    // 筛选"餐饮"分类
    await ledger.filterByCategory('餐饮');

    // 等待数据重新加载
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
    await page.waitForLoadState('networkidle');
    await ledger.expectTotalItemCount(20);
  });
});
