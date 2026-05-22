// e2e/pages/ledger.page.ts
import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class LedgerPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/ledger');
  }

  async expectOnLedger() {
    await expect(this.page).toHaveURL(/\/#\/ledger/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.getByTestId('ledger-modal')).toBeVisible();
  }

  async selectCategory(name: string) {
    await this.page.getByTestId('category-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  async fillAmount(amount: string) {
    await this.page.getByTestId('amount-input').click();
    await this.page.getByTestId('amount-input').locator('input').fill(amount);
  }

  async fillNote(note: string) {
    await this.page.getByTestId('note-input').fill(note);
  }

  async submit() {
    await this.page.getByTestId('submit-btn').click();
  }

  async expectRecordCount(count: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items).toHaveCount(count);
  }

  async clickRecord(index: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await items.nth(index).click();
  }

  async deleteRecord() {
    await this.page.getByTestId('delete-btn').click();
    // 确认删除弹窗
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
  }

  async goNextMonth() {
    await this.page.locator('[data-testid="ledger-page"] .month-row button').last().click();
  }

  async goPrevMonth() {
    await this.page.locator('[data-testid="ledger-page"] .month-row button').first().click();
  }
}
