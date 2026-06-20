// e2e/tests/backup.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { BackupPage } from '../pages/backup.page';

test.describe('网盘备份', () => {
  test.describe('管理员访问', () => {
    test('页面正确加载', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      // Should display all major sections
      await expect(backup.configCard).toBeVisible();
      await expect(backup.scheduleCard).toBeVisible();
      await expect(backup.historyCard).toBeVisible();
      await expect(backup.cloudFilesCard).toBeVisible();

      // Should show unconfigured status initially
      await expect(page.locator('.ant-tag').filter({ hasText: '未配置' })).toBeVisible();

      // Save config button should be visible
      await expect(backup.saveConfigBtn).toBeVisible();

      // Trigger button should be present but may be disabled (no auth)
      await expect(backup.triggerBackupBtn).toBeVisible();

      // Schedule save button should be visible
      await expect(backup.saveScheduleBtn).toBeVisible();
    });

    test('保存云盘配置后刷新验证持久化', async ({ authenticated }) => {
      const { page } = authenticated;

      const backup = new BackupPage(page);
      await backup.goto();

      // Fill and save
      await backup.fillAppId('78bb8cec0642495d8c71f1c8aaeee4b6');
      await backup.fillRedirectUri('http://127.0.0.1:8080/api/backup/callback');
      await backup.saveConfig();
      await page.waitForLoadState('networkidle');

      // Reload and check
      await page.reload();
      await page.waitForSelector('[data-testid="backup-page"]');
      await page.waitForLoadState('networkidle');

      const appIdInput = backup.configCard.locator('input').nth(0);
      await expect(appIdInput).toHaveValue('78bb8cec0642495d8c71f1c8aaeee4b6');
    });

    test('定时备份配置开关默认关闭', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      const switchEl = backup.scheduleCard.locator('.ant-switch');
      await expect(switchEl).toHaveAttribute('aria-checked', 'false');
    });

    test('开启定时备份配置', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      // Toggle schedule on
      await backup.toggleSchedule(true);
      await backup.setRetentionDays(60);
      await backup.saveSchedule();
      await page.waitForLoadState('networkidle');

      // Verify switch is on immediately after save
      let switchEl = backup.scheduleCard.locator('.ant-switch');
      await expect(switchEl).toHaveAttribute('aria-checked', 'true');

      // Toggle off
      await backup.toggleSchedule(false);
      await backup.saveSchedule();
      await page.waitForLoadState('networkidle');

      // Reload and verify persistence
      await page.reload();
      await page.waitForSelector('[data-testid="backup-page"]');
      await page.waitForLoadState('networkidle');

      switchEl = backup.scheduleCard.locator('.ant-switch');
      await expect(switchEl).toHaveAttribute('aria-checked', 'false');
    });

    test('获取授权链接', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      // First save config to enable authorize button
      await backup.fillAppId('78bb8cec0642495d8c71f1c8aaeee4b6');
      await backup.fillRedirectUri('http://127.0.0.1:8080/api/backup/callback');
      await backup.saveConfig();
      await page.waitForLoadState('networkidle');

      // Authorize button should now be enabled
      await expect(backup.authorizeBtn).toBeEnabled();
    });

    test('未授权时触发备份应失败', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      // Save config but don't authorize
      await backup.fillAppId('78bb8cec0642495d8c71f1c8aaeee4b6');
      await backup.fillRedirectUri('http://127.0.0.1:8080/api/backup/callback');
      await backup.saveConfig();
      await page.waitForLoadState('networkidle');

      // Trigger button may be disabled since not authorized
      const isDisabled = await backup.triggerBackupBtn.getAttribute('disabled');
      if (isDisabled === null || isDisabled === 'false') {
        await backup.triggerBackup();
        // Should show an error message (not authorized)
        await expect(page.getByTestId('error-toast')).toBeVisible({ timeout: 3000 });
      } else {
        // Button is disabled as expected
        await expect(backup.triggerBackupBtn).toBeDisabled();
      }
    });

    test('备份历史初始为空', async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      await backup.expectHistoryEmpty();
    });
  });

  test.describe('权限控制', () => {
    test('普通成员无法访问备份页面', async ({ memberContext }) => {
      const { page } = memberContext;
      await page.goto('/#/backup');

      // Non-admin should be redirected away from backup page
      await expect(page).not.toHaveURL(/\/#\/backup/, { timeout: 5000 });
      // Should see dashboard (redirect destination for non-admin)
      await expect(page.getByTestId('dashboard-page')).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe('响应式', () => {
    test('移动端备份页面正常加载', { tag: '@mobile' }, async ({ authenticated }) => {
      const { page } = authenticated;
      const backup = new BackupPage(page);
      await backup.goto();

      await expect(backup.configCard).toBeVisible();
      await expect(backup.scheduleCard).toBeVisible();
    });
  });
});
