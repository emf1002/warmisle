// e2e/tests/categories.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { CategoriesPage } from '../pages/categories.page';

test.describe('分类管理', () => {
  test('页面加载显示预置分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.expectExpenseCategoryCount(15);
    await categories.expectIncomeCategoryCount(5);
  });

  test('新增支出分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('测试分类');
    await categories.submit();
    // submit() 内部已等待弹窗关闭
    // toHaveCount 会自动重试直到超时，如果列表不刷新会报错
    await categories.expectExpenseCategoryCount(16);
    // 验证新分类名称出现在页面上
    await expect(page.getByTestId('expense-categories')).toContainText('测试分类');
  });

  test('删除自定义分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('待删除');
    await categories.submit();
    await categories.expectExpenseCategoryCount(16);
    // 删除最后一个（刚创建的）
    await categories.deleteCategory(15);
    await categories.expectExpenseCategoryCount(15);
  });

  test('分类页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.screenshot('categories-page.png');
  });

  // === 功能场景 ===

  test('编辑分类', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.editCategory(0);
    await categories.fillEditName('改名后的餐饮');
    await categories.submitEdit();
    await expect(page.getByTestId('expense-categories')).toContainText('改名后的餐饮');
  });

  test('修改分类 type', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('可转换分类');
    await categories.submit();
    await categories.expectExpenseCategoryCount(16);
    await categories.editCategory(15);
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
    await page.goto('/#/ledger');
    await page.getByTestId('add-btn').click();
    await page.locator('.category-pick-item', { hasText: '餐饮' }).click();
    await page.getByTestId('amount-input').click();
    await page.getByTestId('amount-input').locator('input').fill('10');
    await page.getByTestId('submit-btn').click();
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.deleteCategory(0);
    await categories.expectExpenseCategoryCount(15);
  });

  test('同类型重复名称被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const categories = new CategoriesPage(page);
    await categories.goto();
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('测试分类');
    await categories.submit();
    await categories.expectExpenseCategoryCount(16);
    await categories.openCreate();
    await categories.selectType('支出');
    await categories.fillName('测试分类');
    await categories.submit();
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
});
