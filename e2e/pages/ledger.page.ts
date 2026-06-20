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
    await this.expectModalVisible();
  }

  /** 在弹窗网格选择器中点击分类（按文本匹配，自动切换 tab） */
  async pickCategory(name: string) {
    // Try expense tab first
    let item = this.page.locator('.category-pick-item', { hasText: name });
    if (await item.count() === 0) {
      // Switch to income tab
      await this.page.locator('.ant-modal:visible').locator('.ant-tabs-tab', { hasText: '收入' }).click();
    }
    await this.page.locator('.category-pick-item', { hasText: name }).click();
  }

  async fillAmount(amount: string) {
    const input = this.page.locator('.ant-modal:visible').getByRole('spinbutton');
    await input.click();
    await input.fill(amount);
  }

  async fillNote(note: string) {
    await this.page.getByTestId('note-input').fill(note);
  }

  async selectMember(name: string) {
    await this.page.getByTestId('member-select').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: name }).click();
    // Close the select dropdown by clicking the modal title area
    await this.page.locator('.ant-modal-header').click();
    await this.page.waitForTimeout(200);
  }

  async submit() {
    await this.page.getByTestId('submit-btn').click();
    await this.expectModalHidden();
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
    await this.confirmModal();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 通过滤选栏选择分类 */
  async filterByCategory(name: string) {
    await this.page.getByTestId('filter-category').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: name }).click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 通过滤选栏选择创建者 */
  async filterByCreator(name: string) {
    await this.page.getByTestId('filter-creator').click();
    await this.page.locator('.ant-select-item-option').filter({ hasText: name }).click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 滚动到 sentinel 触发无限加载 */
  async scrollToLoadMore() {
    const sentinel = this.page.getByTestId('load-sentinel');
    await sentinel.scrollIntoViewIfNeeded();
    // Wait for new data to load
    await this.page.waitForFunction(
      () => document.querySelectorAll('[data-testid^="ledger-item-"]').length > 20,
      { timeout: 5000 }
    ).catch(() => {});
  }

  /** 断言汇总栏文本 */
  async expectSummary(values: { income?: string; expense?: string; balance?: string }) {
    if (values.income !== undefined) {
      await expect(this.page.getByTestId('summary-income')).toContainText(values.income);
    }
    if (values.expense !== undefined) {
      await expect(this.page.getByTestId('summary-expense')).toContainText(values.expense);
    }
    if (values.balance !== undefined) {
      await expect(this.page.getByTestId('summary-balance')).toContainText(values.balance);
    }
  }

  /** 断言日期分组数量 */
  async expectDateGroupCount(count: number) {
    await expect(this.page.getByTestId('date-group')).toHaveCount(count);
  }

  /** 断言所有记录总数 */
  async expectTotalItemCount(count: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items).toHaveCount(count);
  }

  /** 断言第 n 个日期组的每日小计包含指定文本 */
  async expectDailyTotal(index: number, text: string) {
    await expect(this.page.getByTestId('daily-total').nth(index)).toContainText(text);
  }

  /** 编辑第 n 条记录：点击记录直接打开编辑弹窗 */
  async editRecord(index: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await items.nth(index).click();
    await this.expectModalVisible();
  }

  /** 断言第 n 条记录显示指定的金额文本（如 "+35.50" 或 "-35.50"） */
  async expectRecordAmount(index: number, amountText: string) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items.nth(index).locator('.item-amount')).toContainText(amountText);
  }

  /** 断言第 n 条记录显示指定的创建者名称 */
  async expectRecordCreator(index: number, name: string) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items.nth(index).getByTestId('creator-name')).toContainText(name);
  }

  /** 切换到上个月 */
  async goPrevMonth() {
    await this.page.getByPlaceholder('请选择月份').click();
    await this.page.locator('.ant-picker-dropdown:visible').waitFor();
    const prevYearBtn = this.page.locator('.ant-picker-dropdown:visible .ant-picker-super-prev-btn');
    if (await prevYearBtn.isVisible()) {
      await prevYearBtn.click();
    }
    const cells = this.page.locator('.ant-picker-dropdown:visible .ant-picker-cell-inner');
    const count = await cells.count();
    if (count > 0) {
      await cells.nth(count - 1).click();
    }
    await this.page.waitForLoadState('networkidle');
  }

  /** 切换到下个月 */
  async goNextMonth() {
    await this.page.getByPlaceholder('请选择月份').click();
    await this.page.locator('.ant-picker-dropdown:visible').waitFor();
    const cells = this.page.locator('.ant-picker-dropdown:visible .ant-picker-cell-inner');
    const count = await cells.count();
    if (count > 0) {
      await cells.nth(Math.min(1, count - 1)).click();
    }
    await this.page.waitForLoadState('networkidle');
  }

  /** 断言当前月份显示（如"2026年6月"） */
  async expectMonthText(text: string) {
    await expect(this.page.getByPlaceholder('请选择月份')).toContainText(text);
  }

  /** 断言金额输入校验错误 */
  async expectAmountError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('请输入正数金额');
  }
}
