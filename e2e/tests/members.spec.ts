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
});
