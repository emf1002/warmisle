// e2e/tests/auth.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { LoginPage } from '../pages/auth/login.page';
import { InitPage } from '../pages/auth/init.page';
import { resetDatabase, initAdmin } from '../fixtures/db.fixture';

test.describe('认证流程', () => {
  test('未登录访问受保护页面应重定向到登录页', async ({ page }) => {
    await resetDatabase(page);
    await initAdmin();
    await page.goto('/#/');
    await expect(page).toHaveURL(/\/#\/login/);
  });

  test('首次访问应跳转到初始化页面', async ({ page }) => {
    await resetDatabase(page);
    await page.goto('/');
    await expect(page).toHaveURL(/\/#\/init/);
  });

  test('初始化系统创建管理员', async ({ page }) => {
    await resetDatabase(page);
    const initPage = new InitPage(page);
    await initPage.goto();
    await initPage.setup('测试管理员', 'admin', 'test123');
    await expect(page).toHaveURL(/\/#\/$/);
  });

  test('登录成功跳转仪表盘', async ({ page }) => {
    await resetDatabase(page);
    await initAdmin();
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login('admin', 'test123');
    await expect(page).toHaveURL(/\/#\/$/);
  });

  test('登录失败显示错误提示', async ({ page }) => {
    await resetDatabase(page);
    await initAdmin();
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login('admin', 'wrongpassword');
    await loginPage.expectOnLoginPage();
  });

  test('登录页视觉回归', async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.screenshot('login-page.png');
  });

  test('移动端登录页视觉回归', { tag: '@mobile' }, async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.screenshot('login-page-mobile.png');
  });

  test('初始化页视觉回归', async ({ page }) => {
    await resetDatabase(page);
    const initPage = new InitPage(page);
    await initPage.goto();
    await initPage.screenshot('init-page.png');
  });

  // === 错误路径 ===

  test('连续登录失败锁定', async ({ page }) => {
    await resetDatabase(page);
    await initAdmin();
    const loginPage = new LoginPage(page);
    for (let i = 0; i < 5; i++) {
      await loginPage.goto();
      await loginPage.login('admin', 'wrongpassword');
    }
    await loginPage.goto();
    await loginPage.login('admin', 'test123');
    await loginPage.expectOnLoginPage();
  });

  test('用户名重复被拒绝', async ({ page }) => {
    await resetDatabase(page);
    const initPage = new InitPage(page);
    await initPage.goto();
    await initPage.setup('管理员', 'admin', 'test123');
    await expect(page).toHaveURL(/\/#\/$/);
    await resetDatabase(page);
    await initPage.goto();
    await initPage.setup('另一个管理员', 'admin', 'test456');
    await initPage.expectOnInitPage();
  });
});
