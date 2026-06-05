# Task 4: 愿望清单 — E2E 测试扩展 (13 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐愿望模块的投票、状态流转、评论、删除和权限测试。

**Files:**
- Modify: `e2e/pages/wish.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/wish.spec.ts` — 新增 13 个测试用例

**Dependencies:** Task 0 完成（memberContext fixture 可用）

---

## Step 1: 扩展 WishPage POM

**Modify:** `e2e/pages/wish.page.ts` — 在 `voteWish` 方法之后新增：

```typescript
  /** 取消投票（再次点击投票按钮） */
  async unvoteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('vote-btn').click();
  }

  /** 断言投票人数 */
  async expectVoteCount(index: number, count: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items.nth(index).getByTestId('vote-count')).toContainText(String(count));
  }

  /** 变更愿望状态（管理员操作） */
  async changeWishStatus(index: number, status: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('status-action').click();
    await this.page.locator('.ant-dropdown', { hasText: status }).click();
    // 等待状态更新
    await this.page.waitForTimeout(300);
  }

  /** 断言愿望状态 */
  async expectStatus(index: number, status: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await expect(items.nth(index).getByTestId('status-tag')).toContainText(status);
  }

  /** 创建者放弃愿望 */
  async abandonWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('abandon-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  /** 删除愿望 */
  async deleteWish(index: number) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  /** 评论愿望 */
  async commentOnWish(index: number, text: string) {
    const items = this.page.getByTestId(/^wish-card-/);
    await items.nth(index).getByTestId('comment-btn').click();
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await this.page.waitForTimeout(300);
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/wish.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 投票 ===

  test('投票', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建家庭愿望
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    await wish.openCreate();
    await wish.fillTitle('家庭旅行');
    await wish.submit();
    await wish.expectWishCount(1);
    // member 投票
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await memberWish.switchType('家庭');
    await memberWish.voteWish(0);
    await memberWish.expectVoteCount(0, 1);
  });

  test('取消投票', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建家庭愿望
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    await wish.openCreate();
    await wish.fillTitle('家庭旅行');
    await wish.submit();
    // member 投票再取消
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await memberWish.switchType('家庭');
    await memberWish.voteWish(0);
    await memberWish.expectVoteCount(0, 1);
    await memberWish.unvoteWish(0);
    await memberWish.expectVoteCount(0, 0);
  });

  test('重复投票被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    await wish.openCreate();
    await wish.fillTitle('家庭旅行');
    await wish.submit();
    // 投票
    await wish.voteWish(0);
    await wish.expectVoteCount(0, 1);
    // 再次投票（应被拒绝或无效果）
    await wish.voteWish(0);
    await wish.expectVoteCount(0, 1);
  });

  // === 状态流转 ===

  test('管理员变更愿望状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('新手机');
    await wish.submit();
    await wish.expectStatus(0, '待讨论');
    await wish.changeWishStatus(0, '已同意');
    await wish.expectStatus(0, '已同意');
  });

  test('创建者放弃自己的愿望', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('放弃的愿望');
    await wish.submit();
    await wish.abandonWish(0);
    await wish.expectStatus(0, '已放弃');
  });

  // === 筛选 ===

  test('按状态筛选', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    // 创建两个愿望
    await wish.openCreate();
    await wish.fillTitle('愿望A');
    await wish.submit();
    await wish.openCreate();
    await wish.fillTitle('愿望B');
    await wish.submit();
    // 变更一个状态
    await wish.changeWishStatus(0, '已同意');
    await wish.expectWishCount(2);
    // 筛选"已同意"
    await wish.filterByStatus('已同意');
    await page.waitForTimeout(500);
    await wish.expectWishCount(1);
  });

  // === 评论 ===

  test('评论愿望', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('可评论的愿望');
    await wish.submit();
    await wish.commentOnWish(0, '这个愿望很好！');
    // 评论应显示
    await expect(page.getByTestId('comment-list')).toContainText('这个愿望很好！');
  });

  // === 删除 ===

  test('创建者删除自己的愿望', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('待删除');
    await wish.submit();
    await wish.expectWishCount(1);
    await wish.deleteWish(0);
    await wish.expectWishCount(0);
  });

  test('管理员删除他人的愿望', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 创建愿望
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await memberWish.openCreate();
    await memberWish.fillTitle('member的愿望');
    await memberWish.submit();
    // admin 删除
    const wish = new WishPage(page);
    await wish.goto();
    await wish.deleteWish(0);
    await wish.expectWishCount(0);
  });

  // === 权限 ===

  test('成员不可删除他人的愿望', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建愿望
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('admin的愿望');
    await wish.submit();
    // member 不应看到删除按钮
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await expect(memberPage.getByTestId('delete-btn')).not.toBeVisible();
  });

  test('非创建者不可放弃愿望', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建愿望
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('admin的愿望');
    await wish.submit();
    // member 不应看到放弃按钮
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await expect(memberPage.getByTestId('abandon-btn')).not.toBeVisible();
  });

  // === 错误路径 ===

  test('标题为空被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('');
    await wish.submit();
    await expect(page.getByTestId('wish-modal')).toBeVisible();
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test wish.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/wish.page.ts e2e/tests/wish.spec.ts
git commit -m "test(e2e): expand wish tests — vote, status flow, comment, delete, permissions"
```
