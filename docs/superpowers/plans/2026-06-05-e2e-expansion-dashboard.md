# Task 3: 仪表盘 — E2E 测试扩展 (7 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐仪表盘的月份切换、数据验证和各区域数据展示测试。

**Files:**
- Modify: `e2e/pages/dashboard.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/dashboard.spec.ts` — 新增 7 个测试用例

**Dependencies:** Task 0 完成

---

## Step 1: 扩展 DashboardPage POM

**Modify:** `e2e/pages/dashboard.page.ts` — 在 `expectForumHotVisible` 方法之后新增：

```typescript
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
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/dashboard.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  test('月份切换', async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    // 记录当前月份文本
    const monthText = await page.getByTestId('current-month').textContent();
    // 切换到上月
    await dashboard.goPrevMonth();
    // 月份文本应变化
    const newMonthText = await page.getByTestId('current-month').textContent();
    expect(newMonthText).not.toBe(monthText);
  });

  test('汇总数据正确性', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 20 });
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectSummaryVisible();
    // 验证收入和支出有数据
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
    // 创建待办
    await page.goto('/#/todo');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试待办');
    await page.getByTestId('submit-btn').click();
    // 回到仪表盘验证
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectUpcomingTodosVisible();
    await dashboard.expectUpcomingTodoCount(1);
  });

  test('愿望动态数据', async ({ authenticated }) => {
    const { page } = authenticated;
    // 创建愿望
    await page.goto('/#/wish');
    await page.getByTestId('add-btn').click();
    await page.getByTestId('title-input').fill('测试愿望');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
    // 回到仪表盘验证
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.expectWishTrendsVisible();
  });

  test('论坛热点数据', async ({ authenticated }) => {
    const { page } = authenticated;
    // 创建话题
    await page.goto('/#/forum');
    await page.getByTestId('create-topic-btn').click();
    await page.getByTestId('topic-title').fill('测试话题');
    await page.getByTestId('topic-content').fill('内容');
    await page.locator('.ant-modal-footer .ant-btn-primary').click();
    // 回到仪表盘验证
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
```

## Step 3: 验证

```bash
cd e2e && npx playwright test dashboard.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/dashboard.page.ts e2e/tests/dashboard.spec.ts
git commit -m "test(e2e): expand dashboard tests — month switch, data verification, widget data"
```
