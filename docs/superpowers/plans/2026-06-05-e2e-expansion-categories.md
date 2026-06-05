# Task 7: 分类管理 — E2E 测试扩展 (6 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐分类管理的编辑、type 变更、删除保护和错误路径测试。

**Files:**
- Modify: `e2e/pages/categories.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/categories.spec.ts` — 新增 6 个测试用例

**Dependencies:** Task 0 完成

---

## Step 1: 扩展 CategoriesPage POM

**Modify:** `e2e/pages/categories.page.ts` — 在 `deleteCategory` 方法之后新增：

```typescript
  /** 在编辑弹窗中修改分类 type */
  async changeCategoryType(type: string) {
    await this.page.getByTestId('category-modal').getByTestId('type-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: type }).click();
  }

  /** 在编辑弹窗中填写分类名称 */
  async fillEditName(name: string) {
    await this.page.getByTestId('category-modal').getByTestId('name-input').fill(name);
  }

  /** 提交编辑弹窗 */
  async submitEdit() {
    await this.page.getByTestId('category-modal').locator('.ant-modal-footer .ant-btn-primary').click();
    await expect(this.page.getByTestId('category-modal')).not.toBeVisible();
  }

  /** 断言删除按钮被禁用（有关联记录的分类） */
  async expectDeleteDisabled(index: number) {
    const cards = this.page.getByTestId(/^category-card-/);
    await expect(cards.nth(index).getByTestId('delete-btn')).toBeDisabled();
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/categories.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 功能场景 ===

  test('编辑分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    // 编辑第一个预置分类
    await categories.editCategory(0);
    await categories.fillEditName('改名后的餐饮');
    await categories.submitEdit();
    await expect(page.getByTestId('expense-categories')).toContainText('改名后的餐饮');
  });

  test('修改分类 type', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    // 创建一个支出分类
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('可转换分类');
    await categories.submit();
    await categories.expectExpenseCategoryCount(16);
    // 编辑为收入分类
    const newIndex = 15; // 刚创建的
    await categories.editCategory(newIndex);
    await categories.changeCategoryType('收入');
    await categories.submitEdit();
    await categories.expectExpenseCategoryCount(15);
    await categories.expectIncomeCategoryCount(6);
  });

  test('添加收入分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.openCreate();
    await categories.selectType('收入');
    await categories.fillName('新收入分类');
    await categories.submit();
    await categories.expectIncomeCategoryCount(6);
    await expect(page.getByTestId('income-categories')).toContainText('新收入分类');
  });

  // === 错误路径 ===

  test('有关联记录的分类不可删除', async ({ authenticated }) => {
    const { page } = authenticated;
    // 先创建一条记账记录使用"餐饮"分类
    await page.goto('/#/ledger');
    await page.getByTestId('add-btn').click();
    await page.locator('.category-pick-item', { hasText: '餐饮' }).click();
    await page.getByTestId('amount-input').click();
    await page.getByTestId('amount-input').locator('input').fill('10');
    await page.getByTestId('submit-btn').click();
    // 回到分类管理，尝试删除"餐饮"
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.deleteCategory(0);
    // 应被拒绝
    await categories.expectExpenseCategoryCount(15);
  });

  test('同类型重复名称被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    // 创建"测试分类"
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('测试分类');
    await categories.submit();
    await categories.expectExpenseCategoryCount(16);
    // 再次创建同名
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('测试分类');
    await categories.submit();
    // 应被拒绝
    await categories.expectExpenseCategoryCount(16);
  });

  // === 响应式 ===

  test('移动端分类列表', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.expectExpenseCategoryCount(15);
    await categories.expectIncomeCategoryCount(5);
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test categories.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/categories.page.ts e2e/tests/categories.spec.ts
git commit -m "test(e2e): expand categories tests — edit, type change, delete protection"
```
