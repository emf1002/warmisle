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
});
