# Task 8: 个人中心 — E2E 测试扩展 (6 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐个人中心的头像选择和密码校验边界测试。

**Files:**
- Modify: `e2e/pages/profile.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/profile.spec.ts` — 新增 6 个测试用例

**Dependencies:** Task 0 完成

---

## Step 1: 扩展 ProfilePage POM

**Modify:** `e2e/pages/profile.page.ts` — 在 `expectProfileName` 方法之后新增：

```typescript
  /** 从 emoji 列表选择头像 */
  async selectAvatar(avatar: string) {
    await this.page.getByTestId('avatar-picker').click();
    await this.page.getByTestId('avatar-grid').locator(`text=${avatar}`).click();
  }

  /** 填写确认密码（修改密码弹窗中） */
  async fillConfirmPassword(password: string) {
    await this.page.getByTestId('confirm-pwd-input').fill(password);
  }

  /** 断言旧密码错误提示 */
  async expectOldPasswordError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('原密码错误');
  }

  /** 断言密码不一致提示 */
  async expectPasswordMismatchError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('两次输入的密码不一致');
  }

  /** 断言新旧密码相同提示 */
  async expectSamePasswordError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('新密码不能与旧密码相同');
  }

  /** 断言姓名超长提示 */
  async expectNameTooLongError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('请输入1-20字符');
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/profile.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 功能场景 ===

  test('选择头像', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openEditProfile();
    await profile.selectAvatar('🐶');
    await profile.submitProfile();
    // 头像应更新
    await expect(page.locator('.profile-header')).toContainText('🐶');
  });

  // === 错误路径 ===

  test('旧密码错误', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openChangePassword();
    await profile.fillOldPassword('wrongpassword');
    await profile.fillNewPassword('newpass123');
    await profile.fillConfirmPassword('newpass123');
    await profile.submitPassword();
    await profile.expectOldPasswordError();
  });

  test('新密码不一致', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openChangePassword();
    await profile.fillOldPassword('test123');
    await profile.fillNewPassword('newpass123');
    await profile.fillConfirmPassword('different456');
    await profile.submitPassword();
    await profile.expectPasswordMismatchError();
  });

  test('新密码与旧密码相同', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openChangePassword();
    await profile.fillOldPassword('test123');
    await profile.fillNewPassword('test123');
    await profile.fillConfirmPassword('test123');
    await profile.submitPassword();
    await profile.expectSamePasswordError();
  });

  test('姓名超长被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openEditProfile();
    await profile.fillName('这是一个超过二十个字符的很长很长很长的姓名');
    await profile.submitProfile();
    await profile.expectNameTooLongError();
  });

  // === 响应式 ===

  test('移动端个人中心', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.expectProfileName('管理员');
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test profile.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/profile.page.ts e2e/tests/profile.spec.ts
git commit -m "test(e2e): expand profile tests — avatar, password validation edge cases"
```
