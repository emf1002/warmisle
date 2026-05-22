// e2e/tests/wish.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { WishPage } from '../pages/wish.page';

test.describe('愿望清单', () => {
  test('页面加载显示空状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('新建愿望', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('新手机');
    await wish.fillDescription('iPhone 16');
    await wish.selectCategory('物品');
    await wish.fillAmount('6999');
    await wish.submit();
    await wish.expectWishCount(1);
  });

  test('切换个人/家庭愿望', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    // 应切换到家庭愿望视图
  });

  test('愿望页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.screenshot('wish-empty.png');
  });
});
