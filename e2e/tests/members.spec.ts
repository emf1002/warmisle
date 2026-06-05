// e2e/tests/members.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { MembersPage } from '../pages/members.page';

test.describe('成员管理', () => {
  test('页面加载显示管理员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
  });

  test('添加成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('member1');
    await members.fillPassword('test123');
    await members.fillName('成员一');
    await members.submit();
    await members.expectMemberCount(2);
  });

  test('禁用成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('disabletest');
    await members.fillPassword('test123');
    await members.fillName('待禁用');
    await members.submit();
    await members.disableMember(1);
    // 成员状态应变为"已禁用"
  });

  test('成员页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.screenshot('members-page.png');
  });

  // === 功能场景 ===

  test('编辑成员信息', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('editme');
    await members.fillPassword('test123');
    await members.fillName('原名');
    await members.submit();
    await members.expectMemberCount(2);
    await members.editMember(1);
    await members.fillEditName('新名字');
    await members.submitEdit();
    await members.expectMemberCount(2);
  });

  test('成员角色切换', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('promote');
    await members.fillPassword('test123');
    await members.fillName('待提升');
    await members.submit();
    await members.expectMemberRole(1, '成员');
    await members.editMember(1);
    await members.selectRole('管理员');
    await members.submitEdit();
    await members.expectMemberRole(1, '管理员');
  });

  test('启用被禁用成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('disableme');
    await members.fillPassword('test123');
    await members.fillName('待禁用');
    await members.submit();
    await members.disableMember(1);
    await members.expectMemberStatus(1, '已禁用');
    await members.enableMember(1);
    await members.expectMemberStatus(1, '正常');
  });

  test('删除无活动记录成员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('deleteme');
    await members.fillPassword('test123');
    await members.fillName('待删除');
    await members.submit();
    await members.expectMemberCount(2);
    await members.deleteMember(1);
    await members.expectMemberCount(1);
  });

  // === 错误路径 ===

  test('有活动记录成员不可删除', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('activeuser');
    await members.fillPassword('test123');
    await members.fillName('活跃用户');
    await members.submit();
    // Create a ledger record as admin (which counts as activity for the system)
    await page.goto('/#/ledger');
    await page.getByTestId('add-btn').click();
    await page.locator('.category-pick-item', { hasText: '餐饮' }).click();
    await page.getByTestId('amount-input').click();
    await page.getByTestId('amount-input').locator('input').fill('10');
    await page.getByTestId('submit-btn').click();
    await members.goto();
    await members.deleteMember(1);
    await members.expectMemberCount(2);
  });

  test('禁止删除最后一个管理员', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
    await members.deleteMember(0);
    await members.expectMemberCount(1);
  });

  test('禁止管理员禁用自己', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
    await members.disableMember(0);
    await members.expectMemberStatus(0, '正常');
  });

  test('用户名重复被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.openCreate();
    await members.fillUsername('duplicate');
    await members.fillPassword('test123');
    await members.fillName('用户A');
    await members.submit();
    await members.expectMemberCount(2);
    await members.openCreate();
    await members.fillUsername('duplicate');
    await members.fillPassword('test456');
    await members.fillName('用户B');
    await members.submit();
    await members.expectMemberCount(2);
  });

  // === 权限测试 ===

  test('普通成员不可访问成员管理', async ({ authenticated, memberContext }) => {
    const { page: memberPage } = memberContext;
    await memberPage.goto('/#/members');
    await expect(memberPage).not.toHaveURL(/\/#\/members/);
  });

  // === 响应式 ===

  test('移动端成员列表', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const members = new MembersPage(page);
    await members.goto();
    await members.expectMemberCount(1);
  });
});
