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
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  /** 在弹窗网格选择器中点击分类（按文本匹配） */
  async pickCategory(name: string) {
    await this.page.locator('.category-pick-item', { hasText: name }).click();
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
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
  }

  /** 通过滤选栏选择分类 */
  async filterByCategory(name: string) {
    await this.page.getByTestId('filter-category').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 通过滤选栏选择创建者 */
  async filterByCreator(name: string) {
    await this.page.getByTestId('filter-creator').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 滚动到 sentinel 触发无限加载 */
  async scrollToLoadMore() {
    await this.page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
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

  /** 编辑第 n 条记录：点击记录 → 点击编辑按钮 */
  async editRecord(index: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await items.nth(index).click();
    await this.page.getByTestId('edit-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  /** 断言第 n 条记录显示指定的创建者名称 */
  async expectRecordCreator(index: number, name: string) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items.nth(index).getByTestId('creator-name')).toContainText(name);
  }

  /** 切换到上个月 */
  async goPrevMonth() {
    await this.page.getByTestId('month-prev').click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 切换到下个月 */
  async goNextMonth() {
    await this.page.getByTestId('month-next').click();
    await this.page.waitForLoadState('networkidle');
  }

  /** 断言当前月份显示（如"2026年6月"） */
  async expectMonthText(text: string) {
    await expect(this.page.getByTestId('current-month')).toContainText(text);
  }

  /** 断言金额输入校验错误 */
  async expectAmountError() {
    await expect(this.page.locator('.ant-form-item-explain')).toContainText('请输入正数金额');
  }
}
