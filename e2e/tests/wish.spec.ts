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

  // === 投票 ===

  test('投票', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    await wish.openCreate();
    await wish.fillTitle('家庭旅行');
    await wish.submit();
    await wish.expectWishCount(1);
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await memberWish.switchType('家庭');
    await memberWish.voteWish(0);
    await memberWish.expectVoteCount(0, 1);
  });

  test('取消投票', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.switchType('家庭');
    await wish.openCreate();
    await wish.fillTitle('家庭旅行');
    await wish.submit();
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
    await wish.voteWish(0);
    await wish.expectVoteCount(0, 1);
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
    await wish.expectStatus(0, '待定');
    await wish.changeWishStatus(0, '标记为同意');
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
    await wish.openCreate();
    await wish.fillTitle('愿望A');
    await wish.submit();
    await wish.openCreate();
    await wish.fillTitle('愿望B');
    await wish.submit();
    await wish.changeWishStatus(0, '标记为同意');
    await wish.expectWishCount(2);
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
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await memberWish.openCreate();
    await memberWish.fillTitle('member的愿望');
    await memberWish.submit();
    const wish = new WishPage(page);
    await wish.goto();
    await wish.deleteWish(0);
    await wish.expectWishCount(0);
  });

  // === 权限 ===

  test('成员不可删除他人的愿望', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('admin的愿望');
    await wish.submit();
    const memberWish = new WishPage(memberPage);
    await memberWish.goto();
    await expect(memberPage.getByTestId('delete-btn')).not.toBeVisible();
  });

  test('非创建者不可放弃愿望', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const wish = new WishPage(page);
    await wish.goto();
    await wish.openCreate();
    await wish.fillTitle('admin的愿望');
    await wish.submit();
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
    await expect(page.locator('.ant-modal-wrap:visible')).toBeVisible();
  });
});
