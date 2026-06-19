import { expect, type Page, type Locator } from '@playwright/test';
import { BasePage } from './base.page';

export class BackupPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.page.goto('/#/backup');
    await this.page.waitForSelector('[data-testid="backup-page"]');
  }

  // --- Config ---
  get configCard(): Locator { return this.page.getByTestId('config-card'); }
  get saveConfigBtn(): Locator { return this.page.getByTestId('save-config-btn'); }
  get authorizeBtn(): Locator { return this.page.getByTestId('authorize-btn'); }

  async fillAppId(appId: string) {
    await this.configCard.locator('input').nth(0).fill(appId);
  }

  async fillAppSecret(secret: string) {
    await this.configCard.locator('input').nth(1).fill(secret);
  }

  async fillRedirectUri(uri: string) {
    await this.configCard.locator('input').nth(2).fill(uri);
  }

  async fillBackupDir(dir: string) {
    await this.configCard.locator('input').nth(3).fill(dir);
  }

  async saveConfig() {
    await this.saveConfigBtn.click();
  }

  async clickAuthorize() {
    await this.authorizeBtn.click();
  }

  // --- Schedule ---
  get scheduleCard(): Locator { return this.page.getByTestId('schedule-card'); }
  get saveScheduleBtn(): Locator { return this.page.getByTestId('save-schedule-btn'); }

  async toggleSchedule(enabled: boolean) {
    const switchEl = this.scheduleCard.locator('.ant-switch');
    const isChecked = await switchEl.getAttribute('aria-checked');
    if ((isChecked === 'true') !== enabled) {
      await switchEl.click();
    }
  }

  async setRetentionDays(days: number) {
    const input = this.scheduleCard.locator('.ant-input-number input');
    await input.fill(String(days));
  }

  async saveSchedule() {
    await this.saveScheduleBtn.click();
  }

  // --- Manual Backup ---
  get triggerBackupBtn(): Locator { return this.page.getByTestId('trigger-backup-btn'); }

  async triggerBackup() {
    await this.triggerBackupBtn.click();
  }

  // --- History ---
  get historyCard(): Locator { return this.page.getByTestId('history-card'); }
  get historyTable(): Locator { return this.page.getByTestId('history-table'); }

  async expectHistoryEmpty() {
    await expect(this.historyTable.locator('.ant-empty')).toBeVisible();
  }

  // --- Cloud Files ---
  get cloudFilesCard(): Locator { return this.page.getByTestId('cloud-files-card'); }

  // --- Restore Modal ---
  get restoreModal(): Locator { return this.page.getByTestId('restore-modal'); }

  // --- Status Tags ---
  async expectStatusTag(text: string) {
    await expect(this.configCard.locator('.ant-tag').filter({ hasText: text })).toBeVisible();
  }
}
