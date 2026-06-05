# Task 1: 记账本 — E2E 测试扩展 (14 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐记账本模块的功能场景、权限测试和错误路径测试。

**Files:**
- Modify: `e2e/pages/ledger.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/ledger.spec.ts` — 新增 14 个测试用例

**Dependencies:** Task 0 完成（memberContext fixture 可用）

---

## Step 1: 扩展 LedgerPage POM

**Modify:** `e2e/pages/ledger.page.ts` — 在 `expectDailyTotal` 方法之后、类的结束 `}` 之前新增以下方法：

```typescript
  /** 编辑第 n 条记录：点击记录 → 点击编辑按钮 */
  async editRecord(index: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await items.nth(index).click();
    await this.page.getByTestId('edit-btn').click();
    await expect(this.page.getByTestId('ledger-modal')).toBeVisible();
  }

  /** 断言第 n 条记录显示指定的创建者名称 */
  async expectRecordCreator(index: number, name: string) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items.nth(index).getByTestId('creator-name')).toContainText(name);
  }

  /** 切换到上个月 */
  async goPrevMonth() {
    await this.page.getByTestId('month-prev').click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 切换到下个月 */
  async goNextMonth() {
    await this.page.getByTestId('month-next').click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 断言当前月份显示（如"2026年6月"） */
  async expectMonthText(text: string) {
    await expect(this.page.getByTestId('current-month')).toContainText(text);
  }

  /** 断言金额输入校验错误 */
  async expectAmountError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('请输入正数金额');
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/ledger.spec.ts` — 在现有 `test.describe('记账本', () => {` 块内，最后一个 `});` 之前新增：

```typescript
  // === 功能场景 ===

  test('编辑自己的记录', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    // 创建记录
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('35.5');
    await ledger.fillNote('午饭');
    await ledger.submit();
    await ledger.expectRecordCount(1);
    // 编辑记录
    await ledger.editRecord(0);
    await ledger.fillAmount('50');
    await ledger.fillNote('晚饭');
    await ledger.submit();
    // 验证更新
    await ledger.expectRecordCount(1);
  });

  test('记录者显示', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('10');
    await ledger.submit();
    await ledger.expectRecordCreator(0, '管理员');
  });

  test('月份切换', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 10 });
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.expectTotalItemCount(10);
    // 切换到上个月
    await ledger.goPrevMonth();
    // 上个月应无记录（seed 默认当月数据）
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('按记录者筛选', async ({ authenticated, memberContext }) => {
    const { page, token } = authenticated;
    const { page: memberPage, token: memberToken } = memberContext;
    // admin 创建一条记录
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('10');
    await ledger.submit();
    // member 创建一条记录（通过 API）
    const memberLedger = new LedgerPage(memberPage);
    await memberLedger.goto();
    await memberLedger.openCreate();
    await memberLedger.pickCategory('交通');
    await memberLedger.fillAmount('20');
    await memberLedger.submit();
    // admin 按记录者筛选
    await ledger.goto();
    await ledger.expectTotalItemCount(2);
    await ledger.filterByCreator('管理员');
    await page.waitForTimeout(500);
    await ledger.expectTotalItemCount(1);
  });

  test('清除筛选恢复全部', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 10 });
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.expectTotalItemCount(10);
    // 分类筛选
    await ledger.filterByCategory('餐饮');
    await page.waitForTimeout(500);
    const filteredCount = await page.getByTestId(/^ledger-item-/).count();
    expect(filteredCount).toBeLessThan(10);
    // 清除筛选
    await ledger.clearFilters();
    await page.waitForTimeout(500);
    await ledger.expectTotalItemCount(10);
  });

  test('日期分组和日小计', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 30 });
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.expectTotalItemCount(20);
    // 应有多个日期分组
    const groupCount = await page.getByTestId('date-group').count();
    expect(groupCount).toBeGreaterThan(1);
    // 每个分组应有日小计
    for (let i = 0; i < groupCount; i++) {
      await expect(page.getByTestId('daily-total').nth(i)).toContainText('¥');
    }
  });

  test('汇总数据正确性', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 20 });
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.expectTotalItemCount(20);
    // 验证汇总栏有数据
    await ledger.expectSummary({
      income: '+¥',
      expense: '-¥',
    });
  });

  // === 权限测试 ===

  test('管理员可编辑成员的记录', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 创建记录
    const memberLedger = new LedgerPage(memberPage);
    await memberLedger.goto();
    await memberLedger.openCreate();
    await memberLedger.pickCategory('餐饮');
    await memberLedger.fillAmount('10');
    await memberLedger.submit();
    // admin 编辑 member 的记录
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.editRecord(0);
    await ledger.fillAmount('99');
    await ledger.submit();
    await ledger.expectRecordCount(1);
  });

  test('成员不可编辑他人的记录', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建记录
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('10');
    await ledger.submit();
    // member 访问记账列表，不应看到编辑入口
    const memberLedger = new LedgerPage(memberPage);
    await memberLedger.goto();
    await memberLedger.expectRecordCount(1);
    // 点击记录后不应有编辑按钮
    await memberLedger.clickRecord(0);
    await expect(memberPage.getByTestId('edit-btn')).not.toBeVisible();
  });

  test('管理员可删除成员的记录', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 创建记录
    const memberLedger = new LedgerPage(memberPage);
    await memberLedger.goto();
    await memberLedger.openCreate();
    await memberLedger.pickCategory('餐饮');
    await memberLedger.fillAmount('10');
    await memberLedger.submit();
    // admin 删除 member 的记录
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.clickRecord(0);
    await ledger.deleteRecord();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('成员不可删除他人的记录', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建记录
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('10');
    await ledger.submit();
    // member 不应看到删除按钮
    const memberLedger = new LedgerPage(memberPage);
    await memberLedger.goto();
    await memberLedger.clickRecord(0);
    await expect(memberPage.getByTestId('delete-btn')).not.toBeVisible();
  });

  // === 错误路径 ===

  test('金额为 0 被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('0');
    await ledger.submit();
    // 应显示错误提示或提交被阻止
    await expect(page.getByTestId('ledger-modal')).toBeVisible();
  });

  // === 响应式 ===

  test('移动端记账列表', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.pickCategory('餐饮');
    await ledger.fillAmount('35.5');
    await ledger.submit();
    await ledger.expectRecordCount(1);
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test ledger.spec.ts --project=chromium
```

预期：全部通过（原有 6 个 + 新增 14 个 = 20 个）

## Step 4: Commit

```bash
git add e2e/pages/ledger.page.ts e2e/tests/ledger.spec.ts
git commit -m "test(e2e): expand ledger tests — edit, permissions, filtering, error paths"
```
