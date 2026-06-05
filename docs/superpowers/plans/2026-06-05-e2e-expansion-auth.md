# Task 9: 认证模块补充 — E2E 测试扩展 (2 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐认证模块的登录锁定和用户名重复测试。

**Files:**
- Modify: `e2e/tests/auth.spec.ts` — 新增 2 个测试用例

**Dependencies:** Task 0 完成

---

## Step 1: 新增测试用例

**Modify:** `e2e/tests/auth.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 错误路径 ===

  test('连续登录失败锁定', async ({ page }) => {
    await resetDatabase();
    await initAdmin();
    const loginPage = new LoginPage(page);
    // 连续 5 次错误登录
    for (let i = 0; i < 5; i++) {
      await loginPage.goto();
      await loginPage.login('admin', 'wrongpassword');
    }
    // 第 6 次应被锁定（即使密码正确）
    await loginPage.goto();
    await loginPage.login('admin', 'test123');
    // 应显示锁定提示或仍在登录页
    await loginPage.expectOnLoginPage();
  });

  test('用户名重复被拒绝', async ({ page }) => {
    await resetDatabase();
    const initPage = new InitPage(page);
    // 初始化第一个管理员
    await initPage.goto();
    await initPage.setup('管理员', 'admin', 'test123');
    await expect(page).toHaveURL(/\/#\/$/);
    // 重新初始化（reset 后）尝试同名用户名
    await resetDatabase();
    await initPage.goto();
    await initPage.setup('另一个管理员', 'admin', 'test456');
    // 应提示用户名已存在
    await initPage.expectOnInitPage();
  });
```

## Step 2: 验证

```bash
cd e2e && npx playwright test auth.spec.ts --project=chromium
```

## Step 3: Commit

```bash
git add e2e/tests/auth.spec.ts
git commit -m "test(e2e): expand auth tests — login lockout and duplicate username"
```
