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

  // === 功能场景 ===

  test('选择头像', async ({ authenticated }) => {
    const { page } = authenticated;
    const profile = new ProfilePage(page);
    await profile.goto();
    await profile.openEditProfile();
    await profile.selectAvatar('PawPrint');
    await profile.submitProfile();
    await expect(page.locator('.ant-message')).toContainText('修改成功');
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
});
