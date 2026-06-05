# Task 10: 移动端视觉回归 — E2E 测试 (4 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 为关键页面建立移动端视觉回归基线。

**Files:**
- Modify: `e2e/tests/auth.spec.ts` — 新增移动端登录页截图
- Modify: `e2e/tests/dashboard.spec.ts` — 新增移动端仪表盘截图
- Modify: `e2e/tests/ledger.spec.ts` — 新增移动端记账列表截图
- Modify: `e2e/tests/todo.spec.ts` — 新增移动端待办截图（验证底栏导航）

**Dependencies:** Task 0 中 Playwright mobile project 配置完成

---

## Step 1: 新增移动端视觉回归测试

各 spec 文件中新增带 `@mobile` tag 的视觉回归测试。

### auth.spec.ts

```typescript
  test('移动端登录页视觉回归', { tag: '@mobile' }, async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.screenshot('login-page-mobile.png');
  });
```

### dashboard.spec.ts

```typescript
  test('移动端仪表盘视觉回归', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const dashboard = new DashboardPage(page);
    await dashboard.goto();
    await dashboard.screenshot('dashboard-mobile.png');
  });
```

### ledger.spec.ts

```typescript
  test('移动端记账列表视觉回归', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({ count: 10 });
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.screenshot('ledger-list-mobile.png');
  });
```

### todo.spec.ts

```typescript
  test('移动端底栏导航视觉回归', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    // 验证底栏 TabBar 可见
    await expect(page.getByTestId('mobile-tabbar')).toBeVisible();
    await todo.screenshot('todo-mobile.png');
  });
```

## Step 2: 更新视觉回归基线

```bash
cd e2e && npx playwright test --project=mobile --update-snapshots
```

首次运行自动生成基线截图。检查 `e2e/__snapshots__/` 下的 mobile 截图，确认无误后提交。

## Step 3: 验证

```bash
cd e2e && npx playwright test --project=mobile
```

## Step 4: Commit

```bash
git add e2e/__snapshots__/ e2e/tests/auth.spec.ts e2e/tests/dashboard.spec.ts e2e/tests/ledger.spec.ts e2e/tests/todo.spec.ts
git commit -m "test(e2e): add mobile visual regression baselines for key pages"
```
