// e2e/tests/profile.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { ProfilePage } from '../pages/profile.page';

test.describe('个人中心', () => {
  test('页面加载显示个人信息', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.expectProfileName('管理员');
  });

  test('修改昵称', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openEditProfile();
    await profile.fillName('新名字');
    await profile.submitProfile();
    await profile.expectProfileName('新名字');
  });

  test('修改密码', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openChangePassword();
    await profile.fillOldPassword('test123');
    await profile.fillNewPassword('newpass123');
    await profile.submitPassword();
    // 密码修改成功提示
  });

  test('退出登录', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.logout();
    await expect(page).toHaveURL(/\/#\/login/);
  });

  test('个人中心页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.screenshot('profile-page.png');
  });

  test('修改个人信息弹窗视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openEditProfile();
    await profile.screenshotComponent('[data-testid="profile-modal"]', 'profile-edit-modal.png');
  });
});
