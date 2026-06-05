# Task 6: 成员管理 — E2E 测试扩展 (11 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐成员管理的编辑、启用、删除保护和权限边界测试。

**Files:**
- Modify: `e2e/pages/members.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/members.spec.ts` — 新增 11 个测试用例

**Dependencies:** Task 0 完成（memberContext fixture 可用）

---

## Step 1: 扩展 MembersPage POM

**Modify:** `e2e/pages/members.page.ts` — 在 `deleteMember` 方法之后新增：

```typescript
  /** 断言成员状态标签 */
  async expectMemberStatus(index: number, status: string) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await expect(rows.nth(index).getByTestId('status-tag')).toContainText(status);
  }

  /** 断言成员角色标签 */
  async expectMemberRole(index: number, role: string) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await expect(rows.nth(index).getByTestId('role-tag')).toContainText(role);
  }

  /** 断言删除按钮被禁用（有活动记录的成员） */
  async expectDeleteDisabled(index: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await expect(rows.nth(index).getByTestId('delete-btn')).toBeDisabled();
  }

  /** 填写成员姓名（编辑弹窗中） */
  async fillEditName(name: string) {
    await this.page.getByTestId('member-modal').getByTestId('name-input').fill(name);
  }

  /** 提交编辑弹窗 */
  async submitEdit() {
    await this.page.getByTestId('member-modal').locator('.ant-modal-footer .ant-btn-primary').click();
    await expect(this.page.getByTestId('member-modal')).not.toBeVisible();
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/members.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 功能场景 ===

  test('编辑成员信息', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 添加成员
    await members.openCreate();
    await members.fillUsername('editme');
    await members.fillPassword('test123');
    await members.fillName('原名');
    await members.submit();
    await members.expectMemberCount(2);
    // 编辑
    await members.editMember(1);
    await members.fillEditName('新名字');
    await members.submitEdit();
    await members.expectMemberCount(2);
  });

  test('成员角色切换', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 添加 member
    await members.openCreate();
    await members.fillUsername('promote');
    await members.fillPassword('test123');
    await members.fillName('待提升');
    await members.submit();
    await members.expectMemberRole(1, '成员');
    // 提升为管理员
    await members.editMember(1);
    await members.selectRole('管理员');
    await members.submitEdit();
    await members.expectMemberRole(1, '管理员');
  });

  test('启用被禁用成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 添加并禁用
    await members.openCreate();
    await members.fillUsername('disableme');
    await members.fillPassword('test123');
    await members.fillName('待禁用');
    await members.submit();
    await members.disableMember(1);
    await members.expectMemberStatus(1, '已禁用');
    // 启用
    await members.enableMember(1);
    await members.expectMemberStatus(1, '正常');
  });

  test('删除无活动记录成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 添加成员
    await members.openCreate();
    await members.fillUsername('deleteme');
    await members.fillPassword('test123');
    await members.fillName('待删除');
    await members.submit();
    await members.expectMemberCount(2);
    // 删除
    await members.deleteMember(1);
    await members.expectMemberCount(1);
  });

  // === 错误路径 ===

  test('有活动记录成员不可删除', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 添加成员
    await members.openCreate();
    await members.fillUsername('activeuser');
    await members.fillPassword('test123');
    await members.fillName('活跃用户');
    await members.submit();
    // 用该成员创建一条记账记录（通过 API 或 UI）
    await page.goto('/#/ledger');
    await page.getByTestId('add-btn').click();
    await page.locator('.category-pick-item', { hasText: '餐饮' }).click();
    await page.getByTestId('amount-input').click();
    await page.getByTestId('amount-input').locator('input').fill('10');
    await page.getByTestId('submit-btn').click();
    // 回到成员管理，尝试删除
    await members.goto();
    // 有活动记录的成员删除应被拒绝
    await members.deleteMember(1);
    await members.expectMemberCount(2);
  });

  test('禁止删除最后一个管理员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
    // 尝试删除唯一的管理员
    await members.deleteMember(0);
    // 应被拒绝
    await members.expectMemberCount(1);
  });

  test('禁止管理员禁用自己', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
    // 尝试禁用自己
    await members.disableMember(0);
    // 应被拒绝，状态不变
    await members.expectMemberStatus(0, '正常');
  });

  test('用户名重复被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    // 创建第一个用户
    await members.openCreate();
    await members.fillUsername('duplicate');
    await members.fillPassword('test123');
    await members.fillName('用户A');
    await members.submit();
    await members.expectMemberCount(2);
    // 创建同名用户
    await members.openCreate();
    await members.fillUsername('duplicate');
    await members.fillPassword('test456');
    await members.fillName('用户B');
    await members.submit();
    // 应被拒绝
    await members.expectMemberCount(2);
  });

  // === 权限测试 ===

  test('普通成员不可访问成员管理', async ({ authenticated, memberContext }) => {
    const { page: memberPage } = memberContext;
    await memberPage.goto('/#/members');
    // 应被重定向或无权限
    await expect(memberPage).not.toHaveURL(/\/#\/members/);
  });

  // === 响应式 ===

  test('移动端成员列表', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test members.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/members.page.ts e2e/tests/members.spec.ts
git commit -m "test(e2e): expand members tests — edit, enable, delete protection, permissions"
```
